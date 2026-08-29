package main

import (
	"math"
	"path/filepath"
	"testing"
)

func closeEnough(actual, expected, tolerance float64) bool {
	return math.Abs(actual-expected) <= tolerance
}

func TestPlandalfFSRS7ReferenceValues(t *testing.T) {
	p := defaultPlandalfFSRS7Parameters()
	state := plandalfMemoryState{StabilityDays: 4.1283, Difficulty: 4.194588083372719}
	retention := plandalfForgettingCurve(1.0, state.StabilityDays, p).Retention
	if !closeEnough(retention, 0.9242342483541028, 1e-12) {
		t.Fatalf("retrievability mismatch: got %.15f", retention)
	}

	initial := plandalfInitialMemoryState(plandalfGood, p)
	next, err := plandalfNextMemoryState(initial, 2.0, plandalfGood, p)
	if err != nil { t.Fatal(err) }
	if !closeEnough(next.StabilityDays, 10.362647327728341, 1e-10) {
		t.Fatalf("stability mismatch: got %.15f", next.StabilityDays)
	}
	if !closeEnough(next.Difficulty, 4.180821488255665, 1e-10) {
		t.Fatalf("difficulty mismatch: got %.15f", next.Difficulty)
	}

	interval, err := plandalfIntervalForRetention(4.1283, 0.9, p)
	if err != nil { t.Fatal(err) }
	if !closeEnough(interval, 2.966927162372141, 1e-10) {
		t.Fatalf("interval mismatch: got %.15f", interval)
	}
}

func TestPlandalfStoreRecordsImmutableReviewAndState(t *testing.T) {
	store, err := OpenPlandalfStore(DatabaseConfig{
		Mode: DatabaseModeSQLite,
		Path: filepath.Join(t.TempDir(), "plandalf.db"),
	})
	if err != nil { t.Fatal(err) }
	defer store.Close()

	deckID, err := store.CreateDeck("MongoDB", 0)
	if err != nil { t.Fatal(err) }
	cardID, err := store.CreateCard(deckID, "What does $match do?", "Filters documents in an aggregation pipeline.", 0)
	if err != nil { t.Fatal(err) }

	schedule, count, err := store.Preview(cardID, 0)
	if err != nil { t.Fatal(err) }
	if count != 0 { t.Fatalf("expected new card, got %d reviews", count) }
	if schedule.Good.DueAtMs <= 0 { t.Fatal("expected positive due timestamp") }

	reviewID, candidate, state, err := store.RecordReview(cardID, plandalfGood, 0, 0)
	if err != nil { t.Fatal(err) }
	if reviewID <= 0 { t.Fatal("expected review id") }
	if candidate.DueAtMs != state.DueAtMs { t.Fatal("review candidate and state due times differ") }

	history, err := loadPlandalfHistory(store.db, cardID)
	if err != nil { t.Fatal(err) }
	if len(history) != 1 || history[0].Rating != plandalfGood {
		t.Fatalf("unexpected history: %#v", history)
	}
	if _, err := store.db.Exec(`UPDATE reviews SET rating=1 WHERE id=?`, reviewID); err == nil {
		t.Fatal("review history update should be rejected by immutable trigger")
	}
}

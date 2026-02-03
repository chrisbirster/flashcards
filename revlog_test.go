package main

import (
	"os"
	"testing"
	"time"

	"github.com/open-spaced-repetition/go-fsrs/v3"
)

// TestRevlogPersistence validates that revlog entries are correctly saved with card_id and time_taken_ms
func TestRevlogPersistence(t *testing.T) {
	dbPath := "./test_revlog_persistence.db"
	defer os.Remove(dbPath)

	// Initialize collection
	col, store, err := InitDefaultCollection(dbPath)
	if err != nil {
		t.Fatalf("Failed to initialize collection: %v", err)
	}
	defer store.Close()

	// Create a deck
	deck := col.NewDeck("Revlog Test Deck")
	err = store.CreateDeck(deck)
	if err != nil {
		t.Fatalf("Failed to create deck: %v", err)
	}

	// Add a note
	note, cards, err := col.AddNote(deck.ID, "Basic", map[string]string{
		"Front": "Test Question",
		"Back":  "Test Answer",
	}, time.Now())
	if err != nil {
		t.Fatalf("Failed to add note: %v", err)
	}

	// Persist note
	err = store.CreateNote("default", &note)
	if err != nil {
		t.Fatalf("Failed to persist note: %v", err)
	}

	// Persist card
	card := cards[0]
	err = store.CreateCard(card)
	if err != nil {
		t.Fatalf("Failed to persist card: %v", err)
	}

	cardID := card.ID
	now := time.Now()

	// Answer the card with Good rating and specific time
	timeTakenMs := 3500 // 3.5 seconds
	revlog, err := col.Answer(cardID, fsrs.Good, now, timeTakenMs)
	if err != nil {
		t.Fatalf("Failed to answer card: %v", err)
	}

	// Persist the revlog entry (simulating what server.go does)
	err = store.AddRevlog(revlog, cardID, timeTakenMs)
	if err != nil {
		t.Fatalf("Failed to persist revlog: %v", err)
	}

	// Query the database directly to verify the revlog entry
	query := `SELECT card_id, rating, time_taken_ms FROM revlog WHERE card_id = ?`
	var dbCardID int64
	var dbRating int
	var dbTimeTakenMs int

	err = store.db.QueryRow(query, cardID).Scan(&dbCardID, &dbRating, &dbTimeTakenMs)
	if err != nil {
		t.Fatalf("Failed to query revlog: %v", err)
	}

	// Validate card_id is correct (not 0)
	if dbCardID != cardID {
		t.Errorf("Expected revlog card_id to be %d, got %d", cardID, dbCardID)
	}

	// Validate rating is correct
	if dbRating != int(fsrs.Good) {
		t.Errorf("Expected revlog rating to be %d (Good), got %d", int(fsrs.Good), dbRating)
	}

	// Validate time_taken_ms is correct (not 0)
	if dbTimeTakenMs != timeTakenMs {
		t.Errorf("Expected revlog time_taken_ms to be %d, got %d", timeTakenMs, dbTimeTakenMs)
	}

	// Test with different ratings and times
	ratings := []fsrs.Rating{fsrs.Again, fsrs.Hard, fsrs.Easy}
	times := []int{1200, 2500, 4800}

	for i, rating := range ratings {
		revlog, err := col.Answer(cardID, rating, now.Add(time.Duration(i+1)*time.Minute), times[i])
		if err != nil {
			t.Fatalf("Failed to answer card with rating %v: %v", rating, err)
		}

		err = store.AddRevlog(revlog, cardID, times[i])
		if err != nil {
			t.Fatalf("Failed to persist revlog for rating %v: %v", rating, err)
		}
	}

	// Verify all entries exist with correct data
	countQuery := `SELECT COUNT(*) FROM revlog WHERE card_id = ? AND time_taken_ms > 0`
	var count int
	err = store.db.QueryRow(countQuery, cardID).Scan(&count)
	if err != nil {
		t.Fatalf("Failed to count revlog entries: %v", err)
	}

	expectedCount := 4 // 1 initial + 3 additional
	if count != expectedCount {
		t.Errorf("Expected %d revlog entries with non-zero time, got %d", expectedCount, count)
	}
}

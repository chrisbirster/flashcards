package main

import (
	"path/filepath"
	"testing"
	"time"
)

func setupStoreWithTempDB(t *testing.T) (*SQLiteStore, *Collection) {
	t.Helper()

	dbPath := filepath.Join(t.TempDir(), "storage-additional.db")
	store, err := NewSQLiteStore(dbPath)
	if err != nil {
		t.Fatalf("failed to create sqlite store: %v", err)
	}
	t.Cleanup(func() {
		_ = store.Close()
	})

	col := NewCollection()
	if err := store.CreateCollection(col); err != nil {
		t.Fatalf("failed to create collection: %v", err)
	}

	return store, col
}

func TestSQLiteStore_TransactionsAndCRUDBranches(t *testing.T) {
	store, _ := setupStoreWithTempDB(t)

	tx, err := store.BeginTx()
	if err != nil {
		t.Fatalf("BeginTx failed: %v", err)
	}
	if err := store.RollbackTx(tx); err != nil {
		t.Fatalf("RollbackTx failed: %v", err)
	}

	tx2, err := store.BeginTx()
	if err != nil {
		t.Fatalf("BeginTx (second) failed: %v", err)
	}
	if err := store.CommitTx(tx2); err != nil {
		t.Fatalf("CommitTx failed: %v", err)
	}

	now := time.Now()
	collection := NewCollection()
	collection.USN = 77
	collection.LastSync = now
	if err := store.UpdateCollection(collection); err != nil {
		t.Fatalf("UpdateCollection failed: %v", err)
	}

	loadedCollection, err := store.GetCollection("default")
	if err != nil {
		t.Fatalf("GetCollection failed: %v", err)
	}
	if loadedCollection.USN != 77 {
		t.Fatalf("expected collection USN 77, got %d", loadedCollection.USN)
	}

	deck1 := &Deck{ID: 1, Name: "Deck 1", Cards: []int64{}}
	deck2 := &Deck{ID: 2, Name: "Deck 2", Cards: []int64{}}
	if err := store.CreateDeck(deck1); err != nil {
		t.Fatalf("CreateDeck(deck1) failed: %v", err)
	}
	if err := store.CreateDeck(deck2); err != nil {
		t.Fatalf("CreateDeck(deck2) failed: %v", err)
	}

	newName := "Deck 1 Updated"
	deck1.Name = newName
	if err := store.UpdateDeck(deck1); err != nil {
		t.Fatalf("UpdateDeck failed: %v", err)
	}
	gotDeck1, err := store.GetDeck(deck1.ID)
	if err != nil {
		t.Fatalf("GetDeck(deck1) failed: %v", err)
	}
	if gotDeck1.Name != newName {
		t.Fatalf("expected updated deck name %q, got %q", newName, gotDeck1.Name)
	}

	if err := store.DeleteDeck(deck2.ID); err != nil {
		t.Fatalf("DeleteDeck(deck2) failed: %v", err)
	}
	if _, err := store.GetDeck(deck2.ID); err == nil {
		t.Fatalf("expected deleted deck %d to be missing", deck2.ID)
	}

	basic := NoteType{
		Name:           "Basic",
		Fields:         []string{"Front", "Back"},
		Templates:      []CardTemplate{{Name: "Card 1", QFmt: "{{Front}}", AFmt: "{{Back}}"}},
		SortFieldIndex: 0,
	}
	if err := store.CreateNoteType("default", &basic); err != nil {
		t.Fatalf("CreateNoteType failed: %v", err)
	}

	note := &Note{
		ID:         10,
		Type:       "Basic",
		FieldMap:   map[string]string{"Front": "Question", "Back": "Answer"},
		Tags:       []string{"tag-a"},
		USN:        1,
		CreatedAt:  now,
		ModifiedAt: now,
	}
	if err := store.CreateNote("default", note); err != nil {
		t.Fatalf("CreateNote failed: %v", err)
	}

	note.FieldMap["Front"] = "Question Updated"
	note.Tags = []string{"tag-b"}
	note.USN = 2
	note.ModifiedAt = now.Add(5 * time.Minute)
	if err := store.UpdateNote(note); err != nil {
		t.Fatalf("UpdateNote failed: %v", err)
	}
	gotNote, err := store.GetNote(note.ID)
	if err != nil {
		t.Fatalf("GetNote failed: %v", err)
	}
	if gotNote.FieldMap["Front"] != "Question Updated" {
		t.Fatalf("expected updated note Front field, got %q", gotNote.FieldMap["Front"])
	}

	card := &Card{
		ID:           20,
		NoteID:       note.ID,
		DeckID:       deck1.ID,
		TemplateName: "Card 1",
		Ordinal:      0,
		Front:        "Question Updated",
		Back:         "Answer",
		SRS:          newDueNow(now),
		USN:          1,
	}
	if err := store.CreateCard(card); err != nil {
		t.Fatalf("CreateCard failed: %v", err)
	}

	notesByType, err := store.GetNotesByType("default", "Basic")
	if err != nil {
		t.Fatalf("GetNotesByType failed: %v", err)
	}
	if len(notesByType) != 1 {
		t.Fatalf("expected 1 note by type, got %d", len(notesByType))
	}

	cardsByNote, err := store.GetCardsByNote(note.ID)
	if err != nil {
		t.Fatalf("GetCardsByNote failed: %v", err)
	}
	if len(cardsByNote) != 1 {
		t.Fatalf("expected 1 card by note, got %d", len(cardsByNote))
	}

	// Delete a different note so we don't break card FK expectations.
	secondNote := &Note{
		ID:         11,
		Type:       "Basic",
		FieldMap:   map[string]string{"Front": "Q2", "Back": "A2"},
		Tags:       []string{},
		USN:        1,
		CreatedAt:  now,
		ModifiedAt: now,
	}
	if err := store.CreateNote("default", secondNote); err != nil {
		t.Fatalf("CreateNote(second) failed: %v", err)
	}
	if err := store.DeleteNote(secondNote.ID); err != nil {
		t.Fatalf("DeleteNote(second) failed: %v", err)
	}
	if _, err := store.GetNote(secondNote.ID); err == nil {
		t.Fatalf("expected deleted note %d to be missing", secondNote.ID)
	}
}

func TestSQLiteStore_GetDueCardsForUserScopesReviewStateBackfillToDeck(t *testing.T) {
	store, _ := setupStoreWithTempDB(t)

	now := time.Now()
	user := &User{
		ID:          newID("usr"),
		Email:       "review-state-scope@example.com",
		DisplayName: "Review Scope",
		LastLoginAt: now,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := store.CreateUser(user); err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	deckOne := &Deck{ID: 1, Name: "Deck 1", Cards: []int64{}}
	deckTwo := &Deck{ID: 2, Name: "Deck 2", Cards: []int64{}}
	if err := store.CreateDeck(deckOne); err != nil {
		t.Fatalf("CreateDeck(deckOne) failed: %v", err)
	}
	if err := store.CreateDeck(deckTwo); err != nil {
		t.Fatalf("CreateDeck(deckTwo) failed: %v", err)
	}

	noteType := &NoteType{
		Name:   "Basic",
		Fields: []string{"Front", "Back"},
		Templates: []CardTemplate{
			{Name: "Card 1", QFmt: "{{Front}}", AFmt: "{{Back}}"},
		},
	}
	if err := store.CreateNoteType("default", noteType); err != nil {
		t.Fatalf("CreateNoteType failed: %v", err)
	}

	for i := 1; i <= 3; i++ {
		note := &Note{
			ID:         int64(i),
			Type:       "Basic",
			FieldMap:   map[string]string{"Front": "Q", "Back": "A"},
			Tags:       []string{},
			USN:        int64(i),
			CreatedAt:  now,
			ModifiedAt: now,
		}
		if err := store.CreateNote("default", note); err != nil {
			t.Fatalf("CreateNote(%d) failed: %v", i, err)
		}
	}

	cards := []*Card{
		{ID: 1, NoteID: 1, DeckID: deckOne.ID, TemplateName: "Card 1", Front: "Q1", Back: "A1", SRS: newDueNow(now.Add(-time.Hour)), USN: 1},
		{ID: 2, NoteID: 2, DeckID: deckOne.ID, TemplateName: "Card 1", Front: "Q2", Back: "A2", SRS: newDueNow(now.Add(-time.Hour)), USN: 1},
		{ID: 3, NoteID: 3, DeckID: deckTwo.ID, TemplateName: "Card 1", Front: "Q3", Back: "A3", SRS: newDueNow(now.Add(-time.Hour)), USN: 1},
	}
	for _, card := range cards {
		if err := store.CreateCard(card); err != nil {
			t.Fatalf("CreateCard(%d) failed: %v", card.ID, err)
		}
	}

	dueCards, err := store.GetDueCardsForUser(user.ID, deckOne.ID, 10)
	if err != nil {
		t.Fatalf("GetDueCardsForUser failed: %v", err)
	}
	if len(dueCards) != 2 {
		t.Fatalf("expected 2 due cards for deck one, got %d", len(dueCards))
	}

	var reviewStateCount int
	if err := store.db.QueryRow(`SELECT COUNT(*) FROM card_review_states WHERE user_id = ?`, user.ID).Scan(&reviewStateCount); err != nil {
		t.Fatalf("failed to count scoped review states: %v", err)
	}
	if reviewStateCount != 2 {
		t.Fatalf("expected 2 review states after deck-scoped due load, got %d", reviewStateCount)
	}

	var otherDeckStateCount int
	if err := store.db.QueryRow(`SELECT COUNT(*) FROM card_review_states WHERE user_id = ? AND card_id = ?`, user.ID, int64(3)).Scan(&otherDeckStateCount); err != nil {
		t.Fatalf("failed to count other deck review states: %v", err)
	}
	if otherDeckStateCount != 0 {
		t.Fatalf("expected deck-scoped due load to skip other decks, got %d rows for deck two", otherDeckStateCount)
	}
}

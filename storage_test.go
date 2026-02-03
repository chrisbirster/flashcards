package main

import (
	"os"
	"testing"
	"time"

	fsrs "github.com/open-spaced-repetition/go-fsrs/v3"
)

func setupTestDB(t *testing.T) (*SQLiteStore, func()) {
	// Create temporary database
	dbPath := "./test_microdote.db"
	store, err := NewSQLiteStore(dbPath)
	if err != nil {
		t.Fatalf("Failed to create test store: %v", err)
	}

	cleanup := func() {
		store.Close()
		os.Remove(dbPath)
	}

	return store, cleanup
}

func TestCreateAndGetDeck(t *testing.T) {
	store, cleanup := setupTestDB(t)
	defer cleanup()

	// Create collection first (required for foreign key)
	col := NewCollection()
	if err := store.CreateCollection(col); err != nil {
		t.Fatalf("Failed to create collection: %v", err)
	}

	// Create a deck
	deck := &Deck{
		ID:    1,
		Name:  "Test Deck",
		Cards: []int64{},
	}

	err := store.CreateDeck(deck)
	if err != nil {
		t.Fatalf("Failed to create deck: %v", err)
	}

	// Retrieve the deck
	retrieved, err := store.GetDeck(1)
	if err != nil {
		t.Fatalf("Failed to get deck: %v", err)
	}

	if retrieved.Name != "Test Deck" {
		t.Errorf("Expected deck name 'Test Deck', got '%s'", retrieved.Name)
	}
}

func TestListDecks(t *testing.T) {
	store, cleanup := setupTestDB(t)
	defer cleanup()

	// Create collection
	col := NewCollection()
	store.CreateCollection(col)

	// Create multiple decks
	decks := []*Deck{
		{ID: 1, Name: "Deck A", Cards: []int64{}},
		{ID: 2, Name: "Deck B", Cards: []int64{}},
		{ID: 3, Name: "Deck C", Cards: []int64{}},
	}

	for _, d := range decks {
		if err := store.CreateDeck(d); err != nil {
			t.Fatalf("Failed to create deck: %v", err)
		}
	}

	// List decks
	retrieved, err := store.ListDecks("default")
	if err != nil {
		t.Fatalf("Failed to list decks: %v", err)
	}

	if len(retrieved) != 3 {
		t.Errorf("Expected 3 decks, got %d", len(retrieved))
	}
}

func TestCreateAndGetNote(t *testing.T) {
	store, cleanup := setupTestDB(t)
	defer cleanup()

	// Create collection
	col := NewCollection()
	store.CreateCollection(col)

	// Create note type first
	nt := &NoteType{
		Name:   "Basic",
		Fields: []string{"Front", "Back"},
		Templates: []CardTemplate{
			{Name: "Card 1", QFmt: "{{Front}}", AFmt: "{{Back}}"},
		},
	}
	store.CreateNoteType("default", nt)

	// Create a note
	note := &Note{
		ID:         1,
		Type:       "Basic",
		FieldMap:   map[string]string{"Front": "Question", "Back": "Answer"},
		Tags:       []string{"test", "vocab"},
		USN:        1,
		CreatedAt:  time.Now(),
		ModifiedAt: time.Now(),
	}

	err := store.CreateNote("default", note)
	if err != nil {
		t.Fatalf("Failed to create note: %v", err)
	}

	// Retrieve the note
	retrieved, err := store.GetNote(1)
	if err != nil {
		t.Fatalf("Failed to get note: %v", err)
	}

	if retrieved.FieldMap["Front"] != "Question" {
		t.Errorf("Expected Front='Question', got '%s'", retrieved.FieldMap["Front"])
	}

	if len(retrieved.Tags) != 2 {
		t.Errorf("Expected 2 tags, got %d", len(retrieved.Tags))
	}
}

func TestCreateAndGetCard(t *testing.T) {
	store, cleanup := setupTestDB(t)
	defer cleanup()

	// Create collection, deck, and note
	col := NewCollection()
	store.CreateCollection(col)
	store.CreateDeck(&Deck{ID: 1, Name: "Test", Cards: []int64{}})

	nt := &NoteType{Name: "Basic", Fields: []string{"Front", "Back"}, Templates: []CardTemplate{{Name: "Card 1", QFmt: "{{Front}}", AFmt: "{{Back}}"}}}
	store.CreateNoteType("default", nt)
	store.CreateNote("default", &Note{ID: 1, Type: "Basic", FieldMap: map[string]string{"Front": "Q"}, Tags: []string{}, USN: 1, CreatedAt: time.Now(), ModifiedAt: time.Now()})

	// Create a card
	card := &Card{
		ID:           1,
		NoteID:       1,
		DeckID:       1,
		TemplateName: "Card 1",
		Ordinal:      0,
		Front:        "Question",
		Back:         "Answer",
		SRS:          newDueNow(time.Now()),
		Flag:         0,
		Marked:       false,
		Suspended:    false,
		USN:          1,
	}

	err := store.CreateCard(card)
	if err != nil {
		t.Fatalf("Failed to create card: %v", err)
	}

	// Retrieve the card
	retrieved, err := store.GetCard(1)
	if err != nil {
		t.Fatalf("Failed to get card: %v", err)
	}

	if retrieved.Front != "Question" {
		t.Errorf("Expected Front='Question', got '%s'", retrieved.Front)
	}

	if retrieved.Marked {
		t.Error("Expected card to not be marked")
	}
}

func TestGetDueCards(t *testing.T) {
	store, cleanup := setupTestDB(t)
	defer cleanup()

	// Setup collection, deck, note type, and notes
	col := NewCollection()
	store.CreateCollection(col)
	store.CreateDeck(&Deck{ID: 1, Name: "Test", Cards: []int64{}})

	nt := &NoteType{Name: "Basic", Fields: []string{"Front", "Back"}, Templates: []CardTemplate{{Name: "Card 1", QFmt: "{{Front}}", AFmt: "{{Back}}"}}}
	store.CreateNoteType("default", nt)

	for i := 1; i <= 3; i++ {
		store.CreateNote("default", &Note{ID: int64(i), Type: "Basic", FieldMap: map[string]string{"Front": "Q"}, Tags: []string{}, USN: 1, CreatedAt: time.Now(), ModifiedAt: time.Now()})
	}

	now := time.Now()
	pastDue := now.Add(-1 * time.Hour)

	// Create cards with different due dates
	cards := []*Card{
		{ID: 1, NoteID: 1, DeckID: 1, TemplateName: "T1", Front: "Q1", Back: "A1", SRS: newDueNow(pastDue), USN: 1},
		{ID: 2, NoteID: 2, DeckID: 1, TemplateName: "T1", Front: "Q2", Back: "A2", SRS: newDueNow(pastDue), USN: 1},
		{ID: 3, NoteID: 3, DeckID: 1, TemplateName: "T1", Front: "Q3", Back: "A3", SRS: newDueNow(now.Add(24 * time.Hour)), USN: 1},
	}

	for _, c := range cards {
		if err := store.CreateCard(c); err != nil {
			t.Fatalf("Failed to create card: %v", err)
		}
	}

	// Get due cards
	due, err := store.GetDueCards(1, 10)
	if err != nil {
		t.Fatalf("Failed to get due cards: %v", err)
	}

	if len(due) != 2 {
		t.Errorf("Expected 2 due cards, got %d", len(due))
	}
}

func TestProfileCRUD(t *testing.T) {
	store, cleanup := setupTestDB(t)
	defer cleanup()

	// Create collection first
	col := NewCollection()
	store.CreateCollection(col)

	// Create profile
	profile := &Profile{
		ID:           "test-profile",
		Name:         "Test User",
		CollectionID: "default",
		SyncAccount:  "test@example.com",
		CreatedAt:    time.Now(),
	}

	err := store.CreateProfile(profile)
	if err != nil {
		t.Fatalf("Failed to create profile: %v", err)
	}

	// Get profile
	retrieved, err := store.GetProfile("test-profile")
	if err != nil {
		t.Fatalf("Failed to get profile: %v", err)
	}

	if retrieved.Name != "Test User" {
		t.Errorf("Expected name 'Test User', got '%s'", retrieved.Name)
	}

	// List profiles
	profiles, err := store.ListProfiles()
	if err != nil {
		t.Fatalf("Failed to list profiles: %v", err)
	}

	if len(profiles) != 1 {
		t.Errorf("Expected 1 profile, got %d", len(profiles))
	}

	// Set and get active profile
	err = store.SetActiveProfile("test-profile")
	if err != nil {
		t.Fatalf("Failed to set active profile: %v", err)
	}

	active, err := store.GetActiveProfile()
	if err != nil {
		t.Fatalf("Failed to get active profile: %v", err)
	}

	if active.ID != "test-profile" {
		t.Errorf("Expected active profile 'test-profile', got '%s'", active.ID)
	}
}

func TestCollectionIDCounters(t *testing.T) {
	store, cleanup := setupTestDB(t)
	defer cleanup()

	// Create collection with some data
	col := NewCollection()
	store.CreateCollection(col)

	// Create decks with specific IDs
	store.CreateDeck(&Deck{ID: 1, Name: "Deck 1", Cards: []int64{}})
	store.CreateDeck(&Deck{ID: 5, Name: "Deck 5", Cards: []int64{}})
	store.CreateDeck(&Deck{ID: 3, Name: "Deck 3", Cards: []int64{}})

	// Load collection
	loaded, err := store.GetCollection("default")
	if err != nil {
		t.Fatalf("Failed to load collection: %v", err)
	}

	// Next deck ID should be max(existing IDs) + 1 = 6
	if loaded.nextDeckID != 6 {
		t.Errorf("Expected nextDeckID=6, got %d", loaded.nextDeckID)
	}

	// Test creating a new deck
	newDeck := loaded.NewDeck("New Deck")
	if newDeck.ID != 6 {
		t.Errorf("Expected new deck ID=6, got %d", newDeck.ID)
	}
}

func TestBackupCreation(t *testing.T) {
	store, cleanup := setupTestDB(t)
	defer cleanup()

	// Create collection
	col := NewCollection()
	store.CreateCollection(col)

	// Create backup manager
	backupDir := "./test_backups"
	os.MkdirAll(backupDir, 0755)
	defer os.RemoveAll(backupDir)

	bm := NewBackupManager("./test_microdote.db", backupDir, store)

	// Create backup
	backupPath, err := bm.CreateBackup("default")
	if err != nil {
		t.Fatalf("Failed to create backup: %v", err)
	}

	// Verify backup file exists
	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		t.Errorf("Backup file was not created: %s", backupPath)
	}
}

// M1 Tests - Studying MVP

func TestGetDeckStats(t *testing.T) {
	store, cleanup := setupTestDB(t)
	defer cleanup()

	// Setup
	col := NewCollection()
	store.CreateCollection(col)
	store.CreateDeck(&Deck{ID: 1, Name: "Test Deck", Cards: []int64{}})

	nt := &NoteType{Name: "Basic", Fields: []string{"Front", "Back"}, Templates: []CardTemplate{{Name: "Card 1", QFmt: "{{Front}}", AFmt: "{{Back}}"}}}
	store.CreateNoteType("default", nt)

	// Create cards with different states
	now := time.Now()
	pastDue := now.Add(-1 * time.Hour)
	futureDue := now.Add(24 * time.Hour)

	// Create notes
	for i := 1; i <= 5; i++ {
		store.CreateNote("default", &Note{
			ID:         int64(i),
			Type:       "Basic",
			FieldMap:   map[string]string{"Front": "Q", "Back": "A"},
			Tags:       []string{},
			USN:        1,
			CreatedAt:  now,
			ModifiedAt: now,
		})
	}

	// Create cards with various states
	cards := []*Card{
		{ID: 1, NoteID: 1, DeckID: 1, TemplateName: "Card 1", Front: "Q1", Back: "A1", SRS: newDueNow(pastDue), USN: 1},                  // New, due
		{ID: 2, NoteID: 2, DeckID: 1, TemplateName: "Card 1", Front: "Q2", Back: "A2", SRS: newDueNow(futureDue), USN: 1},                // New, not due
		{ID: 3, NoteID: 3, DeckID: 1, TemplateName: "Card 1", Front: "Q3", Back: "A3", SRS: newDueNow(pastDue), USN: 1},                  // Review, due
		{ID: 4, NoteID: 4, DeckID: 1, TemplateName: "Card 1", Front: "Q4", Back: "A4", SRS: newDueNow(pastDue), Suspended: true, USN: 1}, // Suspended
		{ID: 5, NoteID: 5, DeckID: 1, TemplateName: "Card 1", Front: "Q5", Back: "A5", SRS: newDueNow(futureDue), USN: 1},                // Not due
	}

	// Set card 3 to Review state
	cards[2].SRS.State = 2 // Review state

	for _, c := range cards {
		if err := store.CreateCard(c); err != nil {
			t.Fatalf("Failed to create card: %v", err)
		}
	}

	// Get deck stats
	stats, err := store.GetDeckStats(1)
	if err != nil {
		t.Fatalf("Failed to get deck stats: %v", err)
	}

	// Verify stats
	if stats.TotalCards != 5 {
		t.Errorf("Expected 5 total cards, got %d", stats.TotalCards)
	}

	if stats.Suspended != 1 {
		t.Errorf("Expected 1 suspended card, got %d", stats.Suspended)
	}

	// Cards 1, 2, 5 are New (card 3 is Review, card 4 is suspended)
	if stats.NewCards != 3 {
		t.Errorf("Expected 3 new cards, got %d", stats.NewCards)
	}

	if stats.Review != 1 {
		t.Errorf("Expected 1 review card, got %d", stats.Review)
	}

	// Cards 1 and 3 are due (cards 2 and 5 are future, card 4 is suspended)
	if stats.DueToday != 2 {
		t.Errorf("Expected 2 cards due today, got %d", stats.DueToday)
	}
}

func TestStudyFlow(t *testing.T) {
	store, cleanup := setupTestDB(t)
	defer cleanup()

	// Setup collection
	col := NewCollection()
	store.CreateCollection(col)

	// Add built-in note types
	noteTypes := builtins()
	for _, nt := range noteTypes {
		store.CreateNoteType("default", &nt)
	}

	// Create a deck
	deck := &Deck{ID: 1, Name: "Study Deck", Cards: []int64{}}
	store.CreateDeck(deck)

	// Create a note
	note := &Note{
		ID:         1,
		Type:       "Basic",
		FieldMap:   map[string]string{"Front": "Capital of France?", "Back": "Paris"},
		Tags:       []string{"geography"},
		USN:        1,
		CreatedAt:  time.Now(),
		ModifiedAt: time.Now(),
	}
	store.CreateNote("default", note)

	// Create a card
	now := time.Now()
	card := &Card{
		ID:           1,
		NoteID:       1,
		DeckID:       1,
		TemplateName: "Card 1",
		Ordinal:      0,
		Front:        "Q: Capital of France?",
		Back:         "A: Paris",
		SRS:          newDueNow(now),
		Flag:         0,
		Marked:       false,
		Suspended:    false,
		USN:          1,
	}
	store.CreateCard(card)

	// Get due cards
	dueCards, err := store.GetDueCards(1, 10)
	if err != nil {
		t.Fatalf("Failed to get due cards: %v", err)
	}

	if len(dueCards) != 1 {
		t.Fatalf("Expected 1 due card, got %d", len(dueCards))
	}

	// Answer the card with "Good" (rating 3)
	col.Cards = make(map[int64]*Card)
	col.Cards[1] = dueCards[0]

	_, err = col.Answer(1, 3, now, 5000) // Good rating, 5 seconds
	if err != nil {
		t.Fatalf("Failed to answer card: %v", err)
	}

	// Update card in database
	updatedCard := col.Cards[1]
	err = store.UpdateCard(updatedCard)
	if err != nil {
		t.Fatalf("Failed to update card: %v", err)
	}

	// Verify card was rescheduled (should have future due date)
	retrievedCard, err := store.GetCard(1)
	if err != nil {
		t.Fatalf("Failed to retrieve card: %v", err)
	}

	if !retrievedCard.SRS.Due.After(now) {
		t.Error("Expected card to be rescheduled to future date")
	}

	// Card should no longer be due
	dueCardsAfter, err := store.GetDueCards(1, 10)
	if err != nil {
		t.Fatalf("Failed to get due cards after answer: %v", err)
	}

	if len(dueCardsAfter) != 0 {
		t.Errorf("Expected 0 due cards after answering, got %d", len(dueCardsAfter))
	}
}

func TestAnswerCardMultipleRatings(t *testing.T) {
	store, cleanup := setupTestDB(t)
	defer cleanup()

	// Setup
	col := NewCollection()
	store.CreateCollection(col)
	store.CreateDeck(&Deck{ID: 1, Name: "Test", Cards: []int64{}})

	nt := &NoteType{Name: "Basic", Fields: []string{"Front", "Back"}, Templates: []CardTemplate{{Name: "Card 1", QFmt: "{{Front}}", AFmt: "{{Back}}"}}}
	store.CreateNoteType("default", nt)

	now := time.Now()

	// Test each rating: Again (1), Hard (2), Good (3), Easy (4)
	ratings := []int{1, 2, 3, 4}

	for i, rating := range ratings {
		noteID := int64(i + 1)
		cardID := int64(i + 1)

		// Create note
		store.CreateNote("default", &Note{
			ID:         noteID,
			Type:       "Basic",
			FieldMap:   map[string]string{"Front": "Q", "Back": "A"},
			Tags:       []string{},
			USN:        1,
			CreatedAt:  now,
			ModifiedAt: now,
		})

		// Create card
		card := &Card{
			ID:           cardID,
			NoteID:       noteID,
			DeckID:       1,
			TemplateName: "Card 1",
			Front:        "Q",
			Back:         "A",
			SRS:          newDueNow(now),
			USN:          1,
		}
		store.CreateCard(card)

		// Answer with different ratings
		col.Cards = make(map[int64]*Card)
		col.Cards[cardID] = card

		_, err := col.Answer(cardID, fsrs.Rating(rating), now, rating*1000) // Different timing per rating
		if err != nil {
			t.Fatalf("Failed to answer card with rating %d: %v", rating, err)
		}

		// Verify card state changed
		updatedCard := col.Cards[cardID]
		if updatedCard.SRS.Reps != 1 {
			t.Errorf("Expected card reps to be 1 after rating %d, got %d", rating, updatedCard.SRS.Reps)
		}
	}
}

func TestEmptyDeckStudy(t *testing.T) {
	store, cleanup := setupTestDB(t)
	defer cleanup()

	// Create collection and empty deck
	col := NewCollection()
	store.CreateCollection(col)
	store.CreateDeck(&Deck{ID: 1, Name: "Empty Deck", Cards: []int64{}})

	// Get due cards from empty deck
	dueCards, err := store.GetDueCards(1, 10)
	if err != nil {
		t.Fatalf("Failed to get due cards: %v", err)
	}

	if len(dueCards) != 0 {
		t.Errorf("Expected 0 due cards from empty deck, got %d", len(dueCards))
	}

	// Get stats for empty deck
	stats, err := store.GetDeckStats(1)
	if err != nil {
		t.Fatalf("Failed to get stats for empty deck: %v", err)
	}

	if stats.TotalCards != 0 {
		t.Errorf("Expected 0 total cards, got %d", stats.TotalCards)
	}

	if stats.DueToday != 0 {
		t.Errorf("Expected 0 due today, got %d", stats.DueToday)
	}
}

// M2 Tests

func TestAddNoteWithTags(t *testing.T) {
	dbPath := "./test_add_note.db"
	defer os.Remove(dbPath)

	// Initialize collection with note types
	col, store, err := InitDefaultCollection(dbPath)
	if err != nil {
		t.Fatalf("Failed to initialize collection: %v", err)
	}
	defer store.Close()

	// Create a deck
	deck := col.NewDeck("Test Deck")
	store.CreateDeck(deck)

	// Add a note with tags using Basic note type
	note, cards, err := col.AddNote(deck.ID, "Basic", map[string]string{
		"Front": "What is the capital of France?",
		"Back":  "Paris",
	}, time.Now())

	if err != nil {
		t.Fatalf("Failed to add note: %v", err)
	}

	// Set tags
	note.Tags = []string{"geography", "europe", "capitals"}

	// Persist note
	err = store.CreateNote("default", &note)
	if err != nil {
		t.Fatalf("Failed to persist note: %v", err)
	}

	// Persist cards
	for _, card := range cards {
		err = store.CreateCard(card)
		if err != nil {
			t.Fatalf("Failed to persist card: %v", err)
		}
	}

	// Verify note was created with correct fields
	retrievedNote, err := store.GetNote(note.ID)
	if err != nil {
		t.Fatalf("Failed to retrieve note: %v", err)
	}

	if retrievedNote.FieldMap["Front"] != "What is the capital of France?" {
		t.Errorf("Expected Front field 'What is the capital of France?', got '%s'", retrievedNote.FieldMap["Front"])
	}

	if len(retrievedNote.Tags) != 3 {
		t.Errorf("Expected 3 tags, got %d", len(retrievedNote.Tags))
	}

	// Verify card was created
	if len(cards) != 1 {
		t.Errorf("Expected 1 card for Basic note type, got %d", len(cards))
	}
}

func TestNoteTypesExist(t *testing.T) {
	dbPath := "./test_note_types.db"
	defer os.Remove(dbPath)

	// Initialize collection with built-in note types
	col, store, err := InitDefaultCollection(dbPath)
	if err != nil {
		t.Fatalf("Failed to initialize collection: %v", err)
	}
	defer store.Close()

	// Verify built-in note types exist
	expectedTypes := []string{"Basic", "Basic (and reversed card)", "Basic (optional reversed card)", "Basic (type in the answer)", "Cloze"}

	for _, typeName := range expectedTypes {
		if _, ok := col.NoteTypes[NoteTypeName(typeName)]; !ok {
			t.Errorf("Expected note type '%s' to exist", typeName)
		}
	}

	// Verify Basic has correct fields
	basic := col.NoteTypes["Basic"]
	if len(basic.Fields) != 2 {
		t.Errorf("Expected Basic to have 2 fields, got %d", len(basic.Fields))
	}
	if basic.Fields[0] != "Front" || basic.Fields[1] != "Back" {
		t.Errorf("Expected Basic fields to be [Front, Back], got %v", basic.Fields)
	}
}

func TestClozeNoteCreation(t *testing.T) {
	dbPath := "./test_cloze.db"
	defer os.Remove(dbPath)

	// Initialize collection
	col, store, err := InitDefaultCollection(dbPath)
	if err != nil {
		t.Fatalf("Failed to initialize collection: %v", err)
	}
	defer store.Close()

	// Create a deck
	deck := col.NewDeck("Cloze Test Deck")
	store.CreateDeck(deck)

	// Add a cloze note - {{c1::}} creates cloze deletions
	note, cards, err := col.AddNote(deck.ID, "Cloze", map[string]string{
		"Text":  "The capital of {{c1::France}} is {{c2::Paris}}",
		"Extra": "European geography",
	}, time.Now())

	if err != nil {
		t.Fatalf("Failed to add cloze note: %v", err)
	}

	// Cloze with c1 and c2 should create 2 cards
	if len(cards) != 2 {
		t.Errorf("Expected 2 cards for cloze with c1 and c2, got %d", len(cards))
	}

	// Verify note
	if note.Type != "Cloze" {
		t.Errorf("Expected note type 'Cloze', got '%s'", note.Type)
	}
}

// Task 0202: Duplicate Check Tests

func TestFindDuplicateNotes(t *testing.T) {
	dbPath := "./test_duplicate.db"
	defer os.Remove(dbPath)

	col, store, err := InitDefaultCollection(dbPath)
	if err != nil {
		t.Fatalf("Failed to initialize collection: %v", err)
	}
	defer store.Close()

	// Create a deck
	deck := col.NewDeck("Test Deck")
	store.CreateDeck(deck)

	// Create first note with Front="Hello World"
	note1, cards1, err := col.AddNote(deck.ID, "Basic", map[string]string{
		"Front": "Hello World",
		"Back":  "Test answer",
	}, time.Now())
	if err != nil {
		t.Fatalf("Failed to add note: %v", err)
	}
	store.CreateNote("default", &note1)
	for _, card := range cards1 {
		store.CreateCard(card)
	}

	// Create second note with different content
	note2, cards2, err := col.AddNote(deck.ID, "Basic", map[string]string{
		"Front": "Different Question",
		"Back":  "Different answer",
	}, time.Now())
	if err != nil {
		t.Fatalf("Failed to add note: %v", err)
	}
	store.CreateNote("default", &note2)
	for _, card := range cards2 {
		store.CreateCard(card)
	}

	// Test 1: Check for exact duplicate - should find one
	duplicates, err := store.FindDuplicateNotes("default", "Front", "Hello World", 0)
	if err != nil {
		t.Fatalf("Failed to find duplicates: %v", err)
	}
	if len(duplicates) != 1 {
		t.Errorf("Expected 1 duplicate for 'Hello World', got %d", len(duplicates))
	}

	// Test 2: Check for case-insensitive duplicate
	duplicates, err = store.FindDuplicateNotes("default", "Front", "hello world", 0)
	if err != nil {
		t.Fatalf("Failed to find duplicates: %v", err)
	}
	if len(duplicates) != 1 {
		t.Errorf("Expected 1 duplicate for case-insensitive 'hello world', got %d", len(duplicates))
	}

	// Test 3: Check for non-existent content - should find none
	duplicates, err = store.FindDuplicateNotes("default", "Front", "Non-existent question", 0)
	if err != nil {
		t.Fatalf("Failed to find duplicates: %v", err)
	}
	if len(duplicates) != 0 {
		t.Errorf("Expected 0 duplicates for non-existent content, got %d", len(duplicates))
	}

	// Test 4: Check with whitespace normalization
	duplicates, err = store.FindDuplicateNotes("default", "Front", "  Hello World  ", 0)
	if err != nil {
		t.Fatalf("Failed to find duplicates: %v", err)
	}
	if len(duplicates) != 1 {
		t.Errorf("Expected 1 duplicate with whitespace normalization, got %d", len(duplicates))
	}
}

func TestFindDuplicateNotesWithDeckFilter(t *testing.T) {
	dbPath := "./test_duplicate_deck.db"
	defer os.Remove(dbPath)

	col, store, err := InitDefaultCollection(dbPath)
	if err != nil {
		t.Fatalf("Failed to initialize collection: %v", err)
	}
	defer store.Close()

	// Create two decks
	deck1 := col.NewDeck("Deck 1")
	deck2 := col.NewDeck("Deck 2")
	store.CreateDeck(deck1)
	store.CreateDeck(deck2)

	// Create note in deck1
	note1, cards1, err := col.AddNote(deck1.ID, "Basic", map[string]string{
		"Front": "Same Question",
		"Back":  "Answer 1",
	}, time.Now())
	if err != nil {
		t.Fatalf("Failed to add note: %v", err)
	}
	store.CreateNote("default", &note1)
	for _, card := range cards1 {
		store.CreateCard(card)
	}

	// Create note in deck2
	note2, cards2, err := col.AddNote(deck2.ID, "Basic", map[string]string{
		"Front": "Same Question",
		"Back":  "Answer 2",
	}, time.Now())
	if err != nil {
		t.Fatalf("Failed to add note: %v", err)
	}
	store.CreateNote("default", &note2)
	for _, card := range cards2 {
		store.CreateCard(card)
	}

	// Test 1: Without deck filter - should find 2
	duplicates, err := store.FindDuplicateNotes("default", "Front", "Same Question", 0)
	if err != nil {
		t.Fatalf("Failed to find duplicates: %v", err)
	}
	if len(duplicates) != 2 {
		t.Errorf("Expected 2 duplicates without deck filter, got %d", len(duplicates))
	}

	// Test 2: With deck1 filter - should find 1
	duplicates, err = store.FindDuplicateNotes("default", "Front", "Same Question", deck1.ID)
	if err != nil {
		t.Fatalf("Failed to find duplicates: %v", err)
	}
	if len(duplicates) != 1 {
		t.Errorf("Expected 1 duplicate in deck1, got %d", len(duplicates))
	}

	// Test 3: With deck2 filter - should find 1
	duplicates, err = store.FindDuplicateNotes("default", "Front", "Same Question", deck2.ID)
	if err != nil {
		t.Fatalf("Failed to find duplicates: %v", err)
	}
	if len(duplicates) != 1 {
		t.Errorf("Expected 1 duplicate in deck2, got %d", len(duplicates))
	}
}

func TestFindDuplicateNotesEmptyValue(t *testing.T) {
	store, cleanup := setupTestDB(t)
	defer cleanup()

	col := NewCollection()
	store.CreateCollection(col)

	// Empty value should return empty result
	duplicates, err := store.FindDuplicateNotes("default", "Front", "", 0)
	if err != nil {
		t.Fatalf("Failed to check empty value: %v", err)
	}
	if len(duplicates) != 0 {
		t.Errorf("Expected 0 duplicates for empty value, got %d", len(duplicates))
	}

	// Whitespace-only value should return empty result (normalized to empty)
	duplicates, err = store.FindDuplicateNotes("default", "Front", "   ", 0)
	if err != nil {
		t.Fatalf("Failed to check whitespace value: %v", err)
	}
	if len(duplicates) != 0 {
		t.Errorf("Expected 0 duplicates for whitespace-only value, got %d", len(duplicates))
	}
}

// Task 0212: Field Editor Tests

func TestUpdateNoteType(t *testing.T) {
	store, cleanup := setupTestDB(t)
	defer cleanup()

	col := NewCollection()
	store.CreateCollection(col)

	// Create a note type
	nt := &NoteType{
		Name:   "TestType",
		Fields: []string{"Front", "Back"},
		Templates: []CardTemplate{{
			Name: "Card 1",
			QFmt: "{{Front}}",
			AFmt: "{{Back}}",
		}},
	}
	store.CreateNoteType("default", nt)

	// Update the note type - add a field
	nt.Fields = []string{"Front", "Back", "Extra"}
	err := store.UpdateNoteType("default", nt)
	if err != nil {
		t.Fatalf("Failed to update note type: %v", err)
	}

	// Verify the update
	retrieved, err := store.GetNoteType("default", "TestType")
	if err != nil {
		t.Fatalf("Failed to get note type: %v", err)
	}

	if len(retrieved.Fields) != 3 {
		t.Errorf("Expected 3 fields, got %d", len(retrieved.Fields))
	}

	if retrieved.Fields[2] != "Extra" {
		t.Errorf("Expected third field to be 'Extra', got '%s'", retrieved.Fields[2])
	}
}

func TestFieldReorder(t *testing.T) {
	store, cleanup := setupTestDB(t)
	defer cleanup()

	col := NewCollection()
	store.CreateCollection(col)

	// Create a note type
	nt := &NoteType{
		Name:   "ReorderTest",
		Fields: []string{"Field1", "Field2", "Field3"},
		Templates: []CardTemplate{{
			Name: "Card 1",
			QFmt: "{{Field1}}",
			AFmt: "{{Field2}}",
		}},
	}
	store.CreateNoteType("default", nt)

	// Reorder fields
	nt.Fields = []string{"Field3", "Field1", "Field2"}
	err := store.UpdateNoteType("default", nt)
	if err != nil {
		t.Fatalf("Failed to reorder fields: %v", err)
	}

	// Verify the reorder
	retrieved, err := store.GetNoteType("default", "ReorderTest")
	if err != nil {
		t.Fatalf("Failed to get note type: %v", err)
	}

	expected := []string{"Field3", "Field1", "Field2"}
	for i, f := range expected {
		if retrieved.Fields[i] != f {
			t.Errorf("Expected field %d to be '%s', got '%s'", i, f, retrieved.Fields[i])
		}
	}
}

// Task 0203: Flags/Marked Tests

func TestUpdateCardFlag(t *testing.T) {
	store, cleanup := setupTestDB(t)
	defer cleanup()

	// Setup
	col := NewCollection()
	store.CreateCollection(col)
	store.CreateDeck(&Deck{ID: 1, Name: "Test", Cards: []int64{}})

	nt := &NoteType{Name: "Basic", Fields: []string{"Front", "Back"}, Templates: []CardTemplate{{Name: "Card 1", QFmt: "{{Front}}", AFmt: "{{Back}}"}}}
	store.CreateNoteType("default", nt)

	now := time.Now()
	store.CreateNote("default", &Note{ID: 1, Type: "Basic", FieldMap: map[string]string{"Front": "Q"}, Tags: []string{}, USN: 1, CreatedAt: now, ModifiedAt: now})

	// Create card with flag 0
	card := &Card{
		ID:           1,
		NoteID:       1,
		DeckID:       1,
		TemplateName: "Card 1",
		Front:        "Q",
		Back:         "A",
		SRS:          newDueNow(now),
		Flag:         0,
		USN:          1,
	}
	store.CreateCard(card)

	// Update flag to 3 (green)
	card.Flag = 3
	err := store.UpdateCard(card)
	if err != nil {
		t.Fatalf("Failed to update card flag: %v", err)
	}

	// Verify flag was updated
	retrieved, err := store.GetCard(1)
	if err != nil {
		t.Fatalf("Failed to get card: %v", err)
	}
	if retrieved.Flag != 3 {
		t.Errorf("Expected flag=3, got %d", retrieved.Flag)
	}
}

func TestUpdateCardMarked(t *testing.T) {
	store, cleanup := setupTestDB(t)
	defer cleanup()

	// Setup
	col := NewCollection()
	store.CreateCollection(col)
	store.CreateDeck(&Deck{ID: 1, Name: "Test", Cards: []int64{}})

	nt := &NoteType{Name: "Basic", Fields: []string{"Front", "Back"}, Templates: []CardTemplate{{Name: "Card 1", QFmt: "{{Front}}", AFmt: "{{Back}}"}}}
	store.CreateNoteType("default", nt)

	now := time.Now()
	store.CreateNote("default", &Note{ID: 1, Type: "Basic", FieldMap: map[string]string{"Front": "Q"}, Tags: []string{}, USN: 1, CreatedAt: now, ModifiedAt: now})

	// Create card with marked=false
	card := &Card{
		ID:           1,
		NoteID:       1,
		DeckID:       1,
		TemplateName: "Card 1",
		Front:        "Q",
		Back:         "A",
		SRS:          newDueNow(now),
		Marked:       false,
		USN:          1,
	}
	store.CreateCard(card)

	// Update marked to true
	card.Marked = true
	err := store.UpdateCard(card)
	if err != nil {
		t.Fatalf("Failed to update card marked: %v", err)
	}

	// Verify marked was updated
	retrieved, err := store.GetCard(1)
	if err != nil {
		t.Fatalf("Failed to get card: %v", err)
	}
	if !retrieved.Marked {
		t.Error("Expected marked=true, got false")
	}
}

func TestUpdateCardSuspended(t *testing.T) {
	store, cleanup := setupTestDB(t)
	defer cleanup()

	// Setup
	col := NewCollection()
	store.CreateCollection(col)
	store.CreateDeck(&Deck{ID: 1, Name: "Test", Cards: []int64{}})

	nt := &NoteType{Name: "Basic", Fields: []string{"Front", "Back"}, Templates: []CardTemplate{{Name: "Card 1", QFmt: "{{Front}}", AFmt: "{{Back}}"}}}
	store.CreateNoteType("default", nt)

	now := time.Now()
	store.CreateNote("default", &Note{ID: 1, Type: "Basic", FieldMap: map[string]string{"Front": "Q"}, Tags: []string{}, USN: 1, CreatedAt: now, ModifiedAt: now})

	// Create card with suspended=false
	card := &Card{
		ID:           1,
		NoteID:       1,
		DeckID:       1,
		TemplateName: "Card 1",
		Front:        "Q",
		Back:         "A",
		SRS:          newDueNow(now),
		Suspended:    false,
		USN:          1,
	}
	store.CreateCard(card)

	// Update suspended to true
	card.Suspended = true
	err := store.UpdateCard(card)
	if err != nil {
		t.Fatalf("Failed to update card suspended: %v", err)
	}

	// Verify suspended was updated
	retrieved, err := store.GetCard(1)
	if err != nil {
		t.Fatalf("Failed to get card: %v", err)
	}
	if !retrieved.Suspended {
		t.Error("Expected suspended=true, got false")
	}

	// Suspended cards should not appear in due cards
	dueCards, err := store.GetDueCards(1, 10)
	if err != nil {
		t.Fatalf("Failed to get due cards: %v", err)
	}
	if len(dueCards) != 0 {
		t.Error("Expected suspended card to not appear in due cards")
	}
}

package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
	fsrs "github.com/open-spaced-repetition/go-fsrs/v3"
)

// Store defines the persistence interface for Microdote.
// All business logic should interact with this interface, not directly with SQL.
type Store interface {
	// Collection
	CreateCollection(c *Collection) error
	GetCollection(id string) (*Collection, error)
	UpdateCollection(c *Collection) error

	// Decks
	CreateDeck(d *Deck) error
	GetDeck(id int64) (*Deck, error)
	UpdateDeck(d *Deck) error
	DeleteDeck(id int64) error
	ListDecks(collectionID string) ([]*Deck, error)

	// Note Types
	CreateNoteType(collectionID string, nt *NoteType) error
	GetNoteType(collectionID string, name NoteTypeName) (*NoteType, error)
	UpdateNoteType(collectionID string, nt *NoteType) error
	ListNoteTypes(collectionID string) (map[NoteTypeName]NoteType, error)

	// Notes
	CreateNote(collectionID string, n *Note) error
	GetNote(id int64) (*Note, error)
	UpdateNote(n *Note) error
	DeleteNote(id int64) error
	ListNotes(collectionID string) (map[int64]Note, error)
	FindDuplicateNotes(collectionID, fieldName, value string, deckID int64) ([]NoteBrief, error)

	// Cards
	CreateCard(c *Card) error
	GetCard(id int64) (*Card, error)
	UpdateCard(c *Card) error
	DeleteCard(id int64) error
	GetDueCards(deckID int64, limit int) ([]*Card, error)
	ListCardsInDeck(deckID int64) ([]*Card, error)
	GetDeckStats(deckID int64) (*DeckStats, error)

	// Revlog
	AddRevlog(r *fsrs.ReviewLog, cardID int64, timeTakenMs int) error
	GetRevlogForCard(cardID int64) ([]*fsrs.ReviewLog, error)

	// Media
	AddMedia(collectionID string, m *MediaRef) error
	GetMedia(filename string) (*MediaRef, error)
	DeleteMedia(filename string) error

	// Profiles (Task 0003)
	CreateProfile(p *Profile) error
	GetProfile(id string) (*Profile, error)
	ListProfiles() ([]*Profile, error)
	SetActiveProfile(id string) error
	GetActiveProfile() (*Profile, error)

	// Transactions
	BeginTx() (*sql.Tx, error)
	CommitTx(tx *sql.Tx) error
	RollbackTx(tx *sql.Tx) error

	// Close database connection
	Close() error
}

// SQLiteStore implements Store using SQLite as the backend.
type SQLiteStore struct {
	db *sql.DB
}

// NewSQLiteStore creates a new SQLite store and runs migrations.
func NewSQLiteStore(dbPath string) (*SQLiteStore, error) {
	db, err := sql.Open("sqlite3", dbPath+"?_foreign_keys=on")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	store := &SQLiteStore{db: db}

	// Run migrations
	if err := store.migrate(); err != nil {
		return nil, fmt.Errorf("migration failed: %w", err)
	}

	return store, nil
}

func (s *SQLiteStore) Close() error {
	return s.db.Close()
}

// Transaction methods
func (s *SQLiteStore) BeginTx() (*sql.Tx, error) {
	return s.db.Begin()
}

func (s *SQLiteStore) CommitTx(tx *sql.Tx) error {
	return tx.Commit()
}

func (s *SQLiteStore) RollbackTx(tx *sql.Tx) error {
	return tx.Rollback()
}

// Collection methods
func (s *SQLiteStore) CreateCollection(c *Collection) error {
	query := `
		INSERT INTO collections (id, name, usn, last_sync, created_at)
		VALUES (?, ?, ?, ?, ?)
	`
	_, err := s.db.Exec(query, "default", "Default Collection", c.USN, c.LastSync.Unix(), time.Now().Unix())
	return err
}

func (s *SQLiteStore) GetCollection(id string) (*Collection, error) {
	query := `SELECT id, name, usn, last_sync, created_at FROM collections WHERE id = ?`
	row := s.db.QueryRow(query, id)

	var name string
	var usn int64
	var lastSync, createdAt int64

	err := row.Scan(&id, &name, &usn, &lastSync, &createdAt)
	if err != nil {
		return nil, err
	}

	// Load collection data
	col := NewCollection()
	col.USN = usn
	if lastSync > 0 {
		col.LastSync = time.Unix(lastSync, 0)
	}

	// Load note types
	noteTypes, err := s.ListNoteTypes(id)
	if err != nil {
		return nil, err
	}
	col.NoteTypes = noteTypes

	// Load notes
	notes, err := s.ListNotes(id)
	if err != nil {
		return nil, err
	}
	col.Notes = notes

	// Load decks (which will load their cards)
	decks, err := s.ListDecks(id)
	if err != nil {
		return nil, err
	}
	col.Decks = make(map[int64]*Deck)
	for _, d := range decks {
		col.Decks[d.ID] = d
	}

	// Load all cards
	col.Cards = make(map[int64]*Card)
	for _, deck := range decks {
		cards, err := s.ListCardsInDeck(deck.ID)
		if err != nil {
			return nil, err
		}
		for _, card := range cards {
			col.Cards[card.ID] = card
		}
	}

	// Initialize ID counters based on existing data
	col.nextDeckID = s.getMaxID("decks") + 1
	col.nextNoteID = s.getMaxID("notes") + 1
	col.nextCardID = s.getMaxID("cards") + 1

	return col, nil
}

// getMaxID returns the maximum ID from a table, or 0 if table is empty
func (s *SQLiteStore) getMaxID(tableName string) int64 {
	var maxID sql.NullInt64
	query := fmt.Sprintf("SELECT MAX(id) FROM %s", tableName)
	s.db.QueryRow(query).Scan(&maxID)
	if maxID.Valid {
		return maxID.Int64
	}
	return 0
}

func (s *SQLiteStore) UpdateCollection(c *Collection) error {
	query := `UPDATE collections SET usn = ?, last_sync = ? WHERE id = ?`
	_, err := s.db.Exec(query, c.USN, c.LastSync.Unix(), "default")
	return err
}

// Deck methods
func (s *SQLiteStore) CreateDeck(d *Deck) error {
	query := `
		INSERT INTO decks (id, collection_id, name, parent_id, options_id)
		VALUES (?, ?, ?, ?, ?)
	`
	_, err := s.db.Exec(query, d.ID, "default", d.Name, d.ParentID, d.OptionsID)
	return err
}

func (s *SQLiteStore) GetDeck(id int64) (*Deck, error) {
	query := `SELECT id, name, parent_id, options_id FROM decks WHERE id = ?`
	row := s.db.QueryRow(query, id)

	var deck Deck
	var parentID, optionsID sql.NullInt64

	err := row.Scan(&deck.ID, &deck.Name, &parentID, &optionsID)
	if err != nil {
		return nil, err
	}

	if parentID.Valid {
		pid := parentID.Int64
		deck.ParentID = &pid
	}
	if optionsID.Valid {
		oid := optionsID.Int64
		deck.OptionsID = &oid
	}

	// Load card IDs for this deck
	cardQuery := `SELECT id FROM cards WHERE deck_id = ? ORDER BY id`
	rows, err := s.db.Query(cardQuery, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	deck.Cards = []int64{}
	for rows.Next() {
		var cardID int64
		if err := rows.Scan(&cardID); err != nil {
			return nil, err
		}
		deck.Cards = append(deck.Cards, cardID)
	}

	return &deck, nil
}

func (s *SQLiteStore) UpdateDeck(d *Deck) error {
	query := `UPDATE decks SET name = ?, parent_id = ?, options_id = ? WHERE id = ?`
	_, err := s.db.Exec(query, d.Name, d.ParentID, d.OptionsID, d.ID)
	return err
}

func (s *SQLiteStore) DeleteDeck(id int64) error {
	// In a full implementation, handle card deletion or moving to another deck
	query := `DELETE FROM decks WHERE id = ?`
	_, err := s.db.Exec(query, id)
	return err
}

func (s *SQLiteStore) ListDecks(collectionID string) ([]*Deck, error) {
	query := `SELECT id FROM decks WHERE collection_id = ? ORDER BY name`
	rows, err := s.db.Query(query, collectionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var decks []*Deck
	for rows.Next() {
		var deckID int64
		if err := rows.Scan(&deckID); err != nil {
			return nil, err
		}
		deck, err := s.GetDeck(deckID)
		if err != nil {
			return nil, err
		}
		decks = append(decks, deck)
	}

	return decks, nil
}

// Note Type methods
func (s *SQLiteStore) CreateNoteType(collectionID string, nt *NoteType) error {
	fieldsJSON, err := json.Marshal(nt.Fields)
	if err != nil {
		return err
	}
	templatesJSON, err := json.Marshal(nt.Templates)
	if err != nil {
		return err
	}
	var fieldOptionsJSON []byte
	if nt.FieldOptions != nil {
		fieldOptionsJSON, err = json.Marshal(nt.FieldOptions)
		if err != nil {
			return err
		}
	}

	query := `
		INSERT INTO note_types (id, collection_id, name, fields, templates, sort_field_index, field_options)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	_, err = s.db.Exec(query, string(nt.Name), collectionID, string(nt.Name), fieldsJSON, templatesJSON, nt.SortFieldIndex, fieldOptionsJSON)
	return err
}

func (s *SQLiteStore) GetNoteType(collectionID string, name NoteTypeName) (*NoteType, error) {
	query := `SELECT name, fields, templates, sort_field_index, field_options FROM note_types WHERE collection_id = ? AND name = ?`
	row := s.db.QueryRow(query, collectionID, string(name))

	var ntName string
	var fieldsJSON, templatesJSON []byte
	var sortFieldIndex int
	var fieldOptionsJSON []byte

	err := row.Scan(&ntName, &fieldsJSON, &templatesJSON, &sortFieldIndex, &fieldOptionsJSON)
	if err != nil {
		return nil, err
	}

	var fields []string
	var templates []CardTemplate
	if err := json.Unmarshal(fieldsJSON, &fields); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(templatesJSON, &templates); err != nil {
		return nil, err
	}

	var fieldOptions map[string]FieldOptions
	if len(fieldOptionsJSON) > 0 {
		if err := json.Unmarshal(fieldOptionsJSON, &fieldOptions); err != nil {
			return nil, err
		}
	}

	return &NoteType{
		Name:           NoteTypeName(ntName),
		Fields:         fields,
		Templates:      templates,
		SortFieldIndex: sortFieldIndex,
		FieldOptions:   fieldOptions,
	}, nil
}

func (s *SQLiteStore) UpdateNoteType(collectionID string, nt *NoteType) error {
	fieldsJSON, err := json.Marshal(nt.Fields)
	if err != nil {
		return err
	}
	templatesJSON, err := json.Marshal(nt.Templates)
	if err != nil {
		return err
	}
	var fieldOptionsJSON []byte
	if nt.FieldOptions != nil {
		fieldOptionsJSON, err = json.Marshal(nt.FieldOptions)
		if err != nil {
			return err
		}
	}

	query := `
		UPDATE note_types
		SET fields = ?, templates = ?, sort_field_index = ?, field_options = ?
		WHERE collection_id = ? AND name = ?
	`
	_, err = s.db.Exec(query, fieldsJSON, templatesJSON, nt.SortFieldIndex, fieldOptionsJSON, collectionID, string(nt.Name))
	return err
}

func (s *SQLiteStore) ListNoteTypes(collectionID string) (map[NoteTypeName]NoteType, error) {
	query := `SELECT name, fields, templates, sort_field_index, field_options FROM note_types WHERE collection_id = ?`
	rows, err := s.db.Query(query, collectionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	noteTypes := make(map[NoteTypeName]NoteType)
	for rows.Next() {
		var name string
		var fieldsJSON, templatesJSON []byte
		var sortFieldIndex int
		var fieldOptionsJSON []byte

		if err := rows.Scan(&name, &fieldsJSON, &templatesJSON, &sortFieldIndex, &fieldOptionsJSON); err != nil {
			return nil, err
		}

		var fields []string
		var templates []CardTemplate
		if err := json.Unmarshal(fieldsJSON, &fields); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(templatesJSON, &templates); err != nil {
			return nil, err
		}

		var fieldOptions map[string]FieldOptions
		if len(fieldOptionsJSON) > 0 {
			if err := json.Unmarshal(fieldOptionsJSON, &fieldOptions); err != nil {
				return nil, err
			}
		}

		noteTypes[NoteTypeName(name)] = NoteType{
			Name:           NoteTypeName(name),
			Fields:         fields,
			Templates:      templates,
			SortFieldIndex: sortFieldIndex,
			FieldOptions:   fieldOptions,
		}
	}

	return noteTypes, nil
}

// Note methods
func (s *SQLiteStore) CreateNote(collectionID string, n *Note) error {
	fieldValsJSON, err := json.Marshal(n.FieldMap)
	if err != nil {
		return err
	}
	tagsJSON, err := json.Marshal(n.Tags)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO notes (id, collection_id, type_id, field_vals, tags, usn, created_at, modified_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err = s.db.Exec(query, n.ID, collectionID, string(n.Type), fieldValsJSON, tagsJSON,
		n.USN, n.CreatedAt.Unix(), n.ModifiedAt.Unix())
	return err
}

func (s *SQLiteStore) GetNote(id int64) (*Note, error) {
	query := `SELECT id, type_id, field_vals, tags, usn, created_at, modified_at FROM notes WHERE id = ?`
	row := s.db.QueryRow(query, id)

	var note Note
	var typeID string
	var fieldValsJSON, tagsJSON []byte
	var createdAt, modifiedAt int64

	err := row.Scan(&note.ID, &typeID, &fieldValsJSON, &tagsJSON, &note.USN, &createdAt, &modifiedAt)
	if err != nil {
		return nil, err
	}

	note.Type = NoteTypeName(typeID)
	note.CreatedAt = time.Unix(createdAt, 0)
	note.ModifiedAt = time.Unix(modifiedAt, 0)

	if err := json.Unmarshal(fieldValsJSON, &note.FieldMap); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(tagsJSON, &note.Tags); err != nil {
		return nil, err
	}

	return &note, nil
}

func (s *SQLiteStore) UpdateNote(n *Note) error {
	fieldValsJSON, err := json.Marshal(n.FieldMap)
	if err != nil {
		return err
	}
	tagsJSON, err := json.Marshal(n.Tags)
	if err != nil {
		return err
	}

	query := `
		UPDATE notes
		SET type_id = ?, field_vals = ?, tags = ?, usn = ?, modified_at = ?
		WHERE id = ?
	`
	_, err = s.db.Exec(query, string(n.Type), fieldValsJSON, tagsJSON, n.USN, n.ModifiedAt.Unix(), n.ID)
	return err
}

func (s *SQLiteStore) DeleteNote(id int64) error {
	// In a full implementation, cascade delete cards
	query := `DELETE FROM notes WHERE id = ?`
	_, err := s.db.Exec(query, id)
	return err
}

func (s *SQLiteStore) ListNotes(collectionID string) (map[int64]Note, error) {
	query := `SELECT id FROM notes WHERE collection_id = ?`
	rows, err := s.db.Query(query, collectionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	notes := make(map[int64]Note)
	for rows.Next() {
		var noteID int64
		if err := rows.Scan(&noteID); err != nil {
			return nil, err
		}
		note, err := s.GetNote(noteID)
		if err != nil {
			return nil, err
		}
		notes[noteID] = *note
	}

	return notes, nil
}

// GetNotesByType returns all notes of a specific note type
func (s *SQLiteStore) GetNotesByType(collectionID string, noteTypeName string) ([]Note, error) {
	query := `SELECT id, type_id, field_vals, tags, usn, created_at, modified_at 
	          FROM notes 
	          WHERE collection_id = ? AND type_id = ?`
	rows, err := s.db.Query(query, collectionID, noteTypeName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []Note
	for rows.Next() {
		var note Note
		var typeID string
		var fieldValsJSON, tagsJSON []byte
		var createdAt, modifiedAt int64

		err := rows.Scan(&note.ID, &typeID, &fieldValsJSON, &tagsJSON, &note.USN, &createdAt, &modifiedAt)
		if err != nil {
			return nil, err
		}

		note.Type = NoteTypeName(typeID)
		note.CreatedAt = time.Unix(createdAt, 0)
		note.ModifiedAt = time.Unix(modifiedAt, 0)

		if err := json.Unmarshal(fieldValsJSON, &note.FieldMap); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(tagsJSON, &note.Tags); err != nil {
			return nil, err
		}

		notes = append(notes, note)
	}

	return notes, rows.Err()
}

func (s *SQLiteStore) FindDuplicateNotes(collectionID, fieldName, value string, deckID int64) ([]NoteBrief, error) {
	// Search notes where the specified field contains the value
	query := `SELECT id, type_id, field_vals FROM notes WHERE collection_id = ?`
	rows, err := s.db.Query(query, collectionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var duplicates []NoteBrief
	normalizedValue := strings.ToLower(strings.TrimSpace(value))

	for rows.Next() {
		var noteID int64
		var typeID string
		var fieldValsJSON []byte

		if err := rows.Scan(&noteID, &typeID, &fieldValsJSON); err != nil {
			return nil, err
		}

		var fieldVals map[string]string
		if err := json.Unmarshal(fieldValsJSON, &fieldVals); err != nil {
			continue
		}

		// Check if the specified field matches (case-insensitive)
		if fieldVal, ok := fieldVals[fieldName]; ok {
			normalizedFieldVal := strings.ToLower(strings.TrimSpace(fieldVal))
			if normalizedFieldVal == normalizedValue {
				// If deckID specified, check if any card for this note is in that deck
				if deckID > 0 {
					inDeck, _ := s.noteHasCardInDeck(noteID, deckID)
					if !inDeck {
						continue
					}
				}
				duplicates = append(duplicates, NoteBrief{
					ID:       noteID,
					TypeID:   typeID,
					FieldVal: fieldVals,
				})
			}
		}
	}

	return duplicates, nil
}

func (s *SQLiteStore) noteHasCardInDeck(noteID, deckID int64) (bool, error) {
	query := `SELECT COUNT(*) FROM cards WHERE note_id = ? AND deck_id = ?`
	var count int
	err := s.db.QueryRow(query, noteID, deckID).Scan(&count)
	return count > 0, err
}

// Card methods
func (s *SQLiteStore) CreateCard(c *Card) error {
	fsrsJSON, err := json.Marshal(c.SRS)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO cards (id, note_id, deck_id, template_name, ordinal, front, back,
		                   due, state, fsrs_data, flag, marked, suspended, usn)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err = s.db.Exec(query, c.ID, c.NoteID, c.DeckID, c.TemplateName, c.Ordinal, c.Front, c.Back,
		c.SRS.Due.Unix(), int(c.SRS.State), fsrsJSON, c.Flag, c.Marked, c.Suspended, c.USN)
	return err
}

// GetCardsByNote returns all cards for a given note
func (s *SQLiteStore) GetCardsByNote(noteID int64) ([]Card, error) {
	query := `
		SELECT id, note_id, deck_id, template_name, ordinal, front, back,
		       due, state, fsrs_data, flag, marked, suspended, usn
		FROM cards WHERE note_id = ?
	`
	rows, err := s.db.Query(query, noteID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cards []Card
	for rows.Next() {
		var card Card
		var dueUnix int64
		var state int
		var fsrsJSON []byte
		var marked, suspended int

		err := rows.Scan(&card.ID, &card.NoteID, &card.DeckID, &card.TemplateName, &card.Ordinal,
			&card.Front, &card.Back, &dueUnix, &state, &fsrsJSON, &card.Flag, &marked, &suspended, &card.USN)
		if err != nil {
			return nil, err
		}

		card.SRS.Due = time.Unix(dueUnix, 0)
		card.SRS.State = fsrs.State(state)
		card.Marked = marked != 0
		card.Suspended = suspended != 0

		if err := json.Unmarshal(fsrsJSON, &card.SRS); err != nil {
			return nil, err
		}

		cards = append(cards, card)
	}

	return cards, rows.Err()
}

func (s *SQLiteStore) GetCard(id int64) (*Card, error) {
	query := `
		SELECT id, note_id, deck_id, template_name, ordinal, front, back,
		       due, state, fsrs_data, flag, marked, suspended, usn
		FROM cards WHERE id = ?
	`
	row := s.db.QueryRow(query, id)

	var card Card
	var dueUnix int64
	var state int
	var fsrsJSON []byte
	var marked, suspended int

	err := row.Scan(&card.ID, &card.NoteID, &card.DeckID, &card.TemplateName, &card.Ordinal,
		&card.Front, &card.Back, &dueUnix, &state, &fsrsJSON, &card.Flag, &marked, &suspended, &card.USN)
	if err != nil {
		return nil, err
	}

	card.Marked = marked == 1
	card.Suspended = suspended == 1

	if err := json.Unmarshal(fsrsJSON, &card.SRS); err != nil {
		return nil, err
	}
	card.SRS.Due = time.Unix(dueUnix, 0)
	card.SRS.State = fsrs.State(state)

	return &card, nil
}

func (s *SQLiteStore) UpdateCard(c *Card) error {
	fsrsJSON, err := json.Marshal(c.SRS)
	if err != nil {
		return err
	}

	query := `
		UPDATE cards
		SET note_id = ?, deck_id = ?, template_name = ?, ordinal = ?, front = ?, back = ?,
		    due = ?, state = ?, fsrs_data = ?, flag = ?, marked = ?, suspended = ?, usn = ?
		WHERE id = ?
	`
	_, err = s.db.Exec(query, c.NoteID, c.DeckID, c.TemplateName, c.Ordinal, c.Front, c.Back,
		c.SRS.Due.Unix(), int(c.SRS.State), fsrsJSON, c.Flag, c.Marked, c.Suspended, c.USN, c.ID)
	return err
}

func (s *SQLiteStore) DeleteCard(id int64) error {
	query := `DELETE FROM cards WHERE id = ?`
	_, err := s.db.Exec(query, id)
	return err
}

func (s *SQLiteStore) GetDueCards(deckID int64, limit int) ([]*Card, error) {
	now := time.Now().Unix()
	query := `
		SELECT id FROM cards
		WHERE deck_id = ? AND due <= ? AND suspended = 0
		ORDER BY due
		LIMIT ?
	`
	rows, err := s.db.Query(query, deckID, now, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cards []*Card
	for rows.Next() {
		var cardID int64
		if err := rows.Scan(&cardID); err != nil {
			return nil, err
		}
		card, err := s.GetCard(cardID)
		if err != nil {
			return nil, err
		}
		cards = append(cards, card)
	}

	return cards, nil
}

func (s *SQLiteStore) ListCardsInDeck(deckID int64) ([]*Card, error) {
	query := `SELECT id FROM cards WHERE deck_id = ? ORDER BY id`
	rows, err := s.db.Query(query, deckID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cards []*Card
	for rows.Next() {
		var cardID int64
		if err := rows.Scan(&cardID); err != nil {
			return nil, err
		}
		card, err := s.GetCard(cardID)
		if err != nil {
			return nil, err
		}
		cards = append(cards, card)
	}

	return cards, nil
}

// Revlog methods
func (s *SQLiteStore) AddRevlog(r *fsrs.ReviewLog, cardID int64, timeTakenMs int) error {
	query := `
		INSERT INTO revlog (id, card_id, rating, state, due, reviewed_at, time_taken_ms)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	// Generate ID (in real implementation, use proper ID generation)
	id := time.Now().UnixNano()
	_, err := s.db.Exec(query, id, cardID, int(r.Rating), int(r.State), r.Review.Unix(), r.Review.Unix(), timeTakenMs)
	return err
}

func (s *SQLiteStore) GetRevlogForCard(cardID int64) ([]*fsrs.ReviewLog, error) {
	query := `SELECT rating, state, due, reviewed_at FROM revlog WHERE card_id = ? ORDER BY reviewed_at`
	rows, err := s.db.Query(query, cardID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []*fsrs.ReviewLog
	for rows.Next() {
		var rating, state int
		var dueUnix, reviewedAt int64

		if err := rows.Scan(&rating, &state, &dueUnix, &reviewedAt); err != nil {
			return nil, err
		}

		log := &fsrs.ReviewLog{
			Rating: fsrs.Rating(rating),
			State:  fsrs.State(state),
			Review: time.Unix(reviewedAt, 0),
		}
		logs = append(logs, log)
	}

	return logs, nil
}

// Media methods
func (s *SQLiteStore) AddMedia(collectionID string, m *MediaRef) error {
	query := `
		INSERT INTO media (id, collection_id, filename, data, added_at)
		VALUES (?, ?, ?, ?, ?)
	`
	_, err := s.db.Exec(query, m.ID, collectionID, m.Filename, m.Data, m.AddedAt.Unix())
	return err
}

func (s *SQLiteStore) GetMedia(filename string) (*MediaRef, error) {
	query := `SELECT id, filename, data, added_at FROM media WHERE filename = ?`
	row := s.db.QueryRow(query, filename)

	var m MediaRef
	var addedAt int64

	err := row.Scan(&m.ID, &m.Filename, &m.Data, &addedAt)
	if err != nil {
		return nil, err
	}

	m.AddedAt = time.Unix(addedAt, 0)
	return &m, nil
}

func (s *SQLiteStore) DeleteMedia(filename string) error {
	query := `DELETE FROM media WHERE filename = ?`
	_, err := s.db.Exec(query, filename)
	return err
}

// GetDeckStats returns card counts by state for a deck
func (s *SQLiteStore) GetDeckStats(deckID int64) (*DeckStats, error) {
	stats := &DeckStats{DeckID: deckID}
	now := time.Now().Unix()

	// Get all cards for the deck
	query := `SELECT state, suspended, due FROM cards WHERE deck_id = ?`
	rows, err := s.db.Query(query, deckID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var state, suspended int
		var due int64

		if err := rows.Scan(&state, &suspended, &due); err != nil {
			return nil, err
		}

		stats.TotalCards++

		if suspended == 1 {
			stats.Suspended++
			continue
		}

		// fsrs.State values: New=0, Learning=1, Review=2, Relearning=3
		switch state {
		case 0: // New
			stats.NewCards++
			if due <= now {
				stats.DueToday++
			}
		case 1: // Learning
			stats.Learning++
			if due <= now {
				stats.DueToday++
			}
		case 2: // Review
			stats.Review++
			if due <= now {
				stats.DueToday++
			}
		case 3: // Relearning
			stats.Relearning++
			if due <= now {
				stats.DueToday++
			}
		}
	}

	return stats, nil
}

// Profile methods (Task 0003)

func (s *SQLiteStore) CreateProfile(p *Profile) error {
	query := `
		INSERT INTO profiles (id, name, collection_id, sync_account, created_at)
		VALUES (?, ?, ?, ?, ?)
	`
	_, err := s.db.Exec(query, p.ID, p.Name, p.CollectionID, p.SyncAccount, p.CreatedAt.Unix())
	return err
}

func (s *SQLiteStore) GetProfile(id string) (*Profile, error) {
	query := `SELECT id, name, collection_id, sync_account, created_at FROM profiles WHERE id = ?`
	row := s.db.QueryRow(query, id)

	var p Profile
	var syncAccount sql.NullString
	var createdAt int64

	err := row.Scan(&p.ID, &p.Name, &p.CollectionID, &syncAccount, &createdAt)
	if err != nil {
		return nil, err
	}

	if syncAccount.Valid {
		p.SyncAccount = syncAccount.String
	}
	p.CreatedAt = time.Unix(createdAt, 0)

	return &p, nil
}

func (s *SQLiteStore) ListProfiles() ([]*Profile, error) {
	query := `SELECT id, name, collection_id, sync_account, created_at FROM profiles ORDER BY name`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var profiles []*Profile
	for rows.Next() {
		var p Profile
		var syncAccount sql.NullString
		var createdAt int64

		err := rows.Scan(&p.ID, &p.Name, &p.CollectionID, &syncAccount, &createdAt)
		if err != nil {
			return nil, err
		}

		if syncAccount.Valid {
			p.SyncAccount = syncAccount.String
		}
		p.CreatedAt = time.Unix(createdAt, 0)

		profiles = append(profiles, &p)
	}

	return profiles, nil
}

func (s *SQLiteStore) SetActiveProfile(id string) error {
	query := `INSERT OR REPLACE INTO metadata (key, value) VALUES ('active_profile', ?)`
	_, err := s.db.Exec(query, id)
	return err
}

func (s *SQLiteStore) GetActiveProfile() (*Profile, error) {
	// Get active profile ID from metadata
	var profileID string
	err := s.db.QueryRow("SELECT value FROM metadata WHERE key = 'active_profile'").Scan(&profileID)
	if err == sql.ErrNoRows {
		// No active profile set, return default if exists
		return s.getOrCreateDefaultProfile()
	}
	if err != nil {
		return nil, err
	}

	return s.GetProfile(profileID)
}

func (s *SQLiteStore) getOrCreateDefaultProfile() (*Profile, error) {
	// Try to get "default" profile
	profile, err := s.GetProfile("default")
	if err == nil {
		// Set as active
		s.SetActiveProfile("default")
		return profile, nil
	}

	// Ensure default collection exists first (foreign key constraint)
	_, colErr := s.GetCollection("default")
	if colErr != nil {
		// Create default collection
		defaultCol := NewCollection()
		if err := s.CreateCollection(defaultCol); err != nil {
			return nil, fmt.Errorf("failed to create default collection: %w", err)
		}
	}

	// Create default profile
	profile = &Profile{
		ID:           "default",
		Name:         "Default",
		CollectionID: "default",
		SyncAccount:  "",
		CreatedAt:    time.Now(),
	}

	if err := s.CreateProfile(profile); err != nil {
		return nil, err
	}

	s.SetActiveProfile("default")
	return profile, nil
}

package main

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	fsrs "github.com/open-spaced-repetition/go-fsrs/v3"
)

/*
Microdote single-file prototype:
- NoteTypes define Fields + Templates
- Notes hold field values
- Cards are generated from Notes + Templates
- Scheduling is FSRS (modern Anki-style scheduler)
  - desired retention + weights (parameters) per preset is the plan
*/

type NoteTypeName string

type CardTemplate struct {
	Name            string `json:"name"`
	QFmt            string `json:"qFmt"`
	AFmt            string `json:"aFmt"`
	Styling         string `json:"styling"`
	IfFieldNonEmpty string `json:"ifFieldNonEmpty"`
	IsCloze         bool   `json:"isCloze"`
}

type NoteType struct {
	Name      NoteTypeName   `json:"name"`
	Fields    []string       `json:"fields"`
	Templates []CardTemplate `json:"templates"`
}

type Note struct {
	ID         int64             `json:"id"`
	Type       NoteTypeName      `json:"type"`
	FieldMap   map[string]string `json:"fieldMap"`   // field name -> value
	Tags       []string          `json:"tags"`       // user-defined tags for organization
	USN        int64             `json:"usn"`        // Update Sequence Number for sync
	CreatedAt  time.Time         `json:"createdAt"`  // when note was created
	ModifiedAt time.Time         `json:"modifiedAt"` // when note was last modified
}

type Card struct {
	ID           int64  `json:"id"`
	NoteID       int64  `json:"noteId"`
	DeckID       int64  `json:"deckId"`
	TemplateName string `json:"templateName"` // which template generated this card
	Ordinal      int    `json:"ordinal"`      // e.g. cloze c1,c2 => ordinal=1,2; non-cloze => 0

	Front string `json:"front"`
	Back  string `json:"back"`

	SRS fsrs.Card `json:"srs"` // FSRS state: due, stability, difficulty, reps, lapses, etc.

	// Additional metadata for Anki parity
	Flag      int   `json:"flag"`      // 0=none, 1-7=color flags for marking cards
	Marked    bool  `json:"marked"`    // special "marked" tag for review
	Suspended bool  `json:"suspended"` // whether card is suspended (excluded from study)
	USN       int64 `json:"usn"`       // Update Sequence Number for sync
}

type Deck struct {
	ID        int64 `json:"id"`
	Name      string
	Cards     []int64 // card IDs
	ParentID  *int64  // for deck hierarchy (nil if root deck)
	OptionsID *int64  // reference to DeckOptions preset (nil = use default)
}

// DeckOptions represents scheduling/behavior presets that can be shared across decks.
// This is the foundation for M4 deck options (Tasks 0401-0405).
type DeckOptions struct {
	ID                 int64
	Name               string
	NewCardsPerDay     int   // daily limit for new cards
	ReviewsPerDay      int   // daily limit for reviews
	LearningSteps      []int // learning steps in minutes (e.g. [1, 10])
	GraduatingInterval int   // days until a learning card becomes review card
	EasyInterval       int   // days for "easy" button on new card
	// Future: add more options from Tasks 0402-0405 (lapses, relearning, etc.)
}

// MediaRef represents a media file (image, audio, video) referenced by notes.
type MediaRef struct {
	ID       int64
	Filename string    // unique filename
	Data     []byte    // file contents (or path if stored externally)
	AddedAt  time.Time // when media was added
}

// Profile represents a user profile with its own collection.
// Each profile can have separate decks, notes, and sync settings.
type Profile struct {
	ID           string
	Name         string
	CollectionID string
	SyncAccount  string // optional: linked sync account
	CreatedAt    time.Time
}

// DeckStats represents card counts for a deck by state.
type DeckStats struct {
	DeckID     int64 `json:"deckId"`
	NewCards   int   `json:"newCards"`
	Learning   int   `json:"learning"`
	Review     int   `json:"review"`
	Relearning int   `json:"relearning"`
	Suspended  int   `json:"suspended"`
	Buried     int   `json:"buried"`
	TotalCards int   `json:"totalCards"`
	DueToday   int   `json:"dueToday"`
}

type Collection struct {
	nextNoteID int64
	nextCardID int64
	nextDeckID int64

	NoteTypes map[NoteTypeName]NoteType `json:"noteTypes"`
	Notes     map[int64]Note            `json:"notes"`
	Cards     map[int64]*Card           `json:"cards"`
	Decks     map[int64]*Deck           `json:"decks"`

	// FSRS config (in a real app: per preset)
	Params fsrs.Parameters `json:"params"`

	Revlog []fsrs.ReviewLog `json:"revlog"`

	// Media storage
	Media map[string]*MediaRef `json:"media"` // filename -> media ref

	// Sync metadata
	USN      int64     `json:"usn"`      // Update Sequence Number (increments on each change)
	LastSync time.Time `json:"lastSync"` // when collection was last synced
}

func NewCollection() *Collection {
	p := fsrs.DefaultParam()
	// Match Anki conceptually: desired retention ~ 0.90-0.95 common starting points
	p.RequestRetention = 0.90
	// MaximumInterval is in days in go-fsrs params
	p.MaximumInterval = 36500 // ~100 years; tune later

	return &Collection{
		nextNoteID: 1,
		nextCardID: 1,
		nextDeckID: 1,
		NoteTypes:  make(map[NoteTypeName]NoteType),
		Notes:      make(map[int64]Note),
		Cards:      make(map[int64]*Card),
		Decks:      make(map[int64]*Deck),
		Params:     p,
		Revlog:     nil,
		Media:      make(map[string]*MediaRef),
		USN:        0,
		LastSync:   time.Time{}, // zero time = never synced
	}
}

func (c *Collection) NewDeck(name string) *Deck {
	id := c.nextDeckID
	c.nextDeckID++
	d := &Deck{ID: id, Name: name, Cards: []int64{}}
	c.Decks[id] = d
	return d
}

func (c *Collection) AddNote(deckID int64, ntName NoteTypeName, fields map[string]string, now time.Time) (Note, []*Card, error) {
	nt, ok := c.NoteTypes[ntName]
	if !ok {
		return Note{}, nil, fmt.Errorf("unknown note type: %s", ntName)
	}

	noteID := c.nextNoteID
	c.nextNoteID++
	c.USN++ // increment collection USN on modification

	n := Note{
		ID:         noteID,
		Type:       ntName,
		FieldMap:   fields,
		Tags:       []string{}, // empty tags by default
		USN:        c.USN,      // track when note was created
		CreatedAt:  now,
		ModifiedAt: now,
	}
	c.Notes[noteID] = n

	genCards, err := generateCardsFromNote(nt, n, deckID, now)
	if err != nil {
		return Note{}, nil, err
	}

	var out []*Card
	for _, card := range genCards {
		cardID := c.nextCardID
		c.nextCardID++

		card.ID = cardID
		card.USN = c.USN // track when card was created
		c.Cards[cardID] = card
		out = append(out, card)

		if d, ok := c.Decks[deckID]; ok {
			d.Cards = append(d.Cards, cardID)
		}
	}
	return n, out, nil
}

// GenerateCards creates cards from a note using its note type templates.
// This is used when changing a note's note type (note type migration).
// It regenerates all cards for the note, preserving stable card IDs where possible.
//
// Migration strategy:
// 1. Delete old cards that no longer match any template
// 2. Keep existing cards that still match (preserve scheduling state)
// 3. Create new cards for new templates
// 4. Update card content (Front/Back) for all cards
func (c *Collection) GenerateCards(note *Note, deckID int64, now time.Time) ([]*Card, error) {
	nt, ok := c.NoteTypes[note.Type]
	if !ok {
		return nil, fmt.Errorf("unknown note type: %s", note.Type)
	}

	// Generate fresh cards based on current templates
	newCards, err := generateCardsFromNote(nt, *note, deckID, now)
	if err != nil {
		return nil, err
	}

	// In a full implementation:
	// - Match new cards to existing cards by template name + ordinal
	// - Reuse existing card IDs to preserve scheduling history
	// - Delete orphaned cards
	// For now, this is a placeholder that just generates new cards.
	// Task 0222 (conditional generation logic) will implement full matching.

	return newCards, nil
}

// Answer a card with Again/Hard/Good/Easy and update FSRS state.
func (c *Collection) Answer(cardID int64, rating fsrs.Rating, now time.Time, timeTakenMs int) (*fsrs.ReviewLog, error) {
	card, ok := c.Cards[cardID]
	if !ok {
		return nil, fmt.Errorf("unknown card id: %d", cardID)
	}
	sched := fsrs.NewFSRS(c.Params).Repeat(card.SRS, now)

	info, ok := sched[rating]
	if !ok {
		return nil, fmt.Errorf("no scheduling info for rating %v", rating)
	}

	card.SRS = info.Card
	c.Revlog = append(c.Revlog, info.ReviewLog)
	return &info.ReviewLog, nil
}

// Pick the next due card in a deck (very simplified).
func (c *Collection) NextDue(deckID int64, now time.Time) (*Card, bool) {
	d, ok := c.Decks[deckID]
	if !ok {
		return nil, false
	}

	var dueCards []*Card
	for _, cid := range d.Cards {
		card := c.Cards[cid]
		if card == nil {
			continue
		}
		if !card.SRS.Due.IsZero() && (card.SRS.Due.Before(now) || card.SRS.Due.Equal(now)) {
			dueCards = append(dueCards, card)
		}
	}

	if len(dueCards) == 0 {
		return nil, false
	}

	sort.Slice(dueCards, func(i, j int) bool {
		return dueCards[i].SRS.Due.Before(dueCards[j].SRS.Due)
	})
	return dueCards[0], true
}

/* --------------------------
   Note → Card generation
-------------------------- */

var fieldTokenRe = regexp.MustCompile(`\{\{([^}]+)\}\}`)

// Cloze pattern like {{c1::1969}} or {{c2::answer::hint}}
var clozeRe = regexp.MustCompile(`\{\{c(\d+)::(.*?)(?:::([^}]*))?\}\}`)

func generateCardsFromNote(nt NoteType, n Note, deckID int64, now time.Time) ([]*Card, error) {
	var cards []*Card

	for _, tmpl := range nt.Templates {
		// optional template generation
		if tmpl.IfFieldNonEmpty != "" {
			if strings.TrimSpace(n.FieldMap[tmpl.IfFieldNonEmpty]) == "" {
				continue
			}
		}

		if tmpl.IsCloze {
			textField := n.FieldMap["Text"]
			ordinals := extractClozeOrdinals(textField)
			for _, ord := range ordinals {
				q := renderTemplateWithCloze(tmpl.QFmt, n.FieldMap, ord, false)
				a := renderTemplateWithCloze(tmpl.AFmt, n.FieldMap, ord, true)
				card := &Card{
					NoteID:       n.ID,
					DeckID:       deckID,
					TemplateName: tmpl.Name,
					Ordinal:      ord,
					Front:        q,
					Back:         a,
					SRS:          newDueNow(now),
				}
				cards = append(cards, card)
			}
			continue
		}

		q := renderTemplate(tmpl.QFmt, n.FieldMap)
		a := renderTemplate(tmpl.AFmt, n.FieldMap)

		card := &Card{
			NoteID:       n.ID,
			DeckID:       deckID,
			TemplateName: tmpl.Name,
			Ordinal:      0,
			Front:        q,
			Back:         a,
			SRS:          newDueNow(now),
		}
		cards = append(cards, card)
	}

	return cards, nil
}

func newDueNow(now time.Time) fsrs.Card {
	c := fsrs.NewCard()
	c.Due = now
	return c
}

func renderTemplate(tmpl string, fields map[string]string) string {
	return fieldTokenRe.ReplaceAllStringFunc(tmpl, func(token string) string {
		m := fieldTokenRe.FindStringSubmatch(token)
		if len(m) != 2 {
			return token
		}
		key := strings.TrimSpace(m[1])

		// very small compatibility shim for "type in the answer"
		// Anki uses {{type:Back}} etc. We'll render a placeholder.
		if strings.HasPrefix(key, "type:") {
			fieldName := strings.TrimSpace(strings.TrimPrefix(key, "type:"))
			expected := fields[fieldName]
			if expected == "" {
				return "[type: empty]"
			}
			return "[type your answer here]"
		}

		// normal {{Field}} replacement
		return fields[key]
	})
}

func renderTemplateWithCloze(tmpl string, fields map[string]string, targetOrdinal int, reveal bool) string {
	// First replace {{cloze:Text}} tokens (Anki style)
	out := fieldTokenRe.ReplaceAllStringFunc(tmpl, func(token string) string {
		m := fieldTokenRe.FindStringSubmatch(token)
		if len(m) != 2 {
			return token
		}
		key := strings.TrimSpace(m[1])
		if key == "cloze:Text" {
			return renderCloze(fields["Text"], targetOrdinal, reveal)
		}
		// fallback: normal replacement
		return fields[key]
	})
	return out
}

func extractClozeOrdinals(text string) []int {
	seen := map[int]bool{}
	matches := clozeRe.FindAllStringSubmatch(text, -1)
	for _, m := range matches {
		if len(m) < 2 {
			continue
		}
		n, err := strconv.Atoi(m[1])
		if err != nil || n <= 0 {
			continue
		}
		seen[n] = true
	}
	var ord []int
	for k := range seen {
		ord = append(ord, k)
	}
	sort.Ints(ord)
	return ord
}

func renderCloze(text string, targetOrdinal int, reveal bool) string {
	// Replace each {{cN::answer::hint}}:
	// - On FRONT (reveal=false): target -> "[...]" or "[hint]" ; others -> answer
	// - On BACK  (reveal=true) : all -> answer (target highlighted)
	return clozeRe.ReplaceAllStringFunc(text, func(token string) string {
		m := clozeRe.FindStringSubmatch(token)
		if len(m) < 3 {
			return token
		}
		ord, _ := strconv.Atoi(m[1])
		answer := m[2]
		hint := ""
		if len(m) >= 4 {
			hint = m[3]
		}

		if reveal {
			if ord == targetOrdinal {
				return fmt.Sprintf("**%s**", answer)
			}
			return answer
		}

		// front side
		if ord == targetOrdinal {
			if strings.TrimSpace(hint) != "" {
				return fmt.Sprintf("[%s]", hint)
			}
			return "[...]"
		}
		return answer
	})
}

/* --------------------------
   Built-in Note Types
-------------------------- */

// InitDefaultCollection initializes or loads the default collection from SQLite.
// If the collection doesn't exist, it creates one with built-in note types.
// Also ensures a default profile exists.
func InitDefaultCollection(dbPath string) (*Collection, *SQLiteStore, error) {
	store, err := NewSQLiteStore(dbPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create store: %w", err)
	}

	// Ensure default profile exists and is active
	profile, err := store.GetActiveProfile()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get active profile: %w", err)
	}
	fmt.Printf("Active profile: %s (%s)\n", profile.Name, profile.ID)

	// Try to load existing collection for this profile
	col, err := store.GetCollection(profile.CollectionID)
	if err != nil {
		// Collection doesn't exist, create a new one
		fmt.Println("Creating new collection with built-in note types...")
		col = NewCollection()

		// Create collection record
		if err := store.CreateCollection(col); err != nil {
			return nil, nil, fmt.Errorf("failed to create collection: %w", err)
		}
	}

	// Ensure built-in note types exist (whether new or existing collection)
	noteTypes := builtins()
	if len(col.NoteTypes) == 0 {
		col.NoteTypes = make(map[NoteTypeName]NoteType)
		for _, nt := range noteTypes {
			// Check if note type already exists in DB
			_, err := store.GetNoteType("default", nt.Name)
			if err != nil {
				// Doesn't exist, create it
				if err := store.CreateNoteType("default", &nt); err != nil {
					return nil, nil, fmt.Errorf("failed to create note type %s: %w", nt.Name, err)
				}
			}
			col.NoteTypes[nt.Name] = nt
		}
	}

	return col, store, nil
}

func builtins() map[NoteTypeName]NoteType {
	return map[NoteTypeName]NoteType{
		"Basic": {
			Name:   "Basic",
			Fields: []string{"Front", "Back"},
			Templates: []CardTemplate{
				{
					Name: "Card 1",
					QFmt: "Q: {{Front}}",
					AFmt: "A: {{Back}}",
				},
			},
		},
		"Basic (and reversed card)": {
			Name:   "Basic (and reversed card)",
			Fields: []string{"Front", "Back"},
			Templates: []CardTemplate{
				{Name: "Card 1", QFmt: "Q: {{Front}}", AFmt: "A: {{Back}}"},
				{Name: "Card 2", QFmt: "Q: {{Back}}", AFmt: "A: {{Front}}"},
			},
		},
		"Basic (optional reversed card)": {
			Name:   "Basic (optional reversed card)",
			Fields: []string{"Front", "Back", "Add Reverse"},
			Templates: []CardTemplate{
				{Name: "Card 1", QFmt: "Q: {{Front}}", AFmt: "A: {{Back}}"},
				{
					Name:            "Card 2 (optional reverse)",
					QFmt:            "Q: {{Back}}",
					AFmt:            "A: {{Front}}",
					IfFieldNonEmpty: "Add Reverse",
				},
			},
		},
		"Basic (type in the answer)": {
			Name:   "Basic (type in the answer)",
			Fields: []string{"Front", "Back"},
			Templates: []CardTemplate{
				{
					Name: "Card 1",
					QFmt: "Q: {{Front}}\n\n{{type:Back}}",
					AFmt: "A: {{Back}}",
				},
			},
		},
		"Cloze": {
			Name:   "Cloze",
			Fields: []string{"Text", "Extra"},
			Templates: []CardTemplate{
				{
					Name:    "Cloze",
					QFmt:    "Q: {{cloze:Text}}",
					AFmt:    "A: {{cloze:Text}}\n\nExtra: {{Extra}}",
					IsCloze: true,
				},
			},
		},
	}
}

/* --------------------------
   Demo / Testing Functions
   (moved to separate file or removed - server.go is now the entry point)
-------------------------- */

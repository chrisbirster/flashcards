package main

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type PlandalfStore struct {
	db *sql.DB
}

type plandalfSQL interface {
	Exec(query string, args ...any) (sql.Result, error)
	Query(query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
}

type PlandalfDeckSummary struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	CardCount  int    `json:"card_count"`
	DueCount   int    `json:"due_count"`
	NewCount   int    `json:"new_count"`
}

type PlandalfCardRecord struct {
	ID         int64
	DeckID     int64
	Question   string
	Answer     string
	DueAtMs    *int64
	ReviewCount int
}

type PlandalfSchedulerState struct {
	StabilityDays   float64
	Difficulty      float64
	DueAtMs         int64
	LastReviewedAtMs int64
}

func OpenPlandalfStore(cfg DatabaseConfig) (*PlandalfStore, error) {
	driverName, dsn, err := databaseDSN(cfg)
	if err != nil {
		return nil, err
	}
	db, err := sql.Open(driverName, dsn)
	if err != nil {
		return nil, fmt.Errorf("open Plandalf database: %w", err)
	}
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("ping Plandalf database: %w", err)
	}
	store := &PlandalfStore{db: db}
	if err := store.migrate(); err != nil {
		db.Close()
		return nil, err
	}
	return store, nil
}

func (s *PlandalfStore) Close() error { return s.db.Close() }

func (s *PlandalfStore) migrate() error {
	statements := []string{
		`PRAGMA foreign_keys = ON`,
		`CREATE TABLE IF NOT EXISTS parameter_sets (
			id BLOB PRIMARY KEY CHECK(length(id) = 32),
			algorithm_family TEXT NOT NULL,
			algorithm_major INTEGER NOT NULL,
			implementation_major INTEGER NOT NULL,
			implementation_minor INTEGER NOT NULL,
			implementation_patch INTEGER NOT NULL,
			source TEXT NOT NULL,
			parameters_json TEXT NOT NULL DEFAULT '[]',
			desired_retention REAL NOT NULL,
			minimum_interval_days REAL NOT NULL,
			maximum_interval_days REAL NOT NULL,
			created_at_ms INTEGER NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS parameter_weights (
			parameter_set_id BLOB NOT NULL REFERENCES parameter_sets(id) ON DELETE CASCADE,
			position INTEGER NOT NULL CHECK(position >= 0),
			value REAL NOT NULL,
			PRIMARY KEY(parameter_set_id, position)
		)`,
		`CREATE TABLE IF NOT EXISTS deck_groups (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL UNIQUE,
			algorithm_family TEXT NULL,
			algorithm_major INTEGER NULL,
			parameter_set_id BLOB NULL REFERENCES parameter_sets(id),
			created_at_ms INTEGER NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS scheduler_defaults (
			id INTEGER PRIMARY KEY CHECK(id = 1),
			algorithm_family TEXT NOT NULL DEFAULT 'fsrs',
			algorithm_major INTEGER NOT NULL DEFAULT 7,
			parameter_set_id BLOB NULL REFERENCES parameter_sets(id)
		)`,
		`INSERT OR IGNORE INTO scheduler_defaults(id, algorithm_family, algorithm_major) VALUES (1, 'fsrs', 7)`,
		`CREATE TABLE IF NOT EXISTS decks (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL UNIQUE,
			group_id INTEGER NULL REFERENCES deck_groups(id) ON DELETE SET NULL,
			algorithm_family TEXT NULL DEFAULT 'fsrs',
			algorithm_major INTEGER NULL DEFAULT 7,
			parameter_set_id BLOB NULL REFERENCES parameter_sets(id),
			created_at_ms INTEGER NOT NULL
		)`,
		`CREATE INDEX IF NOT EXISTS decks_group_id_idx ON decks(group_id)`,
		`CREATE TABLE IF NOT EXISTS cards (
			id INTEGER PRIMARY KEY,
			deck_id INTEGER NOT NULL REFERENCES decks(id) ON DELETE CASCADE,
			question TEXT NOT NULL,
			answer TEXT NOT NULL,
			created_at_ms INTEGER NOT NULL
		)`,
		`CREATE INDEX IF NOT EXISTS cards_deck_id_idx ON cards(deck_id)`,
		`CREATE TABLE IF NOT EXISTS reviews (
			id INTEGER PRIMARY KEY,
			card_id INTEGER NOT NULL REFERENCES cards(id) ON DELETE RESTRICT,
			rating INTEGER NOT NULL CHECK(rating BETWEEN 1 AND 4),
			reviewed_at_ms INTEGER NOT NULL,
			algorithm_family TEXT NULL,
			algorithm_major INTEGER NULL,
			implementation_major INTEGER NULL,
			implementation_minor INTEGER NULL,
			implementation_patch INTEGER NULL,
			parameter_set_id BLOB NULL REFERENCES parameter_sets(id),
			scheduled_at_ms INTEGER NULL
		)`,
		`CREATE INDEX IF NOT EXISTS reviews_card_time_idx ON reviews(card_id, reviewed_at_ms, id)`,
		`CREATE INDEX IF NOT EXISTS reviews_time_card_rating_idx ON reviews(reviewed_at_ms, card_id, rating)`,
		`CREATE TRIGGER IF NOT EXISTS reviews_immutable_update BEFORE UPDATE ON reviews BEGIN SELECT RAISE(ABORT, 'review history is immutable'); END`,
		`CREATE TRIGGER IF NOT EXISTS reviews_immutable_delete BEFORE DELETE ON reviews BEGIN SELECT RAISE(ABORT, 'review history is immutable'); END`,
		`CREATE TABLE IF NOT EXISTS scheduler_state (
			card_id INTEGER PRIMARY KEY REFERENCES cards(id) ON DELETE CASCADE,
			algorithm_family TEXT NOT NULL,
			algorithm_major INTEGER NOT NULL,
			implementation_major INTEGER NOT NULL,
			implementation_minor INTEGER NOT NULL,
			implementation_patch INTEGER NOT NULL,
			parameter_set_id BLOB NOT NULL REFERENCES parameter_sets(id),
			stability_days REAL NULL,
			difficulty REAL NULL,
			due_at_ms INTEGER NOT NULL,
			last_reviewed_at_ms INTEGER NULL
		)`,
		`CREATE INDEX IF NOT EXISTS scheduler_state_due_idx ON scheduler_state(due_at_ms)`,
		`CREATE TABLE IF NOT EXISTS note_types (
			id INTEGER PRIMARY KEY,
			slug TEXT NOT NULL UNIQUE,
			name TEXT NOT NULL,
			kind TEXT NOT NULL,
			css TEXT NOT NULL,
			created_at_ms INTEGER NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS note_type_fields (
			note_type_id INTEGER NOT NULL REFERENCES note_types(id) ON DELETE CASCADE,
			ordinal INTEGER NOT NULL,
			name TEXT NOT NULL,
			PRIMARY KEY(note_type_id, ordinal), UNIQUE(note_type_id, name)
		)`,
		`CREATE TABLE IF NOT EXISTS card_templates (
			note_type_id INTEGER NOT NULL REFERENCES note_types(id) ON DELETE CASCADE,
			ordinal INTEGER NOT NULL,
			name TEXT NOT NULL,
			front TEXT NOT NULL,
			back TEXT NOT NULL,
			PRIMARY KEY(note_type_id, ordinal)
		)`,
		`CREATE TABLE IF NOT EXISTS notes (
			id INTEGER PRIMARY KEY,
			note_type_id INTEGER NOT NULL REFERENCES note_types(id),
			tags_json TEXT NOT NULL DEFAULT '[]',
			created_at_ms INTEGER NOT NULL,
			updated_at_ms INTEGER NOT NULL
		)`,
		`CREATE INDEX IF NOT EXISTS notes_note_type_idx ON notes(note_type_id)`,
		`CREATE TABLE IF NOT EXISTS note_fields (
			note_id INTEGER NOT NULL REFERENCES notes(id) ON DELETE CASCADE,
			ordinal INTEGER NOT NULL,
			value TEXT NOT NULL,
			PRIMARY KEY(note_id, ordinal)
		)`,
		`CREATE TABLE IF NOT EXISTS generated_cards (
			card_id INTEGER PRIMARY KEY REFERENCES cards(id) ON DELETE CASCADE,
			note_id INTEGER NOT NULL REFERENCES notes(id) ON DELETE CASCADE,
			template_ordinal INTEGER NOT NULL,
			generation_key TEXT NOT NULL UNIQUE
		)`,
		`CREATE INDEX IF NOT EXISTS generated_cards_note_idx ON generated_cards(note_id, template_ordinal)`,
	}
	for _, statement := range statements {
		if _, err := s.db.Exec(statement); err != nil {
			return fmt.Errorf("migrate Plandalf schema: %w", err)
		}
	}
	return s.ensureDefaultFSRS7()
}

func (s *PlandalfStore) ensureDefaultFSRS7() error {
	p := defaultPlandalfFSRS7Parameters()
	id := plandalfParameterSetID(p)
	nowMs := time.Now().UnixMilli()
	_, err := s.db.Exec(`INSERT OR IGNORE INTO parameter_sets(
		id, algorithm_family, algorithm_major, implementation_major, implementation_minor,
		implementation_patch, source, parameters_json, desired_retention,
		minimum_interval_days, maximum_interval_days, created_at_ms
	) VALUES (?, 'fsrs', 7, 0, 1, 0, 'default', '[]', ?, ?, ?, ?)`,
		id[:], p.DesiredRetention, p.MinimumIntervalDays, p.MaximumIntervalDays, nowMs)
	if err != nil { return err }
	for position, weight := range p.Weights {
		if _, err := s.db.Exec(`INSERT OR IGNORE INTO parameter_weights(parameter_set_id, position, value) VALUES (?, ?, ?)`, id[:], position, weight); err != nil {
			return err
		}
	}
	_, err = s.db.Exec(`UPDATE scheduler_defaults SET algorithm_family='fsrs', algorithm_major=7, parameter_set_id=? WHERE id=1 AND parameter_set_id IS NULL`, id[:])
	return err
}

func (s *PlandalfStore) CreateDeck(name string, nowMs int64) (int64, error) {
	result, err := s.db.Exec(`INSERT INTO decks(name, created_at_ms) VALUES (?, ?)`, name, nowMs)
	if err != nil { return 0, err }
	return result.LastInsertId()
}

func (s *PlandalfStore) CreateCard(deckID int64, question, answer string, nowMs int64) (int64, error) {
	result, err := s.db.Exec(`INSERT INTO cards(deck_id, question, answer, created_at_ms) VALUES (?, ?, ?, ?)`, deckID, question, answer, nowMs)
	if err != nil { return 0, err }
	return result.LastInsertId()
}

func (s *PlandalfStore) ListDecks(nowMs int64) ([]PlandalfDeckSummary, error) {
	rows, err := s.db.Query(`
		SELECT d.id, d.name,
		       COUNT(c.id) AS card_count,
		       SUM(CASE WHEN c.id IS NOT NULL AND ss.card_id IS NULL THEN 1 ELSE 0 END) AS new_count,
		       SUM(CASE WHEN ss.card_id IS NOT NULL AND ss.due_at_ms <= ? THEN 1 ELSE 0 END) AS due_count
		FROM decks d
		LEFT JOIN cards c ON c.deck_id = d.id
		LEFT JOIN scheduler_state ss ON ss.card_id = c.id
		GROUP BY d.id, d.name
		ORDER BY d.name COLLATE NOCASE, d.id`, nowMs)
	if err != nil { return nil, err }
	defer rows.Close()
	var decks []PlandalfDeckSummary
	for rows.Next() {
		var d PlandalfDeckSummary
		var id int64
		if err := rows.Scan(&id, &d.Name, &d.CardCount, &d.NewCount, &d.DueCount); err != nil { return nil, err }
		d.ID = fmt.Sprintf("%d", id)
		decks = append(decks, d)
	}
	return decks, rows.Err()
}

func (s *PlandalfStore) GetCard(cardID int64) (*PlandalfCardRecord, error) {
	var card PlandalfCardRecord
	var due sql.NullInt64
	err := s.db.QueryRow(`
		SELECT c.id, c.deck_id, c.question, c.answer, ss.due_at_ms,
		       (SELECT COUNT(*) FROM reviews r WHERE r.card_id = c.id)
		FROM cards c LEFT JOIN scheduler_state ss ON ss.card_id = c.id
		WHERE c.id = ?`, cardID).Scan(&card.ID, &card.DeckID, &card.Question, &card.Answer, &due, &card.ReviewCount)
	if errors.Is(err, sql.ErrNoRows) { return nil, nil }
	if err != nil { return nil, err }
	if due.Valid { value := due.Int64; card.DueAtMs = &value }
	return &card, nil
}

func (s *PlandalfStore) NextCard(deckID int64, nowMs int64) (*PlandalfCardRecord, error) {
	var card PlandalfCardRecord
	var due sql.NullInt64
	err := s.db.QueryRow(`
		SELECT c.id, c.deck_id, c.question, c.answer, ss.due_at_ms,
		       (SELECT COUNT(*) FROM reviews r WHERE r.card_id = c.id)
		FROM cards c
		LEFT JOIN scheduler_state ss ON ss.card_id = c.id
		WHERE c.deck_id = ? AND (ss.card_id IS NULL OR ss.due_at_ms <= ?)
		ORDER BY CASE WHEN ss.card_id IS NULL THEN 1 ELSE 0 END,
		         COALESCE(ss.due_at_ms, c.created_at_ms), c.id
		LIMIT 1`, deckID, nowMs).Scan(&card.ID, &card.DeckID, &card.Question, &card.Answer, &due, &card.ReviewCount)
	if errors.Is(err, sql.ErrNoRows) { return nil, nil }
	if err != nil { return nil, err }
	if due.Valid { value := due.Int64; card.DueAtMs = &value }
	return &card, nil
}

func loadPlandalfHistory(q plandalfSQL, cardID int64) ([]plandalfHistoryEntry, error) {
	rows, err := q.Query(`SELECT rating, reviewed_at_ms FROM reviews WHERE card_id=? ORDER BY reviewed_at_ms, id`, cardID)
	if err != nil { return nil, err }
	defer rows.Close()
	var history []plandalfHistoryEntry
	for rows.Next() {
		var rating int
		var reviewedAt int64
		if err := rows.Scan(&rating, &reviewedAt); err != nil { return nil, err }
		parsed, err := parsePlandalfRating(rating)
		if err != nil { return nil, err }
		history = append(history, plandalfHistoryEntry{Rating: parsed, ReviewedAtMs: reviewedAt})
	}
	return history, rows.Err()
}

func loadPlandalfParameters(q plandalfSQL, deckID int64) (plandalfFSRS7Parameters, []byte, error) {
	var family string
	var major int
	var parameterID []byte
	err := q.QueryRow(`
		SELECT COALESCE(d.algorithm_family, g.algorithm_family, s.algorithm_family),
		       COALESCE(d.algorithm_major, g.algorithm_major, s.algorithm_major),
		       COALESCE(d.parameter_set_id, g.parameter_set_id, s.parameter_set_id)
		FROM decks d
		LEFT JOIN deck_groups g ON g.id=d.group_id
		CROSS JOIN scheduler_defaults s
		WHERE d.id=? AND s.id=1`, deckID).Scan(&family, &major, &parameterID)
	if err != nil { return plandalfFSRS7Parameters{}, nil, err }
	if family != "fsrs" || major != 7 || len(parameterID) != 32 {
		return plandalfFSRS7Parameters{}, nil, errors.New("unsupported scheduler configuration")
	}
	p := defaultPlandalfFSRS7Parameters()
	if err := q.QueryRow(`SELECT desired_retention, minimum_interval_days, maximum_interval_days FROM parameter_sets WHERE id=? AND algorithm_family='fsrs' AND algorithm_major=7`, parameterID).
		Scan(&p.DesiredRetention, &p.MinimumIntervalDays, &p.MaximumIntervalDays); err != nil {
		return plandalfFSRS7Parameters{}, nil, err
	}
	rows, err := q.Query(`SELECT position, value FROM parameter_weights WHERE parameter_set_id=? ORDER BY position`, parameterID)
	if err != nil { return plandalfFSRS7Parameters{}, nil, err }
	defer rows.Close()
	count := 0
	for rows.Next() {
		var position int
		var value float64
		if err := rows.Scan(&position, &value); err != nil { return plandalfFSRS7Parameters{}, nil, err }
		if position != count || position < 0 || position >= len(p.Weights) { return plandalfFSRS7Parameters{}, nil, errors.New("invalid parameter weights") }
		p.Weights[position] = value
		count++
	}
	if err := rows.Err(); err != nil { return plandalfFSRS7Parameters{}, nil, err }
	if count != len(p.Weights) { return plandalfFSRS7Parameters{}, nil, errors.New("incomplete parameter weights") }
	return p, append([]byte(nil), parameterID...), nil
}

func (s *PlandalfStore) Preview(cardID int64, nowMs int64) (plandalfSchedule, int, error) {
	card, err := s.GetCard(cardID)
	if err != nil { return plandalfSchedule{}, 0, err }
	if card == nil { return plandalfSchedule{}, 0, sql.ErrNoRows }
	history, err := loadPlandalfHistory(s.db, cardID)
	if err != nil { return plandalfSchedule{}, 0, err }
	p, _, err := loadPlandalfParameters(s.db, card.DeckID)
	if err != nil { return plandalfSchedule{}, 0, err }
	schedule, err := plandalfScheduleFor(history, nowMs, p)
	return schedule, len(history), err
}

func candidateForRating(schedule plandalfSchedule, rating plandalfRating) plandalfCandidate {
	switch rating {
	case plandalfAgain: return schedule.Again
	case plandalfHard: return schedule.Hard
	case plandalfGood: return schedule.Good
	default: return schedule.Easy
	}
}

func (s *PlandalfStore) RecordReview(cardID int64, rating plandalfRating, expectedReviewCount int, reviewedAtMs int64) (int64, plandalfCandidate, PlandalfSchedulerState, error) {
	tx, err := s.db.Begin()
	if err != nil { return 0, plandalfCandidate{}, PlandalfSchedulerState{}, err }
	defer tx.Rollback()

	var deckID int64
	if err := tx.QueryRow(`SELECT deck_id FROM cards WHERE id=?`, cardID).Scan(&deckID); err != nil {
		return 0, plandalfCandidate{}, PlandalfSchedulerState{}, err
	}
	history, err := loadPlandalfHistory(tx, cardID)
	if err != nil { return 0, plandalfCandidate{}, PlandalfSchedulerState{}, err }
	if len(history) != expectedReviewCount {
		return 0, plandalfCandidate{}, PlandalfSchedulerState{}, errors.New("stale_review")
	}
	var due sql.NullInt64
	if err := tx.QueryRow(`SELECT due_at_ms FROM scheduler_state WHERE card_id=?`, cardID).Scan(&due); err != nil && !errors.Is(err, sql.ErrNoRows) {
		return 0, plandalfCandidate{}, PlandalfSchedulerState{}, err
	}
	if due.Valid && due.Int64 > reviewedAtMs {
		return 0, plandalfCandidate{}, PlandalfSchedulerState{}, errors.New("card_not_due")
	}
	p, parameterID, err := loadPlandalfParameters(tx, deckID)
	if err != nil { return 0, plandalfCandidate{}, PlandalfSchedulerState{}, err }
	schedule, err := plandalfScheduleFor(history, reviewedAtMs, p)
	if err != nil { return 0, plandalfCandidate{}, PlandalfSchedulerState{}, err }
	candidate := candidateForRating(schedule, rating)
	combined := append(append([]plandalfHistoryEntry(nil), history...), plandalfHistoryEntry{Rating: rating, ReviewedAtMs: reviewedAtMs})
	replayed, err := plandalfReplay(combined, p)
	if err != nil || replayed == nil { return 0, plandalfCandidate{}, PlandalfSchedulerState{}, errors.New("missing replay state") }

	result, err := tx.Exec(`INSERT INTO reviews(card_id, rating, reviewed_at_ms, algorithm_family, algorithm_major, implementation_major, implementation_minor, implementation_patch, parameter_set_id, scheduled_at_ms) VALUES (?, ?, ?, 'fsrs', 7, 0, 1, 0, ?, ?)`,
		cardID, int(rating), reviewedAtMs, parameterID, candidate.DueAtMs)
	if err != nil { return 0, plandalfCandidate{}, PlandalfSchedulerState{}, err }
	reviewID, err := result.LastInsertId()
	if err != nil { return 0, plandalfCandidate{}, PlandalfSchedulerState{}, err }
	state := PlandalfSchedulerState{
		StabilityDays: replayed.Memory.StabilityDays,
		Difficulty: replayed.Memory.Difficulty,
		DueAtMs: candidate.DueAtMs,
		LastReviewedAtMs: reviewedAtMs,
	}
	_, err = tx.Exec(`INSERT INTO scheduler_state(card_id, algorithm_family, algorithm_major, implementation_major, implementation_minor, implementation_patch, parameter_set_id, stability_days, difficulty, due_at_ms, last_reviewed_at_ms)
		VALUES (?, 'fsrs', 7, 0, 1, 0, ?, ?, ?, ?, ?)
		ON CONFLICT(card_id) DO UPDATE SET algorithm_family=excluded.algorithm_family, algorithm_major=excluded.algorithm_major, implementation_major=excluded.implementation_major, implementation_minor=excluded.implementation_minor, implementation_patch=excluded.implementation_patch, parameter_set_id=excluded.parameter_set_id, stability_days=excluded.stability_days, difficulty=excluded.difficulty, due_at_ms=excluded.due_at_ms, last_reviewed_at_ms=excluded.last_reviewed_at_ms`,
		cardID, parameterID, state.StabilityDays, state.Difficulty, state.DueAtMs, state.LastReviewedAtMs)
	if err != nil { return 0, plandalfCandidate{}, PlandalfSchedulerState{}, err }
	if err := tx.Commit(); err != nil { return 0, plandalfCandidate{}, PlandalfSchedulerState{}, err }
	return reviewID, candidate, state, nil
}

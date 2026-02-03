package main

import (
	"database/sql"
	"fmt"
)

// migrate runs database migrations to ensure schema is up to date.
func (s *SQLiteStore) migrate() error {
	// Ensure metadata table exists first
	if err := s.ensureMetadataTable(); err != nil {
		return err
	}

	version, err := s.getSchemaVersion()
	if err != nil {
		return err
	}

	// Run migrations sequentially
	migrations := []struct {
		version int
		name    string
		fn      func() error
	}{
		{1, "initial_schema", s.runMigration001_InitialSchema},
		// Future migrations go here
		// {2, "add_deck_stats", s.runMigration002_AddDeckStats},
	}

	for _, m := range migrations {
		if version < m.version {
			fmt.Printf("Running migration %d: %s\n", m.version, m.name)
			if err := m.fn(); err != nil {
				return fmt.Errorf("migration %d failed: %w", m.version, err)
			}
			if err := s.setSchemaVersion(m.version); err != nil {
				return fmt.Errorf("failed to update schema version: %w", err)
			}
			version = m.version
		}
	}

	fmt.Printf("Database schema up to date (version %d)\n", version)
	return nil
}

func (s *SQLiteStore) ensureMetadataTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS metadata (
			key TEXT PRIMARY KEY,
			value TEXT
		)
	`
	_, err := s.db.Exec(query)
	return err
}

func (s *SQLiteStore) getSchemaVersion() (int, error) {
	var version int
	err := s.db.QueryRow("SELECT value FROM metadata WHERE key = 'schema_version'").Scan(&version)
	if err == sql.ErrNoRows {
		return 0, nil // No version set = version 0
	}
	return version, err
}

func (s *SQLiteStore) setSchemaVersion(version int) error {
	query := `
		INSERT OR REPLACE INTO metadata (key, value)
		VALUES ('schema_version', ?)
	`
	_, err := s.db.Exec(query, fmt.Sprintf("%d", version))
	return err
}

// runMigration001_InitialSchema creates the initial database schema for M0.
func (s *SQLiteStore) runMigration001_InitialSchema() error {
	schema := `
	-- Collections table
	CREATE TABLE IF NOT EXISTS collections (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		usn INTEGER DEFAULT 0,
		last_sync INTEGER,
		created_at INTEGER
	);

	-- Profiles table (for Task 0003)
	CREATE TABLE IF NOT EXISTS profiles (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		collection_id TEXT,
		sync_account TEXT,
		created_at INTEGER,
		FOREIGN KEY (collection_id) REFERENCES collections(id)
	);

	-- Decks table
	CREATE TABLE IF NOT EXISTS decks (
		id INTEGER PRIMARY KEY,
		collection_id TEXT NOT NULL,
		name TEXT NOT NULL,
		parent_id INTEGER,
		options_id INTEGER,
		FOREIGN KEY (collection_id) REFERENCES collections(id),
		FOREIGN KEY (parent_id) REFERENCES decks(id),
		FOREIGN KEY (options_id) REFERENCES deck_options(id)
	);

	-- Deck options table (for M4, foundation in M0)
	CREATE TABLE IF NOT EXISTS deck_options (
		id INTEGER PRIMARY KEY,
		name TEXT NOT NULL,
		new_cards_per_day INTEGER DEFAULT 20,
		reviews_per_day INTEGER DEFAULT 200,
		learning_steps TEXT,
		graduating_interval INTEGER DEFAULT 1,
		easy_interval INTEGER DEFAULT 4
	);

	-- Note types table
	CREATE TABLE IF NOT EXISTS note_types (
		id TEXT PRIMARY KEY,
		collection_id TEXT NOT NULL,
		name TEXT NOT NULL,
		fields TEXT NOT NULL,
		templates TEXT NOT NULL,
		FOREIGN KEY (collection_id) REFERENCES collections(id)
	);

	-- Notes table
	CREATE TABLE IF NOT EXISTS notes (
		id INTEGER PRIMARY KEY,
		collection_id TEXT NOT NULL,
		type_id TEXT NOT NULL,
		field_vals TEXT NOT NULL,
		tags TEXT,
		usn INTEGER DEFAULT 0,
		created_at INTEGER,
		modified_at INTEGER,
		FOREIGN KEY (collection_id) REFERENCES collections(id),
		FOREIGN KEY (type_id) REFERENCES note_types(id)
	);

	-- Cards table
	CREATE TABLE IF NOT EXISTS cards (
		id INTEGER PRIMARY KEY,
		note_id INTEGER NOT NULL,
		deck_id INTEGER NOT NULL,
		template_name TEXT NOT NULL,
		ordinal INTEGER DEFAULT 0,
		front TEXT,
		back TEXT,
		due INTEGER,
		state INTEGER,
		fsrs_data TEXT,
		flag INTEGER DEFAULT 0,
		marked INTEGER DEFAULT 0,
		suspended INTEGER DEFAULT 0,
		usn INTEGER DEFAULT 0,
		FOREIGN KEY (note_id) REFERENCES notes(id) ON DELETE CASCADE,
		FOREIGN KEY (deck_id) REFERENCES decks(id)
	);

	-- Review log table
	CREATE TABLE IF NOT EXISTS revlog (
		id INTEGER PRIMARY KEY,
		card_id INTEGER NOT NULL,
		rating INTEGER NOT NULL,
		state INTEGER,
		due INTEGER,
		reviewed_at INTEGER,
		time_taken_ms INTEGER DEFAULT 0,
		FOREIGN KEY (card_id) REFERENCES cards(id)
	);

	-- Media table
	CREATE TABLE IF NOT EXISTS media (
		id INTEGER PRIMARY KEY,
		collection_id TEXT NOT NULL,
		filename TEXT UNIQUE NOT NULL,
		data BLOB,
		added_at INTEGER,
		FOREIGN KEY (collection_id) REFERENCES collections(id)
	);

	-- Indexes for performance
	CREATE INDEX IF NOT EXISTS idx_cards_due ON cards(due, deck_id);
	CREATE INDEX IF NOT EXISTS idx_cards_note ON cards(note_id);
	CREATE INDEX IF NOT EXISTS idx_cards_deck ON cards(deck_id);
	CREATE INDEX IF NOT EXISTS idx_revlog_card ON revlog(card_id, reviewed_at);
	CREATE INDEX IF NOT EXISTS idx_notes_collection ON notes(collection_id);
	CREATE INDEX IF NOT EXISTS idx_notes_type ON notes(type_id);
	CREATE INDEX IF NOT EXISTS idx_decks_collection ON decks(collection_id);
	CREATE INDEX IF NOT EXISTS idx_decks_parent ON decks(parent_id);
	`

	_, err := s.db.Exec(schema)
	if err != nil {
		return fmt.Errorf("failed to create schema: %w", err)
	}

	return nil
}

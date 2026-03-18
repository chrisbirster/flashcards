package main

import (
	"database/sql"
	"fmt"
	"strings"
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
		{2, "add_field_options", s.runMigration002_AddFieldOptions},
		{3, "add_account_billing_schema", s.runMigration003_AddAccountBillingSchema},
		{4, "add_otp_auth_schema", s.runMigration004_AddOTPAuthSchema},
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

// runMigration002_AddFieldOptions adds sort_field_index and field_options columns to note_types.
func (s *SQLiteStore) runMigration002_AddFieldOptions() error {
	schema := `
	ALTER TABLE note_types ADD COLUMN sort_field_index INTEGER DEFAULT 0;
	ALTER TABLE note_types ADD COLUMN field_options TEXT;
	`

	_, err := s.db.Exec(schema)
	if err != nil {
		return fmt.Errorf("failed to add field options columns: %w", err)
	}

	return nil
}

func (s *SQLiteStore) runMigration003_AddAccountBillingSchema() error {
	schema := `
	CREATE TABLE IF NOT EXISTS users (
		id TEXT PRIMARY KEY,
		email TEXT NOT NULL UNIQUE,
		display_name TEXT NOT NULL,
		avatar_url TEXT,
		created_at INTEGER NOT NULL,
		updated_at INTEGER NOT NULL
	);

	CREATE TABLE IF NOT EXISTS oauth_identities (
		id TEXT PRIMARY KEY,
		user_id TEXT NOT NULL,
		provider TEXT NOT NULL,
		subject TEXT NOT NULL,
		email TEXT NOT NULL,
		created_at INTEGER NOT NULL,
		UNIQUE(provider, subject),
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS organizations (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		slug TEXT NOT NULL UNIQUE,
		created_at INTEGER NOT NULL,
		updated_at INTEGER NOT NULL
	);

	CREATE TABLE IF NOT EXISTS workspaces (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		slug TEXT NOT NULL UNIQUE,
		collection_id TEXT NOT NULL,
		owner_user_id TEXT,
		organization_id TEXT,
		created_at INTEGER NOT NULL,
		updated_at INTEGER NOT NULL,
		FOREIGN KEY (collection_id) REFERENCES collections(id),
		FOREIGN KEY (owner_user_id) REFERENCES users(id),
		FOREIGN KEY (organization_id) REFERENCES organizations(id)
	);

	CREATE TABLE IF NOT EXISTS organization_members (
		id TEXT PRIMARY KEY,
		organization_id TEXT NOT NULL,
		user_id TEXT,
		email TEXT NOT NULL,
		role TEXT NOT NULL,
		status TEXT NOT NULL,
		created_at INTEGER NOT NULL,
		UNIQUE(organization_id, email),
		FOREIGN KEY (organization_id) REFERENCES organizations(id) ON DELETE CASCADE,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL
	);

	CREATE TABLE IF NOT EXISTS sessions (
		id TEXT PRIMARY KEY,
		user_id TEXT NOT NULL,
		workspace_id TEXT,
		plan TEXT NOT NULL,
		guest INTEGER NOT NULL DEFAULT 0,
		expires_at INTEGER,
		created_at INTEGER NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
		FOREIGN KEY (workspace_id) REFERENCES workspaces(id) ON DELETE SET NULL
	);

	CREATE TABLE IF NOT EXISTS subscriptions (
		id TEXT PRIMARY KEY,
		workspace_id TEXT,
		organization_id TEXT,
		plan TEXT NOT NULL,
		status TEXT NOT NULL,
		provider TEXT,
		provider_customer_id TEXT,
		provider_subscription_id TEXT,
		current_period_end INTEGER,
		created_at INTEGER NOT NULL,
		updated_at INTEGER NOT NULL,
		FOREIGN KEY (workspace_id) REFERENCES workspaces(id) ON DELETE CASCADE,
		FOREIGN KEY (organization_id) REFERENCES organizations(id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS subscription_events (
		id TEXT PRIMARY KEY,
		subscription_id TEXT NOT NULL,
		event_type TEXT NOT NULL,
		provider_event_id TEXT,
		payload TEXT,
		created_at INTEGER NOT NULL,
		UNIQUE(provider_event_id),
		FOREIGN KEY (subscription_id) REFERENCES subscriptions(id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS entitlements (
		id TEXT PRIMARY KEY,
		workspace_id TEXT,
		organization_id TEXT,
		plan TEXT NOT NULL,
		max_decks INTEGER NOT NULL,
		max_notes INTEGER NOT NULL,
		max_shared_decks INTEGER NOT NULL,
		max_sync_devices INTEGER NOT NULL,
		max_workspaces INTEGER NOT NULL,
		created_at INTEGER NOT NULL,
		updated_at INTEGER NOT NULL,
		FOREIGN KEY (workspace_id) REFERENCES workspaces(id) ON DELETE CASCADE,
		FOREIGN KEY (organization_id) REFERENCES organizations(id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS deck_shares (
		id TEXT PRIMARY KEY,
		deck_id INTEGER NOT NULL,
		workspace_id TEXT,
		created_by_user_id TEXT,
		token TEXT NOT NULL UNIQUE,
		access_type TEXT NOT NULL,
		created_at INTEGER NOT NULL,
		FOREIGN KEY (deck_id) REFERENCES decks(id) ON DELETE CASCADE,
		FOREIGN KEY (workspace_id) REFERENCES workspaces(id) ON DELETE CASCADE,
		FOREIGN KEY (created_by_user_id) REFERENCES users(id) ON DELETE SET NULL
	);

	CREATE TABLE IF NOT EXISTS sync_devices (
		id TEXT PRIMARY KEY,
		workspace_id TEXT NOT NULL,
		name TEXT NOT NULL,
		platform TEXT,
		created_at INTEGER NOT NULL,
		updated_at INTEGER NOT NULL,
		UNIQUE(workspace_id, name),
		FOREIGN KEY (workspace_id) REFERENCES workspaces(id) ON DELETE CASCADE
	);

	CREATE INDEX IF NOT EXISTS idx_oauth_identities_lookup ON oauth_identities(provider, subject);
	CREATE INDEX IF NOT EXISTS idx_sessions_user ON sessions(user_id);
	CREATE INDEX IF NOT EXISTS idx_workspaces_owner ON workspaces(owner_user_id);
	CREATE INDEX IF NOT EXISTS idx_subscriptions_workspace ON subscriptions(workspace_id);
	CREATE INDEX IF NOT EXISTS idx_deck_shares_deck ON deck_shares(deck_id);
	CREATE INDEX IF NOT EXISTS idx_sync_devices_workspace ON sync_devices(workspace_id);
	`

	_, err := s.db.Exec(schema)
	if err != nil {
		return fmt.Errorf("failed to add account and billing schema: %w", err)
	}

	return nil
}

func (s *SQLiteStore) runMigration004_AddOTPAuthSchema() error {
	statements := []string{
		`ALTER TABLE users ADD COLUMN last_login_at INTEGER`,
		`ALTER TABLE sessions ADD COLUMN last_seen_at INTEGER`,
		`
		CREATE TABLE IF NOT EXISTS otp_challenges (
			id TEXT PRIMARY KEY,
			email TEXT NOT NULL,
			code_hash TEXT NOT NULL,
			expires_at INTEGER NOT NULL,
			attempt_count INTEGER NOT NULL DEFAULT 0,
			max_attempts INTEGER NOT NULL,
			resend_available_at INTEGER NOT NULL,
			consumed_at INTEGER,
			requested_ip TEXT,
			user_agent TEXT,
			created_at INTEGER NOT NULL
		)
		`,
		`CREATE INDEX IF NOT EXISTS idx_otp_challenges_email_created ON otp_challenges(email, created_at DESC)`,
		`CREATE INDEX IF NOT EXISTS idx_otp_challenges_email_consumed ON otp_challenges(email, consumed_at)`,
	}

	for _, statement := range statements {
		if _, err := s.db.Exec(statement); err != nil && !isIgnorableMigrationError(err) {
			return fmt.Errorf("failed to apply OTP auth migration statement: %w", err)
		}
	}

	return nil
}

func isIgnorableMigrationError(err error) bool {
	if err == nil {
		return false
	}
	message := strings.ToLower(err.Error())
	return strings.Contains(message, "duplicate column name")
}

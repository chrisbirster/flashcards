package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"
)

var slugPattern = regexp.MustCompile(`[^a-z0-9]+`)

func newID(prefix string) string {
	bytes := make([]byte, 8)
	if _, err := rand.Read(bytes); err != nil {
		return fmt.Sprintf("%s_%d", prefix, time.Now().UnixNano())
	}
	return fmt.Sprintf("%s_%s", prefix, hex.EncodeToString(bytes))
}

func slugify(value string) string {
	slug := strings.ToLower(strings.TrimSpace(value))
	slug = slugPattern.ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")
	if slug == "" {
		return "workspace"
	}
	return slug
}

func firstFieldPreview(fieldMap map[string]string) string {
	for _, value := range fieldMap {
		trimmed := strings.TrimSpace(value)
		if trimmed != "" {
			return trimmed
		}
	}
	return ""
}

func (s *SQLiteStore) CreateUser(user *User) error {
	query := `
		INSERT INTO users (id, email, display_name, avatar_url, last_login_at, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	_, err := s.db.Exec(query, user.ID, user.Email, user.DisplayName, user.AvatarURL, nullIfZeroTime(user.LastLoginAt), user.CreatedAt.Unix(), user.UpdatedAt.Unix())
	return err
}

func (s *SQLiteStore) GetUserByID(id string) (*User, error) {
	query := `SELECT id, email, display_name, avatar_url, last_login_at, created_at, updated_at FROM users WHERE id = ?`
	row := s.db.QueryRow(query, id)

	var user User
	var avatar sql.NullString
	var lastLoginAt, createdAt, updatedAt sql.NullInt64
	if err := row.Scan(&user.ID, &user.Email, &user.DisplayName, &avatar, &lastLoginAt, &createdAt, &updatedAt); err != nil {
		return nil, err
	}

	if avatar.Valid {
		user.AvatarURL = avatar.String
	}
	if lastLoginAt.Valid {
		user.LastLoginAt = time.Unix(lastLoginAt.Int64, 0)
	}
	if createdAt.Valid {
		user.CreatedAt = time.Unix(createdAt.Int64, 0)
	}
	if updatedAt.Valid {
		user.UpdatedAt = time.Unix(updatedAt.Int64, 0)
	}
	return &user, nil
}

func (s *SQLiteStore) GetUserByEmail(email string) (*User, error) {
	query := `SELECT id FROM users WHERE lower(email) = lower(?)`
	var userID string
	if err := s.db.QueryRow(query, email).Scan(&userID); err != nil {
		return nil, err
	}
	return s.GetUserByID(userID)
}

func (s *SQLiteStore) GetUserByOAuth(provider, subject string) (*User, error) {
	query := `SELECT user_id FROM oauth_identities WHERE provider = ? AND subject = ?`
	var userID string
	if err := s.db.QueryRow(query, provider, subject).Scan(&userID); err != nil {
		return nil, err
	}
	return s.GetUserByID(userID)
}

func (s *SQLiteStore) UpsertOAuthIdentity(identity *OAuthIdentity) error {
	query := `
		INSERT INTO oauth_identities (id, user_id, provider, subject, email, created_at)
		VALUES (?, ?, ?, ?, ?, ?)
		ON CONFLICT(provider, subject) DO UPDATE SET
			user_id = excluded.user_id,
			email = excluded.email
	`
	_, err := s.db.Exec(query, identity.ID, identity.UserID, identity.Provider, identity.Subject, identity.Email, identity.CreatedAt.Unix())
	return err
}

func (s *SQLiteStore) CreateWorkspaceRecord(workspace *Workspace) error {
	query := `
		INSERT INTO workspaces (id, name, slug, collection_id, owner_user_id, organization_id, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err := s.db.Exec(
		query,
		workspace.ID,
		workspace.Name,
		workspace.Slug,
		workspace.CollectionID,
		workspace.OwnerUserID,
		nullIfEmpty(workspace.OrganizationID),
		workspace.CreatedAt.Unix(),
		workspace.UpdatedAt.Unix(),
	)
	return err
}

func (s *SQLiteStore) GetWorkspaceRecord(id string) (*Workspace, error) {
	query := `
		SELECT id, name, slug, collection_id, owner_user_id, organization_id, created_at, updated_at
		FROM workspaces WHERE id = ?
	`
	row := s.db.QueryRow(query, id)

	var workspace Workspace
	var ownerID, orgID sql.NullString
	var createdAt, updatedAt int64
	if err := row.Scan(
		&workspace.ID,
		&workspace.Name,
		&workspace.Slug,
		&workspace.CollectionID,
		&ownerID,
		&orgID,
		&createdAt,
		&updatedAt,
	); err != nil {
		return nil, err
	}

	if ownerID.Valid {
		workspace.OwnerUserID = ownerID.String
	}
	if orgID.Valid {
		workspace.OrganizationID = orgID.String
	}
	workspace.CreatedAt = time.Unix(createdAt, 0)
	workspace.UpdatedAt = time.Unix(updatedAt, 0)
	return &workspace, nil
}

func (s *SQLiteStore) GetFirstWorkspaceForUser(userID string) (*Workspace, error) {
	query := `SELECT id FROM workspaces WHERE owner_user_id = ? ORDER BY created_at ASC LIMIT 1`
	var workspaceID string
	if err := s.db.QueryRow(query, userID).Scan(&workspaceID); err != nil {
		return nil, err
	}
	return s.GetWorkspaceRecord(workspaceID)
}

func (s *SQLiteStore) CountWorkspacesForUser(userID string) (int, error) {
	query := `SELECT COUNT(*) FROM workspaces WHERE owner_user_id = ?`
	var count int
	err := s.db.QueryRow(query, userID).Scan(&count)
	return count, err
}

func (s *SQLiteStore) CreateOrganizationRecord(org *Organization) error {
	query := `
		INSERT INTO organizations (id, name, slug, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)
	`
	_, err := s.db.Exec(query, org.ID, org.Name, org.Slug, org.CreatedAt.Unix(), org.UpdatedAt.Unix())
	return err
}

func (s *SQLiteStore) GetOrganizationRecord(id string) (*Organization, error) {
	query := `SELECT id, name, slug, created_at, updated_at FROM organizations WHERE id = ?`
	row := s.db.QueryRow(query, id)

	var org Organization
	var createdAt, updatedAt int64
	if err := row.Scan(&org.ID, &org.Name, &org.Slug, &createdAt, &updatedAt); err != nil {
		return nil, err
	}
	org.CreatedAt = time.Unix(createdAt, 0)
	org.UpdatedAt = time.Unix(updatedAt, 0)
	return &org, nil
}

func (s *SQLiteStore) CreateOrganizationMemberRecord(member *OrganizationMember) error {
	query := `
		INSERT INTO organization_members (id, organization_id, user_id, email, role, status, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	_, err := s.db.Exec(
		query,
		member.ID,
		member.OrganizationID,
		nullIfEmpty(member.UserID),
		member.Email,
		member.Role,
		member.Status,
		member.CreatedAt.Unix(),
	)
	return err
}

func (s *SQLiteStore) CreateSessionRecord(session *SessionRecord) error {
	query := `
		INSERT INTO sessions (id, user_id, workspace_id, plan, guest, expires_at, last_seen_at, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err := s.db.Exec(
		query,
		session.ID,
		session.UserID,
		nullIfEmpty(session.WorkspaceID),
		string(session.Plan),
		boolToInt(session.Guest),
		session.ExpiresAt.Unix(),
		nullIfZeroTime(session.LastSeenAt),
		session.CreatedAt.Unix(),
	)
	return err
}

func (s *SQLiteStore) GetSessionRecord(id string) (*SessionRecord, error) {
	query := `SELECT id, user_id, workspace_id, plan, guest, expires_at, last_seen_at, created_at FROM sessions WHERE id = ?`
	row := s.db.QueryRow(query, id)

	var session SessionRecord
	var workspaceID sql.NullString
	var guest int
	var expiresAt, lastSeenAt, createdAt sql.NullInt64
	var plan string
	if err := row.Scan(&session.ID, &session.UserID, &workspaceID, &plan, &guest, &expiresAt, &lastSeenAt, &createdAt); err != nil {
		return nil, err
	}

	if workspaceID.Valid {
		session.WorkspaceID = workspaceID.String
	}
	session.Plan = parsePlan(plan)
	session.Guest = guest == 1
	if expiresAt.Valid {
		session.ExpiresAt = time.Unix(expiresAt.Int64, 0)
	}
	if lastSeenAt.Valid {
		session.LastSeenAt = time.Unix(lastSeenAt.Int64, 0)
	}
	if createdAt.Valid {
		session.CreatedAt = time.Unix(createdAt.Int64, 0)
	}
	return &session, nil
}

func (s *SQLiteStore) DeleteSessionRecord(id string) error {
	_, err := s.db.Exec(`DELETE FROM sessions WHERE id = ?`, id)
	return err
}

func (s *SQLiteStore) TouchSessionRecord(id string, expiresAt, lastSeenAt time.Time) error {
	_, err := s.db.Exec(
		`UPDATE sessions SET expires_at = ?, last_seen_at = ? WHERE id = ?`,
		expiresAt.Unix(),
		lastSeenAt.Unix(),
		id,
	)
	return err
}

func (s *SQLiteStore) UpdateUserLastLogin(userID string, at time.Time) error {
	_, err := s.db.Exec(`UPDATE users SET last_login_at = ?, updated_at = ? WHERE id = ?`, at.Unix(), at.Unix(), userID)
	return err
}

func (s *SQLiteStore) UpsertSubscription(subscription *Subscription) error {
	query := `
		INSERT INTO subscriptions (
			id, workspace_id, organization_id, plan, status, provider, provider_customer_id,
			provider_subscription_id, current_period_end, created_at, updated_at
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			plan = excluded.plan,
			status = excluded.status,
			provider = excluded.provider,
			provider_customer_id = excluded.provider_customer_id,
			provider_subscription_id = excluded.provider_subscription_id,
			current_period_end = excluded.current_period_end,
			updated_at = excluded.updated_at
	`
	_, err := s.db.Exec(
		query,
		subscription.ID,
		nullIfEmpty(subscription.WorkspaceID),
		nullIfEmpty(subscription.OrganizationID),
		string(subscription.Plan),
		subscription.Status,
		nullIfEmpty(subscription.Provider),
		nullIfEmpty(subscription.ProviderCustomerID),
		nullIfEmpty(subscription.ProviderSubscriptionID),
		nullIfZeroTime(subscription.CurrentPeriodEnd),
		subscription.CreatedAt.Unix(),
		subscription.UpdatedAt.Unix(),
	)
	return err
}

func (s *SQLiteStore) GetSubscriptionForWorkspace(workspaceID string) (*Subscription, error) {
	query := `
		SELECT id, workspace_id, organization_id, plan, status, provider, provider_customer_id,
		       provider_subscription_id, current_period_end, created_at, updated_at
		FROM subscriptions
		WHERE workspace_id = ?
		ORDER BY updated_at DESC
		LIMIT 1
	`
	row := s.db.QueryRow(query, workspaceID)

	var subscription Subscription
	var orgID, provider, customerID, subscriptionID sql.NullString
	var currentPeriodEnd, createdAt, updatedAt sql.NullInt64
	var plan string
	if err := row.Scan(
		&subscription.ID,
		&subscription.WorkspaceID,
		&orgID,
		&plan,
		&subscription.Status,
		&provider,
		&customerID,
		&subscriptionID,
		&currentPeriodEnd,
		&createdAt,
		&updatedAt,
	); err != nil {
		return nil, err
	}

	if orgID.Valid {
		subscription.OrganizationID = orgID.String
	}
	if provider.Valid {
		subscription.Provider = provider.String
	}
	if customerID.Valid {
		subscription.ProviderCustomerID = customerID.String
	}
	if subscriptionID.Valid {
		subscription.ProviderSubscriptionID = subscriptionID.String
	}
	if currentPeriodEnd.Valid {
		subscription.CurrentPeriodEnd = time.Unix(currentPeriodEnd.Int64, 0)
	}
	if createdAt.Valid {
		subscription.CreatedAt = time.Unix(createdAt.Int64, 0)
	}
	if updatedAt.Valid {
		subscription.UpdatedAt = time.Unix(updatedAt.Int64, 0)
	}
	subscription.Plan = parsePlan(plan)
	return &subscription, nil
}

func (s *SQLiteStore) CreateSubscriptionEvent(event *SubscriptionEvent) error {
	query := `
		INSERT INTO subscription_events (id, subscription_id, event_type, provider_event_id, payload, created_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	_, err := s.db.Exec(
		query,
		event.ID,
		event.SubscriptionID,
		event.EventType,
		nullIfEmpty(event.ProviderEventID),
		nullIfEmpty(event.Payload),
		event.CreatedAt.Unix(),
	)
	return err
}

func (s *SQLiteStore) CreateDeckShareRecord(share *DeckShare) error {
	query := `
		INSERT INTO deck_shares (id, deck_id, workspace_id, created_by_user_id, token, access_type, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	_, err := s.db.Exec(
		query,
		share.ID,
		share.DeckID,
		nullIfEmpty(share.WorkspaceID),
		nullIfEmpty(share.CreatedByUserID),
		share.Token,
		share.AccessType,
		share.CreatedAt.Unix(),
	)
	return err
}

func (s *SQLiteStore) GetDeckShareByDeckID(deckID int64) (*DeckShare, error) {
	query := `
		SELECT id, deck_id, workspace_id, created_by_user_id, token, access_type, created_at
		FROM deck_shares WHERE deck_id = ?
		ORDER BY created_at DESC
		LIMIT 1
	`
	row := s.db.QueryRow(query, deckID)

	var share DeckShare
	var workspaceID, createdBy sql.NullString
	var createdAt int64
	if err := row.Scan(&share.ID, &share.DeckID, &workspaceID, &createdBy, &share.Token, &share.AccessType, &createdAt); err != nil {
		return nil, err
	}
	if workspaceID.Valid {
		share.WorkspaceID = workspaceID.String
	}
	if createdBy.Valid {
		share.CreatedByUserID = createdBy.String
	}
	share.CreatedAt = time.Unix(createdAt, 0)
	return &share, nil
}

func (s *SQLiteStore) DeleteDeckShareByDeckID(deckID int64) error {
	_, err := s.db.Exec(`DELETE FROM deck_shares WHERE deck_id = ?`, deckID)
	return err
}

func (s *SQLiteStore) CountDeckSharesForWorkspace(workspaceID string) (int, error) {
	var count int
	err := s.db.QueryRow(`SELECT COUNT(*) FROM deck_shares WHERE workspace_id = ?`, workspaceID).Scan(&count)
	return count, err
}

func (s *SQLiteStore) CountSyncDevicesForWorkspace(workspaceID string) (int, error) {
	var count int
	err := s.db.QueryRow(`SELECT COUNT(*) FROM sync_devices WHERE workspace_id = ?`, workspaceID).Scan(&count)
	return count, err
}

func (s *SQLiteStore) ListRecentDeckNotes(collectionID string, deckID int64, limit int, cursorCreatedAt int64, cursorNoteID int64) ([]RecentDeckNoteSummary, error) {
	if limit <= 0 {
		limit = 20
	}

	query := `
		SELECT
			n.id,
			n.type_id,
			n.field_vals,
			n.tags,
			n.created_at,
			n.modified_at,
			COUNT(c.id) as card_count
		FROM notes n
		INNER JOIN cards c ON c.note_id = n.id
		WHERE n.collection_id = ?
		  AND c.deck_id = ?
	`
	args := []interface{}{collectionID, deckID}

	if cursorCreatedAt > 0 {
		query += `
		  AND (
		    n.created_at < ?
		    OR (n.created_at = ? AND n.id < ?)
		  )
		`
		args = append(args, cursorCreatedAt, cursorCreatedAt, cursorNoteID)
	}

	query += `
		GROUP BY n.id, n.type_id, n.field_vals, n.tags, n.created_at, n.modified_at
		ORDER BY n.created_at DESC, n.id DESC
		LIMIT ?
	`
	args = append(args, limit)

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var summaries []RecentDeckNoteSummary
	for rows.Next() {
		var summary RecentDeckNoteSummary
		var fieldValsJSON, tagsJSON []byte
		var createdAt, modifiedAt int64

		if err := rows.Scan(
			&summary.NoteID,
			&summary.NoteType,
			&fieldValsJSON,
			&tagsJSON,
			&createdAt,
			&modifiedAt,
			&summary.CardCountInDeck,
		); err != nil {
			return nil, err
		}

		var fieldMap map[string]string
		if err := json.Unmarshal(fieldValsJSON, &fieldMap); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(tagsJSON, &summary.Tags); err != nil {
			return nil, err
		}

		summary.CreatedAt = time.Unix(createdAt, 0)
		summary.ModifiedAt = time.Unix(modifiedAt, 0)
		summary.FieldPreview = firstFieldPreview(fieldMap)
		summaries = append(summaries, summary)
	}

	return summaries, rows.Err()
}

func nullIfEmpty(value string) interface{} {
	if strings.TrimSpace(value) == "" {
		return nil
	}
	return value
}

func nullIfZeroTime(value time.Time) interface{} {
	if value.IsZero() {
		return nil
	}
	return value.Unix()
}

func boolToInt(value bool) int {
	if value {
		return 1
	}
	return 0
}

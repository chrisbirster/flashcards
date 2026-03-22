package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

func unixTimeOrZero(value sql.NullInt64) time.Time {
	if !value.Valid || value.Int64 == 0 {
		return time.Time{}
	}
	return time.Unix(value.Int64, 0)
}

func (s *SQLiteStore) ListWorkspacesForUser(userID string) ([]Workspace, error) {
	rows, err := s.db.Query(`
		SELECT id, name, slug, collection_id, owner_user_id, organization_id, created_at, updated_at
		FROM workspaces
		WHERE owner_user_id = ?
		ORDER BY created_at ASC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var workspaces []Workspace
	for rows.Next() {
		workspace, err := scanWorkspaceRow(rows)
		if err != nil {
			return nil, err
		}
		workspaces = append(workspaces, *workspace)
	}
	return workspaces, rows.Err()
}

func (s *SQLiteStore) GetWorkspaceForUser(userID, workspaceID string) (*Workspace, error) {
	row := s.db.QueryRow(`
		SELECT id, name, slug, collection_id, owner_user_id, organization_id, created_at, updated_at
		FROM workspaces
		WHERE id = ? AND owner_user_id = ?
	`, workspaceID, userID)
	return scanWorkspaceRow(row)
}

func scanWorkspaceRow(scanner interface{ Scan(dest ...any) error }) (*Workspace, error) {
	var workspace Workspace
	var ownerID, orgID sql.NullString
	var createdAt, updatedAt int64
	if err := scanner.Scan(
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

func (s *SQLiteStore) GetDeckCollectionID(deckID int64) (string, error) {
	var collectionID string
	if err := s.db.QueryRow(`SELECT collection_id FROM decks WHERE id = ?`, deckID).Scan(&collectionID); err != nil {
		return "", err
	}
	return collectionID, nil
}

func (s *SQLiteStore) GetDeckContentSummary(deckID int64) (noteCount, cardCount int, err error) {
	if err = s.db.QueryRow(`SELECT COUNT(DISTINCT note_id), COUNT(*) FROM cards WHERE deck_id = ?`, deckID).Scan(&noteCount, &cardCount); err != nil {
		return 0, 0, err
	}
	return noteCount, cardCount, nil
}

func (s *SQLiteStore) CreateStudyGroup(group *StudyGroup) error {
	_, err := s.db.Exec(`
		INSERT INTO study_groups (
			id, workspace_id, primary_deck_id, name, description, visibility, join_policy,
			created_by_user_id, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`,
		group.ID,
		group.WorkspaceID,
		group.PrimaryDeckID,
		group.Name,
		group.Description,
		group.Visibility,
		group.JoinPolicy,
		group.CreatedByUserID,
		group.CreatedAt.Unix(),
		group.UpdatedAt.Unix(),
	)
	return err
}

func (s *SQLiteStore) GetStudyGroup(id string) (*StudyGroup, error) {
	row := s.db.QueryRow(`
		SELECT id, workspace_id, primary_deck_id, name, description, visibility, join_policy, created_by_user_id, created_at, updated_at
		FROM study_groups
		WHERE id = ?
	`, id)
	return scanStudyGroup(row)
}

func scanStudyGroup(scanner interface{ Scan(dest ...any) error }) (*StudyGroup, error) {
	var group StudyGroup
	var createdAt, updatedAt int64
	if err := scanner.Scan(
		&group.ID,
		&group.WorkspaceID,
		&group.PrimaryDeckID,
		&group.Name,
		&group.Description,
		&group.Visibility,
		&group.JoinPolicy,
		&group.CreatedByUserID,
		&createdAt,
		&updatedAt,
	); err != nil {
		return nil, err
	}
	group.CreatedAt = time.Unix(createdAt, 0)
	group.UpdatedAt = time.Unix(updatedAt, 0)
	return &group, nil
}

func (s *SQLiteStore) UpdateStudyGroup(group *StudyGroup) error {
	_, err := s.db.Exec(`
		UPDATE study_groups
		SET name = ?, description = ?, visibility = ?, join_policy = ?, updated_at = ?
		WHERE id = ?
	`, group.Name, group.Description, group.Visibility, group.JoinPolicy, group.UpdatedAt.Unix(), group.ID)
	return err
}

func (s *SQLiteStore) DeleteStudyGroup(id string) error {
	_, err := s.db.Exec(`DELETE FROM study_groups WHERE id = ?`, id)
	return err
}

func (s *SQLiteStore) CreateStudyGroupMember(member *StudyGroupMember) error {
	_, err := s.db.Exec(`
		INSERT INTO study_group_members (
			id, study_group_id, user_id, email, role, status, invite_token, invite_expires_at, joined_at, removed_at, created_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`,
		member.ID,
		member.StudyGroupID,
		nullIfEmpty(member.UserID),
		member.Email,
		member.Role,
		member.Status,
		nullIfEmpty(member.InviteToken),
		nullIfZeroTime(member.InviteExpiresAt),
		nullIfZeroTime(member.JoinedAt),
		nullIfZeroTime(member.RemovedAt),
		member.CreatedAt.Unix(),
	)
	return err
}

func (s *SQLiteStore) UpsertStudyGroupInvitation(member *StudyGroupMember) error {
	_, err := s.db.Exec(`
		INSERT INTO study_group_members (
			id, study_group_id, user_id, email, role, status, invite_token, invite_expires_at, joined_at, removed_at, created_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(study_group_id, email) DO UPDATE SET
			role = excluded.role,
			status = excluded.status,
			invite_token = excluded.invite_token,
			invite_expires_at = excluded.invite_expires_at,
			removed_at = excluded.removed_at
	`,
		member.ID,
		member.StudyGroupID,
		nullIfEmpty(member.UserID),
		member.Email,
		member.Role,
		member.Status,
		nullIfEmpty(member.InviteToken),
		nullIfZeroTime(member.InviteExpiresAt),
		nullIfZeroTime(member.JoinedAt),
		nullIfZeroTime(member.RemovedAt),
		member.CreatedAt.Unix(),
	)
	return err
}

func (s *SQLiteStore) GetStudyGroupMember(id string) (*StudyGroupMember, error) {
	row := s.db.QueryRow(`
		SELECT id, study_group_id, user_id, email, role, status, invite_token, invite_expires_at, joined_at, removed_at, created_at
		FROM study_group_members
		WHERE id = ?
	`, id)
	return scanStudyGroupMember(row)
}

func (s *SQLiteStore) GetStudyGroupMemberByGroupAndEmail(groupID, email string) (*StudyGroupMember, error) {
	row := s.db.QueryRow(`
		SELECT id, study_group_id, user_id, email, role, status, invite_token, invite_expires_at, joined_at, removed_at, created_at
		FROM study_group_members
		WHERE study_group_id = ? AND lower(email) = lower(?)
	`, groupID, email)
	return scanStudyGroupMember(row)
}

func (s *SQLiteStore) GetStudyGroupMemberByUser(groupID, userID string) (*StudyGroupMember, error) {
	row := s.db.QueryRow(`
		SELECT id, study_group_id, user_id, email, role, status, invite_token, invite_expires_at, joined_at, removed_at, created_at
		FROM study_group_members
		WHERE study_group_id = ? AND user_id = ?
	`, groupID, userID)
	return scanStudyGroupMember(row)
}

func (s *SQLiteStore) GetStudyGroupMemberByInviteToken(token string) (*StudyGroupMember, error) {
	row := s.db.QueryRow(`
		SELECT id, study_group_id, user_id, email, role, status, invite_token, invite_expires_at, joined_at, removed_at, created_at
		FROM study_group_members
		WHERE invite_token = ?
	`, token)
	return scanStudyGroupMember(row)
}

func (s *SQLiteStore) ListStudyGroupMembers(groupID string) ([]StudyGroupMember, error) {
	rows, err := s.db.Query(`
		SELECT id, study_group_id, user_id, email, role, status, invite_token, invite_expires_at, joined_at, removed_at, created_at
		FROM study_group_members
		WHERE study_group_id = ?
		ORDER BY created_at ASC, lower(email)
	`, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []StudyGroupMember
	for rows.Next() {
		member, err := scanStudyGroupMember(rows)
		if err != nil {
			return nil, err
		}
		members = append(members, *member)
	}
	return members, rows.Err()
}

func scanStudyGroupMember(scanner interface{ Scan(dest ...any) error }) (*StudyGroupMember, error) {
	var member StudyGroupMember
	var userID, inviteToken sql.NullString
	var inviteExpiresAt, joinedAt, removedAt sql.NullInt64
	var createdAt int64
	if err := scanner.Scan(
		&member.ID,
		&member.StudyGroupID,
		&userID,
		&member.Email,
		&member.Role,
		&member.Status,
		&inviteToken,
		&inviteExpiresAt,
		&joinedAt,
		&removedAt,
		&createdAt,
	); err != nil {
		return nil, err
	}
	if userID.Valid {
		member.UserID = userID.String
	}
	if inviteToken.Valid {
		member.InviteToken = inviteToken.String
	}
	member.InviteExpiresAt = unixTimeOrZero(inviteExpiresAt)
	member.JoinedAt = unixTimeOrZero(joinedAt)
	member.RemovedAt = unixTimeOrZero(removedAt)
	member.CreatedAt = time.Unix(createdAt, 0)
	return &member, nil
}

func (s *SQLiteStore) UpdateStudyGroupMember(member *StudyGroupMember) error {
	_, err := s.db.Exec(`
		UPDATE study_group_members
		SET user_id = ?, role = ?, status = ?, invite_token = ?, invite_expires_at = ?, joined_at = ?, removed_at = ?
		WHERE id = ?
	`,
		nullIfEmpty(member.UserID),
		member.Role,
		member.Status,
		nullIfEmpty(member.InviteToken),
		nullIfZeroTime(member.InviteExpiresAt),
		nullIfZeroTime(member.JoinedAt),
		nullIfZeroTime(member.RemovedAt),
		member.ID,
	)
	return err
}

func (s *SQLiteStore) CreateStudyGroupVersion(version *StudyGroupVersion) error {
	_, err := s.db.Exec(`
		INSERT INTO study_group_versions (
			id, study_group_id, version_number, source_deck_id, published_by_user_id,
			change_summary, note_count, card_count, created_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`,
		version.ID,
		version.StudyGroupID,
		version.VersionNumber,
		version.SourceDeckID,
		version.PublishedByUserID,
		version.ChangeSummary,
		version.NoteCount,
		version.CardCount,
		version.CreatedAt.Unix(),
	)
	return err
}

func (s *SQLiteStore) ListStudyGroupVersions(groupID string) ([]StudyGroupVersion, error) {
	rows, err := s.db.Query(`
		SELECT id, study_group_id, version_number, source_deck_id, published_by_user_id,
		       change_summary, note_count, card_count, created_at
		FROM study_group_versions
		WHERE study_group_id = ?
		ORDER BY version_number DESC
	`, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var versions []StudyGroupVersion
	for rows.Next() {
		version, err := scanStudyGroupVersion(rows)
		if err != nil {
			return nil, err
		}
		versions = append(versions, *version)
	}
	return versions, rows.Err()
}

func (s *SQLiteStore) GetLatestStudyGroupVersion(groupID string) (*StudyGroupVersion, error) {
	row := s.db.QueryRow(`
		SELECT id, study_group_id, version_number, source_deck_id, published_by_user_id,
		       change_summary, note_count, card_count, created_at
		FROM study_group_versions
		WHERE study_group_id = ?
		ORDER BY version_number DESC
		LIMIT 1
	`, groupID)
	return scanStudyGroupVersion(row)
}

func scanStudyGroupVersion(scanner interface{ Scan(dest ...any) error }) (*StudyGroupVersion, error) {
	var version StudyGroupVersion
	var createdAt int64
	if err := scanner.Scan(
		&version.ID,
		&version.StudyGroupID,
		&version.VersionNumber,
		&version.SourceDeckID,
		&version.PublishedByUserID,
		&version.ChangeSummary,
		&version.NoteCount,
		&version.CardCount,
		&createdAt,
	); err != nil {
		return nil, err
	}
	version.CreatedAt = time.Unix(createdAt, 0)
	return &version, nil
}

func (s *SQLiteStore) CreateStudyGroupInstall(install *StudyGroupInstall) error {
	_, err := s.db.Exec(`
		INSERT INTO study_group_installs (
			id, study_group_id, study_group_member_id, destination_workspace_id,
			installed_deck_id, source_version_number, status, sync_state,
			superseded_by_install_id, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`,
		install.ID,
		install.StudyGroupID,
		install.StudyGroupMemberID,
		install.DestinationWorkspaceID,
		install.InstalledDeckID,
		install.SourceVersionNumber,
		install.Status,
		install.SyncState,
		nullIfEmpty(install.SupersededByInstallID),
		install.CreatedAt.Unix(),
		install.UpdatedAt.Unix(),
	)
	return err
}

func (s *SQLiteStore) GetStudyGroupInstall(id string) (*StudyGroupInstall, error) {
	row := s.db.QueryRow(`
		SELECT i.id, i.study_group_id, i.study_group_member_id, i.destination_workspace_id,
		       i.installed_deck_id, d.name, i.source_version_number, i.status, i.sync_state,
		       i.superseded_by_install_id, i.created_at, i.updated_at
		FROM study_group_installs i
		LEFT JOIN decks d ON d.id = i.installed_deck_id
		WHERE i.id = ?
	`, id)
	return scanStudyGroupInstall(row)
}

func (s *SQLiteStore) GetStudyGroupInstallByDeckID(deckID int64) (*StudyGroupInstall, error) {
	row := s.db.QueryRow(`
		SELECT i.id, i.study_group_id, i.study_group_member_id, i.destination_workspace_id,
		       i.installed_deck_id, d.name, i.source_version_number, i.status, i.sync_state,
		       i.superseded_by_install_id, i.created_at, i.updated_at
		FROM study_group_installs i
		LEFT JOIN decks d ON d.id = i.installed_deck_id
		WHERE i.installed_deck_id = ? AND i.status != 'removed'
		ORDER BY i.updated_at DESC
		LIMIT 1
	`, deckID)
	return scanStudyGroupInstall(row)
}

func (s *SQLiteStore) GetCurrentStudyGroupInstall(groupID, memberID string) (*StudyGroupInstall, error) {
	row := s.db.QueryRow(`
		SELECT i.id, i.study_group_id, i.study_group_member_id, i.destination_workspace_id,
		       i.installed_deck_id, d.name, i.source_version_number, i.status, i.sync_state,
		       i.superseded_by_install_id, i.created_at, i.updated_at
		FROM study_group_installs i
		LEFT JOIN decks d ON d.id = i.installed_deck_id
		WHERE i.study_group_id = ? AND i.study_group_member_id = ? AND i.status = 'active'
		ORDER BY i.updated_at DESC
		LIMIT 1
	`, groupID, memberID)
	return scanStudyGroupInstall(row)
}

func scanStudyGroupInstall(scanner interface{ Scan(dest ...any) error }) (*StudyGroupInstall, error) {
	var install StudyGroupInstall
	var installedDeckID sql.NullInt64
	var deckName, supersededBy sql.NullString
	var createdAt, updatedAt int64
	if err := scanner.Scan(
		&install.ID,
		&install.StudyGroupID,
		&install.StudyGroupMemberID,
		&install.DestinationWorkspaceID,
		&installedDeckID,
		&deckName,
		&install.SourceVersionNumber,
		&install.Status,
		&install.SyncState,
		&supersededBy,
		&createdAt,
		&updatedAt,
	); err != nil {
		return nil, err
	}
	if deckName.Valid {
		install.InstalledDeckName = deckName.String
	}
	if installedDeckID.Valid {
		install.InstalledDeckID = installedDeckID.Int64
	}
	if supersededBy.Valid {
		install.SupersededByInstallID = supersededBy.String
	}
	install.CreatedAt = time.Unix(createdAt, 0)
	install.UpdatedAt = time.Unix(updatedAt, 0)
	return &install, nil
}

func (s *SQLiteStore) UpdateStudyGroupInstall(install *StudyGroupInstall) error {
	_, err := s.db.Exec(`
		UPDATE study_group_installs
		SET destination_workspace_id = ?, installed_deck_id = ?, source_version_number = ?,
		    status = ?, sync_state = ?, superseded_by_install_id = ?, updated_at = ?
		WHERE id = ?
	`,
		install.DestinationWorkspaceID,
		nullableDeckID(install.InstalledDeckID),
		install.SourceVersionNumber,
		install.Status,
		install.SyncState,
		nullIfEmpty(install.SupersededByInstallID),
		install.UpdatedAt.Unix(),
		install.ID,
	)
	return err
}

func nullableDeckID(deckID int64) any {
	if deckID <= 0 {
		return nil
	}
	return deckID
}

func (s *SQLiteStore) MarkStudyGroupInstallForkedByDeckID(deckID int64) error {
	_, err := s.db.Exec(`
		UPDATE study_group_installs
		SET sync_state = 'forked', updated_at = ?
		WHERE installed_deck_id = ? AND status != 'removed' AND sync_state != 'forked'
	`, time.Now().Unix(), deckID)
	return err
}

func (s *SQLiteStore) MarkStudyGroupInstallsForkedByNoteType(collectionID string, noteTypeName string) error {
	_, err := s.db.Exec(`
		UPDATE study_group_installs
		SET sync_state = 'forked', updated_at = ?
		WHERE status != 'removed'
		  AND sync_state != 'forked'
		  AND installed_deck_id IN (
			SELECT DISTINCT c.deck_id
			FROM cards c
			JOIN notes n ON n.id = c.note_id
			JOIN decks d ON d.id = c.deck_id
			WHERE d.collection_id = ? AND n.type_id = ?
		  )
	`, time.Now().Unix(), collectionID, noteTypeRecordID(collectionID, NoteTypeName(noteTypeName)))
	return err
}

func (s *SQLiteStore) CreateStudyGroupEvent(event *StudyGroupEvent) error {
	_, err := s.db.Exec(`
		INSERT INTO study_group_events (id, study_group_id, actor_user_id, event_type, payload, created_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`,
		event.ID,
		event.StudyGroupID,
		nullIfEmpty(event.ActorUserID),
		event.EventType,
		event.Payload,
		event.CreatedAt.Unix(),
	)
	return err
}

func (s *SQLiteStore) ListStudyGroupEvents(groupID string, limit int) ([]StudyGroupEvent, error) {
	if limit <= 0 {
		limit = 20
	}
	rows, err := s.db.Query(`
		SELECT id, study_group_id, actor_user_id, event_type, payload, created_at
		FROM study_group_events
		WHERE study_group_id = ?
		ORDER BY created_at DESC
		LIMIT ?
	`, groupID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []StudyGroupEvent
	for rows.Next() {
		event, err := scanStudyGroupEvent(rows)
		if err != nil {
			return nil, err
		}
		events = append(events, *event)
	}
	return events, rows.Err()
}

func scanStudyGroupEvent(scanner interface{ Scan(dest ...any) error }) (*StudyGroupEvent, error) {
	var event StudyGroupEvent
	var actorUserID sql.NullString
	var createdAt int64
	if err := scanner.Scan(&event.ID, &event.StudyGroupID, &actorUserID, &event.EventType, &event.Payload, &createdAt); err != nil {
		return nil, err
	}
	if actorUserID.Valid {
		event.ActorUserID = actorUserID.String
	}
	event.CreatedAt = time.Unix(createdAt, 0)
	return &event, nil
}

func (s *SQLiteStore) GetStudyGroupDashboard(groupID string) (StudyGroupDashboard, error) {
	dashboard := StudyGroupDashboard{}

	if err := s.db.QueryRow(`
		SELECT COUNT(*)
		FROM study_group_members
		WHERE study_group_id = ? AND status = 'active'
	`, groupID).Scan(&dashboard.MemberCount); err != nil {
		return dashboard, err
	}

	since := time.Now().Add(-7 * 24 * time.Hour).Unix()
	if err := s.db.QueryRow(`
		SELECT COUNT(DISTINCT i.study_group_member_id)
		FROM revlog r
		JOIN cards c ON c.id = r.card_id
		JOIN study_group_installs i ON i.installed_deck_id = c.deck_id
		WHERE i.study_group_id = ? AND i.status != 'removed' AND r.reviewed_at >= ?
	`, groupID, since).Scan(&dashboard.ActiveMembers7D); err != nil {
		return dashboard, err
	}
	if err := s.db.QueryRow(`
		SELECT COUNT(*)
		FROM revlog r
		JOIN cards c ON c.id = r.card_id
		JOIN study_group_installs i ON i.installed_deck_id = c.deck_id
		WHERE i.study_group_id = ? AND i.status != 'removed' AND r.reviewed_at >= ?
	`, groupID, since).Scan(&dashboard.Reviews7D); err != nil {
		return dashboard, err
	}

	if latestVersion, err := s.GetLatestStudyGroupVersion(groupID); err == nil {
		dashboard.LatestVersionNumber = latestVersion.VersionNumber
		if err := s.db.QueryRow(`
			SELECT COUNT(*)
			FROM study_group_installs
			WHERE study_group_id = ? AND source_version_number = ? AND status = 'active'
		`, groupID, latestVersion.VersionNumber).Scan(&dashboard.LatestVersionAdoption); err != nil {
			return dashboard, err
		}
	}

	rows, err := s.db.Query(`
		SELECT m.user_id, m.email, COALESCE(u.display_name, ''), COUNT(*)
		FROM revlog r
		JOIN cards c ON c.id = r.card_id
		JOIN study_group_installs i ON i.installed_deck_id = c.deck_id
		JOIN study_group_members m ON m.id = i.study_group_member_id
		LEFT JOIN users u ON u.id = m.user_id
		WHERE i.study_group_id = ? AND i.status != 'removed' AND r.reviewed_at >= ?
		GROUP BY m.id, m.user_id, m.email, u.display_name
		ORDER BY COUNT(*) DESC, lower(m.email) ASC
		LIMIT 10
	`, groupID, since)
	if err != nil {
		return dashboard, err
	}
	defer rows.Close()

	for rows.Next() {
		var entry StudyGroupLeaderboardEntry
		if err := rows.Scan(&entry.UserID, &entry.Email, &entry.DisplayName, &entry.Reviews7D); err != nil {
			return dashboard, err
		}
		dashboard.Leaderboard = append(dashboard.Leaderboard, entry)
	}
	return dashboard, rows.Err()
}

func (s *SQLiteStore) ListStudyGroupSummariesForUser(userID string) ([]StudyGroupSummary, error) {
	user, err := s.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	rows, err := s.db.Query(`
		SELECT DISTINCT sg.id
		FROM study_groups sg
		JOIN study_group_members m ON m.study_group_id = sg.id
		WHERE m.user_id = ? OR lower(m.email) = lower(?)
		ORDER BY sg.created_at DESC
	`, userID, user.Email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var summaries []StudyGroupSummary
	for rows.Next() {
		var groupID string
		if err := rows.Scan(&groupID); err != nil {
			return nil, err
		}
		summary, err := s.BuildStudyGroupSummary(groupID, userID, user.Email)
		if err != nil {
			return nil, err
		}
		summaries = append(summaries, summary)
	}
	return summaries, rows.Err()
}

func (s *SQLiteStore) BuildStudyGroupSummary(groupID, userID, email string) (StudyGroupSummary, error) {
	group, err := s.GetStudyGroup(groupID)
	if err != nil {
		return StudyGroupSummary{}, err
	}

	summary := StudyGroupSummary{
		ID:           group.ID,
		Name:         group.Name,
		Description:  group.Description,
		SourceDeckID: group.PrimaryDeckID,
	}

	if deck, err := s.GetDeck(group.PrimaryDeckID); err == nil {
		summary.SourceDeckName = deck.Name
	}

	member, err := s.getStudyGroupMembership(group.ID, userID, email)
	if err == nil {
		summary.Role = member.Role
		summary.MembershipStatus = member.Status
		if install, err := s.GetCurrentStudyGroupInstall(group.ID, member.ID); err == nil {
			summary.CurrentUserInstall = install
		}
	}

	if latestVersion, err := s.GetLatestStudyGroupVersion(group.ID); err == nil {
		summary.LatestVersionNumber = latestVersion.VersionNumber
		if summary.CurrentUserInstall != nil && summary.CurrentUserInstall.SourceVersionNumber < latestVersion.VersionNumber {
			summary.UpdateAvailable = true
		}
	}

	dashboard, err := s.GetStudyGroupDashboard(group.ID)
	if err != nil {
		return StudyGroupSummary{}, err
	}
	summary.MemberCount = dashboard.MemberCount
	summary.ActiveMembers7D = dashboard.ActiveMembers7D

	return summary, nil
}

func (s *SQLiteStore) getStudyGroupMembership(groupID, userID, email string) (*StudyGroupMember, error) {
	if strings.TrimSpace(userID) != "" {
		member, err := s.GetStudyGroupMemberByUser(groupID, userID)
		if err == nil {
			return member, nil
		}
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
	}
	if strings.TrimSpace(email) == "" {
		return nil, sql.ErrNoRows
	}
	return s.GetStudyGroupMemberByGroupAndEmail(groupID, email)
}

func (s *SQLiteStore) BuildStudyGroupDetail(groupID, userID string) (*StudyGroupDetail, error) {
	user, err := s.GetUserByID(userID)
	if err != nil {
		return nil, err
	}
	group, err := s.GetStudyGroup(groupID)
	if err != nil {
		return nil, err
	}
	member, err := s.getStudyGroupMembership(groupID, userID, user.Email)
	if err != nil {
		return nil, err
	}

	detail := &StudyGroupDetail{
		Group:            *group,
		Role:             member.Role,
		MembershipStatus: member.Status,
	}

	if deck, err := s.GetDeck(group.PrimaryDeckID); err == nil {
		detail.SourceDeckName = deck.Name
	}
	if latestVersion, err := s.GetLatestStudyGroupVersion(groupID); err == nil {
		detail.LatestVersion = latestVersion
	}
	if versions, err := s.ListStudyGroupVersions(groupID); err == nil {
		detail.Versions = versions
	}
	if members, err := s.ListStudyGroupMembers(groupID); err == nil {
		detail.Members = members
	}
	if install, err := s.GetCurrentStudyGroupInstall(groupID, member.ID); err == nil {
		detail.CurrentUserInstall = install
		if detail.LatestVersion != nil && install.SourceVersionNumber < detail.LatestVersion.VersionNumber {
			detail.UpdateAvailable = true
		}
	}
	detail.CanEdit = member.Role == "owner" || member.Role == "admin"
	detail.CanInvite = detail.CanEdit
	detail.CanPublishVersion = detail.CanEdit
	if dashboard, err := s.GetStudyGroupDashboard(groupID); err == nil {
		detail.Dashboard = dashboard
	}
	if events, err := s.ListStudyGroupEvents(groupID, 20); err == nil {
		detail.RecentEvents = events
	}
	if workspaces, err := s.ListWorkspacesForUser(userID); err == nil {
		detail.AvailableWorkspaces = workspaces
	}
	return detail, nil
}

func (s *SQLiteStore) CopyDeckToCollection(sourceDeckID int64, destinationCollectionID, installedDeckName string) (*Deck, error) {
	sourceCollectionID, err := s.GetDeckCollectionID(sourceDeckID)
	if err != nil {
		return nil, err
	}

	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	nextID := func(table string) (int64, error) {
		var next int64
		query := fmt.Sprintf("SELECT COALESCE(MAX(id), 0) + 1 FROM %s", table)
		if err := tx.QueryRow(query).Scan(&next); err != nil {
			return 0, err
		}
		return next, nil
	}

	sourceDeck, err := s.GetDeck(sourceDeckID)
	if err != nil {
		return nil, err
	}
	newDeckID, err := nextID("decks")
	if err != nil {
		return nil, err
	}
	deckName := strings.TrimSpace(installedDeckName)
	if deckName == "" {
		deckName = sourceDeck.Name
	}
	if _, err = tx.Exec(`
		INSERT INTO decks (id, collection_id, name, parent_id, options_id)
		VALUES (?, ?, ?, ?, ?)
	`, newDeckID, destinationCollectionID, deckName, nil, sourceDeck.OptionsID); err != nil {
		return nil, err
	}

	rows, err := tx.Query(`
		SELECT c.id, c.note_id, c.template_name, c.ordinal, c.front, c.back,
		       c.due, c.state, c.fsrs_data, c.flag, c.marked, c.suspended, c.usn,
		       n.type_id, n.field_vals, n.tags, n.usn, n.created_at, n.modified_at
		FROM cards c
		JOIN notes n ON n.id = c.note_id
		WHERE c.deck_id = ?
		ORDER BY c.note_id ASC, c.id ASC
	`, sourceDeckID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	noteTypeEnsured := make(map[string]bool)
	noteIDMap := make(map[int64]int64)
	for rows.Next() {
		var (
			sourceCardID       int64
			sourceNoteID       int64
			templateName       string
			ordinal            int
			front              string
			back               string
			dueUnix            int64
			state              int
			fsrsData           []byte
			flag               int
			marked             int
			suspended          int
			cardUSN            int64
			noteTypeID         string
			fieldValsJSON      []byte
			tagsJSON           []byte
			noteUSN            int64
			noteCreatedAtUnix  int64
			noteModifiedAtUnix int64
		)
		if err := rows.Scan(
			&sourceCardID,
			&sourceNoteID,
			&templateName,
			&ordinal,
			&front,
			&back,
			&dueUnix,
			&state,
			&fsrsData,
			&flag,
			&marked,
			&suspended,
			&cardUSN,
			&noteTypeID,
			&fieldValsJSON,
			&tagsJSON,
			&noteUSN,
			&noteCreatedAtUnix,
			&noteModifiedAtUnix,
		); err != nil {
			return nil, err
		}

		noteTypeName := noteTypeNameFromRecordID(noteTypeID)
		if !noteTypeEnsured[string(noteTypeName)] {
			var (
				typeFields     []byte
				typeTemplates  []byte
				sortFieldIndex int
				fieldOptions   []byte
			)
			if err := tx.QueryRow(`
				SELECT fields, templates, sort_field_index, field_options
				FROM note_types
				WHERE collection_id = ? AND name = ?
			`, sourceCollectionID, string(noteTypeName)).Scan(&typeFields, &typeTemplates, &sortFieldIndex, &fieldOptions); err != nil {
				return nil, err
			}
			if _, err := tx.Exec(`
				INSERT INTO note_types (id, collection_id, name, fields, templates, sort_field_index, field_options)
				VALUES (?, ?, ?, ?, ?, ?, ?)
				ON CONFLICT(collection_id, name) DO NOTHING
			`, noteTypeRecordID(destinationCollectionID, noteTypeName), destinationCollectionID, string(noteTypeName), typeFields, typeTemplates, sortFieldIndex, fieldOptions); err != nil {
				return nil, err
			}
			noteTypeEnsured[string(noteTypeName)] = true
		}

		newNoteID, ok := noteIDMap[sourceNoteID]
		if !ok {
			newNoteID, err = nextID("notes")
			if err != nil {
				return nil, err
			}
			if _, err := tx.Exec(`
				INSERT INTO notes (id, collection_id, type_id, field_vals, tags, usn, created_at, modified_at)
				VALUES (?, ?, ?, ?, ?, ?, ?, ?)
			`, newNoteID, destinationCollectionID, noteTypeRecordID(destinationCollectionID, noteTypeName), fieldValsJSON, tagsJSON, noteUSN, noteCreatedAtUnix, noteModifiedAtUnix); err != nil {
				return nil, err
			}
			noteIDMap[sourceNoteID] = newNoteID
		}

		newCardID, err := nextID("cards")
		if err != nil {
			return nil, err
		}
		if _, err := tx.Exec(`
			INSERT INTO cards (
				id, note_id, deck_id, template_name, ordinal, front, back, due, state, fsrs_data, flag, marked, suspended, usn
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`, newCardID, newNoteID, newDeckID, templateName, ordinal, front, back, dueUnix, state, fsrsData, flag, marked, suspended, cardUSN); err != nil {
			return nil, err
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return s.GetDeck(newDeckID)
}

func (s *SQLiteStore) DeleteCopiedDeck(deckID int64) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	rows, err := tx.Query(`SELECT DISTINCT note_id FROM cards WHERE deck_id = ?`, deckID)
	if err != nil {
		return err
	}
	defer rows.Close()

	var noteIDs []int64
	for rows.Next() {
		var noteID int64
		if err := rows.Scan(&noteID); err != nil {
			return err
		}
		noteIDs = append(noteIDs, noteID)
	}
	if err := rows.Err(); err != nil {
		return err
	}

	if _, err := tx.Exec(`DELETE FROM cards WHERE deck_id = ?`, deckID); err != nil {
		return err
	}
	for _, noteID := range noteIDs {
		if _, err := tx.Exec(`DELETE FROM notes WHERE id = ?`, noteID); err != nil {
			return err
		}
	}
	if _, err := tx.Exec(`DELETE FROM decks WHERE id = ?`, deckID); err != nil {
		return err
	}
	return tx.Commit()
}

func (s *SQLiteStore) encodeStudyGroupEventPayload(payload any) string {
	if payload == nil {
		return ""
	}
	bytes, err := json.Marshal(payload)
	if err != nil {
		return ""
	}
	return string(bytes)
}

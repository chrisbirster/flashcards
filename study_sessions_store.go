package main

import (
	"database/sql"
	"time"
)

func (s *SQLiteStore) CreateStudySessionRecord(session *StudySession) error {
	_, err := s.db.Exec(`
		INSERT INTO study_sessions (
			id, user_id, workspace_id, deck_id, mode, protocol, target_minutes, break_minutes, status, started_at, ended_at,
			cards_reviewed, again_count, hard_count, good_count, easy_count, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`,
		session.ID,
		session.UserID,
		session.WorkspaceID,
		nullIfZeroInt64(session.DeckID),
		session.Mode,
		session.Protocol,
		session.TargetMinutes,
		session.BreakMinutes,
		session.Status,
		session.StartedAt.Unix(),
		nullIfZeroTime(session.EndedAt),
		session.CardsReviewed,
		session.AgainCount,
		session.HardCount,
		session.GoodCount,
		session.EasyCount,
		session.CreatedAt.Unix(),
		session.UpdatedAt.Unix(),
	)
	return err
}

func (s *SQLiteStore) GetStudySession(id string) (*StudySession, error) {
	row := s.db.QueryRow(`
		SELECT id, user_id, workspace_id, deck_id, mode, protocol, target_minutes, break_minutes, status, started_at, ended_at,
			cards_reviewed, again_count, hard_count, good_count, easy_count, created_at, updated_at
		FROM study_sessions
		WHERE id = ?
	`, id)
	return scanStudySession(row)
}

func (s *SQLiteStore) GetStudySessionForUser(id, userID string) (*StudySession, error) {
	row := s.db.QueryRow(`
		SELECT id, user_id, workspace_id, deck_id, mode, protocol, target_minutes, break_minutes, status, started_at, ended_at,
			cards_reviewed, again_count, hard_count, good_count, easy_count, created_at, updated_at
		FROM study_sessions
		WHERE id = ? AND user_id = ?
	`, id, userID)
	return scanStudySession(row)
}

func (s *SQLiteStore) UpdateStudySessionRecord(session *StudySession) error {
	_, err := s.db.Exec(`
		UPDATE study_sessions
		SET status = ?, ended_at = ?, cards_reviewed = ?, again_count = ?, hard_count = ?,
			good_count = ?, easy_count = ?, updated_at = ?
		WHERE id = ?
	`,
		session.Status,
		nullIfZeroTime(session.EndedAt),
		session.CardsReviewed,
		session.AgainCount,
		session.HardCount,
		session.GoodCount,
		session.EasyCount,
		session.UpdatedAt.Unix(),
		session.ID,
	)
	return err
}

func scanStudySession(scanner interface{ Scan(dest ...any) error }) (*StudySession, error) {
	var (
		session   StudySession
		deckID    sql.NullInt64
		endedAt   sql.NullInt64
		startedAt int64
		createdAt int64
		updatedAt int64
	)

	if err := scanner.Scan(
		&session.ID,
		&session.UserID,
		&session.WorkspaceID,
		&deckID,
		&session.Mode,
		&session.Protocol,
		&session.TargetMinutes,
		&session.BreakMinutes,
		&session.Status,
		&startedAt,
		&endedAt,
		&session.CardsReviewed,
		&session.AgainCount,
		&session.HardCount,
		&session.GoodCount,
		&session.EasyCount,
		&createdAt,
		&updatedAt,
	); err != nil {
		return nil, err
	}

	if deckID.Valid {
		session.DeckID = deckID.Int64
	}
	session.StartedAt = time.Unix(startedAt, 0)
	session.EndedAt = unixTimeOrZero(endedAt)
	session.CreatedAt = time.Unix(createdAt, 0)
	session.UpdatedAt = time.Unix(updatedAt, 0)
	return &session, nil
}

func nullIfZeroInt64(value int64) interface{} {
	if value == 0 {
		return nil
	}
	return value
}

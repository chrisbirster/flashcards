package main

import (
	"database/sql"
	"time"
)

func (s *SQLiteStore) CreateOTPChallenge(challenge *OTPChallenge) error {
	query := `
		INSERT INTO otp_challenges (
			id, email, code_hash, expires_at, attempt_count, max_attempts,
			resend_available_at, consumed_at, requested_ip, user_agent, created_at
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err := s.db.Exec(
		query,
		challenge.ID,
		challenge.Email,
		challenge.CodeHash,
		challenge.ExpiresAt.Unix(),
		challenge.AttemptCount,
		challenge.MaxAttempts,
		challenge.ResendAvailableAt.Unix(),
		nullIfZeroTime(challenge.ConsumedAt),
		nullIfEmpty(challenge.RequestedIP),
		nullIfEmpty(challenge.UserAgent),
		challenge.CreatedAt.Unix(),
	)
	return err
}

func (s *SQLiteStore) InvalidateOTPChallenges(email string) error {
	_, err := s.db.Exec(
		`UPDATE otp_challenges SET consumed_at = ? WHERE lower(email) = lower(?) AND consumed_at IS NULL`,
		time.Now().Unix(),
		email,
	)
	return err
}

func (s *SQLiteStore) GetLatestOTPChallenge(email string) (*OTPChallenge, error) {
	query := `
		SELECT id, email, code_hash, expires_at, attempt_count, max_attempts,
		       resend_available_at, consumed_at, requested_ip, user_agent, created_at
		FROM otp_challenges
		WHERE lower(email) = lower(?)
		ORDER BY created_at DESC
		LIMIT 1
	`
	row := s.db.QueryRow(query, email)

	var challenge OTPChallenge
	var consumedAt sql.NullInt64
	var requestedIP, userAgent sql.NullString
	var expiresAt, resendAvailableAt, createdAt int64
	if err := row.Scan(
		&challenge.ID,
		&challenge.Email,
		&challenge.CodeHash,
		&expiresAt,
		&challenge.AttemptCount,
		&challenge.MaxAttempts,
		&resendAvailableAt,
		&consumedAt,
		&requestedIP,
		&userAgent,
		&createdAt,
	); err != nil {
		return nil, err
	}

	challenge.ExpiresAt = time.Unix(expiresAt, 0)
	challenge.ResendAvailableAt = time.Unix(resendAvailableAt, 0)
	challenge.CreatedAt = time.Unix(createdAt, 0)
	if consumedAt.Valid {
		challenge.ConsumedAt = time.Unix(consumedAt.Int64, 0)
	}
	if requestedIP.Valid {
		challenge.RequestedIP = requestedIP.String
	}
	if userAgent.Valid {
		challenge.UserAgent = userAgent.String
	}

	return &challenge, nil
}

func (s *SQLiteStore) IncrementOTPChallengeAttempts(id string) error {
	_, err := s.db.Exec(`UPDATE otp_challenges SET attempt_count = attempt_count + 1 WHERE id = ?`, id)
	return err
}

func (s *SQLiteStore) ConsumeOTPChallenge(id string, consumedAt time.Time) error {
	_, err := s.db.Exec(`UPDATE otp_challenges SET consumed_at = ? WHERE id = ?`, consumedAt.Unix(), id)
	return err
}

func (s *SQLiteStore) CountRecentOTPChallengesByEmail(email string, since time.Time) (int, error) {
	var count int
	err := s.db.QueryRow(
		`SELECT COUNT(*) FROM otp_challenges WHERE lower(email) = lower(?) AND created_at >= ?`,
		email,
		since.Unix(),
	).Scan(&count)
	return count, err
}

func (s *SQLiteStore) CountRecentOTPChallengesByIP(ip string, since time.Time) (int, error) {
	if ip == "" {
		return 0, nil
	}
	var count int
	err := s.db.QueryRow(
		`SELECT COUNT(*) FROM otp_challenges WHERE requested_ip = ? AND created_at >= ?`,
		ip,
		since.Unix(),
	).Scan(&count)
	return count, err
}

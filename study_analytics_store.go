package main

import (
	"database/sql"
	"time"
)

func (s *SQLiteStore) GetStudyAnalyticsOverview(userID, workspaceID string) (StudyAnalyticsOverview, error) {
	overview := StudyAnalyticsOverview{}
	if userID == "" || workspaceID == "" {
		return overview, nil
	}

	now := time.Now().UTC()
	windowStart := startOfUTCDay(now).AddDate(0, 0, -6)
	since := windowStart.Unix()

	var (
		minutesStudiedSeconds int64
		lastStudiedAt         sql.NullInt64
	)
	if err := s.db.QueryRow(`
		SELECT
			COUNT(*),
			COALESCE(SUM(cards_reviewed), 0),
			COALESCE(SUM(CASE
				WHEN ended_at IS NOT NULL AND ended_at > started_at THEN ended_at - started_at
				WHEN updated_at > started_at THEN updated_at - started_at
				ELSE 0
			END), 0),
			COALESCE(SUM(again_count), 0),
			COALESCE(SUM(hard_count), 0),
			COALESCE(SUM(good_count), 0),
			COALESCE(SUM(easy_count), 0),
			MAX(COALESCE(ended_at, updated_at, started_at))
		FROM study_sessions
		WHERE user_id = ? AND workspace_id = ? AND cards_reviewed > 0 AND COALESCE(ended_at, updated_at, started_at) >= ?
	`, userID, workspaceID, since).Scan(
		&overview.Sessions7D,
		&overview.CardsReviewed7D,
		&minutesStudiedSeconds,
		&overview.AnswerBreakdown.Again,
		&overview.AnswerBreakdown.Hard,
		&overview.AnswerBreakdown.Good,
		&overview.AnswerBreakdown.Easy,
		&lastStudiedAt,
	); err != nil {
		return overview, err
	}

	overview.MinutesStudied7D = int(minutesStudiedSeconds / 60)

	var focusSecondsStudied int64
	if err := s.db.QueryRow(`
		SELECT
			COUNT(*),
			COALESCE(SUM(CASE
				WHEN ended_at IS NOT NULL AND ended_at > started_at THEN ended_at - started_at
				WHEN updated_at > started_at THEN updated_at - started_at
				ELSE 0
			END), 0)
		FROM study_sessions
		WHERE user_id = ? AND workspace_id = ? AND mode = 'focus' AND status = 'completed' AND COALESCE(ended_at, updated_at, started_at) >= ?
	`, userID, workspaceID, since).Scan(
		&overview.FocusSessions7D,
		&focusSecondsStudied,
	); err != nil {
		return overview, err
	}
	overview.FocusMinutes7D = int(focusSecondsStudied / 60)

	if err := s.db.QueryRow(`
		SELECT MAX(COALESCE(ended_at, updated_at, started_at))
		FROM study_sessions
		WHERE user_id = ? AND workspace_id = ? AND (
			cards_reviewed > 0 OR (mode = 'focus' AND status = 'completed')
		)
	`, userID, workspaceID).Scan(&lastStudiedAt); err != nil {
		return overview, err
	}
	overview.LastStudiedAt = unixTimeOrZero(lastStudiedAt)

	streak, err := s.currentStudyStreak(userID, workspaceID, now)
	if err != nil {
		return overview, err
	}
	overview.CurrentStreak = streak

	dailyActivity, err := s.getStudyAnalyticsDailyActivity(userID, workspaceID, windowStart, 7)
	if err != nil {
		return overview, err
	}
	overview.DailyActivity = dailyActivity

	recentSessions, err := s.getRecentStudySessionSummaries(userID, workspaceID, 5)
	if err != nil {
		return overview, err
	}
	overview.RecentSessions = recentSessions

	return overview, nil
}

func (s *SQLiteStore) GetDeckStudyAnalyticsSummary(userID, workspaceID string) (map[int64]DeckStudyAnalytics, error) {
	summaries := make(map[int64]DeckStudyAnalytics)
	if userID == "" || workspaceID == "" {
		return summaries, nil
	}

	since := startOfUTCDay(time.Now().UTC()).AddDate(0, 0, -6).Unix()
	rows, err := s.db.Query(`
		SELECT
			deck_id,
			COALESCE(SUM(CASE WHEN COALESCE(ended_at, updated_at, started_at) >= ? THEN 1 ELSE 0 END), 0) AS sessions_7d,
			COALESCE(SUM(CASE WHEN COALESCE(ended_at, updated_at, started_at) >= ? THEN cards_reviewed ELSE 0 END), 0) AS cards_reviewed_7d,
			COALESCE(SUM(CASE
				WHEN COALESCE(ended_at, updated_at, started_at) >= ? AND ended_at IS NOT NULL AND ended_at > started_at THEN ended_at - started_at
				WHEN COALESCE(ended_at, updated_at, started_at) >= ? AND updated_at > started_at THEN updated_at - started_at
				ELSE 0
			END), 0) AS seconds_studied_7d,
			COALESCE(SUM(CASE WHEN COALESCE(ended_at, updated_at, started_at) >= ? THEN again_count ELSE 0 END), 0) AS again_count_7d,
			COALESCE(SUM(CASE WHEN COALESCE(ended_at, updated_at, started_at) >= ? THEN hard_count ELSE 0 END), 0) AS hard_count_7d,
			COALESCE(SUM(CASE WHEN COALESCE(ended_at, updated_at, started_at) >= ? THEN good_count ELSE 0 END), 0) AS good_count_7d,
			COALESCE(SUM(CASE WHEN COALESCE(ended_at, updated_at, started_at) >= ? THEN easy_count ELSE 0 END), 0) AS easy_count_7d,
			MAX(COALESCE(ended_at, updated_at, started_at)) AS last_studied_at
		FROM study_sessions
		WHERE user_id = ? AND workspace_id = ? AND deck_id IS NOT NULL AND cards_reviewed > 0
		GROUP BY deck_id
	`, since, since, since, since, since, since, since, since, userID, workspaceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			deckID         int64
			summary        DeckStudyAnalytics
			secondsStudied int64
			lastStudiedAt  sql.NullInt64
		)
		if err := rows.Scan(
			&deckID,
			&summary.Sessions7D,
			&summary.CardsReviewed7D,
			&secondsStudied,
			&summary.AgainCount7D,
			&summary.HardCount7D,
			&summary.GoodCount7D,
			&summary.EasyCount7D,
			&lastStudiedAt,
		); err != nil {
			return nil, err
		}
		summary.MinutesStudied7D = int(secondsStudied / 60)
		if summary.Sessions7D > 0 {
			summary.AverageCardsPerSession7D = float64(summary.CardsReviewed7D) / float64(summary.Sessions7D)
		}
		summary.LastStudiedAt = unixTimeOrZero(lastStudiedAt)
		summaries[deckID] = summary
	}

	return summaries, rows.Err()
}

func (s *SQLiteStore) getStudyAnalyticsDailyActivity(userID, workspaceID string, windowStart time.Time, days int) ([]StudyAnalyticsDay, error) {
	rows, err := s.db.Query(`
		SELECT
			date(COALESCE(ended_at, updated_at, started_at), 'unixepoch') AS study_day,
			COUNT(*),
			COALESCE(SUM(cards_reviewed), 0),
			COALESCE(SUM(CASE
				WHEN ended_at IS NOT NULL AND ended_at > started_at THEN ended_at - started_at
				WHEN updated_at > started_at THEN updated_at - started_at
				ELSE 0
			END), 0)
		FROM study_sessions
		WHERE user_id = ? AND workspace_id = ? AND cards_reviewed > 0 AND COALESCE(ended_at, updated_at, started_at) >= ?
		GROUP BY study_day
	`, userID, workspaceID, windowStart.Unix())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	byDay := make(map[string]StudyAnalyticsDay, days)
	for rows.Next() {
		var (
			day            string
			summary        StudyAnalyticsDay
			secondsStudied int64
		)
		if err := rows.Scan(&day, &summary.Sessions, &summary.CardsReviewed, &secondsStudied); err != nil {
			return nil, err
		}
		summary.Date = day
		summary.MinutesStudied = int(secondsStudied / 60)
		byDay[day] = summary
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	dailyActivity := make([]StudyAnalyticsDay, 0, days)
	for offset := 0; offset < days; offset++ {
		currentDay := windowStart.AddDate(0, 0, offset).Format("2006-01-02")
		summary, ok := byDay[currentDay]
		if !ok {
			summary = StudyAnalyticsDay{Date: currentDay}
		}
		dailyActivity = append(dailyActivity, summary)
	}
	return dailyActivity, nil
}

func (s *SQLiteStore) getRecentStudySessionSummaries(userID, workspaceID string, limit int) ([]StudySessionSummary, error) {
	rows, err := s.db.Query(`
		SELECT
			id,
			deck_id,
			mode,
			protocol,
			target_minutes,
			break_minutes,
			status,
			started_at,
			ended_at,
			cards_reviewed,
			again_count,
			hard_count,
			good_count,
			easy_count,
			updated_at
		FROM study_sessions
		WHERE user_id = ? AND workspace_id = ? AND (cards_reviewed > 0 OR mode = 'focus')
		ORDER BY COALESCE(ended_at, updated_at, started_at) DESC
		LIMIT ?
	`, userID, workspaceID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	recent := make([]StudySessionSummary, 0, limit)
	for rows.Next() {
		var (
			summary   StudySessionSummary
			deckID    sql.NullInt64
			startedAt int64
			endedAt   sql.NullInt64
			updatedAt int64
		)
		if err := rows.Scan(
			&summary.ID,
			&deckID,
			&summary.Mode,
			&summary.Protocol,
			&summary.TargetMinutes,
			&summary.BreakMinutes,
			&summary.Status,
			&startedAt,
			&endedAt,
			&summary.CardsReviewed,
			&summary.AgainCount,
			&summary.HardCount,
			&summary.GoodCount,
			&summary.EasyCount,
			&updatedAt,
		); err != nil {
			return nil, err
		}
		if deckID.Valid {
			summary.DeckID = deckID.Int64
		}
		summary.StartedAt = time.Unix(startedAt, 0).UTC()
		summary.EndedAt = unixTimeOrZero(endedAt)
		summary.UpdatedAt = time.Unix(updatedAt, 0).UTC()
		summary.MinutesStudied = int(studySessionDurationSeconds(summary.StartedAt, summary.EndedAt, summary.UpdatedAt) / 60)
		recent = append(recent, summary)
	}

	return recent, rows.Err()
}

func (s *SQLiteStore) currentStudyStreak(userID, workspaceID string, now time.Time) (int, error) {
	rows, err := s.db.Query(`
		SELECT DISTINCT date(COALESCE(ended_at, updated_at, started_at), 'unixepoch') AS study_day
		FROM study_sessions
		WHERE user_id = ? AND workspace_id = ? AND (
			cards_reviewed > 0 OR (mode = 'focus' AND status = 'completed')
		)
		ORDER BY study_day DESC
	`, userID, workspaceID)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	var studyDays []string
	for rows.Next() {
		var studyDay string
		if err := rows.Scan(&studyDay); err != nil {
			return 0, err
		}
		studyDays = append(studyDays, studyDay)
	}
	if err := rows.Err(); err != nil {
		return 0, err
	}
	if len(studyDays) == 0 {
		return 0, nil
	}

	expected := now.Format("2006-01-02")
	streak := 0
	for _, studyDay := range studyDays {
		if studyDay != expected {
			break
		}
		streak++
		now = now.AddDate(0, 0, -1)
		expected = now.Format("2006-01-02")
	}
	return streak, nil
}

func startOfUTCDay(value time.Time) time.Time {
	utc := value.UTC()
	return time.Date(utc.Year(), utc.Month(), utc.Day(), 0, 0, 0, 0, time.UTC)
}

func studySessionDurationSeconds(startedAt, endedAt, updatedAt time.Time) int64 {
	if !endedAt.IsZero() && endedAt.After(startedAt) {
		return int64(endedAt.Sub(startedAt).Seconds())
	}
	if !updatedAt.IsZero() && updatedAt.After(startedAt) {
		return int64(updatedAt.Sub(startedAt).Seconds())
	}
	return 0
}

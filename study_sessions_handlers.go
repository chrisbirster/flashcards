package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
)

func validStudySessionStatus(status string) bool {
	switch status {
	case "active", "completed", "abandoned":
		return true
	default:
		return false
	}
}

func (h *APIHandler) CreateStudySession(w http.ResponseWriter, r *http.Request) {
	session := h.sessionFromRequest(r)
	if session == nil || strings.TrimSpace(session.UserID) == "" {
		respondAPIError(w, http.StatusUnauthorized, "study_session_unauthorized", "Authentication is required.")
		return
	}

	workspace, err := h.workspaceForSession(session)
	if err != nil {
		respondAPIError(w, http.StatusBadRequest, "workspace_not_found", "Workspace not found.")
		return
	}

	var req CreateStudySessionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondAPIError(w, http.StatusBadRequest, "invalid_request", "Invalid request body.")
		return
	}

	req.Mode = strings.TrimSpace(req.Mode)
	if req.Mode == "" {
		req.Mode = "review"
	}

	if req.DeckID != 0 {
		deckCollectionID, err := h.store.GetDeckCollectionID(req.DeckID)
		if err != nil {
			respondAPIError(w, http.StatusBadRequest, "deck_not_found", "Deck not found.")
			return
		}
		if deckCollectionID != workspace.CollectionID {
			respondAPIError(w, http.StatusBadRequest, "invalid_deck_workspace", "Deck must belong to the current workspace.")
			return
		}
	}

	now := time.Now()
	studySession := &StudySession{
		ID:          newID("sts"),
		UserID:      session.UserID,
		WorkspaceID: workspace.ID,
		DeckID:      req.DeckID,
		Mode:        req.Mode,
		Status:      "active",
		StartedAt:   now,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := h.store.CreateStudySessionRecord(studySession); err != nil {
		respondAPIError(w, http.StatusInternalServerError, "study_session_create_failed", err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, studySession)
}

func (h *APIHandler) UpdateStudySession(w http.ResponseWriter, r *http.Request) {
	session := h.sessionFromRequest(r)
	if session == nil || strings.TrimSpace(session.UserID) == "" {
		respondAPIError(w, http.StatusUnauthorized, "study_session_unauthorized", "Authentication is required.")
		return
	}

	studySessionID := strings.TrimSpace(chi.URLParam(r, "id"))
	if studySessionID == "" {
		respondAPIError(w, http.StatusBadRequest, "invalid_study_session", "Study session id is required.")
		return
	}

	studySession, err := h.store.GetStudySessionForUser(studySessionID, session.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			respondAPIError(w, http.StatusNotFound, "study_session_not_found", "Study session not found.")
			return
		}
		respondAPIError(w, http.StatusInternalServerError, "study_session_load_failed", err.Error())
		return
	}

	if studySession.Status != "active" {
		respondAPIError(w, http.StatusConflict, "study_session_closed", "Study session is already closed.")
		return
	}

	var req UpdateStudySessionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondAPIError(w, http.StatusBadRequest, "invalid_request", "Invalid request body.")
		return
	}

	if req.Status != "" && !validStudySessionStatus(req.Status) {
		respondAPIError(w, http.StatusBadRequest, "invalid_study_session_status", "Status must be active, completed, or abandoned.")
		return
	}

	for _, value := range []*int{req.CardsReviewed, req.AgainCount, req.HardCount, req.GoodCount, req.EasyCount} {
		if value != nil && *value < 0 {
			respondAPIError(w, http.StatusBadRequest, "invalid_study_session_counts", "Study session counts must be zero or greater.")
			return
		}
	}

	if req.CardsReviewed != nil {
		studySession.CardsReviewed = *req.CardsReviewed
	}
	if req.AgainCount != nil {
		studySession.AgainCount = *req.AgainCount
	}
	if req.HardCount != nil {
		studySession.HardCount = *req.HardCount
	}
	if req.GoodCount != nil {
		studySession.GoodCount = *req.GoodCount
	}
	if req.EasyCount != nil {
		studySession.EasyCount = *req.EasyCount
	}

	if req.Status != "" {
		studySession.Status = req.Status
		if req.Status != "active" {
			if !req.EndedAt.IsZero() {
				studySession.EndedAt = req.EndedAt
			} else {
				studySession.EndedAt = time.Now()
			}
		}
	}

	studySession.UpdatedAt = time.Now()
	if err := h.store.UpdateStudySessionRecord(studySession); err != nil {
		respondAPIError(w, http.StatusInternalServerError, "study_session_update_failed", err.Error())
		return
	}

	reloaded, err := h.store.GetStudySessionForUser(studySession.ID, session.UserID)
	if err != nil {
		respondAPIError(w, http.StatusInternalServerError, "study_session_reload_failed", err.Error())
		return
	}

	respondJSON(w, http.StatusOK, reloaded)
}

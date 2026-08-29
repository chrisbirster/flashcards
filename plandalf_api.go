package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

type plandalfAPI struct {
	store *PlandalfStore
	cfg   PlandalfWebConfig
}

type plandalfCreateDeckRequest struct {
	Name string `json:"name"`
}

type plandalfCreateCardRequest struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

type plandalfReviewRequest struct {
	Rating              int `json:"rating"`
	ExpectedReviewCount int `json:"expected_review_count"`
}

type plandalfCardResponse struct {
	ID          string `json:"id"`
	DeckID      string `json:"deck_id"`
	Question    string `json:"question"`
	Answer      string `json:"answer"`
	DueAtMs     *int64 `json:"due_at_ms"`
	ReviewCount int    `json:"review_count"`
}

type plandalfStudyNextResponse struct {
	Card     *plandalfCardResponse `json:"card"`
	Schedule *plandalfSchedule     `json:"schedule,omitempty"`
}

func NewPlandalfServer(cfg PlandalfWebConfig, store *PlandalfStore, frontend fs.FS) http.Handler {
	api := &plandalfAPI{store: store, cfg: cfg}
	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: cfg.AllowedOrigins,
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type"},
		MaxAge: 300,
	}))

	router.Get("/api/v1/health", api.health)
	router.Route("/api/v1", func(r chi.Router) {
		r.Use(api.auth)
		r.Get("/decks", api.listDecks)
		r.Post("/decks", api.createDeck)
		r.Post("/decks/{deckID}/cards", api.createCard)
		r.Get("/decks/{deckID}/study/next", api.studyNext)
		r.Get("/cards/{cardID}", api.getCard)
		r.Get("/cards/{cardID}/study/preview", api.studyPreview)
		r.Post("/cards/{cardID}/reviews", api.review)
	})

	fileServer := http.FileServer(http.FS(frontend))
	router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/")
		if path != "" {
			if file, err := frontend.Open(path); err == nil {
				_ = file.Close()
				fileServer.ServeHTTP(w, r)
				return
			}
		}
		index, err := fs.ReadFile(frontend, "index.html")
		if err != nil {
			http.Error(w, "web app not built", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write(index)
	})
	return router
}

func (a *plandalfAPI) auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if a.cfg.APIToken == "" {
			next.ServeHTTP(w, r)
			return
		}
		const prefix = "Bearer "
		header := r.Header.Get("Authorization")
		if !strings.HasPrefix(header, prefix) || strings.TrimSpace(strings.TrimPrefix(header, prefix)) != a.cfg.APIToken {
			writePlandalfError(w, http.StatusUnauthorized, "unauthorized", "A valid Plandalf API token is required")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func writePlandalfJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func writePlandalfError(w http.ResponseWriter, status int, code, message string) {
	writePlandalfJSON(w, status, map[string]any{"error": map[string]string{"code": code, "message": message}})
}

func decodePlandalfJSON(w http.ResponseWriter, r *http.Request, target any) bool {
	decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		writePlandalfError(w, http.StatusBadRequest, "invalid_json", "Request body must be valid JSON")
		return false
	}
	return true
}

func routeInt64(r *http.Request, name string) (int64, error) {
	value := chi.URLParam(r, name)
	parsed, err := strconv.ParseInt(value, 10, 64)
	if err != nil || parsed <= 0 { return 0, errors.New("invalid id") }
	return parsed, nil
}

func toPlandalfCardResponse(card *PlandalfCardRecord) *plandalfCardResponse {
	if card == nil { return nil }
	return &plandalfCardResponse{
		ID: fmt.Sprintf("%d", card.ID),
		DeckID: fmt.Sprintf("%d", card.DeckID),
		Question: card.Question,
		Answer: card.Answer,
		DueAtMs: card.DueAtMs,
		ReviewCount: card.ReviewCount,
	}
}

func (a *plandalfAPI) health(w http.ResponseWriter, r *http.Request) {
	writePlandalfJSON(w, http.StatusOK, map[string]string{"status": "ok", "scheduler": "fsrs/7"})
}

func (a *plandalfAPI) listDecks(w http.ResponseWriter, r *http.Request) {
	decks, err := a.store.ListDecks(time.Now().UnixMilli())
	if err != nil { writePlandalfError(w, 500, "storage_error", err.Error()); return }
	if decks == nil { decks = []PlandalfDeckSummary{} }
	writePlandalfJSON(w, 200, map[string]any{"decks": decks})
}

func (a *plandalfAPI) createDeck(w http.ResponseWriter, r *http.Request) {
	var input plandalfCreateDeckRequest
	if !decodePlandalfJSON(w, r, &input) { return }
	input.Name = strings.TrimSpace(input.Name)
	if input.Name == "" { writePlandalfError(w, 422, "invalid_deck", "Deck name is required"); return }
	id, err := a.store.CreateDeck(input.Name, time.Now().UnixMilli())
	if err != nil { writePlandalfError(w, 409, "deck_create_failed", err.Error()); return }
	writePlandalfJSON(w, 201, map[string]string{"id": fmt.Sprintf("%d", id), "name": input.Name})
}

func (a *plandalfAPI) createCard(w http.ResponseWriter, r *http.Request) {
	deckID, err := routeInt64(r, "deckID")
	if err != nil { writePlandalfError(w, 400, "invalid_id", "deckID must be a positive integer"); return }
	var input plandalfCreateCardRequest
	if !decodePlandalfJSON(w, r, &input) { return }
	input.Question = strings.TrimSpace(input.Question)
	input.Answer = strings.TrimSpace(input.Answer)
	if input.Question == "" || input.Answer == "" { writePlandalfError(w, 422, "invalid_card", "Question and answer are required"); return }
	id, err := a.store.CreateCard(deckID, input.Question, input.Answer, time.Now().UnixMilli())
	if err != nil { writePlandalfError(w, 409, "card_create_failed", err.Error()); return }
	writePlandalfJSON(w, 201, map[string]string{"id": fmt.Sprintf("%d", id), "deck_id": fmt.Sprintf("%d", deckID)})
}

func (a *plandalfAPI) getCard(w http.ResponseWriter, r *http.Request) {
	cardID, err := routeInt64(r, "cardID")
	if err != nil { writePlandalfError(w, 400, "invalid_id", "cardID must be a positive integer"); return }
	card, err := a.store.GetCard(cardID)
	if err != nil { writePlandalfError(w, 500, "storage_error", err.Error()); return }
	if card == nil { writePlandalfError(w, 404, "card_not_found", "Card not found"); return }
	writePlandalfJSON(w, 200, map[string]any{"card": toPlandalfCardResponse(card)})
}

func (a *plandalfAPI) studyNext(w http.ResponseWriter, r *http.Request) {
	deckID, err := routeInt64(r, "deckID")
	if err != nil { writePlandalfError(w, 400, "invalid_id", "deckID must be a positive integer"); return }
	nowMs := time.Now().UnixMilli()
	card, err := a.store.NextCard(deckID, nowMs)
	if err != nil { writePlandalfError(w, 500, "storage_error", err.Error()); return }
	if card == nil { writePlandalfJSON(w, 200, plandalfStudyNextResponse{Card: nil}); return }
	schedule, _, err := a.store.Preview(card.ID, nowMs)
	if err != nil { writePlandalfError(w, 500, "schedule_error", err.Error()); return }
	writePlandalfJSON(w, 200, plandalfStudyNextResponse{Card: toPlandalfCardResponse(card), Schedule: &schedule})
}

func (a *plandalfAPI) studyPreview(w http.ResponseWriter, r *http.Request) {
	cardID, err := routeInt64(r, "cardID")
	if err != nil { writePlandalfError(w, 400, "invalid_id", "cardID must be a positive integer"); return }
	schedule, reviewCount, err := a.store.Preview(cardID, time.Now().UnixMilli())
	if errors.Is(err, sql.ErrNoRows) { writePlandalfError(w, 404, "card_not_found", "Card not found"); return }
	if err != nil { writePlandalfError(w, 500, "schedule_error", err.Error()); return }
	writePlandalfJSON(w, 200, map[string]any{"card_id": fmt.Sprintf("%d", cardID), "review_count": reviewCount, "schedule": schedule})
}

func (a *plandalfAPI) review(w http.ResponseWriter, r *http.Request) {
	cardID, err := routeInt64(r, "cardID")
	if err != nil { writePlandalfError(w, 400, "invalid_id", "cardID must be a positive integer"); return }
	var input plandalfReviewRequest
	if !decodePlandalfJSON(w, r, &input) { return }
	rating, err := parsePlandalfRating(input.Rating)
	if err != nil { writePlandalfError(w, 400, "invalid_rating", err.Error()); return }
	reviewID, candidate, state, err := a.store.RecordReview(cardID, rating, input.ExpectedReviewCount, time.Now().UnixMilli())
	if err != nil {
		switch err.Error() {
		case "stale_review": writePlandalfError(w, 409, "stale_review", "Card history changed; refresh before rating")
		case "card_not_due": writePlandalfError(w, 409, "card_not_due", "Card is not due yet")
		default:
			if errors.Is(err, sql.ErrNoRows) { writePlandalfError(w, 404, "card_not_found", "Card not found") } else { writePlandalfError(w, 500, "review_failed", err.Error()) }
		}
		return
	}
	writePlandalfJSON(w, 201, map[string]any{
		"review_id": fmt.Sprintf("%d", reviewID),
		"card_id": fmt.Sprintf("%d", cardID),
		"rating": int(rating),
		"due_at_ms": candidate.DueAtMs,
		"interval_days": candidate.IntervalDays,
		"scheduler": map[string]any{
			"stability_days": state.StabilityDays,
			"difficulty": state.Difficulty,
			"due_at_ms": state.DueAtMs,
			"last_reviewed_at_ms": state.LastReviewedAtMs,
		},
	})
}

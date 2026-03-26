package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/microcosm-cc/bluemonday"
	fsrs "github.com/open-spaced-repetition/go-fsrs/v3"
)

var htmlPolicy = bluemonday.UGCPolicy()

func sanitizeHTML(input string) string {
	return htmlPolicy.Sanitize(input)
}

// APIHandler wraps the store and provides HTTP handlers
type APIHandler struct {
	store         *SQLiteStore
	collectionID  string
	collection    *Collection
	backupManager *BackupManager
	config        AppConfig
	emailSender   EmailSender
}

func NewAPIHandler(store *SQLiteStore, collection *Collection, backupMgr *BackupManager) *APIHandler {
	cfg, err := LoadAppConfig()
	if err != nil {
		cfg = mustLocalAppConfig()
	}
	return NewAPIHandlerWithConfig(store, collection, backupMgr, cfg, NewEmailSender(cfg))
}

func NewAPIHandlerWithConfig(store *SQLiteStore, collection *Collection, backupMgr *BackupManager, cfg AppConfig, emailSender EmailSender) *APIHandler {
	return &APIHandler{
		store:         store,
		collectionID:  "default",
		collection:    collection,
		backupManager: backupMgr,
		config:        cfg,
		emailSender:   emailSender,
	}
}

// Request/Response types
type CreateDeckRequest struct {
	Name string `json:"name"`
}

type DeckResponse struct {
	ID                  int64              `json:"id"`
	Name                string             `json:"name"`
	ParentID            *int64             `json:"parentId,omitempty"`
	CardIDs             []int64            `json:"cardIds"`
	DueToday            int                `json:"dueToday"`
	DueReviewBacklog    int                `json:"dueReviewBacklog"`
	NewCardsPerDay      int                `json:"newCardsPerDay"`
	ReviewsPerDay       int                `json:"reviewsPerDay"`
	PriorityOrder       int                `json:"priorityOrder"`
	NewCardsPaused      bool               `json:"newCardsPaused"`
	NoteCount           int                `json:"noteCount"`
	CardCount           int                `json:"cardCount"`
	CanDelete           bool               `json:"canDelete"`
	DeleteBlockedReason string             `json:"deleteBlockedReason,omitempty"`
	Analytics           DeckStudyAnalytics `json:"analytics"`
}

type DashboardResponse struct {
	TotalDecks     int                    `json:"totalDecks"`
	TotalNotes     int                    `json:"totalNotes"`
	DueToday       int                    `json:"dueToday"`
	Plan           Plan                   `json:"plan"`
	Usage          EntitlementUsage       `json:"usage"`
	Limits         PlanLimits             `json:"limits"`
	Features       EntitlementFeatures    `json:"features"`
	StudyAnalytics StudyAnalyticsOverview `json:"studyAnalytics"`
	RecentNotes    []NoteListItemResponse `json:"recentNotes"`
}

type CreateNoteRequest struct {
	TypeID         string            `json:"typeId"`
	DeckID         int64             `json:"deckId"`
	FieldVals      map[string]string `json:"fieldVals"`
	Tags           []string          `json:"tags"`
	AllowDuplicate bool              `json:"allowDuplicate"` // Override duplicate check
}

type CheckDuplicateRequest struct {
	TypeID    string `json:"typeId"`
	FieldName string `json:"fieldName"` // Field to check for duplicates (usually "Front" or first field)
	Value     string `json:"value"`
	DeckID    int64  `json:"deckId,omitempty"` // Optional: limit scope to deck
}

type DuplicateResult struct {
	IsDuplicate bool        `json:"isDuplicate"`
	Duplicates  []NoteBrief `json:"duplicates,omitempty"`
}

type NoteBrief struct {
	ID       int64             `json:"id"`
	TypeID   string            `json:"typeId"`
	FieldVal map[string]string `json:"fieldVals"`
	DeckID   int64             `json:"deckId,omitempty"`
}

type AnswerCardRequest struct {
	Rating      int `json:"rating"`      // 1=Again, 2=Hard, 3=Good, 4=Easy
	TimeTakenMs int `json:"timeTakenMs"` // Time spent on the card in milliseconds
}

type UpdateCardRequest struct {
	Flag      *int  `json:"flag,omitempty"`      // 0-7 color flags
	Marked    *bool `json:"marked,omitempty"`    // toggle marked status
	Suspended *bool `json:"suspended,omitempty"` // toggle suspended status
}

type ImportNotesJSONRequest struct {
	Content  string `json:"content"`
	Filename string `json:"filename"`
	Source   string `json:"source,omitempty"`
	Format   string `json:"format,omitempty"`
	DeckName string `json:"deckName,omitempty"`
	NoteType string `json:"noteType,omitempty"`
}

type ImportNotesResponse struct {
	Imported     int      `json:"imported"`
	Skipped      int      `json:"skipped"`
	Source       string   `json:"source"`
	Format       string   `json:"format"`
	DecksCreated []string `json:"decksCreated,omitempty"`
	Errors       []string `json:"errors,omitempty"`
}

// Handler methods

func (h *APIHandler) GetCollection(w http.ResponseWriter, r *http.Request) {
	col, _, err := h.collectionForRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusOK, col)
}

func (h *APIHandler) ListDecks(w http.ResponseWriter, r *http.Request) {
	col, collectionID, err := h.collectionForRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	decks, err := h.store.ListDecks(collectionID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	userID := h.userIDFromRequest(r)
	session := h.sessionFromRequest(r)
	if err := h.store.EnsureReviewStatesForUser(userID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	workspaceID := ""
	if session != nil {
		workspaceID = session.WorkspaceID
	}
	analyticsByDeck, err := h.store.GetDeckStudyAnalyticsSummary(userID, workspaceID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert to response format
	var response []DeckResponse
	for _, d := range decks {
		response = append(response, h.deckResponse(userID, d, col, analyticsByDeck))
	}

	respondJSON(w, http.StatusOK, response)
}

func (h *APIHandler) CreateDeck(w http.ResponseWriter, r *http.Request) {
	if !h.requireWorkspaceWritePermission(w, r) {
		return
	}
	col, collectionID, err := h.collectionForRequest(r)
	if err != nil {
		respondAPIError(w, http.StatusInternalServerError, "collection_load_failed", err.Error())
		return
	}

	var req CreateDeckRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondAPIError(w, http.StatusBadRequest, "invalid_request", "Invalid request body")
		return
	}

	if req.Name == "" {
		respondAPIError(w, http.StatusBadRequest, "invalid_name", "Deck name is required")
		return
	}

	session := h.sessionFromRequest(r)
	plan := h.planForRequest(r, session)
	usage := h.usageForSession(session)
	if err := validateDeckLimit(plan, usage); err != nil {
		respondAPIError(w, http.StatusForbidden, "plan_limit_exceeded", err.Error())
		return
	}

	// Sanitize deck name to prevent XSS
	sanitizedName := sanitizeHTML(req.Name)

	// Create deck using collection method
	deck := col.NewDeck(sanitizedName)

	// Persist to database
	if err := h.store.CreateDeckInCollection(collectionID, deck); err != nil {
		respondAPIError(w, http.StatusInternalServerError, "deck_create_failed", err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, h.deckResponse(h.userIDFromRequest(r), deck, col, nil))
}

func (h *APIHandler) GetDeck(w http.ResponseWriter, r *http.Request) {
	col, _, err := h.collectionForRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, err := parseIDParam(r, "id")
	if err != nil {
		http.Error(w, "Invalid deck ID", http.StatusBadRequest)
		return
	}

	deck, err := h.store.GetDeck(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	session := h.sessionFromRequest(r)
	workspaceID := ""
	if session != nil {
		workspaceID = session.WorkspaceID
	}
	analyticsByDeck, err := h.store.GetDeckStudyAnalyticsSummary(h.userIDFromRequest(r), workspaceID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get deck stats
	stats, err := h.store.GetDeckStatsForUser(h.userIDFromRequest(r), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"deck":  h.deckResponse(h.userIDFromRequest(r), deck, col, analyticsByDeck),
		"stats": stats,
	})
}

func (h *APIHandler) buildNoteCardIndex(col *Collection) map[int64][]Card {
	index := make(map[int64][]Card, len(col.Notes))
	for _, card := range col.Cards {
		index[card.NoteID] = append(index[card.NoteID], *card)
	}
	return index
}

func (h *APIHandler) userIDFromRequest(r *http.Request) string {
	session := h.sessionFromRequest(r)
	if session == nil {
		return ""
	}
	return strings.TrimSpace(session.UserID)
}

func (h *APIHandler) deckDeleteBlockedReason(deck *Deck, cardCount int, col *Collection) string {
	if cardCount > 0 {
		return "Only empty decks can be deleted right now. Move or delete the cards in this deck first."
	}
	for _, candidate := range col.Decks {
		if candidate.ParentID != nil && *candidate.ParentID == deck.ID {
			return "This deck has child decks. Move or delete those child decks first."
		}
	}
	return ""
}

func (h *APIHandler) deckResponse(userID string, deck *Deck, col *Collection, analyticsByDeck map[int64]DeckStudyAnalytics) DeckResponse {
	dueToday := 0
	dueReviewBacklog := 0
	newCardsPerDay := defaultNewCardsPerDay
	reviewsPerDay := defaultReviewsPerDay
	noteIDs := make(map[int64]struct{})
	cardCount := 0

	for _, cardID := range deck.Cards {
		card, ok := col.Cards[cardID]
		if !ok {
			continue
		}
		cardCount++
		noteIDs[card.NoteID] = struct{}{}
	}

	if stats, err := h.store.GetDeckStatsForUser(userID, deck.ID); err == nil {
		dueToday = stats.DueToday
		dueReviewBacklog = stats.DueReviewBacklog
	}
	if configuredNew, configuredReview, err := h.store.getDeckDailyLimits(deck.ID); err == nil {
		newCardsPerDay = configuredNew
		reviewsPerDay = configuredReview
	}

	deleteBlockedReason := h.deckDeleteBlockedReason(deck, cardCount, col)
	analytics := DeckStudyAnalytics{}
	if analyticsByDeck != nil {
		analytics = analyticsByDeck[deck.ID]
	}

	return DeckResponse{
		ID:                  deck.ID,
		Name:                deck.Name,
		ParentID:            deck.ParentID,
		CardIDs:             deck.Cards,
		DueToday:            dueToday,
		DueReviewBacklog:    dueReviewBacklog,
		NewCardsPerDay:      newCardsPerDay,
		ReviewsPerDay:       reviewsPerDay,
		PriorityOrder:       deck.PriorityOrder,
		NewCardsPaused:      dueReviewBacklog > reviewsPerDay,
		NoteCount:           len(noteIDs),
		CardCount:           cardCount,
		CanDelete:           deleteBlockedReason == "",
		DeleteBlockedReason: deleteBlockedReason,
		Analytics:           analytics,
	}
}

func (h *APIHandler) buildDashboardResponse(r *http.Request) DashboardResponse {
	col, _, err := h.collectionForRequest(r)
	if err != nil {
		col = h.collection
	}
	studyAnalytics := StudyAnalyticsOverview{}
	if session := h.sessionFromRequest(r); session != nil {
		if analytics, err := h.store.GetStudyAnalyticsOverview(session.UserID, session.WorkspaceID); err == nil {
			studyAnalytics = analytics
		}
	}
	sessionResponse := h.buildSessionResponse(r)
	noteCards := h.buildNoteCardIndex(col)
	recentNotes := make([]NoteListItemResponse, 0, len(col.Notes))
	dueToday := 0
	if count, err := h.store.CountDueCardsForUser(h.userIDFromRequest(r)); err == nil {
		dueToday = count
	}

	for _, note := range col.Notes {
		cards := noteCards[note.ID]
		primaryDeckID, primaryDeckName := h.primaryDeckDetails(cards, col)
		recentNotes = append(recentNotes, NoteListItemResponse{
			ID:           note.ID,
			TypeID:       string(note.Type),
			FieldVals:    note.FieldMap,
			FieldPreview: h.noteFieldPreview(note, col),
			Tags:         note.Tags,
			CreatedAt:    note.CreatedAt,
			ModifiedAt:   note.ModifiedAt,
			DeckID:       primaryDeckID,
			DeckName:     primaryDeckName,
			CardCount:    len(cards),
		})
	}

	sort.Slice(recentNotes, func(i, j int) bool {
		if recentNotes[i].ModifiedAt.Equal(recentNotes[j].ModifiedAt) {
			return recentNotes[i].ID > recentNotes[j].ID
		}
		return recentNotes[i].ModifiedAt.After(recentNotes[j].ModifiedAt)
	})
	if len(recentNotes) > 5 {
		recentNotes = recentNotes[:5]
	}

	for i := range studyAnalytics.RecentSessions {
		deckID := studyAnalytics.RecentSessions[i].DeckID
		if deckID == 0 {
			continue
		}
		if deck, ok := col.Decks[deckID]; ok {
			studyAnalytics.RecentSessions[i].DeckName = deck.Name
		}
	}

	return DashboardResponse{
		TotalDecks:     len(col.Decks),
		TotalNotes:     len(col.Notes),
		DueToday:       dueToday,
		Plan:           sessionResponse.Entitlements.Plan,
		Usage:          sessionResponse.Entitlements.Usage,
		Limits:         sessionResponse.Entitlements.Limits,
		Features:       sessionResponse.Entitlements.Features,
		StudyAnalytics: studyAnalytics,
		RecentNotes:    recentNotes,
	}
}

func (h *APIHandler) GetDeckStats(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDParam(r, "id")
	if err != nil {
		http.Error(w, "Invalid deck ID", http.StatusBadRequest)
		return
	}

	stats, err := h.store.GetDeckStatsForUser(h.userIDFromRequest(r), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusOK, stats)
}

func (h *APIHandler) CreateNote(w http.ResponseWriter, r *http.Request) {
	if !h.requireWorkspaceWritePermission(w, r) {
		return
	}
	col, collectionID, err := h.collectionForRequest(r)
	if err != nil {
		respondAPIError(w, http.StatusInternalServerError, "collection_load_failed", err.Error())
		return
	}

	var req CreateNoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondAPIError(w, http.StatusBadRequest, "invalid_request", "Invalid request body")
		return
	}

	if req.TypeID == "" || req.DeckID == 0 {
		respondAPIError(w, http.StatusBadRequest, "invalid_note_request", "TypeID and DeckID are required")
		return
	}

	session := h.sessionFromRequest(r)
	plan := h.planForRequest(r, session)
	usage := h.usageForSession(session)
	if err := validateNoteLimit(plan, usage); err != nil {
		respondAPIError(w, http.StatusForbidden, "plan_limit_exceeded", err.Error())
		return
	}

	sanitizedFieldVals := sanitizeFieldVals(req.FieldVals)
	sanitizedTags := sanitizeTags(req.Tags)

	noteType, ok := col.NoteTypes[NoteTypeName(req.TypeID)]
	if !ok {
		respondAPIError(w, http.StatusBadRequest, "invalid_note_type", "Note type not found")
		return
	}
	previewAt := time.Now()
	previewNote := Note{
		Type:       NoteTypeName(req.TypeID),
		FieldMap:   sanitizedFieldVals,
		Tags:       sanitizedTags,
		CreatedAt:  previewAt,
		ModifiedAt: previewAt,
	}
	previewCards, err := col.generateCardsFromNote(noteType, previewNote, req.DeckID, previewAt)
	if err != nil {
		respondAPIError(w, http.StatusBadRequest, "note_create_failed", err.Error())
		return
	}
	if err := validateCardsTotalLimit(plan, usage, len(previewCards)); err != nil {
		respondAPIError(w, http.StatusForbidden, "plan_limit_exceeded", err.Error())
		return
	}

	// Use Collection.AddNote to create note and generate cards
	note, cards, err := col.AddNote(req.DeckID, NoteTypeName(req.TypeID), sanitizedFieldVals, time.Now())
	if err != nil {
		respondAPIError(w, http.StatusBadRequest, "note_create_failed", err.Error())
		return
	}

	// Set tags if provided (use sanitized tags)
	note.Tags = sanitizedTags
	// Persist note to database
	if err := h.store.CreateNote(collectionID, &note); err != nil {
		respondAPIError(w, http.StatusInternalServerError, "note_persist_failed", err.Error())
		return
	}

	// Persist generated cards to database
	for _, card := range cards {
		if err := h.store.CreateCard(card); err != nil {
			respondAPIError(w, http.StatusInternalServerError, "card_persist_failed", fmt.Sprintf("Failed to save card: %v", err))
			return
		}
	}

	responseCards := make([]Card, 0, len(cards))
	for _, card := range cards {
		responseCards = append(responseCards, *card)
	}
	h.markStudyGroupInstallsForkedByDeckIDs(req.DeckID)

	respondJSON(w, http.StatusCreated, map[string]interface{}{
		"note":  h.noteToResponse(&note, responseCards),
		"cards": responseCards,
	})
}

func (h *APIHandler) ImportNotes(w http.ResponseWriter, r *http.Request) {
	if !h.requireWorkspaceWritePermission(w, r) {
		return
	}
	fileData, opts, err := parseImportRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	opts.DefaultDeckName = firstNonEmpty(opts.DefaultDeckName, "Default")
	opts.DefaultNoteType = firstNonEmpty(opts.DefaultNoteType, "Basic")

	parsed, err := parseImportData(fileData, opts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	col, collectionID, err := h.collectionForRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	importResult := h.applyImportedNotesToCollection(collectionID, col, parsed.Notes, opts.DefaultDeckName)
	importResult.Source = parsed.Source
	importResult.Format = parsed.Format

	if importResult.Imported == 0 {
		respondJSON(w, http.StatusBadRequest, importResult)
		return
	}

	respondJSON(w, http.StatusOK, importResult)
}

func (h *APIHandler) GetNote(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDParam(r, "id")
	if err != nil {
		http.Error(w, "Invalid note ID", http.StatusBadRequest)
		return
	}

	note, err := h.store.GetNote(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	cards, err := h.store.GetCardsByNote(id)
	if err != nil {
		respondAPIError(w, http.StatusInternalServerError, "note_cards_failed", err.Error())
		return
	}

	respondJSON(w, http.StatusOK, h.noteToResponse(note, cards))
}

func (h *APIHandler) CheckDuplicate(w http.ResponseWriter, r *http.Request) {
	var req CheckDuplicateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Value == "" {
		respondJSON(w, http.StatusOK, DuplicateResult{IsDuplicate: false})
		return
	}

	// Check for duplicates in the collection
	duplicates, err := h.store.FindDuplicateNotes(h.collectionIDForRequest(r), req.FieldName, req.Value, req.DeckID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	result := DuplicateResult{
		IsDuplicate: len(duplicates) > 0,
		Duplicates:  duplicates,
	}

	respondJSON(w, http.StatusOK, result)
}

func (h *APIHandler) GetDueCards(w http.ResponseWriter, r *http.Request) {
	deckID, err := parseIDParam(r, "deckId")
	if err != nil {
		http.Error(w, "Invalid deck ID", http.StatusBadRequest)
		return
	}

	limit := 10
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	cards, err := h.store.GetDueCardsForUser(h.userIDFromRequest(r), deckID, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusOK, cards)
}

func (h *APIHandler) GetCard(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDParam(r, "id")
	if err != nil {
		http.Error(w, "Invalid card ID", http.StatusBadRequest)
		return
	}

	card, err := h.store.GetCardForUser(h.userIDFromRequest(r), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	respondJSON(w, http.StatusOK, card)
}

func (h *APIHandler) AnswerCard(w http.ResponseWriter, r *http.Request) {
	col, _, err := h.collectionForRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, err := parseIDParam(r, "id")
	if err != nil {
		http.Error(w, "Invalid card ID", http.StatusBadRequest)
		return
	}

	var req AnswerCardRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Rating < 1 || req.Rating > 4 {
		http.Error(w, "Rating must be 1-4 (Again/Hard/Good/Easy)", http.StatusBadRequest)
		return
	}

	userID := h.userIDFromRequest(r)
	card, err := h.store.GetCardForUser(userID, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sched := fsrs.NewFSRS(col.Params).Repeat(card.SRS, time.Now())
	info, ok := sched[fsrs.Rating(req.Rating)]
	if !ok {
		http.Error(w, "Unable to schedule card review", http.StatusInternalServerError)
		return
	}
	card.SRS = info.Card

	if err := h.store.UpdateCardReviewState(userID, card); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := h.store.AddRevlogForUser(userID, &info.ReviewLog, id, req.TimeTakenMs); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusOK, card)
}

func (h *APIHandler) UpdateCard(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDParam(r, "id")
	if err != nil {
		http.Error(w, "Invalid card ID", http.StatusBadRequest)
		return
	}

	var req UpdateCardRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userID := h.userIDFromRequest(r)

	// Get card from store
	card, err := h.store.GetCardForUser(userID, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Update fields if provided
	if req.Flag != nil {
		if *req.Flag < 0 || *req.Flag > 7 {
			http.Error(w, "Flag must be 0-7", http.StatusBadRequest)
			return
		}
		card.Flag = *req.Flag
	}
	if req.Marked != nil {
		card.Marked = *req.Marked
	}
	if req.Suspended != nil {
		card.Suspended = *req.Suspended
	}

	// Persist changes
	if err := h.store.UpdateCardReviewState(userID, card); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusOK, card)
}

func (h *APIHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, map[string]string{
		"status":  "ok",
		"service": "vutadex-api",
		"version": "M2",
	})
}

// Note type response for API
type NoteTypeResponse struct {
	Name           string                  `json:"name"`
	Fields         []string                `json:"fields"`
	Templates      []TemplateInfo          `json:"templates"`
	SortFieldIndex int                     `json:"sortFieldIndex"`
	FieldOptions   map[string]FieldOptions `json:"fieldOptions,omitempty"`
}

type TemplateInfo struct {
	Name            string `json:"name"`
	QFmt            string `json:"qFmt"`
	AFmt            string `json:"aFmt"`
	Styling         string `json:"styling"`
	IfFieldNonEmpty string `json:"ifFieldNonEmpty,omitempty"`
	IsCloze         bool   `json:"isCloze"`
	DeckOverride    string `json:"deckOverride,omitempty"`
	BrowserQFmt     string `json:"browserQFmt,omitempty"`
	BrowserAFmt     string `json:"browserAFmt,omitempty"`
}

type SetSortFieldRequest struct {
	FieldIndex int `json:"fieldIndex"` // Index of the field to use as sort field
}

type SetFieldOptionsRequest struct {
	FieldName string       `json:"fieldName"`
	Options   FieldOptions `json:"options"`
}

type UpdateTemplateRequest struct {
	Name            *string `json:"name,omitempty"`
	QFmt            *string `json:"qFmt,omitempty"`
	AFmt            *string `json:"aFmt,omitempty"`
	Styling         *string `json:"styling,omitempty"`
	IfFieldNonEmpty *string `json:"ifFieldNonEmpty,omitempty"`
	DeckOverride    *string `json:"deckOverride,omitempty"`
	BrowserQFmt     *string `json:"browserQFmt,omitempty"`
	BrowserAFmt     *string `json:"browserAFmt,omitempty"`
}

type TemplatesResponse struct {
	Message   string         `json:"message"`
	Templates []TemplateInfo `json:"templates"`
}

func (h *APIHandler) ListNoteTypes(w http.ResponseWriter, r *http.Request) {
	col, _, err := h.collectionForRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var noteTypes []NoteTypeResponse
	for _, nt := range col.NoteTypes {
		var templates []TemplateInfo
		for _, t := range nt.Templates {
			templates = append(templates, TemplateInfo{
				Name:            t.Name,
				QFmt:            t.QFmt,
				AFmt:            t.AFmt,
				Styling:         t.Styling,
				IfFieldNonEmpty: t.IfFieldNonEmpty,
				IsCloze:         t.IsCloze,
				DeckOverride:    t.DeckOverride,
				BrowserQFmt:     t.BrowserQFmt,
				BrowserAFmt:     t.BrowserAFmt,
			})
		}
		noteTypes = append(noteTypes, NoteTypeResponse{
			Name:           string(nt.Name),
			Fields:         nt.Fields,
			Templates:      templates,
			SortFieldIndex: nt.SortFieldIndex,
			FieldOptions:   nt.FieldOptions,
		})
	}

	sort.Slice(noteTypes, func(i, j int) bool {
		return noteTypes[i].Name < noteTypes[j].Name
	})

	respondJSON(w, http.StatusOK, noteTypes)
}

func (h *APIHandler) GetNoteType(w http.ResponseWriter, r *http.Request) {
	col, _, err := h.collectionForRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	name := chi.URLParam(r, "name")
	nt, ok := col.NoteTypes[NoteTypeName(name)]
	if !ok {
		http.Error(w, "Note type not found", http.StatusNotFound)
		return
	}

	var templates []TemplateInfo
	for _, t := range nt.Templates {
		templates = append(templates, TemplateInfo{
			Name:            t.Name,
			QFmt:            t.QFmt,
			AFmt:            t.AFmt,
			Styling:         t.Styling,
			IfFieldNonEmpty: t.IfFieldNonEmpty,
			IsCloze:         t.IsCloze,
			DeckOverride:    t.DeckOverride,
			BrowserQFmt:     t.BrowserQFmt,
			BrowserAFmt:     t.BrowserAFmt,
		})
	}

	respondJSON(w, http.StatusOK, NoteTypeResponse{
		Name:           string(nt.Name),
		Fields:         nt.Fields,
		Templates:      templates,
		SortFieldIndex: nt.SortFieldIndex,
		FieldOptions:   nt.FieldOptions,
	})
}

// Reserved field names that cannot be used
var reservedFieldNames = map[string]bool{
	"Tags":      true,
	"Type":      true,
	"Deck":      true,
	"Card":      true,
	"FrontSide": true,
}

// Field management request types
type AddFieldRequest struct {
	FieldName string `json:"fieldName"`
	Position  *int   `json:"position,omitempty"` // Optional: insert at specific position
}

type RenameFieldRequest struct {
	OldName string `json:"oldName"`
	NewName string `json:"newName"`
}

type RemoveFieldRequest struct {
	FieldName string `json:"fieldName"`
}

type ReorderFieldsRequest struct {
	Fields []string `json:"fields"` // New field order
}

func (h *APIHandler) AddField(w http.ResponseWriter, r *http.Request) {
	if !h.requireWorkspaceWritePermission(w, r) {
		return
	}
	col, collectionID, err := h.collectionForRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	name := chi.URLParam(r, "name")
	nt, ok := col.NoteTypes[NoteTypeName(name)]
	if !ok {
		http.Error(w, "Note type not found", http.StatusNotFound)
		return
	}

	var req AddFieldRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.FieldName == "" {
		http.Error(w, "fieldName is required", http.StatusBadRequest)
		return
	}

	// Sanitize field name to prevent XSS
	sanitizedFieldName := sanitizeHTML(req.FieldName)

	// Check for reserved field names
	if reservedFieldNames[sanitizedFieldName] {
		http.Error(w, fmt.Sprintf("'%s' is a reserved field name", sanitizedFieldName), http.StatusBadRequest)
		return
	}

	// Check for duplicate field name
	for _, f := range nt.Fields {
		if f == sanitizedFieldName {
			http.Error(w, "Field name already exists", http.StatusBadRequest)
			return
		}
	}

	// Add field at specified position or end
	if req.Position != nil && *req.Position >= 0 && *req.Position < len(nt.Fields) {
		newFields := make([]string, 0, len(nt.Fields)+1)
		newFields = append(newFields, nt.Fields[:*req.Position]...)
		newFields = append(newFields, sanitizedFieldName)
		newFields = append(newFields, nt.Fields[*req.Position:]...)
		nt.Fields = newFields
	} else {
		nt.Fields = append(nt.Fields, sanitizedFieldName)
	}

	// Update in store
	if err := h.store.UpdateNoteType(collectionID, &nt); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Update collection cache
	col.NoteTypes[NoteTypeName(name)] = nt
	h.markStudyGroupInstallsForkedByNoteType(name)

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Field added successfully",
		"fields":  nt.Fields,
	})
}

func (h *APIHandler) RenameField(w http.ResponseWriter, r *http.Request) {
	if !h.requireWorkspaceWritePermission(w, r) {
		return
	}
	col, collectionID, err := h.collectionForRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	name := chi.URLParam(r, "name")
	nt, ok := col.NoteTypes[NoteTypeName(name)]
	if !ok {
		http.Error(w, "Note type not found", http.StatusNotFound)
		return
	}

	var req RenameFieldRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.OldName == "" || req.NewName == "" {
		http.Error(w, "oldName and newName are required", http.StatusBadRequest)
		return
	}

	// Sanitize new field name to prevent XSS
	sanitizedNewName := sanitizeHTML(req.NewName)

	// Check for reserved field names
	if reservedFieldNames[sanitizedNewName] {
		http.Error(w, fmt.Sprintf("'%s' is a reserved field name", sanitizedNewName), http.StatusBadRequest)
		return
	}

	// Find and rename the field
	found := false
	for i, f := range nt.Fields {
		if f == req.OldName {
			nt.Fields[i] = sanitizedNewName
			found = true
			break
		}
	}

	if !found {
		http.Error(w, "Field not found", http.StatusNotFound)
		return
	}

	// Check for duplicate with new name
	count := 0
	for _, f := range nt.Fields {
		if f == sanitizedNewName {
			count++
		}
	}
	if count > 1 {
		// Revert and return error
		for i, f := range nt.Fields {
			if f == sanitizedNewName {
				nt.Fields[i] = req.OldName
				break
			}
		}
		http.Error(w, "Field name already exists", http.StatusBadRequest)
		return
	}

	// Update templates to use new field name
	for i := range nt.Templates {
		nt.Templates[i].QFmt = strings.ReplaceAll(nt.Templates[i].QFmt, "{{"+req.OldName+"}}", "{{"+sanitizedNewName+"}}")
		nt.Templates[i].AFmt = strings.ReplaceAll(nt.Templates[i].AFmt, "{{"+req.OldName+"}}", "{{"+sanitizedNewName+"}}")
		if nt.Templates[i].IfFieldNonEmpty == req.OldName {
			nt.Templates[i].IfFieldNonEmpty = sanitizedNewName
		}
	}

	// Update in store
	if err := h.store.UpdateNoteType(collectionID, &nt); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Update collection cache
	col.NoteTypes[NoteTypeName(name)] = nt
	h.markStudyGroupInstallsForkedByNoteType(name)

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Field renamed successfully",
		"fields":  nt.Fields,
	})
}

func (h *APIHandler) RemoveField(w http.ResponseWriter, r *http.Request) {
	if !h.requireWorkspaceWritePermission(w, r) {
		return
	}
	col, collectionID, err := h.collectionForRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	name := chi.URLParam(r, "name")
	nt, ok := col.NoteTypes[NoteTypeName(name)]
	if !ok {
		http.Error(w, "Note type not found", http.StatusNotFound)
		return
	}

	var req RemoveFieldRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.FieldName == "" {
		http.Error(w, "fieldName is required", http.StatusBadRequest)
		return
	}

	// Must have at least one field
	if len(nt.Fields) <= 1 {
		http.Error(w, "Cannot remove the last field", http.StatusBadRequest)
		return
	}

	// Find and remove the field
	found := false
	newFields := make([]string, 0, len(nt.Fields)-1)
	for _, f := range nt.Fields {
		if f == req.FieldName {
			found = true
		} else {
			newFields = append(newFields, f)
		}
	}

	if !found {
		http.Error(w, "Field not found", http.StatusNotFound)
		return
	}

	nt.Fields = newFields

	// Update in store
	if err := h.store.UpdateNoteType(collectionID, &nt); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Update collection cache
	col.NoteTypes[NoteTypeName(name)] = nt
	h.markStudyGroupInstallsForkedByNoteType(name)

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Field removed successfully",
		"fields":  nt.Fields,
	})
}

func (h *APIHandler) ReorderFields(w http.ResponseWriter, r *http.Request) {
	if !h.requireWorkspaceWritePermission(w, r) {
		return
	}
	col, collectionID, err := h.collectionForRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	name := chi.URLParam(r, "name")
	nt, ok := col.NoteTypes[NoteTypeName(name)]
	if !ok {
		http.Error(w, "Note type not found", http.StatusNotFound)
		return
	}

	var req ReorderFieldsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate that the new order contains the same fields
	if len(req.Fields) != len(nt.Fields) {
		http.Error(w, "Field count mismatch", http.StatusBadRequest)
		return
	}

	// Check that all existing fields are present
	existingFields := make(map[string]bool)
	for _, f := range nt.Fields {
		existingFields[f] = true
	}

	for _, f := range req.Fields {
		if !existingFields[f] {
			http.Error(w, fmt.Sprintf("Unknown field: %s", f), http.StatusBadRequest)
			return
		}
		delete(existingFields, f)
	}

	if len(existingFields) > 0 {
		http.Error(w, "Some fields are missing from the new order", http.StatusBadRequest)
		return
	}

	nt.Fields = req.Fields

	// Update in store
	if err := h.store.UpdateNoteType(collectionID, &nt); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Update collection cache
	col.NoteTypes[NoteTypeName(name)] = nt
	h.markStudyGroupInstallsForkedByNoteType(name)

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Fields reordered successfully",
		"fields":  nt.Fields,
	})
}

func (h *APIHandler) SetSortField(w http.ResponseWriter, r *http.Request) {
	if !h.requireWorkspaceWritePermission(w, r) {
		return
	}
	col, collectionID, err := h.collectionForRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	name := chi.URLParam(r, "name")
	nt, ok := col.NoteTypes[NoteTypeName(name)]
	if !ok {
		http.Error(w, "Note type not found", http.StatusNotFound)
		return
	}

	var req SetSortFieldRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate field index
	if req.FieldIndex < 0 || req.FieldIndex >= len(nt.Fields) {
		http.Error(w, "Invalid field index", http.StatusBadRequest)
		return
	}

	nt.SortFieldIndex = req.FieldIndex

	// Update in store
	if err := h.store.UpdateNoteType(collectionID, &nt); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Update collection cache
	col.NoteTypes[NoteTypeName(name)] = nt
	h.markStudyGroupInstallsForkedByNoteType(name)

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"message":        "Sort field updated successfully",
		"sortFieldIndex": nt.SortFieldIndex,
		"sortFieldName":  nt.Fields[nt.SortFieldIndex],
	})
}

func (h *APIHandler) SetFieldOptions(w http.ResponseWriter, r *http.Request) {
	if !h.requireWorkspaceWritePermission(w, r) {
		return
	}
	col, collectionID, err := h.collectionForRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	name := chi.URLParam(r, "name")
	nt, ok := col.NoteTypes[NoteTypeName(name)]
	if !ok {
		http.Error(w, "Note type not found", http.StatusNotFound)
		return
	}

	var req SetFieldOptionsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate field exists
	fieldExists := false
	for _, f := range nt.Fields {
		if f == req.FieldName {
			fieldExists = true
			break
		}
	}
	if !fieldExists {
		http.Error(w, "Field not found", http.StatusBadRequest)
		return
	}

	// Initialize FieldOptions map if nil
	if nt.FieldOptions == nil {
		nt.FieldOptions = make(map[string]FieldOptions)
	}

	// Update field options
	nt.FieldOptions[req.FieldName] = req.Options

	// Update in store
	if err := h.store.UpdateNoteType(collectionID, &nt); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Update collection cache
	col.NoteTypes[NoteTypeName(name)] = nt
	h.markStudyGroupInstallsForkedByNoteType(name)

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"message":      "Field options updated successfully",
		"fieldOptions": nt.FieldOptions,
	})
}

func (h *APIHandler) UpdateTemplate(w http.ResponseWriter, r *http.Request) {
	if !h.requireWorkspaceWritePermission(w, r) {
		return
	}
	col, collectionID, err := h.collectionForRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	noteTypeName := chi.URLParam(r, "name")
	templateName := chi.URLParam(r, "templateName")

	nt, ok := col.NoteTypes[NoteTypeName(noteTypeName)]
	if !ok {
		http.Error(w, "Note type not found", http.StatusNotFound)
		return
	}

	var req UpdateTemplateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Find and update the template
	templateIndex := -1
	for i, t := range nt.Templates {
		if t.Name == templateName {
			templateIndex = i
			break
		}
	}

	if templateIndex == -1 {
		http.Error(w, "Template not found", http.StatusNotFound)
		return
	}

	templateAliases := map[string]string{}
	// Update fields if provided
	if req.Name != nil {
		nextName := sanitizeHTML(strings.TrimSpace(*req.Name))
		if nextName == "" {
			http.Error(w, "Template name is required", http.StatusBadRequest)
			return
		}
		if !strings.EqualFold(nextName, templateName) {
			for i, candidate := range nt.Templates {
				if i != templateIndex && strings.EqualFold(candidate.Name, nextName) {
					http.Error(w, "Template name already exists", http.StatusBadRequest)
					return
				}
			}
			templateAliases[templateName] = nextName
			nt.Templates[templateIndex].Name = nextName
		}
	}
	if req.QFmt != nil {
		nt.Templates[templateIndex].QFmt = sanitizeHTML(*req.QFmt)
	}
	if req.AFmt != nil {
		nt.Templates[templateIndex].AFmt = sanitizeHTML(*req.AFmt)
	}
	if req.Styling != nil {
		nt.Templates[templateIndex].Styling = sanitizeHTML(*req.Styling)
	}
	if req.IfFieldNonEmpty != nil {
		nt.Templates[templateIndex].IfFieldNonEmpty = sanitizeHTML(*req.IfFieldNonEmpty)
	}
	if req.DeckOverride != nil {
		nt.Templates[templateIndex].DeckOverride = sanitizeHTML(*req.DeckOverride)
	}
	if req.BrowserQFmt != nil {
		nt.Templates[templateIndex].BrowserQFmt = sanitizeHTML(*req.BrowserQFmt)
	}
	if req.BrowserAFmt != nil {
		nt.Templates[templateIndex].BrowserAFmt = sanitizeHTML(*req.BrowserAFmt)
	}

	// Update in store
	if err := h.store.UpdateNoteType(collectionID, &nt); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Update collection cache
	col.NoteTypes[NoteTypeName(noteTypeName)] = nt
	h.markStudyGroupInstallsForkedByNoteType(noteTypeName)

	// Regenerate cards for all notes of this type
	// This ensures cards reflect the updated templates
	if err := h.regenerateCardsForNoteTypeWithAliases(collectionID, col, noteTypeName, templateAliases); err != nil {
		log.Printf("Warning: Failed to regenerate cards after template update: %v", err)
		// Don't fail the request - template was updated successfully
	}

	respondJSON(w, http.StatusOK, buildTemplatesResponse(nt, "Template updated successfully"))
}

// regenerateCardsForNoteType regenerates cards for all notes of a given note type.
// This preserves existing card scheduling data (SRS state, flags, etc.) while updating content.
func (h *APIHandler) regenerateCardsForNoteType(noteTypeName string) error {
	return h.regenerateCardsForNoteTypeWithAliases(h.collectionID, h.collection, noteTypeName, nil)
}

func (h *APIHandler) regenerateCardsForNoteTypeInCollection(collectionID string, col *Collection, noteTypeName string) error {
	return h.regenerateCardsForNoteTypeWithAliases(collectionID, col, noteTypeName, nil)
}

func (h *APIHandler) regenerateCardsForNoteTypeWithAliases(collectionID string, col *Collection, noteTypeName string, templateAliases map[string]string) error {
	// Get all notes of this type
	notes, err := h.store.GetNotesByType(collectionID, noteTypeName)
	if err != nil {
		return fmt.Errorf("failed to get notes: %w", err)
	}

	for _, note := range notes {
		// Get the deck ID from one of the note's existing cards
		// If the note has no cards, we'll use the default deck
		deckID := int64(1) // default deck
		existingCards, err := h.store.GetCardsByNote(note.ID)
		if err == nil && len(existingCards) > 0 {
			deckID = existingCards[0].DeckID
		}

		if _, err := h.regenerateCardsForSingleNote(col, &note, deckID, templateAliases); err != nil {
			log.Printf("Warning: Failed to regenerate cards for note %d: %v", note.ID, err)
		}
	}

	return nil
}

// Empty cards detection and cleanup

type EmptyCardInfo struct {
	CardID       int64  `json:"cardId"`
	NoteID       int64  `json:"noteId"`
	DeckID       int64  `json:"deckId"`
	TemplateName string `json:"templateName"`
	Ordinal      int    `json:"ordinal"`
	Front        string `json:"front"`
	Back         string `json:"back"`
	Reason       string `json:"reason"` // Why this card is considered empty
}

type EmptyCardsResponse struct {
	Count      int             `json:"count"`
	EmptyCards []EmptyCardInfo `json:"emptyCards"`
}

// FindEmptyCards detects cards that have no meaningful content
// This primarily targets cloze cards where the cloze deletion was removed
func (h *APIHandler) FindEmptyCards(w http.ResponseWriter, r *http.Request) {
	col, collectionID, err := h.collectionForRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get all notes
	notes, err := h.store.ListNotes(collectionID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var emptyCards []EmptyCardInfo

	for _, note := range notes {
		// Get note type to check if it's cloze
		nt, ok := col.NoteTypes[note.Type]
		if !ok {
			continue
		}

		// Get cards for this note
		cards, err := h.store.GetCardsByNote(note.ID)
		if err != nil {
			log.Printf("Warning: Failed to get cards for note %d: %v", note.ID, err)
			continue
		}

		for _, card := range cards {
			// Check if this is a cloze card
			isCloze := false
			for _, tmpl := range nt.Templates {
				if tmpl.Name == card.TemplateName && tmpl.IsCloze {
					isCloze = true
					break
				}
			}

			if isCloze {
				// Check if the cloze ordinal still exists in the text
				textField := note.FieldMap["Text"]
				ordinals := extractClozeOrdinals(textField)

				hasOrdinal := false
				for _, ord := range ordinals {
					if ord == card.Ordinal {
						hasOrdinal = true
						break
					}
				}

				if !hasOrdinal {
					emptyCards = append(emptyCards, EmptyCardInfo{
						CardID:       card.ID,
						NoteID:       note.ID,
						DeckID:       card.DeckID,
						TemplateName: card.TemplateName,
						Ordinal:      card.Ordinal,
						Front:        card.Front,
						Back:         card.Back,
						Reason:       fmt.Sprintf("Cloze deletion c%d no longer exists in note", card.Ordinal),
					})
				}
			} else {
				// For non-cloze cards, check if front and back are essentially empty
				// Strip HTML tags and check if there's meaningful content
				frontStripped := stripHTML(card.Front)
				backStripped := stripHTML(card.Back)

				if strings.TrimSpace(frontStripped) == "" && strings.TrimSpace(backStripped) == "" {
					emptyCards = append(emptyCards, EmptyCardInfo{
						CardID:       card.ID,
						NoteID:       note.ID,
						DeckID:       card.DeckID,
						TemplateName: card.TemplateName,
						Ordinal:      card.Ordinal,
						Front:        card.Front,
						Back:         card.Back,
						Reason:       "Card has no content (both front and back are empty)",
					})
				}
			}
		}
	}

	respondJSON(w, http.StatusOK, EmptyCardsResponse{
		Count:      len(emptyCards),
		EmptyCards: emptyCards,
	})
}

type DeleteEmptyCardsRequest struct {
	CardIDs []int64 `json:"cardIds"` // List of card IDs to delete
}

type DeleteEmptyCardsResponse struct {
	Deleted int      `json:"deleted"`
	Failed  []string `json:"failed,omitempty"` // Error messages for failed deletions
}

// DeleteEmptyCards deletes specified empty cards
func (h *APIHandler) DeleteEmptyCards(w http.ResponseWriter, r *http.Request) {
	col, _, err := h.collectionForRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var req DeleteEmptyCardsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if len(req.CardIDs) == 0 {
		http.Error(w, "No card IDs provided", http.StatusBadRequest)
		return
	}

	deleted := 0
	var failed []string

	for _, cardID := range req.CardIDs {
		if err := h.store.DeleteCard(cardID); err != nil {
			failed = append(failed, fmt.Sprintf("Card %d: %v", cardID, err))
			log.Printf("Failed to delete card %d: %v", cardID, err)
		} else {
			deleted++
			// Remove from collection cache
			delete(col.Cards, cardID)
		}
	}

	respondJSON(w, http.StatusOK, DeleteEmptyCardsResponse{
		Deleted: deleted,
		Failed:  failed,
	})
}

func parseImportRequest(r *http.Request) ([]byte, importParseOptions, error) {
	contentType := strings.ToLower(r.Header.Get("Content-Type"))
	if strings.HasPrefix(contentType, "multipart/form-data") {
		if err := r.ParseMultipartForm(32 << 20); err != nil {
			return nil, importParseOptions{}, fmt.Errorf("invalid multipart form: %w", err)
		}

		file, header, err := r.FormFile("file")
		if err != nil {
			return nil, importParseOptions{}, fmt.Errorf("file is required: %w", err)
		}
		defer file.Close()

		fileData, err := io.ReadAll(file)
		if err != nil {
			return nil, importParseOptions{}, fmt.Errorf("failed to read file: %w", err)
		}
		if len(strings.TrimSpace(string(fileData))) == 0 {
			return nil, importParseOptions{}, fmt.Errorf("import file is empty")
		}

		return fileData, importParseOptions{
			Source:          r.FormValue("source"),
			FormatHint:      r.FormValue("format"),
			Filename:        header.Filename,
			DefaultDeckName: r.FormValue("deckName"),
			DefaultNoteType: r.FormValue("noteType"),
		}, nil
	}

	var req ImportNotesJSONRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, importParseOptions{}, fmt.Errorf("invalid request body")
	}
	if strings.TrimSpace(req.Content) == "" {
		return nil, importParseOptions{}, fmt.Errorf("content is required when file upload is not used")
	}

	return []byte(req.Content), importParseOptions{
		Source:          req.Source,
		FormatHint:      req.Format,
		Filename:        req.Filename,
		DefaultDeckName: req.DeckName,
		DefaultNoteType: req.NoteType,
	}, nil
}

func (h *APIHandler) applyImportedNotes(notes []importNormalizedNote, defaultDeckName string) ImportNotesResponse {
	return h.applyImportedNotesToCollection(h.collectionID, h.collection, notes, defaultDeckName)
}

func (h *APIHandler) applyImportedNotesToCollection(collectionID string, col *Collection, notes []importNormalizedNote, defaultDeckName string) ImportNotesResponse {
	result := ImportNotesResponse{}
	deckCache := make(map[string]int64)
	createdDecks := make(map[string]struct{})

	for id, deck := range col.Decks {
		deckCache[strings.ToLower(deck.Name)] = id
	}

	for i, importedNote := range notes {
		noteTypeName := importedNote.NoteType
		if noteTypeName == "" {
			noteTypeName = "Basic"
		}

		noteType, ok := col.NoteTypes[noteTypeName]
		if !ok {
			result.Skipped++
			result.Errors = append(result.Errors, fmt.Sprintf("row %d: unknown note type %q", i+1, noteTypeName))
			continue
		}

		deckName := firstNonEmpty(importedNote.DeckName, defaultDeckName, "Default")
		deckID, err := h.ensureDeckByName(collectionID, col, deckName, deckCache, createdDecks)
		if err != nil {
			result.Skipped++
			result.Errors = append(result.Errors, fmt.Sprintf("row %d: failed to resolve deck %q: %v", i+1, deckName, err))
			continue
		}

		fieldVals := make(map[string]string, len(noteType.Fields))
		allEmpty := true
		for _, fieldName := range noteType.Fields {
			value, _ := getFieldValueCaseInsensitive(importedNote.Fields, fieldName)
			sanitizedValue := sanitizeHTML(value)
			fieldVals[fieldName] = sanitizedValue
			if strings.TrimSpace(sanitizedValue) != "" {
				allEmpty = false
			}
		}

		if allEmpty {
			result.Skipped++
			result.Errors = append(result.Errors, fmt.Sprintf("row %d: note has no content after field mapping", i+1))
			continue
		}

		note, cards, err := col.AddNote(deckID, noteTypeName, fieldVals, time.Now())
		if err != nil {
			result.Skipped++
			result.Errors = append(result.Errors, fmt.Sprintf("row %d: failed to create note: %v", i+1, err))
			continue
		}

		note.Tags = sanitizeImportTags(importedNote.Tags)
		if note.Tags == nil {
			note.Tags = []string{}
		}

		if err := h.store.CreateNote(collectionID, &note); err != nil {
			result.Skipped++
			result.Errors = append(result.Errors, fmt.Sprintf("row %d: failed to persist note: %v", i+1, err))
			continue
		}

		cardErr := false
		for _, card := range cards {
			if err := h.store.CreateCard(card); err != nil {
				cardErr = true
				result.Errors = append(result.Errors, fmt.Sprintf("row %d: failed to persist card: %v", i+1, err))
				break
			}
		}
		if cardErr {
			result.Skipped++
			continue
		}

		result.Imported++
	}

	result.DecksCreated = sortedKeys(createdDecks)
	return result
}

func (h *APIHandler) ensureDeckByName(collectionID string, col *Collection, deckName string, deckCache map[string]int64, createdDecks map[string]struct{}) (int64, error) {
	name := firstNonEmpty(deckName, "Default")
	key := strings.ToLower(name)
	if id, ok := deckCache[key]; ok {
		return id, nil
	}

	for id, deck := range col.Decks {
		if strings.EqualFold(deck.Name, name) {
			deckCache[key] = id
			return id, nil
		}
	}

	sanitized := sanitizeHTML(name)
	if strings.TrimSpace(sanitized) == "" {
		sanitized = "Imported"
	}

	newDeck := col.NewDeck(sanitized)
	if err := h.store.CreateDeckInCollection(collectionID, newDeck); err != nil {
		return 0, err
	}

	deckCache[key] = newDeck.ID
	createdDecks[newDeck.Name] = struct{}{}
	return newDeck.ID, nil
}

func sanitizeImportTags(tags []string) []string {
	if len(tags) == 0 {
		return nil
	}

	seen := make(map[string]struct{}, len(tags))
	out := make([]string, 0, len(tags))
	for _, raw := range tags {
		sanitized := strings.TrimSpace(sanitizeHTML(raw))
		if sanitized == "" {
			continue
		}
		if _, exists := seen[sanitized]; exists {
			continue
		}
		seen[sanitized] = struct{}{}
		out = append(out, sanitized)
	}
	return out
}

// stripHTML removes HTML tags and returns plain text
func stripHTML(html string) string {
	return htmlPolicy.Sanitize(html)
}

// Backup endpoints

func (h *APIHandler) CreateBackup(w http.ResponseWriter, r *http.Request) {
	backupPath, err := h.backupManager.CreateBackup(h.collectionIDForRequest(r))
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create backup: %v", err), http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusCreated, map[string]string{
		"message":    "Backup created successfully",
		"backupPath": backupPath,
		"timestamp":  time.Now().Format(time.RFC3339),
	})
}

type RestoreBackupRequest struct {
	BackupPath string `json:"backupPath"`
}

func (h *APIHandler) RestoreBackup(w http.ResponseWriter, r *http.Request) {
	var req RestoreBackupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.BackupPath == "" {
		http.Error(w, "backupPath is required", http.StatusBadRequest)
		return
	}

	// Warning: This will close the database connection and replace the database
	// In a production system, this would need more careful handling
	respondJSON(w, http.StatusOK, map[string]string{
		"message": "Restore functionality requires server restart. Use with caution.",
		"warning": "This operation will replace the current database.",
		"note":    "Implement with /api/backups/restore endpoint after server architecture supports it.",
	})
}

func (h *APIHandler) ListBackups(w http.ResponseWriter, r *http.Request) {
	files, err := filepath.Glob(filepath.Join(h.backupManager.backupDir, "microdote-backup-*.zip"))
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to list backups: %v", err), http.StatusInternalServerError)
		return
	}

	type backupInfo struct {
		Path     string    `json:"path"`
		Filename string    `json:"filename"`
		Size     int64     `json:"size"`
		Modified time.Time `json:"modified"`
	}

	var backups []backupInfo
	for _, path := range files {
		info, err := os.Stat(path)
		if err != nil {
			continue
		}

		backups = append(backups, backupInfo{
			Path:     path,
			Filename: filepath.Base(path),
			Size:     info.Size(),
			Modified: info.ModTime(),
		})
	}

	respondJSON(w, http.StatusOK, backups)
}

// Helper functions

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func parseIDParam(r *http.Request, paramName string) (int64, error) {
	idStr := chi.URLParam(r, paramName)
	return strconv.ParseInt(idStr, 10, 64)
}

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
)

type apiTestEnv struct {
	store      *SQLiteStore
	collection *Collection
	handler    *APIHandler
	router     http.Handler
	dbPath     string
	backupDir  string
	authCookie string
}

type otpEmailStub struct {
	lastTo      string
	lastCode    string
	lastExpires time.Time
}

func (s *otpEmailStub) SendOTP(_ context.Context, to, code string, expiresAt time.Time) error {
	s.lastTo = to
	s.lastCode = code
	s.lastExpires = expiresAt
	return nil
}

func setupAPITestEnv(t *testing.T) *apiTestEnv {
	t.Helper()

	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "microdote-test.db")
	backupDir := filepath.Join(tempDir, "backups")

	store, err := NewSQLiteStore(dbPath)
	if err != nil {
		t.Fatalf("failed to create sqlite store: %v", err)
	}
	t.Cleanup(func() {
		_ = store.Close()
	})

	col := NewCollection()
	col.NoteTypes = builtins()

	if err := store.CreateCollection(col); err != nil {
		t.Fatalf("failed to create collection: %v", err)
	}
	for _, nt := range col.NoteTypes {
		ntCopy := nt
		if err := store.CreateNoteType("default", &ntCopy); err != nil {
			t.Fatalf("failed to create note type %s: %v", nt.Name, err)
		}
	}

	// Ensure deck 1 always exists for tests that assume a valid deck target.
	deck := col.NewDeck("Default")
	if err := store.CreateDeck(deck); err != nil {
		t.Fatalf("failed to create default deck: %v", err)
	}

	now := time.Now()
	user := &User{
		ID:          newID("usr"),
		Email:       "test@example.com",
		DisplayName: "Test User",
		LastLoginAt: now,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := store.CreateUser(user); err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}

	workspace := &Workspace{
		ID:           newID("ws"),
		Name:         "Test Workspace",
		Slug:         "test-workspace",
		CollectionID: "default",
		OwnerUserID:  user.ID,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	if err := store.CreateWorkspaceRecord(workspace); err != nil {
		t.Fatalf("failed to create test workspace: %v", err)
	}

	session := &SessionRecord{
		ID:          newID("sess"),
		UserID:      user.ID,
		WorkspaceID: workspace.ID,
		Plan:        PlanFree,
		ExpiresAt:   now.Add(24 * time.Hour),
		LastSeenAt:  now,
		CreatedAt:   now,
	}
	if err := store.CreateSessionRecord(session); err != nil {
		t.Fatalf("failed to create test session: %v", err)
	}

	handler := NewAPIHandler(store, col, NewBackupManager(dbPath, backupDir, store))
	authCookie := fmt.Sprintf("%s=%s", sessionCookieName, session.ID)
	router := newTestAPIRouter(handler, authCookie)

	return &apiTestEnv{
		store:      store,
		collection: col,
		handler:    handler,
		router:     router,
		dbPath:     dbPath,
		backupDir:  backupDir,
		authCookie: authCookie,
	}
}

func newTestAPIRouter(handler *APIHandler, authCookie string) http.Handler {
	r := chi.NewRouter()
	if authCookie != "" {
		r.Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				if req.Header.Get("X-Test-No-Auth") == "1" {
					next.ServeHTTP(w, req)
					return
				}
				if req.Header.Get("Cookie") == "" {
					req.Header.Set("Cookie", authCookie)
				}
				next.ServeHTTP(w, req)
			})
		})
	}
	r.Route("/api", func(r chi.Router) {
		registerAPIRoutes(r, handler)
	})
	return r
}

func doJSONRequest(t *testing.T, router http.Handler, method, path string, payload interface{}) *httptest.ResponseRecorder {
	t.Helper()
	return doJSONRequestWithHeaders(t, router, method, path, payload, nil)
}

func doJSONRequestWithHeaders(t *testing.T, router http.Handler, method, path string, payload interface{}, headers map[string]string) *httptest.ResponseRecorder {
	t.Helper()
	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("failed to marshal request payload: %v", err)
	}
	req := httptest.NewRequest(method, path, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	return rr
}

func doRawRequest(router http.Handler, method, path string, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	return rr
}

func decodeJSON[T any](t *testing.T, rr *httptest.ResponseRecorder) T {
	t.Helper()
	var out T
	if err := json.Unmarshal(rr.Body.Bytes(), &out); err != nil {
		t.Fatalf("failed to decode response JSON (%d): %v\nbody=%s", rr.Code, err, rr.Body.String())
	}
	return out
}

func TestAPI_ProtectedEndpointsRequireAuth(t *testing.T) {
	env := setupAPITestEnv(t)

	rr := doJSONRequestWithHeaders(t, env.router, http.MethodGet, "/api/decks", map[string]string{}, map[string]string{
		"X-Test-No-Auth": "1",
	})

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 for unauthenticated decks request, got %d (%s)", rr.Code, rr.Body.String())
	}

	var apiErr APIErrorResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &apiErr); err != nil {
		t.Fatalf("failed to decode API error: %v", err)
	}
	if apiErr.Code != "auth_required" {
		t.Fatalf("expected auth_required code, got %+v", apiErr)
	}
}

func TestAPI_OTPRequestAndVerifyCreatesSession(t *testing.T) {
	env := setupAPITestEnv(t)
	emailStub := &otpEmailStub{}
	env.handler.emailSender = emailStub

	requestRR := doJSONRequestWithHeaders(t, env.router, http.MethodPost, "/api/auth/otp/request", map[string]string{
		"email": "otp@example.com",
	}, map[string]string{
		"X-Test-No-Auth": "1",
	})
	if requestRR.Code != http.StatusOK {
		t.Fatalf("expected OTP request 200, got %d (%s)", requestRR.Code, requestRR.Body.String())
	}
	if emailStub.lastCode == "" {
		t.Fatalf("expected OTP code to be captured by stub")
	}

	verifyRR := doJSONRequestWithHeaders(t, env.router, http.MethodPost, "/api/auth/otp/verify", map[string]string{
		"email": "otp@example.com",
		"code":  emailStub.lastCode,
	}, map[string]string{
		"X-Test-No-Auth": "1",
	})
	if verifyRR.Code != http.StatusOK {
		t.Fatalf("expected OTP verify 200, got %d (%s)", verifyRR.Code, verifyRR.Body.String())
	}

	session := decodeJSON[AuthSessionResponse](t, verifyRR)
	if !session.Authenticated {
		t.Fatalf("expected authenticated session response, got %+v", session)
	}
	if session.User == nil || session.User.Email != "otp@example.com" {
		t.Fatalf("expected created user in session response, got %+v", session.User)
	}
	if session.Workspace == nil {
		t.Fatalf("expected workspace in session response")
	}

	cookies := verifyRR.Result().Cookies()
	if len(cookies) == 0 || cookies[0].Name != sessionCookieName {
		t.Fatalf("expected session cookie to be set, got %+v", cookies)
	}
}

type createNoteAPIResponse struct {
	Note  Note   `json:"note"`
	Cards []Card `json:"cards"`
}

func createNoteForTest(t *testing.T, env *apiTestEnv, req CreateNoteRequest, headers map[string]string) createNoteAPIResponse {
	t.Helper()

	rr := doJSONRequestWithHeaders(t, env.router, http.MethodPost, "/api/notes", req, headers)
	if rr.Code != http.StatusCreated {
		t.Fatalf("expected create note 201, got %d (%s)", rr.Code, rr.Body.String())
	}

	return decodeJSON[createNoteAPIResponse](t, rr)
}

func TestAPI_DeckAndCollectionEndpoints(t *testing.T) {
	env := setupAPITestEnv(t)

	health := doRawRequest(env.router, http.MethodGet, "/api/health", "")
	if health.Code != http.StatusOK {
		t.Fatalf("expected health 200, got %d", health.Code)
	}
	healthMap := decodeJSON[map[string]string](t, health)
	if healthMap["status"] != "ok" {
		t.Fatalf("expected status ok, got %q", healthMap["status"])
	}

	collectionResp := doRawRequest(env.router, http.MethodGet, "/api/collection", "")
	if collectionResp.Code != http.StatusOK {
		t.Fatalf("expected collection 200, got %d", collectionResp.Code)
	}

	badDeckBody := doRawRequest(env.router, http.MethodPost, "/api/decks", "{")
	if badDeckBody.Code != http.StatusBadRequest {
		t.Fatalf("expected create deck invalid body 400, got %d", badDeckBody.Code)
	}

	emptyDeckName := doJSONRequest(t, env.router, http.MethodPost, "/api/decks", CreateDeckRequest{Name: ""})
	if emptyDeckName.Code != http.StatusBadRequest {
		t.Fatalf("expected create deck empty name 400, got %d", emptyDeckName.Code)
	}

	createDeck := doJSONRequest(t, env.router, http.MethodPost, "/api/decks", CreateDeckRequest{Name: "<script>alert(1)</script>API Deck"})
	if createDeck.Code != http.StatusCreated {
		t.Fatalf("expected create deck 201, got %d (%s)", createDeck.Code, createDeck.Body.String())
	}
	createdDeck := decodeJSON[DeckResponse](t, createDeck)
	if strings.Contains(createdDeck.Name, "<script>") {
		t.Fatalf("expected deck name to be sanitized, got %q", createdDeck.Name)
	}

	listDecks := doRawRequest(env.router, http.MethodGet, "/api/decks", "")
	if listDecks.Code != http.StatusOK {
		t.Fatalf("expected list decks 200, got %d", listDecks.Code)
	}
	decks := decodeJSON[[]DeckResponse](t, listDecks)
	if len(decks) < 2 {
		t.Fatalf("expected at least 2 decks, got %d", len(decks))
	}

	getDeckBadID := doRawRequest(env.router, http.MethodGet, "/api/decks/not-a-number", "")
	if getDeckBadID.Code != http.StatusBadRequest {
		t.Fatalf("expected get deck invalid id 400, got %d", getDeckBadID.Code)
	}

	getDeckMissing := doRawRequest(env.router, http.MethodGet, "/api/decks/9999999", "")
	if getDeckMissing.Code != http.StatusNotFound {
		t.Fatalf("expected get missing deck 404, got %d", getDeckMissing.Code)
	}

	getDeck := doRawRequest(env.router, http.MethodGet, fmt.Sprintf("/api/decks/%d", createdDeck.ID), "")
	if getDeck.Code != http.StatusOK {
		t.Fatalf("expected get deck 200, got %d (%s)", getDeck.Code, getDeck.Body.String())
	}

	getStatsBadID := doRawRequest(env.router, http.MethodGet, "/api/decks/not-a-number/stats", "")
	if getStatsBadID.Code != http.StatusBadRequest {
		t.Fatalf("expected get stats invalid id 400, got %d", getStatsBadID.Code)
	}

	getStats := doRawRequest(env.router, http.MethodGet, fmt.Sprintf("/api/decks/%d/stats", createdDeck.ID), "")
	if getStats.Code != http.StatusOK {
		t.Fatalf("expected get stats 200, got %d", getStats.Code)
	}
}

func TestAPI_NoteAndCardEndpoints(t *testing.T) {
	env := setupAPITestEnv(t)

	createNoteBadBody := doRawRequest(env.router, http.MethodPost, "/api/notes", "{")
	if createNoteBadBody.Code != http.StatusBadRequest {
		t.Fatalf("expected create note invalid body 400, got %d", createNoteBadBody.Code)
	}

	createNoteMissingFields := doJSONRequest(t, env.router, http.MethodPost, "/api/notes", CreateNoteRequest{})
	if createNoteMissingFields.Code != http.StatusBadRequest {
		t.Fatalf("expected create note missing fields 400, got %d", createNoteMissingFields.Code)
	}

	createUnknownType := doJSONRequest(t, env.router, http.MethodPost, "/api/notes", CreateNoteRequest{
		TypeID:    "UnknownType",
		DeckID:    1,
		FieldVals: map[string]string{"Front": "Q", "Back": "A"},
	})
	if createUnknownType.Code != http.StatusBadRequest {
		t.Fatalf("expected create note unknown type 400, got %d", createUnknownType.Code)
	}

	createGoodNote := doJSONRequest(t, env.router, http.MethodPost, "/api/notes", CreateNoteRequest{
		TypeID: "Basic",
		DeckID: 1,
		FieldVals: map[string]string{
			"Front": "<script>alert(1)</script>Front",
			"Back":  "Back",
		},
		Tags: []string{"<b>tag1</b>", "tag2"},
	})
	if createGoodNote.Code != http.StatusCreated {
		t.Fatalf("expected create note 201, got %d (%s)", createGoodNote.Code, createGoodNote.Body.String())
	}

	var created struct {
		Note  Note   `json:"note"`
		Cards []Card `json:"cards"`
	}
	if err := json.Unmarshal(createGoodNote.Body.Bytes(), &created); err != nil {
		t.Fatalf("failed to decode create note response: %v", err)
	}
	if created.Note.ID == 0 || len(created.Cards) == 0 {
		t.Fatalf("expected note and cards to be created, note=%d cards=%d", created.Note.ID, len(created.Cards))
	}
	if strings.Contains(created.Note.FieldMap["Front"], "<script>") {
		t.Fatalf("expected note front field to be sanitized, got %q", created.Note.FieldMap["Front"])
	}

	getNoteBadID := doRawRequest(env.router, http.MethodGet, "/api/notes/not-a-number", "")
	if getNoteBadID.Code != http.StatusBadRequest {
		t.Fatalf("expected get note invalid id 400, got %d", getNoteBadID.Code)
	}
	getNoteMissing := doRawRequest(env.router, http.MethodGet, "/api/notes/9999999", "")
	if getNoteMissing.Code != http.StatusNotFound {
		t.Fatalf("expected get missing note 404, got %d", getNoteMissing.Code)
	}
	getNote := doRawRequest(env.router, http.MethodGet, fmt.Sprintf("/api/notes/%d", created.Note.ID), "")
	if getNote.Code != http.StatusOK {
		t.Fatalf("expected get note 200, got %d", getNote.Code)
	}

	checkDupBadBody := doRawRequest(env.router, http.MethodPost, "/api/notes/check-duplicate", "{")
	if checkDupBadBody.Code != http.StatusBadRequest {
		t.Fatalf("expected check duplicate bad body 400, got %d", checkDupBadBody.Code)
	}
	checkDupEmpty := doJSONRequest(t, env.router, http.MethodPost, "/api/notes/check-duplicate", CheckDuplicateRequest{
		TypeID:    "Basic",
		FieldName: "Front",
		Value:     "",
	})
	if checkDupEmpty.Code != http.StatusOK {
		t.Fatalf("expected check duplicate empty value 200, got %d", checkDupEmpty.Code)
	}
	dupEmpty := decodeJSON[DuplicateResult](t, checkDupEmpty)
	if dupEmpty.IsDuplicate {
		t.Fatalf("expected empty value duplicate check to be false")
	}

	checkDup := doJSONRequest(t, env.router, http.MethodPost, "/api/notes/check-duplicate", CheckDuplicateRequest{
		TypeID:    "Basic",
		FieldName: "Front",
		Value:     created.Note.FieldMap["Front"],
	})
	if checkDup.Code != http.StatusOK {
		t.Fatalf("expected check duplicate 200, got %d", checkDup.Code)
	}
	dupResult := decodeJSON[DuplicateResult](t, checkDup)
	if !dupResult.IsDuplicate || len(dupResult.Duplicates) == 0 {
		t.Fatalf("expected duplicate to be detected, got %+v", dupResult)
	}

	getDueBadDeck := doRawRequest(env.router, http.MethodGet, "/api/decks/not-a-number/due", "")
	if getDueBadDeck.Code != http.StatusBadRequest {
		t.Fatalf("expected get due invalid deck id 400, got %d", getDueBadDeck.Code)
	}
	getDue := doRawRequest(env.router, http.MethodGet, "/api/decks/1/due?limit=5", "")
	if getDue.Code != http.StatusOK {
		t.Fatalf("expected get due 200, got %d", getDue.Code)
	}

	cardID := created.Cards[0].ID
	getCardBadID := doRawRequest(env.router, http.MethodGet, "/api/cards/not-a-number", "")
	if getCardBadID.Code != http.StatusBadRequest {
		t.Fatalf("expected get card bad id 400, got %d", getCardBadID.Code)
	}
	getCardMissing := doRawRequest(env.router, http.MethodGet, "/api/cards/9999999", "")
	if getCardMissing.Code != http.StatusNotFound {
		t.Fatalf("expected get card missing 404, got %d", getCardMissing.Code)
	}
	getCard := doRawRequest(env.router, http.MethodGet, fmt.Sprintf("/api/cards/%d", cardID), "")
	if getCard.Code != http.StatusOK {
		t.Fatalf("expected get card 200, got %d", getCard.Code)
	}

	answerBadID := doRawRequest(env.router, http.MethodPost, "/api/cards/not-a-number/answer", "{}")
	if answerBadID.Code != http.StatusBadRequest {
		t.Fatalf("expected answer card bad id 400, got %d", answerBadID.Code)
	}
	answerBadBody := doRawRequest(env.router, http.MethodPost, fmt.Sprintf("/api/cards/%d/answer", cardID), "{")
	if answerBadBody.Code != http.StatusBadRequest {
		t.Fatalf("expected answer card bad body 400, got %d", answerBadBody.Code)
	}
	answerBadRating := doJSONRequest(t, env.router, http.MethodPost, fmt.Sprintf("/api/cards/%d/answer", cardID), AnswerCardRequest{Rating: 9})
	if answerBadRating.Code != http.StatusBadRequest {
		t.Fatalf("expected answer card invalid rating 400, got %d", answerBadRating.Code)
	}
	answerMissingCard := doJSONRequest(t, env.router, http.MethodPost, "/api/cards/9999999/answer", AnswerCardRequest{Rating: 3})
	if answerMissingCard.Code != http.StatusInternalServerError {
		t.Fatalf("expected answer card missing 500, got %d", answerMissingCard.Code)
	}
	answerGood := doJSONRequest(t, env.router, http.MethodPost, fmt.Sprintf("/api/cards/%d/answer", cardID), AnswerCardRequest{Rating: 3, TimeTakenMs: 1500})
	if answerGood.Code != http.StatusOK {
		t.Fatalf("expected answer card 200, got %d (%s)", answerGood.Code, answerGood.Body.String())
	}

	updateCardBadID := doRawRequest(env.router, http.MethodPatch, "/api/cards/not-a-number", "{}")
	if updateCardBadID.Code != http.StatusBadRequest {
		t.Fatalf("expected update card bad id 400, got %d", updateCardBadID.Code)
	}
	updateCardBadBody := doRawRequest(env.router, http.MethodPatch, fmt.Sprintf("/api/cards/%d", cardID), "{")
	if updateCardBadBody.Code != http.StatusBadRequest {
		t.Fatalf("expected update card bad body 400, got %d", updateCardBadBody.Code)
	}
	updateCardMissing := doJSONRequest(t, env.router, http.MethodPatch, "/api/cards/9999999", UpdateCardRequest{})
	if updateCardMissing.Code != http.StatusNotFound {
		t.Fatalf("expected update missing card 404, got %d", updateCardMissing.Code)
	}
	badFlag := 9
	updateCardBadFlag := doJSONRequest(t, env.router, http.MethodPatch, fmt.Sprintf("/api/cards/%d", cardID), UpdateCardRequest{
		Flag: &badFlag,
	})
	if updateCardBadFlag.Code != http.StatusBadRequest {
		t.Fatalf("expected update card bad flag 400, got %d", updateCardBadFlag.Code)
	}
	flag := 3
	marked := true
	suspended := true
	updateCard := doJSONRequest(t, env.router, http.MethodPatch, fmt.Sprintf("/api/cards/%d", cardID), UpdateCardRequest{
		Flag:      &flag,
		Marked:    &marked,
		Suspended: &suspended,
	})
	if updateCard.Code != http.StatusOK {
		t.Fatalf("expected update card 200, got %d (%s)", updateCard.Code, updateCard.Body.String())
	}
}

func TestAPI_GetDeckNotesReturnsRecentDistinctNotesForDeck(t *testing.T) {
	env := setupAPITestEnv(t)

	otherDeck := env.collection.NewDeck("Other Deck")
	if err := env.store.CreateDeck(otherDeck); err != nil {
		t.Fatalf("failed to create other deck: %v", err)
	}

	first := createNoteForTest(t, env, CreateNoteRequest{
		TypeID: "Basic",
		DeckID: 1,
		FieldVals: map[string]string{
			"Front": "Older front",
			"Back":  "Older back",
		},
		Tags: []string{"older"},
	}, nil)

	cloze := createNoteForTest(t, env, CreateNoteRequest{
		TypeID: "Cloze",
		DeckID: 1,
		FieldVals: map[string]string{
			"Text":  "{{c1::Newest}} and {{c2::Second}}",
			"Extra": "extra context",
		},
		Tags: []string{"latest", "cloze"},
	}, nil)

	createNoteForTest(t, env, CreateNoteRequest{
		TypeID: "Basic",
		DeckID: otherDeck.ID,
		FieldVals: map[string]string{
			"Front": "Other deck front",
			"Back":  "Other deck back",
		},
	}, nil)

	now := time.Now()
	olderUnix := now.Add(-2 * time.Hour).Unix()
	newerUnix := now.Add(-1 * time.Hour).Unix()
	if _, err := env.store.db.Exec(`UPDATE notes SET created_at = ?, modified_at = ? WHERE id = ?`, olderUnix, olderUnix, first.Note.ID); err != nil {
		t.Fatalf("failed to update first note timestamps: %v", err)
	}
	if _, err := env.store.db.Exec(`UPDATE notes SET created_at = ?, modified_at = ? WHERE id = ?`, newerUnix, newerUnix, cloze.Note.ID); err != nil {
		t.Fatalf("failed to update cloze note timestamps: %v", err)
	}

	rr := doRawRequest(env.router, http.MethodGet, "/api/decks/1/notes?limit=10", "")
	if rr.Code != http.StatusOK {
		t.Fatalf("expected deck notes 200, got %d (%s)", rr.Code, rr.Body.String())
	}

	var response struct {
		Notes []RecentDeckNoteSummary `json:"notes"`
	}
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to decode deck notes response: %v", err)
	}

	if len(response.Notes) != 2 {
		t.Fatalf("expected 2 recent notes for deck 1, got %d", len(response.Notes))
	}
	if response.Notes[0].NoteID != cloze.Note.ID {
		t.Fatalf("expected newest note %d first, got %d", cloze.Note.ID, response.Notes[0].NoteID)
	}
	if response.Notes[0].CardCountInDeck != 2 {
		t.Fatalf("expected cloze note to count as 2 cards, got %d", response.Notes[0].CardCountInDeck)
	}
	if response.Notes[1].NoteID != first.Note.ID {
		t.Fatalf("expected older note %d second, got %d", first.Note.ID, response.Notes[1].NoteID)
	}
	for _, note := range response.Notes {
		if note.FieldPreview == "Other deck front" {
			t.Fatalf("expected notes from other decks to be excluded, got %+v", response.Notes)
		}
	}
}

func TestAPI_GuestPlanLimitsDecksAndNotes(t *testing.T) {
	env := setupAPITestEnv(t)
	guestHeaders := map[string]string{"X-Vutadex-Plan": "guest"}

	createSecondDeck := doJSONRequestWithHeaders(t, env.router, http.MethodPost, "/api/decks", CreateDeckRequest{Name: "Second Deck"}, guestHeaders)
	if createSecondDeck.Code != http.StatusCreated {
		t.Fatalf("expected second deck create to succeed, got %d (%s)", createSecondDeck.Code, createSecondDeck.Body.String())
	}

	thirdDeck := doJSONRequestWithHeaders(t, env.router, http.MethodPost, "/api/decks", CreateDeckRequest{Name: "Third Deck"}, guestHeaders)
	if thirdDeck.Code != http.StatusForbidden {
		t.Fatalf("expected third deck create to fail with 403, got %d (%s)", thirdDeck.Code, thirdDeck.Body.String())
	}
	deckErr := decodeJSON[APIErrorResponse](t, thirdDeck)
	if deckErr.Code != "plan_limit_exceeded" {
		t.Fatalf("expected plan_limit_exceeded for deck limit, got %+v", deckErr)
	}

	for i := 0; i < 10; i++ {
		createNote := doJSONRequestWithHeaders(t, env.router, http.MethodPost, "/api/notes", CreateNoteRequest{
			TypeID: "Basic",
			DeckID: 1,
			FieldVals: map[string]string{
				"Front": fmt.Sprintf("Front %d", i),
				"Back":  fmt.Sprintf("Back %d", i),
			},
		}, guestHeaders)
		if createNote.Code != http.StatusCreated {
			t.Fatalf("expected note %d create to succeed, got %d (%s)", i+1, createNote.Code, createNote.Body.String())
		}
	}

	eleventhNote := doJSONRequestWithHeaders(t, env.router, http.MethodPost, "/api/notes", CreateNoteRequest{
		TypeID: "Basic",
		DeckID: 1,
		FieldVals: map[string]string{
			"Front": "Limit front",
			"Back":  "Limit back",
		},
	}, guestHeaders)
	if eleventhNote.Code != http.StatusForbidden {
		t.Fatalf("expected eleventh note to fail with 403, got %d (%s)", eleventhNote.Code, eleventhNote.Body.String())
	}
	noteErr := decodeJSON[APIErrorResponse](t, eleventhNote)
	if noteErr.Code != "plan_limit_exceeded" {
		t.Fatalf("expected plan_limit_exceeded for note limit, got %+v", noteErr)
	}
}

func TestAPI_NoteTypeFieldTemplateEmptyAndBackupEndpoints(t *testing.T) {
	env := setupAPITestEnv(t)

	listNoteTypes := doRawRequest(env.router, http.MethodGet, "/api/note-types", "")
	if listNoteTypes.Code != http.StatusOK {
		t.Fatalf("expected list note types 200, got %d", listNoteTypes.Code)
	}
	nts := decodeJSON[[]NoteTypeResponse](t, listNoteTypes)
	if len(nts) == 0 {
		t.Fatal("expected built-in note types")
	}
	names := make([]string, 0, len(nts))
	for _, nt := range nts {
		names = append(names, nt.Name)
	}
	if !sort.StringsAreSorted(names) {
		t.Fatalf("expected note types sorted, got %v", names)
	}

	getNoteTypeMissing := doRawRequest(env.router, http.MethodGet, "/api/note-types/NoSuchType", "")
	if getNoteTypeMissing.Code != http.StatusNotFound {
		t.Fatalf("expected get note type missing 404, got %d", getNoteTypeMissing.Code)
	}
	getNoteType := doRawRequest(env.router, http.MethodGet, "/api/note-types/Basic", "")
	if getNoteType.Code != http.StatusOK {
		t.Fatalf("expected get note type 200, got %d", getNoteType.Code)
	}

	addFieldMissingType := doJSONRequest(t, env.router, http.MethodPost, "/api/note-types/NoSuchType/fields", AddFieldRequest{FieldName: "F"})
	if addFieldMissingType.Code != http.StatusNotFound {
		t.Fatalf("expected add field missing note type 404, got %d", addFieldMissingType.Code)
	}
	addFieldBadBody := doRawRequest(env.router, http.MethodPost, "/api/note-types/Basic/fields", "{")
	if addFieldBadBody.Code != http.StatusBadRequest {
		t.Fatalf("expected add field bad body 400, got %d", addFieldBadBody.Code)
	}
	addFieldEmpty := doJSONRequest(t, env.router, http.MethodPost, "/api/note-types/Basic/fields", AddFieldRequest{})
	if addFieldEmpty.Code != http.StatusBadRequest {
		t.Fatalf("expected add field empty name 400, got %d", addFieldEmpty.Code)
	}
	addFieldReserved := doJSONRequest(t, env.router, http.MethodPost, "/api/note-types/Basic/fields", AddFieldRequest{FieldName: "Tags"})
	if addFieldReserved.Code != http.StatusBadRequest {
		t.Fatalf("expected add field reserved name 400, got %d", addFieldReserved.Code)
	}
	addFieldDuplicate := doJSONRequest(t, env.router, http.MethodPost, "/api/note-types/Basic/fields", AddFieldRequest{FieldName: "Front"})
	if addFieldDuplicate.Code != http.StatusBadRequest {
		t.Fatalf("expected add field duplicate 400, got %d", addFieldDuplicate.Code)
	}
	pos := 1
	addField := doJSONRequest(t, env.router, http.MethodPost, "/api/note-types/Basic/fields", AddFieldRequest{
		FieldName: "MiddleField",
		Position:  &pos,
	})
	if addField.Code != http.StatusOK {
		t.Fatalf("expected add field success 200, got %d (%s)", addField.Code, addField.Body.String())
	}

	renameBadBody := doRawRequest(env.router, http.MethodPatch, "/api/note-types/Basic/fields/rename", "{")
	if renameBadBody.Code != http.StatusBadRequest {
		t.Fatalf("expected rename bad body 400, got %d", renameBadBody.Code)
	}
	renameMissing := doJSONRequest(t, env.router, http.MethodPatch, "/api/note-types/Basic/fields/rename", RenameFieldRequest{})
	if renameMissing.Code != http.StatusBadRequest {
		t.Fatalf("expected rename missing fields 400, got %d", renameMissing.Code)
	}
	renameReserved := doJSONRequest(t, env.router, http.MethodPatch, "/api/note-types/Basic/fields/rename", RenameFieldRequest{
		OldName: "MiddleField",
		NewName: "Deck",
	})
	if renameReserved.Code != http.StatusBadRequest {
		t.Fatalf("expected rename reserved target 400, got %d", renameReserved.Code)
	}
	renameMissingField := doJSONRequest(t, env.router, http.MethodPatch, "/api/note-types/Basic/fields/rename", RenameFieldRequest{
		OldName: "DoesNotExist",
		NewName: "X",
	})
	if renameMissingField.Code != http.StatusNotFound {
		t.Fatalf("expected rename missing field 404, got %d", renameMissingField.Code)
	}
	renameDuplicate := doJSONRequest(t, env.router, http.MethodPatch, "/api/note-types/Basic/fields/rename", RenameFieldRequest{
		OldName: "Back",
		NewName: "Front",
	})
	if renameDuplicate.Code != http.StatusBadRequest {
		t.Fatalf("expected rename duplicate target 400, got %d", renameDuplicate.Code)
	}
	renameOK := doJSONRequest(t, env.router, http.MethodPatch, "/api/note-types/Basic/fields/rename", RenameFieldRequest{
		OldName: "MiddleField",
		NewName: "MiddleRenamed",
	})
	if renameOK.Code != http.StatusOK {
		t.Fatalf("expected rename success 200, got %d (%s)", renameOK.Code, renameOK.Body.String())
	}

	removeBadBody := doRawRequest(env.router, http.MethodDelete, "/api/note-types/Basic/fields", "{")
	if removeBadBody.Code != http.StatusBadRequest {
		t.Fatalf("expected remove field bad body 400, got %d", removeBadBody.Code)
	}
	removeEmpty := doJSONRequest(t, env.router, http.MethodDelete, "/api/note-types/Basic/fields", RemoveFieldRequest{})
	if removeEmpty.Code != http.StatusBadRequest {
		t.Fatalf("expected remove field empty name 400, got %d", removeEmpty.Code)
	}
	removeMissing := doJSONRequest(t, env.router, http.MethodDelete, "/api/note-types/Basic/fields", RemoveFieldRequest{FieldName: "Nope"})
	if removeMissing.Code != http.StatusNotFound {
		t.Fatalf("expected remove missing field 404, got %d", removeMissing.Code)
	}

	// Exercise "Cannot remove the last field" branch with a single-field note type.
	single := NoteType{
		Name:      "SingleFieldType",
		Fields:    []string{"Only"},
		Templates: []CardTemplate{{Name: "Card 1", QFmt: "{{Only}}", AFmt: "{{Only}}"}},
	}
	env.collection.NoteTypes[single.Name] = single
	if err := env.store.CreateNoteType("default", &single); err != nil {
		t.Fatalf("failed to create single-field note type: %v", err)
	}
	removeLast := doJSONRequest(t, env.router, http.MethodDelete, "/api/note-types/SingleFieldType/fields", RemoveFieldRequest{FieldName: "Only"})
	if removeLast.Code != http.StatusBadRequest {
		t.Fatalf("expected remove last field 400, got %d", removeLast.Code)
	}

	removeOK := doJSONRequest(t, env.router, http.MethodDelete, "/api/note-types/Basic/fields", RemoveFieldRequest{FieldName: "MiddleRenamed"})
	if removeOK.Code != http.StatusOK {
		t.Fatalf("expected remove field success 200, got %d", removeOK.Code)
	}

	reorderBadBody := doRawRequest(env.router, http.MethodPut, "/api/note-types/Basic/fields/reorder", "{")
	if reorderBadBody.Code != http.StatusBadRequest {
		t.Fatalf("expected reorder bad body 400, got %d", reorderBadBody.Code)
	}
	reorderMismatch := doJSONRequest(t, env.router, http.MethodPut, "/api/note-types/Basic/fields/reorder", ReorderFieldsRequest{
		Fields: []string{"Front"},
	})
	if reorderMismatch.Code != http.StatusBadRequest {
		t.Fatalf("expected reorder mismatch 400, got %d", reorderMismatch.Code)
	}
	reorderUnknown := doJSONRequest(t, env.router, http.MethodPut, "/api/note-types/Basic/fields/reorder", ReorderFieldsRequest{
		Fields: []string{"Front", "Unknown"},
	})
	if reorderUnknown.Code != http.StatusBadRequest {
		t.Fatalf("expected reorder unknown field 400, got %d", reorderUnknown.Code)
	}
	reorderOK := doJSONRequest(t, env.router, http.MethodPut, "/api/note-types/Basic/fields/reorder", ReorderFieldsRequest{
		Fields: []string{"Back", "Front"},
	})
	if reorderOK.Code != http.StatusOK {
		t.Fatalf("expected reorder success 200, got %d", reorderOK.Code)
	}

	sortBadBody := doRawRequest(env.router, http.MethodPut, "/api/note-types/Basic/sort-field", "{")
	if sortBadBody.Code != http.StatusBadRequest {
		t.Fatalf("expected sort field bad body 400, got %d", sortBadBody.Code)
	}
	sortInvalid := doJSONRequest(t, env.router, http.MethodPut, "/api/note-types/Basic/sort-field", SetSortFieldRequest{FieldIndex: 9})
	if sortInvalid.Code != http.StatusBadRequest {
		t.Fatalf("expected sort field invalid index 400, got %d", sortInvalid.Code)
	}
	sortOK := doJSONRequest(t, env.router, http.MethodPut, "/api/note-types/Basic/sort-field", SetSortFieldRequest{FieldIndex: 1})
	if sortOK.Code != http.StatusOK {
		t.Fatalf("expected sort field success 200, got %d", sortOK.Code)
	}

	setFieldOptionsBadBody := doRawRequest(env.router, http.MethodPut, "/api/note-types/Basic/fields/options", "{")
	if setFieldOptionsBadBody.Code != http.StatusBadRequest {
		t.Fatalf("expected set field options bad body 400, got %d", setFieldOptionsBadBody.Code)
	}
	setFieldOptionsMissing := doJSONRequest(t, env.router, http.MethodPut, "/api/note-types/Basic/fields/options", SetFieldOptionsRequest{
		FieldName: "Nope",
		Options:   FieldOptions{Font: "Arial"},
	})
	if setFieldOptionsMissing.Code != http.StatusBadRequest {
		t.Fatalf("expected set field options missing field 400, got %d", setFieldOptionsMissing.Code)
	}
	setFieldOptionsOK := doJSONRequest(t, env.router, http.MethodPut, "/api/note-types/Basic/fields/options", SetFieldOptionsRequest{
		FieldName: "Back",
		Options:   FieldOptions{Font: "Arial", FontSize: 20, RTL: true, HTMLEditor: true},
	})
	if setFieldOptionsOK.Code != http.StatusOK {
		t.Fatalf("expected set field options success 200, got %d", setFieldOptionsOK.Code)
	}

	// Create a note so template regeneration has data to process.
	createNote := doJSONRequest(t, env.router, http.MethodPost, "/api/notes", CreateNoteRequest{
		TypeID: "Basic",
		DeckID: 1,
		FieldVals: map[string]string{
			"Front": "Regenerate Me",
			"Back":  "Answer",
		},
	})
	if createNote.Code != http.StatusCreated {
		t.Fatalf("expected create note for template update 201, got %d (%s)", createNote.Code, createNote.Body.String())
	}

	var noteCreated struct {
		Note  Note   `json:"note"`
		Cards []Card `json:"cards"`
	}
	if err := json.Unmarshal(createNote.Body.Bytes(), &noteCreated); err != nil {
		t.Fatalf("failed to decode create note: %v", err)
	}

	// Add an orphan card to exercise card cleanup path during regeneration.
	orphan := &Card{
		ID:           777777,
		NoteID:       noteCreated.Note.ID,
		DeckID:       1,
		TemplateName: "Ghost Template",
		Ordinal:      0,
		Front:        "ghost",
		Back:         "ghost",
		SRS:          newDueNow(time.Now()),
		USN:          1,
	}
	if err := env.store.CreateCard(orphan); err != nil {
		t.Fatalf("failed to create orphan card: %v", err)
	}

	updateTemplateMissingType := doJSONRequest(t, env.router, http.MethodPatch, "/api/note-types/Nope/templates/Card%201", UpdateTemplateRequest{})
	if updateTemplateMissingType.Code != http.StatusNotFound {
		t.Fatalf("expected update template missing type 404, got %d", updateTemplateMissingType.Code)
	}
	updateTemplateBadBody := doRawRequest(env.router, http.MethodPatch, "/api/note-types/Basic/templates/Card%201", "{")
	if updateTemplateBadBody.Code != http.StatusBadRequest {
		t.Fatalf("expected update template bad body 400, got %d", updateTemplateBadBody.Code)
	}
	updateTemplateMissing := doJSONRequest(t, env.router, http.MethodPatch, "/api/note-types/Basic/templates/Nope", UpdateTemplateRequest{})
	if updateTemplateMissing.Code != http.StatusNotFound {
		t.Fatalf("expected update missing template 404, got %d", updateTemplateMissing.Code)
	}
	newQ := "<div>New {{Front}}</div>"
	newA := "<div>New {{Back}}</div>"
	updateTemplate := doJSONRequest(t, env.router, http.MethodPatch, "/api/note-types/Basic/templates/Card%201", UpdateTemplateRequest{
		QFmt: &newQ,
		AFmt: &newA,
	})
	if updateTemplate.Code != http.StatusOK {
		t.Fatalf("expected update template success 200, got %d (%s)", updateTemplate.Code, updateTemplate.Body.String())
	}

	// Build deterministic empty-card scenarios:
	// 1) A cloze card where c1 no longer exists in note text.
	clozeNote, clozeCards, err := env.collection.AddNote(1, "Cloze", map[string]string{
		"Text":  "{{c1::gap}}",
		"Extra": "extra",
	}, time.Now())
	if err != nil {
		t.Fatalf("failed to add cloze note: %v", err)
	}
	if err := env.store.CreateNote("default", &clozeNote); err != nil {
		t.Fatalf("failed to persist cloze note: %v", err)
	}
	for _, c := range clozeCards {
		if err := env.store.CreateCard(c); err != nil {
			t.Fatalf("failed to persist cloze card: %v", err)
		}
	}
	clozeNote.FieldMap["Text"] = "no cloze now"
	clozeNote.ModifiedAt = time.Now()
	if err := env.store.UpdateNote(&clozeNote); err != nil {
		t.Fatalf("failed to update cloze note text: %v", err)
	}
	env.collection.Notes[clozeNote.ID] = clozeNote

	// 2) A non-cloze card with effectively empty content.
	emptyCard := &Card{
		ID:           888888,
		NoteID:       noteCreated.Note.ID,
		DeckID:       1,
		TemplateName: "Card 1",
		Ordinal:      0,
		Front:        "<div> </div>",
		Back:         "<p></p>",
		SRS:          newDueNow(time.Now()),
		USN:          1,
	}
	if err := env.store.CreateCard(emptyCard); err != nil {
		t.Fatalf("failed to create empty content card: %v", err)
	}

	emptyCards := doRawRequest(env.router, http.MethodGet, "/api/cards/empty", "")
	if emptyCards.Code != http.StatusOK {
		t.Fatalf("expected find empty cards 200, got %d (%s)", emptyCards.Code, emptyCards.Body.String())
	}
	emptyResp := decodeJSON[EmptyCardsResponse](t, emptyCards)
	if emptyResp.Count == 0 {
		t.Fatalf("expected at least one empty card, got %+v", emptyResp)
	}

	deleteEmptyBadBody := doRawRequest(env.router, http.MethodPost, "/api/cards/empty/delete", "{")
	if deleteEmptyBadBody.Code != http.StatusBadRequest {
		t.Fatalf("expected delete empty cards bad body 400, got %d", deleteEmptyBadBody.Code)
	}
	deleteEmptyNoIDs := doJSONRequest(t, env.router, http.MethodPost, "/api/cards/empty/delete", DeleteEmptyCardsRequest{})
	if deleteEmptyNoIDs.Code != http.StatusBadRequest {
		t.Fatalf("expected delete empty cards no IDs 400, got %d", deleteEmptyNoIDs.Code)
	}
	deleteEmpty := doJSONRequest(t, env.router, http.MethodPost, "/api/cards/empty/delete", DeleteEmptyCardsRequest{
		CardIDs: []int64{emptyResp.EmptyCards[0].CardID, 4242424242},
	})
	if deleteEmpty.Code != http.StatusOK {
		t.Fatalf("expected delete empty cards 200, got %d (%s)", deleteEmpty.Code, deleteEmpty.Body.String())
	}
	deleteResp := decodeJSON[DeleteEmptyCardsResponse](t, deleteEmpty)
	if deleteResp.Deleted == 0 {
		t.Fatalf("expected at least one deleted empty card, got %+v", deleteResp)
	}
	// SQLite delete is idempotent; deleting missing IDs does not return an error.
	if len(deleteResp.Failed) != 0 {
		t.Fatalf("expected no delete errors, got %+v", deleteResp)
	}

	createBackup := doRawRequest(env.router, http.MethodPost, "/api/backups", "{}")
	if createBackup.Code != http.StatusCreated {
		t.Fatalf("expected create backup 201, got %d (%s)", createBackup.Code, createBackup.Body.String())
	}
	var backupResp map[string]string
	if err := json.Unmarshal(createBackup.Body.Bytes(), &backupResp); err != nil {
		t.Fatalf("failed to decode backup response: %v", err)
	}
	backupPath := backupResp["backupPath"]
	if backupPath == "" {
		t.Fatalf("expected backup path in response, got %+v", backupResp)
	}

	listBackups := doRawRequest(env.router, http.MethodGet, "/api/backups", "")
	if listBackups.Code != http.StatusOK {
		t.Fatalf("expected list backups 200, got %d", listBackups.Code)
	}

	restoreBadBody := doRawRequest(env.router, http.MethodPost, "/api/backups/restore", "{")
	if restoreBadBody.Code != http.StatusBadRequest {
		t.Fatalf("expected restore backup bad body 400, got %d", restoreBadBody.Code)
	}
	restoreNoPath := doJSONRequest(t, env.router, http.MethodPost, "/api/backups/restore", RestoreBackupRequest{})
	if restoreNoPath.Code != http.StatusBadRequest {
		t.Fatalf("expected restore backup missing path 400, got %d", restoreNoPath.Code)
	}
	restore := doJSONRequest(t, env.router, http.MethodPost, "/api/backups/restore", RestoreBackupRequest{BackupPath: backupPath})
	if restore.Code != http.StatusOK {
		t.Fatalf("expected restore backup warning response 200, got %d (%s)", restore.Code, restore.Body.String())
	}
}

func TestRegenerateCardsForNoteType_CoversUpdateCreateDeletePaths(t *testing.T) {
	env := setupAPITestEnv(t)

	// Create a basic note + persisted card.
	note, cards, err := env.collection.AddNote(1, "Basic", map[string]string{
		"Front": "Front A",
		"Back":  "Back A",
	}, time.Now())
	if err != nil {
		t.Fatalf("failed to add note: %v", err)
	}
	if err := env.store.CreateNote("default", &note); err != nil {
		t.Fatalf("failed to persist note: %v", err)
	}
	for _, c := range cards {
		if err := env.store.CreateCard(c); err != nil {
			t.Fatalf("failed to persist card: %v", err)
		}
	}
	if len(cards) == 0 {
		t.Fatal("expected at least one generated card")
	}

	// Add orphan card that should be removed by regeneration.
	orphan := &Card{
		ID:           999991,
		NoteID:       note.ID,
		DeckID:       1,
		TemplateName: "Orphan",
		Ordinal:      0,
		Front:        "orphan-front",
		Back:         "orphan-back",
		SRS:          newDueNow(time.Now()),
		USN:          1,
	}
	if err := env.store.CreateCard(orphan); err != nil {
		t.Fatalf("failed to create orphan card: %v", err)
	}

	// Modify the note type templates to:
	// - update existing card 1 content
	// - add a second card template (create path)
	basic := env.collection.NoteTypes["Basic"]
	basic.Templates = []CardTemplate{
		{
			Name: "Card 1",
			QFmt: "Updated Front: {{Front}}",
			AFmt: "Updated Back: {{Back}}",
		},
		{
			Name: "Card 2",
			QFmt: "Second Card: {{Back}}",
			AFmt: "Second Answer: {{Front}}",
		},
	}
	env.collection.NoteTypes["Basic"] = basic
	if err := env.store.UpdateNoteType("default", &basic); err != nil {
		t.Fatalf("failed to persist updated note type: %v", err)
	}

	if err := env.handler.regenerateCardsForNoteType("Basic"); err != nil {
		t.Fatalf("expected regeneration to succeed, got %v", err)
	}

	updatedCard, err := env.store.GetCard(cards[0].ID)
	if err != nil {
		t.Fatalf("failed to get updated existing card: %v", err)
	}
	if !strings.Contains(updatedCard.Front, "Updated Front") {
		t.Fatalf("expected existing card content to be refreshed, got %q", updatedCard.Front)
	}

	if _, err := env.store.GetCard(orphan.ID); err == nil {
		t.Fatalf("expected orphan card %d to be deleted", orphan.ID)
	}

	byNote, err := env.store.GetCardsByNote(note.ID)
	if err != nil {
		t.Fatalf("failed to read cards by note after regeneration: %v", err)
	}
	foundCard2 := false
	for _, c := range byNote {
		if c.TemplateName == "Card 2" {
			foundCard2 = true
			break
		}
	}
	if !foundCard2 {
		t.Fatalf("expected regeneration to create Card 2 template output, got %#v", byNote)
	}
}

func TestParseIDParamAndRespondJSONHelpers(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/unused", nil)
	routeCtx := chi.NewRouteContext()
	routeCtx.URLParams.Add("id", "42")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))

	id, err := parseIDParam(req, "id")
	if err != nil {
		t.Fatalf("expected parse id to succeed: %v", err)
	}
	if id != 42 {
		t.Fatalf("expected parsed id=42, got %d", id)
	}

	rr := httptest.NewRecorder()
	respondJSON(rr, http.StatusCreated, map[string]string{"ok": "yes"})
	if rr.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", rr.Code)
	}
	if got := rr.Header().Get("Content-Type"); got != "application/json" {
		t.Fatalf("expected content type json, got %q", got)
	}
}

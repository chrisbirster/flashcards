package main

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	fsrs "github.com/open-spaced-repetition/go-fsrs/v3"
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

type authenticatedTestClient struct {
	router     http.Handler
	authCookie string
	user       *User
	workspace  *Workspace
	session    *SessionRecord
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
	return setupAPITestEnvWithConfig(t, mustLocalAppConfig())
}

func setupAPITestEnvWithConfig(t *testing.T, cfg AppConfig) *apiTestEnv {
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

	handler := NewAPIHandlerWithConfig(store, col, NewBackupManager(dbPath, backupDir, store), cfg, NewEmailSender(cfg))
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

type roundTripFunc func(*http.Request) (*http.Response, error)

func (fn roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return fn(req)
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

func createAuthenticatedTestClient(t *testing.T, env *apiTestEnv, email, displayName string) authenticatedTestClient {
	t.Helper()

	now := time.Now()
	user := &User{
		ID:          newID("usr"),
		Email:       email,
		DisplayName: displayName,
		LastLoginAt: now,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := env.store.CreateUser(user); err != nil {
		t.Fatalf("failed to create test user %s: %v", email, err)
	}

	workspace := &Workspace{
		ID:           newID("ws"),
		Name:         displayName + " Workspace",
		Slug:         strings.ToLower(strings.ReplaceAll(displayName, " ", "-")),
		CollectionID: "default",
		OwnerUserID:  user.ID,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	if err := env.store.CreateWorkspaceRecord(workspace); err != nil {
		t.Fatalf("failed to create test workspace for %s: %v", email, err)
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
	if err := env.store.CreateSessionRecord(session); err != nil {
		t.Fatalf("failed to create test session for %s: %v", email, err)
	}

	authCookie := fmt.Sprintf("%s=%s", sessionCookieName, session.ID)
	return authenticatedTestClient{
		router:     newTestAPIRouter(env.handler, authCookie),
		authCookie: authCookie,
		user:       user,
		workspace:  workspace,
		session:    session,
	}
}

func createAuthenticatedIsolatedTestClient(t *testing.T, env *apiTestEnv, email, displayName string) authenticatedTestClient {
	t.Helper()

	now := time.Now()
	user := &User{
		ID:          newID("usr"),
		Email:       email,
		DisplayName: displayName,
		LastLoginAt: now,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := env.store.CreateUser(user); err != nil {
		t.Fatalf("failed to create isolated test user %s: %v", email, err)
	}

	collectionID := newID("col")
	collection := NewCollection()
	if err := env.store.CreateCollectionRecord(collectionID, displayName+" Collection", collection); err != nil {
		t.Fatalf("failed to create isolated collection for %s: %v", email, err)
	}

	workspace := &Workspace{
		ID:           newID("ws"),
		Name:         displayName + " Workspace",
		Slug:         strings.ToLower(strings.ReplaceAll(displayName, " ", "-")) + "-isolated",
		CollectionID: collectionID,
		OwnerUserID:  user.ID,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	if err := env.store.CreateWorkspaceRecord(workspace); err != nil {
		t.Fatalf("failed to create isolated workspace for %s: %v", email, err)
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
	if err := env.store.CreateSessionRecord(session); err != nil {
		t.Fatalf("failed to create isolated session for %s: %v", email, err)
	}

	authCookie := fmt.Sprintf("%s=%s", sessionCookieName, session.ID)
	return authenticatedTestClient{
		router:     newTestAPIRouter(env.handler, authCookie),
		authCookie: authCookie,
		user:       user,
		workspace:  workspace,
		session:    session,
	}
}

func activateWorkspaceSubscriptionForTest(t *testing.T, env *apiTestEnv, workspaceID string, plan Plan) {
	t.Helper()

	now := time.Now()
	subscription := &Subscription{
		ID:               newID("sub"),
		WorkspaceID:      workspaceID,
		Plan:             plan,
		Status:           "active",
		Provider:         "test",
		CreatedAt:        now,
		UpdatedAt:        now,
		CurrentPeriodEnd: now.Add(30 * 24 * time.Hour),
	}
	if err := env.store.UpsertSubscription(subscription); err != nil {
		t.Fatalf("failed to create active workspace subscription for %s: %v", workspaceID, err)
	}
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
	return doRawRequestWithHeaders(router, method, path, body, nil)
}

func doRawRequestWithHeaders(router http.Handler, method, path string, body string, headers map[string]string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	for key, value := range headers {
		req.Header.Set(key, value)
	}
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
	if !session.User.Onboarding {
		t.Fatalf("expected new OTP user to start in onboarding, got %+v", session.User)
	}
}

func TestAPI_OnboardingPlanSelectionClearsFlagAndCreatesOrganization(t *testing.T) {
	env := setupAPITestEnv(t)
	emailStub := &otpEmailStub{}
	env.handler.emailSender = emailStub

	requestRR := doJSONRequestWithHeaders(t, env.router, http.MethodPost, "/api/auth/otp/request", map[string]string{
		"email": "onboarding@example.com",
	}, map[string]string{
		"X-Test-No-Auth": "1",
	})
	if requestRR.Code != http.StatusOK {
		t.Fatalf("expected OTP request 200, got %d (%s)", requestRR.Code, requestRR.Body.String())
	}

	verifyRR := doJSONRequestWithHeaders(t, env.router, http.MethodPost, "/api/auth/otp/verify", map[string]string{
		"email": "onboarding@example.com",
		"code":  emailStub.lastCode,
	}, map[string]string{
		"X-Test-No-Auth": "1",
	})
	if verifyRR.Code != http.StatusOK {
		t.Fatalf("expected OTP verify 200, got %d (%s)", verifyRR.Code, verifyRR.Body.String())
	}
	session := decodeJSON[AuthSessionResponse](t, verifyRR)
	if session.User == nil || !session.User.Onboarding {
		t.Fatalf("expected verified user to remain in onboarding before plan selection, got %+v", session.User)
	}

	cookies := verifyRR.Result().Cookies()
	if len(cookies) == 0 || cookies[0].Name != sessionCookieName {
		t.Fatalf("expected onboarding session cookie to be set, got %+v", cookies)
	}
	authCookie := fmt.Sprintf("%s=%s", sessionCookieName, cookies[0].Value)
	userRouter := newTestAPIRouter(env.handler, authCookie)

	selectPlanRR := doJSONRequest(t, userRouter, http.MethodPost, "/api/onboarding/plan", UpdateWorkspacePlanRequest{
		Plan: PlanTeam,
	})
	if selectPlanRR.Code != http.StatusOK {
		t.Fatalf("expected onboarding plan selection 200, got %d (%s)", selectPlanRR.Code, selectPlanRR.Body.String())
	}
	updated := decodeJSON[AuthSessionResponse](t, selectPlanRR)
	if updated.User == nil || updated.User.Onboarding {
		t.Fatalf("expected onboarding to clear after plan selection, got %+v", updated.User)
	}
	if updated.Workspace == nil || updated.Workspace.OrganizationID == "" {
		t.Fatalf("expected team selection to attach an organization to the workspace, got %+v", updated.Workspace)
	}
	if updated.Organization == nil || updated.OrganizationMember == nil || updated.OrganizationMember.Role != "owner" {
		t.Fatalf("expected team selection to create owner membership, got org=%+v member=%+v", updated.Organization, updated.OrganizationMember)
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
	if createdDeck.NewCardsPerDay != defaultNewCardsPerDay {
		t.Fatalf("expected default newCardsPerDay=%d, got %d", defaultNewCardsPerDay, createdDeck.NewCardsPerDay)
	}
	if createdDeck.ReviewsPerDay != defaultReviewsPerDay {
		t.Fatalf("expected default reviewsPerDay=%d, got %d", defaultReviewsPerDay, createdDeck.ReviewsPerDay)
	}
	if createdDeck.PriorityOrder <= 0 {
		t.Fatalf("expected created deck to have positive priority order, got %d", createdDeck.PriorityOrder)
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

func TestAPI_DeckWorkloadPolicy_DefaultCapPauseRuleAndPriority(t *testing.T) {
	env := setupAPITestEnv(t)
	sessionID := strings.TrimPrefix(env.authCookie, sessionCookieName+"=")
	sessionRecord, err := env.store.GetSessionRecord(sessionID)
	if err != nil {
		t.Fatalf("failed to load current session: %v", err)
	}
	activateWorkspaceSubscriptionForTest(t, env, sessionRecord.WorkspaceID, PlanPro)

	for i := 0; i < 25; i++ {
		createNoteForTest(t, env, CreateNoteRequest{
			TypeID: "Basic",
			DeckID: 1,
			FieldVals: map[string]string{
				"Front": fmt.Sprintf("Default cap %d", i),
				"Back":  "Answer",
			},
		}, nil)
	}

	dueWithDefaults := doRawRequest(env.router, http.MethodGet, "/api/decks/1/due?limit=100", "")
	if dueWithDefaults.Code != http.StatusOK {
		t.Fatalf("expected due request 200, got %d (%s)", dueWithDefaults.Code, dueWithDefaults.Body.String())
	}
	defaultDueCards := decodeJSON[[]Card](t, dueWithDefaults)
	if len(defaultDueCards) != defaultNewCardsPerDay {
		t.Fatalf("expected default new-card cap to return %d cards, got %d", defaultNewCardsPerDay, len(defaultDueCards))
	}

	for i := 0; i < 5; i++ {
		if _, err := env.store.db.Exec(`
			UPDATE card_review_states
			SET state = ?, due = ?
			WHERE user_id = ? AND card_id = ?
		`, int(fsrs.Review), time.Now().Add(-time.Hour).Unix(), sessionRecord.UserID, defaultDueCards[i].ID); err != nil {
			t.Fatalf("failed to promote card %d into review backlog: %v", defaultDueCards[i].ID, err)
		}
	}

	reviewsPerDay := 3
	newCardsPerDay := 20
	priorityOrder := 1
	updateDeckRR := doJSONRequest(t, env.router, http.MethodPatch, "/api/decks/1", UpdateDeckRequest{
		ReviewsPerDay:  &reviewsPerDay,
		NewCardsPerDay: &newCardsPerDay,
		PriorityOrder:  &priorityOrder,
	})
	if updateDeckRR.Code != http.StatusOK {
		t.Fatalf("expected deck settings update 200, got %d (%s)", updateDeckRR.Code, updateDeckRR.Body.String())
	}
	updatedDeck := decodeJSON[DeckResponse](t, updateDeckRR)
	if !updatedDeck.NewCardsPaused {
		t.Fatalf("expected deck to pause new cards when backlog exceeds review cap, got %+v", updatedDeck)
	}
	if updatedDeck.DueReviewBacklog != 5 {
		t.Fatalf("expected due review backlog=5, got %d", updatedDeck.DueReviewBacklog)
	}

	dueWithBacklog := doRawRequest(env.router, http.MethodGet, "/api/decks/1/due?limit=100", "")
	if dueWithBacklog.Code != http.StatusOK {
		t.Fatalf("expected due request with backlog 200, got %d (%s)", dueWithBacklog.Code, dueWithBacklog.Body.String())
	}
	backlogCards := decodeJSON[[]Card](t, dueWithBacklog)
	if len(backlogCards) != reviewsPerDay {
		t.Fatalf("expected due cards to stop at review cap=%d while new cards are paused, got %d", reviewsPerDay, len(backlogCards))
	}

	for _, card := range backlogCards {
		var state int
		if err := env.store.db.QueryRow(`SELECT state FROM card_review_states WHERE user_id = ? AND card_id = ?`, sessionRecord.UserID, card.ID).Scan(&state); err != nil {
			t.Fatalf("failed to read review state for card %d: %v", card.ID, err)
		}
		if state != int(fsrs.Review) {
			t.Fatalf("expected paused-new queue to return only review cards, got state=%d for card %d", state, card.ID)
		}
	}

	secondDeckRR := doJSONRequest(t, env.router, http.MethodPost, "/api/decks", CreateDeckRequest{Name: "Later deck"})
	if secondDeckRR.Code != http.StatusCreated {
		t.Fatalf("expected second deck create 201, got %d (%s)", secondDeckRR.Code, secondDeckRR.Body.String())
	}
	secondDeck := decodeJSON[DeckResponse](t, secondDeckRR)
	latePriority := 9
	updateSecondDeckRR := doJSONRequest(t, env.router, http.MethodPatch, fmt.Sprintf("/api/decks/%d", secondDeck.ID), UpdateDeckRequest{
		PriorityOrder: &latePriority,
	})
	if updateSecondDeckRR.Code != http.StatusOK {
		t.Fatalf("expected second deck priority update 200, got %d (%s)", updateSecondDeckRR.Code, updateSecondDeckRR.Body.String())
	}

	listDecksRR := doRawRequest(env.router, http.MethodGet, "/api/decks", "")
	if listDecksRR.Code != http.StatusOK {
		t.Fatalf("expected list decks 200, got %d (%s)", listDecksRR.Code, listDecksRR.Body.String())
	}
	listedDecks := decodeJSON[[]DeckResponse](t, listDecksRR)
	if len(listedDecks) < 2 {
		t.Fatalf("expected at least two decks, got %d", len(listedDecks))
	}
	if listedDecks[0].ID != 1 {
		t.Fatalf("expected lower priority order deck to surface first, got deck ID %d", listedDecks[0].ID)
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

	cfg := mustLocalAppConfig()
	cfg.OpenAI.APIKey = ""
	aiEnv := setupAPITestEnvWithConfig(t, cfg)
	suggestionRR := doJSONRequest(t, aiEnv.router, http.MethodPost, "/api/ai/card-suggestions", GenerateAICardSuggestionsRequest{
		SourceText: "Mitochondria: The powerhouse of the cell\nATP: The cell's main energy currency",
		NoteType:   "Basic",
	})
	if suggestionRR.Code != http.StatusOK {
		t.Fatalf("expected AI suggestions 200, got %d (%s)", suggestionRR.Code, suggestionRR.Body.String())
	}
	suggestionResp := decodeJSON[AICardSuggestionsResponse](t, suggestionRR)
	if suggestionResp.Provider != "dev" {
		t.Fatalf("expected dev provider in test config, got %q", suggestionResp.Provider)
	}
	if len(suggestionResp.Suggestions) != 2 {
		t.Fatalf("expected 2 AI suggestions, got %d", len(suggestionResp.Suggestions))
	}
	if got := suggestionResp.Suggestions[0].FieldVals["Front"]; got != "Mitochondria" {
		t.Fatalf("expected first Front field to match source prompt, got %q", got)
	}
	if got := suggestionResp.Suggestions[0].FieldVals["Back"]; got != "The powerhouse of the cell" {
		t.Fatalf("expected first Back field to match source answer, got %q", got)
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

func TestAPI_StudySessionLifecycle(t *testing.T) {
	env := setupAPITestEnv(t)

	createRR := doJSONRequest(t, env.router, http.MethodPost, "/api/study-sessions", CreateStudySessionRequest{
		DeckID: 1,
		Mode:   "review",
	})
	if createRR.Code != http.StatusCreated {
		t.Fatalf("expected study session create 201, got %d (%s)", createRR.Code, createRR.Body.String())
	}

	session := decodeJSON[StudySession](t, createRR)
	if session.Status != "active" {
		t.Fatalf("expected active study session, got %s", session.Status)
	}
	if session.DeckID != 1 {
		t.Fatalf("expected deckId=1, got %d", session.DeckID)
	}

	cardsReviewed := 2
	againCount := 1
	goodCount := 1
	progressRR := doJSONRequest(t, env.router, http.MethodPatch, fmt.Sprintf("/api/study-sessions/%s", session.ID), UpdateStudySessionRequest{
		CardsReviewed: &cardsReviewed,
		AgainCount:    &againCount,
		GoodCount:     &goodCount,
	})
	if progressRR.Code != http.StatusOK {
		t.Fatalf("expected study session update 200, got %d (%s)", progressRR.Code, progressRR.Body.String())
	}

	progressed := decodeJSON[StudySession](t, progressRR)
	if progressed.CardsReviewed != 2 || progressed.AgainCount != 1 || progressed.GoodCount != 1 {
		t.Fatalf("expected persisted progress counts, got %+v", progressed)
	}

	completedRR := doJSONRequest(t, env.router, http.MethodPatch, fmt.Sprintf("/api/study-sessions/%s", session.ID), UpdateStudySessionRequest{
		Status:        "completed",
		CardsReviewed: &cardsReviewed,
		AgainCount:    &againCount,
		GoodCount:     &goodCount,
	})
	if completedRR.Code != http.StatusOK {
		t.Fatalf("expected study session completion 200, got %d (%s)", completedRR.Code, completedRR.Body.String())
	}

	completed := decodeJSON[StudySession](t, completedRR)
	if completed.Status != "completed" {
		t.Fatalf("expected completed study session, got %s", completed.Status)
	}
	if completed.EndedAt.IsZero() {
		t.Fatalf("expected completed study session to have endedAt set")
	}

	reloaded, err := env.store.GetStudySession(session.ID)
	if err != nil {
		t.Fatalf("failed to reload study session: %v", err)
	}
	if reloaded.Status != "completed" || reloaded.CardsReviewed != 2 || reloaded.GoodCount != 1 {
		t.Fatalf("unexpected reloaded study session: %+v", reloaded)
	}

	dashboardRR := doRawRequest(env.router, http.MethodGet, "/api/dashboard", "")
	if dashboardRR.Code != http.StatusOK {
		t.Fatalf("expected dashboard 200, got %d (%s)", dashboardRR.Code, dashboardRR.Body.String())
	}
	dashboard := decodeJSON[DashboardResponse](t, dashboardRR)
	if dashboard.StudyAnalytics.Sessions7D != 1 || dashboard.StudyAnalytics.CardsReviewed7D != 2 {
		t.Fatalf("expected dashboard study analytics to reflect session progress, got %+v", dashboard.StudyAnalytics)
	}
	if dashboard.StudyAnalytics.CurrentStreak != 1 {
		t.Fatalf("expected dashboard current streak to be 1 after a completed session, got %d", dashboard.StudyAnalytics.CurrentStreak)
	}
	if len(dashboard.StudyAnalytics.RecentSessions) != 1 {
		t.Fatalf("expected dashboard to include one recent session, got %+v", dashboard.StudyAnalytics.RecentSessions)
	}
	if dashboard.StudyAnalytics.RecentSessions[0].DeckName != "Default" {
		t.Fatalf("expected dashboard recent session to resolve deck name, got %+v", dashboard.StudyAnalytics.RecentSessions[0])
	}

	decksRR := doRawRequest(env.router, http.MethodGet, "/api/decks", "")
	if decksRR.Code != http.StatusOK {
		t.Fatalf("expected decks 200, got %d (%s)", decksRR.Code, decksRR.Body.String())
	}
	decks := decodeJSON[[]DeckResponse](t, decksRR)
	if len(decks) == 0 {
		t.Fatalf("expected at least one deck in response")
	}
	if decks[0].Analytics.Sessions7D != 1 || decks[0].Analytics.CardsReviewed7D != 2 {
		t.Fatalf("expected deck analytics to reflect completed session, got %+v", decks[0].Analytics)
	}
	if decks[0].Analytics.GoodCount7D != 1 || decks[0].Analytics.AgainCount7D != 1 {
		t.Fatalf("expected deck answer analytics to reflect session ratings, got %+v", decks[0].Analytics)
	}

	analyticsRR := doRawRequest(env.router, http.MethodGet, "/api/analytics/overview", "")
	if analyticsRR.Code != http.StatusOK {
		t.Fatalf("expected analytics overview 200, got %d (%s)", analyticsRR.Code, analyticsRR.Body.String())
	}
	analytics := decodeJSON[StudyAnalyticsOverview](t, analyticsRR)
	if analytics.AnswerBreakdown.Again != 1 || analytics.AnswerBreakdown.Good != 1 {
		t.Fatalf("expected analytics answer breakdown to reflect session answers, got %+v", analytics.AnswerBreakdown)
	}
	if len(analytics.DailyActivity) != 7 {
		t.Fatalf("expected analytics daily activity to include 7 days, got %+v", analytics.DailyActivity)
	}
	if len(analytics.RecentSessions) != 1 || analytics.RecentSessions[0].CardsReviewed != 2 {
		t.Fatalf("expected analytics recent session payload to reflect session progress, got %+v", analytics.RecentSessions)
	}

	rejectedRR := doJSONRequest(t, env.router, http.MethodPatch, fmt.Sprintf("/api/study-sessions/%s", session.ID), UpdateStudySessionRequest{
		CardsReviewed: &cardsReviewed,
	})
	if rejectedRR.Code != http.StatusConflict {
		t.Fatalf("expected closed session update to return 409, got %d (%s)", rejectedRR.Code, rejectedRR.Body.String())
	}
}

func TestAPI_FocusSessionLifecycleAndAnalytics(t *testing.T) {
	env := setupAPITestEnv(t)

	createRR := doJSONRequest(t, env.router, http.MethodPost, "/api/study-sessions", CreateStudySessionRequest{
		Mode:          "focus",
		Protocol:      "pomodoro",
		TargetMinutes: 25,
		BreakMinutes:  5,
	})
	if createRR.Code != http.StatusCreated {
		t.Fatalf("expected focus session create 201, got %d (%s)", createRR.Code, createRR.Body.String())
	}

	session := decodeJSON[StudySession](t, createRR)
	if session.Mode != "focus" || session.Protocol != "pomodoro" {
		t.Fatalf("expected focus pomodoro session, got %+v", session)
	}
	if session.TargetMinutes != 25 || session.BreakMinutes != 5 {
		t.Fatalf("expected focus timing to persist, got %+v", session)
	}

	endedAt := time.Now().UTC().Add(25 * time.Minute)
	completedRR := doJSONRequest(t, env.router, http.MethodPatch, fmt.Sprintf("/api/study-sessions/%s", session.ID), UpdateStudySessionRequest{
		Status:  "completed",
		EndedAt: endedAt,
	})
	if completedRR.Code != http.StatusOK {
		t.Fatalf("expected focus session completion 200, got %d (%s)", completedRR.Code, completedRR.Body.String())
	}

	completed := decodeJSON[StudySession](t, completedRR)
	if completed.Status != "completed" {
		t.Fatalf("expected completed focus session, got %+v", completed)
	}

	analyticsRR := doRawRequest(env.router, http.MethodGet, "/api/analytics/overview", "")
	if analyticsRR.Code != http.StatusOK {
		t.Fatalf("expected analytics overview 200, got %d (%s)", analyticsRR.Code, analyticsRR.Body.String())
	}
	analytics := decodeJSON[StudyAnalyticsOverview](t, analyticsRR)
	if analytics.FocusSessions7D != 1 {
		t.Fatalf("expected one completed focus session, got %+v", analytics)
	}
	if analytics.FocusMinutes7D < 20 {
		t.Fatalf("expected focus minutes to reflect the completed block, got %+v", analytics)
	}
	if analytics.CurrentStreak != 1 {
		t.Fatalf("expected focus session to count toward streak, got %+v", analytics)
	}
	if len(analytics.RecentSessions) != 1 {
		t.Fatalf("expected focus session to appear in recent sessions, got %+v", analytics.RecentSessions)
	}
	if analytics.RecentSessions[0].Mode != "focus" || analytics.RecentSessions[0].Protocol != "pomodoro" {
		t.Fatalf("expected recent session to include focus metadata, got %+v", analytics.RecentSessions[0])
	}
	if analytics.RecentSessions[0].TargetMinutes != 25 || analytics.RecentSessions[0].BreakMinutes != 5 {
		t.Fatalf("expected recent focus session timing metadata, got %+v", analytics.RecentSessions[0])
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

func TestAPI_SharedCardsKeepPerUserDueQueuesSeparate(t *testing.T) {
	env := setupAPITestEnv(t)
	secondClient := createAuthenticatedTestClient(t, env, "second@example.com", "Second User")

	created := createNoteForTest(t, env, CreateNoteRequest{
		TypeID: "Basic",
		DeckID: 1,
		FieldVals: map[string]string{
			"Front": "Per-user review state question",
			"Back":  "Per-user review state answer",
		},
	}, nil)

	cardID := created.Cards[0].ID

	firstDue := doRawRequest(env.router, http.MethodGet, "/api/decks/1/due?limit=10", "")
	if firstDue.Code != http.StatusOK {
		t.Fatalf("expected first user due cards 200, got %d (%s)", firstDue.Code, firstDue.Body.String())
	}
	firstDueCards := decodeJSON[[]Card](t, firstDue)
	if len(firstDueCards) != 1 {
		t.Fatalf("expected first user to have 1 due card, got %d", len(firstDueCards))
	}

	secondDue := doRawRequest(secondClient.router, http.MethodGet, "/api/decks/1/due?limit=10", "")
	if secondDue.Code != http.StatusOK {
		t.Fatalf("expected second user due cards 200, got %d (%s)", secondDue.Code, secondDue.Body.String())
	}
	secondDueCards := decodeJSON[[]Card](t, secondDue)
	if len(secondDueCards) != 1 {
		t.Fatalf("expected second user to have 1 due card, got %d", len(secondDueCards))
	}

	answer := doJSONRequest(t, env.router, http.MethodPost, fmt.Sprintf("/api/cards/%d/answer", cardID), AnswerCardRequest{
		Rating:      3,
		TimeTakenMs: 900,
	})
	if answer.Code != http.StatusOK {
		t.Fatalf("expected first user answer 200, got %d (%s)", answer.Code, answer.Body.String())
	}

	firstDueAfter := doRawRequest(env.router, http.MethodGet, "/api/decks/1/due?limit=10", "")
	if firstDueAfter.Code != http.StatusOK {
		t.Fatalf("expected first user due cards after answer 200, got %d (%s)", firstDueAfter.Code, firstDueAfter.Body.String())
	}
	firstDueCardsAfter := decodeJSON[[]Card](t, firstDueAfter)
	if len(firstDueCardsAfter) != 0 {
		t.Fatalf("expected first user to have 0 due cards after answering, got %d", len(firstDueCardsAfter))
	}

	secondDueAfter := doRawRequest(secondClient.router, http.MethodGet, "/api/decks/1/due?limit=10", "")
	if secondDueAfter.Code != http.StatusOK {
		t.Fatalf("expected second user due cards after first answer 200, got %d (%s)", secondDueAfter.Code, secondDueAfter.Body.String())
	}
	secondDueCardsAfter := decodeJSON[[]Card](t, secondDueAfter)
	if len(secondDueCardsAfter) != 1 {
		t.Fatalf("expected second user due queue to remain unchanged, got %d", len(secondDueCardsAfter))
	}

	firstStats := doRawRequest(env.router, http.MethodGet, "/api/decks/1/stats", "")
	if firstStats.Code != http.StatusOK {
		t.Fatalf("expected first user stats 200, got %d (%s)", firstStats.Code, firstStats.Body.String())
	}
	firstDeckStats := decodeJSON[DeckStats](t, firstStats)
	if firstDeckStats.DueToday != 0 {
		t.Fatalf("expected first user dueToday=0 after answer, got %d", firstDeckStats.DueToday)
	}

	secondStats := doRawRequest(secondClient.router, http.MethodGet, "/api/decks/1/stats", "")
	if secondStats.Code != http.StatusOK {
		t.Fatalf("expected second user stats 200, got %d (%s)", secondStats.Code, secondStats.Body.String())
	}
	secondDeckStats := decodeJSON[DeckStats](t, secondStats)
	if secondDeckStats.DueToday != 1 {
		t.Fatalf("expected second user dueToday=1 to remain pending, got %d", secondDeckStats.DueToday)
	}

	var revlogUserID string
	if err := env.store.db.QueryRow(`SELECT user_id FROM revlog WHERE card_id = ? ORDER BY reviewed_at DESC LIMIT 1`, cardID).Scan(&revlogUserID); err != nil {
		t.Fatalf("failed to query revlog user_id: %v", err)
	}
	sessionID := strings.TrimPrefix(env.authCookie, sessionCookieName+"=")
	sessionRecord, err := env.store.GetSessionRecord(sessionID)
	if err != nil {
		t.Fatalf("failed to load primary test session: %v", err)
	}
	if revlogUserID != sessionRecord.UserID {
		t.Fatalf("expected revlog user_id=%q, got %q", sessionRecord.UserID, revlogUserID)
	}

	var reviewStateCount int
	if err := env.store.db.QueryRow(`SELECT COUNT(*) FROM card_review_states WHERE card_id = ?`, cardID).Scan(&reviewStateCount); err != nil {
		t.Fatalf("failed to count card review states: %v", err)
	}
	if reviewStateCount < 2 {
		t.Fatalf("expected per-user review states for both users, got %d", reviewStateCount)
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

func TestAPI_StudyGroupCreationRequiresTeamOrEnterprise(t *testing.T) {
	env := setupAPITestEnv(t)

	createFree := doJSONRequest(t, env.router, http.MethodPost, "/api/study-groups", CreateStudyGroupRequest{
		Name:          "Free plan group",
		PrimaryDeckID: 1,
	})
	if createFree.Code != http.StatusForbidden {
		t.Fatalf("expected free plan study group create to return 403, got %d (%s)", createFree.Code, createFree.Body.String())
	}
	freeErr := decodeJSON[APIErrorResponse](t, createFree)
	if freeErr.Code != "study_groups_not_available" {
		t.Fatalf("expected study_groups_not_available code, got %+v", freeErr)
	}

	createPro := doJSONRequestWithHeaders(t, env.router, http.MethodPost, "/api/study-groups", CreateStudyGroupRequest{
		Name:          "Pro plan group",
		PrimaryDeckID: 1,
	}, map[string]string{"X-Vutadex-Plan": "pro"})
	if createPro.Code != http.StatusForbidden {
		t.Fatalf("expected pro plan study group create to return 403, got %d (%s)", createPro.Code, createPro.Body.String())
	}

	createTeam := doJSONRequestWithHeaders(t, env.router, http.MethodPost, "/api/study-groups", CreateStudyGroupRequest{
		Name:          "Team plan group",
		PrimaryDeckID: 1,
	}, map[string]string{"X-Vutadex-Plan": "team"})
	if createTeam.Code != http.StatusCreated {
		t.Fatalf("expected team plan study group create to return 201, got %d (%s)", createTeam.Code, createTeam.Body.String())
	}
}

func TestAPI_StudyGroupsUsePublishedVersionsAndPersonalInstalls(t *testing.T) {
	env := setupAPITestEnv(t)
	memberClient := createAuthenticatedTestClient(t, env, "group-member@example.com", "Group Member")
	teamHeaders := map[string]string{"X-Vutadex-Plan": "team"}

	createNoteForTest(t, env, CreateNoteRequest{
		TypeID: "Basic",
		DeckID: 1,
		FieldVals: map[string]string{
			"Front": "Source deck card 1",
			"Back":  "Answer 1",
		},
	}, nil)

	createGroup := doJSONRequestWithHeaders(t, env.router, http.MethodPost, "/api/study-groups", CreateStudyGroupRequest{
		Name:          "Anatomy Cohort",
		Description:   "Canonical source deck + personal installs",
		PrimaryDeckID: 1,
		Visibility:    "private",
		JoinPolicy:    "invite",
	}, teamHeaders)
	if createGroup.Code != http.StatusCreated {
		t.Fatalf("expected create study group 201, got %d (%s)", createGroup.Code, createGroup.Body.String())
	}
	groupDetail := decodeJSON[StudyGroupDetail](t, createGroup)
	groupID := groupDetail.Group.ID
	if groupDetail.Role != "owner" || groupDetail.MembershipStatus != "active" {
		t.Fatalf("expected creator to be active owner, got role=%q status=%q", groupDetail.Role, groupDetail.MembershipStatus)
	}

	publishV1 := doJSONRequestWithHeaders(t, env.router, http.MethodPost, fmt.Sprintf("/api/study-groups/%s/versions", groupID), PublishStudyGroupVersionRequest{
		ChangeSummary: "Initial release",
	}, teamHeaders)
	if publishV1.Code != http.StatusCreated {
		t.Fatalf("expected publish version 1 to return 201, got %d (%s)", publishV1.Code, publishV1.Body.String())
	}
	version1 := decodeJSON[StudyGroupVersion](t, publishV1)
	if version1.VersionNumber != 1 || version1.NoteCount != 1 || version1.CardCount != 1 {
		t.Fatalf("expected initial version metadata to reflect source deck, got %+v", version1)
	}

	inviteMember := doJSONRequestWithHeaders(t, env.router, http.MethodPost, fmt.Sprintf("/api/study-groups/%s/members", groupID), InviteStudyGroupMemberRequest{
		Email: memberClient.user.Email,
		Role:  "read",
	}, teamHeaders)
	if inviteMember.Code != http.StatusCreated {
		t.Fatalf("expected member invite to return 201, got %d (%s)", inviteMember.Code, inviteMember.Body.String())
	}
	invite := decodeJSON[StudyGroupMember](t, inviteMember)
	if invite.InviteToken == "" {
		t.Fatalf("expected invite token to be returned, got %+v", invite)
	}

	joinGroup := doJSONRequest(t, memberClient.router, http.MethodPost, "/api/study-groups/join", JoinStudyGroupRequest{
		Token:                  invite.InviteToken,
		DestinationWorkspaceID: memberClient.workspace.ID,
		InstallLatest:          false,
	})
	if joinGroup.Code != http.StatusOK {
		t.Fatalf("expected join group 200, got %d (%s)", joinGroup.Code, joinGroup.Body.String())
	}
	memberDetailAfterJoin := decodeJSON[StudyGroupDetail](t, joinGroup)
	if memberDetailAfterJoin.Role != "read" || memberDetailAfterJoin.MembershipStatus != "active" {
		t.Fatalf("expected joined member to be active read member, got role=%q status=%q", memberDetailAfterJoin.Role, memberDetailAfterJoin.MembershipStatus)
	}
	if memberDetailAfterJoin.CurrentUserInstall != nil {
		t.Fatalf("expected join without installLatest to leave current install nil, got %+v", memberDetailAfterJoin.CurrentUserInstall)
	}

	ownerInstallRR := doJSONRequest(t, env.router, http.MethodPost, fmt.Sprintf("/api/study-groups/%s/installs", groupID), InstallStudyGroupDeckRequest{})
	if ownerInstallRR.Code != http.StatusCreated {
		t.Fatalf("expected owner install 201, got %d (%s)", ownerInstallRR.Code, ownerInstallRR.Body.String())
	}
	ownerInstall := decodeJSON[StudyGroupInstall](t, ownerInstallRR)

	memberInstallRR := doJSONRequest(t, memberClient.router, http.MethodPost, fmt.Sprintf("/api/study-groups/%s/installs", groupID), InstallStudyGroupDeckRequest{
		DestinationWorkspaceID: memberClient.workspace.ID,
	})
	if memberInstallRR.Code != http.StatusCreated {
		t.Fatalf("expected member install 201, got %d (%s)", memberInstallRR.Code, memberInstallRR.Body.String())
	}
	memberInstall := decodeJSON[StudyGroupInstall](t, memberInstallRR)
	if memberInstall.SourceVersionNumber != 1 || memberInstall.Status != "active" || memberInstall.SyncState != "clean" {
		t.Fatalf("expected clean active install on version 1, got %+v", memberInstall)
	}
	if memberInstall.InstalledDeckID == ownerInstall.InstalledDeckID {
		t.Fatalf("expected personal installs to create separate deck copies, got same deck id %d", memberInstall.InstalledDeckID)
	}

	var ownerCardID int64
	if err := env.store.db.QueryRow(`SELECT id FROM cards WHERE deck_id = ? ORDER BY id ASC LIMIT 1`, ownerInstall.InstalledDeckID).Scan(&ownerCardID); err != nil {
		t.Fatalf("failed to load owner installed deck card: %v", err)
	}
	var memberCardID int64
	if err := env.store.db.QueryRow(`SELECT id FROM cards WHERE deck_id = ? ORDER BY id ASC LIMIT 1`, memberInstall.InstalledDeckID).Scan(&memberCardID); err != nil {
		t.Fatalf("failed to load member installed deck card: %v", err)
	}

	ownerDueBefore := doRawRequest(env.router, http.MethodGet, fmt.Sprintf("/api/decks/%d/due?limit=10", ownerInstall.InstalledDeckID), "")
	if ownerDueBefore.Code != http.StatusOK {
		t.Fatalf("expected owner due queue 200, got %d (%s)", ownerDueBefore.Code, ownerDueBefore.Body.String())
	}
	if cards := decodeJSON[[]Card](t, ownerDueBefore); len(cards) != 1 {
		t.Fatalf("expected owner installed deck due count to start at 1, got %d", len(cards))
	}

	memberDueBefore := doRawRequest(memberClient.router, http.MethodGet, fmt.Sprintf("/api/decks/%d/due?limit=10", memberInstall.InstalledDeckID), "")
	if memberDueBefore.Code != http.StatusOK {
		t.Fatalf("expected member due queue 200, got %d (%s)", memberDueBefore.Code, memberDueBefore.Body.String())
	}
	if cards := decodeJSON[[]Card](t, memberDueBefore); len(cards) != 1 {
		t.Fatalf("expected member installed deck due count to start at 1, got %d", len(cards))
	}

	memberAnswer := doJSONRequest(t, memberClient.router, http.MethodPost, fmt.Sprintf("/api/cards/%d/answer", memberCardID), AnswerCardRequest{
		Rating:      3,
		TimeTakenMs: 850,
	})
	if memberAnswer.Code != http.StatusOK {
		t.Fatalf("expected member answer 200, got %d (%s)", memberAnswer.Code, memberAnswer.Body.String())
	}

	ownerDueAfterMemberAnswer := doRawRequest(env.router, http.MethodGet, fmt.Sprintf("/api/decks/%d/due?limit=10", ownerInstall.InstalledDeckID), "")
	if ownerDueAfterMemberAnswer.Code != http.StatusOK {
		t.Fatalf("expected owner due queue after member answer 200, got %d (%s)", ownerDueAfterMemberAnswer.Code, ownerDueAfterMemberAnswer.Body.String())
	}
	if cards := decodeJSON[[]Card](t, ownerDueAfterMemberAnswer); len(cards) != 1 {
		t.Fatalf("expected owner due queue to stay unchanged after member studies personal copy, got %d", len(cards))
	}

	ownerAnswer := doJSONRequest(t, env.router, http.MethodPost, fmt.Sprintf("/api/cards/%d/answer", ownerCardID), AnswerCardRequest{
		Rating:      3,
		TimeTakenMs: 900,
	})
	if ownerAnswer.Code != http.StatusOK {
		t.Fatalf("expected owner answer 200, got %d (%s)", ownerAnswer.Code, ownerAnswer.Body.String())
	}

	createNoteForTest(t, env, CreateNoteRequest{
		TypeID: "Basic",
		DeckID: 1,
		FieldVals: map[string]string{
			"Front": "Source deck card 2",
			"Back":  "Answer 2",
		},
	}, nil)

	publishV2 := doJSONRequestWithHeaders(t, env.router, http.MethodPost, fmt.Sprintf("/api/study-groups/%s/versions", groupID), PublishStudyGroupVersionRequest{
		ChangeSummary: "Added a second source card",
	}, teamHeaders)
	if publishV2.Code != http.StatusCreated {
		t.Fatalf("expected publish version 2 to return 201, got %d (%s)", publishV2.Code, publishV2.Body.String())
	}
	version2 := decodeJSON[StudyGroupVersion](t, publishV2)
	if version2.VersionNumber != 2 || version2.NoteCount != 2 || version2.CardCount != 2 {
		t.Fatalf("expected version 2 metadata to reflect updated source deck, got %+v", version2)
	}

	memberDetailBeforeUpdateRR := doRawRequest(memberClient.router, http.MethodGet, fmt.Sprintf("/api/study-groups/%s", groupID), "")
	if memberDetailBeforeUpdateRR.Code != http.StatusOK {
		t.Fatalf("expected member detail before update 200, got %d (%s)", memberDetailBeforeUpdateRR.Code, memberDetailBeforeUpdateRR.Body.String())
	}
	memberDetailBeforeUpdate := decodeJSON[StudyGroupDetail](t, memberDetailBeforeUpdateRR)
	if !memberDetailBeforeUpdate.UpdateAvailable {
		t.Fatalf("expected updateAvailable=true after publishing a newer source version")
	}
	if memberDetailBeforeUpdate.CurrentUserInstall == nil || memberDetailBeforeUpdate.CurrentUserInstall.SourceVersionNumber != 1 {
		t.Fatalf("expected member current install to still point at version 1 before update, got %+v", memberDetailBeforeUpdate.CurrentUserInstall)
	}

	updateInstallRR := doJSONRequest(t, memberClient.router, http.MethodPost, fmt.Sprintf("/api/study-groups/%s/installs/%s/update", groupID, memberInstall.ID), UpdateStudyGroupInstallRequest{})
	if updateInstallRR.Code != http.StatusOK {
		t.Fatalf("expected install update 200, got %d (%s)", updateInstallRR.Code, updateInstallRR.Body.String())
	}
	memberInstallV2 := decodeJSON[StudyGroupInstall](t, updateInstallRR)
	if memberInstallV2.ID == memberInstall.ID {
		t.Fatalf("expected install update to create a fresh install record")
	}
	if memberInstallV2.SourceVersionNumber != 2 || memberInstallV2.Status != "active" || memberInstallV2.SyncState != "clean" {
		t.Fatalf("expected active clean version 2 install after update, got %+v", memberInstallV2)
	}

	oldInstall, err := env.store.GetStudyGroupInstall(memberInstall.ID)
	if err != nil {
		t.Fatalf("failed to reload superseded install: %v", err)
	}
	if oldInstall.Status != "superseded" || oldInstall.SupersededByInstallID != memberInstallV2.ID {
		t.Fatalf("expected original install to be superseded by the new install, got %+v", oldInstall)
	}

	oldNotes, oldCards, err := env.store.GetDeckContentSummary(memberInstall.InstalledDeckID)
	if err != nil {
		t.Fatalf("failed to read original install summary: %v", err)
	}
	if oldNotes != 1 || oldCards != 1 {
		t.Fatalf("expected original install copy to remain intact on version 1, got notes=%d cards=%d", oldNotes, oldCards)
	}
	newNotes, newCards, err := env.store.GetDeckContentSummary(memberInstallV2.InstalledDeckID)
	if err != nil {
		t.Fatalf("failed to read updated install summary: %v", err)
	}
	if newNotes != 2 || newCards != 2 {
		t.Fatalf("expected updated install copy to reflect version 2 content, got notes=%d cards=%d", newNotes, newCards)
	}

	forkedName := "Forked Personal Copy"
	renameInstallDeck := doJSONRequest(t, memberClient.router, http.MethodPatch, fmt.Sprintf("/api/decks/%d", memberInstallV2.InstalledDeckID), UpdateDeckRequest{
		Name: &forkedName,
	})
	if renameInstallDeck.Code != http.StatusOK {
		t.Fatalf("expected renaming installed deck to succeed, got %d (%s)", renameInstallDeck.Code, renameInstallDeck.Body.String())
	}

	memberDetailAfterForkRR := doRawRequest(memberClient.router, http.MethodGet, fmt.Sprintf("/api/study-groups/%s", groupID), "")
	if memberDetailAfterForkRR.Code != http.StatusOK {
		t.Fatalf("expected member detail after fork 200, got %d (%s)", memberDetailAfterForkRR.Code, memberDetailAfterForkRR.Body.String())
	}
	memberDetailAfterFork := decodeJSON[StudyGroupDetail](t, memberDetailAfterForkRR)
	if memberDetailAfterFork.CurrentUserInstall == nil || memberDetailAfterFork.CurrentUserInstall.SyncState != "forked" {
		t.Fatalf("expected renamed installed copy to be marked forked, got %+v", memberDetailAfterFork.CurrentUserInstall)
	}

	dashboardRR := doRawRequest(env.router, http.MethodGet, fmt.Sprintf("/api/study-groups/%s/dashboard", groupID), "")
	if dashboardRR.Code != http.StatusOK {
		t.Fatalf("expected study group dashboard 200, got %d (%s)", dashboardRR.Code, dashboardRR.Body.String())
	}
	dashboard := decodeJSON[StudyGroupDashboard](t, dashboardRR)
	if dashboard.MemberCount != 2 {
		t.Fatalf("expected dashboard memberCount=2, got %d", dashboard.MemberCount)
	}
	if dashboard.ActiveMembers7D != 2 {
		t.Fatalf("expected dashboard activeMembers7d=2 after both members studied, got %d", dashboard.ActiveMembers7D)
	}
	if dashboard.Reviews7D != 2 {
		t.Fatalf("expected dashboard reviews7d=2 after two answers, got %d", dashboard.Reviews7D)
	}
	if dashboard.LatestVersionNumber != 2 {
		t.Fatalf("expected dashboard latestVersionNumber=2, got %d", dashboard.LatestVersionNumber)
	}
	if dashboard.LatestVersionAdoption != 1 {
		t.Fatalf("expected dashboard latestVersionAdoption=1 with one active v2 install, got %d", dashboard.LatestVersionAdoption)
	}

	removeInstallRR := doJSONRequest(t, memberClient.router, http.MethodDelete, fmt.Sprintf("/api/study-groups/%s/installs/%s", groupID, memberInstallV2.ID), struct{}{})
	if removeInstallRR.Code != http.StatusNoContent {
		t.Fatalf("expected remove install 204, got %d (%s)", removeInstallRR.Code, removeInstallRR.Body.String())
	}
	removedInstall, err := env.store.GetStudyGroupInstall(memberInstallV2.ID)
	if err != nil {
		t.Fatalf("failed to reload removed install: %v", err)
	}
	if removedInstall.Status != "removed" {
		t.Fatalf("expected removed install status=removed, got %+v", removedInstall)
	}
	if _, err := env.store.GetDeck(memberInstallV2.InstalledDeckID); err == nil {
		t.Fatalf("expected removing an install to delete its copied deck %d", memberInstallV2.InstalledDeckID)
	}
}

func TestAPI_StudyGroupsInstallAcrossDifferentCollections(t *testing.T) {
	env := setupAPITestEnv(t)
	memberClient := createAuthenticatedIsolatedTestClient(t, env, "isolated-member@example.com", "Isolated Member")
	teamHeaders := map[string]string{"X-Vutadex-Plan": "team"}

	createNoteForTest(t, env, CreateNoteRequest{
		TypeID: "Basic",
		DeckID: 1,
		FieldVals: map[string]string{
			"Front": "Cross collection source card",
			"Back":  "Cross collection answer",
		},
	}, nil)

	createGroup := doJSONRequestWithHeaders(t, env.router, http.MethodPost, "/api/study-groups", CreateStudyGroupRequest{
		Name:          "Cross Collection Cohort",
		Description:   "Ensure installs work across real workspace collections",
		PrimaryDeckID: 1,
		Visibility:    "private",
		JoinPolicy:    "invite",
	}, teamHeaders)
	if createGroup.Code != http.StatusCreated {
		t.Fatalf("expected create study group 201, got %d (%s)", createGroup.Code, createGroup.Body.String())
	}
	groupDetail := decodeJSON[StudyGroupDetail](t, createGroup)
	groupID := groupDetail.Group.ID

	publishV1 := doJSONRequestWithHeaders(t, env.router, http.MethodPost, fmt.Sprintf("/api/study-groups/%s/versions", groupID), PublishStudyGroupVersionRequest{
		ChangeSummary: "Initial cross-collection release",
	}, teamHeaders)
	if publishV1.Code != http.StatusCreated {
		t.Fatalf("expected publish version 1 to return 201, got %d (%s)", publishV1.Code, publishV1.Body.String())
	}

	inviteMember := doJSONRequestWithHeaders(t, env.router, http.MethodPost, fmt.Sprintf("/api/study-groups/%s/members", groupID), InviteStudyGroupMemberRequest{
		Email: memberClient.user.Email,
		Role:  "read",
	}, teamHeaders)
	if inviteMember.Code != http.StatusCreated {
		t.Fatalf("expected member invite to return 201, got %d (%s)", inviteMember.Code, inviteMember.Body.String())
	}
	invite := decodeJSON[StudyGroupMember](t, inviteMember)

	joinGroup := doJSONRequest(t, memberClient.router, http.MethodPost, "/api/study-groups/join", JoinStudyGroupRequest{
		Token:                  invite.InviteToken,
		DestinationWorkspaceID: memberClient.workspace.ID,
		InstallLatest:          false,
	})
	if joinGroup.Code != http.StatusOK {
		t.Fatalf("expected join group 200, got %d (%s)", joinGroup.Code, joinGroup.Body.String())
	}

	memberInstallRR := doJSONRequest(t, memberClient.router, http.MethodPost, fmt.Sprintf("/api/study-groups/%s/installs", groupID), InstallStudyGroupDeckRequest{
		DestinationWorkspaceID: memberClient.workspace.ID,
	})
	if memberInstallRR.Code != http.StatusCreated {
		t.Fatalf("expected isolated member install 201, got %d (%s)", memberInstallRR.Code, memberInstallRR.Body.String())
	}
	memberInstall := decodeJSON[StudyGroupInstall](t, memberInstallRR)

	var destinationCollectionID string
	if err := env.store.db.QueryRow(`SELECT collection_id FROM decks WHERE id = ?`, memberInstall.InstalledDeckID).Scan(&destinationCollectionID); err != nil {
		t.Fatalf("failed to load installed deck collection: %v", err)
	}
	if destinationCollectionID != memberClient.workspace.CollectionID {
		t.Fatalf("expected installed deck to land in destination workspace collection %q, got %q", memberClient.workspace.CollectionID, destinationCollectionID)
	}

	var sourceCollectionID string
	if err := env.store.db.QueryRow(`SELECT collection_id FROM decks WHERE id = 1`).Scan(&sourceCollectionID); err != nil {
		t.Fatalf("failed to load source deck collection: %v", err)
	}
	if sourceCollectionID == destinationCollectionID {
		t.Fatalf("expected source and destination collections to differ, both were %q", sourceCollectionID)
	}

	memberDueBefore := doRawRequest(memberClient.router, http.MethodGet, fmt.Sprintf("/api/decks/%d/due?limit=10", memberInstall.InstalledDeckID), "")
	if memberDueBefore.Code != http.StatusOK {
		t.Fatalf("expected isolated member due queue 200, got %d (%s)", memberDueBefore.Code, memberDueBefore.Body.String())
	}
	memberCards := decodeJSON[[]Card](t, memberDueBefore)
	if len(memberCards) != 1 {
		t.Fatalf("expected isolated member due count to start at 1, got %d", len(memberCards))
	}

	memberAnswer := doJSONRequest(t, memberClient.router, http.MethodPost, fmt.Sprintf("/api/cards/%d/answer", memberCards[0].ID), AnswerCardRequest{
		Rating:      3,
		TimeTakenMs: 600,
	})
	if memberAnswer.Code != http.StatusOK {
		t.Fatalf("expected isolated member answer 200, got %d (%s)", memberAnswer.Code, memberAnswer.Body.String())
	}

	ownerDue := doRawRequest(env.router, http.MethodGet, "/api/decks/1/due?limit=10", "")
	if ownerDue.Code != http.StatusOK {
		t.Fatalf("expected owner source deck due queue 200, got %d (%s)", ownerDue.Code, ownerDue.Body.String())
	}
	if cards := decodeJSON[[]Card](t, ownerDue); len(cards) != 1 {
		t.Fatalf("expected source deck due queue to remain unchanged after isolated member studies installed copy, got %d", len(cards))
	}
}

func TestAPI_MarketplacePublishingRequiresProOrHigher(t *testing.T) {
	env := setupAPITestEnv(t)

	createFree := doJSONRequest(t, env.router, http.MethodPost, "/api/marketplace/listings", CreateMarketplaceListingRequest{
		DeckID:      1,
		Title:       "Free Plan Listing",
		Summary:     "Should be blocked",
		Description: "Free users cannot publish marketplace listings.",
		PriceMode:   "free",
	})
	if createFree.Code != http.StatusForbidden {
		t.Fatalf("expected free plan marketplace create to return 403, got %d (%s)", createFree.Code, createFree.Body.String())
	}
	freeErr := decodeJSON[APIErrorResponse](t, createFree)
	if freeErr.Code != "marketplace_publish_not_available" {
		t.Fatalf("expected marketplace_publish_not_available code, got %+v", freeErr)
	}

	createPro := doJSONRequestWithHeaders(t, env.router, http.MethodPost, "/api/marketplace/listings", CreateMarketplaceListingRequest{
		DeckID:      1,
		Title:       "Pro Plan Listing",
		Summary:     "Creator listing",
		Description: "Pro users can create marketplace drafts.",
		Category:    "Medicine",
		Tags:        []string{"anki-alternative", "exam-prep"},
		PriceMode:   "free",
	}, map[string]string{"X-Vutadex-Plan": "pro"})
	if createPro.Code != http.StatusCreated {
		t.Fatalf("expected pro plan marketplace create to return 201, got %d (%s)", createPro.Code, createPro.Body.String())
	}
	detail := decodeJSON[MarketplaceListingDetail](t, createPro)
	if detail.Listing.Status != "draft" || !detail.CanEdit {
		t.Fatalf("expected creator draft detail after pro create, got %+v", detail)
	}

	publish := doJSONRequestWithHeaders(t, env.router, http.MethodPost, fmt.Sprintf("/api/marketplace/listings/%s/publish", detail.Listing.ID), PublishMarketplaceListingRequest{
		ChangeSummary: "Initial catalog release",
	}, map[string]string{"X-Vutadex-Plan": "pro"})
	if publish.Code != http.StatusCreated {
		t.Fatalf("expected pro publish 201, got %d (%s)", publish.Code, publish.Body.String())
	}
	version := decodeJSON[MarketplaceListingVersion](t, publish)
	if version.VersionNumber != 1 {
		t.Fatalf("expected first marketplace version number to be 1, got %+v", version)
	}
}

func TestAPI_MarketplacePremiumListingsRequireCreatorSetupAndCheckout(t *testing.T) {
	env := setupAPITestEnv(t)
	memberClient := createAuthenticatedIsolatedTestClient(t, env, "marketplace-buyer@example.com", "Marketplace Buyer")
	proHeaders := map[string]string{"X-Vutadex-Plan": "pro"}

	createNoteForTest(t, env, CreateNoteRequest{
		TypeID: "Basic",
		DeckID: 1,
		FieldVals: map[string]string{
			"Front": "Premium source card",
			"Back":  "Premium answer",
		},
	}, nil)

	createListing := doJSONRequestWithHeaders(t, env.router, http.MethodPost, "/api/marketplace/listings", CreateMarketplaceListingRequest{
		DeckID:      1,
		Title:       "Premium Listing",
		Summary:     "Premium marketplace metadata",
		Description: "Premium checkout requires creator setup and buyer licensing.",
		Category:    "Certifications",
		PriceMode:   "premium",
		PriceCents:  4900,
		Currency:    "USD",
	}, proHeaders)
	if createListing.Code != http.StatusCreated {
		t.Fatalf("expected premium listing create 201, got %d (%s)", createListing.Code, createListing.Body.String())
	}
	detail := decodeJSON[MarketplaceListingDetail](t, createListing)

	creatorStatus := doRawRequestWithHeaders(env.router, http.MethodGet, "/api/marketplace/creator-account/status", "", proHeaders)
	if creatorStatus.Code != http.StatusOK {
		t.Fatalf("expected creator status 200, got %d (%s)", creatorStatus.Code, creatorStatus.Body.String())
	}
	initialStatus := decodeJSON[MarketplaceCreatorAccountStatusResponse](t, creatorStatus)
	if initialStatus.CanSellPremium || initialStatus.Account != nil {
		t.Fatalf("expected creator premium selling to be unavailable before onboarding, got %+v", initialStatus)
	}

	startCreator := doRawRequestWithHeaders(env.router, http.MethodPost, "/api/marketplace/creator-account/start", "", proHeaders)
	if startCreator.Code != http.StatusOK {
		t.Fatalf("expected creator account start 200, got %d (%s)", startCreator.Code, startCreator.Body.String())
	}
	creatorAccount := decodeJSON[MarketplaceCreatorAccountStatusResponse](t, startCreator)
	if !creatorAccount.CanSellPremium || creatorAccount.Account == nil || !creatorAccount.Account.ChargesEnabled || !creatorAccount.Account.PayoutsEnabled {
		t.Fatalf("expected creator account to be premium-ready in development, got %+v", creatorAccount)
	}

	publish := doJSONRequestWithHeaders(t, env.router, http.MethodPost, fmt.Sprintf("/api/marketplace/listings/%s/publish", detail.Listing.ID), PublishMarketplaceListingRequest{
		ChangeSummary: "Premium v1",
	}, proHeaders)
	if publish.Code != http.StatusCreated {
		t.Fatalf("expected premium listing publish 201, got %d (%s)", publish.Code, publish.Body.String())
	}

	install := doJSONRequest(t, memberClient.router, http.MethodPost, fmt.Sprintf("/api/marketplace/listings/%s/installs", detail.Listing.Slug), InstallMarketplaceListingRequest{
		DestinationWorkspaceID: memberClient.workspace.ID,
	})
	if install.Code != http.StatusConflict {
		t.Fatalf("expected premium listing install to return 409, got %d (%s)", install.Code, install.Body.String())
	}
	apiErr := decodeJSON[APIErrorResponse](t, install)
	if apiErr.Code != "marketplace_purchase_required" {
		t.Fatalf("expected marketplace_purchase_required code, got %+v", apiErr)
	}

	checkout := doRawRequest(memberClient.router, http.MethodPost, fmt.Sprintf("/api/marketplace/listings/%s/checkout", detail.Listing.Slug), "")
	if checkout.Code != http.StatusOK {
		t.Fatalf("expected premium listing checkout 200, got %d (%s)", checkout.Code, checkout.Body.String())
	}
	checkoutResp := decodeJSON[MarketplaceCheckoutResponse](t, checkout)
	if !checkoutResp.Completed || checkoutResp.License == nil || checkoutResp.Order.Status != "paid" {
		t.Fatalf("expected development checkout to complete immediately with a license, got %+v", checkoutResp)
	}

	installAfterPurchase := doJSONRequest(t, memberClient.router, http.MethodPost, fmt.Sprintf("/api/marketplace/listings/%s/installs", detail.Listing.Slug), InstallMarketplaceListingRequest{
		DestinationWorkspaceID: memberClient.workspace.ID,
	})
	if installAfterPurchase.Code != http.StatusCreated {
		t.Fatalf("expected premium listing install after purchase to return 201, got %d (%s)", installAfterPurchase.Code, installAfterPurchase.Body.String())
	}
	installed := decodeJSON[MarketplaceInstall](t, installAfterPurchase)
	if installed.SourceVersionNumber != 1 {
		t.Fatalf("expected premium install to target v1, got %+v", installed)
	}
}

func TestAPI_MarketplaceCheckoutSessionSyncCompletesStripeOrder(t *testing.T) {
	cfg := mustLocalAppConfig()
	cfg.Environment = "production"
	cfg.Stripe.SecretKey = "sk_test_marketplace"
	cfg.Stripe.WebhookSecret = "whsec_marketplace_test"
	env := setupAPITestEnvWithConfig(t, cfg)
	memberClient := createAuthenticatedIsolatedTestClient(t, env, "stripe-buyer@example.com", "Stripe Buyer")
	proHeaders := map[string]string{"X-Vutadex-Plan": "pro"}

	createNoteForTest(t, env, CreateNoteRequest{
		TypeID: "Basic",
		DeckID: 1,
		FieldVals: map[string]string{
			"Front": "Stripe source card",
			"Back":  "Stripe answer",
		},
	}, nil)

	ownerSessionReq := httptest.NewRequest(http.MethodGet, "/", nil)
	ownerSessionReq.Header.Set("Cookie", env.authCookie)
	ownerSession := env.handler.sessionFromRequest(ownerSessionReq)

	creatorAccount := &MarketplaceCreatorAccount{
		ID:                newID("mca"),
		UserID:            ownerSession.UserID,
		WorkspaceID:       ownerSession.WorkspaceID,
		Provider:          "stripe",
		ProviderAccountID: "acct_creator_ready",
		OnboardingStatus:  "active",
		DetailsSubmitted:  true,
		ChargesEnabled:    true,
		PayoutsEnabled:    true,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}
	if err := env.store.UpsertMarketplaceCreatorAccount(creatorAccount); err != nil {
		t.Fatalf("failed to seed creator account: %v", err)
	}

	originalTransport := http.DefaultTransport
	http.DefaultTransport = roundTripFunc(func(req *http.Request) (*http.Response, error) {
		if req.URL.Host != "api.stripe.com" {
			return nil, fmt.Errorf("unexpected host %s", req.URL.Host)
		}
		switch {
		case req.Method == http.MethodGet && req.URL.Path == "/v1/accounts/acct_creator_ready":
			return &http.Response{
				StatusCode: http.StatusOK,
				Header:     make(http.Header),
				Body:       io.NopCloser(strings.NewReader(`{"id":"acct_creator_ready","details_submitted":true,"charges_enabled":true,"payouts_enabled":true}`)),
			}, nil
		case req.Method == http.MethodPost && req.URL.Path == "/v1/accounts/acct_creator_ready/login_links":
			return &http.Response{
				StatusCode: http.StatusOK,
				Header:     make(http.Header),
				Body:       io.NopCloser(strings.NewReader(`{"url":"https://dashboard.stripe.test/acct_creator_ready"}`)),
			}, nil
		case req.Method == http.MethodGet && req.URL.Path == "/v1/checkout/sessions/cs_test_sync_paid":
			if got := req.Header.Get("Stripe-Account"); got != "acct_creator_ready" {
				return nil, fmt.Errorf("expected Stripe-Account acct_creator_ready, got %q", got)
			}
			return &http.Response{
				StatusCode: http.StatusOK,
				Header:     make(http.Header),
				Body:       io.NopCloser(strings.NewReader(`{"id":"cs_test_sync_paid","status":"complete","payment_status":"paid","payment_intent":"pi_sync_paid"}`)),
			}, nil
		default:
			return nil, fmt.Errorf("unexpected stripe request %s %s", req.Method, req.URL.Path)
		}
	})
	t.Cleanup(func() {
		http.DefaultTransport = originalTransport
	})

	createListing := doJSONRequestWithHeaders(t, env.router, http.MethodPost, "/api/marketplace/listings", CreateMarketplaceListingRequest{
		DeckID:      1,
		Title:       "Stripe Premium Listing",
		Summary:     "Premium via stripe sync",
		Description: "Checkout session sync should complete the order.",
		Category:    "Medicine",
		PriceMode:   "premium",
		PriceCents:  1900,
		Currency:    "USD",
	}, proHeaders)
	if createListing.Code != http.StatusCreated {
		t.Fatalf("expected premium listing create 201, got %d (%s)", createListing.Code, createListing.Body.String())
	}
	detail := decodeJSON[MarketplaceListingDetail](t, createListing)

	publish := doJSONRequestWithHeaders(t, env.router, http.MethodPost, fmt.Sprintf("/api/marketplace/listings/%s/publish", detail.Listing.ID), PublishMarketplaceListingRequest{
		ChangeSummary: "Premium v1",
	}, proHeaders)
	if publish.Code != http.StatusCreated {
		t.Fatalf("expected premium listing publish 201, got %d (%s)", publish.Code, publish.Body.String())
	}

	checkoutSessionID := "cs_test_sync_paid"
	order := &MarketplaceOrder{
		ID:                        newID("mord"),
		ListingID:                 detail.Listing.ID,
		ListingVersionNumber:      1,
		BuyerUserID:               memberClient.user.ID,
		BuyerWorkspaceID:          memberClient.workspace.ID,
		CreatorUserID:             ownerSession.UserID,
		CreatorAccountID:          creatorAccount.ID,
		Provider:                  "stripe",
		ProviderCheckoutSessionID: checkoutSessionID,
		Status:                    "pending",
		AmountCents:               1900,
		Currency:                  "USD",
		PlatformFeeCents:          marketplacePlatformFeeCents(1900, cfg.Stripe.PlatformFeeBasisPts),
		CreatedAt:                 time.Now(),
		UpdatedAt:                 time.Now(),
	}
	order.CreatorAmountCents = order.AmountCents - order.PlatformFeeCents
	if err := env.store.CreateMarketplaceOrder(order); err != nil {
		t.Fatalf("failed to seed marketplace order: %v", err)
	}

	syncRR := doRawRequest(memberClient.router, http.MethodPost, fmt.Sprintf("/api/marketplace/checkout/sessions/%s/sync", checkoutSessionID), "")
	if syncRR.Code != http.StatusOK {
		t.Fatalf("expected checkout session sync 200, got %d (%s)", syncRR.Code, syncRR.Body.String())
	}
	syncResp := decodeJSON[MarketplaceCheckoutResponse](t, syncRR)
	if !syncResp.Completed || syncResp.License == nil || syncResp.Order.Status != "paid" {
		t.Fatalf("expected sync to complete order and grant license, got %+v", syncResp)
	}

	installAfterSync := doJSONRequest(t, memberClient.router, http.MethodPost, fmt.Sprintf("/api/marketplace/listings/%s/installs", detail.Listing.Slug), InstallMarketplaceListingRequest{
		DestinationWorkspaceID: memberClient.workspace.ID,
	})
	if installAfterSync.Code != http.StatusCreated {
		t.Fatalf("expected premium install after sync 201, got %d (%s)", installAfterSync.Code, installAfterSync.Body.String())
	}
}

func TestAPI_MarketplaceWebhookCompletesOrderAndUpdatesCreatorAccount(t *testing.T) {
	cfg := mustLocalAppConfig()
	cfg.Environment = "production"
	cfg.Stripe.SecretKey = "sk_test_marketplace"
	cfg.Stripe.WebhookSecret = "whsec_marketplace_test"
	env := setupAPITestEnvWithConfig(t, cfg)
	memberClient := createAuthenticatedIsolatedTestClient(t, env, "webhook-buyer@example.com", "Webhook Buyer")
	proHeaders := map[string]string{"X-Vutadex-Plan": "pro"}

	createNoteForTest(t, env, CreateNoteRequest{
		TypeID: "Basic",
		DeckID: 1,
		FieldVals: map[string]string{
			"Front": "Webhook source card",
			"Back":  "Webhook answer",
		},
	}, nil)

	ownerSessionReq := httptest.NewRequest(http.MethodGet, "/", nil)
	ownerSessionReq.Header.Set("Cookie", env.authCookie)
	ownerSession := env.handler.sessionFromRequest(ownerSessionReq)

	creatorAccount := &MarketplaceCreatorAccount{
		ID:                newID("mca"),
		UserID:            ownerSession.UserID,
		WorkspaceID:       ownerSession.WorkspaceID,
		Provider:          "stripe",
		ProviderAccountID: "acct_marketplace_webhook",
		OnboardingStatus:  "pending",
		DetailsSubmitted:  false,
		ChargesEnabled:    false,
		PayoutsEnabled:    false,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}
	if err := env.store.UpsertMarketplaceCreatorAccount(creatorAccount); err != nil {
		t.Fatalf("failed to seed creator account: %v", err)
	}

	originalTransport := http.DefaultTransport
	http.DefaultTransport = roundTripFunc(func(req *http.Request) (*http.Response, error) {
		if req.URL.Host != "api.stripe.com" {
			return nil, fmt.Errorf("unexpected host %s", req.URL.Host)
		}
		switch {
		case req.Method == http.MethodGet && req.URL.Path == "/v1/accounts/acct_marketplace_webhook":
			return &http.Response{
				StatusCode: http.StatusOK,
				Header:     make(http.Header),
				Body:       io.NopCloser(strings.NewReader(`{"id":"acct_marketplace_webhook","details_submitted":true,"charges_enabled":true,"payouts_enabled":true}`)),
			}, nil
		case req.Method == http.MethodPost && req.URL.Path == "/v1/accounts/acct_marketplace_webhook/login_links":
			return &http.Response{
				StatusCode: http.StatusOK,
				Header:     make(http.Header),
				Body:       io.NopCloser(strings.NewReader(`{"url":"https://dashboard.stripe.test/acct_marketplace_webhook"}`)),
			}, nil
		default:
			return nil, fmt.Errorf("unexpected stripe request %s %s", req.Method, req.URL.Path)
		}
	})
	t.Cleanup(func() {
		http.DefaultTransport = originalTransport
	})

	createListing := doJSONRequestWithHeaders(t, env.router, http.MethodPost, "/api/marketplace/listings", CreateMarketplaceListingRequest{
		DeckID:      1,
		Title:       "Webhook Premium Listing",
		Summary:     "Premium via webhook completion",
		Description: "Webhook should complete checkout and update creator state.",
		Category:    "Engineering",
		PriceMode:   "premium",
		PriceCents:  2900,
		Currency:    "USD",
	}, proHeaders)
	if createListing.Code != http.StatusCreated {
		t.Fatalf("expected premium listing create 201, got %d (%s)", createListing.Code, createListing.Body.String())
	}
	detail := decodeJSON[MarketplaceListingDetail](t, createListing)

	creatorAccount.DetailsSubmitted = true
	creatorAccount.ChargesEnabled = true
	creatorAccount.PayoutsEnabled = true
	creatorAccount.OnboardingStatus = "active"
	creatorAccount.UpdatedAt = time.Now()
	if err := env.store.UpsertMarketplaceCreatorAccount(creatorAccount); err != nil {
		t.Fatalf("failed to activate creator account: %v", err)
	}

	publish := doJSONRequestWithHeaders(t, env.router, http.MethodPost, fmt.Sprintf("/api/marketplace/listings/%s/publish", detail.Listing.ID), PublishMarketplaceListingRequest{
		ChangeSummary: "Premium v1",
	}, proHeaders)
	if publish.Code != http.StatusCreated {
		t.Fatalf("expected premium listing publish 201, got %d (%s)", publish.Code, publish.Body.String())
	}

	checkoutSessionID := "cs_test_webhook_paid"
	order := &MarketplaceOrder{
		ID:                        newID("mord"),
		ListingID:                 detail.Listing.ID,
		ListingVersionNumber:      1,
		BuyerUserID:               memberClient.user.ID,
		BuyerWorkspaceID:          memberClient.workspace.ID,
		CreatorUserID:             ownerSession.UserID,
		CreatorAccountID:          creatorAccount.ID,
		Provider:                  "stripe",
		ProviderCheckoutSessionID: checkoutSessionID,
		Status:                    "pending",
		AmountCents:               2900,
		Currency:                  "USD",
		PlatformFeeCents:          marketplacePlatformFeeCents(2900, cfg.Stripe.PlatformFeeBasisPts),
		CreatedAt:                 time.Now(),
		UpdatedAt:                 time.Now(),
	}
	order.CreatorAmountCents = order.AmountCents - order.PlatformFeeCents
	if err := env.store.CreateMarketplaceOrder(order); err != nil {
		t.Fatalf("failed to seed marketplace order: %v", err)
	}

	makeStripeSignature := func(payload []byte, secret string) string {
		timestamp := strconv.FormatInt(time.Now().Unix(), 10)
		mac := hmac.New(sha256.New, []byte(secret))
		mac.Write([]byte(timestamp))
		mac.Write([]byte("."))
		mac.Write(payload)
		return fmt.Sprintf("t=%s,v1=%s", timestamp, hex.EncodeToString(mac.Sum(nil)))
	}

	accountPayload := []byte(`{"id":"evt_account","type":"account.updated","data":{"object":{"id":"acct_marketplace_webhook","details_submitted":true,"charges_enabled":true,"payouts_enabled":true}}}`)
	accountReq := httptest.NewRequest(http.MethodPost, "/api/marketplace/webhook", bytes.NewReader(accountPayload))
	accountReq.Header.Set("Stripe-Signature", makeStripeSignature(accountPayload, cfg.Stripe.WebhookSecret))
	accountRR := httptest.NewRecorder()
	env.router.ServeHTTP(accountRR, accountReq)
	if accountRR.Code != http.StatusOK {
		t.Fatalf("expected account webhook 200, got %d (%s)", accountRR.Code, accountRR.Body.String())
	}

	creatorAfterWebhook, err := env.store.GetMarketplaceCreatorAccount(creatorAccount.ID)
	if err != nil {
		t.Fatalf("failed to reload creator account: %v", err)
	}
	if !creatorAfterWebhook.DetailsSubmitted || !creatorAfterWebhook.ChargesEnabled || !creatorAfterWebhook.PayoutsEnabled {
		t.Fatalf("expected creator account to be updated by webhook, got %+v", creatorAfterWebhook)
	}

	checkoutPayload := []byte(fmt.Sprintf(`{"id":"evt_checkout","type":"checkout.session.completed","data":{"object":{"id":"%s","payment_intent":"pi_webhook_paid","status":"complete","payment_status":"paid"}}}`, checkoutSessionID))
	checkoutReq := httptest.NewRequest(http.MethodPost, "/api/marketplace/webhook", bytes.NewReader(checkoutPayload))
	checkoutReq.Header.Set("Stripe-Signature", makeStripeSignature(checkoutPayload, cfg.Stripe.WebhookSecret))
	checkoutRR := httptest.NewRecorder()
	env.router.ServeHTTP(checkoutRR, checkoutReq)
	if checkoutRR.Code != http.StatusOK {
		t.Fatalf("expected checkout webhook 200, got %d (%s)", checkoutRR.Code, checkoutRR.Body.String())
	}

	reloadedOrder, err := env.store.GetMarketplaceOrder(order.ID)
	if err != nil {
		t.Fatalf("failed to reload order: %v", err)
	}
	if reloadedOrder.Status != "paid" || reloadedOrder.ProviderPaymentIntentID != "pi_webhook_paid" {
		t.Fatalf("expected webhook to mark order paid, got %+v", reloadedOrder)
	}

	license, err := env.store.GetMarketplaceLicense(detail.Listing.ID, memberClient.user.ID)
	if err != nil {
		t.Fatalf("failed to load marketplace license: %v", err)
	}
	if license.Status != "active" || license.OrderID != order.ID {
		t.Fatalf("expected active marketplace license after webhook, got %+v", license)
	}

	payout, err := env.store.GetMarketplacePayoutByOrder(order.ID)
	if err != nil {
		t.Fatalf("failed to load marketplace payout: %v", err)
	}
	if payout.Status != "pending" || payout.AmountCents != order.CreatorAmountCents {
		t.Fatalf("expected pending payout after webhook, got %+v", payout)
	}

	installAfterWebhook := doJSONRequest(t, memberClient.router, http.MethodPost, fmt.Sprintf("/api/marketplace/listings/%s/installs", detail.Listing.Slug), InstallMarketplaceListingRequest{
		DestinationWorkspaceID: memberClient.workspace.ID,
	})
	if installAfterWebhook.Code != http.StatusCreated {
		t.Fatalf("expected premium install after webhook 201, got %d (%s)", installAfterWebhook.Code, installAfterWebhook.Body.String())
	}
}

func TestAPI_MarketplaceFreeListingsPublishInstallAndUpdateAcrossCollections(t *testing.T) {
	env := setupAPITestEnv(t)
	memberClient := createAuthenticatedIsolatedTestClient(t, env, "marketplace-member@example.com", "Marketplace Member")
	proHeaders := map[string]string{"X-Vutadex-Plan": "pro"}

	createNoteForTest(t, env, CreateNoteRequest{
		TypeID: "Basic",
		DeckID: 1,
		FieldVals: map[string]string{
			"Front": "Marketplace source card 1",
			"Back":  "Marketplace answer 1",
		},
	}, nil)

	createListing := doJSONRequestWithHeaders(t, env.router, http.MethodPost, "/api/marketplace/listings", CreateMarketplaceListingRequest{
		DeckID:      1,
		Title:       "USMLE Foundations",
		Summary:     "A versioned source deck for marketplace installs.",
		Description: "Free listing installs should create workspace-local copies with source attribution.",
		Category:    "Medicine",
		Tags:        []string{"usmle", "anki-alternative"},
		PriceMode:   "free",
	}, proHeaders)
	if createListing.Code != http.StatusCreated {
		t.Fatalf("expected create marketplace listing 201, got %d (%s)", createListing.Code, createListing.Body.String())
	}
	draftDetail := decodeJSON[MarketplaceListingDetail](t, createListing)
	listingID := draftDetail.Listing.ID
	listingSlug := draftDetail.Listing.Slug

	publicBeforePublish := doRawRequest(memberClient.router, http.MethodGet, "/api/marketplace/listings", "")
	if publicBeforePublish.Code != http.StatusOK {
		t.Fatalf("expected marketplace public list 200, got %d (%s)", publicBeforePublish.Code, publicBeforePublish.Body.String())
	}
	if listings := decodeJSON[[]MarketplaceListingSummary](t, publicBeforePublish); len(listings) != 0 {
		t.Fatalf("expected draft listing to be hidden from public catalog, got %+v", listings)
	}

	creatorMineBeforePublish := doRawRequest(env.router, http.MethodGet, "/api/marketplace/listings?scope=mine", "")
	if creatorMineBeforePublish.Code != http.StatusOK {
		t.Fatalf("expected marketplace mine list 200, got %d (%s)", creatorMineBeforePublish.Code, creatorMineBeforePublish.Body.String())
	}
	mineListings := decodeJSON[[]MarketplaceListingSummary](t, creatorMineBeforePublish)
	if len(mineListings) != 1 || mineListings[0].Status != "draft" {
		t.Fatalf("expected creator mine list to include one draft listing, got %+v", mineListings)
	}

	draftDetailForMember := doRawRequest(memberClient.router, http.MethodGet, fmt.Sprintf("/api/marketplace/listings/%s", listingSlug), "")
	if draftDetailForMember.Code != http.StatusNotFound {
		t.Fatalf("expected draft listing detail to be hidden from non-creators, got %d (%s)", draftDetailForMember.Code, draftDetailForMember.Body.String())
	}

	publishV1 := doJSONRequestWithHeaders(t, env.router, http.MethodPost, fmt.Sprintf("/api/marketplace/listings/%s/publish", listingID), PublishMarketplaceListingRequest{
		ChangeSummary: "Initial free catalog release",
	}, proHeaders)
	if publishV1.Code != http.StatusCreated {
		t.Fatalf("expected publish v1 201, got %d (%s)", publishV1.Code, publishV1.Body.String())
	}
	version1 := decodeJSON[MarketplaceListingVersion](t, publishV1)
	if version1.VersionNumber != 1 || version1.NoteCount != 1 || version1.CardCount != 1 {
		t.Fatalf("expected v1 marketplace metadata to reflect source deck, got %+v", version1)
	}

	publicAfterPublish := doRawRequest(memberClient.router, http.MethodGet, "/api/marketplace/listings", "")
	if publicAfterPublish.Code != http.StatusOK {
		t.Fatalf("expected marketplace public list after publish 200, got %d (%s)", publicAfterPublish.Code, publicAfterPublish.Body.String())
	}
	publicListings := decodeJSON[[]MarketplaceListingSummary](t, publicAfterPublish)
	if len(publicListings) != 1 || publicListings[0].Slug != listingSlug || publicListings[0].LatestVersionNumber != 1 {
		t.Fatalf("expected published marketplace listing in public catalog, got %+v", publicListings)
	}

	memberDetailRR := doRawRequest(memberClient.router, http.MethodGet, fmt.Sprintf("/api/marketplace/listings/%s", listingSlug), "")
	if memberDetailRR.Code != http.StatusOK {
		t.Fatalf("expected published marketplace detail 200, got %d (%s)", memberDetailRR.Code, memberDetailRR.Body.String())
	}
	memberDetail := decodeJSON[MarketplaceListingDetail](t, memberDetailRR)
	if memberDetail.LatestVersion == nil || memberDetail.LatestVersion.VersionNumber != 1 {
		t.Fatalf("expected published detail to expose latest version metadata, got %+v", memberDetail.LatestVersion)
	}

	installRR := doJSONRequest(t, memberClient.router, http.MethodPost, fmt.Sprintf("/api/marketplace/listings/%s/installs", listingSlug), InstallMarketplaceListingRequest{
		DestinationWorkspaceID: memberClient.workspace.ID,
	})
	if installRR.Code != http.StatusCreated {
		t.Fatalf("expected marketplace install 201, got %d (%s)", installRR.Code, installRR.Body.String())
	}
	memberInstall := decodeJSON[MarketplaceInstall](t, installRR)
	if memberInstall.SourceVersionNumber != 1 || memberInstall.Status != "active" {
		t.Fatalf("expected active marketplace install on version 1, got %+v", memberInstall)
	}

	var destinationCollectionID string
	if err := env.store.db.QueryRow(`SELECT collection_id FROM decks WHERE id = ?`, memberInstall.InstalledDeckID).Scan(&destinationCollectionID); err != nil {
		t.Fatalf("failed to load installed marketplace deck collection: %v", err)
	}
	if destinationCollectionID != memberClient.workspace.CollectionID {
		t.Fatalf("expected marketplace install to land in destination workspace collection %q, got %q", memberClient.workspace.CollectionID, destinationCollectionID)
	}

	var sourceCollectionID string
	if err := env.store.db.QueryRow(`SELECT collection_id FROM decks WHERE id = 1`).Scan(&sourceCollectionID); err != nil {
		t.Fatalf("failed to load marketplace source deck collection: %v", err)
	}
	if sourceCollectionID == destinationCollectionID {
		t.Fatalf("expected marketplace source and destination collections to differ, both were %q", sourceCollectionID)
	}

	memberDetailAfterInstallRR := doRawRequest(memberClient.router, http.MethodGet, fmt.Sprintf("/api/marketplace/listings/%s", listingSlug), "")
	if memberDetailAfterInstallRR.Code != http.StatusOK {
		t.Fatalf("expected marketplace detail after install 200, got %d (%s)", memberDetailAfterInstallRR.Code, memberDetailAfterInstallRR.Body.String())
	}
	memberDetailAfterInstall := decodeJSON[MarketplaceListingDetail](t, memberDetailAfterInstallRR)
	if memberDetailAfterInstall.CurrentUserInstall == nil || memberDetailAfterInstall.CurrentUserInstall.SourceVersionNumber != 1 {
		t.Fatalf("expected current user install metadata after install, got %+v", memberDetailAfterInstall.CurrentUserInstall)
	}

	var memberCardID int64
	if err := env.store.db.QueryRow(`SELECT id FROM cards WHERE deck_id = ? ORDER BY id ASC LIMIT 1`, memberInstall.InstalledDeckID).Scan(&memberCardID); err != nil {
		t.Fatalf("failed to load installed marketplace deck card: %v", err)
	}
	memberAnswer := doJSONRequest(t, memberClient.router, http.MethodPost, fmt.Sprintf("/api/cards/%d/answer", memberCardID), AnswerCardRequest{
		Rating:      3,
		TimeTakenMs: 700,
	})
	if memberAnswer.Code != http.StatusOK {
		t.Fatalf("expected marketplace installed deck answer 200, got %d (%s)", memberAnswer.Code, memberAnswer.Body.String())
	}

	sourceDueAfterInstallAnswer := doRawRequest(env.router, http.MethodGet, "/api/decks/1/due?limit=10", "")
	if sourceDueAfterInstallAnswer.Code != http.StatusOK {
		t.Fatalf("expected source deck due queue 200, got %d (%s)", sourceDueAfterInstallAnswer.Code, sourceDueAfterInstallAnswer.Body.String())
	}
	if cards := decodeJSON[[]Card](t, sourceDueAfterInstallAnswer); len(cards) != 1 {
		t.Fatalf("expected source deck due queue to remain unchanged after marketplace install study, got %d", len(cards))
	}

	createNoteForTest(t, env, CreateNoteRequest{
		TypeID: "Basic",
		DeckID: 1,
		FieldVals: map[string]string{
			"Front": "Marketplace source card 2",
			"Back":  "Marketplace answer 2",
		},
	}, nil)

	publishV2 := doJSONRequestWithHeaders(t, env.router, http.MethodPost, fmt.Sprintf("/api/marketplace/listings/%s/publish", listingID), PublishMarketplaceListingRequest{
		ChangeSummary: "Added a second marketplace source card",
	}, proHeaders)
	if publishV2.Code != http.StatusCreated {
		t.Fatalf("expected publish v2 201, got %d (%s)", publishV2.Code, publishV2.Body.String())
	}
	version2 := decodeJSON[MarketplaceListingVersion](t, publishV2)
	if version2.VersionNumber != 2 || version2.NoteCount != 2 || version2.CardCount != 2 {
		t.Fatalf("expected v2 marketplace metadata to reflect updated source deck, got %+v", version2)
	}

	memberDetailBeforeUpdateRR := doRawRequest(memberClient.router, http.MethodGet, fmt.Sprintf("/api/marketplace/listings/%s", listingSlug), "")
	if memberDetailBeforeUpdateRR.Code != http.StatusOK {
		t.Fatalf("expected marketplace detail before update 200, got %d (%s)", memberDetailBeforeUpdateRR.Code, memberDetailBeforeUpdateRR.Body.String())
	}
	memberDetailBeforeUpdate := decodeJSON[MarketplaceListingDetail](t, memberDetailBeforeUpdateRR)
	if !memberDetailBeforeUpdate.UpdateAvailable {
		t.Fatalf("expected marketplace detail updateAvailable=true after publishing v2")
	}
	if memberDetailBeforeUpdate.CurrentUserInstall == nil || memberDetailBeforeUpdate.CurrentUserInstall.SourceVersionNumber != 1 {
		t.Fatalf("expected current user install to still point at version 1 before update, got %+v", memberDetailBeforeUpdate.CurrentUserInstall)
	}

	updateInstallRR := doJSONRequest(t, memberClient.router, http.MethodPost, fmt.Sprintf("/api/marketplace/listings/%s/installs/%s/update", listingSlug, memberInstall.ID), UpdateMarketplaceInstallRequest{})
	if updateInstallRR.Code != http.StatusOK {
		t.Fatalf("expected marketplace install update 200, got %d (%s)", updateInstallRR.Code, updateInstallRR.Body.String())
	}
	memberInstallV2 := decodeJSON[MarketplaceInstall](t, updateInstallRR)
	if memberInstallV2.ID == memberInstall.ID {
		t.Fatalf("expected marketplace install update to create a fresh install record")
	}
	if memberInstallV2.SourceVersionNumber != 2 || memberInstallV2.Status != "active" {
		t.Fatalf("expected active version 2 marketplace install after update, got %+v", memberInstallV2)
	}

	oldInstall, err := env.store.GetMarketplaceInstall(memberInstall.ID)
	if err != nil {
		t.Fatalf("failed to reload superseded marketplace install: %v", err)
	}
	if oldInstall.Status != "superseded" || oldInstall.SupersededByInstall != memberInstallV2.ID {
		t.Fatalf("expected original marketplace install to be superseded by the new install, got %+v", oldInstall)
	}

	oldNotes, oldCards, err := env.store.GetDeckContentSummary(memberInstall.InstalledDeckID)
	if err != nil {
		t.Fatalf("failed to read original marketplace install summary: %v", err)
	}
	if oldNotes != 1 || oldCards != 1 {
		t.Fatalf("expected original marketplace install copy to remain on v1, got notes=%d cards=%d", oldNotes, oldCards)
	}
	newNotes, newCards, err := env.store.GetDeckContentSummary(memberInstallV2.InstalledDeckID)
	if err != nil {
		t.Fatalf("failed to read updated marketplace install summary: %v", err)
	}
	if newNotes != 2 || newCards != 2 {
		t.Fatalf("expected updated marketplace install copy to reflect v2 content, got notes=%d cards=%d", newNotes, newCards)
	}

	removeInstallRR := doJSONRequest(t, memberClient.router, http.MethodDelete, fmt.Sprintf("/api/marketplace/listings/%s/installs/%s", listingSlug, memberInstallV2.ID), struct{}{})
	if removeInstallRR.Code != http.StatusNoContent {
		t.Fatalf("expected remove marketplace install 204, got %d (%s)", removeInstallRR.Code, removeInstallRR.Body.String())
	}
	removedInstall, err := env.store.GetMarketplaceInstall(memberInstallV2.ID)
	if err != nil {
		t.Fatalf("failed to reload removed marketplace install: %v", err)
	}
	if removedInstall.Status != "removed" {
		t.Fatalf("expected removed marketplace install status=removed, got %+v", removedInstall)
	}
	if _, err := env.store.GetDeck(memberInstallV2.InstalledDeckID); err == nil {
		t.Fatalf("expected removing a marketplace install to delete its copied deck %d", memberInstallV2.InstalledDeckID)
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

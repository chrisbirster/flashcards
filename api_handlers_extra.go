package main

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const (
	sessionCookieName       = "vutadex_session"
	oauthStateCookieName    = "vutadex_oauth_state"
	oauthVerifierCookieName = "vutadex_oauth_verifier"
	oauthCookieTTL          = 10 * time.Minute
)

type contextKey string

const sessionContextKey contextKey = "vutadex_session"

type billingCheckoutRequest struct {
	Plan Plan `json:"plan"`
}

type billingWebhookRequest struct {
	WorkspaceID        string `json:"workspaceId,omitempty"`
	OrganizationID     string `json:"organizationId,omitempty"`
	Plan               Plan   `json:"plan"`
	Status             string `json:"status"`
	Provider           string `json:"provider,omitempty"`
	ProviderCustomerID string `json:"providerCustomerId,omitempty"`
	ProviderSubID      string `json:"providerSubscriptionId,omitempty"`
	ProviderEventID    string `json:"providerEventId,omitempty"`
	CurrentPeriodEnd   string `json:"currentPeriodEnd,omitempty"`
}

type googleUserInfo struct {
	Sub     string `json:"sub"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

func registerAPIRoutes(r chi.Router, handler *APIHandler) {
	r.Use(handler.SessionMiddleware)

	r.Get("/health", handler.HealthCheck)
	r.Get("/auth/session", handler.GetAuthSession)
	r.Post("/auth/otp/request", handler.RequestOTP)
	r.Post("/auth/otp/verify", handler.VerifyOTP)
	r.Post("/auth/logout", handler.Logout)
	r.Post("/marketplace/webhook", handler.MarketplaceWebhook)

	r.Group(func(r chi.Router) {
		r.Use(handler.RequireAuthenticatedUser)

		r.Get("/collection", handler.GetCollection)
		r.Get("/dashboard", handler.GetDashboard)
		r.Post("/import", handler.ImportNotes)

		r.Get("/decks", handler.ListDecks)
		r.Post("/decks", handler.CreateDeck)
		r.Get("/decks/{id}", handler.GetDeck)
		r.Patch("/decks/{id}", handler.UpdateDeck)
		r.Delete("/decks/{id}", handler.DeleteDeck)
		r.Get("/decks/{id}/stats", handler.GetDeckStats)
		r.Get("/decks/{deckId}/notes", handler.GetDeckNotes)
		r.Get("/decks/{deckId}/due", handler.GetDueCards)
		r.Post("/decks/{deckId}/share", handler.CreateDeckShare)
		r.Delete("/decks/{deckId}/share", handler.DeleteDeckShare)

		r.Get("/note-types", handler.ListNoteTypes)
		r.Get("/note-types/{name}", handler.GetNoteType)
		r.Post("/note-types/{name}/fields", handler.AddField)
		r.Patch("/note-types/{name}/fields/rename", handler.RenameField)
		r.Delete("/note-types/{name}/fields", handler.RemoveField)
		r.Put("/note-types/{name}/fields/reorder", handler.ReorderFields)
		r.Put("/note-types/{name}/sort-field", handler.SetSortField)
		r.Put("/note-types/{name}/fields/options", handler.SetFieldOptions)
		r.Post("/note-types/{name}/templates", handler.CreateTemplate)
		r.Patch("/note-types/{name}/templates/{templateName}", handler.UpdateTemplate)
		r.Delete("/note-types/{name}/templates/{templateName}", handler.DeleteTemplate)

		r.Get("/notes", handler.ListNotes)
		r.Post("/notes", handler.CreateNote)
		r.Get("/notes/{id}", handler.GetNote)
		r.Patch("/notes/{id}", handler.UpdateNote)
		r.Delete("/notes/{id}", handler.DeleteNote)
		r.Post("/notes/check-duplicate", handler.CheckDuplicate)

		r.Get("/cards/{id}", handler.GetCard)
		r.Post("/cards/{id}/answer", handler.AnswerCard)
		r.Patch("/cards/{id}", handler.UpdateCard)
		r.Get("/cards/empty", handler.FindEmptyCards)
		r.Post("/cards/empty/delete", handler.DeleteEmptyCards)

		r.Get("/entitlements", handler.GetEntitlements)
		r.Post("/onboarding/import-local-collection", handler.ImportLocalCollection)
		r.Post("/ai/card-suggestions", handler.GenerateCardSuggestions)
		r.Post("/study-sessions", handler.CreateStudySession)
		r.Patch("/study-sessions/{id}", handler.UpdateStudySession)
		r.Get("/analytics/overview", handler.GetStudyAnalyticsOverview)

		r.Post("/billing/checkout", handler.BillingCheckout)
		r.Post("/billing/portal", handler.BillingPortal)
		r.Post("/billing/webhook", handler.BillingWebhook)

		r.Post("/orgs", handler.CreateOrganization)
		r.Post("/orgs/{orgId}/members", handler.AddOrganizationMember)

		r.Get("/study-groups", handler.ListStudyGroups)
		r.Post("/study-groups", handler.CreateStudyGroup)
		r.Post("/study-groups/join", handler.JoinStudyGroup)
		r.Route("/study-groups/{id}", func(r chi.Router) {
			r.Get("/", handler.GetStudyGroup)
			r.Patch("/", handler.UpdateStudyGroup)
			r.Delete("/", handler.DeleteStudyGroup)
			r.Post("/members", handler.InviteStudyGroupMember)
			r.Patch("/members/{memberId}", handler.UpdateStudyGroupMember)
			r.Delete("/members/{memberId}", handler.DeleteStudyGroupMember)
			r.Get("/versions", handler.ListStudyGroupVersions)
			r.Post("/versions", handler.PublishStudyGroupVersion)
			r.Post("/installs", handler.InstallStudyGroupDeck)
			r.Post("/installs/{installId}/update", handler.UpdateStudyGroupInstall)
			r.Delete("/installs/{installId}", handler.RemoveStudyGroupInstall)
			r.Get("/dashboard", handler.GetStudyGroupDashboard)
		})

		r.Route("/marketplace", func(r chi.Router) {
			r.Get("/creator-account/status", handler.GetMarketplaceCreatorAccountStatus)
			r.Post("/creator-account/start", handler.StartMarketplaceCreatorAccount)
			r.Post("/checkout/sessions/{sessionId}/sync", handler.SyncMarketplaceCheckoutSession)
			r.Get("/listings", handler.ListMarketplaceListings)
			r.Post("/listings", handler.CreateMarketplaceListing)
			r.Route("/listings/{ref}", func(r chi.Router) {
				r.Get("/", handler.GetMarketplaceListing)
				r.Patch("/", handler.UpdateMarketplaceListing)
				r.Delete("/", handler.DeleteMarketplaceListing)
				r.Post("/publish", handler.PublishMarketplaceListing)
				r.Post("/checkout", handler.CheckoutMarketplaceListing)
				r.Post("/installs", handler.InstallMarketplaceListing)
				r.Post("/installs/{installId}/update", handler.UpdateMarketplaceInstall)
				r.Delete("/installs/{installId}", handler.RemoveMarketplaceInstall)
			})
		})

		r.Post("/backups", handler.CreateBackup)
		r.Get("/backups", handler.ListBackups)
		r.Post("/backups/restore", handler.RestoreBackup)
	})
}

func googleOAuthConfigured() bool {
	return strings.TrimSpace(os.Getenv("VUTADEX_GOOGLE_CLIENT_ID")) != "" &&
		strings.TrimSpace(os.Getenv("VUTADEX_GOOGLE_CLIENT_SECRET")) != "" &&
		strings.TrimSpace(os.Getenv("VUTADEX_GOOGLE_REDIRECT_URL")) != ""
}

func googleOAuthConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     strings.TrimSpace(os.Getenv("VUTADEX_GOOGLE_CLIENT_ID")),
		ClientSecret: strings.TrimSpace(os.Getenv("VUTADEX_GOOGLE_CLIENT_SECRET")),
		RedirectURL:  strings.TrimSpace(os.Getenv("VUTADEX_GOOGLE_REDIRECT_URL")),
		Scopes:       []string{"openid", "profile", "email"},
		Endpoint:     google.Endpoint,
	}
}

func randomToken() string {
	return newID("token")
}

func pkceChallenge(verifier string) string {
	sum := sha256.Sum256([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(sum[:])
}

func (h *APIHandler) writeCookie(w http.ResponseWriter, name, value string, expiresAt time.Time) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   h.config.Cookie.Secure,
		Domain:   h.config.Cookie.Domain,
		Expires:  expiresAt,
		MaxAge:   int(time.Until(expiresAt).Seconds()),
	})
}

func (h *APIHandler) clearCookie(w http.ResponseWriter, name string) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   h.config.Cookie.Secure,
		Domain:   h.config.Cookie.Domain,
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
	})
}

func (h *APIHandler) sessionFromRequest(r *http.Request) *SessionRecord {
	if session, ok := r.Context().Value(sessionContextKey).(*SessionRecord); ok {
		return session
	}

	cookie, err := r.Cookie(sessionCookieName)
	if err != nil || cookie.Value == "" {
		return nil
	}

	session, err := h.store.GetSessionRecord(cookie.Value)
	if err != nil || session == nil {
		return nil
	}
	if !session.ExpiresAt.IsZero() && session.ExpiresAt.Before(time.Now()) {
		_ = h.store.DeleteSessionRecord(session.ID)
		return nil
	}
	return session
}

func (h *APIHandler) SessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session := h.sessionFromRequest(r)
		if session != nil && session.UserID != "" {
			now := time.Now()
			session.LastSeenAt = now
			session.ExpiresAt = now.Add(h.config.SessionTTL)
			if err := h.store.TouchSessionRecord(session.ID, session.ExpiresAt, session.LastSeenAt); err == nil {
				h.writeCookie(w, sessionCookieName, session.ID, session.ExpiresAt)
			}
		}
		ctx := context.WithValue(r.Context(), sessionContextKey, session)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h *APIHandler) RequireAuthenticatedUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session := h.sessionFromRequest(r)
		if session == nil || session.UserID == "" {
			respondAPIError(w, http.StatusUnauthorized, "auth_required", "You must be signed in to access this resource")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (h *APIHandler) planForRequest(r *http.Request, session *SessionRecord) Plan {
	if session != nil && session.WorkspaceID != "" {
		subscription, err := h.store.GetSubscriptionForWorkspace(session.WorkspaceID)
		if err == nil && subscription.Status == "active" {
			return subscription.Plan
		}
	}
	return resolvePlanFromRequest(r, session)
}

func (h *APIHandler) workspaceForSession(session *SessionRecord) (*Workspace, error) {
	if session == nil || strings.TrimSpace(session.WorkspaceID) == "" {
		return nil, nil
	}
	return h.store.GetWorkspaceRecord(session.WorkspaceID)
}

func (h *APIHandler) collectionIDForRequest(r *http.Request) string {
	session := h.sessionFromRequest(r)
	if workspace, err := h.workspaceForSession(session); err == nil && workspace != nil && strings.TrimSpace(workspace.CollectionID) != "" {
		return workspace.CollectionID
	}
	return h.collectionID
}

func (h *APIHandler) collectionForRequest(r *http.Request) (*Collection, string, error) {
	collectionID := h.collectionIDForRequest(r)
	col, err := h.store.GetCollection(collectionID)
	if err != nil {
		return nil, collectionID, err
	}
	if collectionID == h.collectionID {
		if h.collection != nil {
			*h.collection = *col
			return h.collection, collectionID, nil
		}
		h.collection = col
	}
	return col, collectionID, nil
}

func (h *APIHandler) usageForSession(session *SessionRecord) EntitlementUsage {
	usage := EntitlementUsage{
		Decks:      len(h.collection.Decks),
		Notes:      len(h.collection.Notes),
		CardsTotal: len(h.collection.Cards),
	}
	if session == nil || session.WorkspaceID == "" {
		return usage
	}
	if workspace, err := h.workspaceForSession(session); err == nil && workspace != nil {
		if col, err := h.store.GetCollection(workspace.CollectionID); err == nil {
			usage.Decks = len(col.Decks)
			usage.Notes = len(col.Notes)
			usage.CardsTotal = len(col.Cards)
		}
	}

	if count, err := h.store.CountDeckSharesForWorkspace(session.WorkspaceID); err == nil {
		usage.SharedDecks = count
	}
	if count, err := h.store.CountSyncDevicesForWorkspace(session.WorkspaceID); err == nil {
		usage.SyncDevices = count
	}
	if session.UserID != "" {
		if count, err := h.store.CountWorkspacesForUser(session.UserID); err == nil {
			usage.Workspaces = count
		}
	}
	return usage
}

func (h *APIHandler) buildSessionResponse(r *http.Request) AuthSessionResponse {
	session := h.sessionFromRequest(r)
	plan := PlanGuest
	if session != nil {
		plan = h.planForRequest(r, session)
	}
	entitlements := entitlementsForPlan(plan, h.usageForSession(session))

	response := AuthSessionResponse{
		Authenticated:        session != nil && session.UserID != "",
		GoogleAuthConfigured: false,
		OTPAuthEnabled:       true,
		Entitlements:         entitlements,
	}

	if session == nil {
		return response
	}

	if user, err := h.store.GetUserByID(session.UserID); err == nil {
		response.User = user
	}
	if session.WorkspaceID != "" {
		if workspace, err := h.store.GetWorkspaceRecord(session.WorkspaceID); err == nil {
			response.Workspace = workspace
		}
	}

	return response
}

func (h *APIHandler) ensureDefaultWorkspaceForUser(user *User) (*Workspace, error) {
	workspace, err := h.store.GetFirstWorkspaceForUser(user.ID)
	if err == nil {
		return workspace, nil
	}
	if err != sql.ErrNoRows {
		return nil, err
	}

	now := time.Now()
	collectionID := newID("col")
	collectionName := fmt.Sprintf("%s Collection", firstNonEmpty(user.DisplayName, "Vutadex"))
	collection := NewCollection()
	if err := h.store.CreateCollectionRecord(collectionID, collectionName, collection); err != nil {
		return nil, err
	}
	seedCollection, err := h.store.GetCollection(collectionID)
	if err != nil {
		return nil, err
	}
	for _, nt := range builtins() {
		ntCopy := nt
		if err := h.store.CreateNoteType(collectionID, &ntCopy); err != nil {
			return nil, err
		}
	}
	defaultDeck := seedCollection.NewDeck("Default")
	if err := h.store.CreateDeckInCollection(collectionID, defaultDeck); err != nil {
		return nil, err
	}

	workspace = &Workspace{
		ID:           newID("ws"),
		Name:         fmt.Sprintf("%s Workspace", firstNonEmpty(user.DisplayName, "Vutadex")),
		Slug:         fmt.Sprintf("%s-%s", slugify(firstNonEmpty(user.DisplayName, user.Email)), slugify(user.ID[len(user.ID)-6:])),
		CollectionID: collectionID,
		OwnerUserID:  user.ID,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	return workspace, h.store.CreateWorkspaceRecord(workspace)
}

func (h *APIHandler) GetAuthSession(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, h.buildSessionResponse(r))
}

func (h *APIHandler) GetEntitlements(w http.ResponseWriter, r *http.Request) {
	response := h.buildSessionResponse(r)
	respondJSON(w, http.StatusOK, response.Entitlements)
}

func (h *APIHandler) GetDashboard(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, h.buildDashboardResponse(r))
}

func (h *APIHandler) GetDeckNotes(w http.ResponseWriter, r *http.Request) {
	deckID, err := parseIDParam(r, "deckId")
	if err != nil {
		respondAPIError(w, http.StatusBadRequest, "invalid_deck_id", "Invalid deck ID")
		return
	}

	limit := 20
	if limitStr := strings.TrimSpace(r.URL.Query().Get("limit")); limitStr != "" {
		if parsed, err := strconv.Atoi(limitStr); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	var cursorCreatedAt, cursorNoteID int64
	if cursorRaw := strings.TrimSpace(r.URL.Query().Get("cursor")); cursorRaw != "" {
		parts := strings.Split(cursorRaw, ",")
		if len(parts) == 2 {
			cursorCreatedAt, _ = strconv.ParseInt(parts[0], 10, 64)
			cursorNoteID, _ = strconv.ParseInt(parts[1], 10, 64)
		}
	}

	summaries, err := h.store.ListRecentDeckNotes(h.collectionIDForRequest(r), deckID, limit, cursorCreatedAt, cursorNoteID)
	if err != nil {
		respondAPIError(w, http.StatusInternalServerError, "deck_notes_failed", err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"notes": summaries,
	})
}

func (h *APIHandler) StartGoogleAuth(w http.ResponseWriter, r *http.Request) {
	if !googleOAuthConfigured() {
		respondAPIError(w, http.StatusNotImplemented, "google_auth_not_configured", "Google OAuth is not configured")
		return
	}

	config := googleOAuthConfig()
	state := randomToken()
	verifier := randomToken()
	now := time.Now()
	h.writeCookie(w, oauthStateCookieName, state, now.Add(oauthCookieTTL))
	h.writeCookie(w, oauthVerifierCookieName, verifier, now.Add(oauthCookieTTL))

	authURL := config.AuthCodeURL(
		state,
		oauth2.AccessTypeOffline,
		oauth2.SetAuthURLParam("prompt", "consent"),
		oauth2.SetAuthURLParam("code_challenge_method", "S256"),
		oauth2.SetAuthURLParam("code_challenge", pkceChallenge(verifier)),
	)

	if r.Method == http.MethodGet {
		http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"url": authURL})
}

func (h *APIHandler) GoogleAuthCallback(w http.ResponseWriter, r *http.Request) {
	if !googleOAuthConfigured() {
		respondAPIError(w, http.StatusNotImplemented, "google_auth_not_configured", "Google OAuth is not configured")
		return
	}

	stateCookie, err := r.Cookie(oauthStateCookieName)
	if err != nil || stateCookie.Value == "" || stateCookie.Value != r.URL.Query().Get("state") {
		respondAPIError(w, http.StatusBadRequest, "oauth_state_invalid", "OAuth state validation failed")
		return
	}
	verifierCookie, err := r.Cookie(oauthVerifierCookieName)
	if err != nil || verifierCookie.Value == "" {
		respondAPIError(w, http.StatusBadRequest, "oauth_verifier_missing", "OAuth verifier is missing")
		return
	}

	code := strings.TrimSpace(r.URL.Query().Get("code"))
	if code == "" {
		respondAPIError(w, http.StatusBadRequest, "oauth_code_missing", "Google OAuth code is missing")
		return
	}

	config := googleOAuthConfig()
	token, err := config.Exchange(context.Background(), code, oauth2.SetAuthURLParam("code_verifier", verifierCookie.Value))
	if err != nil {
		respondAPIError(w, http.StatusBadGateway, "oauth_exchange_failed", err.Error())
		return
	}

	client := config.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		respondAPIError(w, http.StatusBadGateway, "oauth_userinfo_failed", err.Error())
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		respondAPIError(w, http.StatusBadGateway, "oauth_userinfo_failed", strings.TrimSpace(string(body)))
		return
	}

	var googleUser googleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&googleUser); err != nil {
		respondAPIError(w, http.StatusBadGateway, "oauth_userinfo_invalid", "Failed to decode Google user profile")
		return
	}

	if googleUser.Email == "" || googleUser.Sub == "" {
		respondAPIError(w, http.StatusBadGateway, "oauth_userinfo_invalid", "Google user profile is incomplete")
		return
	}

	now := time.Now()
	user, err := h.store.GetUserByOAuth("google", googleUser.Sub)
	if err == sql.ErrNoRows {
		user = &User{
			ID:          newID("usr"),
			Email:       googleUser.Email,
			DisplayName: firstNonEmpty(googleUser.Name, googleUser.Email),
			AvatarURL:   googleUser.Picture,
			CreatedAt:   now,
			UpdatedAt:   now,
		}
		if err := h.store.CreateUser(user); err != nil {
			respondAPIError(w, http.StatusInternalServerError, "user_create_failed", err.Error())
			return
		}
	} else if err != nil {
		respondAPIError(w, http.StatusInternalServerError, "user_lookup_failed", err.Error())
		return
	}

	if err := h.store.UpsertOAuthIdentity(&OAuthIdentity{
		ID:        newID("oauth"),
		UserID:    user.ID,
		Provider:  "google",
		Subject:   googleUser.Sub,
		Email:     googleUser.Email,
		CreatedAt: now,
	}); err != nil {
		respondAPIError(w, http.StatusInternalServerError, "oauth_identity_failed", err.Error())
		return
	}

	workspace, err := h.ensureDefaultWorkspaceForUser(user)
	if err != nil {
		respondAPIError(w, http.StatusInternalServerError, "workspace_create_failed", err.Error())
		return
	}

	session := &SessionRecord{
		ID:          newID("sess"),
		UserID:      user.ID,
		WorkspaceID: workspace.ID,
		Plan:        PlanFree,
		ExpiresAt:   now.Add(h.config.SessionTTL),
		LastSeenAt:  now,
		CreatedAt:   now,
	}
	if err := h.store.CreateSessionRecord(session); err != nil {
		respondAPIError(w, http.StatusInternalServerError, "session_create_failed", err.Error())
		return
	}

	h.clearCookie(w, oauthStateCookieName)
	h.clearCookie(w, oauthVerifierCookieName)
	h.writeCookie(w, sessionCookieName, session.ID, session.ExpiresAt)

	redirectPath := h.config.AuthSuccessPath
	http.Redirect(w, r, redirectPath, http.StatusTemporaryRedirect)
}

func (h *APIHandler) Logout(w http.ResponseWriter, r *http.Request) {
	session := h.sessionFromRequest(r)
	if session != nil {
		_ = h.store.DeleteSessionRecord(session.ID)
	}
	h.clearCookie(w, sessionCookieName)
	respondJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func (h *APIHandler) ImportLocalCollection(w http.ResponseWriter, r *http.Request) {
	session := h.sessionFromRequest(r)
	if session == nil || session.UserID == "" {
		respondAPIError(w, http.StatusUnauthorized, "auth_required", "You must be signed in to import a local collection")
		return
	}

	var req ImportLocalCollectionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondAPIError(w, http.StatusBadRequest, "invalid_request", "Invalid request body")
		return
	}

	if len(req.Collection.Decks) > 0 || len(req.Collection.Notes) > 0 || len(req.Collection.Cards) > 0 {
		respondJSON(w, http.StatusAccepted, map[string]string{
			"message": "Local collection import has been accepted but full browser-local migration is not wired in this tranche.",
		})
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{
		"message": "No local collection content was provided.",
	})
}

func (h *APIHandler) CreateOrganization(w http.ResponseWriter, r *http.Request) {
	session := h.sessionFromRequest(r)
	if session == nil || session.UserID == "" {
		respondAPIError(w, http.StatusUnauthorized, "auth_required", "You must be signed in to create an organization")
		return
	}

	var req CreateOrganizationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondAPIError(w, http.StatusBadRequest, "invalid_request", "Invalid request body")
		return
	}
	if strings.TrimSpace(req.Name) == "" {
		respondAPIError(w, http.StatusBadRequest, "invalid_name", "Organization name is required")
		return
	}

	now := time.Now()
	org := &Organization{
		ID:        newID("org"),
		Name:      strings.TrimSpace(req.Name),
		Slug:      firstNonEmpty(strings.TrimSpace(req.Slug), fmt.Sprintf("%s-%s", slugify(req.Name), slugify(newID("o")[2:8]))),
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := h.store.CreateOrganizationRecord(org); err != nil {
		respondAPIError(w, http.StatusInternalServerError, "organization_create_failed", err.Error())
		return
	}

	member := &OrganizationMember{
		ID:             newID("orgmem"),
		OrganizationID: org.ID,
		UserID:         session.UserID,
		Email:          firstNonEmpty(h.buildSessionResponse(r).User.Email, ""),
		Role:           "owner",
		Status:         "active",
		CreatedAt:      now,
	}
	if err := h.store.CreateOrganizationMemberRecord(member); err != nil {
		respondAPIError(w, http.StatusInternalServerError, "organization_member_failed", err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, map[string]interface{}{
		"organization": org,
		"member":       member,
	})
}

func (h *APIHandler) AddOrganizationMember(w http.ResponseWriter, r *http.Request) {
	session := h.sessionFromRequest(r)
	if session == nil || session.UserID == "" {
		respondAPIError(w, http.StatusUnauthorized, "auth_required", "You must be signed in to invite organization members")
		return
	}

	orgID := chi.URLParam(r, "orgId")
	if _, err := h.store.GetOrganizationRecord(orgID); err != nil {
		respondAPIError(w, http.StatusNotFound, "organization_not_found", "Organization not found")
		return
	}

	var req AddOrganizationMemberRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondAPIError(w, http.StatusBadRequest, "invalid_request", "Invalid request body")
		return
	}
	if strings.TrimSpace(req.Email) == "" {
		respondAPIError(w, http.StatusBadRequest, "invalid_email", "Invite email is required")
		return
	}

	member := &OrganizationMember{
		ID:             newID("orgmem"),
		OrganizationID: orgID,
		Email:          strings.TrimSpace(req.Email),
		Role:           firstNonEmpty(strings.TrimSpace(req.Role), "member"),
		Status:         "pending",
		CreatedAt:      time.Now(),
	}
	if err := h.store.CreateOrganizationMemberRecord(member); err != nil {
		respondAPIError(w, http.StatusInternalServerError, "organization_member_failed", err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, map[string]interface{}{"member": member})
}

func (h *APIHandler) CreateDeckShare(w http.ResponseWriter, r *http.Request) {
	session := h.sessionFromRequest(r)
	plan := h.planForRequest(r, session)
	usage := h.usageForSession(session)
	if !entitlementsForPlan(plan, usage).Features.ShareDecks {
		respondAPIError(w, http.StatusForbidden, "plan_limit_exceeded", "Deck sharing requires a Pro or Team plan")
		return
	}
	if err := validateDeckShareLimit(plan, usage); err != nil {
		respondAPIError(w, http.StatusForbidden, "plan_limit_exceeded", err.Error())
		return
	}

	deckID, err := parseIDParam(r, "deckId")
	if err != nil {
		respondAPIError(w, http.StatusBadRequest, "invalid_deck_id", "Invalid deck ID")
		return
	}
	if _, err := h.store.GetDeck(deckID); err != nil {
		respondAPIError(w, http.StatusNotFound, "deck_not_found", "Deck not found")
		return
	}

	var req ShareDeckRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil && err != io.EOF {
		respondAPIError(w, http.StatusBadRequest, "invalid_request", "Invalid request body")
		return
	}

	share := &DeckShare{
		ID:              newID("share"),
		DeckID:          deckID,
		CreatedByUserID: "",
		AccessType:      firstNonEmpty(strings.TrimSpace(req.AccessType), "read"),
		Token:           newID("share"),
		CreatedAt:       time.Now(),
	}
	if session != nil {
		share.WorkspaceID = session.WorkspaceID
		share.CreatedByUserID = session.UserID
	}
	if err := h.store.CreateDeckShareRecord(share); err != nil {
		respondAPIError(w, http.StatusInternalServerError, "deck_share_failed", err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, map[string]interface{}{
		"share": share,
		"url":   fmt.Sprintf("/shared/%s", share.Token),
	})
}

func (h *APIHandler) DeleteDeckShare(w http.ResponseWriter, r *http.Request) {
	deckID, err := parseIDParam(r, "deckId")
	if err != nil {
		respondAPIError(w, http.StatusBadRequest, "invalid_deck_id", "Invalid deck ID")
		return
	}

	if err := h.store.DeleteDeckShareByDeckID(deckID); err != nil {
		respondAPIError(w, http.StatusInternalServerError, "deck_share_delete_failed", err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func (h *APIHandler) BillingCheckout(w http.ResponseWriter, r *http.Request) {
	session := h.sessionFromRequest(r)
	if session == nil || session.UserID == "" {
		respondAPIError(w, http.StatusUnauthorized, "auth_required", "You must be signed in to start checkout")
		return
	}

	var req billingCheckoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondAPIError(w, http.StatusBadRequest, "invalid_request", "Invalid checkout request")
		return
	}

	if strings.TrimSpace(os.Getenv("VUTADEX_STRIPE_SECRET_KEY")) == "" {
		respondAPIError(w, http.StatusNotImplemented, "billing_not_configured", "Stripe checkout is not configured")
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{
		"message": "Stripe checkout is configured but the hosted checkout handoff still needs provider wiring.",
		"plan":    string(parsePlan(string(req.Plan))),
	})
}

func (h *APIHandler) BillingPortal(w http.ResponseWriter, r *http.Request) {
	session := h.sessionFromRequest(r)
	if session == nil || session.UserID == "" {
		respondAPIError(w, http.StatusUnauthorized, "auth_required", "You must be signed in to open billing portal")
		return
	}
	if strings.TrimSpace(os.Getenv("VUTADEX_STRIPE_SECRET_KEY")) == "" {
		respondAPIError(w, http.StatusNotImplemented, "billing_not_configured", "Stripe billing portal is not configured")
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{
		"message": "Stripe billing portal is configured but provider handoff still needs wiring.",
	})
}

func (h *APIHandler) BillingWebhook(w http.ResponseWriter, r *http.Request) {
	expectedSecret := strings.TrimSpace(os.Getenv("VUTADEX_BILLING_WEBHOOK_SECRET"))
	if expectedSecret != "" && strings.TrimSpace(r.Header.Get("X-Vutadex-Billing-Secret")) != expectedSecret {
		respondAPIError(w, http.StatusUnauthorized, "webhook_unauthorized", "Invalid billing webhook secret")
		return
	}

	var req billingWebhookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondAPIError(w, http.StatusBadRequest, "invalid_request", "Invalid billing webhook payload")
		return
	}

	now := time.Now()
	subscription := &Subscription{
		ID:                     newID("sub"),
		WorkspaceID:            req.WorkspaceID,
		OrganizationID:         req.OrganizationID,
		Plan:                   parsePlan(string(req.Plan)),
		Status:                 firstNonEmpty(strings.TrimSpace(req.Status), "active"),
		Provider:               firstNonEmpty(strings.TrimSpace(req.Provider), "stripe"),
		ProviderCustomerID:     req.ProviderCustomerID,
		ProviderSubscriptionID: req.ProviderSubID,
		CreatedAt:              now,
		UpdatedAt:              now,
	}
	if req.CurrentPeriodEnd != "" {
		if parsed, err := time.Parse(time.RFC3339, req.CurrentPeriodEnd); err == nil {
			subscription.CurrentPeriodEnd = parsed
		}
	}

	if err := h.store.UpsertSubscription(subscription); err != nil {
		respondAPIError(w, http.StatusInternalServerError, "billing_subscription_failed", err.Error())
		return
	}

	event := &SubscriptionEvent{
		ID:              newID("subevt"),
		SubscriptionID:  subscription.ID,
		EventType:       "billing.webhook",
		ProviderEventID: req.ProviderEventID,
		Payload:         string(mustJSON(req)),
		CreatedAt:       now,
	}
	if err := h.store.CreateSubscriptionEvent(event); err != nil {
		respondAPIError(w, http.StatusInternalServerError, "billing_event_failed", err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func mustJSON(value interface{}) []byte {
	data, _ := json.Marshal(value)
	return data
}

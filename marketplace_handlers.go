package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
)

func validMarketplacePriceMode(raw string) bool {
	switch strings.TrimSpace(strings.ToLower(raw)) {
	case "free", "premium":
		return true
	default:
		return false
	}
}

func (h *APIHandler) requireMarketplacePublishEntitlement(w http.ResponseWriter, r *http.Request) (*SessionRecord, bool) {
	session := h.sessionFromRequest(r)
	plan := h.planForRequest(r, session)
	usage := h.usageForSession(session)
	if !entitlementsForPlan(plan, usage).Features.MarketplacePublish {
		respondAPIError(w, http.StatusForbidden, "marketplace_publish_not_available", "Marketplace publishing is reserved for Pro, Team, and Enterprise workspaces.")
		return nil, false
	}
	return session, true
}

func (h *APIHandler) uniqueMarketplaceSlug(raw, title, excludeID string) (string, error) {
	base := slugify(firstNonEmpty(strings.TrimSpace(raw), strings.TrimSpace(title)))
	if base == "workspace" {
		base = "listing"
	}
	candidate := base
	for attempt := 0; attempt < 100; attempt++ {
		exists, err := h.store.MarketplaceListingSlugExists(candidate, excludeID)
		if err != nil {
			return "", err
		}
		if !exists {
			return candidate, nil
		}
		candidate = fmt.Sprintf("%s-%d", base, attempt+2)
	}
	return "", fmt.Errorf("unable to allocate unique marketplace slug")
}

func (h *APIHandler) loadMarketplaceListingForEdit(r *http.Request, ref string) (*MarketplaceListing, *User, *SessionRecord, error) {
	session := h.sessionFromRequest(r)
	if session == nil || session.UserID == "" {
		return nil, nil, nil, sql.ErrNoRows
	}
	user, err := h.store.GetUserByID(session.UserID)
	if err != nil {
		return nil, nil, session, err
	}
	listing, err := h.store.resolveMarketplaceListing(ref)
	if err != nil {
		return nil, user, session, err
	}
	return listing, user, session, nil
}

func (h *APIHandler) installMarketplaceVersion(listing *MarketplaceListing, version *MarketplaceListingVersion, destinationWorkspace *Workspace, user *User) (*MarketplaceInstall, error) {
	sourceDeck, err := h.store.GetDeck(version.SourceDeckID)
	if err != nil {
		return nil, err
	}
	installedName := fmt.Sprintf("%s (Marketplace v%d)", sourceDeck.Name, version.VersionNumber)
	newDeck, err := h.store.CopyDeckToCollection(version.SourceDeckID, destinationWorkspace.CollectionID, installedName)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	install := &MarketplaceInstall{
		ID:                  newID("mki"),
		ListingID:           listing.ID,
		WorkspaceID:         destinationWorkspace.ID,
		InstalledByUserID:   user.ID,
		InstalledDeckID:     newDeck.ID,
		InstalledDeckName:   newDeck.Name,
		SourceVersionNumber: version.VersionNumber,
		Status:              "active",
		CreatedAt:           now,
		UpdatedAt:           now,
	}
	if err := h.store.CreateMarketplaceInstall(install); err != nil {
		return nil, err
	}
	if err := h.reloadCollectionSnapshot(destinationWorkspace.CollectionID); err != nil {
		return nil, err
	}
	return h.store.GetMarketplaceInstall(install.ID)
}

func (h *APIHandler) ListMarketplaceListings(w http.ResponseWriter, r *http.Request) {
	session := h.sessionFromRequest(r)
	scope := strings.TrimSpace(r.URL.Query().Get("scope"))
	workspaceID := ""
	userID := ""
	if session != nil {
		workspaceID = session.WorkspaceID
		userID = session.UserID
	}
	listings, err := h.store.ListMarketplaceListings(scope, userID, workspaceID)
	if err != nil {
		respondAPIError(w, http.StatusInternalServerError, "marketplace_list_failed", err.Error())
		return
	}
	respondJSON(w, http.StatusOK, listings)
}

func (h *APIHandler) CreateMarketplaceListing(w http.ResponseWriter, r *http.Request) {
	session, ok := h.requireMarketplacePublishEntitlement(w, r)
	if !ok {
		return
	}

	var req CreateMarketplaceListingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondAPIError(w, http.StatusBadRequest, "invalid_request", "Invalid request body")
		return
	}
	if req.DeckID == 0 || strings.TrimSpace(req.Title) == "" {
		respondAPIError(w, http.StatusBadRequest, "invalid_marketplace_listing", "Deck and title are required.")
		return
	}

	workspace, err := h.store.GetWorkspaceRecord(session.WorkspaceID)
	if err != nil {
		respondAPIError(w, http.StatusBadRequest, "workspace_not_found", "Workspace not found.")
		return
	}
	sourceCollectionID, err := h.store.GetDeckCollectionID(req.DeckID)
	if err != nil {
		respondAPIError(w, http.StatusBadRequest, "source_deck_not_found", "Source deck not found.")
		return
	}
	if sourceCollectionID != workspace.CollectionID {
		respondAPIError(w, http.StatusBadRequest, "invalid_source_deck", "Source deck must belong to the current workspace.")
		return
	}

	slug, err := h.uniqueMarketplaceSlug(req.Slug, req.Title, "")
	if err != nil {
		respondAPIError(w, http.StatusInternalServerError, "marketplace_slug_failed", err.Error())
		return
	}
	priceMode := strings.ToLower(strings.TrimSpace(req.PriceMode))
	if priceMode == "" {
		priceMode = "free"
	}
	if !validMarketplacePriceMode(priceMode) {
		respondAPIError(w, http.StatusBadRequest, "invalid_marketplace_price_mode", "Price mode must be free or premium.")
		return
	}
	currency := strings.ToUpper(strings.TrimSpace(req.Currency))
	if currency == "" {
		currency = "USD"
	}

	now := time.Now()
	listing := &MarketplaceListing{
		ID:            newID("mkt"),
		WorkspaceID:   workspace.ID,
		DeckID:        req.DeckID,
		Slug:          slug,
		Title:         sanitizeHTML(strings.TrimSpace(req.Title)),
		Summary:       sanitizeHTML(strings.TrimSpace(req.Summary)),
		Description:   sanitizeHTML(strings.TrimSpace(req.Description)),
		Category:      sanitizeHTML(strings.TrimSpace(req.Category)),
		Tags:          sanitizeMarketplaceTags(req.Tags),
		CoverImageURL: strings.TrimSpace(req.CoverImageURL),
		CreatorUserID: session.UserID,
		PriceMode:     priceMode,
		PriceCents:    req.PriceCents,
		Currency:      currency,
		Status:        "draft",
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	if err := h.store.CreateMarketplaceListing(listing); err != nil {
		respondAPIError(w, http.StatusInternalServerError, "marketplace_create_failed", err.Error())
		return
	}
	detail, err := h.store.BuildMarketplaceListingDetail(listing.ID, session.UserID, session.WorkspaceID)
	if err != nil {
		respondAPIError(w, http.StatusInternalServerError, "marketplace_detail_failed", err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, detail)
}

func (h *APIHandler) GetMarketplaceListing(w http.ResponseWriter, r *http.Request) {
	ref := chi.URLParam(r, "ref")
	session := h.sessionFromRequest(r)
	userID := ""
	workspaceID := ""
	if session != nil {
		userID = session.UserID
		workspaceID = session.WorkspaceID
	}
	detail, err := h.store.BuildMarketplaceListingDetail(ref, userID, workspaceID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondAPIError(w, http.StatusNotFound, "marketplace_listing_not_found", "Marketplace listing not found.")
			return
		}
		respondAPIError(w, http.StatusInternalServerError, "marketplace_detail_failed", err.Error())
		return
	}
	respondJSON(w, http.StatusOK, detail)
}

func (h *APIHandler) UpdateMarketplaceListing(w http.ResponseWriter, r *http.Request) {
	session, ok := h.requireMarketplacePublishEntitlement(w, r)
	if !ok {
		return
	}

	ref := chi.URLParam(r, "ref")
	listing, _, _, err := h.loadMarketplaceListingForEdit(r, ref)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondAPIError(w, http.StatusNotFound, "marketplace_listing_not_found", "Marketplace listing not found.")
			return
		}
		respondAPIError(w, http.StatusInternalServerError, "marketplace_access_failed", err.Error())
		return
	}
	if listing.CreatorUserID != session.UserID || listing.WorkspaceID != session.WorkspaceID {
		respondAPIError(w, http.StatusForbidden, "marketplace_forbidden", "You can only edit listings in your current workspace.")
		return
	}

	var req UpdateMarketplaceListingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondAPIError(w, http.StatusBadRequest, "invalid_request", "Invalid request body")
		return
	}
	if req.DeckID == 0 || strings.TrimSpace(req.Title) == "" {
		respondAPIError(w, http.StatusBadRequest, "invalid_marketplace_listing", "Deck and title are required.")
		return
	}

	workspace, err := h.store.GetWorkspaceRecord(session.WorkspaceID)
	if err != nil {
		respondAPIError(w, http.StatusBadRequest, "workspace_not_found", "Workspace not found.")
		return
	}
	sourceCollectionID, err := h.store.GetDeckCollectionID(req.DeckID)
	if err != nil {
		respondAPIError(w, http.StatusBadRequest, "source_deck_not_found", "Source deck not found.")
		return
	}
	if sourceCollectionID != workspace.CollectionID {
		respondAPIError(w, http.StatusBadRequest, "invalid_source_deck", "Source deck must belong to the current workspace.")
		return
	}

	slug, err := h.uniqueMarketplaceSlug(req.Slug, req.Title, listing.ID)
	if err != nil {
		respondAPIError(w, http.StatusInternalServerError, "marketplace_slug_failed", err.Error())
		return
	}
	priceMode := strings.ToLower(strings.TrimSpace(req.PriceMode))
	if priceMode == "" {
		priceMode = "free"
	}
	if !validMarketplacePriceMode(priceMode) {
		respondAPIError(w, http.StatusBadRequest, "invalid_marketplace_price_mode", "Price mode must be free or premium.")
		return
	}
	currency := strings.ToUpper(strings.TrimSpace(req.Currency))
	if currency == "" {
		currency = "USD"
	}

	listing.DeckID = req.DeckID
	listing.Slug = slug
	listing.Title = sanitizeHTML(strings.TrimSpace(req.Title))
	listing.Summary = sanitizeHTML(strings.TrimSpace(req.Summary))
	listing.Description = sanitizeHTML(strings.TrimSpace(req.Description))
	listing.Category = sanitizeHTML(strings.TrimSpace(req.Category))
	listing.Tags = sanitizeMarketplaceTags(req.Tags)
	listing.CoverImageURL = strings.TrimSpace(req.CoverImageURL)
	listing.PriceMode = priceMode
	listing.PriceCents = req.PriceCents
	listing.Currency = currency
	listing.UpdatedAt = time.Now()
	if err := h.store.UpdateMarketplaceListing(listing); err != nil {
		respondAPIError(w, http.StatusInternalServerError, "marketplace_update_failed", err.Error())
		return
	}
	detail, err := h.store.BuildMarketplaceListingDetail(listing.ID, session.UserID, session.WorkspaceID)
	if err != nil {
		respondAPIError(w, http.StatusInternalServerError, "marketplace_detail_failed", err.Error())
		return
	}
	respondJSON(w, http.StatusOK, detail)
}

func (h *APIHandler) DeleteMarketplaceListing(w http.ResponseWriter, r *http.Request) {
	session, ok := h.requireMarketplacePublishEntitlement(w, r)
	if !ok {
		return
	}

	ref := chi.URLParam(r, "ref")
	listing, _, _, err := h.loadMarketplaceListingForEdit(r, ref)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondAPIError(w, http.StatusNotFound, "marketplace_listing_not_found", "Marketplace listing not found.")
			return
		}
		respondAPIError(w, http.StatusInternalServerError, "marketplace_access_failed", err.Error())
		return
	}
	if listing.CreatorUserID != session.UserID || listing.WorkspaceID != session.WorkspaceID {
		respondAPIError(w, http.StatusForbidden, "marketplace_forbidden", "You can only delete listings in your current workspace.")
		return
	}
	installCount, err := h.store.CountMarketplaceInstalls(listing.ID)
	if err != nil {
		respondAPIError(w, http.StatusInternalServerError, "marketplace_delete_failed", err.Error())
		return
	}
	if installCount > 0 {
		respondAPIError(w, http.StatusConflict, "marketplace_listing_has_installs", "Remove active installs before deleting this listing.")
		return
	}
	if err := h.store.DeleteMarketplaceListing(listing.ID); err != nil {
		respondAPIError(w, http.StatusInternalServerError, "marketplace_delete_failed", err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *APIHandler) PublishMarketplaceListing(w http.ResponseWriter, r *http.Request) {
	session, ok := h.requireMarketplacePublishEntitlement(w, r)
	if !ok {
		return
	}

	ref := chi.URLParam(r, "ref")
	listing, _, _, err := h.loadMarketplaceListingForEdit(r, ref)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondAPIError(w, http.StatusNotFound, "marketplace_listing_not_found", "Marketplace listing not found.")
			return
		}
		respondAPIError(w, http.StatusInternalServerError, "marketplace_access_failed", err.Error())
		return
	}
	if listing.CreatorUserID != session.UserID || listing.WorkspaceID != session.WorkspaceID {
		respondAPIError(w, http.StatusForbidden, "marketplace_forbidden", "You can only publish listings in your current workspace.")
		return
	}

	var req PublishMarketplaceListingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil && !errors.Is(err, http.ErrBodyNotAllowed) {
		respondAPIError(w, http.StatusBadRequest, "invalid_request", "Invalid request body")
		return
	}

	nextVersion := 1
	if latestVersion, err := h.store.GetLatestMarketplaceListingVersion(listing.ID); err == nil {
		nextVersion = latestVersion.VersionNumber + 1
	}
	noteCount, cardCount, err := h.store.GetDeckContentSummary(listing.DeckID)
	if err != nil {
		respondAPIError(w, http.StatusInternalServerError, "marketplace_publish_failed", err.Error())
		return
	}
	version := &MarketplaceListingVersion{
		ID:                newID("mkv"),
		ListingID:         listing.ID,
		VersionNumber:     nextVersion,
		SourceDeckID:      listing.DeckID,
		PublishedByUserID: session.UserID,
		ChangeSummary:     strings.TrimSpace(req.ChangeSummary),
		NoteCount:         noteCount,
		CardCount:         cardCount,
		CreatedAt:         time.Now(),
	}
	if err := h.store.CreateMarketplaceListingVersion(version); err != nil {
		respondAPIError(w, http.StatusInternalServerError, "marketplace_publish_failed", err.Error())
		return
	}
	listing.Status = "published"
	listing.UpdatedAt = time.Now()
	if err := h.store.UpdateMarketplaceListing(listing); err != nil {
		respondAPIError(w, http.StatusInternalServerError, "marketplace_publish_failed", err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, version)
}

func (h *APIHandler) InstallMarketplaceListing(w http.ResponseWriter, r *http.Request) {
	ref := chi.URLParam(r, "ref")
	session := h.sessionFromRequest(r)
	user, err := h.store.GetUserByID(session.UserID)
	if err != nil {
		respondAPIError(w, http.StatusInternalServerError, "user_not_found", err.Error())
		return
	}
	listing, err := h.store.resolveMarketplaceListing(ref)
	if err != nil || listing.Status != "published" {
		respondAPIError(w, http.StatusNotFound, "marketplace_listing_not_found", "Marketplace listing not found.")
		return
	}
	if listing.PriceMode != "free" {
		respondAPIError(w, http.StatusConflict, "marketplace_purchase_required", "Premium marketplace checkout arrives in Phase 4. Only free listings are installable right now.")
		return
	}

	var req InstallMarketplaceListingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondAPIError(w, http.StatusBadRequest, "invalid_request", "Invalid request body")
		return
	}
	if req.DestinationWorkspaceID == "" {
		req.DestinationWorkspaceID = session.WorkspaceID
	}
	workspace, err := h.store.GetWorkspaceForUser(user.ID, req.DestinationWorkspaceID)
	if err != nil {
		respondAPIError(w, http.StatusBadRequest, "invalid_destination_workspace", "Destination workspace not found.")
		return
	}
	latestVersion, err := h.store.GetLatestMarketplaceListingVersion(listing.ID)
	if err != nil {
		respondAPIError(w, http.StatusConflict, "marketplace_not_published", "Publish a listing version before it can be installed.")
		return
	}
	if currentInstall, err := h.store.GetCurrentMarketplaceInstall(listing.ID, user.ID); err == nil {
		if currentInstall.SourceVersionNumber == latestVersion.VersionNumber {
			respondJSON(w, http.StatusOK, currentInstall)
			return
		}
		respondAPIError(w, http.StatusConflict, "marketplace_update_required", "Use the install update flow to move to a newer marketplace version.")
		return
	}
	install, err := h.installMarketplaceVersion(listing, latestVersion, workspace, user)
	if err != nil {
		respondAPIError(w, http.StatusInternalServerError, "marketplace_install_failed", err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, install)
}

func (h *APIHandler) UpdateMarketplaceInstall(w http.ResponseWriter, r *http.Request) {
	ref := chi.URLParam(r, "ref")
	installID := chi.URLParam(r, "installId")
	session := h.sessionFromRequest(r)
	user, err := h.store.GetUserByID(session.UserID)
	if err != nil {
		respondAPIError(w, http.StatusInternalServerError, "user_not_found", err.Error())
		return
	}
	listing, err := h.store.resolveMarketplaceListing(ref)
	if err != nil || listing.Status != "published" {
		respondAPIError(w, http.StatusNotFound, "marketplace_listing_not_found", "Marketplace listing not found.")
		return
	}
	currentInstall, err := h.store.GetMarketplaceInstall(installID)
	if err != nil || currentInstall.ListingID != listing.ID || currentInstall.InstalledByUserID != user.ID {
		respondAPIError(w, http.StatusNotFound, "marketplace_install_not_found", "Marketplace install not found.")
		return
	}
	latestVersion, err := h.store.GetLatestMarketplaceListingVersion(listing.ID)
	if err != nil {
		respondAPIError(w, http.StatusConflict, "marketplace_not_published", "No published marketplace version is available.")
		return
	}
	if latestVersion.VersionNumber <= currentInstall.SourceVersionNumber {
		respondAPIError(w, http.StatusConflict, "marketplace_already_current", "This marketplace install already uses the latest published version.")
		return
	}

	var req UpdateMarketplaceInstallRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil && !errors.Is(err, http.ErrBodyNotAllowed) {
		respondAPIError(w, http.StatusBadRequest, "invalid_request", "Invalid request body")
		return
	}
	destinationWorkspaceID := currentInstall.WorkspaceID
	if strings.TrimSpace(req.DestinationWorkspaceID) != "" {
		destinationWorkspaceID = strings.TrimSpace(req.DestinationWorkspaceID)
	}
	workspace, err := h.store.GetWorkspaceForUser(user.ID, destinationWorkspaceID)
	if err != nil {
		respondAPIError(w, http.StatusBadRequest, "invalid_destination_workspace", "Destination workspace not found.")
		return
	}

	nextInstall, err := h.installMarketplaceVersion(listing, latestVersion, workspace, user)
	if err != nil {
		respondAPIError(w, http.StatusInternalServerError, "marketplace_install_update_failed", err.Error())
		return
	}
	currentInstall.Status = "superseded"
	currentInstall.SupersededByInstall = nextInstall.ID
	currentInstall.UpdatedAt = time.Now()
	if err := h.store.UpdateMarketplaceInstall(currentInstall); err != nil {
		respondAPIError(w, http.StatusInternalServerError, "marketplace_install_update_failed", err.Error())
		return
	}
	respondJSON(w, http.StatusOK, nextInstall)
}

func (h *APIHandler) RemoveMarketplaceInstall(w http.ResponseWriter, r *http.Request) {
	installID := chi.URLParam(r, "installId")
	session := h.sessionFromRequest(r)
	install, err := h.store.GetMarketplaceInstall(installID)
	if err != nil || install.InstalledByUserID != session.UserID {
		respondAPIError(w, http.StatusNotFound, "marketplace_install_not_found", "Marketplace install not found.")
		return
	}
	if install.InstalledDeckID != 0 {
		if err := h.store.DeleteCopiedDeck(install.InstalledDeckID); err != nil {
			respondAPIError(w, http.StatusInternalServerError, "marketplace_install_delete_failed", err.Error())
			return
		}
	}
	install.Status = "removed"
	install.InstalledDeckID = 0
	install.InstalledDeckName = ""
	install.UpdatedAt = time.Now()
	if err := h.store.UpdateMarketplaceInstall(install); err != nil {
		respondAPIError(w, http.StatusInternalServerError, "marketplace_install_delete_failed", err.Error())
		return
	}
	if workspace, err := h.store.GetWorkspaceRecord(install.WorkspaceID); err == nil {
		if reloadErr := h.reloadCollectionSnapshot(workspace.CollectionID); reloadErr != nil {
			respondAPIError(w, http.StatusInternalServerError, "marketplace_install_delete_failed", reloadErr.Error())
			return
		}
	}
	w.WriteHeader(http.StatusNoContent)
}

package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
)

func (h *APIHandler) getMarketplaceCreatorAccountStatus(ctx context.Context, userID string) MarketplaceCreatorAccountStatusResponse {
	provider := newMarketplaceBillingProvider(h.config)
	response := MarketplaceCreatorAccountStatusResponse{
		Provider: provider.ProviderName(),
	}
	account, err := h.store.GetMarketplaceCreatorAccountByUser(userID)
	if err == nil {
		if refreshed, refreshErr := provider.RefreshCreatorAccount(ctx, account); refreshErr == nil && refreshed != nil {
			refreshed.UserID = account.UserID
			refreshed.WorkspaceID = account.WorkspaceID
			refreshed.ID = account.ID
			refreshed.CreatedAt = account.CreatedAt
			if err := h.store.UpsertMarketplaceCreatorAccount(refreshed); err == nil {
				account = refreshed
			}
		}
		response.Account = account
		response.CanSellPremium = account.ChargesEnabled && account.PayoutsEnabled && account.DetailsSubmitted
	}
	return response
}

func (h *APIHandler) GetMarketplaceCreatorAccountStatus(w http.ResponseWriter, r *http.Request) {
	session := h.sessionFromRequest(r)
	respondJSON(w, http.StatusOK, h.getMarketplaceCreatorAccountStatus(r.Context(), session.UserID))
}

func (h *APIHandler) StartMarketplaceCreatorAccount(w http.ResponseWriter, r *http.Request) {
	session := h.sessionFromRequest(r)
	plan := h.planForRequest(r, session)
	if !entitlementsForPlan(plan, h.usageForSession(session)).Features.MarketplacePublish {
		respondAPIError(w, http.StatusForbidden, "marketplace_publish_not_available", "Marketplace creator accounts are reserved for Pro, Team, and Enterprise workspaces.")
		return
	}

	user, err := h.store.GetUserByID(session.UserID)
	if err != nil {
		respondAPIError(w, http.StatusInternalServerError, "user_not_found", err.Error())
		return
	}
	workspace, err := h.store.GetWorkspaceForUser(user.ID, session.WorkspaceID)
	if err != nil {
		respondAPIError(w, http.StatusBadRequest, "invalid_workspace", "Current workspace not found.")
		return
	}

	var existing *MarketplaceCreatorAccount
	if account, err := h.store.GetMarketplaceCreatorAccountByUser(user.ID); err == nil {
		existing = account
	} else if !errors.Is(err, sql.ErrNoRows) {
		respondAPIError(w, http.StatusInternalServerError, "marketplace_creator_account_failed", err.Error())
		return
	}

	provider := newMarketplaceBillingProvider(h.config)
	account, err := provider.StartCreatorOnboarding(r.Context(), user, workspace, existing)
	if err != nil {
		respondAPIError(w, http.StatusConflict, "marketplace_creator_account_failed", err.Error())
		return
	}
	if err := h.store.UpsertMarketplaceCreatorAccount(account); err != nil {
		respondAPIError(w, http.StatusInternalServerError, "marketplace_creator_account_failed", err.Error())
		return
	}
	respondJSON(w, http.StatusOK, MarketplaceCreatorAccountStatusResponse{
		Account:        account,
		Provider:       provider.ProviderName(),
		CanSellPremium: account.ChargesEnabled && account.PayoutsEnabled && account.DetailsSubmitted,
	})
}

func (h *APIHandler) activeMarketplaceLicense(listingID, userID string) (*MarketplaceLicense, error) {
	license, err := h.store.GetMarketplaceLicense(listingID, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	if license.Status != "active" {
		return nil, nil
	}
	return license, nil
}

func (h *APIHandler) canAccessPremiumMarketplaceListing(listing *MarketplaceListing, userID string) (*MarketplaceLicense, bool, error) {
	if listing.PriceMode == "free" || listing.CreatorUserID == userID {
		return nil, true, nil
	}
	license, err := h.activeMarketplaceLicense(listing.ID, userID)
	if err != nil {
		return nil, false, err
	}
	return license, license != nil, nil
}

func (h *APIHandler) ensureMarketplaceCreatorReady(ctx context.Context, listing *MarketplaceListing) (*MarketplaceCreatorAccount, error) {
	account, err := h.store.GetMarketplaceCreatorAccountByUser(listing.CreatorUserID)
	if err != nil {
		return nil, err
	}
	provider := newMarketplaceBillingProvider(h.config)
	if refreshed, refreshErr := provider.RefreshCreatorAccount(ctx, account); refreshErr == nil && refreshed != nil {
		refreshed.UserID = account.UserID
		refreshed.WorkspaceID = account.WorkspaceID
		refreshed.ID = account.ID
		refreshed.CreatedAt = account.CreatedAt
		if err := h.store.UpsertMarketplaceCreatorAccount(refreshed); err == nil {
			account = refreshed
		}
	}
	if !account.DetailsSubmitted || !account.ChargesEnabled || !account.PayoutsEnabled {
		return nil, errors.New("creator onboarding is incomplete")
	}
	return account, nil
}

func (h *APIHandler) completeMarketplaceOrder(order *MarketplaceOrder) (*MarketplaceLicense, error) {
	if order.Status == "paid" {
		license, err := h.store.GetMarketplaceLicense(order.ListingID, order.BuyerUserID)
		if err == nil {
			return license, nil
		}
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
	}

	now := time.Now()
	order.Status = "paid"
	order.CompletedAt = now
	order.UpdatedAt = now
	if err := h.store.UpdateMarketplaceOrder(order); err != nil {
		return nil, err
	}

	license := &MarketplaceLicense{
		ID:                   newID("mlic"),
		ListingID:            order.ListingID,
		BuyerUserID:          order.BuyerUserID,
		OrderID:              order.ID,
		Status:               "active",
		GrantedVersionNumber: order.ListingVersionNumber,
		CreatedAt:            now,
		UpdatedAt:            now,
	}
	if existing, err := h.store.GetMarketplaceLicense(order.ListingID, order.BuyerUserID); err == nil {
		license.ID = existing.ID
		license.CreatedAt = existing.CreatedAt
	}
	if err := h.store.UpsertMarketplaceLicense(license); err != nil {
		return nil, err
	}

	payout := &MarketplacePayout{
		ID:               newID("mpo"),
		OrderID:          order.ID,
		CreatorUserID:    order.CreatorUserID,
		CreatorAccountID: order.CreatorAccountID,
		Provider:         order.Provider,
		Status:           "pending",
		AmountCents:      order.CreatorAmountCents,
		Currency:         order.Currency,
		PlatformFeeCents: order.PlatformFeeCents,
		CreatedAt:        now,
		UpdatedAt:        now,
	}
	if existing, err := h.store.GetMarketplacePayoutByOrder(order.ID); err == nil {
		payout.ID = existing.ID
		payout.ProviderTransferID = existing.ProviderTransferID
		payout.CreatedAt = existing.CreatedAt
	}
	if err := h.store.UpsertMarketplacePayout(payout); err != nil {
		return nil, err
	}

	return license, nil
}

func (h *APIHandler) failMarketplaceOrder(order *MarketplaceOrder) error {
	if order == nil {
		return nil
	}
	if order.Status == "paid" || order.Status == "failed" || order.Status == "canceled" {
		return nil
	}
	order.Status = "failed"
	order.UpdatedAt = time.Now()
	return h.store.UpdateMarketplaceOrder(order)
}

func (h *APIHandler) SyncMarketplaceCheckoutSession(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionId")
	if strings.TrimSpace(sessionID) == "" {
		respondAPIError(w, http.StatusBadRequest, "invalid_checkout_session", "Checkout session is required.")
		return
	}

	session := h.sessionFromRequest(r)
	order, err := h.store.GetMarketplaceOrderByCheckoutSession("stripe", sessionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondAPIError(w, http.StatusNotFound, "marketplace_order_not_found", "Marketplace order not found.")
			return
		}
		respondAPIError(w, http.StatusInternalServerError, "marketplace_order_lookup_failed", err.Error())
		return
	}
	if order.BuyerUserID != session.UserID {
		respondAPIError(w, http.StatusForbidden, "marketplace_forbidden", "You can only sync your own marketplace checkout.")
		return
	}

	license, response, err := h.syncMarketplaceCheckoutSession(r.Context(), order)
	if err != nil {
		respondAPIError(w, http.StatusConflict, "marketplace_checkout_sync_failed", err.Error())
		return
	}
	response.License = license
	respondJSON(w, http.StatusOK, response)
}

func (h *APIHandler) syncMarketplaceCheckoutSession(ctx context.Context, order *MarketplaceOrder) (*MarketplaceLicense, MarketplaceCheckoutResponse, error) {
	response := MarketplaceCheckoutResponse{
		Provider: firstNonEmpty(order.Provider, "stripe"),
		Order:    *order,
	}

	if order.Status == "paid" {
		license, err := h.store.GetMarketplaceLicense(order.ListingID, order.BuyerUserID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, response, nil
			}
			return nil, response, err
		}
		response.Completed = true
		return license, response, nil
	}

	creatorAccount, err := h.store.GetMarketplaceCreatorAccount(order.CreatorAccountID)
	if err != nil {
		return nil, response, err
	}
	provider := newMarketplaceBillingProvider(h.config)
	state, err := provider.GetCheckoutSession(ctx, order, creatorAccount)
	if err != nil {
		return nil, response, err
	}
	if state.ProviderPaymentIntentID != "" {
		order.ProviderPaymentIntentID = state.ProviderPaymentIntentID
	}
	response.Provider = state.Provider
	response.CheckoutURL = state.ProviderCheckoutURL
	response.Order = *order

	if state.Completed {
		license, err := h.completeMarketplaceOrder(order)
		if err != nil {
			return nil, response, err
		}
		response.Completed = true
		response.Order = *order
		return license, response, nil
	}
	if strings.EqualFold(state.PaymentStatus, "unpaid") && strings.EqualFold(state.Status, "expired") {
		if err := h.failMarketplaceOrder(order); err != nil {
			return nil, response, err
		}
		response.Order = *order
	}
	return nil, response, nil
}

func (h *APIHandler) MarketplaceWebhook(w http.ResponseWriter, r *http.Request) {
	provider := newMarketplaceBillingProvider(h.config)
	if provider.ProviderName() != "stripe" {
		respondAPIError(w, http.StatusNotImplemented, "marketplace_billing_not_configured", "Marketplace billing webhooks are not configured.")
		return
	}

	payload, err := io.ReadAll(r.Body)
	if err != nil {
		respondAPIError(w, http.StatusBadRequest, "invalid_request", "Unable to read webhook payload.")
		return
	}
	if err := verifyStripeWebhookSignature(payload, r.Header.Get("Stripe-Signature"), h.config.Stripe.WebhookSecret, time.Now()); err != nil {
		respondAPIError(w, http.StatusUnauthorized, "webhook_unauthorized", err.Error())
		return
	}

	var event stripeEventEnvelope
	if err := json.Unmarshal(payload, &event); err != nil {
		respondAPIError(w, http.StatusBadRequest, "invalid_request", "Invalid webhook payload.")
		return
	}

	switch event.Type {
	case "checkout.session.completed", "checkout.session.async_payment_succeeded":
		var session stripeCheckoutSessionResponse
		if err := json.Unmarshal(event.Data.Object, &session); err != nil {
			respondAPIError(w, http.StatusBadRequest, "invalid_request", "Invalid checkout session payload.")
			return
		}
		order, err := h.store.GetMarketplaceOrderByCheckoutSession("stripe", session.ID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				respondJSON(w, http.StatusOK, map[string]bool{"ok": true})
				return
			}
			respondAPIError(w, http.StatusInternalServerError, "marketplace_order_lookup_failed", err.Error())
			return
		}
		if session.PaymentIntent != "" {
			order.ProviderPaymentIntentID = session.PaymentIntent
		}
		if _, err := h.completeMarketplaceOrder(order); err != nil {
			respondAPIError(w, http.StatusInternalServerError, "marketplace_checkout_complete_failed", err.Error())
			return
		}
	case "checkout.session.expired", "checkout.session.async_payment_failed":
		var session stripeCheckoutSessionResponse
		if err := json.Unmarshal(event.Data.Object, &session); err != nil {
			respondAPIError(w, http.StatusBadRequest, "invalid_request", "Invalid checkout session payload.")
			return
		}
		order, err := h.store.GetMarketplaceOrderByCheckoutSession("stripe", session.ID)
		if err == nil {
			if session.PaymentIntent != "" {
				order.ProviderPaymentIntentID = session.PaymentIntent
			}
			if err := h.failMarketplaceOrder(order); err != nil {
				respondAPIError(w, http.StatusInternalServerError, "marketplace_order_update_failed", err.Error())
				return
			}
		} else if !errors.Is(err, sql.ErrNoRows) {
			respondAPIError(w, http.StatusInternalServerError, "marketplace_order_lookup_failed", err.Error())
			return
		}
	case "account.updated":
		var accountObject stripeAccountWebhookObject
		if err := json.Unmarshal(event.Data.Object, &accountObject); err != nil {
			respondAPIError(w, http.StatusBadRequest, "invalid_request", "Invalid account payload.")
			return
		}
		creatorAccount, err := h.store.GetMarketplaceCreatorAccountByProviderAccount("stripe", firstNonEmpty(accountObject.ID, event.Account))
		if err == nil {
			updated := *creatorAccount
			updated.UpdatedAt = time.Now()
			applyStripeAccountState(&updated, &stripeAccountResponse{
				ID:               creatorAccount.ProviderAccountID,
				DetailsSubmitted: accountObject.DetailsSubmitted,
				ChargesEnabled:   accountObject.ChargesEnabled,
				PayoutsEnabled:   accountObject.PayoutsEnabled,
			})
			if err := h.store.UpsertMarketplaceCreatorAccount(&updated); err != nil {
				respondAPIError(w, http.StatusInternalServerError, "marketplace_creator_account_failed", err.Error())
				return
			}
		} else if !errors.Is(err, sql.ErrNoRows) {
			respondAPIError(w, http.StatusInternalServerError, "marketplace_creator_account_failed", err.Error())
			return
		}
	}

	respondJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func (h *APIHandler) CheckoutMarketplaceListing(w http.ResponseWriter, r *http.Request) {
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
	if listing.PriceMode != "premium" {
		respondAPIError(w, http.StatusConflict, "marketplace_listing_is_free", "This listing is free. Install it directly instead of starting checkout.")
		return
	}
	if listing.CreatorUserID == session.UserID {
		respondAPIError(w, http.StatusConflict, "marketplace_creator_owns_listing", "Listing creators do not need to purchase their own premium listings.")
		return
	}

	license, err := h.activeMarketplaceLicense(listing.ID, user.ID)
	if err != nil {
		respondAPIError(w, http.StatusInternalServerError, "marketplace_license_lookup_failed", err.Error())
		return
	}
	if license != nil {
		order, err := h.store.GetMarketplaceOrder(license.OrderID)
		if err != nil {
			respondAPIError(w, http.StatusInternalServerError, "marketplace_order_lookup_failed", err.Error())
			return
		}
		respondJSON(w, http.StatusOK, MarketplaceCheckoutResponse{
			Provider:  order.Provider,
			Completed: true,
			Order:     *order,
			License:   license,
		})
		return
	}

	creatorAccount, err := h.ensureMarketplaceCreatorReady(r.Context(), listing)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || strings.Contains(strings.ToLower(err.Error()), "incomplete") {
			respondAPIError(w, http.StatusConflict, "marketplace_creator_onboarding_required", "The listing creator still needs to complete payout onboarding before premium checkout is available.")
			return
		}
		respondAPIError(w, http.StatusInternalServerError, "marketplace_creator_account_failed", err.Error())
		return
	}
	latestVersion, err := h.store.GetLatestMarketplaceListingVersion(listing.ID)
	if err != nil {
		respondAPIError(w, http.StatusConflict, "marketplace_not_published", "No published marketplace version is available.")
		return
	}

	now := time.Now()
	order := &MarketplaceOrder{
		ID:                        newID("mord"),
		ListingID:                 listing.ID,
		ListingVersionNumber:      latestVersion.VersionNumber,
		BuyerUserID:               user.ID,
		BuyerWorkspaceID:          session.WorkspaceID,
		CreatorUserID:             listing.CreatorUserID,
		CreatorAccountID:          creatorAccount.ID,
		Provider:                  "",
		ProviderCheckoutSessionID: "",
		Status:                    "pending",
		AmountCents:               listing.PriceCents,
		Currency:                  firstNonEmpty(strings.TrimSpace(listing.Currency), "USD"),
		PlatformFeeCents:          marketplacePlatformFeeCents(listing.PriceCents, h.config.Stripe.PlatformFeeBasisPts),
		CreatedAt:                 now,
		UpdatedAt:                 now,
	}
	order.CreatorAmountCents = order.AmountCents - order.PlatformFeeCents

	provider := newMarketplaceBillingProvider(h.config)
	response, err := provider.CreateCheckoutSession(r.Context(), listing, order, user, creatorAccount)
	if err != nil {
		respondAPIError(w, http.StatusConflict, "marketplace_checkout_failed", err.Error())
		return
	}
	if order.Provider == "" {
		order.Provider = response.Provider
	}
	if order.ProviderCheckoutSessionID == "" {
		order.ProviderCheckoutSessionID = response.Order.ProviderCheckoutSessionID
	}
	if order.ProviderPaymentIntentID == "" {
		order.ProviderPaymentIntentID = response.Order.ProviderPaymentIntentID
	}
	if err := h.store.CreateMarketplaceOrder(order); err != nil {
		respondAPIError(w, http.StatusInternalServerError, "marketplace_order_create_failed", err.Error())
		return
	}

	response.Order = *order
	if response.Completed {
		license, err := h.completeMarketplaceOrder(order)
		if err != nil {
			respondAPIError(w, http.StatusInternalServerError, "marketplace_checkout_complete_failed", err.Error())
			return
		}
		response.Order = *order
		response.License = license
	}
	respondJSON(w, http.StatusOK, response)
}

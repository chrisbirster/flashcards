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

func (h *APIHandler) handleBillingCheckout(w http.ResponseWriter, r *http.Request) {
	session := h.sessionFromRequest(r)
	if session == nil || session.UserID == "" {
		respondAPIError(w, http.StatusUnauthorized, "auth_required", "You must be signed in to start checkout.")
		return
	}

	var req billingCheckoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondAPIError(w, http.StatusBadRequest, "invalid_request", "Invalid checkout request.")
		return
	}
	if req.Plan != PlanPro && req.Plan != PlanTeam {
		if req.Plan == PlanEnterprise {
			respondAPIError(w, http.StatusConflict, "enterprise_sales_required", "Enterprise plans are handled manually. Contact sales to continue.")
			return
		}
		respondAPIError(w, http.StatusBadRequest, "invalid_plan", "Checkout is only available for Pro and Team plans.")
		return
	}

	user, err := h.store.GetUserByID(session.UserID)
	if err != nil {
		respondAPIError(w, http.StatusInternalServerError, "user_lookup_failed", err.Error())
		return
	}
	workspace, err := h.workspaceForSession(session)
	if err != nil || workspace == nil {
		respondAPIError(w, http.StatusInternalServerError, "workspace_not_found", "Workspace not found.")
		return
	}
	org, orgMember, err := h.billingWorkspaceContext(session, workspace)
	if err != nil {
		respondAPIError(w, http.StatusInternalServerError, "workspace_context_failed", err.Error())
		return
	}
	if !h.canManageBillingPlan(session, workspace, orgMember) {
		respondAPIError(w, http.StatusForbidden, "plan_forbidden", "You do not have permission to start billing checkout for this workspace.")
		return
	}

	currentPlan := h.planForRequest(r, session)
	if currentPlan != PlanFree && currentPlan != PlanGuest {
		respondAPIError(w, http.StatusConflict, "billing_portal_required", "This workspace already has a paid subscription. Use the billing portal to manage it.")
		return
	}

	if _, disabled := h.subscriptionBilling.(*disabledSubscriptionBillingProvider); disabled {
		if h.config.IsDevelopment() {
			orgForManual := org
			if req.Plan == PlanTeam {
				orgForManual, err = h.ensureTeamOrganizationForBilling(workspace, user, session)
				if err != nil {
					respondAPIError(w, http.StatusInternalServerError, "organization_create_failed", err.Error())
					return
				}
			}
			if orgForManual != nil {
				workspace.OrganizationID = orgForManual.ID
			}
			response, err := h.applyWorkspacePlanChange(r, session, workspace, req.Plan, user.Onboarding)
			if err != nil {
				respondAPIError(w, http.StatusInternalServerError, "billing_fallback_failed", err.Error())
				return
			}
			respondJSON(w, http.StatusOK, BillingCheckoutResponse{
				Provider:     subscriptionProviderManual,
				Plan:         req.Plan,
				Completed:    true,
				Message:      "Applied locally in development mode.",
				Session:      response,
				Subscription: response.Subscription,
			})
			return
		}

		respondAPIError(w, http.StatusNotImplemented, "billing_not_configured", "Stripe subscription billing is not configured.")
		return
	}

	existingSubscription := h.organizationSubscription(workspace.OrganizationID, workspace.ID)
	quantity := 1
	if req.Plan == PlanTeam {
		activeMembers := 0
		if org != nil {
			activeMembers, err = h.store.CountActiveOrganizationMembers(org.ID)
			if err != nil {
				respondAPIError(w, http.StatusInternalServerError, "organization_members_failed", err.Error())
				return
			}
		}
		quantity = teamSeatQuantity(activeMembers)
	}

	result, err := h.subscriptionBilling.CreateCheckoutSession(r.Context(), subscriptionCheckoutParams{
		User:                 user,
		Workspace:            workspace,
		Organization:         org,
		ExistingSubscription: existingSubscription,
		Plan:                 req.Plan,
		Quantity:             quantity,
		ClearOnboarding:      user.Onboarding,
	})
	if err != nil {
		respondAPIError(w, http.StatusBadGateway, "billing_checkout_failed", err.Error())
		return
	}

	now := time.Now()
	subscription := existingSubscription
	if subscription == nil {
		subscription = &Subscription{
			ID:             newID("sub"),
			WorkspaceID:    workspace.ID,
			OrganizationID: workspace.OrganizationID,
			CreatedAt:      now,
		}
	}
	subscription.WorkspaceID = workspace.ID
	subscription.OrganizationID = workspace.OrganizationID
	subscription.Provider = result.Provider
	subscription.Plan = req.Plan
	subscription.Status = "pending"
	subscription.ProviderCustomerID = firstNonEmpty(result.CheckoutState.CustomerID, subscription.ProviderCustomerID)
	subscription.ProviderSubscriptionID = firstNonEmpty(result.CheckoutState.SubscriptionID, subscription.ProviderSubscriptionID)
	subscription.ProviderCheckoutSessionID = result.CheckoutState.CheckoutSessionID
	subscription.BilledQuantity = quantity
	subscription.CancelAtPeriodEnd = false
	subscription.ScheduledPlan = ""
	subscription.UpdatedAt = now
	if err := h.store.UpsertSubscription(subscription); err != nil {
		respondAPIError(w, http.StatusInternalServerError, "billing_subscription_failed", err.Error())
		return
	}

	respondJSON(w, http.StatusOK, BillingCheckoutResponse{
		Provider:     result.Provider,
		Plan:         req.Plan,
		CheckoutURL:  result.CheckoutURL,
		Completed:    false,
		Subscription: subscription,
	})
}

func (h *APIHandler) handleBillingPortal(w http.ResponseWriter, r *http.Request) {
	session := h.sessionFromRequest(r)
	if session == nil || session.UserID == "" {
		respondAPIError(w, http.StatusUnauthorized, "auth_required", "You must be signed in to open billing portal.")
		return
	}

	var req billingPortalRequest
	if r.Body != nil {
		_ = json.NewDecoder(r.Body).Decode(&req)
	}

	workspace, err := h.workspaceForSession(session)
	if err != nil || workspace == nil {
		respondAPIError(w, http.StatusInternalServerError, "workspace_not_found", "Workspace not found.")
		return
	}
	_, orgMember, err := h.billingWorkspaceContext(session, workspace)
	if err != nil {
		respondAPIError(w, http.StatusInternalServerError, "workspace_context_failed", err.Error())
		return
	}
	if !h.canManageBillingPlan(session, workspace, orgMember) {
		respondAPIError(w, http.StatusForbidden, "plan_forbidden", "You do not have permission to manage billing for this workspace.")
		return
	}

	subscription := h.organizationSubscription(workspace.OrganizationID, workspace.ID)
	if !subscriptionGrantsPlan(subscription, time.Now()) || subscription == nil || strings.TrimSpace(subscription.ProviderCustomerID) == "" {
		respondAPIError(w, http.StatusConflict, "billing_portal_unavailable", "This workspace does not have an active Stripe subscription to manage.")
		return
	}
	if strings.TrimSpace(subscription.Provider) != subscriptionProviderStripe {
		respondAPIError(w, http.StatusConflict, "billing_portal_unavailable", "Billing portal is only available for Stripe-managed subscriptions.")
		return
	}

	targetPlan := req.Plan
	if targetPlan == "" {
		targetPlan = subscription.Plan
	}
	if targetPlan == PlanEnterprise {
		respondAPIError(w, http.StatusConflict, "enterprise_sales_required", "Enterprise plan changes are handled manually. Contact sales to continue.")
		return
	}
	if subscription.Plan == PlanTeam && (targetPlan == PlanFree || targetPlan == PlanPro) {
		activeMembers, err := h.store.CountActiveOrganizationMembers(workspace.OrganizationID)
		if err != nil {
			respondAPIError(w, http.StatusInternalServerError, "organization_members_failed", err.Error())
			return
		}
		if activeMembers > 1 {
			respondAPIError(w, http.StatusConflict, "team_downgrade_blocked", "Reduce the team to one active member before downgrading from Team.")
			return
		}
	}

	if _, disabled := h.subscriptionBilling.(*disabledSubscriptionBillingProvider); disabled {
		respondAPIError(w, http.StatusNotImplemented, "billing_not_configured", "Stripe subscription billing is not configured.")
		return
	}

	url, err := h.subscriptionBilling.CreatePortalSession(r.Context(), subscription.ProviderCustomerID)
	if err != nil {
		respondAPIError(w, http.StatusBadGateway, "billing_portal_failed", err.Error())
		return
	}

	respondJSON(w, http.StatusOK, BillingPortalResponse{
		Provider: subscription.Provider,
		URL:      url,
	})
}

func (h *APIHandler) handleBillingCheckoutSync(w http.ResponseWriter, r *http.Request) {
	session := h.sessionFromRequest(r)
	if session == nil || session.UserID == "" {
		respondAPIError(w, http.StatusUnauthorized, "auth_required", "You must be signed in to confirm checkout.")
		return
	}
	if _, disabled := h.subscriptionBilling.(*disabledSubscriptionBillingProvider); disabled {
		respondAPIError(w, http.StatusNotImplemented, "billing_not_configured", "Stripe subscription billing is not configured.")
		return
	}

	checkoutSessionID := strings.TrimSpace(chi.URLParam(r, "sessionId"))
	if checkoutSessionID == "" {
		respondAPIError(w, http.StatusBadRequest, "invalid_session_id", "Checkout session ID is required.")
		return
	}

	subscription, err := h.store.GetSubscriptionByProviderCheckoutSessionID(checkoutSessionID)
	if err != nil {
		respondAPIError(w, http.StatusNotFound, "billing_checkout_not_found", "Billing checkout session not found.")
		return
	}
	if !h.subscriptionBelongsToWorkspaceSession(subscription, session) {
		respondAPIError(w, http.StatusForbidden, "billing_checkout_forbidden", "That checkout session does not belong to the current workspace.")
		return
	}

	checkoutState, err := h.subscriptionBilling.GetCheckoutSession(r.Context(), checkoutSessionID)
	if err != nil {
		respondAPIError(w, http.StatusBadGateway, "billing_checkout_sync_failed", err.Error())
		return
	}
	if !checkoutState.Completed || strings.TrimSpace(checkoutState.SubscriptionID) == "" {
		respondJSON(w, http.StatusOK, BillingCheckoutSyncResponse{
			Provider:     checkoutState.Provider,
			Completed:    false,
			Subscription: subscription,
		})
		return
	}

	user, err := h.store.GetUserByID(session.UserID)
	if err != nil {
		respondAPIError(w, http.StatusInternalServerError, "user_lookup_failed", err.Error())
		return
	}
	workspace, err := h.store.GetWorkspaceRecord(subscription.WorkspaceID)
	if err != nil {
		respondAPIError(w, http.StatusInternalServerError, "workspace_not_found", err.Error())
		return
	}
	state, err := h.subscriptionBilling.GetSubscription(r.Context(), checkoutState.SubscriptionID)
	if err != nil {
		respondAPIError(w, http.StatusBadGateway, "billing_subscription_sync_failed", err.Error())
		return
	}

	applied, err := h.applyStripeSubscriptionState(r.Context(), subscription, workspace, user, state, checkoutSessionID)
	if err != nil {
		respondAPIError(w, http.StatusInternalServerError, "billing_apply_failed", err.Error())
		return
	}

	response := h.buildSessionResponse(r.WithContext(context.WithValue(r.Context(), sessionContextKey, session)))
	respondJSON(w, http.StatusOK, BillingCheckoutSyncResponse{
		Provider:     state.Provider,
		Completed:    true,
		Session:      &response,
		Subscription: applied,
	})
}

func (h *APIHandler) handleBillingWebhook(w http.ResponseWriter, r *http.Request) {
	if strings.TrimSpace(h.config.Stripe.WebhookSecret) == "" {
		respondAPIError(w, http.StatusNotImplemented, "billing_not_configured", "Stripe subscription webhook is not configured.")
		return
	}

	payload, err := io.ReadAll(r.Body)
	if err != nil {
		respondAPIError(w, http.StatusBadRequest, "invalid_request", "Failed to read billing webhook payload.")
		return
	}
	if err := verifyStripeWebhookSignature(payload, r.Header.Get("Stripe-Signature"), h.config.Stripe.WebhookSecret, time.Now()); err != nil {
		respondAPIError(w, http.StatusUnauthorized, "webhook_unauthorized", err.Error())
		return
	}

	var envelope stripeEventEnvelope
	if err := json.Unmarshal(payload, &envelope); err != nil {
		respondAPIError(w, http.StatusBadRequest, "invalid_request", "Invalid billing webhook payload.")
		return
	}

	eventSubscriptionID := ""
	switch envelope.Type {
	case "checkout.session.completed":
		var checkoutSession stripeBillingCheckoutSessionResponse
		if err := json.Unmarshal(envelope.Data.Object, &checkoutSession); err != nil {
			respondAPIError(w, http.StatusBadRequest, "invalid_request", "Invalid checkout session payload.")
			return
		}
		if strings.TrimSpace(checkoutSession.Subscription) != "" {
			state, err := h.subscriptionBilling.GetSubscription(r.Context(), checkoutSession.Subscription)
			if err != nil {
				respondAPIError(w, http.StatusBadGateway, "billing_subscription_sync_failed", err.Error())
				return
			}
			subscription, workspace, user, err := h.resolveSubscriptionRecordForStripe(checkoutSession.Metadata, checkoutSession.Subscription, checkoutSession.ID)
			if err != nil {
				respondAPIError(w, http.StatusInternalServerError, "billing_subscription_failed", err.Error())
				return
			}
			if _, err := h.applyStripeSubscriptionState(r.Context(), subscription, workspace, user, state, checkoutSession.ID); err != nil {
				respondAPIError(w, http.StatusInternalServerError, "billing_apply_failed", err.Error())
				return
			}
			eventSubscriptionID = subscription.ID
		}
	case "customer.subscription.created", "customer.subscription.updated", "customer.subscription.deleted":
		var stripeSubscription stripeBillingSubscriptionResponse
		if err := json.Unmarshal(envelope.Data.Object, &stripeSubscription); err != nil {
			respondAPIError(w, http.StatusBadRequest, "invalid_request", "Invalid subscription payload.")
			return
		}
		state := h.subscriptionBilling.(*stripeSubscriptionBillingProvider).stateFromStripeSubscription(&stripeSubscription)
		subscription, workspace, user, err := h.resolveSubscriptionRecordForStripe(stripeSubscription.Metadata, stripeSubscription.ID, "")
		if err != nil {
			respondAPIError(w, http.StatusInternalServerError, "billing_subscription_failed", err.Error())
			return
		}
		if _, err := h.applyStripeSubscriptionState(r.Context(), subscription, workspace, user, state, subscription.ProviderCheckoutSessionID); err != nil {
			respondAPIError(w, http.StatusInternalServerError, "billing_apply_failed", err.Error())
			return
		}
		eventSubscriptionID = subscription.ID
	default:
		respondJSON(w, http.StatusOK, map[string]any{"ok": true, "ignored": envelope.Type})
		return
	}

	if strings.TrimSpace(eventSubscriptionID) != "" {
		if err := h.store.CreateSubscriptionEvent(&SubscriptionEvent{
			ID:              newID("subevt"),
			SubscriptionID:  eventSubscriptionID,
			EventType:       envelope.Type,
			ProviderEventID: envelope.ID,
			Payload:         string(payload),
			CreatedAt:       time.Now(),
		}); err != nil {
			respondAPIError(w, http.StatusInternalServerError, "billing_event_failed", err.Error())
			return
		}
	}

	respondJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func (h *APIHandler) billingWorkspaceContext(session *SessionRecord, workspace *Workspace) (*Organization, *OrganizationMember, error) {
	if workspace == nil || strings.TrimSpace(workspace.OrganizationID) == "" {
		return nil, nil, nil
	}
	org, err := h.store.GetOrganizationRecord(workspace.OrganizationID)
	if err != nil {
		return nil, nil, err
	}
	member, err := h.store.GetOrganizationMemberByUser(workspace.OrganizationID, session.UserID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, nil, err
	}
	if errors.Is(err, sql.ErrNoRows) {
		member = nil
	}
	return org, member, nil
}

func (h *APIHandler) canManageBillingPlan(session *SessionRecord, workspace *Workspace, orgMember *OrganizationMember) bool {
	if session == nil || workspace == nil {
		return false
	}
	if strings.TrimSpace(workspace.OrganizationID) == "" {
		return workspace.OwnerUserID == session.UserID
	}
	return orgMember != nil && orgMember.Status == "active" && canManageOrganizationPlan(orgMember.Role)
}

func (h *APIHandler) ensureTeamOrganizationForBilling(workspace *Workspace, user *User, session *SessionRecord) (*Organization, error) {
	if workspace == nil {
		return nil, sql.ErrNoRows
	}
	if strings.TrimSpace(workspace.OrganizationID) != "" {
		return h.store.GetOrganizationRecord(workspace.OrganizationID)
	}
	if user == nil {
		switch {
		case session != nil && strings.TrimSpace(session.UserID) != "":
			resolvedUser, err := h.store.GetUserByID(session.UserID)
			if err != nil {
				return nil, err
			}
			user = resolvedUser
		case strings.TrimSpace(workspace.OwnerUserID) != "":
			resolvedUser, err := h.store.GetUserByID(workspace.OwnerUserID)
			if err != nil {
				return nil, err
			}
			user = resolvedUser
		default:
			return nil, sql.ErrNoRows
		}
	}

	now := time.Now()
	org := &Organization{
		ID:        newID("org"),
		Name:      firstNonEmpty(strings.TrimSpace(user.DisplayName), workspace.Name, "Vutadex Team"),
		Slug:      slugify(firstNonEmpty(strings.TrimSpace(workspace.Slug), user.Email)),
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := h.store.CreateOrganizationRecord(org); err != nil {
		return nil, err
	}
	member := &OrganizationMember{
		ID:             newID("orgmem"),
		OrganizationID: org.ID,
		UserID:         user.ID,
		Email:          user.Email,
		Role:           "owner",
		Status:         "active",
		JoinedAt:       now,
		CreatedAt:      now,
	}
	if err := h.store.CreateOrganizationMemberRecord(member); err != nil {
		return nil, err
	}
	if err := h.ensureWorkspaceAttachedToOrganization(workspace, org); err != nil {
		return nil, err
	}
	workspace.OrganizationID = org.ID
	return org, nil
}

func (h *APIHandler) resolveSubscriptionRecordForStripe(metadata map[string]string, providerSubscriptionID string, checkoutSessionID string) (*Subscription, *Workspace, *User, error) {
	if strings.TrimSpace(providerSubscriptionID) != "" {
		if subscription, err := h.store.GetSubscriptionByProviderSubscriptionID(providerSubscriptionID); err == nil {
			workspace, err := h.store.GetWorkspaceRecord(subscription.WorkspaceID)
			if err != nil {
				return nil, nil, nil, err
			}
			var user *User
			if userID := strings.TrimSpace(metadata["user_id"]); userID != "" {
				user, _ = h.store.GetUserByID(userID)
			} else if strings.TrimSpace(workspace.OwnerUserID) != "" {
				user, _ = h.store.GetUserByID(workspace.OwnerUserID)
			}
			return subscription, workspace, user, nil
		}
	}
	if strings.TrimSpace(checkoutSessionID) != "" {
		if subscription, err := h.store.GetSubscriptionByProviderCheckoutSessionID(checkoutSessionID); err == nil {
			workspace, err := h.store.GetWorkspaceRecord(subscription.WorkspaceID)
			if err != nil {
				return nil, nil, nil, err
			}
			var user *User
			if userID := strings.TrimSpace(metadata["user_id"]); userID != "" {
				user, _ = h.store.GetUserByID(userID)
			} else if strings.TrimSpace(workspace.OwnerUserID) != "" {
				user, _ = h.store.GetUserByID(workspace.OwnerUserID)
			}
			return subscription, workspace, user, nil
		}
	}

	workspaceID := strings.TrimSpace(metadata["workspace_id"])
	if workspaceID == "" {
		return nil, nil, nil, sql.ErrNoRows
	}
	workspace, err := h.store.GetWorkspaceRecord(workspaceID)
	if err != nil {
		return nil, nil, nil, err
	}
	var subscription *Subscription
	if existing, err := h.store.GetSubscriptionForWorkspace(workspace.ID); err == nil {
		subscription = existing
	} else {
		subscription = &Subscription{
			ID:          newID("sub"),
			WorkspaceID: workspace.ID,
			CreatedAt:   time.Now(),
		}
	}
	userID := strings.TrimSpace(metadata["user_id"])
	var user *User
	if userID != "" {
		user, _ = h.store.GetUserByID(userID)
	} else if strings.TrimSpace(workspace.OwnerUserID) != "" {
		user, _ = h.store.GetUserByID(workspace.OwnerUserID)
	}
	return subscription, workspace, user, nil
}

func (h *APIHandler) applyStripeSubscriptionState(ctx context.Context, subscription *Subscription, workspace *Workspace, user *User, state *subscriptionBillingState, checkoutSessionID string) (*Subscription, error) {
	if subscription == nil || workspace == nil || state == nil {
		return nil, sql.ErrNoRows
	}

	if state.Plan == PlanTeam {
		org, err := h.ensureTeamOrganizationForBilling(workspace, user, nil)
		if err != nil {
			return nil, err
		}
		workspace.OrganizationID = org.ID
	}

	subscription.WorkspaceID = workspace.ID
	subscription.OrganizationID = workspace.OrganizationID
	applySubscriptionState(subscription, state)
	if strings.TrimSpace(checkoutSessionID) != "" {
		subscription.ProviderCheckoutSessionID = checkoutSessionID
	}
	if state.Plan == PlanFree || state.Plan == PlanGuest {
		subscription.Plan = subscription.Plan
	}
	subscription.UpdatedAt = time.Now()
	if subscription.CreatedAt.IsZero() {
		subscription.CreatedAt = subscription.UpdatedAt
	}
	if err := h.store.UpsertSubscription(subscription); err != nil {
		return nil, err
	}
	if user != nil && user.Onboarding && (state.Plan == PlanPro || state.Plan == PlanTeam) {
		if err := h.store.UpdateUserOnboarding(user.ID, false); err != nil {
			return nil, err
		}
	}
	return subscription, nil
}

func (h *APIHandler) subscriptionBelongsToWorkspaceSession(subscription *Subscription, session *SessionRecord) bool {
	if subscription == nil || session == nil {
		return false
	}
	if subscription.WorkspaceID == session.WorkspaceID {
		return true
	}
	if strings.TrimSpace(subscription.OrganizationID) == "" {
		return false
	}
	workspace, err := h.workspaceForSession(session)
	if err != nil || workspace == nil {
		return false
	}
	return strings.TrimSpace(workspace.OrganizationID) != "" && workspace.OrganizationID == subscription.OrganizationID
}

func (h *APIHandler) reconcileOrganizationSeatBilling(ctx context.Context, organizationID, workspaceID string) error {
	if strings.TrimSpace(organizationID) == "" {
		return nil
	}
	subscription := h.organizationSubscription(organizationID, workspaceID)
	if subscription == nil || strings.TrimSpace(subscription.Provider) != subscriptionProviderStripe || subscription.Plan != PlanTeam || strings.TrimSpace(subscription.ProviderSubscriptionID) == "" {
		return nil
	}
	activeMembers, err := h.store.CountActiveOrganizationMembers(organizationID)
	if err != nil {
		return err
	}
	targetQuantity := teamSeatQuantity(activeMembers)
	if targetQuantity == subscription.BilledQuantity {
		return nil
	}
	prorationBehavior := "none"
	if targetQuantity > subscription.BilledQuantity {
		prorationBehavior = "create_prorations"
	}
	state, err := h.subscriptionBilling.UpdateSubscriptionQuantity(ctx, subscription, targetQuantity, prorationBehavior)
	if err != nil {
		return err
	}
	applySubscriptionState(subscription, state)
	subscription.UpdatedAt = time.Now()
	return h.store.UpsertSubscription(subscription)
}

package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
)

const organizationInviteTTL = 7 * 24 * time.Hour

func validOrganizationRole(role string) bool {
	switch role {
	case "read", "edit", "admin", "owner":
		return true
	default:
		return false
	}
}

func validOrganizationMemberStatus(status string) bool {
	switch status {
	case "invited", "active", "removed", "declined", "expired":
		return true
	default:
		return false
	}
}

func canManageOrganizationMembers(role string) bool {
	return role == "admin" || role == "owner"
}

func canManageOrganizationPlan(role string) bool {
	return role == "admin" || role == "owner"
}

func canWriteWorkspaceContent(role string) bool {
	return role == "edit" || role == "admin" || role == "owner"
}

func canEditOrganizationMetadata(role string) bool {
	return role == "admin" || role == "owner"
}

func canDeleteOrganization(role string) bool {
	return role == "owner"
}

func validManagedPlan(plan Plan) bool {
	switch plan {
	case PlanFree, PlanPro, PlanTeam, PlanEnterprise:
		return true
	default:
		return false
	}
}

func (h *APIHandler) organizationAccessForRequest(r *http.Request, orgID string) (*Organization, *OrganizationMember, *Workspace, *SessionRecord, error) {
	session := h.sessionFromRequest(r)
	if session == nil || session.UserID == "" {
		return nil, nil, nil, nil, sql.ErrNoRows
	}
	org, err := h.store.GetOrganizationRecord(orgID)
	if err != nil {
		return nil, nil, nil, session, err
	}
	member, err := h.store.GetOrganizationMemberByUser(orgID, session.UserID)
	if err != nil {
		return nil, nil, nil, session, err
	}
	workspace, err := h.store.GetWorkspaceForOrganization(orgID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, nil, nil, session, err
	}
	if errors.Is(err, sql.ErrNoRows) {
		workspace = nil
	}
	return org, member, workspace, session, nil
}

func (h *APIHandler) organizationSubscription(orgID string, workspaceID string) *Subscription {
	if workspaceID != "" {
		if subscription, err := h.store.GetSubscriptionForWorkspace(workspaceID); err == nil {
			return subscription
		}
	}
	if orgID != "" {
		if subscription, err := h.store.GetSubscriptionForOrganization(orgID); err == nil {
			return subscription
		}
	}
	return nil
}

func (h *APIHandler) buildOrganizationDetail(orgID, userID string) (*OrganizationDetail, error) {
	org, err := h.store.GetOrganizationRecord(orgID)
	if err != nil {
		return nil, err
	}
	member, err := h.store.GetOrganizationMemberByUser(orgID, userID)
	if err != nil {
		return nil, err
	}
	members, err := h.store.ListOrganizationMembers(orgID)
	if err != nil {
		return nil, err
	}
	workspace, err := h.store.GetWorkspaceForOrganization(orgID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}
	if errors.Is(err, sql.ErrNoRows) {
		workspace = nil
	}

	detail := &OrganizationDetail{
		Organization:     *org,
		Workspace:        workspace,
		Membership:       *member,
		Members:          members,
		CanManagePlan:    canManageOrganizationPlan(member.Role),
		CanManageMembers: canManageOrganizationMembers(member.Role),
		CanEdit:          canWriteWorkspaceContent(member.Role),
	}
	if workspace != nil {
		detail.Subscription = h.organizationSubscription(org.ID, workspace.ID)
	}
	return detail, nil
}

func (h *APIHandler) ensureOrganizationWorkspace(org *Organization, session *SessionRecord) (*Workspace, error) {
	if workspace, err := h.store.GetWorkspaceForOrganization(org.ID); err == nil {
		return workspace, nil
	} else if !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	collectionID := newID("col")
	collection := NewCollection()
	if err := h.store.CreateCollectionRecord(collectionID, org.Name, collection); err != nil {
		return nil, err
	}
	for _, nt := range builtins() {
		ntCopy := nt
		if err := h.store.CreateNoteType(collectionID, &ntCopy); err != nil {
			return nil, err
		}
	}
	defaultDeck := collection.NewDeck("Default")
	if err := h.store.CreateDeckInCollection(collectionID, defaultDeck); err != nil {
		return nil, err
	}

	now := time.Now()
	workspace := &Workspace{
		ID:             newID("ws"),
		Name:           org.Name,
		Slug:           org.Slug,
		CollectionID:   collectionID,
		OwnerUserID:    session.UserID,
		OrganizationID: org.ID,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
	if err := h.store.CreateWorkspaceRecord(workspace); err != nil {
		return nil, err
	}
	return workspace, nil
}

func (h *APIHandler) ensureWorkspaceAttachedToOrganization(workspace *Workspace, org *Organization) error {
	workspace.OrganizationID = org.ID
	workspace.UpdatedAt = time.Now()
	return h.store.UpdateWorkspaceRecord(workspace)
}

func (h *APIHandler) activeOrganizationMemberForCurrentWorkspace(r *http.Request) (*Workspace, *OrganizationMember, *SessionRecord, error) {
	session := h.sessionFromRequest(r)
	if session == nil || session.UserID == "" {
		return nil, nil, nil, sql.ErrNoRows
	}
	workspace, err := h.workspaceForSession(session)
	if err != nil || workspace == nil {
		return nil, nil, session, err
	}
	if workspace.OrganizationID == "" {
		return workspace, nil, session, nil
	}
	member, err := h.store.GetOrganizationMemberByUser(workspace.OrganizationID, session.UserID)
	if err != nil {
		return workspace, nil, session, err
	}
	return workspace, member, session, nil
}

func (h *APIHandler) requireWorkspaceWritePermission(w http.ResponseWriter, r *http.Request) bool {
	workspace, member, _, err := h.activeOrganizationMemberForCurrentWorkspace(r)
	if errors.Is(err, sql.ErrNoRows) {
		respondAPIError(w, http.StatusUnauthorized, "auth_required", "You must be signed in to change workspace content.")
		return false
	}
	if err != nil {
		respondAPIError(w, http.StatusInternalServerError, "workspace_access_failed", err.Error())
		return false
	}
	if workspace == nil || workspace.OrganizationID == "" {
		return true
	}
	if member == nil || member.Status != "active" || !canWriteWorkspaceContent(member.Role) {
		respondAPIError(w, http.StatusForbidden, "workspace_write_forbidden", "Your team role does not allow editing content in this workspace.")
		return false
	}
	return true
}

func (h *APIHandler) CompleteOnboardingPlanSelection(w http.ResponseWriter, r *http.Request) {
	session := h.sessionFromRequest(r)
	if session == nil || session.UserID == "" {
		respondAPIError(w, http.StatusUnauthorized, "auth_required", "You must be signed in to choose a plan.")
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

	var req UpdateWorkspacePlanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondAPIError(w, http.StatusBadRequest, "invalid_request", "Invalid request body")
		return
	}
	if !validManagedPlan(req.Plan) {
		respondAPIError(w, http.StatusBadRequest, "invalid_plan", "Choose free, pro, team, or enterprise.")
		return
	}
	if _, disabled := h.subscriptionBilling.(*disabledSubscriptionBillingProvider); !disabled {
		switch req.Plan {
		case PlanFree:
			// Free onboarding remains local and immediate.
		case PlanPro, PlanTeam:
			respondAPIError(w, http.StatusConflict, "billing_checkout_required", "Paid onboarding uses Stripe checkout. Start checkout from the billing API instead.")
			return
		case PlanEnterprise:
			respondAPIError(w, http.StatusConflict, "enterprise_sales_required", "Enterprise onboarding is handled manually. Contact sales to continue.")
			return
		}
	}

	if (req.Plan == PlanTeam || req.Plan == PlanEnterprise) && workspace.OrganizationID == "" {
		org := &Organization{
			ID:        newID("org"),
			Name:      firstNonEmpty(strings.TrimSpace(user.DisplayName), workspace.Name, "Vutadex Team"),
			Slug:      slugify(firstNonEmpty(strings.TrimSpace(workspace.Slug), user.Email)),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		if err := h.store.CreateOrganizationRecord(org); err != nil {
			respondAPIError(w, http.StatusInternalServerError, "organization_create_failed", err.Error())
			return
		}
		member := &OrganizationMember{
			ID:             newID("orgmem"),
			OrganizationID: org.ID,
			UserID:         user.ID,
			Email:          user.Email,
			Role:           "owner",
			Status:         "active",
			JoinedAt:       time.Now(),
			CreatedAt:      time.Now(),
		}
		if err := h.store.CreateOrganizationMemberRecord(member); err != nil {
			respondAPIError(w, http.StatusInternalServerError, "organization_member_failed", err.Error())
			return
		}
		if err := h.ensureWorkspaceAttachedToOrganization(workspace, org); err != nil {
			respondAPIError(w, http.StatusInternalServerError, "workspace_attach_failed", err.Error())
			return
		}
		workspace.OrganizationID = org.ID
	}

	response, err := h.applyWorkspacePlanChange(r, session, workspace, req.Plan, true)
	if err != nil {
		respondAPIError(w, http.StatusInternalServerError, "plan_update_failed", err.Error())
		return
	}
	respondJSON(w, http.StatusOK, response)
}

func (h *APIHandler) applyWorkspacePlanChange(r *http.Request, session *SessionRecord, workspace *Workspace, plan Plan, clearOnboarding bool) (*AuthSessionResponse, error) {
	now := time.Now()
	subscription := h.organizationSubscription(workspace.OrganizationID, workspace.ID)
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
	subscription.Plan = plan
	subscription.Status = "active"
	if strings.TrimSpace(subscription.Provider) == "" {
		subscription.Provider = "manual"
	}
	subscription.UpdatedAt = now
	if err := h.store.UpsertSubscription(subscription); err != nil {
		return nil, err
	}
	if clearOnboarding {
		if err := h.store.UpdateUserOnboarding(session.UserID, false); err != nil {
			return nil, err
		}
	}
	response := h.buildSessionResponse(r.WithContext(context.WithValue(r.Context(), sessionContextKey, session)))
	return &response, nil
}

func (h *APIHandler) UpdateWorkspacePlan(w http.ResponseWriter, r *http.Request) {
	session := h.sessionFromRequest(r)
	if session == nil || session.UserID == "" {
		respondAPIError(w, http.StatusUnauthorized, "auth_required", "You must be signed in to manage plans.")
		return
	}
	workspaceID := chi.URLParam(r, "workspaceId")
	if workspaceID == "" || workspaceID != session.WorkspaceID {
		respondAPIError(w, http.StatusForbidden, "workspace_mismatch", "Plan changes are only supported for the current workspace.")
		return
	}
	workspace, err := h.workspaceForSession(session)
	if err != nil || workspace == nil {
		respondAPIError(w, http.StatusInternalServerError, "workspace_not_found", "Workspace not found.")
		return
	}

	var req UpdateWorkspacePlanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondAPIError(w, http.StatusBadRequest, "invalid_request", "Invalid request body")
		return
	}
	if !validManagedPlan(req.Plan) {
		respondAPIError(w, http.StatusBadRequest, "invalid_plan", "Choose free, pro, team, or enterprise.")
		return
	}
	if _, disabled := h.subscriptionBilling.(*disabledSubscriptionBillingProvider); !disabled {
		currentPlan := h.planForRequest(r, session)
		switch {
		case req.Plan == PlanEnterprise:
			respondAPIError(w, http.StatusConflict, "enterprise_sales_required", "Enterprise plan changes are handled manually. Contact sales to continue.")
			return
		case (currentPlan == PlanFree || currentPlan == PlanGuest) && (req.Plan == PlanPro || req.Plan == PlanTeam):
			respondAPIError(w, http.StatusConflict, "billing_checkout_required", "Use Stripe checkout to upgrade this workspace to a paid plan.")
			return
		case currentPlan != PlanFree && currentPlan != PlanGuest && currentPlan != req.Plan:
			respondAPIError(w, http.StatusConflict, "billing_portal_required", "Use the billing portal to manage an existing paid subscription.")
			return
		}
	}

	if workspace.OrganizationID == "" {
		if workspace.OwnerUserID != session.UserID {
			respondAPIError(w, http.StatusForbidden, "plan_forbidden", "Only the workspace owner can manage this plan.")
			return
		}
		if req.Plan == PlanTeam || req.Plan == PlanEnterprise {
			user, err := h.store.GetUserByID(session.UserID)
			if err != nil {
				respondAPIError(w, http.StatusInternalServerError, "user_lookup_failed", err.Error())
				return
			}
			org := &Organization{
				ID:        newID("org"),
				Name:      firstNonEmpty(strings.TrimSpace(user.DisplayName), workspace.Name, "Vutadex Team"),
				Slug:      slugify(firstNonEmpty(strings.TrimSpace(workspace.Slug), user.Email)),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			if err := h.store.CreateOrganizationRecord(org); err != nil {
				respondAPIError(w, http.StatusInternalServerError, "organization_create_failed", err.Error())
				return
			}
			member := &OrganizationMember{
				ID:             newID("orgmem"),
				OrganizationID: org.ID,
				UserID:         user.ID,
				Email:          user.Email,
				Role:           "owner",
				Status:         "active",
				JoinedAt:       time.Now(),
				CreatedAt:      time.Now(),
			}
			if err := h.store.CreateOrganizationMemberRecord(member); err != nil {
				respondAPIError(w, http.StatusInternalServerError, "organization_member_failed", err.Error())
				return
			}
			if err := h.ensureWorkspaceAttachedToOrganization(workspace, org); err != nil {
				respondAPIError(w, http.StatusInternalServerError, "workspace_attach_failed", err.Error())
				return
			}
			workspace.OrganizationID = org.ID
		}
	} else {
		member, err := h.store.GetOrganizationMemberByUser(workspace.OrganizationID, session.UserID)
		if err != nil || member.Status != "active" || !canManageOrganizationPlan(member.Role) {
			respondAPIError(w, http.StatusForbidden, "plan_forbidden", "Only team admins and owners can manage this plan.")
			return
		}
	}

	response, err := h.applyWorkspacePlanChange(r, session, workspace, req.Plan, false)
	if err != nil {
		respondAPIError(w, http.StatusInternalServerError, "plan_update_failed", err.Error())
		return
	}
	respondJSON(w, http.StatusOK, response)
}

func (h *APIHandler) GetOrganization(w http.ResponseWriter, r *http.Request) {
	orgID := chi.URLParam(r, "orgId")
	org, member, _, session, err := h.organizationAccessForRequest(r, orgID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondAPIError(w, http.StatusNotFound, "organization_not_found", "Organization not found.")
			return
		}
		respondAPIError(w, http.StatusInternalServerError, "organization_access_failed", err.Error())
		return
	}
	if member.Status != "active" {
		respondAPIError(w, http.StatusForbidden, "organization_forbidden", "Your membership is not active.")
		return
	}
	detail, err := h.buildOrganizationDetail(org.ID, session.UserID)
	if err != nil {
		respondAPIError(w, http.StatusInternalServerError, "organization_detail_failed", err.Error())
		return
	}
	respondJSON(w, http.StatusOK, detail)
}

func (h *APIHandler) ListOrganizationMembers(w http.ResponseWriter, r *http.Request) {
	orgID := chi.URLParam(r, "orgId")
	_, member, _, _, err := h.organizationAccessForRequest(r, orgID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondAPIError(w, http.StatusNotFound, "organization_not_found", "Organization not found.")
			return
		}
		respondAPIError(w, http.StatusInternalServerError, "organization_access_failed", err.Error())
		return
	}
	if member.Status != "active" {
		respondAPIError(w, http.StatusForbidden, "organization_forbidden", "Your membership is not active.")
		return
	}
	members, err := h.store.ListOrganizationMembers(orgID)
	if err != nil {
		respondAPIError(w, http.StatusInternalServerError, "organization_members_failed", err.Error())
		return
	}
	respondJSON(w, http.StatusOK, map[string]any{"members": members})
}

func (h *APIHandler) UpdateOrganization(w http.ResponseWriter, r *http.Request) {
	orgID := chi.URLParam(r, "orgId")
	org, member, workspace, session, err := h.organizationAccessForRequest(r, orgID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondAPIError(w, http.StatusNotFound, "organization_not_found", "Organization not found.")
			return
		}
		respondAPIError(w, http.StatusInternalServerError, "organization_access_failed", err.Error())
		return
	}
	if member.Status != "active" || !canEditOrganizationMetadata(member.Role) {
		respondAPIError(w, http.StatusForbidden, "organization_forbidden", "Only team admins and owners can update team settings.")
		return
	}

	var req UpdateOrganizationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondAPIError(w, http.StatusBadRequest, "invalid_request", "Invalid request body")
		return
	}
	if trimmed := strings.TrimSpace(req.Name); trimmed != "" {
		org.Name = sanitizeHTML(trimmed)
	}
	if trimmed := strings.TrimSpace(req.Slug); trimmed != "" {
		org.Slug = slugify(trimmed)
	}
	org.UpdatedAt = time.Now()
	if err := h.store.UpdateOrganizationRecord(org); err != nil {
		respondAPIError(w, http.StatusInternalServerError, "organization_update_failed", err.Error())
		return
	}
	if workspace != nil {
		workspace.Name = org.Name
		workspace.UpdatedAt = time.Now()
		if err := h.store.UpdateWorkspaceRecord(workspace); err != nil {
			respondAPIError(w, http.StatusInternalServerError, "workspace_update_failed", err.Error())
			return
		}
	}
	detail, err := h.buildOrganizationDetail(org.ID, session.UserID)
	if err != nil {
		respondAPIError(w, http.StatusInternalServerError, "organization_detail_failed", err.Error())
		return
	}
	respondJSON(w, http.StatusOK, detail)
}

func (h *APIHandler) DeleteOrganization(w http.ResponseWriter, r *http.Request) {
	orgID := chi.URLParam(r, "orgId")
	org, member, workspace, _, err := h.organizationAccessForRequest(r, orgID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondAPIError(w, http.StatusNotFound, "organization_not_found", "Organization not found.")
			return
		}
		respondAPIError(w, http.StatusInternalServerError, "organization_access_failed", err.Error())
		return
	}
	if member.Status != "active" || !canDeleteOrganization(member.Role) {
		respondAPIError(w, http.StatusForbidden, "organization_forbidden", "Only the team owner can delete this team.")
		return
	}
	members, err := h.store.ListOrganizationMembers(org.ID)
	if err != nil {
		respondAPIError(w, http.StatusInternalServerError, "organization_members_failed", err.Error())
		return
	}
	now := time.Now()
	for i := range members {
		members[i].Status = "removed"
		members[i].RemovedAt = now
		if err := h.store.UpdateOrganizationMember(&members[i]); err != nil {
			respondAPIError(w, http.StatusInternalServerError, "organization_member_update_failed", err.Error())
			return
		}
	}
	if workspace != nil {
		workspace.OrganizationID = ""
		workspace.UpdatedAt = now
		if err := h.store.UpdateWorkspaceRecord(workspace); err != nil {
			respondAPIError(w, http.StatusInternalServerError, "workspace_update_failed", err.Error())
			return
		}
	}
	if err := h.store.DeleteOrganizationRecord(org.ID); err != nil {
		respondAPIError(w, http.StatusInternalServerError, "organization_delete_failed", err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *APIHandler) UpdateOrganizationMember(w http.ResponseWriter, r *http.Request) {
	orgID := chi.URLParam(r, "orgId")
	targetID := chi.URLParam(r, "memberId")
	_, currentMember, _, _, err := h.organizationAccessForRequest(r, orgID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondAPIError(w, http.StatusNotFound, "organization_not_found", "Organization not found.")
			return
		}
		respondAPIError(w, http.StatusInternalServerError, "organization_access_failed", err.Error())
		return
	}
	if currentMember.Status != "active" || !canManageOrganizationMembers(currentMember.Role) {
		respondAPIError(w, http.StatusForbidden, "organization_forbidden", "Only team admins and owners can manage members.")
		return
	}
	target, err := h.store.GetOrganizationMember(targetID)
	if err != nil || target.OrganizationID != orgID {
		respondAPIError(w, http.StatusNotFound, "organization_member_not_found", "Team member not found.")
		return
	}

	var req UpdateOrganizationMemberRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondAPIError(w, http.StatusBadRequest, "invalid_request", "Invalid request body")
		return
	}

	if strings.TrimSpace(req.Role) != "" {
		if !validOrganizationRole(req.Role) {
			respondAPIError(w, http.StatusBadRequest, "invalid_role", "Role must be one of read, edit, admin, or owner.")
			return
		}
		if req.Role == "owner" && currentMember.Role != "owner" {
			respondAPIError(w, http.StatusForbidden, "organization_forbidden", "Only the current owner can assign ownership.")
			return
		}
		if target.Role == "owner" && currentMember.Role != "owner" {
			respondAPIError(w, http.StatusForbidden, "organization_forbidden", "Only the current owner can change the owner role.")
			return
		}
		target.Role = req.Role
	}
	if strings.TrimSpace(req.Status) != "" {
		if !validOrganizationMemberStatus(req.Status) {
			respondAPIError(w, http.StatusBadRequest, "invalid_status", "Invalid team membership status.")
			return
		}
		if target.Role == "owner" && currentMember.Role != "owner" && req.Status != target.Status {
			respondAPIError(w, http.StatusForbidden, "organization_forbidden", "Only the current owner can change the owner membership.")
			return
		}
		target.Status = req.Status
		now := time.Now()
		if req.Status == "active" && target.JoinedAt.IsZero() {
			target.JoinedAt = now
		}
		if req.Status == "removed" {
			target.RemovedAt = now
		}
	}
	if err := h.store.UpdateOrganizationMember(target); err != nil {
		respondAPIError(w, http.StatusInternalServerError, "organization_member_update_failed", err.Error())
		return
	}
	workspaceID := ""
	if workspace, err := h.store.GetWorkspaceForOrganization(orgID); err == nil {
		workspaceID = workspace.ID
	}
	if err := h.reconcileOrganizationSeatBilling(r.Context(), orgID, workspaceID); err != nil {
		respondAPIError(w, http.StatusInternalServerError, "organization_billing_failed", err.Error())
		return
	}
	respondJSON(w, http.StatusOK, map[string]any{"member": target})
}

func (h *APIHandler) DeleteOrganizationMember(w http.ResponseWriter, r *http.Request) {
	orgID := chi.URLParam(r, "orgId")
	targetID := chi.URLParam(r, "memberId")
	_, currentMember, _, _, err := h.organizationAccessForRequest(r, orgID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondAPIError(w, http.StatusNotFound, "organization_not_found", "Organization not found.")
			return
		}
		respondAPIError(w, http.StatusInternalServerError, "organization_access_failed", err.Error())
		return
	}
	if currentMember.Status != "active" || !canManageOrganizationMembers(currentMember.Role) {
		respondAPIError(w, http.StatusForbidden, "organization_forbidden", "Only team admins and owners can manage members.")
		return
	}
	target, err := h.store.GetOrganizationMember(targetID)
	if err != nil || target.OrganizationID != orgID {
		respondAPIError(w, http.StatusNotFound, "organization_member_not_found", "Team member not found.")
		return
	}
	if target.Role == "owner" && currentMember.Role != "owner" {
		respondAPIError(w, http.StatusForbidden, "organization_forbidden", "Only the current owner can remove the owner.")
		return
	}
	target.Status = "removed"
	target.RemovedAt = time.Now()
	if err := h.store.UpdateOrganizationMember(target); err != nil {
		respondAPIError(w, http.StatusInternalServerError, "organization_member_update_failed", err.Error())
		return
	}
	workspaceID := ""
	if workspace, err := h.store.GetWorkspaceForOrganization(orgID); err == nil {
		workspaceID = workspace.ID
	}
	if err := h.reconcileOrganizationSeatBilling(r.Context(), orgID, workspaceID); err != nil {
		respondAPIError(w, http.StatusInternalServerError, "organization_billing_failed", err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *APIHandler) JoinOrganization(w http.ResponseWriter, r *http.Request) {
	session := h.sessionFromRequest(r)
	if session == nil || session.UserID == "" {
		respondAPIError(w, http.StatusUnauthorized, "auth_required", "You must be signed in to join a team.")
		return
	}

	var req JoinOrganizationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondAPIError(w, http.StatusBadRequest, "invalid_request", "Invalid request body")
		return
	}
	member, err := h.store.GetOrganizationMemberByInviteToken(strings.TrimSpace(req.Token))
	if err != nil {
		respondAPIError(w, http.StatusNotFound, "organization_invite_not_found", "Invite not found.")
		return
	}
	if member.Status != "invited" {
		respondAPIError(w, http.StatusBadRequest, "organization_invite_invalid", "That invite is no longer active.")
		return
	}
	if !member.InviteExpiresAt.IsZero() && member.InviteExpiresAt.Before(time.Now()) {
		member.Status = "expired"
		_ = h.store.UpdateOrganizationMember(member)
		respondAPIError(w, http.StatusBadRequest, "organization_invite_expired", "That invite has expired.")
		return
	}
	user, err := h.store.GetUserByID(session.UserID)
	if err != nil {
		respondAPIError(w, http.StatusInternalServerError, "user_lookup_failed", err.Error())
		return
	}
	if !strings.EqualFold(member.Email, user.Email) {
		respondAPIError(w, http.StatusForbidden, "organization_invite_email_mismatch", "Sign in with the invited email address to join this team.")
		return
	}

	orgWorkspace, err := h.store.GetWorkspaceForOrganization(member.OrganizationID)
	if err != nil {
		respondAPIError(w, http.StatusInternalServerError, "organization_workspace_failed", "Team workspace not found.")
		return
	}
	member.UserID = user.ID
	member.Status = "active"
	member.JoinedAt = time.Now()
	member.InviteToken = ""
	member.InviteExpiresAt = time.Time{}
	if err := h.store.UpdateOrganizationMember(member); err != nil {
		respondAPIError(w, http.StatusInternalServerError, "organization_join_failed", err.Error())
		return
	}
	if err := h.reconcileOrganizationSeatBilling(r.Context(), member.OrganizationID, orgWorkspace.ID); err != nil {
		respondAPIError(w, http.StatusInternalServerError, "organization_billing_failed", err.Error())
		return
	}
	if err := h.store.UpdateSessionWorkspace(session.ID, orgWorkspace.ID); err != nil {
		respondAPIError(w, http.StatusInternalServerError, "session_workspace_failed", err.Error())
		return
	}
	session.WorkspaceID = orgWorkspace.ID
	detail, err := h.buildOrganizationDetail(member.OrganizationID, session.UserID)
	if err != nil {
		respondAPIError(w, http.StatusInternalServerError, "organization_detail_failed", err.Error())
		return
	}
	respondJSON(w, http.StatusOK, detail)
}

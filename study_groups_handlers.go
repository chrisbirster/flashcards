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

func validStudyGroupRole(role string) bool {
	switch role {
	case "owner", "admin", "edit", "read":
		return true
	default:
		return false
	}
}

func normalizeStudyGroupRole(role string) string {
	switch strings.TrimSpace(role) {
	case "member":
		return "read"
	default:
		return strings.TrimSpace(role)
	}
}

func validStudyGroupMemberStatus(status string) bool {
	switch status {
	case "invited", "active", "removed", "declined", "expired":
		return true
	default:
		return false
	}
}

func canManageStudyGroupMembers(role string) bool {
	return role == "owner" || role == "admin"
}

func canEditStudyGroupSource(role string) bool {
	return role == "owner" || role == "admin" || role == "edit"
}

func (h *APIHandler) reloadCollectionSnapshot(collectionID string) error {
	if strings.TrimSpace(collectionID) == "" || collectionID != h.collectionID {
		return nil
	}
	col, err := h.store.GetCollection(collectionID)
	if err != nil {
		return err
	}
	h.collection = col
	return nil
}

func (h *APIHandler) requireStudyGroupsEntitlement(w http.ResponseWriter, r *http.Request) (*SessionRecord, Plan, bool) {
	session := h.sessionFromRequest(r)
	plan := h.planForRequest(r, session)
	usage := h.usageForSession(session)
	if !entitlementsForPlan(plan, usage).Features.StudyGroups {
		respondAPIError(w, http.StatusForbidden, "study_groups_not_available", "Study group creation and management are reserved for Team and Enterprise workspaces.")
		return nil, "", false
	}
	return session, plan, true
}

func (h *APIHandler) loadStudyGroupAccess(r *http.Request, groupID string) (*StudyGroup, *StudyGroupMember, *User, *SessionRecord, error) {
	session := h.sessionFromRequest(r)
	if session == nil || session.UserID == "" {
		return nil, nil, nil, nil, sql.ErrNoRows
	}
	user, err := h.store.GetUserByID(session.UserID)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	group, err := h.store.GetStudyGroup(groupID)
	if err != nil {
		return nil, nil, nil, session, err
	}
	member, err := h.store.getStudyGroupMembership(groupID, user.ID, user.Email)
	if err != nil {
		return nil, nil, user, session, err
	}
	return group, member, user, session, nil
}

func (h *APIHandler) installStudyGroupVersion(group *StudyGroup, member *StudyGroupMember, version *StudyGroupVersion, destinationWorkspace *Workspace) (*StudyGroupInstall, error) {
	sourceDeck, err := h.store.GetDeck(version.SourceDeckID)
	if err != nil {
		return nil, err
	}
	installedName := fmt.Sprintf("%s (Group v%d)", sourceDeck.Name, version.VersionNumber)
	newDeck, err := h.store.CopyDeckToCollection(version.SourceDeckID, destinationWorkspace.CollectionID, installedName)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	install := &StudyGroupInstall{
		ID:                     newID("sgi"),
		StudyGroupID:           group.ID,
		StudyGroupMemberID:     member.ID,
		DestinationWorkspaceID: destinationWorkspace.ID,
		InstalledDeckID:        newDeck.ID,
		InstalledDeckName:      newDeck.Name,
		SourceVersionNumber:    version.VersionNumber,
		Status:                 "active",
		SyncState:              "clean",
		CreatedAt:              now,
		UpdatedAt:              now,
	}
	if err := h.store.CreateStudyGroupInstall(install); err != nil {
		return nil, err
	}
	if err := h.reloadCollectionSnapshot(destinationWorkspace.CollectionID); err != nil {
		return nil, err
	}
	return h.store.GetStudyGroupInstall(install.ID)
}

func (h *APIHandler) ListStudyGroups(w http.ResponseWriter, r *http.Request) {
	session := h.sessionFromRequest(r)
	summaries, err := h.store.ListStudyGroupSummariesForUser(session.UserID)
	if err != nil {
		respondAPIError(w, http.StatusInternalServerError, "study_groups_list_failed", err.Error())
		return
	}
	respondJSON(w, http.StatusOK, summaries)
}

func (h *APIHandler) CreateStudyGroup(w http.ResponseWriter, r *http.Request) {
	if !h.requireWorkspaceWritePermission(w, r) {
		return
	}
	session, _, ok := h.requireStudyGroupsEntitlement(w, r)
	if !ok {
		return
	}

	var req CreateStudyGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondAPIError(w, http.StatusBadRequest, "invalid_request", "Invalid request body")
		return
	}
	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" || req.PrimaryDeckID == 0 {
		respondAPIError(w, http.StatusBadRequest, "invalid_study_group_request", "Name and primaryDeckId are required.")
		return
	}

	workspace, err := h.store.GetWorkspaceRecord(session.WorkspaceID)
	if err != nil {
		respondAPIError(w, http.StatusBadRequest, "workspace_not_found", "Workspace not found.")
		return
	}
	sourceCollectionID, err := h.store.GetDeckCollectionID(req.PrimaryDeckID)
	if err != nil {
		respondAPIError(w, http.StatusBadRequest, "source_deck_not_found", "Primary deck not found.")
		return
	}
	if sourceCollectionID != workspace.CollectionID {
		respondAPIError(w, http.StatusBadRequest, "invalid_source_deck", "Primary deck must belong to the current workspace.")
		return
	}

	now := time.Now()
	group := &StudyGroup{
		ID:              newID("sg"),
		WorkspaceID:     workspace.ID,
		PrimaryDeckID:   req.PrimaryDeckID,
		Name:            sanitizeHTML(req.Name),
		Description:     sanitizeHTML(strings.TrimSpace(req.Description)),
		Visibility:      strings.TrimSpace(req.Visibility),
		JoinPolicy:      strings.TrimSpace(req.JoinPolicy),
		CreatedByUserID: session.UserID,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
	if group.Visibility == "" {
		group.Visibility = "private"
	}
	if group.JoinPolicy == "" {
		group.JoinPolicy = "invite"
	}
	if err := h.store.CreateStudyGroup(group); err != nil {
		respondAPIError(w, http.StatusInternalServerError, "study_group_create_failed", err.Error())
		return
	}

	user, err := h.store.GetUserByID(session.UserID)
	if err != nil {
		respondAPIError(w, http.StatusInternalServerError, "study_group_owner_failed", err.Error())
		return
	}
	member := &StudyGroupMember{
		ID:           newID("sgm"),
		StudyGroupID: group.ID,
		UserID:       user.ID,
		Email:        user.Email,
		Role:         "owner",
		Status:       "active",
		JoinedAt:     now,
		CreatedAt:    now,
	}
	if err := h.store.CreateStudyGroupMember(member); err != nil {
		respondAPIError(w, http.StatusInternalServerError, "study_group_owner_failed", err.Error())
		return
	}
	_ = h.store.CreateStudyGroupEvent(&StudyGroupEvent{
		ID:           newID("sge"),
		StudyGroupID: group.ID,
		ActorUserID:  session.UserID,
		EventType:    "group_created",
		Payload:      h.store.encodeStudyGroupEventPayload(map[string]any{"primaryDeckId": req.PrimaryDeckID}),
		CreatedAt:    now,
	})

	detail, err := h.store.BuildStudyGroupDetail(group.ID, session.UserID)
	if err != nil {
		respondAPIError(w, http.StatusInternalServerError, "study_group_detail_failed", err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, detail)
}

func (h *APIHandler) GetStudyGroup(w http.ResponseWriter, r *http.Request) {
	groupID := chi.URLParam(r, "id")
	session := h.sessionFromRequest(r)
	detail, err := h.store.BuildStudyGroupDetail(groupID, session.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondAPIError(w, http.StatusNotFound, "study_group_not_found", "Study group not found.")
			return
		}
		respondAPIError(w, http.StatusInternalServerError, "study_group_detail_failed", err.Error())
		return
	}
	respondJSON(w, http.StatusOK, detail)
}

func (h *APIHandler) UpdateStudyGroup(w http.ResponseWriter, r *http.Request) {
	groupID := chi.URLParam(r, "id")
	group, member, _, _, err := h.loadStudyGroupAccess(r, groupID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondAPIError(w, http.StatusNotFound, "study_group_not_found", "Study group not found.")
			return
		}
		respondAPIError(w, http.StatusInternalServerError, "study_group_access_failed", err.Error())
		return
	}
	if member.Status != "active" || !canManageStudyGroupMembers(member.Role) {
		respondAPIError(w, http.StatusForbidden, "study_group_forbidden", "Only owners and admins can update this study group.")
		return
	}

	var req UpdateStudyGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondAPIError(w, http.StatusBadRequest, "invalid_request", "Invalid request body")
		return
	}

	if trimmed := strings.TrimSpace(req.Name); trimmed != "" {
		group.Name = sanitizeHTML(trimmed)
	}
	group.Description = sanitizeHTML(strings.TrimSpace(req.Description))
	if visibility := strings.TrimSpace(req.Visibility); visibility != "" {
		group.Visibility = visibility
	}
	if joinPolicy := strings.TrimSpace(req.JoinPolicy); joinPolicy != "" {
		group.JoinPolicy = joinPolicy
	}
	group.UpdatedAt = time.Now()
	if err := h.store.UpdateStudyGroup(group); err != nil {
		respondAPIError(w, http.StatusInternalServerError, "study_group_update_failed", err.Error())
		return
	}
	_ = h.store.CreateStudyGroupEvent(&StudyGroupEvent{
		ID:           newID("sge"),
		StudyGroupID: group.ID,
		ActorUserID:  member.UserID,
		EventType:    "group_updated",
		Payload:      h.store.encodeStudyGroupEventPayload(map[string]any{"visibility": group.Visibility, "joinPolicy": group.JoinPolicy}),
		CreatedAt:    time.Now(),
	})
	detail, err := h.store.BuildStudyGroupDetail(group.ID, member.UserID)
	if err != nil {
		respondAPIError(w, http.StatusInternalServerError, "study_group_detail_failed", err.Error())
		return
	}
	respondJSON(w, http.StatusOK, detail)
}

func (h *APIHandler) DeleteStudyGroup(w http.ResponseWriter, r *http.Request) {
	groupID := chi.URLParam(r, "id")
	group, member, _, _, err := h.loadStudyGroupAccess(r, groupID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondAPIError(w, http.StatusNotFound, "study_group_not_found", "Study group not found.")
			return
		}
		respondAPIError(w, http.StatusInternalServerError, "study_group_access_failed", err.Error())
		return
	}
	if member.Status != "active" || !canManageStudyGroupMembers(member.Role) {
		respondAPIError(w, http.StatusForbidden, "study_group_forbidden", "Only owners and admins can delete this study group.")
		return
	}
	if err := h.store.DeleteStudyGroup(group.ID); err != nil {
		respondAPIError(w, http.StatusInternalServerError, "study_group_delete_failed", err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *APIHandler) InviteStudyGroupMember(w http.ResponseWriter, r *http.Request) {
	groupID := chi.URLParam(r, "id")
	group, member, _, _, err := h.loadStudyGroupAccess(r, groupID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondAPIError(w, http.StatusNotFound, "study_group_not_found", "Study group not found.")
			return
		}
		respondAPIError(w, http.StatusInternalServerError, "study_group_access_failed", err.Error())
		return
	}
	if member.Status != "active" || !canManageStudyGroupMembers(member.Role) {
		respondAPIError(w, http.StatusForbidden, "study_group_forbidden", "Only owners and admins can invite members.")
		return
	}

	var req InviteStudyGroupMemberRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondAPIError(w, http.StatusBadRequest, "invalid_request", "Invalid request body")
		return
	}
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	req.Role = normalizeStudyGroupRole(req.Role)
	if req.Role == "" {
		req.Role = "read"
	}
	if req.Email == "" || !validStudyGroupRole(req.Role) || req.Role == "owner" {
		respondAPIError(w, http.StatusBadRequest, "invalid_member_request", "Invite requires an email and role of admin, edit, or read.")
		return
	}

	now := time.Now()
	invite := &StudyGroupMember{
		ID:              newID("sgm"),
		StudyGroupID:    group.ID,
		Email:           req.Email,
		Role:            req.Role,
		Status:          "invited",
		InviteToken:     newID("invite"),
		InviteExpiresAt: now.Add(7 * 24 * time.Hour),
		CreatedAt:       now,
	}
	if err := h.store.UpsertStudyGroupInvitation(invite); err != nil {
		respondAPIError(w, http.StatusInternalServerError, "study_group_invite_failed", err.Error())
		return
	}
	savedInvite, err := h.store.GetStudyGroupMemberByGroupAndEmail(group.ID, req.Email)
	if err != nil {
		respondAPIError(w, http.StatusInternalServerError, "study_group_invite_failed", err.Error())
		return
	}
	_ = h.store.CreateStudyGroupEvent(&StudyGroupEvent{
		ID:           newID("sge"),
		StudyGroupID: group.ID,
		ActorUserID:  member.UserID,
		EventType:    "member_invited",
		Payload:      h.store.encodeStudyGroupEventPayload(map[string]any{"email": savedInvite.Email, "role": savedInvite.Role}),
		CreatedAt:    now,
	})
	respondJSON(w, http.StatusCreated, savedInvite)
}

func (h *APIHandler) UpdateStudyGroupMember(w http.ResponseWriter, r *http.Request) {
	groupID := chi.URLParam(r, "id")
	targetMemberID := chi.URLParam(r, "memberId")
	group, member, _, _, err := h.loadStudyGroupAccess(r, groupID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondAPIError(w, http.StatusNotFound, "study_group_not_found", "Study group not found.")
			return
		}
		respondAPIError(w, http.StatusInternalServerError, "study_group_access_failed", err.Error())
		return
	}
	if member.Status != "active" || !canManageStudyGroupMembers(member.Role) {
		respondAPIError(w, http.StatusForbidden, "study_group_forbidden", "Only owners and admins can update members.")
		return
	}

	targetMember, err := h.store.GetStudyGroupMember(targetMemberID)
	if err != nil || targetMember.StudyGroupID != group.ID {
		respondAPIError(w, http.StatusNotFound, "study_group_member_not_found", "Study group member not found.")
		return
	}
	if targetMember.Role == "owner" {
		respondAPIError(w, http.StatusForbidden, "study_group_owner_locked", "Owner membership cannot be edited from this endpoint.")
		return
	}

	var req UpdateStudyGroupMemberRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondAPIError(w, http.StatusBadRequest, "invalid_request", "Invalid request body")
		return
	}
	if role := normalizeStudyGroupRole(req.Role); role != "" {
		if !validStudyGroupRole(role) || role == "owner" {
			respondAPIError(w, http.StatusBadRequest, "invalid_member_role", "Role must be admin, edit, or read.")
			return
		}
		targetMember.Role = role
	}
	if status := strings.TrimSpace(req.Status); status != "" {
		if !validStudyGroupMemberStatus(status) {
			respondAPIError(w, http.StatusBadRequest, "invalid_member_status", "Invalid study group member status.")
			return
		}
		targetMember.Status = status
		if status == "active" && targetMember.JoinedAt.IsZero() {
			targetMember.JoinedAt = time.Now()
		}
		if status == "removed" {
			targetMember.RemovedAt = time.Now()
			targetMember.InviteToken = ""
		}
	}
	if err := h.store.UpdateStudyGroupMember(targetMember); err != nil {
		respondAPIError(w, http.StatusInternalServerError, "study_group_member_update_failed", err.Error())
		return
	}
	_ = h.store.CreateStudyGroupEvent(&StudyGroupEvent{
		ID:           newID("sge"),
		StudyGroupID: group.ID,
		ActorUserID:  member.UserID,
		EventType:    "member_updated",
		Payload:      h.store.encodeStudyGroupEventPayload(map[string]any{"memberId": targetMember.ID, "role": targetMember.Role, "status": targetMember.Status}),
		CreatedAt:    time.Now(),
	})
	respondJSON(w, http.StatusOK, targetMember)
}

func (h *APIHandler) DeleteStudyGroupMember(w http.ResponseWriter, r *http.Request) {
	groupID := chi.URLParam(r, "id")
	targetMemberID := chi.URLParam(r, "memberId")
	group, member, _, _, err := h.loadStudyGroupAccess(r, groupID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondAPIError(w, http.StatusNotFound, "study_group_not_found", "Study group not found.")
			return
		}
		respondAPIError(w, http.StatusInternalServerError, "study_group_access_failed", err.Error())
		return
	}
	if member.Status != "active" || !canManageStudyGroupMembers(member.Role) {
		respondAPIError(w, http.StatusForbidden, "study_group_forbidden", "Only owners and admins can remove members.")
		return
	}

	targetMember, err := h.store.GetStudyGroupMember(targetMemberID)
	if err != nil || targetMember.StudyGroupID != group.ID {
		respondAPIError(w, http.StatusNotFound, "study_group_member_not_found", "Study group member not found.")
		return
	}
	if targetMember.Role == "owner" {
		respondAPIError(w, http.StatusForbidden, "study_group_owner_locked", "Owner membership cannot be removed.")
		return
	}
	targetMember.Status = "removed"
	targetMember.RemovedAt = time.Now()
	targetMember.InviteToken = ""
	targetMember.InviteExpiresAt = time.Time{}
	if err := h.store.UpdateStudyGroupMember(targetMember); err != nil {
		respondAPIError(w, http.StatusInternalServerError, "study_group_member_delete_failed", err.Error())
		return
	}
	_ = h.store.CreateStudyGroupEvent(&StudyGroupEvent{
		ID:           newID("sge"),
		StudyGroupID: group.ID,
		ActorUserID:  member.UserID,
		EventType:    "member_removed",
		Payload:      h.store.encodeStudyGroupEventPayload(map[string]any{"memberId": targetMember.ID, "email": targetMember.Email}),
		CreatedAt:    time.Now(),
	})
	w.WriteHeader(http.StatusNoContent)
}

func (h *APIHandler) JoinStudyGroup(w http.ResponseWriter, r *http.Request) {
	session := h.sessionFromRequest(r)
	user, err := h.store.GetUserByID(session.UserID)
	if err != nil {
		respondAPIError(w, http.StatusInternalServerError, "user_not_found", err.Error())
		return
	}

	var req JoinStudyGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondAPIError(w, http.StatusBadRequest, "invalid_request", "Invalid request body")
		return
	}
	req.Token = strings.TrimSpace(req.Token)
	if req.Token == "" {
		respondAPIError(w, http.StatusBadRequest, "invalid_invite_token", "Invite token is required.")
		return
	}
	member, err := h.store.GetStudyGroupMemberByInviteToken(req.Token)
	if err != nil {
		respondAPIError(w, http.StatusNotFound, "study_group_invite_not_found", "Invite not found.")
		return
	}
	if member.Status != "invited" {
		respondAPIError(w, http.StatusConflict, "study_group_invite_used", "This invite is no longer active.")
		return
	}
	if !member.InviteExpiresAt.IsZero() && member.InviteExpiresAt.Before(time.Now()) {
		member.Status = "expired"
		_ = h.store.UpdateStudyGroupMember(member)
		respondAPIError(w, http.StatusConflict, "study_group_invite_expired", "This invite has expired.")
		return
	}
	if !strings.EqualFold(member.Email, user.Email) {
		respondAPIError(w, http.StatusForbidden, "study_group_invite_email_mismatch", "This invite was issued for a different email address.")
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

	group, err := h.store.GetStudyGroup(member.StudyGroupID)
	if err != nil {
		respondAPIError(w, http.StatusNotFound, "study_group_not_found", "Study group not found.")
		return
	}

	member.UserID = user.ID
	member.Status = "active"
	member.JoinedAt = time.Now()
	member.InviteToken = ""
	member.InviteExpiresAt = time.Time{}
	if err := h.store.UpdateStudyGroupMember(member); err != nil {
		respondAPIError(w, http.StatusInternalServerError, "study_group_join_failed", err.Error())
		return
	}

	if req.InstallLatest {
		if latestVersion, err := h.store.GetLatestStudyGroupVersion(group.ID); err == nil {
			if _, err := h.installStudyGroupVersion(group, member, latestVersion, workspace); err != nil {
				respondAPIError(w, http.StatusInternalServerError, "study_group_install_failed", err.Error())
				return
			}
		}
	}

	_ = h.store.CreateStudyGroupEvent(&StudyGroupEvent{
		ID:           newID("sge"),
		StudyGroupID: group.ID,
		ActorUserID:  user.ID,
		EventType:    "member_joined",
		Payload:      h.store.encodeStudyGroupEventPayload(map[string]any{"workspaceId": workspace.ID, "installLatest": req.InstallLatest}),
		CreatedAt:    time.Now(),
	})
	detail, err := h.store.BuildStudyGroupDetail(group.ID, user.ID)
	if err != nil {
		respondAPIError(w, http.StatusInternalServerError, "study_group_detail_failed", err.Error())
		return
	}
	respondJSON(w, http.StatusOK, detail)
}

func (h *APIHandler) ListStudyGroupVersions(w http.ResponseWriter, r *http.Request) {
	groupID := chi.URLParam(r, "id")
	session := h.sessionFromRequest(r)
	detail, err := h.store.BuildStudyGroupDetail(groupID, session.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondAPIError(w, http.StatusNotFound, "study_group_not_found", "Study group not found.")
			return
		}
		respondAPIError(w, http.StatusInternalServerError, "study_group_detail_failed", err.Error())
		return
	}
	respondJSON(w, http.StatusOK, detail.Versions)
}

func (h *APIHandler) PublishStudyGroupVersion(w http.ResponseWriter, r *http.Request) {
	groupID := chi.URLParam(r, "id")
	group, member, _, _, err := h.loadStudyGroupAccess(r, groupID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondAPIError(w, http.StatusNotFound, "study_group_not_found", "Study group not found.")
			return
		}
		respondAPIError(w, http.StatusInternalServerError, "study_group_access_failed", err.Error())
		return
	}
	if member.Status != "active" || !canManageStudyGroupMembers(member.Role) {
		respondAPIError(w, http.StatusForbidden, "study_group_forbidden", "Only owners and admins can publish versions.")
		return
	}

	var req PublishStudyGroupVersionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil && !errors.Is(err, http.ErrBodyNotAllowed) {
		respondAPIError(w, http.StatusBadRequest, "invalid_request", "Invalid request body")
		return
	}

	nextVersion := 1
	if latestVersion, err := h.store.GetLatestStudyGroupVersion(group.ID); err == nil {
		nextVersion = latestVersion.VersionNumber + 1
	}
	noteCount, cardCount, err := h.store.GetDeckContentSummary(group.PrimaryDeckID)
	if err != nil {
		respondAPIError(w, http.StatusInternalServerError, "study_group_version_failed", err.Error())
		return
	}
	version := &StudyGroupVersion{
		ID:                newID("sgv"),
		StudyGroupID:      group.ID,
		VersionNumber:     nextVersion,
		SourceDeckID:      group.PrimaryDeckID,
		PublishedByUserID: member.UserID,
		ChangeSummary:     strings.TrimSpace(req.ChangeSummary),
		NoteCount:         noteCount,
		CardCount:         cardCount,
		CreatedAt:         time.Now(),
	}
	if err := h.store.CreateStudyGroupVersion(version); err != nil {
		respondAPIError(w, http.StatusInternalServerError, "study_group_version_failed", err.Error())
		return
	}
	_ = h.store.CreateStudyGroupEvent(&StudyGroupEvent{
		ID:           newID("sge"),
		StudyGroupID: group.ID,
		ActorUserID:  member.UserID,
		EventType:    "version_published",
		Payload:      h.store.encodeStudyGroupEventPayload(map[string]any{"versionNumber": version.VersionNumber, "changeSummary": version.ChangeSummary}),
		CreatedAt:    time.Now(),
	})
	respondJSON(w, http.StatusCreated, version)
}

func (h *APIHandler) InstallStudyGroupDeck(w http.ResponseWriter, r *http.Request) {
	groupID := chi.URLParam(r, "id")
	group, member, user, session, err := h.loadStudyGroupAccess(r, groupID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondAPIError(w, http.StatusNotFound, "study_group_not_found", "Study group not found.")
			return
		}
		respondAPIError(w, http.StatusInternalServerError, "study_group_access_failed", err.Error())
		return
	}
	if member.Status != "active" {
		respondAPIError(w, http.StatusForbidden, "study_group_membership_inactive", "Only active members can install group decks.")
		return
	}

	var req InstallStudyGroupDeckRequest
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
	latestVersion, err := h.store.GetLatestStudyGroupVersion(group.ID)
	if err != nil {
		respondAPIError(w, http.StatusConflict, "study_group_no_published_version", "Publish a source version before members can install it.")
		return
	}
	if currentInstall, err := h.store.GetCurrentStudyGroupInstall(group.ID, member.ID); err == nil {
		if currentInstall.SourceVersionNumber == latestVersion.VersionNumber {
			respondJSON(w, http.StatusOK, currentInstall)
			return
		}
		respondAPIError(w, http.StatusConflict, "study_group_update_required", "Use the install update flow to move to a newer version.")
		return
	}
	install, err := h.installStudyGroupVersion(group, member, latestVersion, workspace)
	if err != nil {
		respondAPIError(w, http.StatusInternalServerError, "study_group_install_failed", err.Error())
		return
	}
	_ = h.store.CreateStudyGroupEvent(&StudyGroupEvent{
		ID:           newID("sge"),
		StudyGroupID: group.ID,
		ActorUserID:  user.ID,
		EventType:    "install_created",
		Payload:      h.store.encodeStudyGroupEventPayload(map[string]any{"installId": install.ID, "workspaceId": workspace.ID, "versionNumber": latestVersion.VersionNumber}),
		CreatedAt:    time.Now(),
	})
	respondJSON(w, http.StatusCreated, install)
}

func (h *APIHandler) UpdateStudyGroupInstall(w http.ResponseWriter, r *http.Request) {
	groupID := chi.URLParam(r, "id")
	installID := chi.URLParam(r, "installId")
	group, member, user, _, err := h.loadStudyGroupAccess(r, groupID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondAPIError(w, http.StatusNotFound, "study_group_not_found", "Study group not found.")
			return
		}
		respondAPIError(w, http.StatusInternalServerError, "study_group_access_failed", err.Error())
		return
	}
	if member.Status != "active" {
		respondAPIError(w, http.StatusForbidden, "study_group_membership_inactive", "Only active members can update installs.")
		return
	}
	currentInstall, err := h.store.GetStudyGroupInstall(installID)
	if err != nil || currentInstall.StudyGroupID != group.ID || currentInstall.StudyGroupMemberID != member.ID {
		respondAPIError(w, http.StatusNotFound, "study_group_install_not_found", "Study group install not found.")
		return
	}
	latestVersion, err := h.store.GetLatestStudyGroupVersion(group.ID)
	if err != nil {
		respondAPIError(w, http.StatusConflict, "study_group_no_published_version", "No published version is available.")
		return
	}
	if latestVersion.VersionNumber <= currentInstall.SourceVersionNumber {
		respondAPIError(w, http.StatusConflict, "study_group_already_current", "This install already uses the latest published version.")
		return
	}

	var req UpdateStudyGroupInstallRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil && !errors.Is(err, http.ErrBodyNotAllowed) {
		respondAPIError(w, http.StatusBadRequest, "invalid_request", "Invalid request body")
		return
	}
	destinationWorkspaceID := currentInstall.DestinationWorkspaceID
	if strings.TrimSpace(req.DestinationWorkspaceID) != "" {
		destinationWorkspaceID = strings.TrimSpace(req.DestinationWorkspaceID)
	}
	workspace, err := h.store.GetWorkspaceForUser(user.ID, destinationWorkspaceID)
	if err != nil {
		respondAPIError(w, http.StatusBadRequest, "invalid_destination_workspace", "Destination workspace not found.")
		return
	}

	nextInstall, err := h.installStudyGroupVersion(group, member, latestVersion, workspace)
	if err != nil {
		respondAPIError(w, http.StatusInternalServerError, "study_group_install_update_failed", err.Error())
		return
	}
	currentInstall.Status = "superseded"
	currentInstall.SupersededByInstallID = nextInstall.ID
	currentInstall.UpdatedAt = time.Now()
	if err := h.store.UpdateStudyGroupInstall(currentInstall); err != nil {
		respondAPIError(w, http.StatusInternalServerError, "study_group_install_update_failed", err.Error())
		return
	}
	_ = h.store.CreateStudyGroupEvent(&StudyGroupEvent{
		ID:           newID("sge"),
		StudyGroupID: group.ID,
		ActorUserID:  user.ID,
		EventType:    "install_updated",
		Payload:      h.store.encodeStudyGroupEventPayload(map[string]any{"oldInstallId": currentInstall.ID, "newInstallId": nextInstall.ID, "versionNumber": latestVersion.VersionNumber}),
		CreatedAt:    time.Now(),
	})
	respondJSON(w, http.StatusOK, nextInstall)
}

func (h *APIHandler) RemoveStudyGroupInstall(w http.ResponseWriter, r *http.Request) {
	groupID := chi.URLParam(r, "id")
	installID := chi.URLParam(r, "installId")
	group, member, _, _, err := h.loadStudyGroupAccess(r, groupID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondAPIError(w, http.StatusNotFound, "study_group_not_found", "Study group not found.")
			return
		}
		respondAPIError(w, http.StatusInternalServerError, "study_group_access_failed", err.Error())
		return
	}
	install, err := h.store.GetStudyGroupInstall(installID)
	if err != nil || install.StudyGroupID != group.ID {
		respondAPIError(w, http.StatusNotFound, "study_group_install_not_found", "Study group install not found.")
		return
	}
	if install.StudyGroupMemberID != member.ID && !canManageStudyGroupMembers(member.Role) {
		respondAPIError(w, http.StatusForbidden, "study_group_forbidden", "You can only remove your own install.")
		return
	}
	if err := h.store.DeleteCopiedDeck(install.InstalledDeckID); err != nil {
		respondAPIError(w, http.StatusInternalServerError, "study_group_install_delete_failed", err.Error())
		return
	}
	install.Status = "removed"
	install.InstalledDeckID = 0
	install.InstalledDeckName = ""
	install.UpdatedAt = time.Now()
	if err := h.store.UpdateStudyGroupInstall(install); err != nil {
		respondAPIError(w, http.StatusInternalServerError, "study_group_install_delete_failed", err.Error())
		return
	}
	destinationWorkspace, err := h.store.GetWorkspaceRecord(install.DestinationWorkspaceID)
	if err == nil {
		if reloadErr := h.reloadCollectionSnapshot(destinationWorkspace.CollectionID); reloadErr != nil {
			respondAPIError(w, http.StatusInternalServerError, "study_group_install_delete_failed", reloadErr.Error())
			return
		}
	}
	_ = h.store.CreateStudyGroupEvent(&StudyGroupEvent{
		ID:           newID("sge"),
		StudyGroupID: group.ID,
		ActorUserID:  member.UserID,
		EventType:    "install_removed",
		Payload:      h.store.encodeStudyGroupEventPayload(map[string]any{"installId": install.ID}),
		CreatedAt:    time.Now(),
	})
	w.WriteHeader(http.StatusNoContent)
}

func (h *APIHandler) GetStudyGroupDashboard(w http.ResponseWriter, r *http.Request) {
	groupID := chi.URLParam(r, "id")
	_, _, _, session, err := h.loadStudyGroupAccess(r, groupID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondAPIError(w, http.StatusNotFound, "study_group_not_found", "Study group not found.")
			return
		}
		respondAPIError(w, http.StatusInternalServerError, "study_group_access_failed", err.Error())
		return
	}
	dashboard, err := h.store.GetStudyGroupDashboard(groupID)
	if err != nil {
		respondAPIError(w, http.StatusInternalServerError, "study_group_dashboard_failed", err.Error())
		return
	}
	_ = session
	respondJSON(w, http.StatusOK, dashboard)
}

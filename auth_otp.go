package main

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/mail"
	"strings"
	"time"
)

type otpRequestBody struct {
	Email string `json:"email"`
}

type otpVerifyBody struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

func (h *APIHandler) RequestOTP(w http.ResponseWriter, r *http.Request) {
	var req otpRequestBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondAPIError(w, http.StatusBadRequest, "invalid_request", "Invalid OTP request body")
		return
	}

	email, err := normalizeEmail(req.Email)
	if err != nil {
		respondAPIError(w, http.StatusBadRequest, "invalid_email", err.Error())
		return
	}

	now := time.Now()
	if tooMany, err := h.tooManyOTPRequests(email, requestIP(r), now); err != nil {
		respondAPIError(w, http.StatusInternalServerError, "otp_rate_limit_failed", err.Error())
		return
	} else if tooMany {
		respondAPIError(w, http.StatusTooManyRequests, "otp_rate_limited", "Too many OTP requests. Please wait a few minutes and try again.")
		return
	}

	if latest, err := h.store.GetLatestOTPChallenge(email); err == nil {
		if latest.ConsumedAt.IsZero() && latest.ResendAvailableAt.After(now) {
			respondAPIError(w, http.StatusTooManyRequests, "otp_retry_later", fmt.Sprintf("Please wait %d seconds before requesting another code.", int(time.Until(latest.ResendAvailableAt).Seconds())+1))
			return
		}
	}

	if err := h.store.InvalidateOTPChallenges(email); err != nil {
		respondAPIError(w, http.StatusInternalServerError, "otp_invalidate_failed", err.Error())
		return
	}

	code, err := generateOTPCode()
	if err != nil {
		respondAPIError(w, http.StatusInternalServerError, "otp_generate_failed", "Failed to generate OTP code")
		return
	}

	challenge := &OTPChallenge{
		ID:                newID("otp"),
		Email:             email,
		ExpiresAt:         now.Add(h.config.OTP.TTL),
		AttemptCount:      0,
		MaxAttempts:       h.config.OTP.MaxAttempts,
		ResendAvailableAt: now.Add(h.config.OTP.ResendCooldown),
		RequestedIP:       requestIP(r),
		UserAgent:         strings.TrimSpace(r.UserAgent()),
		CreatedAt:         now,
	}
	challenge.CodeHash = hashOTPCode(h.config.SessionSecret, challenge.ID, code)

	if err := h.store.CreateOTPChallenge(challenge); err != nil {
		respondAPIError(w, http.StatusInternalServerError, "otp_store_failed", err.Error())
		return
	}

	if err := h.emailSender.SendOTP(context.Background(), email, code, challenge.ExpiresAt); err != nil {
		respondAPIError(w, http.StatusBadGateway, "otp_send_failed", err.Error())
		return
	}

	response := map[string]interface{}{
		"ok":                true,
		"expiresAt":         challenge.ExpiresAt.Format(time.RFC3339),
		"retryAfterSeconds": int(h.config.OTP.ResendCooldown.Seconds()),
		"delivery":          "email",
	}
	if h.config.IsDevelopment() {
		response["delivery"] = "dev-inline"
		response["devCode"] = code
	}

	respondJSON(w, http.StatusOK, response)
}

func (h *APIHandler) VerifyOTP(w http.ResponseWriter, r *http.Request) {
	var req otpVerifyBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondAPIError(w, http.StatusBadRequest, "invalid_request", "Invalid OTP verify body")
		return
	}

	email, err := normalizeEmail(req.Email)
	if err != nil {
		respondAPIError(w, http.StatusBadRequest, "invalid_email", err.Error())
		return
	}

	code := strings.TrimSpace(req.Code)
	if len(code) != 6 {
		respondAPIError(w, http.StatusBadRequest, "invalid_code", "OTP code must be 6 digits")
		return
	}

	challenge, err := h.store.GetLatestOTPChallenge(email)
	if err != nil {
		respondAPIError(w, http.StatusBadRequest, "otp_not_found", "No active code found for that email address")
		return
	}

	now := time.Now()
	if !challenge.ConsumedAt.IsZero() || challenge.ExpiresAt.Before(now) {
		respondAPIError(w, http.StatusBadRequest, "otp_expired", "That code has expired. Request a new one.")
		return
	}
	if challenge.AttemptCount >= challenge.MaxAttempts {
		respondAPIError(w, http.StatusTooManyRequests, "otp_locked", "Too many incorrect code attempts. Request a new one.")
		return
	}

	expectedHash := hashOTPCode(h.config.SessionSecret, challenge.ID, code)
	if subtle.ConstantTimeCompare([]byte(expectedHash), []byte(challenge.CodeHash)) != 1 {
		_ = h.store.IncrementOTPChallengeAttempts(challenge.ID)
		respondAPIError(w, http.StatusBadRequest, "invalid_code", "That code is incorrect")
		return
	}

	user, err := h.store.GetUserByEmail(email)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		respondAPIError(w, http.StatusInternalServerError, "user_lookup_failed", err.Error())
		return
	}

	if user == nil || errors.Is(err, sql.ErrNoRows) {
		user = &User{
			ID:          newID("usr"),
			Email:       email,
			DisplayName: displayNameForEmail(email),
			Onboarding:  true,
			LastLoginAt: now,
			CreatedAt:   now,
			UpdatedAt:   now,
		}
		if err := h.store.CreateUser(user); err != nil {
			respondAPIError(w, http.StatusInternalServerError, "user_create_failed", err.Error())
			return
		}
	} else {
		_ = h.store.UpdateUserLastLogin(user.ID, now)
		user.LastLoginAt = now
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
	if err := h.store.ConsumeOTPChallenge(challenge.ID, now); err != nil {
		respondAPIError(w, http.StatusInternalServerError, "otp_consume_failed", err.Error())
		return
	}

	h.writeCookie(w, sessionCookieName, session.ID, session.ExpiresAt)

	requestWithSession := r.WithContext(context.WithValue(r.Context(), sessionContextKey, session))
	respondJSON(w, http.StatusOK, h.buildSessionResponse(requestWithSession))
}

func normalizeEmail(raw string) (string, error) {
	email := strings.ToLower(strings.TrimSpace(raw))
	if email == "" {
		return "", fmt.Errorf("email is required")
	}
	if _, err := mail.ParseAddress(email); err != nil {
		return "", fmt.Errorf("enter a valid email address")
	}
	return email, nil
}

func displayNameForEmail(email string) string {
	localPart := email
	if idx := strings.Index(email, "@"); idx > 0 {
		localPart = email[:idx]
	}
	localPart = strings.ReplaceAll(localPart, ".", " ")
	localPart = strings.ReplaceAll(localPart, "_", " ")
	localPart = strings.TrimSpace(localPart)
	if localPart == "" {
		return "Vutadex User"
	}
	parts := strings.Fields(localPart)
	for i, part := range parts {
		if len(part) > 0 {
			parts[i] = strings.ToUpper(part[:1]) + strings.ToLower(part[1:])
		}
	}
	return strings.Join(parts, " ")
}

func generateOTPCode() (string, error) {
	bytes := make([]byte, 3)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	value := int(bytes[0])<<16 | int(bytes[1])<<8 | int(bytes[2])
	return fmt.Sprintf("%06d", value%1000000), nil
}

func hashOTPCode(secret, challengeID, code string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(challengeID))
	mac.Write([]byte(":"))
	mac.Write([]byte(code))
	return hex.EncodeToString(mac.Sum(nil))
}

func requestIP(r *http.Request) string {
	host, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr))
	if err == nil && host != "" {
		return host
	}
	return strings.TrimSpace(r.RemoteAddr)
}

func (h *APIHandler) tooManyOTPRequests(email, ip string, now time.Time) (bool, error) {
	emailCount, err := h.store.CountRecentOTPChallengesByEmail(email, now.Add(-15*time.Minute))
	if err != nil {
		return false, err
	}
	if emailCount >= 5 {
		return true, nil
	}

	ipCount, err := h.store.CountRecentOTPChallengesByIP(ip, now.Add(-15*time.Minute))
	if err != nil {
		return false, err
	}
	return ipCount >= 20, nil
}

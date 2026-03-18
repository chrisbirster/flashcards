package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type EmailSender interface {
	SendOTP(ctx context.Context, to, code string, expiresAt time.Time) error
}

type HTTPEmailSender struct {
	client *http.Client
	config EmailConfig
}

type LogEmailSender struct{}

func NewEmailSender(cfg AppConfig) EmailSender {
	if cfg.Email.SendURL == "" {
		return LogEmailSender{}
	}
	return &HTTPEmailSender{
		client: &http.Client{Timeout: 15 * time.Second},
		config: cfg.Email,
	}
}

func (s *HTTPEmailSender) SendOTP(ctx context.Context, to, code string, expiresAt time.Time) error {
	payload := map[string]string{
		"to":        to,
		"subject":   "Your Vutadex login code",
		"otpCode":   code,
		"expiresAt": expiresAt.Format(time.RFC3339),
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal email payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.config.SendURL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to create email request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if s.config.AuthHeaderName != "" && s.config.AuthHeaderValue != "" {
		req.Header.Set(s.config.AuthHeaderName, s.config.AuthHeaderValue)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send OTP email: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("email API returned %s", resp.Status)
	}

	return nil
}

func (LogEmailSender) SendOTP(_ context.Context, to, code string, expiresAt time.Time) error {
	log.Printf("OTP for %s: %s (expires %s)\n", to, code, expiresAt.Format(time.RFC3339))
	return nil
}

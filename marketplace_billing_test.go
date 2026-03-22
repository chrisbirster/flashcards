package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"testing"
	"time"
)

func TestVerifyStripeWebhookSignature(t *testing.T) {
	secret := "whsec_test_secret"
	payload := []byte(`{"id":"evt_123","type":"checkout.session.completed"}`)
	now := time.Unix(1_710_000_000, 0)
	timestamp := "1710000000"

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(timestamp))
	mac.Write([]byte("."))
	mac.Write(payload)
	signature := hex.EncodeToString(mac.Sum(nil))

	header := "t=" + timestamp + ",v1=" + signature
	if err := verifyStripeWebhookSignature(payload, header, secret, now); err != nil {
		t.Fatalf("expected signature verification to succeed, got %v", err)
	}

	if err := verifyStripeWebhookSignature(payload, header, "wrong_secret", now); err == nil {
		t.Fatalf("expected verification to fail with wrong secret")
	}
}

func TestAppendURLQueryPreservesCheckoutPlaceholder(t *testing.T) {
	got := appendURLQuery("https://app.vutadex.com/marketplace?checkout=success", map[string]string{
		"checkout_session_id": "{CHECKOUT_SESSION_ID}",
	})
	want := "https://app.vutadex.com/marketplace?checkout=success&checkout_session_id={CHECKOUT_SESSION_ID}"
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

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
	"net/url"
	"strconv"
	"strings"
	"time"
)

type marketplaceBillingProvider interface {
	ProviderName() string
	StartCreatorOnboarding(ctx context.Context, user *User, workspace *Workspace, existing *MarketplaceCreatorAccount) (*MarketplaceCreatorAccount, error)
	RefreshCreatorAccount(ctx context.Context, existing *MarketplaceCreatorAccount) (*MarketplaceCreatorAccount, error)
	CreateCheckoutSession(ctx context.Context, listing *MarketplaceListing, order *MarketplaceOrder, buyer *User, creatorAccount *MarketplaceCreatorAccount) (*MarketplaceCheckoutResponse, error)
	GetCheckoutSession(ctx context.Context, order *MarketplaceOrder, creatorAccount *MarketplaceCreatorAccount) (*MarketplaceCheckoutSessionState, error)
}

type devMarketplaceBillingProvider struct {
	cfg AppConfig
}

type MarketplaceCheckoutSessionState struct {
	Provider                string
	ProviderCheckoutURL     string
	ProviderCheckoutID      string
	ProviderPaymentIntentID string
	Status                  string
	PaymentStatus           string
	Completed               bool
}

func (p *devMarketplaceBillingProvider) ProviderName() string {
	return "dev"
}

func (p *devMarketplaceBillingProvider) StartCreatorOnboarding(_ context.Context, user *User, workspace *Workspace, existing *MarketplaceCreatorAccount) (*MarketplaceCreatorAccount, error) {
	now := time.Now()
	account := &MarketplaceCreatorAccount{
		ID:                    newID("mca"),
		UserID:                user.ID,
		WorkspaceID:           workspace.ID,
		Provider:              p.ProviderName(),
		ProviderAccountID:     newID("devacct"),
		OnboardingStatus:      "active",
		DetailsSubmitted:      true,
		ChargesEnabled:        true,
		PayoutsEnabled:        true,
		DashboardURL:          strings.TrimRight(p.cfg.AppOrigin, "/") + "/marketplace/publish",
		OnboardingCompletedAt: now,
		CreatedAt:             now,
		UpdatedAt:             now,
	}
	if existing != nil {
		account.ID = existing.ID
		account.ProviderAccountID = firstNonEmpty(existing.ProviderAccountID, account.ProviderAccountID)
		account.CreatedAt = existing.CreatedAt
	}
	return account, nil
}

func (p *devMarketplaceBillingProvider) RefreshCreatorAccount(_ context.Context, existing *MarketplaceCreatorAccount) (*MarketplaceCreatorAccount, error) {
	if existing == nil {
		return nil, fmt.Errorf("creator account not found")
	}
	return existing, nil
}

func (p *devMarketplaceBillingProvider) CreateCheckoutSession(_ context.Context, _ *MarketplaceListing, order *MarketplaceOrder, _ *User, _ *MarketplaceCreatorAccount) (*MarketplaceCheckoutResponse, error) {
	order.Provider = p.ProviderName()
	order.ProviderCheckoutSessionID = newID("devchk")
	order.ProviderPaymentIntentID = newID("devpi")
	return &MarketplaceCheckoutResponse{
		Provider:  p.ProviderName(),
		Completed: true,
		Order:     *order,
	}, nil
}

func (p *devMarketplaceBillingProvider) GetCheckoutSession(_ context.Context, order *MarketplaceOrder, _ *MarketplaceCreatorAccount) (*MarketplaceCheckoutSessionState, error) {
	return &MarketplaceCheckoutSessionState{
		Provider:                p.ProviderName(),
		ProviderCheckoutID:      order.ProviderCheckoutSessionID,
		ProviderPaymentIntentID: firstNonEmpty(order.ProviderPaymentIntentID, newID("devpi")),
		Status:                  "complete",
		PaymentStatus:           "paid",
		Completed:               true,
	}, nil
}

type disabledMarketplaceBillingProvider struct {
	provider string
	reason   string
}

func (p *disabledMarketplaceBillingProvider) ProviderName() string {
	return p.provider
}

func (p *disabledMarketplaceBillingProvider) StartCreatorOnboarding(context.Context, *User, *Workspace, *MarketplaceCreatorAccount) (*MarketplaceCreatorAccount, error) {
	return nil, fmt.Errorf("%s", p.reason)
}

func (p *disabledMarketplaceBillingProvider) RefreshCreatorAccount(context.Context, *MarketplaceCreatorAccount) (*MarketplaceCreatorAccount, error) {
	return nil, fmt.Errorf("%s", p.reason)
}

func (p *disabledMarketplaceBillingProvider) CreateCheckoutSession(context.Context, *MarketplaceListing, *MarketplaceOrder, *User, *MarketplaceCreatorAccount) (*MarketplaceCheckoutResponse, error) {
	return nil, fmt.Errorf("%s", p.reason)
}

func (p *disabledMarketplaceBillingProvider) GetCheckoutSession(context.Context, *MarketplaceOrder, *MarketplaceCreatorAccount) (*MarketplaceCheckoutSessionState, error) {
	return nil, fmt.Errorf("%s", p.reason)
}

type stripeMarketplaceBillingProvider struct {
	cfg    AppConfig
	client *http.Client
}

type stripeAPIErrorResponse struct {
	Error struct {
		Message string `json:"message"`
		Type    string `json:"type"`
	} `json:"error"`
}

type stripeAccountResponse struct {
	ID               string `json:"id"`
	DetailsSubmitted bool   `json:"details_submitted"`
	ChargesEnabled   bool   `json:"charges_enabled"`
	PayoutsEnabled   bool   `json:"payouts_enabled"`
}

type stripeAccountLinkResponse struct {
	URL string `json:"url"`
}

type stripeLoginLinkResponse struct {
	URL string `json:"url"`
}

type stripeCheckoutSessionResponse struct {
	ID            string            `json:"id"`
	URL           string            `json:"url"`
	Status        string            `json:"status"`
	PaymentStatus string            `json:"payment_status"`
	PaymentIntent string            `json:"payment_intent"`
	Metadata      map[string]string `json:"metadata"`
}

type stripeEventEnvelope struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Account string `json:"account"`
	Data    struct {
		Object json.RawMessage `json:"object"`
	} `json:"data"`
}

type stripeAccountWebhookObject struct {
	ID               string `json:"id"`
	DetailsSubmitted bool   `json:"details_submitted"`
	ChargesEnabled   bool   `json:"charges_enabled"`
	PayoutsEnabled   bool   `json:"payouts_enabled"`
}

func (p *stripeMarketplaceBillingProvider) ProviderName() string {
	return "stripe"
}

func (p *stripeMarketplaceBillingProvider) StartCreatorOnboarding(ctx context.Context, user *User, workspace *Workspace, existing *MarketplaceCreatorAccount) (*MarketplaceCreatorAccount, error) {
	accountID := ""
	createdAt := time.Now()
	if existing != nil {
		accountID = strings.TrimSpace(existing.ProviderAccountID)
		createdAt = existing.CreatedAt
	}
	if accountID == "" {
		values := url.Values{}
		values.Set("type", "express")
		values.Set("country", firstNonEmpty(strings.TrimSpace(p.cfg.Stripe.ConnectCountry), "US"))
		values.Set("email", user.Email)
		values.Set("business_type", "individual")
		values.Set("business_profile[product_description]", "Premium flashcard marketplace creator payouts for Vutadex")
		values.Set("business_profile[url]", strings.TrimRight(p.cfg.AppOrigin, "/"))
		values.Set("metadata[user_id]", user.ID)
		values.Set("metadata[workspace_id]", workspace.ID)
		values.Set("capabilities[card_payments][requested]", "true")
		values.Set("capabilities[transfers][requested]", "true")

		var created stripeAccountResponse
		if err := p.doStripeRequest(ctx, http.MethodPost, "/v1/accounts", "", values, &created); err != nil {
			return nil, err
		}
		accountID = created.ID
	}

	account, err := p.fetchStripeAccount(ctx, accountID)
	if err != nil {
		return nil, err
	}

	result := &MarketplaceCreatorAccount{
		ID:                newID("mca"),
		UserID:            user.ID,
		WorkspaceID:       workspace.ID,
		Provider:          p.ProviderName(),
		ProviderAccountID: accountID,
		CreatedAt:         createdAt,
		UpdatedAt:         time.Now(),
	}
	if existing != nil {
		result.ID = existing.ID
	}
	applyStripeAccountState(result, account)
	if result.ChargesEnabled && result.PayoutsEnabled && result.DetailsSubmitted {
		if dashboardURL, err := p.createLoginLink(ctx, accountID); err == nil {
			result.DashboardURL = dashboardURL
		}
		return result, nil
	}

	onboardingURL, err := p.createAccountLink(ctx, accountID)
	if err != nil {
		return nil, err
	}
	result.OnboardingURL = onboardingURL
	return result, nil
}

func (p *stripeMarketplaceBillingProvider) RefreshCreatorAccount(ctx context.Context, existing *MarketplaceCreatorAccount) (*MarketplaceCreatorAccount, error) {
	if existing == nil || strings.TrimSpace(existing.ProviderAccountID) == "" {
		return nil, fmt.Errorf("creator account not found")
	}

	account, err := p.fetchStripeAccount(ctx, existing.ProviderAccountID)
	if err != nil {
		return nil, err
	}

	refreshed := *existing
	refreshed.UpdatedAt = time.Now()
	refreshed.OnboardingURL = ""
	refreshed.DashboardURL = ""
	applyStripeAccountState(&refreshed, account)
	if refreshed.ChargesEnabled && refreshed.PayoutsEnabled && refreshed.DetailsSubmitted {
		if dashboardURL, err := p.createLoginLink(ctx, refreshed.ProviderAccountID); err == nil {
			refreshed.DashboardURL = dashboardURL
		}
	} else {
		if onboardingURL, err := p.createAccountLink(ctx, refreshed.ProviderAccountID); err == nil {
			refreshed.OnboardingURL = onboardingURL
		}
	}
	return &refreshed, nil
}

func (p *stripeMarketplaceBillingProvider) CreateCheckoutSession(ctx context.Context, listing *MarketplaceListing, order *MarketplaceOrder, buyer *User, creatorAccount *MarketplaceCreatorAccount) (*MarketplaceCheckoutResponse, error) {
	connectedAccountID := ""
	if creatorAccount != nil {
		connectedAccountID = strings.TrimSpace(creatorAccount.ProviderAccountID)
	}
	if connectedAccountID == "" {
		return nil, fmt.Errorf("creator account is missing")
	}

	values := url.Values{}
	values.Set("mode", "payment")
	values.Set("client_reference_id", order.ID)
	values.Set("success_url", appendURLQuery(p.cfg.Stripe.CheckoutSuccessURL, map[string]string{
		"checkout":            "success",
		"checkout_session_id": "{CHECKOUT_SESSION_ID}",
	}))
	values.Set("cancel_url", appendURLQuery(p.cfg.Stripe.CheckoutCancelURL, map[string]string{
		"checkout":            "cancelled",
		"checkout_session_id": "{CHECKOUT_SESSION_ID}",
	}))
	values.Set("line_items[0][price_data][currency]", strings.ToLower(firstNonEmpty(strings.TrimSpace(listing.Currency), "USD")))
	values.Set("line_items[0][price_data][product_data][name]", listing.Title)
	values.Set("line_items[0][price_data][product_data][description]", firstNonEmpty(strings.TrimSpace(listing.Summary), strings.TrimSpace(listing.Description), "Premium Vutadex marketplace deck"))
	values.Set("line_items[0][price_data][unit_amount]", strconv.Itoa(listing.PriceCents))
	values.Set("line_items[0][quantity]", "1")
	values.Set("payment_intent_data[application_fee_amount]", strconv.Itoa(order.PlatformFeeCents))
	values.Set("metadata[order_id]", order.ID)
	values.Set("metadata[listing_id]", listing.ID)
	values.Set("metadata[listing_version_number]", strconv.Itoa(order.ListingVersionNumber))
	values.Set("payment_intent_data[metadata][order_id]", order.ID)
	values.Set("payment_intent_data[metadata][listing_id]", listing.ID)
	values.Set("payment_intent_data[metadata][listing_version_number]", strconv.Itoa(order.ListingVersionNumber))
	if strings.TrimSpace(buyer.Email) != "" {
		values.Set("customer_email", buyer.Email)
	}

	var created stripeCheckoutSessionResponse
	if err := p.doStripeRequest(ctx, http.MethodPost, "/v1/checkout/sessions", connectedAccountID, values, &created); err != nil {
		return nil, err
	}

	order.Provider = p.ProviderName()
	order.ProviderCheckoutSessionID = created.ID
	order.ProviderPaymentIntentID = created.PaymentIntent

	return &MarketplaceCheckoutResponse{
		Provider:    p.ProviderName(),
		CheckoutURL: created.URL,
		Completed:   false,
		Order:       *order,
	}, nil
}

func (p *stripeMarketplaceBillingProvider) GetCheckoutSession(ctx context.Context, order *MarketplaceOrder, creatorAccount *MarketplaceCreatorAccount) (*MarketplaceCheckoutSessionState, error) {
	if order == nil || strings.TrimSpace(order.ProviderCheckoutSessionID) == "" {
		return nil, fmt.Errorf("checkout session is missing")
	}
	if creatorAccount == nil || strings.TrimSpace(creatorAccount.ProviderAccountID) == "" {
		return nil, fmt.Errorf("creator account is missing")
	}

	var session stripeCheckoutSessionResponse
	if err := p.doStripeRequest(ctx, http.MethodGet, "/v1/checkout/sessions/"+url.PathEscape(order.ProviderCheckoutSessionID), creatorAccount.ProviderAccountID, nil, &session); err != nil {
		return nil, err
	}

	return &MarketplaceCheckoutSessionState{
		Provider:                p.ProviderName(),
		ProviderCheckoutURL:     session.URL,
		ProviderCheckoutID:      session.ID,
		ProviderPaymentIntentID: session.PaymentIntent,
		Status:                  session.Status,
		PaymentStatus:           session.PaymentStatus,
		Completed:               strings.EqualFold(session.PaymentStatus, "paid") || strings.EqualFold(session.Status, "complete"),
	}, nil
}

func (p *stripeMarketplaceBillingProvider) fetchStripeAccount(ctx context.Context, accountID string) (*stripeAccountResponse, error) {
	var account stripeAccountResponse
	if err := p.doStripeRequest(ctx, http.MethodGet, "/v1/accounts/"+url.PathEscape(accountID), "", nil, &account); err != nil {
		return nil, err
	}
	return &account, nil
}

func (p *stripeMarketplaceBillingProvider) createAccountLink(ctx context.Context, accountID string) (string, error) {
	values := url.Values{}
	values.Set("account", accountID)
	values.Set("refresh_url", p.cfg.Stripe.ConnectRefreshURL)
	values.Set("return_url", p.cfg.Stripe.ConnectReturnURL)
	values.Set("type", "account_onboarding")

	var link stripeAccountLinkResponse
	if err := p.doStripeRequest(ctx, http.MethodPost, "/v1/account_links", "", values, &link); err != nil {
		return "", err
	}
	return link.URL, nil
}

func (p *stripeMarketplaceBillingProvider) createLoginLink(ctx context.Context, accountID string) (string, error) {
	var link stripeLoginLinkResponse
	if err := p.doStripeRequest(ctx, http.MethodPost, "/v1/accounts/"+url.PathEscape(accountID)+"/login_links", "", url.Values{}, &link); err != nil {
		return "", err
	}
	return link.URL, nil
}

func (p *stripeMarketplaceBillingProvider) doStripeRequest(ctx context.Context, method, path, connectedAccountID string, values url.Values, out any) error {
	endpoint := "https://api.stripe.com" + path
	var body io.Reader
	if method == http.MethodGet && len(values) > 0 {
		endpoint += "?" + values.Encode()
	} else if method != http.MethodGet {
		body = bytes.NewBufferString(values.Encode())
	}

	req, err := http.NewRequestWithContext(ctx, method, endpoint, body)
	if err != nil {
		return err
	}
	req.SetBasicAuth(p.cfg.Stripe.SecretKey, "")
	if method != http.MethodGet {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	if strings.TrimSpace(connectedAccountID) != "" {
		req.Header.Set("Stripe-Account", connectedAccountID)
	}

	client := p.client
	if client == nil {
		client = http.DefaultClient
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	payload, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var apiErr stripeAPIErrorResponse
		if json.Unmarshal(payload, &apiErr) == nil && strings.TrimSpace(apiErr.Error.Message) != "" {
			return fmt.Errorf("stripe %s failed: %s", path, apiErr.Error.Message)
		}
		return fmt.Errorf("stripe %s failed: status %d", path, resp.StatusCode)
	}
	if out == nil || len(payload) == 0 {
		return nil
	}
	if err := json.Unmarshal(payload, out); err != nil {
		return fmt.Errorf("failed to decode stripe response for %s: %w", path, err)
	}
	return nil
}

func newMarketplaceBillingProvider(cfg AppConfig) marketplaceBillingProvider {
	if cfg.IsDevelopment() && strings.TrimSpace(cfg.Stripe.SecretKey) == "" {
		return &devMarketplaceBillingProvider{cfg: cfg}
	}
	if strings.TrimSpace(cfg.Stripe.SecretKey) != "" {
		return &stripeMarketplaceBillingProvider{cfg: cfg}
	}
	return &disabledMarketplaceBillingProvider{
		provider: "disabled",
		reason:   "Marketplace billing is not configured. Set Stripe configuration or run in development mode.",
	}
}

func marketplacePlatformFeeCents(totalCents int, basisPoints int) int {
	if totalCents <= 0 || basisPoints <= 0 {
		return 0
	}
	return (totalCents*basisPoints + 9999) / 10000
}

func applyStripeAccountState(account *MarketplaceCreatorAccount, stripeAccount *stripeAccountResponse) {
	if account == nil || stripeAccount == nil {
		return
	}
	account.ProviderAccountID = stripeAccount.ID
	account.DetailsSubmitted = stripeAccount.DetailsSubmitted
	account.ChargesEnabled = stripeAccount.ChargesEnabled
	account.PayoutsEnabled = stripeAccount.PayoutsEnabled
	if account.DetailsSubmitted && account.ChargesEnabled && account.PayoutsEnabled {
		account.OnboardingStatus = "active"
		if account.OnboardingCompletedAt.IsZero() {
			account.OnboardingCompletedAt = time.Now()
		}
		return
	}
	account.OnboardingStatus = "pending"
	account.OnboardingCompletedAt = time.Time{}
}

func appendURLQuery(raw string, values map[string]string) string {
	parsed, err := url.Parse(strings.TrimSpace(raw))
	if err != nil || parsed == nil {
		return raw
	}
	query := parsed.Query()
	for key, value := range values {
		if strings.TrimSpace(value) == "" {
			continue
		}
		query.Set(key, value)
	}
	parsed.RawQuery = strings.ReplaceAll(query.Encode(), url.QueryEscape("{CHECKOUT_SESSION_ID}"), "{CHECKOUT_SESSION_ID}")
	return parsed.String()
}

func verifyStripeWebhookSignature(payload []byte, signatureHeader, secret string, now time.Time) error {
	signatureHeader = strings.TrimSpace(signatureHeader)
	secret = strings.TrimSpace(secret)
	if signatureHeader == "" || secret == "" {
		return fmt.Errorf("missing stripe webhook signature")
	}

	var timestamp string
	signatures := make([]string, 0, 2)
	for _, part := range strings.Split(signatureHeader, ",") {
		part = strings.TrimSpace(part)
		switch {
		case strings.HasPrefix(part, "t="):
			timestamp = strings.TrimPrefix(part, "t=")
		case strings.HasPrefix(part, "v1="):
			signatures = append(signatures, strings.TrimPrefix(part, "v1="))
		}
	}
	if timestamp == "" || len(signatures) == 0 {
		return fmt.Errorf("invalid stripe signature header")
	}

	issuedAt, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid stripe signature timestamp")
	}
	if now.Sub(time.Unix(issuedAt, 0)) > 5*time.Minute || time.Unix(issuedAt, 0).Sub(now) > 5*time.Minute {
		return fmt.Errorf("stripe webhook signature expired")
	}

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(timestamp))
	mac.Write([]byte("."))
	mac.Write(payload)
	expected := mac.Sum(nil)

	for _, signature := range signatures {
		decoded, err := hex.DecodeString(signature)
		if err != nil {
			continue
		}
		if hmac.Equal(decoded, expected) {
			return nil
		}
	}
	return fmt.Errorf("invalid stripe webhook signature")
}

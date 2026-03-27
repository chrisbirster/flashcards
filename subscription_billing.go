package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	subscriptionProviderStripe = "stripe"
	subscriptionProviderManual = "manual"
)

type subscriptionBillingProvider interface {
	ProviderName() string
	CreateCheckoutSession(ctx context.Context, params subscriptionCheckoutParams) (*subscriptionCheckoutResult, error)
	CreatePortalSession(ctx context.Context, customerID string) (string, error)
	GetCheckoutSession(ctx context.Context, checkoutSessionID string) (*subscriptionCheckoutSessionState, error)
	GetSubscription(ctx context.Context, subscriptionID string) (*subscriptionBillingState, error)
	UpdateSubscriptionQuantity(ctx context.Context, subscription *Subscription, quantity int, prorationBehavior string) (*subscriptionBillingState, error)
}

type disabledSubscriptionBillingProvider struct {
	reason string
}

type stripeSubscriptionBillingProvider struct {
	cfg    AppConfig
	client *http.Client
}

type subscriptionCheckoutParams struct {
	User                 *User
	Workspace            *Workspace
	Organization         *Organization
	ExistingSubscription *Subscription
	Plan                 Plan
	Quantity             int
	ClearOnboarding      bool
}

type subscriptionCheckoutResult struct {
	Provider      string
	CheckoutURL   string
	CheckoutState subscriptionCheckoutSessionState
}

type subscriptionCheckoutSessionState struct {
	Provider          string
	CheckoutSessionID string
	CheckoutURL       string
	CustomerID        string
	SubscriptionID    string
	Status            string
	PaymentStatus     string
	Completed         bool
}

type subscriptionBillingState struct {
	Provider               string
	Plan                   Plan
	Status                 string
	CustomerID             string
	SubscriptionID         string
	SubscriptionItemID     string
	CurrentPeriodEnd       time.Time
	CancelAtPeriodEnd      bool
	BilledQuantity         int
	CheckoutSessionID      string
}

type stripeBillingCustomerResponse struct {
	ID string `json:"id"`
}

type stripeBillingPortalSessionResponse struct {
	URL string `json:"url"`
}

type stripeBillingCheckoutSessionResponse struct {
	ID            string            `json:"id"`
	URL           string            `json:"url"`
	Customer      string            `json:"customer"`
	Subscription  string            `json:"subscription"`
	Status        string            `json:"status"`
	PaymentStatus string            `json:"payment_status"`
	Metadata      map[string]string `json:"metadata"`
}

type stripeBillingSubscriptionResponse struct {
	ID                string            `json:"id"`
	Customer          string            `json:"customer"`
	Status            string            `json:"status"`
	CancelAtPeriodEnd bool              `json:"cancel_at_period_end"`
	CurrentPeriodEnd  int64             `json:"current_period_end"`
	Metadata          map[string]string `json:"metadata"`
	Items             struct {
		Data []struct {
			ID       string `json:"id"`
			Quantity int    `json:"quantity"`
			Price    struct {
				ID string `json:"id"`
			} `json:"price"`
		} `json:"data"`
	} `json:"items"`
}

func (p *disabledSubscriptionBillingProvider) ProviderName() string {
	return "disabled"
}

func (p *disabledSubscriptionBillingProvider) CreateCheckoutSession(context.Context, subscriptionCheckoutParams) (*subscriptionCheckoutResult, error) {
	return nil, fmt.Errorf("%s", p.reason)
}

func (p *disabledSubscriptionBillingProvider) CreatePortalSession(context.Context, string) (string, error) {
	return "", fmt.Errorf("%s", p.reason)
}

func (p *disabledSubscriptionBillingProvider) GetCheckoutSession(context.Context, string) (*subscriptionCheckoutSessionState, error) {
	return nil, fmt.Errorf("%s", p.reason)
}

func (p *disabledSubscriptionBillingProvider) GetSubscription(context.Context, string) (*subscriptionBillingState, error) {
	return nil, fmt.Errorf("%s", p.reason)
}

func (p *disabledSubscriptionBillingProvider) UpdateSubscriptionQuantity(context.Context, *Subscription, int, string) (*subscriptionBillingState, error) {
	return nil, fmt.Errorf("%s", p.reason)
}

func (p *stripeSubscriptionBillingProvider) ProviderName() string {
	return subscriptionProviderStripe
}

func (p *stripeSubscriptionBillingProvider) CreateCheckoutSession(ctx context.Context, params subscriptionCheckoutParams) (*subscriptionCheckoutResult, error) {
	if params.User == nil || params.Workspace == nil {
		return nil, fmt.Errorf("checkout requires user and workspace")
	}
	priceID := billingPriceIDForPlan(p.cfg, params.Plan)
	if strings.TrimSpace(priceID) == "" {
		return nil, fmt.Errorf("no Stripe price configured for %s", params.Plan)
	}

	customerID, err := p.ensureCustomer(ctx, params)
	if err != nil {
		return nil, err
	}

	quantity := params.Quantity
	if quantity <= 0 {
		quantity = 1
	}

	values := url.Values{}
	values.Set("mode", "subscription")
	values.Set("customer", customerID)
	values.Set("client_reference_id", params.Workspace.ID)
	values.Set("success_url", appendURLQuery(p.cfg.Stripe.BillingCheckoutSuccessURL, map[string]string{
		"checkout":            "success",
		"checkout_session_id": "{CHECKOUT_SESSION_ID}",
	}))
	values.Set("cancel_url", appendURLQuery(p.cfg.Stripe.BillingCheckoutCancelURL, map[string]string{
		"checkout":            "cancelled",
		"checkout_session_id": "{CHECKOUT_SESSION_ID}",
	}))
	values.Set("line_items[0][price]", priceID)
	values.Set("line_items[0][quantity]", strconv.Itoa(quantity))
	values.Set("metadata[workspace_id]", params.Workspace.ID)
	values.Set("metadata[user_id]", params.User.ID)
	values.Set("metadata[plan]", string(params.Plan))
	values.Set("metadata[clear_onboarding]", strconv.FormatBool(params.ClearOnboarding))
	values.Set("subscription_data[metadata][workspace_id]", params.Workspace.ID)
	values.Set("subscription_data[metadata][user_id]", params.User.ID)
	values.Set("subscription_data[metadata][plan]", string(params.Plan))
	values.Set("subscription_data[metadata][clear_onboarding]", strconv.FormatBool(params.ClearOnboarding))
	if params.Organization != nil {
		values.Set("metadata[organization_id]", params.Organization.ID)
		values.Set("subscription_data[metadata][organization_id]", params.Organization.ID)
	}

	var created stripeBillingCheckoutSessionResponse
	if err := p.doStripeRequest(ctx, http.MethodPost, "/v1/checkout/sessions", values, &created); err != nil {
		return nil, err
	}

	return &subscriptionCheckoutResult{
		Provider: p.ProviderName(),
		CheckoutURL: created.URL,
		CheckoutState: subscriptionCheckoutSessionState{
			Provider:          p.ProviderName(),
			CheckoutSessionID: created.ID,
			CheckoutURL:       created.URL,
			CustomerID:        created.Customer,
			SubscriptionID:    created.Subscription,
			Status:            created.Status,
			PaymentStatus:     created.PaymentStatus,
			Completed:         strings.EqualFold(created.Status, "complete") || strings.EqualFold(created.PaymentStatus, "paid"),
		},
	}, nil
}

func (p *stripeSubscriptionBillingProvider) CreatePortalSession(ctx context.Context, customerID string) (string, error) {
	values := url.Values{}
	values.Set("customer", strings.TrimSpace(customerID))
	values.Set("return_url", p.cfg.Stripe.BillingPortalReturnURL)

	var session stripeBillingPortalSessionResponse
	if err := p.doStripeRequest(ctx, http.MethodPost, "/v1/billing_portal/sessions", values, &session); err != nil {
		return "", err
	}
	return session.URL, nil
}

func (p *stripeSubscriptionBillingProvider) GetCheckoutSession(ctx context.Context, checkoutSessionID string) (*subscriptionCheckoutSessionState, error) {
	var session stripeBillingCheckoutSessionResponse
	if err := p.doStripeRequest(ctx, http.MethodGet, "/v1/checkout/sessions/"+url.PathEscape(strings.TrimSpace(checkoutSessionID)), nil, &session); err != nil {
		return nil, err
	}
	return &subscriptionCheckoutSessionState{
		Provider:          p.ProviderName(),
		CheckoutSessionID: session.ID,
		CheckoutURL:       session.URL,
		CustomerID:        session.Customer,
		SubscriptionID:    session.Subscription,
		Status:            session.Status,
		PaymentStatus:     session.PaymentStatus,
		Completed:         strings.EqualFold(session.Status, "complete") || strings.EqualFold(session.PaymentStatus, "paid"),
	}, nil
}

func (p *stripeSubscriptionBillingProvider) GetSubscription(ctx context.Context, subscriptionID string) (*subscriptionBillingState, error) {
	var subscription stripeBillingSubscriptionResponse
	if err := p.doStripeRequest(ctx, http.MethodGet, "/v1/subscriptions/"+url.PathEscape(strings.TrimSpace(subscriptionID)), nil, &subscription); err != nil {
		return nil, err
	}
	return p.stateFromStripeSubscription(&subscription), nil
}

func (p *stripeSubscriptionBillingProvider) UpdateSubscriptionQuantity(ctx context.Context, subscription *Subscription, quantity int, prorationBehavior string) (*subscriptionBillingState, error) {
	if subscription == nil || strings.TrimSpace(subscription.ProviderSubscriptionID) == "" || strings.TrimSpace(subscription.ProviderSubscriptionItemID) == "" {
		return nil, fmt.Errorf("subscription is missing Stripe identifiers")
	}
	if quantity <= 0 {
		quantity = 1
	}
	if strings.TrimSpace(prorationBehavior) == "" {
		prorationBehavior = "none"
	}

	values := url.Values{}
	values.Set("items[0][id]", subscription.ProviderSubscriptionItemID)
	values.Set("items[0][quantity]", strconv.Itoa(quantity))
	values.Set("proration_behavior", prorationBehavior)

	var updated stripeBillingSubscriptionResponse
	if err := p.doStripeRequest(ctx, http.MethodPost, "/v1/subscriptions/"+url.PathEscape(subscription.ProviderSubscriptionID), values, &updated); err != nil {
		return nil, err
	}
	return p.stateFromStripeSubscription(&updated), nil
}

func (p *stripeSubscriptionBillingProvider) ensureCustomer(ctx context.Context, params subscriptionCheckoutParams) (string, error) {
	if params.ExistingSubscription != nil && strings.TrimSpace(params.ExistingSubscription.ProviderCustomerID) != "" {
		return strings.TrimSpace(params.ExistingSubscription.ProviderCustomerID), nil
	}

	values := url.Values{}
	values.Set("email", params.User.Email)
	values.Set("name", firstNonEmpty(strings.TrimSpace(params.OrganizationName()), strings.TrimSpace(params.Workspace.Name), strings.TrimSpace(params.User.DisplayName), params.User.Email))
	values.Set("metadata[workspace_id]", params.Workspace.ID)
	values.Set("metadata[user_id]", params.User.ID)
	values.Set("metadata[plan]", string(params.Plan))
	if params.Organization != nil {
		values.Set("metadata[organization_id]", params.Organization.ID)
	}

	var customer stripeBillingCustomerResponse
	if err := p.doStripeRequest(ctx, http.MethodPost, "/v1/customers", values, &customer); err != nil {
		return "", err
	}
	return customer.ID, nil
}

func (params subscriptionCheckoutParams) OrganizationName() string {
	if params.Organization == nil {
		return ""
	}
	return params.Organization.Name
}

func (p *stripeSubscriptionBillingProvider) stateFromStripeSubscription(subscription *stripeBillingSubscriptionResponse) *subscriptionBillingState {
	if subscription == nil {
		return nil
	}
	state := &subscriptionBillingState{
		Provider:          p.ProviderName(),
		SubscriptionID:    subscription.ID,
		CustomerID:        subscription.Customer,
		Status:            subscription.Status,
		CancelAtPeriodEnd: subscription.CancelAtPeriodEnd,
		BilledQuantity:    1,
	}
	if subscription.CurrentPeriodEnd > 0 {
		state.CurrentPeriodEnd = time.Unix(subscription.CurrentPeriodEnd, 0)
	}
	if len(subscription.Items.Data) > 0 {
		item := subscription.Items.Data[0]
		state.SubscriptionItemID = item.ID
		if item.Quantity > 0 {
			state.BilledQuantity = item.Quantity
		}
		state.Plan = billingPlanForPrice(p.cfg, item.Price.ID)
	}
	if state.Plan == "" {
		state.Plan = parsePlan(subscription.Metadata["plan"])
	}
	return state
}

func (p *stripeSubscriptionBillingProvider) doStripeRequest(ctx context.Context, method, path string, values url.Values, out any) error {
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

func newSubscriptionBillingProvider(cfg AppConfig) subscriptionBillingProvider {
	if strings.TrimSpace(cfg.Stripe.SecretKey) == "" {
		return &disabledSubscriptionBillingProvider{
			reason: "Subscription billing is not configured. Set Stripe billing env vars or use development fallback.",
		}
	}
	if strings.TrimSpace(cfg.Stripe.BillingPriceProMonthly) == "" || strings.TrimSpace(cfg.Stripe.BillingPriceTeamMonthly) == "" {
		return &disabledSubscriptionBillingProvider{
			reason: "Subscription billing price IDs are missing. Set VUTADEX_STRIPE_BILLING_PRICE_PRO_MONTHLY and VUTADEX_STRIPE_BILLING_PRICE_TEAM_MONTHLY.",
		}
	}
	return &stripeSubscriptionBillingProvider{cfg: cfg}
}

func billingPriceIDForPlan(cfg AppConfig, plan Plan) string {
	switch plan {
	case PlanPro:
		return strings.TrimSpace(cfg.Stripe.BillingPriceProMonthly)
	case PlanTeam:
		return strings.TrimSpace(cfg.Stripe.BillingPriceTeamMonthly)
	default:
		return ""
	}
}

func billingPlanForPrice(cfg AppConfig, priceID string) Plan {
	priceID = strings.TrimSpace(priceID)
	switch priceID {
	case strings.TrimSpace(cfg.Stripe.BillingPriceProMonthly):
		return PlanPro
	case strings.TrimSpace(cfg.Stripe.BillingPriceTeamMonthly):
		return PlanTeam
	default:
		return ""
	}
}

func teamSeatQuantity(activeMembers int) int {
	if activeMembers < 3 {
		return 3
	}
	return activeMembers
}

func subscriptionGrantsPlan(subscription *Subscription, now time.Time) bool {
	if subscription == nil {
		return false
	}
	status := strings.ToLower(strings.TrimSpace(subscription.Status))
	switch status {
	case "active", "trialing", "past_due":
		return true
	case "canceled", "cancelled":
		return !subscription.CurrentPeriodEnd.IsZero() && !now.After(subscription.CurrentPeriodEnd)
	default:
		return false
	}
}

func applySubscriptionState(subscription *Subscription, state *subscriptionBillingState) {
	if subscription == nil || state == nil {
		return
	}
	subscription.Provider = state.Provider
	subscription.Plan = state.Plan
	subscription.Status = state.Status
	subscription.ProviderCustomerID = state.CustomerID
	subscription.ProviderSubscriptionID = state.SubscriptionID
	subscription.ProviderSubscriptionItemID = state.SubscriptionItemID
	subscription.ProviderCheckoutSessionID = state.CheckoutSessionID
	subscription.CurrentPeriodEnd = state.CurrentPeriodEnd
	subscription.CancelAtPeriodEnd = state.CancelAtPeriodEnd
	if state.BilledQuantity > 0 {
		subscription.BilledQuantity = state.BilledQuantity
	}
}

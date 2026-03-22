package main

import (
	"database/sql"
	"time"
)

func (s *SQLiteStore) GetMarketplaceCreatorAccount(id string) (*MarketplaceCreatorAccount, error) {
	row := s.db.QueryRow(`
		SELECT id, user_id, workspace_id, provider, provider_account_id, onboarding_status,
		       details_submitted, charges_enabled, payouts_enabled, onboarding_url, dashboard_url,
		       onboarding_completed_at, created_at, updated_at
		FROM marketplace_creator_accounts
		WHERE id = ?
	`, id)
	return scanMarketplaceCreatorAccount(row)
}

func (s *SQLiteStore) GetMarketplaceCreatorAccountByUser(userID string) (*MarketplaceCreatorAccount, error) {
	row := s.db.QueryRow(`
		SELECT id, user_id, workspace_id, provider, provider_account_id, onboarding_status,
		       details_submitted, charges_enabled, payouts_enabled, onboarding_url, dashboard_url,
		       onboarding_completed_at, created_at, updated_at
		FROM marketplace_creator_accounts
		WHERE user_id = ?
	`, userID)
	return scanMarketplaceCreatorAccount(row)
}

func (s *SQLiteStore) GetMarketplaceCreatorAccountByProviderAccount(provider, providerAccountID string) (*MarketplaceCreatorAccount, error) {
	row := s.db.QueryRow(`
		SELECT id, user_id, workspace_id, provider, provider_account_id, onboarding_status,
		       details_submitted, charges_enabled, payouts_enabled, onboarding_url, dashboard_url,
		       onboarding_completed_at, created_at, updated_at
		FROM marketplace_creator_accounts
		WHERE provider = ? AND provider_account_id = ?
	`, provider, providerAccountID)
	return scanMarketplaceCreatorAccount(row)
}

func (s *SQLiteStore) UpsertMarketplaceCreatorAccount(account *MarketplaceCreatorAccount) error {
	_, err := s.db.Exec(`
		INSERT INTO marketplace_creator_accounts (
			id, user_id, workspace_id, provider, provider_account_id, onboarding_status,
			details_submitted, charges_enabled, payouts_enabled, onboarding_url, dashboard_url,
			onboarding_completed_at, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(user_id) DO UPDATE SET
			workspace_id = excluded.workspace_id,
			provider = excluded.provider,
			provider_account_id = excluded.provider_account_id,
			onboarding_status = excluded.onboarding_status,
			details_submitted = excluded.details_submitted,
			charges_enabled = excluded.charges_enabled,
			payouts_enabled = excluded.payouts_enabled,
			onboarding_url = excluded.onboarding_url,
			dashboard_url = excluded.dashboard_url,
			onboarding_completed_at = excluded.onboarding_completed_at,
			updated_at = excluded.updated_at
	`,
		account.ID,
		account.UserID,
		account.WorkspaceID,
		account.Provider,
		account.ProviderAccountID,
		account.OnboardingStatus,
		boolToInt(account.DetailsSubmitted),
		boolToInt(account.ChargesEnabled),
		boolToInt(account.PayoutsEnabled),
		account.OnboardingURL,
		account.DashboardURL,
		nullIfZeroTime(account.OnboardingCompletedAt),
		account.CreatedAt.Unix(),
		account.UpdatedAt.Unix(),
	)
	return err
}

func scanMarketplaceCreatorAccount(scanner interface{ Scan(dest ...any) error }) (*MarketplaceCreatorAccount, error) {
	var (
		account                                          MarketplaceCreatorAccount
		detailsSubmitted, chargesEnabled, payoutsEnabled int
		onboardingCompletedAt, createdAt, updatedAt      sql.NullInt64
	)
	if err := scanner.Scan(
		&account.ID,
		&account.UserID,
		&account.WorkspaceID,
		&account.Provider,
		&account.ProviderAccountID,
		&account.OnboardingStatus,
		&detailsSubmitted,
		&chargesEnabled,
		&payoutsEnabled,
		&account.OnboardingURL,
		&account.DashboardURL,
		&onboardingCompletedAt,
		&createdAt,
		&updatedAt,
	); err != nil {
		return nil, err
	}
	account.DetailsSubmitted = detailsSubmitted == 1
	account.ChargesEnabled = chargesEnabled == 1
	account.PayoutsEnabled = payoutsEnabled == 1
	account.OnboardingCompletedAt = unixTimeOrZero(onboardingCompletedAt)
	account.CreatedAt = unixTimeOrZero(createdAt)
	account.UpdatedAt = unixTimeOrZero(updatedAt)
	return &account, nil
}

func (s *SQLiteStore) CreateMarketplaceOrder(order *MarketplaceOrder) error {
	_, err := s.db.Exec(`
		INSERT INTO marketplace_orders (
			id, listing_id, listing_version_number, buyer_user_id, buyer_workspace_id, creator_user_id,
			creator_account_id, provider, provider_checkout_session_id, provider_payment_intent_id,
			status, amount_cents, currency, platform_fee_cents, creator_amount_cents, completed_at,
			created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`,
		order.ID,
		order.ListingID,
		order.ListingVersionNumber,
		order.BuyerUserID,
		order.BuyerWorkspaceID,
		order.CreatorUserID,
		nullIfEmpty(order.CreatorAccountID),
		order.Provider,
		order.ProviderCheckoutSessionID,
		order.ProviderPaymentIntentID,
		order.Status,
		order.AmountCents,
		order.Currency,
		order.PlatformFeeCents,
		order.CreatorAmountCents,
		nullIfZeroTime(order.CompletedAt),
		order.CreatedAt.Unix(),
		order.UpdatedAt.Unix(),
	)
	return err
}

func (s *SQLiteStore) GetMarketplaceOrderByCheckoutSession(provider, checkoutSessionID string) (*MarketplaceOrder, error) {
	row := s.db.QueryRow(`
		SELECT id, listing_id, listing_version_number, buyer_user_id, buyer_workspace_id, creator_user_id,
		       creator_account_id, provider, provider_checkout_session_id, provider_payment_intent_id,
		       status, amount_cents, currency, platform_fee_cents, creator_amount_cents, completed_at,
		       created_at, updated_at
		FROM marketplace_orders
		WHERE provider = ? AND provider_checkout_session_id = ?
	`, provider, checkoutSessionID)
	return scanMarketplaceOrder(row)
}

func (s *SQLiteStore) GetMarketplaceOrder(id string) (*MarketplaceOrder, error) {
	row := s.db.QueryRow(`
		SELECT id, listing_id, listing_version_number, buyer_user_id, buyer_workspace_id, creator_user_id,
		       creator_account_id, provider, provider_checkout_session_id, provider_payment_intent_id,
		       status, amount_cents, currency, platform_fee_cents, creator_amount_cents, completed_at,
		       created_at, updated_at
		FROM marketplace_orders
		WHERE id = ?
	`, id)
	return scanMarketplaceOrder(row)
}

func (s *SQLiteStore) UpdateMarketplaceOrder(order *MarketplaceOrder) error {
	_, err := s.db.Exec(`
		UPDATE marketplace_orders
		SET status = ?, provider_payment_intent_id = ?, completed_at = ?, updated_at = ?
		WHERE id = ?
	`,
		order.Status,
		order.ProviderPaymentIntentID,
		nullIfZeroTime(order.CompletedAt),
		order.UpdatedAt.Unix(),
		order.ID,
	)
	return err
}

func scanMarketplaceOrder(scanner interface{ Scan(dest ...any) error }) (*MarketplaceOrder, error) {
	var (
		order                             MarketplaceOrder
		creatorAccountID                  sql.NullString
		completedAt, createdAt, updatedAt sql.NullInt64
	)
	if err := scanner.Scan(
		&order.ID,
		&order.ListingID,
		&order.ListingVersionNumber,
		&order.BuyerUserID,
		&order.BuyerWorkspaceID,
		&order.CreatorUserID,
		&creatorAccountID,
		&order.Provider,
		&order.ProviderCheckoutSessionID,
		&order.ProviderPaymentIntentID,
		&order.Status,
		&order.AmountCents,
		&order.Currency,
		&order.PlatformFeeCents,
		&order.CreatorAmountCents,
		&completedAt,
		&createdAt,
		&updatedAt,
	); err != nil {
		return nil, err
	}
	if creatorAccountID.Valid {
		order.CreatorAccountID = creatorAccountID.String
	}
	order.CompletedAt = unixTimeOrZero(completedAt)
	order.CreatedAt = unixTimeOrZero(createdAt)
	order.UpdatedAt = unixTimeOrZero(updatedAt)
	return &order, nil
}

func (s *SQLiteStore) GetMarketplaceLicense(listingID, buyerUserID string) (*MarketplaceLicense, error) {
	row := s.db.QueryRow(`
		SELECT id, listing_id, buyer_user_id, order_id, status, granted_version_number, created_at, updated_at
		FROM marketplace_licenses
		WHERE listing_id = ? AND buyer_user_id = ?
	`, listingID, buyerUserID)
	return scanMarketplaceLicense(row)
}

func (s *SQLiteStore) UpsertMarketplaceLicense(license *MarketplaceLicense) error {
	_, err := s.db.Exec(`
		INSERT INTO marketplace_licenses (
			id, listing_id, buyer_user_id, order_id, status, granted_version_number, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(listing_id, buyer_user_id) DO UPDATE SET
			order_id = excluded.order_id,
			status = excluded.status,
			granted_version_number = excluded.granted_version_number,
			updated_at = excluded.updated_at
	`,
		license.ID,
		license.ListingID,
		license.BuyerUserID,
		license.OrderID,
		license.Status,
		license.GrantedVersionNumber,
		license.CreatedAt.Unix(),
		license.UpdatedAt.Unix(),
	)
	return err
}

func scanMarketplaceLicense(scanner interface{ Scan(dest ...any) error }) (*MarketplaceLicense, error) {
	var (
		license              MarketplaceLicense
		createdAt, updatedAt int64
	)
	if err := scanner.Scan(
		&license.ID,
		&license.ListingID,
		&license.BuyerUserID,
		&license.OrderID,
		&license.Status,
		&license.GrantedVersionNumber,
		&createdAt,
		&updatedAt,
	); err != nil {
		return nil, err
	}
	license.CreatedAt = time.Unix(createdAt, 0)
	license.UpdatedAt = time.Unix(updatedAt, 0)
	return &license, nil
}

func (s *SQLiteStore) GetMarketplacePayoutByOrder(orderID string) (*MarketplacePayout, error) {
	row := s.db.QueryRow(`
		SELECT id, order_id, creator_user_id, creator_account_id, provider, provider_transfer_id,
		       status, amount_cents, currency, platform_fee_cents, created_at, updated_at
		FROM marketplace_payouts
		WHERE order_id = ?
	`, orderID)
	return scanMarketplacePayout(row)
}

func (s *SQLiteStore) UpsertMarketplacePayout(payout *MarketplacePayout) error {
	_, err := s.db.Exec(`
		INSERT INTO marketplace_payouts (
			id, order_id, creator_user_id, creator_account_id, provider, provider_transfer_id,
			status, amount_cents, currency, platform_fee_cents, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(order_id) DO UPDATE SET
			provider_transfer_id = excluded.provider_transfer_id,
			status = excluded.status,
			updated_at = excluded.updated_at
	`,
		payout.ID,
		payout.OrderID,
		payout.CreatorUserID,
		payout.CreatorAccountID,
		payout.Provider,
		payout.ProviderTransferID,
		payout.Status,
		payout.AmountCents,
		payout.Currency,
		payout.PlatformFeeCents,
		payout.CreatedAt.Unix(),
		payout.UpdatedAt.Unix(),
	)
	return err
}

func scanMarketplacePayout(scanner interface{ Scan(dest ...any) error }) (*MarketplacePayout, error) {
	var (
		payout               MarketplacePayout
		createdAt, updatedAt int64
	)
	if err := scanner.Scan(
		&payout.ID,
		&payout.OrderID,
		&payout.CreatorUserID,
		&payout.CreatorAccountID,
		&payout.Provider,
		&payout.ProviderTransferID,
		&payout.Status,
		&payout.AmountCents,
		&payout.Currency,
		&payout.PlatformFeeCents,
		&createdAt,
		&updatedAt,
	); err != nil {
		return nil, err
	}
	payout.CreatedAt = time.Unix(createdAt, 0)
	payout.UpdatedAt = time.Unix(updatedAt, 0)
	return &payout, nil
}

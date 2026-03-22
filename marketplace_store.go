package main

import (
	"database/sql"
	"encoding/json"
	"strings"
	"time"
)

func sanitizeMarketplaceTags(tags []string) []string {
	seen := make(map[string]struct{}, len(tags))
	sanitized := make([]string, 0, len(tags))
	for _, raw := range tags {
		tag := strings.TrimSpace(sanitizeHTML(raw))
		if tag == "" {
			continue
		}
		key := strings.ToLower(tag)
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		sanitized = append(sanitized, tag)
	}
	return sanitized
}

func marketplaceTagsJSON(tags []string) string {
	bytes, err := json.Marshal(sanitizeMarketplaceTags(tags))
	if err != nil {
		return "[]"
	}
	return string(bytes)
}

func parseMarketplaceTags(raw string) []string {
	if strings.TrimSpace(raw) == "" {
		return nil
	}
	var tags []string
	if err := json.Unmarshal([]byte(raw), &tags); err != nil {
		return nil
	}
	return sanitizeMarketplaceTags(tags)
}

func (s *SQLiteStore) MarketplaceListingSlugExists(slug, excludeID string) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM marketplace_listings WHERE slug = ?`
	args := []any{slug}
	if strings.TrimSpace(excludeID) != "" {
		query += ` AND id <> ?`
		args = append(args, excludeID)
	}
	if err := s.db.QueryRow(query, args...).Scan(&count); err != nil {
		return false, err
	}
	return count > 0, nil
}

func (s *SQLiteStore) CreateMarketplaceListing(listing *MarketplaceListing) error {
	_, err := s.db.Exec(`
		INSERT INTO marketplace_listings (
			id, workspace_id, deck_id, slug, title, summary, description, category, tags, cover_image_url,
			creator_user_id, price_mode, price_cents, currency, status, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`,
		listing.ID,
		listing.WorkspaceID,
		listing.DeckID,
		listing.Slug,
		listing.Title,
		listing.Summary,
		listing.Description,
		listing.Category,
		marketplaceTagsJSON(listing.Tags),
		listing.CoverImageURL,
		listing.CreatorUserID,
		listing.PriceMode,
		listing.PriceCents,
		listing.Currency,
		listing.Status,
		listing.CreatedAt.Unix(),
		listing.UpdatedAt.Unix(),
	)
	return err
}

func (s *SQLiteStore) UpdateMarketplaceListing(listing *MarketplaceListing) error {
	_, err := s.db.Exec(`
		UPDATE marketplace_listings
		SET deck_id = ?, slug = ?, title = ?, summary = ?, description = ?, category = ?, tags = ?,
			cover_image_url = ?, price_mode = ?, price_cents = ?, currency = ?, status = ?, updated_at = ?
		WHERE id = ?
	`,
		listing.DeckID,
		listing.Slug,
		listing.Title,
		listing.Summary,
		listing.Description,
		listing.Category,
		marketplaceTagsJSON(listing.Tags),
		listing.CoverImageURL,
		listing.PriceMode,
		listing.PriceCents,
		listing.Currency,
		listing.Status,
		listing.UpdatedAt.Unix(),
		listing.ID,
	)
	return err
}

func (s *SQLiteStore) DeleteMarketplaceListing(id string) error {
	_, err := s.db.Exec(`DELETE FROM marketplace_listings WHERE id = ?`, id)
	return err
}

func (s *SQLiteStore) GetMarketplaceListingByID(id string) (*MarketplaceListing, error) {
	row := s.db.QueryRow(`
		SELECT id, workspace_id, deck_id, slug, title, summary, description, category, tags, cover_image_url,
		       creator_user_id, price_mode, price_cents, currency, status, created_at, updated_at
		FROM marketplace_listings
		WHERE id = ?
	`, id)
	return scanMarketplaceListing(row)
}

func (s *SQLiteStore) GetMarketplaceListingBySlug(slug string) (*MarketplaceListing, error) {
	row := s.db.QueryRow(`
		SELECT id, workspace_id, deck_id, slug, title, summary, description, category, tags, cover_image_url,
		       creator_user_id, price_mode, price_cents, currency, status, created_at, updated_at
		FROM marketplace_listings
		WHERE slug = ?
	`, slug)
	return scanMarketplaceListing(row)
}

func scanMarketplaceListing(scanner interface{ Scan(dest ...any) error }) (*MarketplaceListing, error) {
	var (
		listing       MarketplaceListing
		tagsJSON      string
		createdAtUnix int64
		updatedAtUnix int64
	)
	if err := scanner.Scan(
		&listing.ID,
		&listing.WorkspaceID,
		&listing.DeckID,
		&listing.Slug,
		&listing.Title,
		&listing.Summary,
		&listing.Description,
		&listing.Category,
		&tagsJSON,
		&listing.CoverImageURL,
		&listing.CreatorUserID,
		&listing.PriceMode,
		&listing.PriceCents,
		&listing.Currency,
		&listing.Status,
		&createdAtUnix,
		&updatedAtUnix,
	); err != nil {
		return nil, err
	}
	listing.Tags = parseMarketplaceTags(tagsJSON)
	listing.CreatedAt = time.Unix(createdAtUnix, 0)
	listing.UpdatedAt = time.Unix(updatedAtUnix, 0)
	return &listing, nil
}

func (s *SQLiteStore) CreateMarketplaceListingVersion(version *MarketplaceListingVersion) error {
	_, err := s.db.Exec(`
		INSERT INTO marketplace_listing_versions (
			id, listing_id, version_number, source_deck_id, published_by_user_id,
			change_summary, note_count, card_count, created_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`,
		version.ID,
		version.ListingID,
		version.VersionNumber,
		version.SourceDeckID,
		version.PublishedByUserID,
		version.ChangeSummary,
		version.NoteCount,
		version.CardCount,
		version.CreatedAt.Unix(),
	)
	return err
}

func (s *SQLiteStore) ListMarketplaceListingVersions(listingID string) ([]MarketplaceListingVersion, error) {
	rows, err := s.db.Query(`
		SELECT id, listing_id, version_number, source_deck_id, published_by_user_id,
		       change_summary, note_count, card_count, created_at
		FROM marketplace_listing_versions
		WHERE listing_id = ?
		ORDER BY version_number DESC
	`, listingID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var versions []MarketplaceListingVersion
	for rows.Next() {
		version, err := scanMarketplaceListingVersion(rows)
		if err != nil {
			return nil, err
		}
		versions = append(versions, *version)
	}
	return versions, rows.Err()
}

func (s *SQLiteStore) GetLatestMarketplaceListingVersion(listingID string) (*MarketplaceListingVersion, error) {
	row := s.db.QueryRow(`
		SELECT id, listing_id, version_number, source_deck_id, published_by_user_id,
		       change_summary, note_count, card_count, created_at
		FROM marketplace_listing_versions
		WHERE listing_id = ?
		ORDER BY version_number DESC
		LIMIT 1
	`, listingID)
	return scanMarketplaceListingVersion(row)
}

func scanMarketplaceListingVersion(scanner interface{ Scan(dest ...any) error }) (*MarketplaceListingVersion, error) {
	var version MarketplaceListingVersion
	var createdAtUnix int64
	if err := scanner.Scan(
		&version.ID,
		&version.ListingID,
		&version.VersionNumber,
		&version.SourceDeckID,
		&version.PublishedByUserID,
		&version.ChangeSummary,
		&version.NoteCount,
		&version.CardCount,
		&createdAtUnix,
	); err != nil {
		return nil, err
	}
	version.CreatedAt = time.Unix(createdAtUnix, 0)
	return &version, nil
}

func (s *SQLiteStore) CreateMarketplaceInstall(install *MarketplaceInstall) error {
	_, err := s.db.Exec(`
		INSERT INTO marketplace_installs (
			id, listing_id, workspace_id, installed_by_user_id, installed_deck_id,
			source_version_number, status, superseded_by_install_id, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`,
		install.ID,
		install.ListingID,
		install.WorkspaceID,
		install.InstalledByUserID,
		nullableDeckID(install.InstalledDeckID),
		install.SourceVersionNumber,
		install.Status,
		nullIfEmpty(install.SupersededByInstall),
		install.CreatedAt.Unix(),
		install.UpdatedAt.Unix(),
	)
	return err
}

func (s *SQLiteStore) GetMarketplaceInstall(id string) (*MarketplaceInstall, error) {
	row := s.db.QueryRow(`
		SELECT i.id, i.listing_id, i.workspace_id, i.installed_by_user_id, i.installed_deck_id,
		       d.name, i.source_version_number, i.status, i.superseded_by_install_id, i.created_at, i.updated_at
		FROM marketplace_installs i
		LEFT JOIN decks d ON d.id = i.installed_deck_id
		WHERE i.id = ?
	`, id)
	return scanMarketplaceInstall(row)
}

func (s *SQLiteStore) GetCurrentMarketplaceInstall(listingID, userID string) (*MarketplaceInstall, error) {
	row := s.db.QueryRow(`
		SELECT i.id, i.listing_id, i.workspace_id, i.installed_by_user_id, i.installed_deck_id,
		       d.name, i.source_version_number, i.status, i.superseded_by_install_id, i.created_at, i.updated_at
		FROM marketplace_installs i
		LEFT JOIN decks d ON d.id = i.installed_deck_id
		WHERE i.listing_id = ? AND i.installed_by_user_id = ? AND i.status = 'active'
		ORDER BY i.created_at DESC
		LIMIT 1
	`, listingID, userID)
	return scanMarketplaceInstall(row)
}

func scanMarketplaceInstall(scanner interface{ Scan(dest ...any) error }) (*MarketplaceInstall, error) {
	var (
		install         MarketplaceInstall
		installedDeckID sql.NullInt64
		deckName        sql.NullString
		supersededByID  sql.NullString
		createdAtUnix   int64
		updatedAtUnix   int64
	)
	if err := scanner.Scan(
		&install.ID,
		&install.ListingID,
		&install.WorkspaceID,
		&install.InstalledByUserID,
		&installedDeckID,
		&deckName,
		&install.SourceVersionNumber,
		&install.Status,
		&supersededByID,
		&createdAtUnix,
		&updatedAtUnix,
	); err != nil {
		return nil, err
	}
	if installedDeckID.Valid {
		install.InstalledDeckID = installedDeckID.Int64
	}
	if deckName.Valid {
		install.InstalledDeckName = deckName.String
	}
	if supersededByID.Valid {
		install.SupersededByInstall = supersededByID.String
	}
	install.CreatedAt = time.Unix(createdAtUnix, 0)
	install.UpdatedAt = time.Unix(updatedAtUnix, 0)
	return &install, nil
}

func (s *SQLiteStore) UpdateMarketplaceInstall(install *MarketplaceInstall) error {
	_, err := s.db.Exec(`
		UPDATE marketplace_installs
		SET workspace_id = ?, installed_deck_id = ?, source_version_number = ?, status = ?,
		    superseded_by_install_id = ?, updated_at = ?
		WHERE id = ?
	`,
		install.WorkspaceID,
		nullableDeckID(install.InstalledDeckID),
		install.SourceVersionNumber,
		install.Status,
		nullIfEmpty(install.SupersededByInstall),
		install.UpdatedAt.Unix(),
		install.ID,
	)
	return err
}

func (s *SQLiteStore) CountMarketplaceInstalls(listingID string) (int, error) {
	var count int
	if err := s.db.QueryRow(`
		SELECT COUNT(*)
		FROM marketplace_installs
		WHERE listing_id = ? AND status = 'active'
	`, listingID).Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func (s *SQLiteStore) resolveMarketplaceListing(ref string) (*MarketplaceListing, error) {
	ref = strings.TrimSpace(ref)
	if ref == "" {
		return nil, sql.ErrNoRows
	}
	if listing, err := s.GetMarketplaceListingByID(ref); err == nil {
		return listing, nil
	}
	return s.GetMarketplaceListingBySlug(ref)
}

func (s *SQLiteStore) BuildMarketplaceListingSummary(listing *MarketplaceListing, userID, workspaceID string) (MarketplaceListingSummary, error) {
	creator, err := s.GetUserByID(listing.CreatorUserID)
	if err != nil {
		return MarketplaceListingSummary{}, err
	}
	deck, err := s.GetDeck(listing.DeckID)
	if err != nil {
		return MarketplaceListingSummary{}, err
	}
	installCount, err := s.CountMarketplaceInstalls(listing.ID)
	if err != nil {
		return MarketplaceListingSummary{}, err
	}

	summary := MarketplaceListingSummary{
		ID:                 listing.ID,
		Slug:               listing.Slug,
		Title:              listing.Title,
		Summary:            listing.Summary,
		Description:        listing.Description,
		Category:           listing.Category,
		Tags:               append([]string(nil), listing.Tags...),
		CoverImageURL:      listing.CoverImageURL,
		CreatorUserID:      listing.CreatorUserID,
		CreatorDisplayName: creator.DisplayName,
		CreatorEmail:       creator.Email,
		WorkspaceID:        listing.WorkspaceID,
		SourceDeckID:       listing.DeckID,
		SourceDeckName:     deck.Name,
		PriceMode:          listing.PriceMode,
		PriceCents:         listing.PriceCents,
		Currency:           listing.Currency,
		Status:             listing.Status,
		InstallCount:       installCount,
		CanEdit:            listing.CreatorUserID == userID && listing.WorkspaceID == workspaceID,
		CreatedAt:          listing.CreatedAt,
		UpdatedAt:          listing.UpdatedAt,
	}

	if latestVersion, err := s.GetLatestMarketplaceListingVersion(listing.ID); err == nil {
		summary.LatestVersionNumber = latestVersion.VersionNumber
	}
	if userID != "" {
		if license, err := s.GetMarketplaceLicense(listing.ID, userID); err == nil && license.Status == "active" {
			summary.CurrentUserLicense = license
		}
		if install, err := s.GetCurrentMarketplaceInstall(listing.ID, userID); err == nil {
			summary.CurrentUserInstall = install
			if summary.LatestVersionNumber > 0 && install.SourceVersionNumber < summary.LatestVersionNumber {
				summary.UpdateAvailable = true
			}
		}
	}

	return summary, nil
}

func (s *SQLiteStore) ListMarketplaceListings(scope, userID, workspaceID string) ([]MarketplaceListingSummary, error) {
	scope = strings.TrimSpace(scope)
	query := `
		SELECT id, workspace_id, deck_id, slug, title, summary, description, category, tags, cover_image_url,
		       creator_user_id, price_mode, price_cents, currency, status, created_at, updated_at
		FROM marketplace_listings
	`
	args := make([]any, 0, 2)
	switch scope {
	case "mine":
		query += ` WHERE creator_user_id = ? AND workspace_id = ?`
		args = append(args, userID, workspaceID)
	default:
		query += ` WHERE status = 'published'`
	}
	query += ` ORDER BY updated_at DESC, created_at DESC`

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var listings []MarketplaceListingSummary
	for rows.Next() {
		listing, err := scanMarketplaceListing(rows)
		if err != nil {
			return nil, err
		}
		summary, err := s.BuildMarketplaceListingSummary(listing, userID, workspaceID)
		if err != nil {
			return nil, err
		}
		listings = append(listings, summary)
	}
	return listings, rows.Err()
}

func (s *SQLiteStore) BuildMarketplaceListingDetail(ref, userID, workspaceID string) (*MarketplaceListingDetail, error) {
	listing, err := s.resolveMarketplaceListing(ref)
	if err != nil {
		return nil, err
	}
	summary, err := s.BuildMarketplaceListingSummary(listing, userID, workspaceID)
	if err != nil {
		return nil, err
	}
	if summary.Status != "published" && !summary.CanEdit {
		return nil, sql.ErrNoRows
	}

	detail := &MarketplaceListingDetail{
		Listing:            summary,
		CurrentUserLicense: summary.CurrentUserLicense,
		CurrentUserInstall: summary.CurrentUserInstall,
		UpdateAvailable:    summary.UpdateAvailable,
		CanEdit:            summary.CanEdit,
		CanPublish:         summary.CanEdit,
	}
	if latestVersion, err := s.GetLatestMarketplaceListingVersion(listing.ID); err == nil {
		detail.LatestVersion = latestVersion
	}
	if versions, err := s.ListMarketplaceListingVersions(listing.ID); err == nil {
		detail.Versions = versions
	}
	if userID != "" {
		if workspaces, err := s.ListWorkspacesForUser(userID); err == nil {
			detail.AvailableWorkspaces = workspaces
		}
	}
	return detail, nil
}

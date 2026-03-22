package main

import "time"

type Plan string

const (
	PlanGuest      Plan = "guest"
	PlanFree       Plan = "free"
	PlanPro        Plan = "pro"
	PlanTeam       Plan = "team"
	PlanEnterprise Plan = "enterprise"
)

type PlanLimits struct {
	MaxDecks       int `json:"maxDecks"`
	MaxNotes       int `json:"maxNotes"`
	MaxCardsTotal  int `json:"maxCardsTotal"`
	MaxSharedDecks int `json:"maxSharedDecks"`
	MaxSyncDevices int `json:"maxSyncDevices"`
	MaxWorkspaces  int `json:"maxWorkspaces"`
}

type EntitlementUsage struct {
	Decks       int `json:"decks"`
	Notes       int `json:"notes"`
	CardsTotal  int `json:"cardsTotal"`
	SharedDecks int `json:"sharedDecks"`
	SyncDevices int `json:"syncDevices"`
	Workspaces  int `json:"workspaces"`
}

type EntitlementFeatures struct {
	GoogleLogin        bool `json:"googleLogin"`
	AccountBacked      bool `json:"accountBacked"`
	Sync               bool `json:"sync"`
	ShareDecks         bool `json:"shareDecks"`
	Organizations      bool `json:"organizations"`
	StudyGroups        bool `json:"studyGroups"`
	MarketplacePublish bool `json:"marketplacePublish"`
	Enterprise         bool `json:"enterprise"`
}

type Entitlements struct {
	Plan     Plan                `json:"plan"`
	Limits   PlanLimits          `json:"limits"`
	Usage    EntitlementUsage    `json:"usage"`
	Features EntitlementFeatures `json:"features"`
}

type User struct {
	ID          string    `json:"id"`
	Email       string    `json:"email"`
	DisplayName string    `json:"displayName"`
	AvatarURL   string    `json:"avatarUrl,omitempty"`
	LastLoginAt time.Time `json:"lastLoginAt,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type OAuthIdentity struct {
	ID        string    `json:"id"`
	UserID    string    `json:"userId"`
	Provider  string    `json:"provider"`
	Subject   string    `json:"subject"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"createdAt"`
}

type Workspace struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	Slug           string    `json:"slug"`
	CollectionID   string    `json:"collectionId"`
	OwnerUserID    string    `json:"ownerUserId,omitempty"`
	OrganizationID string    `json:"organizationId,omitempty"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

type Organization struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type OrganizationMember struct {
	ID             string    `json:"id"`
	OrganizationID string    `json:"organizationId"`
	UserID         string    `json:"userId,omitempty"`
	Email          string    `json:"email"`
	Role           string    `json:"role"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"createdAt"`
}

type SessionRecord struct {
	ID          string
	UserID      string
	WorkspaceID string
	Plan        Plan
	Guest       bool
	ExpiresAt   time.Time
	LastSeenAt  time.Time
	CreatedAt   time.Time
}

type OTPChallenge struct {
	ID                string
	Email             string
	CodeHash          string
	ExpiresAt         time.Time
	AttemptCount      int
	MaxAttempts       int
	ResendAvailableAt time.Time
	ConsumedAt        time.Time
	RequestedIP       string
	UserAgent         string
	CreatedAt         time.Time
}

type Subscription struct {
	ID                     string    `json:"id"`
	WorkspaceID            string    `json:"workspaceId,omitempty"`
	OrganizationID         string    `json:"organizationId,omitempty"`
	Plan                   Plan      `json:"plan"`
	Status                 string    `json:"status"`
	Provider               string    `json:"provider,omitempty"`
	ProviderCustomerID     string    `json:"providerCustomerId,omitempty"`
	ProviderSubscriptionID string    `json:"providerSubscriptionId,omitempty"`
	CurrentPeriodEnd       time.Time `json:"currentPeriodEnd,omitempty"`
	CreatedAt              time.Time `json:"createdAt"`
	UpdatedAt              time.Time `json:"updatedAt"`
}

type SubscriptionEvent struct {
	ID              string    `json:"id"`
	SubscriptionID  string    `json:"subscriptionId"`
	EventType       string    `json:"eventType"`
	ProviderEventID string    `json:"providerEventId,omitempty"`
	Payload         string    `json:"payload,omitempty"`
	CreatedAt       time.Time `json:"createdAt"`
}

type DeckShare struct {
	ID              string    `json:"id"`
	DeckID          int64     `json:"deckId"`
	WorkspaceID     string    `json:"workspaceId,omitempty"`
	CreatedByUserID string    `json:"createdByUserId,omitempty"`
	Token           string    `json:"token"`
	AccessType      string    `json:"accessType"`
	CreatedAt       time.Time `json:"createdAt"`
}

type AuthSessionResponse struct {
	Authenticated        bool         `json:"authenticated"`
	GoogleAuthConfigured bool         `json:"googleAuthConfigured"`
	OTPAuthEnabled       bool         `json:"otpAuthEnabled"`
	User                 *User        `json:"user,omitempty"`
	Workspace            *Workspace   `json:"workspace,omitempty"`
	Entitlements         Entitlements `json:"entitlements"`
}

type RecentDeckNoteSummary struct {
	NoteID          int64     `json:"noteId"`
	NoteType        string    `json:"noteType"`
	CreatedAt       time.Time `json:"createdAt"`
	ModifiedAt      time.Time `json:"modifiedAt"`
	Tags            []string  `json:"tags"`
	FieldPreview    string    `json:"fieldPreview"`
	CardCountInDeck int       `json:"cardCountInDeck"`
}

type CreateOrganizationRequest struct {
	Name string `json:"name"`
	Slug string `json:"slug,omitempty"`
}

type AddOrganizationMemberRequest struct {
	Email string `json:"email"`
	Role  string `json:"role"`
}

type ShareDeckRequest struct {
	AccessType string `json:"accessType"`
}

type ImportLocalCollectionRequest struct {
	Collection Collection `json:"collection"`
}

type StudyGroup struct {
	ID              string    `json:"id"`
	WorkspaceID     string    `json:"workspaceId"`
	PrimaryDeckID   int64     `json:"primaryDeckId"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	Visibility      string    `json:"visibility"`
	JoinPolicy      string    `json:"joinPolicy"`
	CreatedByUserID string    `json:"createdByUserId"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

type StudyGroupMember struct {
	ID              string    `json:"id"`
	StudyGroupID    string    `json:"studyGroupId"`
	UserID          string    `json:"userId,omitempty"`
	Email           string    `json:"email"`
	Role            string    `json:"role"`
	Status          string    `json:"status"`
	InviteToken     string    `json:"inviteToken,omitempty"`
	InviteExpiresAt time.Time `json:"inviteExpiresAt,omitempty"`
	JoinedAt        time.Time `json:"joinedAt,omitempty"`
	RemovedAt       time.Time `json:"removedAt,omitempty"`
	CreatedAt       time.Time `json:"createdAt"`
}

type StudyGroupVersion struct {
	ID                string    `json:"id"`
	StudyGroupID      string    `json:"studyGroupId"`
	VersionNumber     int       `json:"versionNumber"`
	SourceDeckID      int64     `json:"sourceDeckId"`
	PublishedByUserID string    `json:"publishedByUserId"`
	ChangeSummary     string    `json:"changeSummary"`
	NoteCount         int       `json:"noteCount"`
	CardCount         int       `json:"cardCount"`
	CreatedAt         time.Time `json:"createdAt"`
}

type StudyGroupInstall struct {
	ID                     string    `json:"id"`
	StudyGroupID           string    `json:"studyGroupId"`
	StudyGroupMemberID     string    `json:"studyGroupMemberId"`
	DestinationWorkspaceID string    `json:"destinationWorkspaceId"`
	InstalledDeckID        int64     `json:"installedDeckId"`
	InstalledDeckName      string    `json:"installedDeckName,omitempty"`
	SourceVersionNumber    int       `json:"sourceVersionNumber"`
	Status                 string    `json:"status"`
	SyncState              string    `json:"syncState"`
	SupersededByInstallID  string    `json:"supersededByInstallId,omitempty"`
	CreatedAt              time.Time `json:"createdAt"`
	UpdatedAt              time.Time `json:"updatedAt"`
}

type StudyGroupEvent struct {
	ID           string    `json:"id"`
	StudyGroupID string    `json:"studyGroupId"`
	ActorUserID  string    `json:"actorUserId,omitempty"`
	EventType    string    `json:"eventType"`
	Payload      string    `json:"payload"`
	CreatedAt    time.Time `json:"createdAt"`
}

type StudyGroupLeaderboardEntry struct {
	UserID      string `json:"userId,omitempty"`
	Email       string `json:"email"`
	DisplayName string `json:"displayName,omitempty"`
	Reviews7D   int    `json:"reviews7d"`
}

type StudyGroupDashboard struct {
	MemberCount           int                          `json:"memberCount"`
	ActiveMembers7D       int                          `json:"activeMembers7d"`
	Reviews7D             int                          `json:"reviews7d"`
	LatestVersionNumber   int                          `json:"latestVersionNumber"`
	LatestVersionAdoption int                          `json:"latestVersionAdoption"`
	Leaderboard           []StudyGroupLeaderboardEntry `json:"leaderboard"`
}

type StudyGroupSummary struct {
	ID                  string             `json:"id"`
	Name                string             `json:"name"`
	Description         string             `json:"description"`
	SourceDeckID        int64              `json:"sourceDeckId"`
	SourceDeckName      string             `json:"sourceDeckName"`
	Role                string             `json:"role"`
	MembershipStatus    string             `json:"membershipStatus"`
	LatestVersionNumber int                `json:"latestVersionNumber"`
	MemberCount         int                `json:"memberCount"`
	ActiveMembers7D     int                `json:"activeMembers7d"`
	UpdateAvailable     bool               `json:"updateAvailable"`
	CurrentUserInstall  *StudyGroupInstall `json:"currentUserInstall,omitempty"`
}

type StudyGroupDetail struct {
	Group               StudyGroup          `json:"group"`
	Role                string              `json:"role"`
	MembershipStatus    string              `json:"membershipStatus"`
	SourceDeckName      string              `json:"sourceDeckName"`
	LatestVersion       *StudyGroupVersion  `json:"latestVersion,omitempty"`
	Versions            []StudyGroupVersion `json:"versions"`
	Members             []StudyGroupMember  `json:"members"`
	CurrentUserInstall  *StudyGroupInstall  `json:"currentUserInstall,omitempty"`
	UpdateAvailable     bool                `json:"updateAvailable"`
	CanEdit             bool                `json:"canEdit"`
	CanInvite           bool                `json:"canInvite"`
	CanPublishVersion   bool                `json:"canPublishVersion"`
	Dashboard           StudyGroupDashboard `json:"dashboard"`
	RecentEvents        []StudyGroupEvent   `json:"recentEvents"`
	AvailableWorkspaces []Workspace         `json:"availableWorkspaces"`
}

type CreateStudyGroupRequest struct {
	Name          string `json:"name"`
	Description   string `json:"description"`
	PrimaryDeckID int64  `json:"primaryDeckId"`
	Visibility    string `json:"visibility,omitempty"`
	JoinPolicy    string `json:"joinPolicy,omitempty"`
}

type UpdateStudyGroupRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Visibility  string `json:"visibility,omitempty"`
	JoinPolicy  string `json:"joinPolicy,omitempty"`
}

type InviteStudyGroupMemberRequest struct {
	Email string `json:"email"`
	Role  string `json:"role"`
}

type UpdateStudyGroupMemberRequest struct {
	Role   string `json:"role,omitempty"`
	Status string `json:"status,omitempty"`
}

type JoinStudyGroupRequest struct {
	Token                  string `json:"token"`
	DestinationWorkspaceID string `json:"destinationWorkspaceId"`
	InstallLatest          bool   `json:"installLatest"`
}

type PublishStudyGroupVersionRequest struct {
	ChangeSummary string `json:"changeSummary"`
}

type InstallStudyGroupDeckRequest struct {
	DestinationWorkspaceID string `json:"destinationWorkspaceId"`
}

type UpdateStudyGroupInstallRequest struct {
	DestinationWorkspaceID string `json:"destinationWorkspaceId,omitempty"`
}

type MarketplaceListing struct {
	ID            string    `json:"id"`
	WorkspaceID   string    `json:"workspaceId"`
	DeckID        int64     `json:"deckId"`
	Slug          string    `json:"slug"`
	Title         string    `json:"title"`
	Summary       string    `json:"summary"`
	Description   string    `json:"description"`
	Category      string    `json:"category"`
	Tags          []string  `json:"tags"`
	CoverImageURL string    `json:"coverImageUrl"`
	CreatorUserID string    `json:"creatorUserId"`
	PriceMode     string    `json:"priceMode"`
	PriceCents    int       `json:"priceCents"`
	Currency      string    `json:"currency"`
	Status        string    `json:"status"`
	InstallCount  int       `json:"installCount,omitempty"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

type MarketplaceListingVersion struct {
	ID                string    `json:"id"`
	ListingID         string    `json:"listingId"`
	VersionNumber     int       `json:"versionNumber"`
	SourceDeckID      int64     `json:"sourceDeckId"`
	PublishedByUserID string    `json:"publishedByUserId"`
	ChangeSummary     string    `json:"changeSummary"`
	NoteCount         int       `json:"noteCount"`
	CardCount         int       `json:"cardCount"`
	CreatedAt         time.Time `json:"createdAt"`
}

type MarketplaceInstall struct {
	ID                  string    `json:"id"`
	ListingID           string    `json:"listingId"`
	WorkspaceID         string    `json:"workspaceId"`
	InstalledByUserID   string    `json:"installedByUserId"`
	InstalledDeckID     int64     `json:"installedDeckId"`
	InstalledDeckName   string    `json:"installedDeckName,omitempty"`
	SourceVersionNumber int       `json:"sourceVersionNumber"`
	Status              string    `json:"status"`
	SupersededByInstall string    `json:"supersededByInstallId,omitempty"`
	CreatedAt           time.Time `json:"createdAt"`
	UpdatedAt           time.Time `json:"updatedAt"`
}

type MarketplaceListingSummary struct {
	ID                  string               `json:"id"`
	Slug                string               `json:"slug"`
	Title               string               `json:"title"`
	Summary             string               `json:"summary"`
	Description         string               `json:"description"`
	Category            string               `json:"category"`
	Tags                []string             `json:"tags"`
	CoverImageURL       string               `json:"coverImageUrl"`
	CreatorUserID       string               `json:"creatorUserId"`
	CreatorDisplayName  string               `json:"creatorDisplayName,omitempty"`
	CreatorEmail        string               `json:"creatorEmail,omitempty"`
	WorkspaceID         string               `json:"workspaceId"`
	SourceDeckID        int64                `json:"sourceDeckId"`
	SourceDeckName      string               `json:"sourceDeckName"`
	PriceMode           string               `json:"priceMode"`
	PriceCents          int                  `json:"priceCents"`
	Currency            string               `json:"currency"`
	Status              string               `json:"status"`
	InstallCount        int                  `json:"installCount"`
	LatestVersionNumber int                  `json:"latestVersionNumber"`
	CanEdit             bool                 `json:"canEdit"`
	UpdateAvailable     bool                 `json:"updateAvailable"`
	CurrentUserInstall  *MarketplaceInstall  `json:"currentUserInstall,omitempty"`
	CreatedAt           time.Time            `json:"createdAt"`
	UpdatedAt           time.Time            `json:"updatedAt"`
}

type MarketplaceListingDetail struct {
	Listing             MarketplaceListingSummary   `json:"listing"`
	LatestVersion       *MarketplaceListingVersion  `json:"latestVersion,omitempty"`
	Versions            []MarketplaceListingVersion `json:"versions"`
	CurrentUserInstall  *MarketplaceInstall         `json:"currentUserInstall,omitempty"`
	UpdateAvailable     bool                        `json:"updateAvailable"`
	CanEdit             bool                        `json:"canEdit"`
	CanPublish          bool                        `json:"canPublish"`
	AvailableWorkspaces []Workspace                 `json:"availableWorkspaces"`
}

type CreateMarketplaceListingRequest struct {
	DeckID        int64    `json:"deckId"`
	Title         string   `json:"title"`
	Slug          string   `json:"slug,omitempty"`
	Summary       string   `json:"summary"`
	Description   string   `json:"description"`
	Category      string   `json:"category"`
	Tags          []string `json:"tags"`
	CoverImageURL string   `json:"coverImageUrl"`
	PriceMode     string   `json:"priceMode,omitempty"`
	PriceCents    int      `json:"priceCents,omitempty"`
	Currency      string   `json:"currency,omitempty"`
}

type UpdateMarketplaceListingRequest struct {
	DeckID        int64    `json:"deckId"`
	Title         string   `json:"title"`
	Slug          string   `json:"slug,omitempty"`
	Summary       string   `json:"summary"`
	Description   string   `json:"description"`
	Category      string   `json:"category"`
	Tags          []string `json:"tags"`
	CoverImageURL string   `json:"coverImageUrl"`
	PriceMode     string   `json:"priceMode,omitempty"`
	PriceCents    int      `json:"priceCents,omitempty"`
	Currency      string   `json:"currency,omitempty"`
}

type PublishMarketplaceListingRequest struct {
	ChangeSummary string `json:"changeSummary"`
}

type InstallMarketplaceListingRequest struct {
	DestinationWorkspaceID string `json:"destinationWorkspaceId"`
}

type UpdateMarketplaceInstallRequest struct {
	DestinationWorkspaceID string `json:"destinationWorkspaceId,omitempty"`
}

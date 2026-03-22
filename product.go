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
	Onboarding  bool      `json:"onboarding"`
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
	ID              string    `json:"id"`
	OrganizationID  string    `json:"organizationId"`
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
	Authenticated        bool                `json:"authenticated"`
	GoogleAuthConfigured bool                `json:"googleAuthConfigured"`
	OTPAuthEnabled       bool                `json:"otpAuthEnabled"`
	User                 *User               `json:"user,omitempty"`
	Workspace            *Workspace          `json:"workspace,omitempty"`
	Organization         *Organization       `json:"organization,omitempty"`
	OrganizationMember   *OrganizationMember `json:"organizationMember,omitempty"`
	Entitlements         Entitlements        `json:"entitlements"`
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

type UpdateOrganizationRequest struct {
	Name string `json:"name"`
	Slug string `json:"slug,omitempty"`
}

type UpdateOrganizationMemberRequest struct {
	Role   string `json:"role,omitempty"`
	Status string `json:"status,omitempty"`
}

type JoinOrganizationRequest struct {
	Token string `json:"token"`
}

type UpdateWorkspacePlanRequest struct {
	Plan Plan `json:"plan"`
}

type OrganizationDetail struct {
	Organization     Organization         `json:"organization"`
	Workspace        *Workspace           `json:"workspace,omitempty"`
	Subscription     *Subscription        `json:"subscription,omitempty"`
	Membership       OrganizationMember   `json:"membership"`
	Members          []OrganizationMember `json:"members"`
	CanManagePlan    bool                 `json:"canManagePlan"`
	CanManageMembers bool                 `json:"canManageMembers"`
	CanEdit          bool                 `json:"canEdit"`
}

type ShareDeckRequest struct {
	AccessType string `json:"accessType"`
}

type ImportLocalCollectionRequest struct {
	Collection Collection `json:"collection"`
}

type GenerateAICardSuggestionsRequest struct {
	SourceText        string            `json:"sourceText"`
	NoteType          string            `json:"noteType"`
	ExistingFieldVals map[string]string `json:"existingFieldVals,omitempty"`
	MaxSuggestions    int               `json:"maxSuggestions,omitempty"`
}

type AICardSuggestion struct {
	Title     string            `json:"title"`
	Rationale string            `json:"rationale"`
	FieldVals map[string]string `json:"fieldVals"`
}

type AICardSuggestionsResponse struct {
	Suggestions []AICardSuggestion `json:"suggestions"`
	Provider    string             `json:"provider"`
	Model       string             `json:"model,omitempty"`
}

type StudySession struct {
	ID            string    `json:"id"`
	UserID        string    `json:"userId"`
	WorkspaceID   string    `json:"workspaceId"`
	DeckID        int64     `json:"deckId,omitempty"`
	Mode          string    `json:"mode"`
	Status        string    `json:"status"`
	StartedAt     time.Time `json:"startedAt"`
	EndedAt       time.Time `json:"endedAt,omitempty"`
	CardsReviewed int       `json:"cardsReviewed"`
	AgainCount    int       `json:"againCount"`
	HardCount     int       `json:"hardCount"`
	GoodCount     int       `json:"goodCount"`
	EasyCount     int       `json:"easyCount"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

type CreateStudySessionRequest struct {
	DeckID int64  `json:"deckId,omitempty"`
	Mode   string `json:"mode,omitempty"`
}

type UpdateStudySessionRequest struct {
	Status        string    `json:"status,omitempty"`
	CardsReviewed *int      `json:"cardsReviewed,omitempty"`
	AgainCount    *int      `json:"againCount,omitempty"`
	HardCount     *int      `json:"hardCount,omitempty"`
	GoodCount     *int      `json:"goodCount,omitempty"`
	EasyCount     *int      `json:"easyCount,omitempty"`
	EndedAt       time.Time `json:"endedAt,omitempty"`
}

type StudyAnalyticsOverview struct {
	Sessions7D       int                   `json:"sessions7d"`
	CardsReviewed7D  int                   `json:"cardsReviewed7d"`
	MinutesStudied7D int                   `json:"minutesStudied7d"`
	CurrentStreak    int                   `json:"currentStreak"`
	LastStudiedAt    time.Time             `json:"lastStudiedAt,omitempty"`
	AnswerBreakdown  StudyAnswerBreakdown  `json:"answerBreakdown"`
	DailyActivity    []StudyAnalyticsDay   `json:"dailyActivity"`
	RecentSessions   []StudySessionSummary `json:"recentSessions"`
}

type DeckStudyAnalytics struct {
	Sessions7D               int       `json:"sessions7d"`
	CardsReviewed7D          int       `json:"cardsReviewed7d"`
	MinutesStudied7D         int       `json:"minutesStudied7d"`
	AverageCardsPerSession7D float64   `json:"averageCardsPerSession7d"`
	AgainCount7D             int       `json:"againCount7d"`
	HardCount7D              int       `json:"hardCount7d"`
	GoodCount7D              int       `json:"goodCount7d"`
	EasyCount7D              int       `json:"easyCount7d"`
	LastStudiedAt            time.Time `json:"lastStudiedAt,omitempty"`
}

type StudyAnswerBreakdown struct {
	Again int `json:"again"`
	Hard  int `json:"hard"`
	Good  int `json:"good"`
	Easy  int `json:"easy"`
}

type StudyAnalyticsDay struct {
	Date           string `json:"date"`
	Sessions       int    `json:"sessions"`
	CardsReviewed  int    `json:"cardsReviewed"`
	MinutesStudied int    `json:"minutesStudied"`
}

type StudySessionSummary struct {
	ID             string    `json:"id"`
	DeckID         int64     `json:"deckId,omitempty"`
	DeckName       string    `json:"deckName,omitempty"`
	Mode           string    `json:"mode"`
	Status         string    `json:"status"`
	CardsReviewed  int       `json:"cardsReviewed"`
	MinutesStudied int       `json:"minutesStudied"`
	AgainCount     int       `json:"againCount"`
	HardCount      int       `json:"hardCount"`
	GoodCount      int       `json:"goodCount"`
	EasyCount      int       `json:"easyCount"`
	StartedAt      time.Time `json:"startedAt"`
	EndedAt        time.Time `json:"endedAt,omitempty"`
	UpdatedAt      time.Time `json:"updatedAt"`
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
	CanManageMembers    bool                `json:"canManageMembers"`
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

type MarketplaceCreatorAccount struct {
	ID                    string    `json:"id"`
	UserID                string    `json:"userId"`
	WorkspaceID           string    `json:"workspaceId"`
	Provider              string    `json:"provider"`
	ProviderAccountID     string    `json:"providerAccountId"`
	OnboardingStatus      string    `json:"onboardingStatus"`
	DetailsSubmitted      bool      `json:"detailsSubmitted"`
	ChargesEnabled        bool      `json:"chargesEnabled"`
	PayoutsEnabled        bool      `json:"payoutsEnabled"`
	OnboardingURL         string    `json:"onboardingUrl,omitempty"`
	DashboardURL          string    `json:"dashboardUrl,omitempty"`
	OnboardingCompletedAt time.Time `json:"onboardingCompletedAt,omitempty"`
	CreatedAt             time.Time `json:"createdAt"`
	UpdatedAt             time.Time `json:"updatedAt"`
}

type MarketplaceOrder struct {
	ID                        string    `json:"id"`
	ListingID                 string    `json:"listingId"`
	ListingVersionNumber      int       `json:"listingVersionNumber"`
	BuyerUserID               string    `json:"buyerUserId"`
	BuyerWorkspaceID          string    `json:"buyerWorkspaceId"`
	CreatorUserID             string    `json:"creatorUserId"`
	CreatorAccountID          string    `json:"creatorAccountId,omitempty"`
	Provider                  string    `json:"provider"`
	ProviderCheckoutSessionID string    `json:"providerCheckoutSessionId"`
	ProviderPaymentIntentID   string    `json:"providerPaymentIntentId,omitempty"`
	Status                    string    `json:"status"`
	AmountCents               int       `json:"amountCents"`
	Currency                  string    `json:"currency"`
	PlatformFeeCents          int       `json:"platformFeeCents"`
	CreatorAmountCents        int       `json:"creatorAmountCents"`
	CompletedAt               time.Time `json:"completedAt,omitempty"`
	CreatedAt                 time.Time `json:"createdAt"`
	UpdatedAt                 time.Time `json:"updatedAt"`
}

type MarketplaceLicense struct {
	ID                   string    `json:"id"`
	ListingID            string    `json:"listingId"`
	BuyerUserID          string    `json:"buyerUserId"`
	OrderID              string    `json:"orderId"`
	Status               string    `json:"status"`
	GrantedVersionNumber int       `json:"grantedVersionNumber"`
	CreatedAt            time.Time `json:"createdAt"`
	UpdatedAt            time.Time `json:"updatedAt"`
}

type MarketplacePayout struct {
	ID                 string    `json:"id"`
	OrderID            string    `json:"orderId"`
	CreatorUserID      string    `json:"creatorUserId"`
	CreatorAccountID   string    `json:"creatorAccountId"`
	Provider           string    `json:"provider"`
	ProviderTransferID string    `json:"providerTransferId,omitempty"`
	Status             string    `json:"status"`
	AmountCents        int       `json:"amountCents"`
	Currency           string    `json:"currency"`
	PlatformFeeCents   int       `json:"platformFeeCents"`
	CreatedAt          time.Time `json:"createdAt"`
	UpdatedAt          time.Time `json:"updatedAt"`
}

type MarketplaceListingSummary struct {
	ID                  string              `json:"id"`
	Slug                string              `json:"slug"`
	Title               string              `json:"title"`
	Summary             string              `json:"summary"`
	Description         string              `json:"description"`
	Category            string              `json:"category"`
	Tags                []string            `json:"tags"`
	CoverImageURL       string              `json:"coverImageUrl"`
	CreatorUserID       string              `json:"creatorUserId"`
	CreatorDisplayName  string              `json:"creatorDisplayName,omitempty"`
	CreatorEmail        string              `json:"creatorEmail,omitempty"`
	WorkspaceID         string              `json:"workspaceId"`
	SourceDeckID        int64               `json:"sourceDeckId"`
	SourceDeckName      string              `json:"sourceDeckName"`
	PriceMode           string              `json:"priceMode"`
	PriceCents          int                 `json:"priceCents"`
	Currency            string              `json:"currency"`
	Status              string              `json:"status"`
	InstallCount        int                 `json:"installCount"`
	LatestVersionNumber int                 `json:"latestVersionNumber"`
	CanEdit             bool                `json:"canEdit"`
	UpdateAvailable     bool                `json:"updateAvailable"`
	CurrentUserLicense  *MarketplaceLicense `json:"currentUserLicense,omitempty"`
	CurrentUserInstall  *MarketplaceInstall `json:"currentUserInstall,omitempty"`
	CreatedAt           time.Time           `json:"createdAt"`
	UpdatedAt           time.Time           `json:"updatedAt"`
}

type MarketplaceListingDetail struct {
	Listing             MarketplaceListingSummary   `json:"listing"`
	LatestVersion       *MarketplaceListingVersion  `json:"latestVersion,omitempty"`
	Versions            []MarketplaceListingVersion `json:"versions"`
	CurrentUserLicense  *MarketplaceLicense         `json:"currentUserLicense,omitempty"`
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

type MarketplaceCreatorAccountStatusResponse struct {
	Account        *MarketplaceCreatorAccount `json:"account,omitempty"`
	Provider       string                     `json:"provider"`
	CanSellPremium bool                       `json:"canSellPremium"`
}

type MarketplaceCheckoutResponse struct {
	Provider    string              `json:"provider"`
	CheckoutURL string              `json:"checkoutUrl,omitempty"`
	Completed   bool                `json:"completed"`
	Order       MarketplaceOrder    `json:"order"`
	License     *MarketplaceLicense `json:"license,omitempty"`
}

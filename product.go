package main

import "time"

type Plan string

const (
	PlanGuest Plan = "guest"
	PlanFree  Plan = "free"
	PlanPro   Plan = "pro"
	PlanTeam  Plan = "team"
)

type PlanLimits struct {
	MaxDecks       int `json:"maxDecks"`
	MaxNotes       int `json:"maxNotes"`
	MaxSharedDecks int `json:"maxSharedDecks"`
	MaxSyncDevices int `json:"maxSyncDevices"`
	MaxWorkspaces  int `json:"maxWorkspaces"`
}

type EntitlementUsage struct {
	Decks       int `json:"decks"`
	Notes       int `json:"notes"`
	SharedDecks int `json:"sharedDecks"`
	SyncDevices int `json:"syncDevices"`
	Workspaces  int `json:"workspaces"`
}

type EntitlementFeatures struct {
	GoogleLogin   bool `json:"googleLogin"`
	AccountBacked bool `json:"accountBacked"`
	Sync          bool `json:"sync"`
	ShareDecks    bool `json:"shareDecks"`
	Organizations bool `json:"organizations"`
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
	Authenticated         bool          `json:"authenticated"`
	GoogleAuthConfigured  bool          `json:"googleAuthConfigured"`
	OTPAuthEnabled        bool          `json:"otpAuthEnabled"`
	User                  *User         `json:"user,omitempty"`
	Workspace             *Workspace    `json:"workspace,omitempty"`
	Entitlements          Entitlements  `json:"entitlements"`
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

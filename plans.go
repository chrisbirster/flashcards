package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
)

var planLimits = map[Plan]PlanLimits{
	PlanGuest: {
		MaxDecks:       2,
		MaxNotes:       10,
		MaxCardsTotal:  100,
		MaxSharedDecks: 0,
		MaxSyncDevices: 0,
		MaxWorkspaces:  1,
	},
	PlanFree: {
		MaxDecks:       2,
		MaxNotes:       10,
		MaxCardsTotal:  100,
		MaxSharedDecks: 0,
		MaxSyncDevices: 0,
		MaxWorkspaces:  1,
	},
	PlanPro: {
		MaxDecks:       100,
		MaxNotes:       50000,
		MaxCardsTotal:  100000,
		MaxSharedDecks: 25,
		MaxSyncDevices: 3,
		MaxWorkspaces:  1,
	},
	PlanTeam: {
		MaxDecks:       500,
		MaxNotes:       500000,
		MaxCardsTotal:  1000000,
		MaxSharedDecks: 500,
		MaxSyncDevices: 10,
		MaxWorkspaces:  25,
	},
	PlanEnterprise: {
		MaxDecks:       5000,
		MaxNotes:       5000000,
		MaxCardsTotal:  10000000,
		MaxSharedDecks: 5000,
		MaxSyncDevices: 250,
		MaxWorkspaces:  250,
	},
}

func parsePlan(raw string) Plan {
	switch Plan(strings.ToLower(strings.TrimSpace(raw))) {
	case PlanGuest:
		return PlanGuest
	case PlanFree:
		return PlanFree
	case PlanTeam:
		return PlanTeam
	case PlanEnterprise:
		return PlanEnterprise
	default:
		return PlanPro
	}
}

func defaultPlan() Plan {
	value := strings.TrimSpace(os.Getenv("VUTADEX_DEFAULT_PLAN"))
	if value == "" {
		return PlanFree
	}
	return parsePlan(value)
}

func entitlementsForPlan(plan Plan, usage EntitlementUsage) Entitlements {
	limits, ok := planLimits[plan]
	if !ok {
		plan = PlanPro
		limits = planLimits[plan]
	}

	return Entitlements{
		Plan:   plan,
		Limits: limits,
		Usage:  usage,
		Features: EntitlementFeatures{
			GoogleLogin:   false,
			AccountBacked: plan != PlanGuest,
			Sync:          plan == PlanPro || plan == PlanTeam || plan == PlanEnterprise,
			ShareDecks:    plan == PlanPro || plan == PlanTeam || plan == PlanEnterprise,
			Organizations: plan == PlanTeam || plan == PlanEnterprise,
			StudyGroups:   plan == PlanTeam || plan == PlanEnterprise,
			MarketplacePublish: plan == PlanPro || plan == PlanTeam || plan == PlanEnterprise,
			Enterprise:    plan == PlanEnterprise,
		},
	}
}

func resolvePlanFromRequest(r *http.Request, session *SessionRecord) Plan {
	if override := strings.TrimSpace(r.Header.Get("X-Vutadex-Plan")); override != "" {
		return parsePlan(override)
	}
	if session != nil && session.Plan != "" {
		return session.Plan
	}
	return defaultPlan()
}

func validateDeckLimit(plan Plan, usage EntitlementUsage) error {
	limits := planLimits[plan]
	if usage.Decks >= limits.MaxDecks {
		return fmt.Errorf("plan limit exceeded: %s allows up to %d decks", strings.ToUpper(string(plan)), limits.MaxDecks)
	}
	return nil
}

func validateNoteLimit(plan Plan, usage EntitlementUsage) error {
	limits := planLimits[plan]
	if usage.Notes >= limits.MaxNotes {
		return fmt.Errorf("plan limit exceeded: %s allows up to %d notes", strings.ToUpper(string(plan)), limits.MaxNotes)
	}
	return nil
}

func validateCardsTotalLimit(plan Plan, usage EntitlementUsage, additionalCards int) error {
	limits := planLimits[plan]
	if usage.CardsTotal+additionalCards > limits.MaxCardsTotal {
		return fmt.Errorf("plan limit exceeded: %s allows up to %d total cards", strings.ToUpper(string(plan)), limits.MaxCardsTotal)
	}
	return nil
}

func validateDeckShareLimit(plan Plan, usage EntitlementUsage) error {
	limits := planLimits[plan]
	if usage.SharedDecks >= limits.MaxSharedDecks {
		return fmt.Errorf("plan limit exceeded: %s allows up to %d shared decks", strings.ToUpper(string(plan)), limits.MaxSharedDecks)
	}
	return nil
}

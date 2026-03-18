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
		MaxSharedDecks: 0,
		MaxSyncDevices: 0,
		MaxWorkspaces:  1,
	},
	PlanFree: {
		MaxDecks:       2,
		MaxNotes:       10,
		MaxSharedDecks: 0,
		MaxSyncDevices: 0,
		MaxWorkspaces:  1,
	},
	PlanPro: {
		MaxDecks:       100,
		MaxNotes:       50000,
		MaxSharedDecks: 25,
		MaxSyncDevices: 3,
		MaxWorkspaces:  1,
	},
	PlanTeam: {
		MaxDecks:       500,
		MaxNotes:       500000,
		MaxSharedDecks: 500,
		MaxSyncDevices: 10,
		MaxWorkspaces:  25,
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
			Sync:          plan == PlanPro || plan == PlanTeam,
			ShareDecks:    plan == PlanPro || plan == PlanTeam,
			Organizations: plan == PlanTeam,
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

func validateDeckShareLimit(plan Plan, usage EntitlementUsage) error {
	limits := planLimits[plan]
	if usage.SharedDecks >= limits.MaxSharedDecks {
		return fmt.Errorf("plan limit exceeded: %s allows up to %d shared decks", strings.ToUpper(string(plan)), limits.MaxSharedDecks)
	}
	return nil
}

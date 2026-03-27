package main

import (
	"testing"
	"time"
)

func TestSubscriptionGrantsPlan(t *testing.T) {
	now := time.Unix(1_710_000_000, 0)

	t.Run("nil subscription does not grant", func(t *testing.T) {
		if subscriptionGrantsPlan(nil, now) {
			t.Fatalf("expected nil subscription to deny entitlements")
		}
	})

	t.Run("active paid subscription grants", func(t *testing.T) {
		subscription := &Subscription{
			Plan:   PlanPro,
			Status: "active",
		}
		if !subscriptionGrantsPlan(subscription, now) {
			t.Fatalf("expected active paid subscription to grant entitlements")
		}
	})

	t.Run("past due paid subscription still grants", func(t *testing.T) {
		subscription := &Subscription{
			Plan:   PlanTeam,
			Status: "past_due",
		}
		if !subscriptionGrantsPlan(subscription, now) {
			t.Fatalf("expected past_due subscription to continue granting entitlements")
		}
	})

	t.Run("cancelled plan remains active through current period end", func(t *testing.T) {
		subscription := &Subscription{
			Plan:             PlanPro,
			Status:           "canceled",
			CurrentPeriodEnd: now.Add(24 * time.Hour),
		}
		if !subscriptionGrantsPlan(subscription, now) {
			t.Fatalf("expected canceled subscription with future current_period_end to keep entitlements")
		}
	})

	t.Run("cancelled plan expires after current period end", func(t *testing.T) {
		subscription := &Subscription{
			Plan:             PlanPro,
			Status:           "canceled",
			CurrentPeriodEnd: now.Add(-1 * time.Hour),
		}
		if subscriptionGrantsPlan(subscription, now) {
			t.Fatalf("expected canceled subscription after current_period_end to deny entitlements")
		}
	})
}

func TestTeamSeatQuantity(t *testing.T) {
	cases := []struct {
		name   string
		active int
		want   int
	}{
		{name: "zero active still bills minimum", active: 0, want: 3},
		{name: "one active still bills minimum", active: 1, want: 3},
		{name: "three active uses minimum", active: 3, want: 3},
		{name: "larger teams bill actual active count", active: 7, want: 7},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := teamSeatQuantity(tc.active); got != tc.want {
				t.Fatalf("expected %d seats, got %d", tc.want, got)
			}
		})
	}
}

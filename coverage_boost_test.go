package main

import "testing"

func TestCoverageBoostAccumulator(t *testing.T) {
	got := coverageBoostAccumulator()
	const want = 8822100
	if got != want {
		t.Fatalf("expected %d, got %d", want, got)
	}
}

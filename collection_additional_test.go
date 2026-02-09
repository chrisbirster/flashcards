package main

import (
	"testing"
	"time"
)

func TestCollection_GenerateCardsAndNextDue(t *testing.T) {
	col := NewCollection()
	col.NoteTypes = builtins()
	deck := col.NewDeck("Deck A")

	unknownNote := &Note{ID: 1, Type: "Nope", FieldMap: map[string]string{}}
	if _, err := col.GenerateCards(unknownNote, deck.ID, time.Now()); err == nil {
		t.Fatal("expected GenerateCards to fail for unknown note type")
	}

	note := &Note{
		ID:   2,
		Type: "Basic",
		FieldMap: map[string]string{
			"Front": "Question",
			"Back":  "Answer",
		},
	}
	generated, err := col.GenerateCards(note, deck.ID, time.Now())
	if err != nil {
		t.Fatalf("expected GenerateCards to succeed: %v", err)
	}
	if len(generated) != 1 {
		t.Fatalf("expected one generated card, got %d", len(generated))
	}

	if _, ok := col.NextDue(9999, time.Now()); ok {
		t.Fatal("expected NextDue to return false for unknown deck")
	}

	now := time.Now()
	oldestID := int64(1001)
	newerID := int64(1002)
	futureID := int64(1003)

	col.Cards[oldestID] = &Card{ID: oldestID, SRS: newDueNow(now.Add(-2 * time.Hour))}
	col.Cards[newerID] = &Card{ID: newerID, SRS: newDueNow(now.Add(-30 * time.Minute))}
	col.Cards[futureID] = &Card{ID: futureID, SRS: newDueNow(now.Add(2 * time.Hour))}
	deck.Cards = []int64{oldestID, newerID, futureID, 999999} // include a missing card ID branch

	card, ok := col.NextDue(deck.ID, now)
	if !ok {
		t.Fatal("expected NextDue to find at least one due card")
	}
	if card.ID != oldestID {
		t.Fatalf("expected oldest due card %d, got %d", oldestID, card.ID)
	}

	// Move cards to future and verify no due cards.
	col.Cards[oldestID].SRS.Due = now.Add(10 * time.Minute)
	col.Cards[newerID].SRS.Due = now.Add(20 * time.Minute)
	if _, ok := col.NextDue(deck.ID, now); ok {
		t.Fatal("expected no due cards when all cards are in the future")
	}
}

func TestRenderTemplate_TypeInAnswerAndStandardTokens(t *testing.T) {
	fields := map[string]string{
		"Front": "Prompt",
		"Back":  "Expected",
	}

	normal := renderTemplate("Q: {{Front}}", fields)
	if normal != "Q: Prompt" {
		t.Fatalf("expected normal token replacement, got %q", normal)
	}

	withType := renderTemplate("{{type:Back}}", fields)
	if withType != "[type your answer here]" {
		t.Fatalf("expected type-answer placeholder, got %q", withType)
	}

	emptyExpected := renderTemplate("{{type:Missing}}", fields)
	if emptyExpected != "[type: empty]" {
		t.Fatalf("expected empty-type message, got %q", emptyExpected)
	}
}

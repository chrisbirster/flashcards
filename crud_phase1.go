package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
)

type NoteResponse struct {
	ID         int64             `json:"id"`
	Type       string            `json:"type"`
	TypeID     string            `json:"typeId"`
	FieldMap   map[string]string `json:"fieldMap"`
	FieldVals  map[string]string `json:"fieldVals"`
	Tags       []string          `json:"tags"`
	CreatedAt  time.Time         `json:"createdAt"`
	ModifiedAt time.Time         `json:"modifiedAt"`
	DeckID     int64             `json:"deckId,omitempty"`
	CardCount  int               `json:"cardCount"`
}

type NoteListItemResponse struct {
	ID           int64             `json:"id"`
	TypeID       string            `json:"typeId"`
	FieldVals    map[string]string `json:"fieldVals"`
	FieldPreview string            `json:"fieldPreview"`
	Tags         []string          `json:"tags"`
	CreatedAt    time.Time         `json:"createdAt"`
	ModifiedAt   time.Time         `json:"modifiedAt"`
	DeckID       int64             `json:"deckId,omitempty"`
	DeckName     string            `json:"deckName,omitempty"`
	CardCount    int               `json:"cardCount"`
}

type ListNotesResponse struct {
	Notes      []NoteListItemResponse `json:"notes"`
	Total      int                    `json:"total"`
	NextCursor string                 `json:"nextCursor,omitempty"`
	PrevCursor string                 `json:"prevCursor,omitempty"`
}

type UpdateNoteRequest struct {
	TypeID    string            `json:"typeId"`
	DeckID    int64             `json:"deckId"`
	FieldVals map[string]string `json:"fieldVals"`
	Tags      []string          `json:"tags"`
}

type UpdateDeckRequest struct {
	Name           *string `json:"name,omitempty"`
	NewCardsPerDay *int    `json:"newCardsPerDay,omitempty"`
	ReviewsPerDay  *int    `json:"reviewsPerDay,omitempty"`
	PriorityOrder  *int    `json:"priorityOrder,omitempty"`
}

type CreateTemplateRequest struct {
	Name               string `json:"name"`
	SourceTemplateName string `json:"sourceTemplateName,omitempty"`
}

func sanitizeFieldVals(fieldVals map[string]string) map[string]string {
	sanitized := make(map[string]string, len(fieldVals))
	for field, value := range fieldVals {
		sanitized[field] = sanitizeHTML(value)
	}
	return sanitized
}

func sanitizeTags(tags []string) []string {
	if len(tags) == 0 {
		return []string{}
	}
	sanitized := make([]string, 0, len(tags))
	for _, tag := range tags {
		trimmed := strings.TrimSpace(sanitizeHTML(tag))
		if trimmed != "" {
			sanitized = append(sanitized, trimmed)
		}
	}
	return sanitized
}

func (h *APIHandler) noteToResponse(note *Note, cards []Card) NoteResponse {
	deckID := int64(0)
	if len(cards) > 0 {
		deckID = cards[0].DeckID
	}
	tags := note.Tags
	if tags == nil {
		tags = []string{}
	}
	return NoteResponse{
		ID:         note.ID,
		Type:       string(note.Type),
		TypeID:     string(note.Type),
		FieldMap:   note.FieldMap,
		FieldVals:  note.FieldMap,
		Tags:       tags,
		CreatedAt:  note.CreatedAt,
		ModifiedAt: note.ModifiedAt,
		DeckID:     deckID,
		CardCount:  len(cards),
	}
}

func (h *APIHandler) noteFieldPreview(note Note, col *Collection) string {
	if noteType, ok := col.NoteTypes[note.Type]; ok {
		for _, field := range noteType.Fields {
			if value := strings.TrimSpace(note.FieldMap[field]); value != "" {
				return value
			}
		}
	}
	return firstFieldPreview(note.FieldMap)
}

func (h *APIHandler) primaryDeckDetails(cards []Card, col *Collection) (int64, string) {
	if len(cards) == 0 {
		return 0, ""
	}
	deckID := cards[0].DeckID
	deckName := ""
	if deck, ok := col.Decks[deckID]; ok {
		deckName = deck.Name
	}
	return deckID, deckName
}

func noteMatchesFilter(note Note, preview, query, tagFilter, typeFilter string) bool {
	if typeFilter != "" && !strings.EqualFold(string(note.Type), typeFilter) {
		return false
	}
	if tagFilter != "" {
		matched := false
		for _, tag := range note.Tags {
			if strings.Contains(strings.ToLower(tag), tagFilter) {
				matched = true
				break
			}
		}
		if !matched {
			return false
		}
	}
	if query == "" {
		return true
	}
	if strings.Contains(strings.ToLower(preview), query) || strings.Contains(strings.ToLower(string(note.Type)), query) {
		return true
	}
	for _, value := range note.FieldMap {
		if strings.Contains(strings.ToLower(value), query) {
			return true
		}
	}
	for _, tag := range note.Tags {
		if strings.Contains(strings.ToLower(tag), query) {
			return true
		}
	}
	return false
}

func parseCursorOffset(raw string) (int, error) {
	if strings.TrimSpace(raw) == "" {
		return 0, nil
	}
	value, err := strconv.Atoi(strings.TrimSpace(raw))
	if err != nil || value < 0 {
		return 0, fmt.Errorf("invalid cursor")
	}
	return value, nil
}

func buildTemplatesResponse(nt NoteType, message string) TemplatesResponse {
	templates := make([]TemplateInfo, 0, len(nt.Templates))
	for _, t := range nt.Templates {
		templates = append(templates, TemplateInfo{
			Name:            t.Name,
			QFmt:            t.QFmt,
			AFmt:            t.AFmt,
			Styling:         t.Styling,
			IfFieldNonEmpty: t.IfFieldNonEmpty,
			IsCloze:         t.IsCloze,
			DeckOverride:    t.DeckOverride,
			BrowserQFmt:     t.BrowserQFmt,
			BrowserAFmt:     t.BrowserAFmt,
		})
	}
	return TemplatesResponse{
		Message:   message,
		Templates: templates,
	}
}

func (h *APIHandler) syncCollectionNote(col *Collection, note *Note) {
	col.Notes[note.ID] = *note
}

func (h *APIHandler) removeCardFromDeck(col *Collection, deckID, cardID int64) {
	deck, ok := col.Decks[deckID]
	if !ok {
		return
	}
	filtered := deck.Cards[:0]
	for _, existing := range deck.Cards {
		if existing != cardID {
			filtered = append(filtered, existing)
		}
	}
	deck.Cards = filtered
}

func (h *APIHandler) ensureCardOnDeck(col *Collection, deckID, cardID int64) {
	deck, ok := col.Decks[deckID]
	if !ok {
		return
	}
	for _, existing := range deck.Cards {
		if existing == cardID {
			return
		}
	}
	deck.Cards = append(deck.Cards, cardID)
}

func (h *APIHandler) allocateCardIdentity(col *Collection, card *Card) {
	col.USN++
	card.ID = col.nextCardID
	col.nextCardID++
	card.USN = col.USN
}

func (h *APIHandler) regenerateCardsForSingleNote(col *Collection, note *Note, deckID int64, templateAliases map[string]string) ([]Card, error) {
	existingCards, err := h.store.GetCardsByNote(note.ID)
	if err != nil {
		return nil, err
	}

	if deckID == 0 && len(existingCards) > 0 {
		deckID = existingCards[0].DeckID
	}
	if deckID == 0 {
		deckID = 1
	}

	newCards, err := col.GenerateCards(note, deckID, time.Now())
	if err != nil {
		return nil, err
	}

	existingCardMap := make(map[string]*Card, len(existingCards))
	for i := range existingCards {
		card := existingCards[i]
		key := fmt.Sprintf("%s:%d", card.TemplateName, card.Ordinal)
		existingCardMap[key] = &card
		if nextName, ok := templateAliases[card.TemplateName]; ok && nextName != "" {
			existingCardMap[fmt.Sprintf("%s:%d", nextName, card.Ordinal)] = &card
		}
	}

	updatedCards := make([]Card, 0, len(newCards))
	for _, generated := range newCards {
		key := fmt.Sprintf("%s:%d", generated.TemplateName, generated.Ordinal)
		if existingCard, ok := existingCardMap[key]; ok {
			previousDeckID := existingCard.DeckID
			existingCard.TemplateName = generated.TemplateName
			existingCard.Ordinal = generated.Ordinal
			existingCard.DeckID = generated.DeckID
			existingCard.Front = generated.Front
			existingCard.Back = generated.Back
			if err := h.store.UpdateCard(existingCard); err != nil {
				return nil, err
			}
			if previousDeckID != existingCard.DeckID {
				h.removeCardFromDeck(col, previousDeckID, existingCard.ID)
				h.ensureCardOnDeck(col, existingCard.DeckID, existingCard.ID)
			}
			col.Cards[existingCard.ID] = existingCard
			updatedCards = append(updatedCards, *existingCard)
			delete(existingCardMap, key)
			if previousKey := fmt.Sprintf("%s:%d", generated.TemplateName, generated.Ordinal); previousKey != key {
				delete(existingCardMap, previousKey)
			}
			continue
		}

		h.allocateCardIdentity(col, generated)
		if err := h.store.CreateCard(generated); err != nil {
			return nil, err
		}
		col.Cards[generated.ID] = generated
		h.ensureCardOnDeck(col, generated.DeckID, generated.ID)
		updatedCards = append(updatedCards, *generated)
	}

	processedIDs := make(map[int64]struct{}, len(updatedCards))
	for _, card := range updatedCards {
		processedIDs[card.ID] = struct{}{}
	}

	seen := make(map[int64]struct{})
	for _, orphanCard := range existingCardMap {
		if _, ok := seen[orphanCard.ID]; ok {
			continue
		}
		if _, ok := processedIDs[orphanCard.ID]; ok {
			continue
		}
		seen[orphanCard.ID] = struct{}{}
		if err := h.store.DeleteCard(orphanCard.ID); err != nil {
			return nil, err
		}
		h.removeCardFromDeck(col, orphanCard.DeckID, orphanCard.ID)
		delete(col.Cards, orphanCard.ID)
	}

	sort.Slice(updatedCards, func(i, j int) bool {
		if updatedCards[i].TemplateName == updatedCards[j].TemplateName {
			return updatedCards[i].Ordinal < updatedCards[j].Ordinal
		}
		return updatedCards[i].TemplateName < updatedCards[j].TemplateName
	})

	return updatedCards, nil
}

func (h *APIHandler) ListNotes(w http.ResponseWriter, r *http.Request) {
	col, collectionID, err := h.collectionForRequest(r)
	if err != nil {
		respondAPIError(w, http.StatusInternalServerError, "collection_load_failed", err.Error())
		return
	}

	limit := 25
	if rawLimit := strings.TrimSpace(r.URL.Query().Get("limit")); rawLimit != "" {
		parsed, err := strconv.Atoi(rawLimit)
		if err != nil || parsed <= 0 {
			respondAPIError(w, http.StatusBadRequest, "invalid_limit", "Limit must be a positive integer")
			return
		}
		if parsed > 100 {
			parsed = 100
		}
		limit = parsed
	}
	offset, err := parseCursorOffset(r.URL.Query().Get("cursor"))
	if err != nil {
		respondAPIError(w, http.StatusBadRequest, "invalid_cursor", "Cursor must be a non-negative integer")
		return
	}

	deckFilter := int64(0)
	if rawDeckID := strings.TrimSpace(r.URL.Query().Get("deckId")); rawDeckID != "" {
		deckFilter, err = strconv.ParseInt(rawDeckID, 10, 64)
		if err != nil || deckFilter <= 0 {
			respondAPIError(w, http.StatusBadRequest, "invalid_deck_id", "Invalid deck ID")
			return
		}
	}
	typeFilter := strings.TrimSpace(r.URL.Query().Get("typeId"))
	query := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("q")))
	tagFilter := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("tag")))

	notes, err := h.store.ListNotes(collectionID)
	if err != nil {
		respondAPIError(w, http.StatusInternalServerError, "notes_list_failed", err.Error())
		return
	}

	items := make([]NoteListItemResponse, 0, len(notes))
	for _, note := range notes {
		cards, err := h.store.GetCardsByNote(note.ID)
		if err != nil {
			respondAPIError(w, http.StatusInternalServerError, "note_cards_failed", err.Error())
			return
		}
		if deckFilter > 0 {
			inDeck := false
			for _, card := range cards {
				if card.DeckID == deckFilter {
					inDeck = true
					break
				}
			}
			if !inDeck {
				continue
			}
		}

		preview := h.noteFieldPreview(note, col)
		if !noteMatchesFilter(note, preview, query, tagFilter, typeFilter) {
			continue
		}

		primaryDeckID, primaryDeckName := h.primaryDeckDetails(cards, col)
		items = append(items, NoteListItemResponse{
			ID:           note.ID,
			TypeID:       string(note.Type),
			FieldVals:    note.FieldMap,
			FieldPreview: preview,
			Tags:         note.Tags,
			CreatedAt:    note.CreatedAt,
			ModifiedAt:   note.ModifiedAt,
			DeckID:       primaryDeckID,
			DeckName:     primaryDeckName,
			CardCount:    len(cards),
		})
	}

	sort.Slice(items, func(i, j int) bool {
		if items[i].ModifiedAt.Equal(items[j].ModifiedAt) {
			return items[i].ID > items[j].ID
		}
		return items[i].ModifiedAt.After(items[j].ModifiedAt)
	})

	total := len(items)
	if offset > total {
		offset = total
	}
	end := offset + limit
	if end > total {
		end = total
	}

	response := ListNotesResponse{
		Notes: items[offset:end],
		Total: total,
	}
	if end < total {
		response.NextCursor = strconv.Itoa(end)
	}
	if offset > 0 {
		prev := offset - limit
		if prev < 0 {
			prev = 0
		}
		response.PrevCursor = strconv.Itoa(prev)
	}

	respondJSON(w, http.StatusOK, response)
}

func (h *APIHandler) UpdateNote(w http.ResponseWriter, r *http.Request) {
	if !h.requireWorkspaceWritePermission(w, r) {
		return
	}
	col, _, err := h.collectionForRequest(r)
	if err != nil {
		respondAPIError(w, http.StatusInternalServerError, "collection_load_failed", err.Error())
		return
	}

	id, err := parseIDParam(r, "id")
	if err != nil {
		respondAPIError(w, http.StatusBadRequest, "invalid_note_id", "Invalid note ID")
		return
	}

	var req UpdateNoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondAPIError(w, http.StatusBadRequest, "invalid_request", "Invalid request body")
		return
	}
	if req.TypeID == "" || req.DeckID == 0 {
		respondAPIError(w, http.StatusBadRequest, "invalid_note_request", "TypeID and DeckID are required")
		return
	}

	note, err := h.store.GetNote(id)
	if err != nil {
		respondAPIError(w, http.StatusNotFound, "note_not_found", "Note not found")
		return
	}

	existingCards, err := h.store.GetCardsByNote(id)
	if err != nil {
		respondAPIError(w, http.StatusInternalServerError, "note_cards_failed", err.Error())
		return
	}

	note.Type = NoteTypeName(req.TypeID)
	note.FieldMap = sanitizeFieldVals(req.FieldVals)
	note.Tags = sanitizeTags(req.Tags)
	col.USN++
	note.USN = col.USN
	note.ModifiedAt = time.Now()

	noteType, ok := col.NoteTypes[note.Type]
	if !ok {
		respondAPIError(w, http.StatusBadRequest, "invalid_note_type", "Note type not found")
		return
	}
	previewCards, err := col.generateCardsFromNote(noteType, *note, req.DeckID, note.ModifiedAt)
	if err != nil {
		respondAPIError(w, http.StatusBadRequest, "note_update_failed", err.Error())
		return
	}
	session := h.sessionFromRequest(r)
	plan := h.planForRequest(r, session)
	usage := h.usageForSession(session)
	if err := validateCardsTotalLimit(plan, usage, len(previewCards)-len(existingCards)); err != nil {
		respondAPIError(w, http.StatusForbidden, "plan_limit_exceeded", err.Error())
		return
	}

	if err := h.store.UpdateNote(note); err != nil {
		respondAPIError(w, http.StatusInternalServerError, "note_update_failed", err.Error())
		return
	}
	updatedCards, err := h.regenerateCardsForSingleNote(col, note, req.DeckID, nil)
	if err != nil {
		respondAPIError(w, http.StatusInternalServerError, "card_regeneration_failed", err.Error())
		return
	}
	h.syncCollectionNote(col, note)
	h.markStudyGroupInstallsForkedByDeckIDs(req.DeckID)

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"note":  h.noteToResponse(note, updatedCards),
		"cards": updatedCards,
	})
}

func (h *APIHandler) DeleteNote(w http.ResponseWriter, r *http.Request) {
	if !h.requireWorkspaceWritePermission(w, r) {
		return
	}
	col, _, err := h.collectionForRequest(r)
	if err != nil {
		respondAPIError(w, http.StatusInternalServerError, "collection_load_failed", err.Error())
		return
	}

	id, err := parseIDParam(r, "id")
	if err != nil {
		respondAPIError(w, http.StatusBadRequest, "invalid_note_id", "Invalid note ID")
		return
	}

	if _, err := h.store.GetNote(id); err != nil {
		respondAPIError(w, http.StatusNotFound, "note_not_found", "Note not found")
		return
	}

	cards, err := h.store.GetCardsByNote(id)
	if err != nil {
		respondAPIError(w, http.StatusInternalServerError, "note_cards_failed", err.Error())
		return
	}
	deckIDs := make([]int64, 0, len(cards))
	for _, card := range cards {
		deckIDs = append(deckIDs, card.DeckID)
		if err := h.store.DeleteCard(card.ID); err != nil {
			respondAPIError(w, http.StatusInternalServerError, "card_delete_failed", err.Error())
			return
		}
		h.removeCardFromDeck(col, card.DeckID, card.ID)
		delete(col.Cards, card.ID)
	}
	if err := h.store.DeleteNote(id); err != nil {
		respondAPIError(w, http.StatusInternalServerError, "note_delete_failed", err.Error())
		return
	}
	delete(col.Notes, id)
	h.markStudyGroupInstallsForkedByDeckIDs(deckIDs...)
	w.WriteHeader(http.StatusNoContent)
}

func (h *APIHandler) UpdateDeck(w http.ResponseWriter, r *http.Request) {
	if !h.requireWorkspaceWritePermission(w, r) {
		return
	}
	col, _, err := h.collectionForRequest(r)
	if err != nil {
		respondAPIError(w, http.StatusInternalServerError, "collection_load_failed", err.Error())
		return
	}

	id, err := parseIDParam(r, "id")
	if err != nil {
		respondAPIError(w, http.StatusBadRequest, "invalid_deck_id", "Invalid deck ID")
		return
	}

	var req UpdateDeckRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondAPIError(w, http.StatusBadRequest, "invalid_request", "Invalid request body")
		return
	}
	if req.Name == nil && req.NewCardsPerDay == nil && req.ReviewsPerDay == nil && req.PriorityOrder == nil {
		respondAPIError(w, http.StatusBadRequest, "invalid_request", "At least one deck field is required")
		return
	}

	deck, err := h.store.GetDeck(id)
	if err != nil {
		respondAPIError(w, http.StatusNotFound, "deck_not_found", "Deck not found")
		return
	}

	if req.Name != nil {
		trimmed := strings.TrimSpace(*req.Name)
		if trimmed == "" {
			respondAPIError(w, http.StatusBadRequest, "invalid_name", "Deck name is required")
			return
		}
		deck.Name = sanitizeHTML(trimmed)
	}
	if req.PriorityOrder != nil {
		if *req.PriorityOrder <= 0 {
			respondAPIError(w, http.StatusBadRequest, "invalid_priority_order", "Priority must be 1 or greater")
			return
		}
		deck.PriorityOrder = *req.PriorityOrder
	}
	if req.NewCardsPerDay != nil || req.ReviewsPerDay != nil {
		if req.NewCardsPerDay != nil && *req.NewCardsPerDay < 0 {
			respondAPIError(w, http.StatusBadRequest, "invalid_new_cards_per_day", "New cards per day must be 0 or greater")
			return
		}
		if req.ReviewsPerDay != nil && *req.ReviewsPerDay < 0 {
			respondAPIError(w, http.StatusBadRequest, "invalid_reviews_per_day", "Reviews per day must be 0 or greater")
			return
		}

		options, err := h.store.EnsureDeckOptionsForDeck(deck)
		if err != nil {
			respondAPIError(w, http.StatusInternalServerError, "deck_options_failed", err.Error())
			return
		}
		if req.NewCardsPerDay != nil {
			options.NewCardsPerDay = *req.NewCardsPerDay
		}
		if req.ReviewsPerDay != nil {
			options.ReviewsPerDay = *req.ReviewsPerDay
		}
		options.Name = fmt.Sprintf("%s settings", deck.Name)
		if err := h.store.UpdateDeckOptions(options); err != nil {
			respondAPIError(w, http.StatusInternalServerError, "deck_options_failed", err.Error())
			return
		}
	}
	if err := h.store.UpdateDeck(deck); err != nil {
		respondAPIError(w, http.StatusInternalServerError, "deck_update_failed", err.Error())
		return
	}
	if existing, ok := col.Decks[id]; ok {
		existing.Name = deck.Name
		existing.OptionsID = deck.OptionsID
		existing.PriorityOrder = deck.PriorityOrder
	}
	h.markStudyGroupInstallsForkedByDeckIDs(id)

	respondJSON(w, http.StatusOK, h.deckResponse(h.userIDFromRequest(r), deck, col, nil))
}

func (h *APIHandler) DeleteDeck(w http.ResponseWriter, r *http.Request) {
	if !h.requireWorkspaceWritePermission(w, r) {
		return
	}
	col, _, err := h.collectionForRequest(r)
	if err != nil {
		respondAPIError(w, http.StatusInternalServerError, "collection_load_failed", err.Error())
		return
	}

	id, err := parseIDParam(r, "id")
	if err != nil {
		respondAPIError(w, http.StatusBadRequest, "invalid_deck_id", "Invalid deck ID")
		return
	}

	deck, err := h.store.GetDeck(id)
	if err != nil {
		respondAPIError(w, http.StatusNotFound, "deck_not_found", "Deck not found")
		return
	}
	if len(deck.Cards) > 0 {
		respondAPIError(w, http.StatusConflict, "deck_not_empty", "Only empty decks can be deleted right now. Move or delete the cards in this deck first.")
		return
	}
	for _, candidate := range col.Decks {
		if candidate.ParentID != nil && *candidate.ParentID == id {
			respondAPIError(w, http.StatusConflict, "deck_has_children", "This deck has child decks. Move or delete those child decks first.")
			return
		}
	}
	if err := h.store.DeleteDeck(id); err != nil {
		respondAPIError(w, http.StatusInternalServerError, "deck_delete_failed", err.Error())
		return
	}
	delete(col.Decks, id)
	w.WriteHeader(http.StatusNoContent)
}

func uniqueTemplateName(existing []CardTemplate, desired string) string {
	base := strings.TrimSpace(desired)
	if base == "" {
		base = "Card"
	}
	name := base
	index := 2
	for {
		duplicate := false
		for _, template := range existing {
			if strings.EqualFold(template.Name, name) {
				duplicate = true
				break
			}
		}
		if !duplicate {
			return name
		}
		name = fmt.Sprintf("%s %d", base, index)
		index++
	}
}

func defaultTemplateForNoteType(nt NoteType, name string) CardTemplate {
	template := CardTemplate{
		Name:    name,
		Styling: "",
	}
	if len(nt.Templates) > 0 {
		template.IsCloze = nt.Templates[0].IsCloze
		template.Styling = nt.Templates[0].Styling
	}
	if template.IsCloze {
		fieldName := "Text"
		if len(nt.Fields) > 0 {
			fieldName = nt.Fields[0]
		}
		template.QFmt = fmt.Sprintf("{{cloze:%s}}", fieldName)
		template.AFmt = "{{FrontSide}}"
		return template
	}
	frontField := "Front"
	backField := "Back"
	if len(nt.Fields) > 0 {
		frontField = nt.Fields[0]
	}
	if len(nt.Fields) > 1 {
		backField = nt.Fields[1]
	} else {
		backField = frontField
	}
	template.QFmt = fmt.Sprintf("{{%s}}", frontField)
	template.AFmt = fmt.Sprintf("{{FrontSide}}\n\n<hr id=\"answer\">\n\n{{%s}}", backField)
	return template
}

func (h *APIHandler) CreateTemplate(w http.ResponseWriter, r *http.Request) {
	if !h.requireWorkspaceWritePermission(w, r) {
		return
	}
	col, collectionID, err := h.collectionForRequest(r)
	if err != nil {
		respondAPIError(w, http.StatusInternalServerError, "collection_load_failed", err.Error())
		return
	}

	noteTypeName := chi.URLParam(r, "name")
	nt, ok := col.NoteTypes[NoteTypeName(noteTypeName)]
	if !ok {
		respondAPIError(w, http.StatusNotFound, "note_type_not_found", "Note type not found")
		return
	}

	var req CreateTemplateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondAPIError(w, http.StatusBadRequest, "invalid_request", "Invalid request body")
		return
	}

	templateName := uniqueTemplateName(nt.Templates, sanitizeHTML(req.Name))
	template := defaultTemplateForNoteType(nt, templateName)
	if req.SourceTemplateName != "" {
		found := false
		for _, candidate := range nt.Templates {
			if candidate.Name == req.SourceTemplateName {
				template = candidate
				template.Name = templateName
				found = true
				break
			}
		}
		if !found {
			respondAPIError(w, http.StatusNotFound, "template_not_found", "Source template not found")
			return
		}
	}

	nt.Templates = append(nt.Templates, template)
	if err := h.store.UpdateNoteType(collectionID, &nt); err != nil {
		respondAPIError(w, http.StatusInternalServerError, "template_create_failed", err.Error())
		return
	}
	col.NoteTypes[NoteTypeName(noteTypeName)] = nt
	h.markStudyGroupInstallsForkedByNoteType(noteTypeName)
	if err := h.regenerateCardsForNoteTypeInCollection(collectionID, col, noteTypeName); err != nil {
		respondAPIError(w, http.StatusInternalServerError, "card_regeneration_failed", err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, buildTemplatesResponse(nt, "Template created successfully"))
}

func (h *APIHandler) DeleteTemplate(w http.ResponseWriter, r *http.Request) {
	if !h.requireWorkspaceWritePermission(w, r) {
		return
	}
	col, collectionID, err := h.collectionForRequest(r)
	if err != nil {
		respondAPIError(w, http.StatusInternalServerError, "collection_load_failed", err.Error())
		return
	}

	noteTypeName := chi.URLParam(r, "name")
	templateName := chi.URLParam(r, "templateName")
	nt, ok := col.NoteTypes[NoteTypeName(noteTypeName)]
	if !ok {
		respondAPIError(w, http.StatusNotFound, "note_type_not_found", "Note type not found")
		return
	}
	if len(nt.Templates) <= 1 {
		respondAPIError(w, http.StatusBadRequest, "template_last_remaining", "A note type must keep at least one template")
		return
	}

	filtered := make([]CardTemplate, 0, len(nt.Templates)-1)
	found := false
	for _, template := range nt.Templates {
		if template.Name == templateName {
			found = true
			continue
		}
		filtered = append(filtered, template)
	}
	if !found {
		respondAPIError(w, http.StatusNotFound, "template_not_found", "Template not found")
		return
	}
	nt.Templates = filtered
	if err := h.store.UpdateNoteType(collectionID, &nt); err != nil {
		respondAPIError(w, http.StatusInternalServerError, "template_delete_failed", err.Error())
		return
	}
	col.NoteTypes[NoteTypeName(noteTypeName)] = nt
	h.markStudyGroupInstallsForkedByNoteType(noteTypeName)
	if err := h.regenerateCardsForNoteTypeInCollection(collectionID, col, noteTypeName); err != nil {
		respondAPIError(w, http.StatusInternalServerError, "card_regeneration_failed", err.Error())
		return
	}
	respondJSON(w, http.StatusOK, buildTemplatesResponse(nt, "Template deleted successfully"))
}

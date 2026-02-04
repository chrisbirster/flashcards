package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/microcosm-cc/bluemonday"
	fsrs "github.com/open-spaced-repetition/go-fsrs/v3"
)

var htmlPolicy = bluemonday.UGCPolicy()

func sanitizeHTML(input string) string {
	return htmlPolicy.Sanitize(input)
}

// APIHandler wraps the store and provides HTTP handlers
type APIHandler struct {
	store         *SQLiteStore
	collectionID  string
	collection    *Collection
	backupManager *BackupManager
}

func NewAPIHandler(store *SQLiteStore, collection *Collection, backupMgr *BackupManager) *APIHandler {
	return &APIHandler{
		store:         store,
		collectionID:  "default",
		collection:    collection,
		backupManager: backupMgr,
	}
}

// Request/Response types
type CreateDeckRequest struct {
	Name string `json:"name"`
}

type DeckResponse struct {
	ID       int64   `json:"id"`
	Name     string  `json:"name"`
	ParentID *int64  `json:"parentId,omitempty"`
	CardIDs  []int64 `json:"cardIds"`
}

type CreateNoteRequest struct {
	TypeID         string            `json:"typeId"`
	DeckID         int64             `json:"deckId"`
	FieldVals      map[string]string `json:"fieldVals"`
	Tags           []string          `json:"tags"`
	AllowDuplicate bool              `json:"allowDuplicate"` // Override duplicate check
}

type CheckDuplicateRequest struct {
	TypeID    string `json:"typeId"`
	FieldName string `json:"fieldName"` // Field to check for duplicates (usually "Front" or first field)
	Value     string `json:"value"`
	DeckID    int64  `json:"deckId,omitempty"` // Optional: limit scope to deck
}

type DuplicateResult struct {
	IsDuplicate bool        `json:"isDuplicate"`
	Duplicates  []NoteBrief `json:"duplicates,omitempty"`
}

type NoteBrief struct {
	ID       int64             `json:"id"`
	TypeID   string            `json:"typeId"`
	FieldVal map[string]string `json:"fieldVals"`
	DeckID   int64             `json:"deckId,omitempty"`
}

type AnswerCardRequest struct {
	Rating      int `json:"rating"`      // 1=Again, 2=Hard, 3=Good, 4=Easy
	TimeTakenMs int `json:"timeTakenMs"` // Time spent on the card in milliseconds
}

type UpdateCardRequest struct {
	Flag      *int  `json:"flag,omitempty"`      // 0-7 color flags
	Marked    *bool `json:"marked,omitempty"`    // toggle marked status
	Suspended *bool `json:"suspended,omitempty"` // toggle suspended status
}

// Handler methods

func (h *APIHandler) GetCollection(w http.ResponseWriter, r *http.Request) {
	col, err := h.store.GetCollection(h.collectionID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusOK, col)
}

func (h *APIHandler) ListDecks(w http.ResponseWriter, r *http.Request) {
	decks, err := h.store.ListDecks(h.collectionID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert to response format
	var response []DeckResponse
	for _, d := range decks {
		response = append(response, DeckResponse{
			ID:       d.ID,
			Name:     d.Name,
			ParentID: d.ParentID,
			CardIDs:  d.Cards,
		})
	}

	respondJSON(w, http.StatusOK, response)
}

func (h *APIHandler) CreateDeck(w http.ResponseWriter, r *http.Request) {
	var req CreateDeckRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, "Deck name is required", http.StatusBadRequest)
		return
	}

	// Sanitize deck name to prevent XSS
	sanitizedName := sanitizeHTML(req.Name)

	// Create deck using collection method
	deck := h.collection.NewDeck(sanitizedName)

	// Persist to database
	if err := h.store.CreateDeck(deck); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusCreated, DeckResponse{
		ID:       deck.ID,
		Name:     deck.Name,
		ParentID: deck.ParentID,
		CardIDs:  deck.Cards,
	})
}

func (h *APIHandler) GetDeck(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDParam(r, "id")
	if err != nil {
		http.Error(w, "Invalid deck ID", http.StatusBadRequest)
		return
	}

	deck, err := h.store.GetDeck(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Get deck stats
	stats, err := h.store.GetDeckStats(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"deck": DeckResponse{
			ID:       deck.ID,
			Name:     deck.Name,
			ParentID: deck.ParentID,
			CardIDs:  deck.Cards,
		},
		"stats": stats,
	})
}

func (h *APIHandler) GetDeckStats(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDParam(r, "id")
	if err != nil {
		http.Error(w, "Invalid deck ID", http.StatusBadRequest)
		return
	}

	stats, err := h.store.GetDeckStats(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusOK, stats)
}

func (h *APIHandler) CreateNote(w http.ResponseWriter, r *http.Request) {
	var req CreateNoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.TypeID == "" || req.DeckID == 0 {
		http.Error(w, "TypeID and DeckID are required", http.StatusBadRequest)
		return
	}

	// Sanitize field values to prevent XSS
	sanitizedFieldVals := make(map[string]string)
	for field, value := range req.FieldVals {
		sanitizedFieldVals[field] = sanitizeHTML(value)
	}

	// Sanitize tags (strip HTML since tags should be plain text)
	sanitizedTags := make([]string, len(req.Tags))
	for i, tag := range req.Tags {
		sanitizedTags[i] = sanitizeHTML(tag)
	}

	// Use Collection.AddNote to create note and generate cards
	note, cards, err := h.collection.AddNote(req.DeckID, NoteTypeName(req.TypeID), sanitizedFieldVals, time.Now())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Set tags if provided (use sanitized tags)
	note.Tags = sanitizedTags
	if sanitizedTags == nil {
		note.Tags = []string{}
	}

	// Persist note to database
	if err := h.store.CreateNote(h.collectionID, &note); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Persist generated cards to database
	for _, card := range cards {
		if err := h.store.CreateCard(card); err != nil {
			http.Error(w, fmt.Sprintf("Failed to save card: %v", err), http.StatusInternalServerError)
			return
		}
	}

	respondJSON(w, http.StatusCreated, map[string]interface{}{
		"note":  note,
		"cards": cards,
	})
}

func (h *APIHandler) GetNote(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDParam(r, "id")
	if err != nil {
		http.Error(w, "Invalid note ID", http.StatusBadRequest)
		return
	}

	note, err := h.store.GetNote(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	respondJSON(w, http.StatusOK, note)
}

func (h *APIHandler) CheckDuplicate(w http.ResponseWriter, r *http.Request) {
	var req CheckDuplicateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Value == "" {
		respondJSON(w, http.StatusOK, DuplicateResult{IsDuplicate: false})
		return
	}

	// Check for duplicates in the collection
	duplicates, err := h.store.FindDuplicateNotes(h.collectionID, req.FieldName, req.Value, req.DeckID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	result := DuplicateResult{
		IsDuplicate: len(duplicates) > 0,
		Duplicates:  duplicates,
	}

	respondJSON(w, http.StatusOK, result)
}

func (h *APIHandler) GetDueCards(w http.ResponseWriter, r *http.Request) {
	deckID, err := parseIDParam(r, "deckId")
	if err != nil {
		http.Error(w, "Invalid deck ID", http.StatusBadRequest)
		return
	}

	limit := 10
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	cards, err := h.store.GetDueCards(deckID, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusOK, cards)
}

func (h *APIHandler) GetCard(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDParam(r, "id")
	if err != nil {
		http.Error(w, "Invalid card ID", http.StatusBadRequest)
		return
	}

	card, err := h.store.GetCard(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	respondJSON(w, http.StatusOK, card)
}

func (h *APIHandler) AnswerCard(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDParam(r, "id")
	if err != nil {
		http.Error(w, "Invalid card ID", http.StatusBadRequest)
		return
	}

	var req AnswerCardRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Rating < 1 || req.Rating > 4 {
		http.Error(w, "Rating must be 1-4 (Again/Hard/Good/Easy)", http.StatusBadRequest)
		return
	}

	// Use Collection.Answer to update FSRS scheduling
	revlog, err := h.collection.Answer(id, fsrs.Rating(req.Rating), time.Now(), req.TimeTakenMs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get updated card from collection
	card, ok := h.collection.Cards[id]
	if !ok {
		http.Error(w, "Card not found after update", http.StatusInternalServerError)
		return
	}

	// Persist updated card to database
	if err := h.store.UpdateCard(card); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Persist revlog entry with actual card ID and time taken
	if err := h.store.AddRevlog(revlog, id, req.TimeTakenMs); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusOK, card)
}

func (h *APIHandler) UpdateCard(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDParam(r, "id")
	if err != nil {
		http.Error(w, "Invalid card ID", http.StatusBadRequest)
		return
	}

	var req UpdateCardRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Get card from store
	card, err := h.store.GetCard(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Update fields if provided
	if req.Flag != nil {
		if *req.Flag < 0 || *req.Flag > 7 {
			http.Error(w, "Flag must be 0-7", http.StatusBadRequest)
			return
		}
		card.Flag = *req.Flag
	}
	if req.Marked != nil {
		card.Marked = *req.Marked
	}
	if req.Suspended != nil {
		card.Suspended = *req.Suspended
	}

	// Persist changes
	if err := h.store.UpdateCard(card); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Update collection cache if card exists there
	if _, ok := h.collection.Cards[id]; ok {
		h.collection.Cards[id] = card
	}

	respondJSON(w, http.StatusOK, card)
}

func (h *APIHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, map[string]string{
		"status":  "ok",
		"service": "microdote-api",
		"version": "M2",
	})
}

// Note type response for API
type NoteTypeResponse struct {
	Name           string                  `json:"name"`
	Fields         []string                `json:"fields"`
	Templates      []TemplateInfo          `json:"templates"`
	SortFieldIndex int                     `json:"sortFieldIndex"`
	FieldOptions   map[string]FieldOptions `json:"fieldOptions,omitempty"`
}

type TemplateInfo struct {
	Name            string `json:"name"`
	QFmt            string `json:"qFmt"`
	AFmt            string `json:"aFmt"`
	Styling         string `json:"styling"`
	IfFieldNonEmpty string `json:"ifFieldNonEmpty,omitempty"`
	IsCloze         bool   `json:"isCloze"`
	DeckOverride    string `json:"deckOverride,omitempty"`
	BrowserQFmt     string `json:"browserQFmt,omitempty"`
	BrowserAFmt     string `json:"browserAFmt,omitempty"`
}

type SetSortFieldRequest struct {
	FieldIndex int `json:"fieldIndex"` // Index of the field to use as sort field
}

type SetFieldOptionsRequest struct {
	FieldName string       `json:"fieldName"`
	Options   FieldOptions `json:"options"`
}

type UpdateTemplateRequest struct {
	QFmt            *string `json:"qFmt,omitempty"`
	AFmt            *string `json:"aFmt,omitempty"`
	Styling         *string `json:"styling,omitempty"`
	IfFieldNonEmpty *string `json:"ifFieldNonEmpty,omitempty"`
	DeckOverride    *string `json:"deckOverride,omitempty"`
	BrowserQFmt     *string `json:"browserQFmt,omitempty"`
	BrowserAFmt     *string `json:"browserAFmt,omitempty"`
}

type TemplatesResponse struct {
	Message   string         `json:"message"`
	Templates []TemplateInfo `json:"templates"`
}

func (h *APIHandler) ListNoteTypes(w http.ResponseWriter, r *http.Request) {
	var noteTypes []NoteTypeResponse
	for _, nt := range h.collection.NoteTypes {
		var templates []TemplateInfo
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
		noteTypes = append(noteTypes, NoteTypeResponse{
			Name:           string(nt.Name),
			Fields:         nt.Fields,
			Templates:      templates,
			SortFieldIndex: nt.SortFieldIndex,
			FieldOptions:   nt.FieldOptions,
		})
	}
	respondJSON(w, http.StatusOK, noteTypes)
}

func (h *APIHandler) GetNoteType(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	nt, ok := h.collection.NoteTypes[NoteTypeName(name)]
	if !ok {
		http.Error(w, "Note type not found", http.StatusNotFound)
		return
	}

	var templates []TemplateInfo
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

	respondJSON(w, http.StatusOK, NoteTypeResponse{
		Name:           string(nt.Name),
		Fields:         nt.Fields,
		Templates:      templates,
		SortFieldIndex: nt.SortFieldIndex,
		FieldOptions:   nt.FieldOptions,
	})
}

// Reserved field names that cannot be used
var reservedFieldNames = map[string]bool{
	"Tags":      true,
	"Type":      true,
	"Deck":      true,
	"Card":      true,
	"FrontSide": true,
}

// Field management request types
type AddFieldRequest struct {
	FieldName string `json:"fieldName"`
	Position  *int   `json:"position,omitempty"` // Optional: insert at specific position
}

type RenameFieldRequest struct {
	OldName string `json:"oldName"`
	NewName string `json:"newName"`
}

type RemoveFieldRequest struct {
	FieldName string `json:"fieldName"`
}

type ReorderFieldsRequest struct {
	Fields []string `json:"fields"` // New field order
}

func (h *APIHandler) AddField(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	nt, ok := h.collection.NoteTypes[NoteTypeName(name)]
	if !ok {
		http.Error(w, "Note type not found", http.StatusNotFound)
		return
	}

	var req AddFieldRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.FieldName == "" {
		http.Error(w, "fieldName is required", http.StatusBadRequest)
		return
	}

	// Sanitize field name to prevent XSS
	sanitizedFieldName := sanitizeHTML(req.FieldName)

	// Check for reserved field names
	if reservedFieldNames[sanitizedFieldName] {
		http.Error(w, fmt.Sprintf("'%s' is a reserved field name", sanitizedFieldName), http.StatusBadRequest)
		return
	}

	// Check for duplicate field name
	for _, f := range nt.Fields {
		if f == sanitizedFieldName {
			http.Error(w, "Field name already exists", http.StatusBadRequest)
			return
		}
	}

	// Add field at specified position or end
	if req.Position != nil && *req.Position >= 0 && *req.Position < len(nt.Fields) {
		newFields := make([]string, 0, len(nt.Fields)+1)
		newFields = append(newFields, nt.Fields[:*req.Position]...)
		newFields = append(newFields, sanitizedFieldName)
		newFields = append(newFields, nt.Fields[*req.Position:]...)
		nt.Fields = newFields
	} else {
		nt.Fields = append(nt.Fields, sanitizedFieldName)
	}

	// Update in store
	if err := h.store.UpdateNoteType(h.collectionID, &nt); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Update collection cache
	h.collection.NoteTypes[NoteTypeName(name)] = nt

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Field added successfully",
		"fields":  nt.Fields,
	})
}

func (h *APIHandler) RenameField(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	nt, ok := h.collection.NoteTypes[NoteTypeName(name)]
	if !ok {
		http.Error(w, "Note type not found", http.StatusNotFound)
		return
	}

	var req RenameFieldRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.OldName == "" || req.NewName == "" {
		http.Error(w, "oldName and newName are required", http.StatusBadRequest)
		return
	}

	// Sanitize new field name to prevent XSS
	sanitizedNewName := sanitizeHTML(req.NewName)

	// Check for reserved field names
	if reservedFieldNames[sanitizedNewName] {
		http.Error(w, fmt.Sprintf("'%s' is a reserved field name", sanitizedNewName), http.StatusBadRequest)
		return
	}

	// Find and rename the field
	found := false
	for i, f := range nt.Fields {
		if f == req.OldName {
			nt.Fields[i] = sanitizedNewName
			found = true
			break
		}
	}

	if !found {
		http.Error(w, "Field not found", http.StatusNotFound)
		return
	}

	// Check for duplicate with new name
	count := 0
	for _, f := range nt.Fields {
		if f == sanitizedNewName {
			count++
		}
	}
	if count > 1 {
		// Revert and return error
		for i, f := range nt.Fields {
			if f == sanitizedNewName {
				nt.Fields[i] = req.OldName
				break
			}
		}
		http.Error(w, "Field name already exists", http.StatusBadRequest)
		return
	}

	// Update templates to use new field name
	for i := range nt.Templates {
		nt.Templates[i].QFmt = strings.ReplaceAll(nt.Templates[i].QFmt, "{{"+req.OldName+"}}", "{{"+sanitizedNewName+"}}")
		nt.Templates[i].AFmt = strings.ReplaceAll(nt.Templates[i].AFmt, "{{"+req.OldName+"}}", "{{"+sanitizedNewName+"}}")
		if nt.Templates[i].IfFieldNonEmpty == req.OldName {
			nt.Templates[i].IfFieldNonEmpty = sanitizedNewName
		}
	}

	// Update in store
	if err := h.store.UpdateNoteType(h.collectionID, &nt); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Update collection cache
	h.collection.NoteTypes[NoteTypeName(name)] = nt

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Field renamed successfully",
		"fields":  nt.Fields,
	})
}

func (h *APIHandler) RemoveField(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	nt, ok := h.collection.NoteTypes[NoteTypeName(name)]
	if !ok {
		http.Error(w, "Note type not found", http.StatusNotFound)
		return
	}

	var req RemoveFieldRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.FieldName == "" {
		http.Error(w, "fieldName is required", http.StatusBadRequest)
		return
	}

	// Must have at least one field
	if len(nt.Fields) <= 1 {
		http.Error(w, "Cannot remove the last field", http.StatusBadRequest)
		return
	}

	// Find and remove the field
	found := false
	newFields := make([]string, 0, len(nt.Fields)-1)
	for _, f := range nt.Fields {
		if f == req.FieldName {
			found = true
		} else {
			newFields = append(newFields, f)
		}
	}

	if !found {
		http.Error(w, "Field not found", http.StatusNotFound)
		return
	}

	nt.Fields = newFields

	// Update in store
	if err := h.store.UpdateNoteType(h.collectionID, &nt); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Update collection cache
	h.collection.NoteTypes[NoteTypeName(name)] = nt

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Field removed successfully",
		"fields":  nt.Fields,
	})
}

func (h *APIHandler) ReorderFields(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	nt, ok := h.collection.NoteTypes[NoteTypeName(name)]
	if !ok {
		http.Error(w, "Note type not found", http.StatusNotFound)
		return
	}

	var req ReorderFieldsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate that the new order contains the same fields
	if len(req.Fields) != len(nt.Fields) {
		http.Error(w, "Field count mismatch", http.StatusBadRequest)
		return
	}

	// Check that all existing fields are present
	existingFields := make(map[string]bool)
	for _, f := range nt.Fields {
		existingFields[f] = true
	}

	for _, f := range req.Fields {
		if !existingFields[f] {
			http.Error(w, fmt.Sprintf("Unknown field: %s", f), http.StatusBadRequest)
			return
		}
		delete(existingFields, f)
	}

	if len(existingFields) > 0 {
		http.Error(w, "Some fields are missing from the new order", http.StatusBadRequest)
		return
	}

	nt.Fields = req.Fields

	// Update in store
	if err := h.store.UpdateNoteType(h.collectionID, &nt); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Update collection cache
	h.collection.NoteTypes[NoteTypeName(name)] = nt

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Fields reordered successfully",
		"fields":  nt.Fields,
	})
}

func (h *APIHandler) SetSortField(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	nt, ok := h.collection.NoteTypes[NoteTypeName(name)]
	if !ok {
		http.Error(w, "Note type not found", http.StatusNotFound)
		return
	}

	var req SetSortFieldRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate field index
	if req.FieldIndex < 0 || req.FieldIndex >= len(nt.Fields) {
		http.Error(w, "Invalid field index", http.StatusBadRequest)
		return
	}

	nt.SortFieldIndex = req.FieldIndex

	// Update in store
	if err := h.store.UpdateNoteType(h.collectionID, &nt); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Update collection cache
	h.collection.NoteTypes[NoteTypeName(name)] = nt

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"message":        "Sort field updated successfully",
		"sortFieldIndex": nt.SortFieldIndex,
		"sortFieldName":  nt.Fields[nt.SortFieldIndex],
	})
}

func (h *APIHandler) SetFieldOptions(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	nt, ok := h.collection.NoteTypes[NoteTypeName(name)]
	if !ok {
		http.Error(w, "Note type not found", http.StatusNotFound)
		return
	}

	var req SetFieldOptionsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate field exists
	fieldExists := false
	for _, f := range nt.Fields {
		if f == req.FieldName {
			fieldExists = true
			break
		}
	}
	if !fieldExists {
		http.Error(w, "Field not found", http.StatusBadRequest)
		return
	}

	// Initialize FieldOptions map if nil
	if nt.FieldOptions == nil {
		nt.FieldOptions = make(map[string]FieldOptions)
	}

	// Update field options
	nt.FieldOptions[req.FieldName] = req.Options

	// Update in store
	if err := h.store.UpdateNoteType(h.collectionID, &nt); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Update collection cache
	h.collection.NoteTypes[NoteTypeName(name)] = nt

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"message":      "Field options updated successfully",
		"fieldOptions": nt.FieldOptions,
	})
}

func (h *APIHandler) UpdateTemplate(w http.ResponseWriter, r *http.Request) {
	noteTypeName := chi.URLParam(r, "name")
	templateName := chi.URLParam(r, "templateName")

	nt, ok := h.collection.NoteTypes[NoteTypeName(noteTypeName)]
	if !ok {
		http.Error(w, "Note type not found", http.StatusNotFound)
		return
	}

	var req UpdateTemplateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Find and update the template
	templateIndex := -1
	for i, t := range nt.Templates {
		if t.Name == templateName {
			templateIndex = i
			break
		}
	}

	if templateIndex == -1 {
		http.Error(w, "Template not found", http.StatusNotFound)
		return
	}

	// Update fields if provided
	if req.QFmt != nil {
		nt.Templates[templateIndex].QFmt = sanitizeHTML(*req.QFmt)
	}
	if req.AFmt != nil {
		nt.Templates[templateIndex].AFmt = sanitizeHTML(*req.AFmt)
	}
	if req.Styling != nil {
		nt.Templates[templateIndex].Styling = sanitizeHTML(*req.Styling)
	}
	if req.IfFieldNonEmpty != nil {
		nt.Templates[templateIndex].IfFieldNonEmpty = sanitizeHTML(*req.IfFieldNonEmpty)
	}
	if req.DeckOverride != nil {
		nt.Templates[templateIndex].DeckOverride = sanitizeHTML(*req.DeckOverride)
	}
	if req.BrowserQFmt != nil {
		nt.Templates[templateIndex].BrowserQFmt = sanitizeHTML(*req.BrowserQFmt)
	}
	if req.BrowserAFmt != nil {
		nt.Templates[templateIndex].BrowserAFmt = sanitizeHTML(*req.BrowserAFmt)
	}

	// Update in store
	if err := h.store.UpdateNoteType(h.collectionID, &nt); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Update collection cache
	h.collection.NoteTypes[NoteTypeName(noteTypeName)] = nt

	// Regenerate cards for all notes of this type
	// This ensures cards reflect the updated templates
	if err := h.regenerateCardsForNoteType(noteTypeName); err != nil {
		log.Printf("Warning: Failed to regenerate cards after template update: %v", err)
		// Don't fail the request - template was updated successfully
	}

	// Build response with updated templates
	var templates []TemplateInfo
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

	respondJSON(w, http.StatusOK, TemplatesResponse{
		Message:   "Template updated successfully",
		Templates: templates,
	})
}

// regenerateCardsForNoteType regenerates cards for all notes of a given note type.
// This preserves existing card scheduling data (SRS state, flags, etc.) while updating content.
func (h *APIHandler) regenerateCardsForNoteType(noteTypeName string) error {
	// Get all notes of this type
	notes, err := h.store.GetNotesByType(h.collectionID, noteTypeName)
	if err != nil {
		return fmt.Errorf("failed to get notes: %w", err)
	}

	for _, note := range notes {
		// Get the deck ID from one of the note's existing cards
		// If the note has no cards, we'll use the default deck
		deckID := int64(1) // default deck
		existingCards, err := h.store.GetCardsByNote(note.ID)
		if err == nil && len(existingCards) > 0 {
			deckID = existingCards[0].DeckID
		}

		// Generate new cards based on current templates
		newCards, err := h.collection.GenerateCards(&note, deckID, time.Now())
		if err != nil {
			log.Printf("Warning: Failed to regenerate cards for note %d: %v", note.ID, err)
			continue
		}

		// Build a map of existing cards by template name and ordinal
		existingCardMap := make(map[string]*Card)
		for _, card := range existingCards {
			key := fmt.Sprintf("%s:%d", card.TemplateName, card.Ordinal)
			existingCardMap[key] = &card
		}

		// Process each newly generated card
		for _, newCard := range newCards {
			key := fmt.Sprintf("%s:%d", newCard.TemplateName, newCard.Ordinal)

			if existingCard, exists := existingCardMap[key]; exists {
				// Card already exists - update content but preserve SRS state
				existingCard.Front = newCard.Front
				existingCard.Back = newCard.Back
				if err := h.store.UpdateCard(existingCard); err != nil {
					log.Printf("Warning: Failed to update card %d: %v", existingCard.ID, err)
				}
				// Remove from map so we know it was processed
				delete(existingCardMap, key)
			} else {
				// New card needs to be created
				if err := h.store.CreateCard(newCard); err != nil {
					log.Printf("Warning: Failed to create new card for note %d: %v", note.ID, err)
				}
			}
		}

		// Any remaining cards in the map no longer match templates and should be deleted
		// (e.g., ifFieldNonEmpty condition no longer met, or template removed)
		for _, orphanCard := range existingCardMap {
			if err := h.store.DeleteCard(orphanCard.ID); err != nil {
				log.Printf("Warning: Failed to delete orphaned card %d: %v", orphanCard.ID, err)
			}
		}
	}

	return nil
}

// Empty cards detection and cleanup

type EmptyCardInfo struct {
	CardID       int64  `json:"cardId"`
	NoteID       int64  `json:"noteId"`
	DeckID       int64  `json:"deckId"`
	TemplateName string `json:"templateName"`
	Ordinal      int    `json:"ordinal"`
	Front        string `json:"front"`
	Back         string `json:"back"`
	Reason       string `json:"reason"` // Why this card is considered empty
}

type EmptyCardsResponse struct {
	Count      int             `json:"count"`
	EmptyCards []EmptyCardInfo `json:"emptyCards"`
}

// FindEmptyCards detects cards that have no meaningful content
// This primarily targets cloze cards where the cloze deletion was removed
func (h *APIHandler) FindEmptyCards(w http.ResponseWriter, r *http.Request) {
	// Get all notes
	notes, err := h.store.ListNotes(h.collectionID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var emptyCards []EmptyCardInfo

	for _, note := range notes {
		// Get note type to check if it's cloze
		nt, ok := h.collection.NoteTypes[note.Type]
		if !ok {
			continue
		}

		// Get cards for this note
		cards, err := h.store.GetCardsByNote(note.ID)
		if err != nil {
			log.Printf("Warning: Failed to get cards for note %d: %v", note.ID, err)
			continue
		}

		for _, card := range cards {
			// Check if this is a cloze card
			isCloze := false
			for _, tmpl := range nt.Templates {
				if tmpl.Name == card.TemplateName && tmpl.IsCloze {
					isCloze = true
					break
				}
			}

			if isCloze {
				// Check if the cloze ordinal still exists in the text
				textField := note.FieldMap["Text"]
				ordinals := extractClozeOrdinals(textField)

				hasOrdinal := false
				for _, ord := range ordinals {
					if ord == card.Ordinal {
						hasOrdinal = true
						break
					}
				}

				if !hasOrdinal {
					emptyCards = append(emptyCards, EmptyCardInfo{
						CardID:       card.ID,
						NoteID:       note.ID,
						DeckID:       card.DeckID,
						TemplateName: card.TemplateName,
						Ordinal:      card.Ordinal,
						Front:        card.Front,
						Back:         card.Back,
						Reason:       fmt.Sprintf("Cloze deletion c%d no longer exists in note", card.Ordinal),
					})
				}
			} else {
				// For non-cloze cards, check if front and back are essentially empty
				// Strip HTML tags and check if there's meaningful content
				frontStripped := stripHTML(card.Front)
				backStripped := stripHTML(card.Back)

				if strings.TrimSpace(frontStripped) == "" && strings.TrimSpace(backStripped) == "" {
					emptyCards = append(emptyCards, EmptyCardInfo{
						CardID:       card.ID,
						NoteID:       note.ID,
						DeckID:       card.DeckID,
						TemplateName: card.TemplateName,
						Ordinal:      card.Ordinal,
						Front:        card.Front,
						Back:         card.Back,
						Reason:       "Card has no content (both front and back are empty)",
					})
				}
			}
		}
	}

	respondJSON(w, http.StatusOK, EmptyCardsResponse{
		Count:      len(emptyCards),
		EmptyCards: emptyCards,
	})
}

type DeleteEmptyCardsRequest struct {
	CardIDs []int64 `json:"cardIds"` // List of card IDs to delete
}

type DeleteEmptyCardsResponse struct {
	Deleted int      `json:"deleted"`
	Failed  []string `json:"failed,omitempty"` // Error messages for failed deletions
}

// DeleteEmptyCards deletes specified empty cards
func (h *APIHandler) DeleteEmptyCards(w http.ResponseWriter, r *http.Request) {
	var req DeleteEmptyCardsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if len(req.CardIDs) == 0 {
		http.Error(w, "No card IDs provided", http.StatusBadRequest)
		return
	}

	deleted := 0
	var failed []string

	for _, cardID := range req.CardIDs {
		if err := h.store.DeleteCard(cardID); err != nil {
			failed = append(failed, fmt.Sprintf("Card %d: %v", cardID, err))
			log.Printf("Failed to delete card %d: %v", cardID, err)
		} else {
			deleted++
			// Remove from collection cache
			delete(h.collection.Cards, cardID)
		}
	}

	respondJSON(w, http.StatusOK, DeleteEmptyCardsResponse{
		Deleted: deleted,
		Failed:  failed,
	})
}

// stripHTML removes HTML tags and returns plain text
func stripHTML(html string) string {
	return htmlPolicy.Sanitize(html)
}

// Backup endpoints

func (h *APIHandler) CreateBackup(w http.ResponseWriter, r *http.Request) {
	backupPath, err := h.backupManager.CreateBackup(h.collectionID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create backup: %v", err), http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusCreated, map[string]string{
		"message":    "Backup created successfully",
		"backupPath": backupPath,
		"timestamp":  time.Now().Format(time.RFC3339),
	})
}

type RestoreBackupRequest struct {
	BackupPath string `json:"backupPath"`
}

func (h *APIHandler) RestoreBackup(w http.ResponseWriter, r *http.Request) {
	var req RestoreBackupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.BackupPath == "" {
		http.Error(w, "backupPath is required", http.StatusBadRequest)
		return
	}

	// Warning: This will close the database connection and replace the database
	// In a production system, this would need more careful handling
	respondJSON(w, http.StatusOK, map[string]string{
		"message": "Restore functionality requires server restart. Use with caution.",
		"warning": "This operation will replace the current database.",
		"note":    "Implement with /api/backups/restore endpoint after server architecture supports it.",
	})
}

func (h *APIHandler) ListBackups(w http.ResponseWriter, r *http.Request) {
	files, err := filepath.Glob(filepath.Join(h.backupManager.backupDir, "microdote-backup-*.zip"))
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to list backups: %v", err), http.StatusInternalServerError)
		return
	}

	type backupInfo struct {
		Path     string    `json:"path"`
		Filename string    `json:"filename"`
		Size     int64     `json:"size"`
		Modified time.Time `json:"modified"`
	}

	var backups []backupInfo
	for _, path := range files {
		info, err := os.Stat(path)
		if err != nil {
			continue
		}

		backups = append(backups, backupInfo{
			Path:     path,
			Filename: filepath.Base(path),
			Size:     info.Size(),
			Modified: info.ModTime(),
		})
	}

	respondJSON(w, http.StatusOK, backups)
}

// Helper functions

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func parseIDParam(r *http.Request, paramName string) (int64, error) {
	idStr := chi.URLParam(r, paramName)
	return strconv.ParseInt(idStr, 10, 64)
}

// Main function to start the server
func main() {
	// Initialize database and collection
	log.Println("Initializing Microdote server...")
	col, store, err := InitDefaultCollection("./data/microdote.db")
	if err != nil {
		log.Fatalf("Failed to initialize collection: %v", err)
	}
	defer store.Close()

	log.Printf("Collection loaded with %d decks, %d notes, %d cards\n",
		len(col.Decks), len(col.Notes), len(col.Cards))

	// Create backup manager
	backupMgr := NewBackupManager("./data/microdote.db", "./backups", store)

	// Create API handler
	handler := NewAPIHandler(store, col, backupMgr)

	// Set up router
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RealIP)

	// CORS configuration for frontend
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// API routes
	r.Route("/api", func(r chi.Router) {
		// Health check
		r.Get("/health", handler.HealthCheck)

		// Collection
		r.Get("/collection", handler.GetCollection)

		// Decks
		r.Get("/decks", handler.ListDecks)
		r.Post("/decks", handler.CreateDeck)
		r.Get("/decks/{id}", handler.GetDeck)
		r.Get("/decks/{id}/stats", handler.GetDeckStats)
		r.Get("/decks/{deckId}/due", handler.GetDueCards)

		// Note Types
		r.Get("/note-types", handler.ListNoteTypes)
		r.Get("/note-types/{name}", handler.GetNoteType)
		r.Post("/note-types/{name}/fields", handler.AddField)
		r.Patch("/note-types/{name}/fields/rename", handler.RenameField)
		r.Delete("/note-types/{name}/fields", handler.RemoveField)
		r.Put("/note-types/{name}/fields/reorder", handler.ReorderFields)
		r.Put("/note-types/{name}/sort-field", handler.SetSortField)
		r.Put("/note-types/{name}/fields/options", handler.SetFieldOptions)
		r.Patch("/note-types/{name}/templates/{templateName}", handler.UpdateTemplate)

		// Notes
		r.Post("/notes", handler.CreateNote)
		r.Get("/notes/{id}", handler.GetNote)
		r.Post("/notes/check-duplicate", handler.CheckDuplicate)

		// Cards
		r.Get("/cards/{id}", handler.GetCard)
		r.Post("/cards/{id}/answer", handler.AnswerCard)
		r.Patch("/cards/{id}", handler.UpdateCard)
		r.Get("/cards/empty", handler.FindEmptyCards)
		r.Post("/cards/empty/delete", handler.DeleteEmptyCards)

		// Backups
		r.Post("/backups", handler.CreateBackup)
		r.Get("/backups", handler.ListBackups)
		r.Post("/backups/restore", handler.RestoreBackup)
	})

	// Start server
	port := ":8080"
	log.Printf("Server starting on http://localhost%s\n", port)
	log.Printf("API endpoints available at http://localhost%s/api\n", port)
	log.Println("Press Ctrl+C to stop")

	if err := http.ListenAndServe(port, r); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

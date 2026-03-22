package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"slices"
	"strings"
)

const (
	defaultAIMaxSuggestions = 3
	maxAICardSuggestions    = 5
)

type aiSuggestionProvider interface {
	Generate(context.Context, aiSuggestionInput) (*AICardSuggestionsResponse, error)
}

type aiSuggestionInput struct {
	SourceText        string
	NoteType          *NoteType
	ExistingFieldVals map[string]string
	MaxSuggestions    int
}

type disabledAISuggestionProvider struct {
	reason string
}

func (p *disabledAISuggestionProvider) Generate(_ context.Context, _ aiSuggestionInput) (*AICardSuggestionsResponse, error) {
	return nil, fmt.Errorf("%s", p.reason)
}

type devAISuggestionProvider struct{}

func (p *devAISuggestionProvider) Generate(_ context.Context, input aiSuggestionInput) (*AICardSuggestionsResponse, error) {
	pairs := extractStudyPairs(input.SourceText)
	if len(pairs) == 0 {
		trimmed := strings.TrimSpace(input.SourceText)
		if trimmed != "" {
			pairs = append(pairs, studyPair{
				Prompt: answerPromptForText(trimmed),
				Answer: trimmed,
				Raw:    trimmed,
			})
		}
	}

	limit := input.MaxSuggestions
	if limit <= 0 {
		limit = defaultAIMaxSuggestions
	}
	if limit > maxAICardSuggestions {
		limit = maxAICardSuggestions
	}
	if len(pairs) > limit {
		pairs = pairs[:limit]
	}

	suggestions := make([]AICardSuggestion, 0, len(pairs))
	for idx, pair := range pairs {
		fieldVals := mapStudyPairToFieldVals(input.NoteType, pair)
		if fieldValsBlank(fieldVals) {
			continue
		}
		title := strings.TrimSpace(pair.Prompt)
		if title == "" {
			title = fmt.Sprintf("Suggestion %d", idx+1)
		}
		suggestions = append(suggestions, AICardSuggestion{
			Title:     title,
			Rationale: "Generated locally from the pasted study material. Review and refine before saving.",
			FieldVals: fieldVals,
		})
	}

	return &AICardSuggestionsResponse{
		Suggestions: suggestions,
		Provider:    "dev",
		Model:       "heuristic",
	}, nil
}

type openAISuggestionProvider struct {
	cfg OpenAIConfig
}

func (p *openAISuggestionProvider) Generate(ctx context.Context, input aiSuggestionInput) (*AICardSuggestionsResponse, error) {
	payload := map[string]any{
		"model": p.cfg.Model,
		"input": []map[string]any{
			{
				"role": "system",
				"content": []map[string]any{
					{
						"type": "input_text",
						"text": buildAISystemPrompt(input.NoteType),
					},
				},
			},
			{
				"role": "user",
				"content": []map[string]any{
					{
						"type": "input_text",
						"text": buildAIUserPrompt(input),
					},
				},
			},
		},
		"text": map[string]any{
			"format": map[string]any{
				"type":   "json_schema",
				"name":   "vutadex_card_suggestions",
				"strict": true,
				"schema": buildAISuggestionSchema(input.NoteType.Fields, input.MaxSuggestions),
			},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.cfg.BaseURL+"/responses", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+p.cfg.APIKey)
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var response struct {
		OutputText string `json:"output_text"`
		Output     []struct {
			Content []struct {
				Type string `json:"type"`
				Text string `json:"text"`
			} `json:"content"`
		} `json:"output"`
		Error *struct {
			Message string `json:"message"`
		} `json:"error,omitempty"`
	}

	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("decode OpenAI response: %w", err)
	}

	if res.StatusCode >= http.StatusBadRequest {
		if response.Error != nil && strings.TrimSpace(response.Error.Message) != "" {
			return nil, fmt.Errorf("%s", response.Error.Message)
		}
		return nil, fmt.Errorf("OpenAI request failed with status %d", res.StatusCode)
	}

	outputText := strings.TrimSpace(response.OutputText)
	if outputText == "" {
		outputText = strings.TrimSpace(extractResponseOutputText(response.Output))
	}
	if outputText == "" {
		return nil, fmt.Errorf("OpenAI returned no structured output")
	}

	var parsed struct {
		Suggestions []AICardSuggestion `json:"suggestions"`
	}
	if err := json.Unmarshal([]byte(outputText), &parsed); err != nil {
		return nil, fmt.Errorf("parse OpenAI structured output: %w", err)
	}

	return &AICardSuggestionsResponse{
		Suggestions: normalizeAISuggestions(parsed.Suggestions, input.NoteType.Fields, input.MaxSuggestions),
		Provider:    "openai",
		Model:       p.cfg.Model,
	}, nil
}

func newAISuggestionProvider(cfg AppConfig) aiSuggestionProvider {
	if strings.TrimSpace(cfg.OpenAI.APIKey) != "" {
		return &openAISuggestionProvider{cfg: cfg.OpenAI}
	}
	if cfg.IsDevelopment() {
		return &devAISuggestionProvider{}
	}
	return &disabledAISuggestionProvider{
		reason: "AI suggestions are not configured. Set VUTADEX_OPENAI_API_KEY to enable them.",
	}
}

func (h *APIHandler) GenerateCardSuggestions(w http.ResponseWriter, r *http.Request) {
	col, collectionID, err := h.collectionForRequest(r)
	if err != nil {
		respondAPIError(w, http.StatusInternalServerError, "collection_load_failed", err.Error())
		return
	}

	var req GenerateAICardSuggestionsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondAPIError(w, http.StatusBadRequest, "invalid_request", "Invalid AI suggestion request.")
		return
	}

	req.SourceText = strings.TrimSpace(req.SourceText)
	req.NoteType = strings.TrimSpace(req.NoteType)
	if req.SourceText == "" || req.NoteType == "" {
		respondAPIError(w, http.StatusBadRequest, "invalid_ai_request", "Source text and note type are required.")
		return
	}
	if len(req.SourceText) > 20000 {
		respondAPIError(w, http.StatusBadRequest, "ai_source_too_large", "Source text is too large. Keep it under 20,000 characters for now.")
		return
	}

	noteType, ok := col.NoteTypes[NoteTypeName(req.NoteType)]
	if !ok {
		reloaded, err := h.store.GetNoteType(collectionID, NoteTypeName(req.NoteType))
		if err != nil {
			respondAPIError(w, http.StatusBadRequest, "invalid_note_type", "Note type not found.")
			return
		}
		noteType = *reloaded
	}

	provider := newAISuggestionProvider(h.config)
	response, err := provider.Generate(r.Context(), aiSuggestionInput{
		SourceText:        req.SourceText,
		NoteType:          &noteType,
		ExistingFieldVals: req.ExistingFieldVals,
		MaxSuggestions:    clampAISuggestionCount(req.MaxSuggestions),
	})
	if err != nil {
		status := http.StatusBadGateway
		code := "ai_suggestions_failed"
		if _, disabled := provider.(*disabledAISuggestionProvider); disabled {
			status = http.StatusNotImplemented
			code = "ai_suggestions_not_configured"
		}
		respondAPIError(w, status, code, err.Error())
		return
	}

	response.Suggestions = normalizeAISuggestions(response.Suggestions, noteType.Fields, clampAISuggestionCount(req.MaxSuggestions))
	respondJSON(w, http.StatusOK, response)
}

type studyPair struct {
	Prompt string
	Answer string
	Raw    string
}

func extractStudyPairs(source string) []studyPair {
	lines := strings.Split(source, "\n")
	pairs := make([]studyPair, 0)
	for _, rawLine := range lines {
		line := strings.TrimSpace(rawLine)
		if line == "" {
			continue
		}
		for _, separator := range []string{":", " - ", " – ", " — ", "\t"} {
			parts := strings.SplitN(line, separator, 2)
			if len(parts) != 2 {
				continue
			}
			prompt := strings.TrimSpace(parts[0])
			answer := strings.TrimSpace(parts[1])
			if prompt == "" || answer == "" {
				continue
			}
			pairs = append(pairs, studyPair{
				Prompt: prompt,
				Answer: answer,
				Raw:    line,
			})
			break
		}
	}
	return pairs
}

func answerPromptForText(source string) string {
	sentences := regexp.MustCompile(`[.!?]\s+`).Split(source, 2)
	first := strings.TrimSpace(sentences[0])
	if first == "" {
		first = strings.TrimSpace(source)
	}
	if len(first) > 90 {
		first = strings.TrimSpace(first[:90]) + "..."
	}
	return "Explain: " + first
}

func mapStudyPairToFieldVals(noteType *NoteType, pair studyPair) map[string]string {
	fieldVals := make(map[string]string, len(noteType.Fields))
	if len(noteType.Fields) == 0 {
		return fieldVals
	}

	if hasField(noteType.Fields, "Front") && hasField(noteType.Fields, "Back") {
		fieldVals["Front"] = pair.Prompt
		fieldVals["Back"] = pair.Answer
		fillRemainingFields(fieldVals, noteType.Fields, pair)
		return fieldVals
	}

	if strings.EqualFold(string(noteType.Name), "Cloze") && hasField(noteType.Fields, "Text") {
		fieldVals["Text"] = fmt.Sprintf("%s means {{c1::%s}}.", pair.Prompt, pair.Answer)
		if hasField(noteType.Fields, "Extra") {
			fieldVals["Extra"] = pair.Raw
		}
		fillRemainingFields(fieldVals, noteType.Fields, pair)
		return fieldVals
	}

	fillRemainingFields(fieldVals, noteType.Fields, pair)
	return fieldVals
}

func fillRemainingFields(fieldVals map[string]string, fields []string, pair studyPair) {
	seedValues := []string{pair.Prompt, pair.Answer, pair.Raw}
	index := 0
	for _, field := range fields {
		if _, exists := fieldVals[field]; exists {
			continue
		}
		if index < len(seedValues) {
			fieldVals[field] = seedValues[index]
			index++
			continue
		}
		fieldVals[field] = ""
	}
}

func hasField(fields []string, target string) bool {
	return slices.Contains(fields, target)
}

func fieldValsBlank(fieldVals map[string]string) bool {
	for _, value := range fieldVals {
		if strings.TrimSpace(value) != "" {
			return false
		}
	}
	return true
}

func clampAISuggestionCount(value int) int {
	if value <= 0 {
		return defaultAIMaxSuggestions
	}
	if value > maxAICardSuggestions {
		return maxAICardSuggestions
	}
	return value
}

func buildAISystemPrompt(noteType *NoteType) string {
	base := []string{
		"You generate high-signal flashcard note suggestions for Vutadex.",
		"Return only facts supported by the source text.",
		"Keep suggestions concise, study-ready, and non-duplicative.",
		"If the source text does not support a solid flashcard, return an empty suggestions list.",
		"Always return complete field values for every required field.",
	}

	if strings.EqualFold(string(noteType.Name), "Cloze") && hasField(noteType.Fields, "Text") {
		base = append(base, "When the note type is Cloze, use Anki-style cloze markers like {{c1::answer}} inside the Text field.")
	}
	if hasField(noteType.Fields, "Front") && hasField(noteType.Fields, "Back") {
		base = append(base, "For Front/Back note types, make Front a retrieval prompt and Back the concise answer.")
	}

	return strings.Join(base, "\n")
}

func buildAIUserPrompt(input aiSuggestionInput) string {
	var b strings.Builder
	b.WriteString("Create flashcard note suggestions from this source material.\n\n")
	fmt.Fprintf(&b, "Note type: %s\n", input.NoteType.Name)
	fmt.Fprintf(&b, "Fields: %s\n", strings.Join(input.NoteType.Fields, ", "))
	fmt.Fprintf(&b, "Maximum suggestions: %d\n", input.MaxSuggestions)
	if len(input.ExistingFieldVals) > 0 {
		b.WriteString("Current draft values:\n")
		for _, field := range input.NoteType.Fields {
			if value := strings.TrimSpace(input.ExistingFieldVals[field]); value != "" {
				fmt.Fprintf(&b, "- %s: %s\n", field, value)
			}
		}
		b.WriteString("\n")
	}
	b.WriteString("Source material:\n")
	b.WriteString(input.SourceText)
	return b.String()
}

func buildAISuggestionSchema(fields []string, maxSuggestions int) map[string]any {
	fieldProps := make(map[string]any, len(fields))
	requiredFields := make([]string, 0, len(fields))
	for _, field := range fields {
		fieldProps[field] = map[string]any{
			"type": "string",
		}
		requiredFields = append(requiredFields, field)
	}

	return map[string]any{
		"type":                 "object",
		"additionalProperties": false,
		"properties": map[string]any{
			"suggestions": map[string]any{
				"type":     "array",
				"maxItems": maxSuggestions,
				"items": map[string]any{
					"type":                 "object",
					"additionalProperties": false,
					"properties": map[string]any{
						"title": map[string]any{
							"type": "string",
						},
						"rationale": map[string]any{
							"type": "string",
						},
						"fieldVals": map[string]any{
							"type":                 "object",
							"additionalProperties": false,
							"properties":           fieldProps,
							"required":             requiredFields,
						},
					},
					"required": []string{"title", "rationale", "fieldVals"},
				},
			},
		},
		"required": []string{"suggestions"},
	}
}

func extractResponseOutputText(output []struct {
	Content []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"content"`
}) string {
	var parts []string
	for _, item := range output {
		for _, content := range item.Content {
			if strings.TrimSpace(content.Text) == "" {
				continue
			}
			parts = append(parts, content.Text)
		}
	}
	return strings.Join(parts, "\n")
}

func normalizeAISuggestions(raw []AICardSuggestion, fields []string, maxSuggestions int) []AICardSuggestion {
	limit := clampAISuggestionCount(maxSuggestions)
	suggestions := make([]AICardSuggestion, 0, len(raw))
	for _, suggestion := range raw {
		fieldVals := make(map[string]string, len(fields))
		for _, field := range fields {
			fieldVals[field] = strings.TrimSpace(suggestion.FieldVals[field])
		}
		if fieldValsBlank(fieldVals) {
			continue
		}
		title := strings.TrimSpace(suggestion.Title)
		if title == "" {
			for _, field := range fields {
				if value := fieldVals[field]; value != "" {
					title = value
					break
				}
			}
		}
		if title == "" {
			title = "Suggested note"
		}
		suggestions = append(suggestions, AICardSuggestion{
			Title:     title,
			Rationale: strings.TrimSpace(suggestion.Rationale),
			FieldVals: fieldVals,
		})
		if len(suggestions) >= limit {
			break
		}
	}
	return suggestions
}

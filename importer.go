package main

import (
	"archive/zip"
	"bytes"
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

type importParseOptions struct {
	Source          string
	FormatHint      string
	Filename        string
	DefaultDeckName string
	DefaultNoteType string
}

type importNormalizedNote struct {
	DeckName string
	NoteType NoteTypeName
	Fields   map[string]string
	Tags     []string
}

type nativeImportPayload struct {
	Deck     string             `json:"deck" yaml:"deck"`
	NoteType string             `json:"noteType" yaml:"noteType"`
	Notes    []nativeImportNote `json:"notes" yaml:"notes"`
	Decks    []nativeImportDeck `json:"decks" yaml:"decks"`
}

type nativeImportDeck struct {
	Name     string             `json:"name" yaml:"name"`
	NoteType string             `json:"noteType" yaml:"noteType"`
	Notes    []nativeImportNote `json:"notes" yaml:"notes"`
}

type nativeImportNote struct {
	Deck     string            `json:"deck" yaml:"deck"`
	NoteType string            `json:"noteType" yaml:"noteType"`
	Front    string            `json:"front" yaml:"front"`
	Back     string            `json:"back" yaml:"back"`
	Text     string            `json:"text" yaml:"text"`
	Extra    string            `json:"extra" yaml:"extra"`
	Fields   map[string]string `json:"fields" yaml:"fields"`
	Tags     []string          `json:"tags" yaml:"tags"`
}

type importParserResult struct {
	Notes  []importNormalizedNote
	Source string
	Format string
}

type ankiDeckMeta struct {
	Name string `json:"name"`
}

type ankiModelField struct {
	Name string `json:"name"`
}

type ankiModelMeta struct {
	Name string           `json:"name"`
	Type int              `json:"type"`
	Flds []ankiModelField `json:"flds"`
}

func parseImportData(data []byte, opts importParseOptions) (importParserResult, error) {
	source := normalizeImportSource(opts.Source)
	format := normalizeImportFormat(opts.FormatHint)
	if format == "" {
		format = detectImportFormat(opts.Filename, data)
	}

	if source == "auto" {
		switch format {
		case "json", "yaml":
			source = "native"
		case "apkg", "colpkg":
			source = "anki"
		default:
			source = "auto"
		}
	}

	switch source {
	case "native":
		notes, usedFormat, err := parseNativeImport(data, format, opts)
		if err != nil {
			return importParserResult{}, err
		}
		return importParserResult{Notes: notes, Source: "native", Format: usedFormat}, nil
	case "anki":
		if format == "apkg" || format == "colpkg" {
			notes, err := parseAnkiPackageImport(data, opts)
			if err != nil {
				return importParserResult{}, err
			}
			return importParserResult{Notes: notes, Source: "anki", Format: format}, nil
		}
		notes, usedFormat, err := parseDelimitedImport(data, "anki", opts)
		if err != nil {
			return importParserResult{}, err
		}
		return importParserResult{Notes: notes, Source: "anki", Format: usedFormat}, nil
	case "quizlet":
		notes, usedFormat, err := parseDelimitedImport(data, "quizlet", opts)
		if err != nil {
			return importParserResult{}, err
		}
		return importParserResult{Notes: notes, Source: "quizlet", Format: usedFormat}, nil
	case "auto":
		if format == "apkg" || format == "colpkg" {
			notes, err := parseAnkiPackageImport(data, opts)
			if err != nil {
				return importParserResult{}, err
			}
			return importParserResult{Notes: notes, Source: "anki", Format: format}, nil
		}
		notes, usedFormat, err := parseDelimitedImport(data, "auto", opts)
		if err != nil {
			return importParserResult{}, err
		}
		detectedSource := "anki"
		if looksLikeQuizlet(data) {
			detectedSource = "quizlet"
		}
		return importParserResult{Notes: notes, Source: detectedSource, Format: usedFormat}, nil
	default:
		return importParserResult{}, fmt.Errorf("unsupported import source: %s", source)
	}
}

func normalizeImportSource(source string) string {
	s := strings.ToLower(strings.TrimSpace(source))
	if s == "" {
		return "auto"
	}
	switch s {
	case "auto", "native", "anki", "quizlet":
		return s
	default:
		return "auto"
	}
}

func normalizeImportFormat(format string) string {
	s := strings.ToLower(strings.TrimSpace(format))
	s = strings.TrimPrefix(s, ".")
	switch s {
	case "json", "yaml", "yml", "csv", "tsv", "txt", "apkg", "colpkg":
		if s == "yml" {
			return "yaml"
		}
		return s
	default:
		return ""
	}
}

func detectImportFormat(filename string, data []byte) string {
	ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(filename), "."))
	switch ext {
	case "json":
		return "json"
	case "yaml", "yml":
		return "yaml"
	case "csv", "tsv", "txt", "apkg", "colpkg":
		return ext
	}

	trimmed := strings.TrimSpace(string(data))
	if strings.HasPrefix(trimmed, "{") || strings.HasPrefix(trimmed, "[") {
		return "json"
	}
	if strings.Contains(trimmed, "notes:") || strings.Contains(trimmed, "decks:") {
		return "yaml"
	}
	return "txt"
}

func parseNativeImport(data []byte, format string, opts importParseOptions) ([]importNormalizedNote, string, error) {
	if format == "" || format == "txt" || format == "csv" || format == "tsv" {
		format = detectImportFormat(opts.Filename, data)
	}

	var payload nativeImportPayload
	switch format {
	case "json":
		if err := json.Unmarshal(data, &payload); err != nil {
			var notesOnly []nativeImportNote
			if errList := json.Unmarshal(data, &notesOnly); errList != nil {
				return nil, "json", fmt.Errorf("invalid JSON import payload: %w", err)
			}
			payload.Notes = notesOnly
		}
	case "yaml":
		if err := yaml.Unmarshal(data, &payload); err != nil {
			var notesOnly []nativeImportNote
			if errList := yaml.Unmarshal(data, &notesOnly); errList != nil {
				return nil, "yaml", fmt.Errorf("invalid YAML import payload: %w", err)
			}
			payload.Notes = notesOnly
		}
	default:
		return nil, format, fmt.Errorf("native import expects JSON or YAML, got %s", format)
	}

	notes, err := normalizeNativePayload(payload, opts)
	if err != nil {
		return nil, format, err
	}
	return notes, format, nil
}

func normalizeNativePayload(payload nativeImportPayload, opts importParseOptions) ([]importNormalizedNote, error) {
	baseDeck := firstNonEmpty(payload.Deck, opts.DefaultDeckName)
	baseType := firstNonEmpty(payload.NoteType, opts.DefaultNoteType, string(NoteTypeName("Basic")))

	var out []importNormalizedNote
	for _, n := range payload.Notes {
		note, err := normalizeNativeNote(n, baseDeck, baseType)
		if err != nil {
			return nil, err
		}
		out = append(out, note)
	}

	for _, deck := range payload.Decks {
		deckName := firstNonEmpty(deck.Name, baseDeck)
		deckType := firstNonEmpty(deck.NoteType, baseType)
		for _, n := range deck.Notes {
			note, err := normalizeNativeNote(n, deckName, deckType)
			if err != nil {
				return nil, err
			}
			out = append(out, note)
		}
	}

	if len(out) == 0 {
		return nil, errors.New("no notes found in native import file")
	}
	return out, nil
}

func normalizeNativeNote(note nativeImportNote, defaultDeckName, defaultType string) (importNormalizedNote, error) {
	noteType := inferNoteType(firstNonEmpty(note.NoteType, defaultType), note.Text, note.Front)
	deckName := firstNonEmpty(note.Deck, defaultDeckName)

	fields := map[string]string{}
	for k, v := range note.Fields {
		fields[k] = v
	}

	if noteType == "Cloze" {
		if _, ok := getFieldValueCaseInsensitive(fields, "Text"); !ok {
			fields["Text"] = firstNonEmpty(note.Text, note.Front)
		}
		if _, ok := getFieldValueCaseInsensitive(fields, "Extra"); !ok {
			fields["Extra"] = firstNonEmpty(note.Extra, note.Back)
		}
	} else {
		if _, ok := getFieldValueCaseInsensitive(fields, "Front"); !ok {
			fields["Front"] = firstNonEmpty(note.Front, note.Text)
		}
		if _, ok := getFieldValueCaseInsensitive(fields, "Back"); !ok {
			fields["Back"] = firstNonEmpty(note.Back, note.Extra)
		}
	}

	if len(fields) == 0 {
		return importNormalizedNote{}, errors.New("native note has no field content")
	}

	return importNormalizedNote{
		DeckName: deckName,
		NoteType: noteType,
		Fields:   fields,
		Tags:     dedupeTags(note.Tags),
	}, nil
}

func parseDelimitedImport(data []byte, source string, opts importParseOptions) ([]importNormalizedNote, string, error) {
	delimiter := detectDelimitedSeparator(data)
	format := "csv"
	if delimiter == '\t' {
		format = "tsv"
	}

	reader := csv.NewReader(bytes.NewReader(data))
	reader.Comma = delimiter
	reader.FieldsPerRecord = -1
	reader.LazyQuotes = true

	records, err := reader.ReadAll()
	if err != nil {
		return nil, format, fmt.Errorf("failed to parse delimited import: %w", err)
	}
	if len(records) == 0 {
		return nil, format, errors.New("import file is empty")
	}

	start := 0
	columnIndex := map[string]int{}
	if isDelimitedHeader(records[0]) {
		start = 1
		for i, rawCol := range records[0] {
			columnIndex[normalizeColumnName(rawCol)] = i
		}
	}

	var notes []importNormalizedNote
	for _, row := range records[start:] {
		if rowIsEmpty(row) {
			continue
		}

		note, ok := parseDelimitedRow(row, columnIndex, source, opts)
		if ok {
			notes = append(notes, note)
		}
	}

	if len(notes) == 0 {
		return nil, format, errors.New("no importable rows found")
	}

	return notes, format, nil
}

func parseDelimitedRow(row []string, columnIndex map[string]int, source string, opts importParseOptions) (importNormalizedNote, bool) {
	valueFor := func(names ...string) string {
		for _, name := range names {
			if idx, ok := columnIndex[name]; ok && idx >= 0 && idx < len(row) {
				return strings.TrimSpace(row[idx])
			}
		}
		return ""
	}

	valueAt := func(i int) string {
		if i >= 0 && i < len(row) {
			return strings.TrimSpace(row[i])
		}
		return ""
	}

	hasHeader := len(columnIndex) > 0
	front := ""
	back := ""
	text := ""
	extra := ""
	tagsRaw := ""
	deckName := ""
	rawType := ""

	if hasHeader {
		front = valueFor("front", "question", "prompt", "term")
		back = valueFor("back", "answer", "definition")
		text = valueFor("text", "cloze")
		extra = valueFor("extra", "hint")
		tagsRaw = valueFor("tags", "tag")
		deckName = valueFor("deck", "deckname")
		rawType = valueFor("type", "notetype", "model")
	} else {
		front = valueAt(0)
		back = valueAt(1)
		tagsRaw = valueAt(2)
		deckName = valueAt(3)
		rawType = valueAt(4)
	}

	if source == "quizlet" {
		rawType = "Basic"
		if front == "" {
			front = valueAt(0)
		}
		if back == "" {
			back = valueAt(1)
		}
	}

	noteType := inferNoteType(firstNonEmpty(rawType, opts.DefaultNoteType), text, front)
	if source == "quizlet" {
		noteType = "Basic"
	}

	finalDeck := firstNonEmpty(deckName, opts.DefaultDeckName)
	fields := map[string]string{}

	if noteType == "Cloze" {
		fields["Text"] = firstNonEmpty(text, front)
		fields["Extra"] = firstNonEmpty(extra, back)
		if strings.TrimSpace(fields["Text"]) == "" {
			return importNormalizedNote{}, false
		}
	} else {
		fields["Front"] = firstNonEmpty(front, text)
		fields["Back"] = firstNonEmpty(back, extra)
		if strings.TrimSpace(fields["Front"]) == "" && strings.TrimSpace(fields["Back"]) == "" {
			return importNormalizedNote{}, false
		}
	}

	return importNormalizedNote{
		DeckName: finalDeck,
		NoteType: noteType,
		Fields:   fields,
		Tags:     splitTags(tagsRaw),
	}, true
}

func parseAnkiPackageImport(data []byte, opts importParseOptions) ([]importNormalizedNote, error) {
	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return nil, fmt.Errorf("failed to read package: %w", err)
	}

	var collectionEntry *zip.File
	for _, file := range zr.File {
		base := strings.ToLower(filepath.Base(file.Name))
		if base == "collection.anki2" || base == "collection.anki21" {
			collectionEntry = file
			break
		}
	}
	if collectionEntry == nil {
		return nil, errors.New("Anki package missing collection database (collection.anki2/collection.anki21)")
	}

	tempDir, err := os.MkdirTemp("", "microdote-anki-import-")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp dir: %w", err)
	}
	defer os.RemoveAll(tempDir)

	tempDBPath := filepath.Join(tempDir, filepath.Base(collectionEntry.Name))
	rc, err := collectionEntry.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open collection database: %w", err)
	}
	defer rc.Close()

	outFile, err := os.Create(tempDBPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create temp collection file: %w", err)
	}
	if _, err := io.Copy(outFile, rc); err != nil {
		outFile.Close()
		return nil, fmt.Errorf("failed to copy collection database: %w", err)
	}
	if err := outFile.Close(); err != nil {
		return nil, fmt.Errorf("failed to close temp collection file: %w", err)
	}

	db, err := sql.Open("sqlite3", tempDBPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open Anki collection database: %w", err)
	}
	defer db.Close()

	var decksJSON string
	var modelsJSON string
	if err := db.QueryRow(`SELECT decks, models FROM col LIMIT 1`).Scan(&decksJSON, &modelsJSON); err != nil {
		return nil, fmt.Errorf("failed to read Anki collection metadata: %w", err)
	}

	deckNames := map[string]string{}
	if decksJSON != "" {
		if err := json.Unmarshal([]byte(decksJSON), &deckNames); err != nil {
			var rawDecks map[string]ankiDeckMeta
			if errRaw := json.Unmarshal([]byte(decksJSON), &rawDecks); errRaw == nil {
				for id, deck := range rawDecks {
					deckNames[id] = deck.Name
				}
			}
		}
	}

	modelMap := map[string]ankiModelMeta{}
	if modelsJSON != "" {
		if err := json.Unmarshal([]byte(modelsJSON), &modelMap); err != nil {
			return nil, fmt.Errorf("failed to parse Anki model metadata: %w", err)
		}
	}

	noteDeckIDs := map[int64]int64{}
	cardRows, err := db.Query(`SELECT nid, MIN(did) FROM cards GROUP BY nid`)
	if err != nil {
		return nil, fmt.Errorf("failed to read Anki cards: %w", err)
	}
	for cardRows.Next() {
		var nid int64
		var did int64
		if err := cardRows.Scan(&nid, &did); err != nil {
			cardRows.Close()
			return nil, fmt.Errorf("failed to scan Anki cards: %w", err)
		}
		noteDeckIDs[nid] = did
	}
	if err := cardRows.Err(); err != nil {
		cardRows.Close()
		return nil, fmt.Errorf("failed iterating Anki cards: %w", err)
	}
	cardRows.Close()

	noteRows, err := db.Query(`SELECT id, mid, tags, flds FROM notes`)
	if err != nil {
		return nil, fmt.Errorf("failed to read Anki notes: %w", err)
	}
	defer noteRows.Close()

	var out []importNormalizedNote
	for noteRows.Next() {
		var nid int64
		var mid int64
		var rawTags string
		var rawFields string
		if err := noteRows.Scan(&nid, &mid, &rawTags, &rawFields); err != nil {
			return nil, fmt.Errorf("failed to scan Anki note: %w", err)
		}

		fields := strings.Split(rawFields, "\x1f")
		model := modelMap[strconv.FormatInt(mid, 10)]
		noteType := inferAnkiPackageNoteType(model, fields, opts.DefaultNoteType)

		deckName := firstNonEmpty(opts.DefaultDeckName)
		if did, ok := noteDeckIDs[nid]; ok {
			if dName, okDeck := deckNames[strconv.FormatInt(did, 10)]; okDeck {
				deckName = firstNonEmpty(dName, deckName)
			}
		}

		fieldMap := buildFieldMapFromAnkiNote(model, fields, noteType)
		if len(fieldMap) == 0 {
			continue
		}

		out = append(out, importNormalizedNote{
			DeckName: deckName,
			NoteType: noteType,
			Fields:   fieldMap,
			Tags:     splitTags(rawTags),
		})
	}
	if err := noteRows.Err(); err != nil {
		return nil, fmt.Errorf("failed iterating Anki notes: %w", err)
	}

	if len(out) == 0 {
		return nil, errors.New("no notes found in Anki package")
	}
	return out, nil
}

func inferAnkiPackageNoteType(model ankiModelMeta, values []string, defaultType string) NoteTypeName {
	if model.Type == 1 || strings.Contains(strings.ToLower(model.Name), "cloze") {
		return "Cloze"
	}
	if len(values) > 0 && strings.Contains(values[0], "{{c") && strings.Contains(values[0], "::") {
		return "Cloze"
	}
	return inferNoteType(defaultType, firstSliceValue(values, 0), "")
}

func buildFieldMapFromAnkiNote(model ankiModelMeta, values []string, noteType NoteTypeName) map[string]string {
	mapped := map[string]string{}
	for idx, field := range model.Flds {
		if idx >= len(values) {
			break
		}
		mapped[field.Name] = values[idx]
	}

	if noteType == "Cloze" {
		textVal := firstNonEmpty(getFieldCaseInsensitive(mapped, []string{"Text", "Front"}), firstSliceValue(values, 0))
		extraVal := firstNonEmpty(getFieldCaseInsensitive(mapped, []string{"Extra", "Back"}), firstSliceValue(values, 1))
		if strings.TrimSpace(textVal) == "" {
			return nil
		}
		return map[string]string{"Text": textVal, "Extra": extraVal}
	}

	frontVal := firstNonEmpty(getFieldCaseInsensitive(mapped, []string{"Front", "Text", "Question", "Term"}), firstSliceValue(values, 0))
	backVal := firstNonEmpty(getFieldCaseInsensitive(mapped, []string{"Back", "Answer", "Definition", "Extra"}), firstSliceValue(values, 1))
	if strings.TrimSpace(frontVal) == "" && strings.TrimSpace(backVal) == "" {
		return nil
	}
	return map[string]string{"Front": frontVal, "Back": backVal}
}

func inferNoteType(rawType string, text string, fallbackFront string) NoteTypeName {
	typeLower := strings.ToLower(strings.TrimSpace(rawType))
	if typeLower == "cloze" || strings.Contains(typeLower, "cloze") {
		return "Cloze"
	}
	if typeLower == "basic" {
		return "Basic"
	}
	if strings.Contains(text, "{{c") && strings.Contains(text, "::") {
		return "Cloze"
	}
	if strings.Contains(fallbackFront, "{{c") && strings.Contains(fallbackFront, "::") {
		return "Cloze"
	}
	if strings.TrimSpace(rawType) != "" {
		return NoteTypeName(strings.TrimSpace(rawType))
	}
	return "Basic"
}

func detectDelimitedSeparator(data []byte) rune {
	lines := strings.Split(string(data), "\n")
	for _, raw := range lines {
		line := strings.TrimSpace(raw)
		if line == "" {
			continue
		}
		tabCount := strings.Count(line, "\t")
		commaCount := strings.Count(line, ",")
		if tabCount >= commaCount {
			return '\t'
		}
		return ','
	}
	return '\t'
}

func isDelimitedHeader(row []string) bool {
	if len(row) == 0 {
		return false
	}
	known := map[string]bool{
		"front": true, "back": true, "text": true, "extra": true,
		"tags": true, "tag": true, "deck": true, "deckname": true,
		"type": true, "notetype": true, "model": true,
		"term": true, "definition": true, "question": true, "answer": true,
	}

	for _, col := range row {
		if known[normalizeColumnName(col)] {
			return true
		}
	}
	return false
}

func normalizeColumnName(s string) string {
	normalized := strings.ToLower(strings.TrimSpace(s))
	replacer := strings.NewReplacer(" ", "", "_", "", "-", "")
	return replacer.Replace(normalized)
}

func rowIsEmpty(row []string) bool {
	for _, col := range row {
		if strings.TrimSpace(col) != "" {
			return false
		}
	}
	return true
}

func looksLikeQuizlet(data []byte) bool {
	lines := strings.Split(string(data), "\n")
	for _, raw := range lines {
		line := strings.TrimSpace(raw)
		if line == "" {
			continue
		}
		lower := strings.ToLower(line)
		return strings.Contains(lower, "term") && strings.Contains(lower, "definition")
	}
	return false
}

func splitTags(raw string) []string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil
	}

	var parts []string
	if strings.Contains(raw, ",") {
		for _, p := range strings.Split(raw, ",") {
			trimmed := strings.TrimSpace(p)
			if trimmed != "" {
				parts = append(parts, trimmed)
			}
		}
	} else {
		parts = strings.Fields(raw)
	}

	return dedupeTags(parts)
}

func dedupeTags(tags []string) []string {
	if len(tags) == 0 {
		return nil
	}
	seen := make(map[string]struct{}, len(tags))
	out := make([]string, 0, len(tags))
	for _, tag := range tags {
		trimmed := strings.TrimSpace(tag)
		if trimmed == "" {
			continue
		}
		if _, exists := seen[trimmed]; exists {
			continue
		}
		seen[trimmed] = struct{}{}
		out = append(out, trimmed)
	}
	return out
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		trimmed := strings.TrimSpace(v)
		if trimmed != "" {
			return trimmed
		}
	}
	return ""
}

func firstSliceValue(values []string, index int) string {
	if index >= 0 && index < len(values) {
		return strings.TrimSpace(values[index])
	}
	return ""
}

func getFieldValueCaseInsensitive(fields map[string]string, target string) (string, bool) {
	for key, value := range fields {
		if strings.EqualFold(key, target) {
			return value, true
		}
	}
	return "", false
}

func getFieldCaseInsensitive(fields map[string]string, candidates []string) string {
	for _, candidate := range candidates {
		if value, ok := getFieldValueCaseInsensitive(fields, candidate); ok {
			return value
		}
	}
	return ""
}

func sortedKeys(m map[string]struct{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

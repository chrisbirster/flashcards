package main

import (
	"archive/zip"
	"bytes"
	"database/sql"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func doMultipartImportRequest(t *testing.T, router http.Handler, fields map[string]string, filename string, content []byte) *httptest.ResponseRecorder {
	t.Helper()

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	for k, v := range fields {
		if err := writer.WriteField(k, v); err != nil {
			t.Fatalf("failed to write form field %s: %v", k, err)
		}
	}

	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		t.Fatalf("failed to create multipart file part: %v", err)
	}
	if _, err := part.Write(content); err != nil {
		t.Fatalf("failed to write multipart file content: %v", err)
	}

	if err := writer.Close(); err != nil {
		t.Fatalf("failed to close multipart writer: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/import", &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	return resp
}

func findNoteByType(t *testing.T, notes map[int64]Note, noteType NoteTypeName) (Note, bool) {
	t.Helper()
	for _, note := range notes {
		if note.Type == noteType {
			return note, true
		}
	}
	return Note{}, false
}

func containsExact(values []string, target string) bool {
	for _, v := range values {
		if v == target {
			return true
		}
	}
	return false
}

func TestAPI_ImportNotes_NativeJSONAndYAML(t *testing.T) {
	t.Run("json", func(t *testing.T) {
		env := setupAPITestEnv(t)

		jsonPayload := `{
  "deck": "JSON Deck",
  "notes": [
    {"front": "What is BFS?", "back": "Breadth-first search", "tags": ["graph", "bfs"]},
    {"noteType": "Cloze", "deck": "Cloze Deck", "text": "{{c1::BFS}} uses a queue", "extra": "Graph traversal"}
  ]
}`

		resp := doMultipartImportRequest(t, env.router, map[string]string{"source": "native"}, "dsa.json", []byte(jsonPayload))
		if resp.Code != http.StatusOK {
			t.Fatalf("expected import 200, got %d: %s", resp.Code, resp.Body.String())
		}

		result := decodeJSON[ImportNotesResponse](t, resp)
		if result.Imported != 2 {
			t.Fatalf("expected 2 imported notes, got %+v", result)
		}
		if !containsExact(result.DecksCreated, "JSON Deck") || !containsExact(result.DecksCreated, "Cloze Deck") {
			t.Fatalf("expected created decks to include JSON Deck and Cloze Deck, got %+v", result.DecksCreated)
		}

		notes, err := env.store.ListNotes("default")
		if err != nil {
			t.Fatalf("failed to list notes: %v", err)
		}
		if len(notes) != 2 {
			t.Fatalf("expected 2 stored notes, got %d", len(notes))
		}

		clozeNote, ok := findNoteByType(t, notes, "Cloze")
		if !ok {
			t.Fatalf("expected one cloze note in imported set")
		}
		if !strings.Contains(clozeNote.FieldMap["Text"], "{{c1::BFS}}") {
			t.Fatalf("expected cloze text to be preserved, got %q", clozeNote.FieldMap["Text"])
		}
	})

	t.Run("yaml", func(t *testing.T) {
		env := setupAPITestEnv(t)

		yamlPayload := "deck: YAML Deck\nnotes:\n  - front: \"Binary\\tSearch\"\n    back: \"O(log n)\"\n"

		resp := doMultipartImportRequest(t, env.router, map[string]string{"source": "native"}, "dsa.yaml", []byte(yamlPayload))
		if resp.Code != http.StatusOK {
			t.Fatalf("expected import 200, got %d: %s", resp.Code, resp.Body.String())
		}

		result := decodeJSON[ImportNotesResponse](t, resp)
		if result.Imported != 1 {
			t.Fatalf("expected 1 imported note, got %+v", result)
		}

		notes, err := env.store.ListNotes("default")
		if err != nil {
			t.Fatalf("failed to list notes: %v", err)
		}
		if len(notes) != 1 {
			t.Fatalf("expected one note, got %d", len(notes))
		}
		for _, note := range notes {
			if note.FieldMap["Front"] != "Binary\tSearch" {
				t.Fatalf("expected quoted tab to be preserved in YAML import, got %q", note.FieldMap["Front"])
			}
		}
	})
}

func TestAPI_ImportNotes_QuizletAndAnkiText(t *testing.T) {
	t.Run("quizlet quoted tabs", func(t *testing.T) {
		env := setupAPITestEnv(t)

		quizletText := "Term\tDefinition\n\"Binary\tSearch\"\t\"O(log n)\"\n"
		resp := doMultipartImportRequest(t, env.router, map[string]string{
			"source":   "quizlet",
			"deckName": "Quizlet DSA",
		}, "quizlet.tsv", []byte(quizletText))
		if resp.Code != http.StatusOK {
			t.Fatalf("expected import 200, got %d: %s", resp.Code, resp.Body.String())
		}

		result := decodeJSON[ImportNotesResponse](t, resp)
		if result.Imported != 1 {
			t.Fatalf("expected 1 imported note, got %+v", result)
		}

		notes, err := env.store.ListNotes("default")
		if err != nil {
			t.Fatalf("failed to list notes: %v", err)
		}
		for _, note := range notes {
			if note.FieldMap["Front"] != "Binary\tSearch" {
				t.Fatalf("expected embedded tab in quoted quizlet field, got %q", note.FieldMap["Front"])
			}
			if note.FieldMap["Back"] != "O(log n)" {
				t.Fatalf("unexpected back field %q", note.FieldMap["Back"])
			}
		}
	})

	t.Run("anki text", func(t *testing.T) {
		env := setupAPITestEnv(t)

		ankiText := "Front\tBack\tTags\tDeck\tType\nQueue\tFIFO\tlinear ds\tAnki Text Deck\tBasic\n{{c1::Heap}}\tPriority Queue\tgraph\tAnki Text Deck\tCloze\n"
		resp := doMultipartImportRequest(t, env.router, map[string]string{"source": "anki"}, "anki.txt", []byte(ankiText))
		if resp.Code != http.StatusOK {
			t.Fatalf("expected import 200, got %d: %s", resp.Code, resp.Body.String())
		}

		result := decodeJSON[ImportNotesResponse](t, resp)
		if result.Imported != 2 {
			t.Fatalf("expected 2 imported notes, got %+v", result)
		}
		if !containsExact(result.DecksCreated, "Anki Text Deck") {
			t.Fatalf("expected created deck Anki Text Deck, got %+v", result.DecksCreated)
		}

		notes, err := env.store.ListNotes("default")
		if err != nil {
			t.Fatalf("failed to list notes: %v", err)
		}
		if len(notes) != 2 {
			t.Fatalf("expected two notes imported from anki text, got %d", len(notes))
		}

		clozeNote, ok := findNoteByType(t, notes, "Cloze")
		if !ok {
			t.Fatalf("expected cloze note from anki text import")
		}
		if !strings.Contains(clozeNote.FieldMap["Text"], "{{c1::Heap}}") {
			t.Fatalf("expected cloze markers preserved, got %q", clozeNote.FieldMap["Text"])
		}
	})
}

func TestAPI_ImportNotes_AnkiPackage(t *testing.T) {
	env := setupAPITestEnv(t)

	packageBytes := buildAnkiPackage(t)
	resp := doMultipartImportRequest(t, env.router, map[string]string{"source": "anki"}, "dsa.apkg", packageBytes)
	if resp.Code != http.StatusOK {
		t.Fatalf("expected import 200, got %d: %s", resp.Code, resp.Body.String())
	}

	result := decodeJSON[ImportNotesResponse](t, resp)
	if result.Imported != 2 {
		t.Fatalf("expected 2 imported notes from apkg, got %+v", result)
	}
	if !containsExact(result.DecksCreated, "Imported::DSA") {
		t.Fatalf("expected imported deck to be created, got %+v", result.DecksCreated)
	}

	notes, err := env.store.ListNotes("default")
	if err != nil {
		t.Fatalf("failed to list notes: %v", err)
	}
	if len(notes) != 2 {
		t.Fatalf("expected two notes after apkg import, got %d", len(notes))
	}

	if _, ok := findNoteByType(t, notes, "Basic"); !ok {
		t.Fatalf("expected basic note from apkg import")
	}
	if _, ok := findNoteByType(t, notes, "Cloze"); !ok {
		t.Fatalf("expected cloze note from apkg import")
	}
}

func buildAnkiPackage(t *testing.T) []byte {
	t.Helper()

	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "collection.anki2")
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		t.Fatalf("failed to create sqlite db: %v", err)
	}

	createStatements := []string{
		`CREATE TABLE col (id integer primary key, decks text not null, models text not null)`,
		`CREATE TABLE notes (id integer primary key, mid integer not null, tags text not null, flds text not null)`,
		`CREATE TABLE cards (id integer primary key, nid integer not null, did integer not null)`,
	}
	for _, stmt := range createStatements {
		if _, err := db.Exec(stmt); err != nil {
			db.Close()
			t.Fatalf("failed to create test anki schema: %v", err)
		}
	}

	decksJSON := `{"1":{"name":"Default"},"999":{"name":"Imported::DSA"}}`
	modelsJSON := `{
  "100": {"name":"Basic","type":0,"flds":[{"name":"Front"},{"name":"Back"}]},
  "200": {"name":"Cloze","type":1,"flds":[{"name":"Text"},{"name":"Extra"}]}
}`
	if _, err := db.Exec(`INSERT INTO col (id, decks, models) VALUES (1, ?, ?)`, decksJSON, modelsJSON); err != nil {
		db.Close()
		t.Fatalf("failed to seed col table: %v", err)
	}

	if _, err := db.Exec(`INSERT INTO notes (id, mid, tags, flds) VALUES (?, ?, ?, ?)`, 1, 100, " tag1 tag2 ", "Queue\x1fFIFO"); err != nil {
		db.Close()
		t.Fatalf("failed to insert basic note: %v", err)
	}
	if _, err := db.Exec(`INSERT INTO notes (id, mid, tags, flds) VALUES (?, ?, ?, ?)`, 2, 200, " graph ", "{{c1::Heap}}\x1fPriority queue"); err != nil {
		db.Close()
		t.Fatalf("failed to insert cloze note: %v", err)
	}
	if _, err := db.Exec(`INSERT INTO cards (id, nid, did) VALUES (11, 1, 999)`); err != nil {
		db.Close()
		t.Fatalf("failed to insert basic card: %v", err)
	}
	if _, err := db.Exec(`INSERT INTO cards (id, nid, did) VALUES (22, 2, 1)`); err != nil {
		db.Close()
		t.Fatalf("failed to insert cloze card: %v", err)
	}

	if err := db.Close(); err != nil {
		t.Fatalf("failed to close sqlite db: %v", err)
	}

	dbBytes, err := os.ReadFile(dbPath)
	if err != nil {
		t.Fatalf("failed to read sqlite file: %v", err)
	}

	var packageBuf bytes.Buffer
	zipWriter := zip.NewWriter(&packageBuf)
	entry, err := zipWriter.Create("collection.anki2")
	if err != nil {
		t.Fatalf("failed to create zip entry: %v", err)
	}
	if _, err := entry.Write(dbBytes); err != nil {
		t.Fatalf("failed to write zip entry: %v", err)
	}

	mediaEntry, err := zipWriter.Create("media")
	if err != nil {
		t.Fatalf("failed to create media entry: %v", err)
	}
	if err := json.NewEncoder(mediaEntry).Encode(map[string]string{}); err != nil {
		t.Fatalf("failed to write media entry: %v", err)
	}

	if err := zipWriter.Close(); err != nil {
		t.Fatalf("failed to finalize zip: %v", err)
	}
	return packageBuf.Bytes()
}

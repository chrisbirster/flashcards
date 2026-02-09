package main

import (
	"archive/zip"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func writeZipWithCollectionDB(t *testing.T, zipPath string, collectionData []byte) {
	t.Helper()

	f, err := os.Create(zipPath)
	if err != nil {
		t.Fatalf("failed to create zip file: %v", err)
	}
	defer f.Close()

	zw := zip.NewWriter(f)
	entry, err := zw.Create("collection.db")
	if err != nil {
		t.Fatalf("failed to create collection.db entry: %v", err)
	}
	if _, err := entry.Write(collectionData); err != nil {
		t.Fatalf("failed to write collection.db data: %v", err)
	}
	if err := zw.Close(); err != nil {
		t.Fatalf("failed to close zip writer: %v", err)
	}
}

func TestBackupManager_RestoreBackupAndHelpers(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "collection.db")
	backupDir := filepath.Join(tempDir, "backups")
	backupZip := filepath.Join(tempDir, "restore.zip")

	original := []byte("old database bytes")
	restored := []byte("restored database bytes")

	if err := os.WriteFile(dbPath, original, 0644); err != nil {
		t.Fatalf("failed to write original db: %v", err)
	}
	writeZipWithCollectionDB(t, backupZip, restored)

	bm := NewBackupManager(dbPath, backupDir, nil)
	if err := bm.RestoreBackup(backupZip); err != nil {
		t.Fatalf("expected restore backup to succeed, got %v", err)
	}

	got, err := os.ReadFile(dbPath)
	if err != nil {
		t.Fatalf("failed to read restored db: %v", err)
	}
	if string(got) != string(restored) {
		t.Fatalf("expected restored db contents %q, got %q", string(restored), string(got))
	}

	preRestorePath := dbPath + ".pre-restore.backup"
	preRestoreData, err := os.ReadFile(preRestorePath)
	if err != nil {
		t.Fatalf("expected pre-restore backup file to exist: %v", err)
	}
	if string(preRestoreData) != string(original) {
		t.Fatalf("expected pre-restore backup to keep original data %q, got %q", string(original), string(preRestoreData))
	}
}

func TestBackupManager_RestoreBackupErrors(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "collection.db")
	if err := os.WriteFile(dbPath, []byte("db"), 0644); err != nil {
		t.Fatalf("failed to write db fixture: %v", err)
	}

	bm := NewBackupManager(dbPath, filepath.Join(tempDir, "backups"), nil)

	if err := bm.RestoreBackup(filepath.Join(tempDir, "missing.zip")); err == nil {
		t.Fatal("expected missing backup file to return an error")
	}

	badZip := filepath.Join(tempDir, "no-collection-db.zip")
	zf, err := os.Create(badZip)
	if err != nil {
		t.Fatalf("failed to create bad zip: %v", err)
	}
	zw := zip.NewWriter(zf)
	_, _ = zw.Create("not-collection.txt")
	if err := zw.Close(); err != nil {
		t.Fatalf("failed to close bad zip writer: %v", err)
	}
	if err := zf.Close(); err != nil {
		t.Fatalf("failed to close bad zip file: %v", err)
	}

	if err := bm.RestoreBackup(badZip); err == nil {
		t.Fatal("expected restore backup to fail when collection.db is missing")
	}
}

func TestBackupManager_CleanupOldBackups(t *testing.T) {
	tempDir := t.TempDir()
	backupDir := filepath.Join(tempDir, "backups")
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		t.Fatalf("failed to create backup dir: %v", err)
	}

	dbPath := filepath.Join(tempDir, "collection.db")
	if err := os.WriteFile(dbPath, []byte("db"), 0644); err != nil {
		t.Fatalf("failed to write db fixture: %v", err)
	}
	bm := NewBackupManager(dbPath, backupDir, nil)

	files := []string{
		filepath.Join(backupDir, "microdote-backup-20250101-000001.zip"),
		filepath.Join(backupDir, "microdote-backup-20250101-000002.zip"),
		filepath.Join(backupDir, "microdote-backup-20250101-000003.zip"),
		filepath.Join(backupDir, "microdote-backup-20250101-000004.zip"),
	}
	for i, f := range files {
		if err := os.WriteFile(f, []byte("backup"), 0644); err != nil {
			t.Fatalf("failed to write backup fixture %q: %v", f, err)
		}
		mod := time.Now().Add(time.Duration(i) * time.Minute)
		if err := os.Chtimes(f, mod, mod); err != nil {
			t.Fatalf("failed to set backup mtime for %q: %v", f, err)
		}
	}

	if err := bm.CleanupOldBackups(2); err != nil {
		t.Fatalf("expected cleanup to succeed, got %v", err)
	}

	remaining, err := filepath.Glob(filepath.Join(backupDir, "microdote-backup-*.zip"))
	if err != nil {
		t.Fatalf("failed to list remaining backups: %v", err)
	}
	if len(remaining) != 2 {
		t.Fatalf("expected 2 backups after cleanup, got %d (%v)", len(remaining), remaining)
	}
}

func TestBackupManager_CreateBackupErrorsWhenDBMissing(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "missing.db")
	backupDir := filepath.Join(tempDir, "backups")

	bm := NewBackupManager(dbPath, backupDir, nil)
	if _, err := bm.CreateBackup("default"); err == nil {
		t.Fatal("expected create backup to fail when DB file is missing")
	}
}

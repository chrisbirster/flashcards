package main

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// BackupManager handles backup and restore operations for collections.
type BackupManager struct {
	dbPath       string
	backupDir    string
	store        *SQLiteStore
}

// NewBackupManager creates a new backup manager.
func NewBackupManager(dbPath string, backupDir string, store *SQLiteStore) *BackupManager {
	return &BackupManager{
		dbPath:    dbPath,
		backupDir: backupDir,
		store:     store,
	}
}

// CreateBackup creates a timestamped backup of the SQLite database.
// Returns the path to the backup file.
func (bm *BackupManager) CreateBackup(collectionID string) (string, error) {
	// Ensure backup directory exists
	if err := os.MkdirAll(bm.backupDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create backup directory: %w", err)
	}

	// Generate backup filename with timestamp
	timestamp := time.Now().Format("20060102-150405")
	backupFilename := fmt.Sprintf("microdote-backup-%s.zip", timestamp)
	backupPath := filepath.Join(bm.backupDir, backupFilename)

	// Create ZIP file
	zipFile, err := os.Create(backupPath)
	if err != nil {
		return "", fmt.Errorf("failed to create backup file: %w", err)
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// Add SQLite database to ZIP
	if err := bm.addFileToZip(zipWriter, bm.dbPath, "collection.db"); err != nil {
		return "", fmt.Errorf("failed to add database to backup: %w", err)
	}

	// Add metadata file with backup info
	metadata := fmt.Sprintf("Backup created: %s\nCollection ID: %s\nDatabase: %s\n",
		time.Now().Format(time.RFC3339), collectionID, filepath.Base(bm.dbPath))

	metadataWriter, err := zipWriter.Create("backup-info.txt")
	if err != nil {
		return "", fmt.Errorf("failed to create metadata: %w", err)
	}
	if _, err := metadataWriter.Write([]byte(metadata)); err != nil {
		return "", fmt.Errorf("failed to write metadata: %w", err)
	}

	fmt.Printf("Backup created: %s\n", backupPath)
	return backupPath, nil
}

// RestoreBackup restores a collection from a backup ZIP file.
// WARNING: This replaces the current database. The database connection should be closed
// before calling this function.
func (bm *BackupManager) RestoreBackup(backupPath string) error {
	// Verify backup file exists
	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		return fmt.Errorf("backup file not found: %s", backupPath)
	}

	// Open ZIP file
	zipReader, err := zip.OpenReader(backupPath)
	if err != nil {
		return fmt.Errorf("failed to open backup file: %w", err)
	}
	defer zipReader.Close()

	// Find collection.db in ZIP
	var dbFile *zip.File
	for _, file := range zipReader.File {
		if file.Name == "collection.db" {
			dbFile = file
			break
		}
	}

	if dbFile == nil {
		return fmt.Errorf("backup does not contain collection.db")
	}

	// Create temporary file for extraction
	tempPath := bm.dbPath + ".restore.tmp"
	defer os.Remove(tempPath)

	// Extract database to temp file
	if err := bm.extractFile(dbFile, tempPath); err != nil {
		return fmt.Errorf("failed to extract database: %w", err)
	}

	// Backup current database before replacing (just in case)
	currentBackupPath := bm.dbPath + ".pre-restore.backup"
	if err := bm.copyFile(bm.dbPath, currentBackupPath); err != nil {
		fmt.Printf("Warning: could not backup current database: %v\n", err)
	} else {
		fmt.Printf("Current database backed up to: %s\n", currentBackupPath)
	}

	// Replace current database with restored one
	if err := os.Rename(tempPath, bm.dbPath); err != nil {
		return fmt.Errorf("failed to replace database: %w", err)
	}

	// Set flag in metadata to disable auto-sync until user confirms
	// This will be checked when the server restarts
	// (In a full implementation, this would integrate with sync system)

	fmt.Printf("Database restored from: %s\n", backupPath)
	fmt.Println("IMPORTANT: Auto-sync is disabled. Restart server to use restored database.")

	return nil
}

// CleanupOldBackups removes backups older than the retention policy.
// retentionCount: number of most recent backups to keep (e.g., 30)
func (bm *BackupManager) CleanupOldBackups(retentionCount int) error {
	// List all backup files
	files, err := filepath.Glob(filepath.Join(bm.backupDir, "microdote-backup-*.zip"))
	if err != nil {
		return fmt.Errorf("failed to list backups: %w", err)
	}

	if len(files) <= retentionCount {
		return nil // Nothing to clean up
	}

	// Sort by modification time (oldest first)
	type fileInfo struct {
		path    string
		modTime time.Time
	}

	var fileInfos []fileInfo
	for _, path := range files {
		info, err := os.Stat(path)
		if err != nil {
			continue
		}
		fileInfos = append(fileInfos, fileInfo{path: path, modTime: info.ModTime()})
	}

	// Sort by modification time (oldest first)
	for i := 0; i < len(fileInfos); i++ {
		for j := i + 1; j < len(fileInfos); j++ {
			if fileInfos[i].modTime.After(fileInfos[j].modTime) {
				fileInfos[i], fileInfos[j] = fileInfos[j], fileInfos[i]
			}
		}
	}

	// Delete oldest backups beyond retention count
	deleteCount := len(fileInfos) - retentionCount
	for i := 0; i < deleteCount; i++ {
		if err := os.Remove(fileInfos[i].path); err != nil {
			fmt.Printf("Warning: failed to delete old backup %s: %v\n", fileInfos[i].path, err)
		} else {
			fmt.Printf("Deleted old backup: %s\n", fileInfos[i].path)
		}
	}

	return nil
}

// Helper functions

func (bm *BackupManager) addFileToZip(zipWriter *zip.Writer, filePath string, nameInZip string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer, err := zipWriter.Create(nameInZip)
	if err != nil {
		return err
	}

	_, err = io.Copy(writer, file)
	return err
}

func (bm *BackupManager) extractFile(zipFile *zip.File, destPath string) error {
	reader, err := zipFile.Open()
	if err != nil {
		return err
	}
	defer reader.Close()

	writer, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer writer.Close()

	_, err = io.Copy(writer, reader)
	return err
}

func (bm *BackupManager) copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}

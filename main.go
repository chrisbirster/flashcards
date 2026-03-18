package main

import (
	"embed"
	"io/fs"
	"log"
	"net"
	"net/http"
	"os"
)

//go:embed all:web/dist
var embeddedWebDist embed.FS

func main() {
	cfg, err := LoadAppConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	log.Printf("Initializing Vutadex server with %s database mode...", cfg.Database.Mode)
	col, store, err := InitDefaultCollectionWithConfig(cfg.Database)
	if err != nil {
		log.Fatalf("failed to initialize collection: %v", err)
	}
	defer store.Close()

	log.Printf("Collection loaded with %d decks, %d notes, %d cards", len(col.Decks), len(col.Notes), len(col.Cards))

	backupDBPath := ""
	if cfg.Database.Mode == DatabaseModeSQLite {
		backupDBPath = cfg.Database.Path
	}
	backupMgr := NewBackupManager(backupDBPath, "./backups", store)
	handler := NewAPIHandlerWithConfig(store, col, backupMgr, cfg, NewEmailSender(cfg))

	frontendFS, err := fs.Sub(embeddedWebDist, "web/dist")
	if err != nil {
		log.Fatalf("failed to load embedded app assets: %v; build the app with `bun --cwd web run build` first", err)
	}

	server := NewServer(cfg, handler, frontendFS)
	addr := net.JoinHostPort(cfg.Host, cfg.Port)
	appURL := cfg.AppOrigin
	if os.Getenv("PORT") == "" {
		appURL = "http://localhost:" + cfg.Port
	}

	log.Printf("Server starting on %s", addr)
	log.Printf("App available at %s", appURL)
	log.Printf("API available at %s/api", appURL)

	if err := http.ListenAndServe(addr, server); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}

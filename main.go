package main

import (
	"embed"
	"io/fs"
	"log"
	"net"
	"net/http"
)

//go:embed all:web/dist
var embeddedWebDist embed.FS

func main() {
	cfg, err := LoadPlandalfWebConfig()
	if err != nil {
		log.Fatalf("failed to load Plandalf config: %v", err)
	}

	log.Printf("Initializing Plandalf web with %s database mode...", cfg.Database.Mode)
	store, err := OpenPlandalfStore(cfg.Database)
	if err != nil {
		log.Fatalf("failed to initialize Plandalf database: %v", err)
	}
	defer store.Close()

	frontendFS, err := fs.Sub(embeddedWebDist, "web/dist")
	if err != nil {
		log.Fatalf("failed to load embedded app assets: %v; build the app with `bun --cwd web run build` first", err)
	}

	server := NewPlandalfServer(cfg, store, frontendFS)
	addr := net.JoinHostPort(cfg.Host, cfg.Port)
	log.Printf("Plandalf web available at %s", cfg.AppOrigin)
	log.Printf("Plandalf API available at %s/api/v1", cfg.AppOrigin)

	if err := http.ListenAndServe(addr, server); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}

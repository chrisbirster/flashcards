package main

import (
	"io/fs"
	"net/http"
	"path"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func NewServer(cfg AppConfig, handler *APIHandler, frontend fs.FS) http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.RealIP)
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   cfg.AllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization", "X-Vutadex-Plan"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	router.Route("/api", func(r chi.Router) {
		registerAPIRoutes(r, handler)
	})

	spaHandler := NewEmbeddedSPAHandler(frontend)
	router.Handle("/*", spaHandler)
	router.Handle("/", spaHandler)

	return router
}

func NewEmbeddedSPAHandler(frontend fs.FS) http.Handler {
	fileServer := http.FileServer(http.FS(frontend))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			http.NotFound(w, r)
			return
		}

		cleanPath := strings.TrimPrefix(path.Clean("/"+r.URL.Path), "/")
		if cleanPath == "api" || strings.HasPrefix(cleanPath, "api/") {
			http.NotFound(w, r)
			return
		}

		switch cleanPath {
		case ".", "":
			serveEmbeddedIndex(frontend, w, r)
			return
		}

		if stat, err := fs.Stat(frontend, cleanPath); err == nil && !stat.IsDir() {
			fileServer.ServeHTTP(w, r)
			return
		}

		serveEmbeddedIndex(frontend, w, r)
	})
}

func serveEmbeddedIndex(frontend fs.FS, w http.ResponseWriter, r *http.Request) {
	index, err := frontend.Open("index.html")
	if err != nil {
		http.Error(w, "embedded app index is unavailable", http.StatusInternalServerError)
		return
	}
	defer index.Close()

	stat, err := index.Stat()
	if err != nil {
		http.Error(w, "embedded app index metadata is unavailable", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	http.ServeContent(w, r, "index.html", stat.ModTime(), index.(interface {
		Read([]byte) (int, error)
		Seek(int64, int) (int64, error)
	}))
}

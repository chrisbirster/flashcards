package main

import (
	"fmt"
	"os"
	"strings"
)

type PlandalfWebConfig struct {
	Environment    string
	Host           string
	Port           string
	AppOrigin      string
	AllowedOrigins []string
	APIToken       string
	Database       DatabaseConfig
}

func LoadPlandalfWebConfig() (PlandalfWebConfig, error) {
	environment := strings.TrimSpace(os.Getenv("PLANDALF_ENV"))
	if environment == "" { environment = "development" }
	port := strings.TrimSpace(os.Getenv("PORT"))
	if port == "" { port = "8000" }
	host := "localhost"
	if os.Getenv("PORT") != "" { host = "0.0.0.0" }
	appOrigin := strings.TrimSpace(os.Getenv("PLANDALF_APP_ORIGIN"))
	if appOrigin == "" { appOrigin = "http://localhost:" + port }

	database := DatabaseConfig{
		Path: strings.TrimSpace(os.Getenv("PLANDALF_DATABASE_PATH")),
		URL: strings.TrimSpace(os.Getenv("PLANDALF_DATABASE_URL")),
		AuthToken: strings.TrimSpace(os.Getenv("PLANDALF_DATABASE_AUTH_TOKEN")),
	}
	if database.Path == "" { database.Path = "./data/plandalf.db" }
	if database.URL != "" {
		database.Mode = DatabaseModeTurso
		if database.AuthToken == "" {
			return PlandalfWebConfig{}, fmt.Errorf("PLANDALF_DATABASE_AUTH_TOKEN is required when PLANDALF_DATABASE_URL is set")
		}
	} else {
		database.Mode = DatabaseModeSQLite
	}

	allowed := []string{appOrigin, "http://localhost:5173", "http://127.0.0.1:5173"}
	if extra := strings.TrimSpace(os.Getenv("PLANDALF_ALLOWED_ORIGINS")); extra != "" {
		for _, origin := range strings.Split(extra, ",") {
			if origin = strings.TrimSpace(origin); origin != "" { allowed = append(allowed, origin) }
		}
	}
	return PlandalfWebConfig{
		Environment: environment,
		Host: host,
		Port: port,
		AppOrigin: appOrigin,
		AllowedOrigins: allowed,
		APIToken: strings.TrimSpace(os.Getenv("PLANDALF_API_TOKEN")),
		Database: database,
	}, nil
}

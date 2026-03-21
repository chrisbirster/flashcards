package main

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

type DatabaseMode string

const (
	DatabaseModeSQLite DatabaseMode = "sqlite"
	DatabaseModeTurso  DatabaseMode = "turso"
)

type DatabaseConfig struct {
	Mode      DatabaseMode
	URL       string
	AuthToken string
	Path      string
}

type CookieConfig struct {
	Domain string
	Secure bool
}

type OTPConfig struct {
	TTL            time.Duration
	ResendCooldown time.Duration
	MaxAttempts    int
}

type EmailConfig struct {
	SendURL         string
	AuthHeaderName  string
	AuthHeaderValue string
}

type AppConfig struct {
	Environment     string
	Port            string
	Host            string
	AppOrigin       string
	MarketingOrigin string
	AllowedOrigins  []string
	Database        DatabaseConfig
	Cookie          CookieConfig
	OTP             OTPConfig
	SessionTTL      time.Duration
	SessionSecret   string
	Email           EmailConfig
	AuthSuccessPath string
}

func LoadAppConfig() (AppConfig, error) {
	environment := stringEnv("VUTADEX_ENV", "development")
	appOrigin := stringEnv("VUTADEX_APP_ORIGIN", "http://localhost:8000")
	marketingOrigin := stringEnv("VUTADEX_MARKETING_ORIGIN", "http://localhost:4173")
	port := stringEnv("PORT", "8000")
	host := "localhost"
	if os.Getenv("PORT") != "" {
		host = "0.0.0.0"
	}

	database := DatabaseConfig{
		Path:      stringEnv("VUTADEX_DATABASE_PATH", "./data/microdote.db"),
		URL:       strings.TrimSpace(os.Getenv("VUTADEX_DATABASE_URL")),
		AuthToken: strings.TrimSpace(os.Getenv("VUTADEX_DATABASE_AUTH_TOKEN")),
	}
	if database.URL != "" {
		database.Mode = DatabaseModeTurso
	} else {
		database.Mode = DatabaseModeSQLite
	}

	cookieSecureDefault := boolEnvDefault("VUTADEX_COOKIE_SECURE", database.Mode == DatabaseModeTurso || strings.HasPrefix(appOrigin, "https://"))
	cookieDomain := strings.TrimSpace(os.Getenv("VUTADEX_COOKIE_DOMAIN"))
	if cookieDomain == "" && cookieSecureDefault {
		cookieDomain = ".vutadex.com"
	}

	sessionTTLDays := intEnv("VUTADEX_SESSION_TTL_DAYS", 7)
	otpTTLMinutes := intEnv("VUTADEX_OTP_TTL_MINUTES", 10)
	otpResendSeconds := intEnv("VUTADEX_OTP_RESEND_SECONDS", 60)
	otpMaxAttempts := intEnv("VUTADEX_OTP_MAX_ATTEMPTS", 5)
	sessionSecret := strings.TrimSpace(os.Getenv("VUTADEX_SESSION_SECRET"))
	if sessionSecret == "" {
		sessionSecret = "dev-session-secret-change-me"
	}

	cfg := AppConfig{
		Environment:     environment,
		Port:            port,
		Host:            host,
		AppOrigin:       appOrigin,
		MarketingOrigin: marketingOrigin,
		AllowedOrigins:  buildAllowedOrigins(appOrigin, marketingOrigin),
		Database:        database,
		Cookie: CookieConfig{
			Domain: cookieDomain,
			Secure: cookieSecureDefault,
		},
		OTP: OTPConfig{
			TTL:            time.Duration(otpTTLMinutes) * time.Minute,
			ResendCooldown: time.Duration(otpResendSeconds) * time.Second,
			MaxAttempts:    otpMaxAttempts,
		},
		SessionTTL:    time.Duration(sessionTTLDays) * 24 * time.Hour,
		SessionSecret: sessionSecret,
		Email: EmailConfig{
			SendURL:         strings.TrimSpace(os.Getenv("VUTADEX_EMAIL_SEND_URL")),
			AuthHeaderName:  stringEnv("VUTADEX_EMAIL_SEND_AUTH_HEADER", "Authorization"),
			AuthHeaderValue: strings.TrimSpace(os.Getenv("VUTADEX_EMAIL_SEND_AUTH_VALUE")),
		},
		AuthSuccessPath: stringEnv("VUTADEX_AUTH_SUCCESS_URL", "/decks"),
	}

	if cfg.Database.Mode == DatabaseModeTurso && cfg.Database.AuthToken == "" {
		return AppConfig{}, fmt.Errorf("VUTADEX_DATABASE_AUTH_TOKEN is required when VUTADEX_DATABASE_URL is set")
	}

	return cfg, nil
}

func mustLocalAppConfig() AppConfig {
	cfg, err := LoadAppConfig()
	if err == nil {
		return cfg
	}

	return AppConfig{
		Environment:     "development",
		Port:            "8000",
		Host:            "localhost",
		AppOrigin:       "http://localhost:3000",
		MarketingOrigin: "http://localhost:4173",
		AllowedOrigins:  buildAllowedOrigins("http://localhost:3000", "http://localhost:4173"),
		Database: DatabaseConfig{
			Mode: DatabaseModeSQLite,
			Path: "./data/microdote.db",
		},
		Cookie: CookieConfig{
			Secure: false,
		},
		OTP: OTPConfig{
			TTL:            10 * time.Minute,
			ResendCooldown: 60 * time.Second,
			MaxAttempts:    5,
		},
		SessionTTL:    7 * 24 * time.Hour,
		SessionSecret: "dev-session-secret-change-me",
		Email: EmailConfig{
			AuthHeaderName: "Authorization",
		},
		AuthSuccessPath: "/decks",
	}
}

func (cfg AppConfig) IsDevelopment() bool {
	return strings.ToLower(strings.TrimSpace(cfg.Environment)) != "production"
}

func buildAllowedOrigins(appOrigin, marketingOrigin string) []string {
	origins := []string{
		appOrigin,
		marketingOrigin,
		"https://www.vutadex.com",
		"http://localhost:3000",
		"http://127.0.0.1:3000",
		"http://localhost:4173",
		"http://127.0.0.1:4173",
		"http://marketing.lvh.me:4174",
		"http://app.lvh.me:3000",
		"http://localhost:4317",
		"http://127.0.0.1:4317",
	}

	if extra := strings.TrimSpace(os.Getenv("VUTADEX_ALLOWED_ORIGINS")); extra != "" {
		for _, origin := range strings.Split(extra, ",") {
			origins = append(origins, strings.TrimSpace(origin))
		}
	}

	seen := make(map[string]struct{}, len(origins))
	filtered := make([]string, 0, len(origins))
	for _, origin := range origins {
		origin = strings.TrimSpace(origin)
		if origin == "" {
			continue
		}
		parsed, err := url.Parse(origin)
		if err != nil || parsed.Scheme == "" || parsed.Host == "" {
			continue
		}
		if _, ok := seen[origin]; ok {
			continue
		}
		seen[origin] = struct{}{}
		filtered = append(filtered, origin)
	}

	return filtered
}

func stringEnv(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}

func intEnv(key string, fallback int) int {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}

func boolEnvDefault(key string, fallback bool) bool {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return fallback
	}
	return parsed
}

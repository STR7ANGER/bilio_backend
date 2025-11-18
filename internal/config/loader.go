package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	App struct {
		Name string
		Env  string
		Port int
	}
	CORS struct {
		AllowedOrigins []string
	}
	Server struct {
		ReadTimeout     time.Duration
		WriteTimeout    time.Duration
		IdleTimeout     time.Duration
		ShutdownTimeout time.Duration
	}
	Database struct {
		URL string
	}
	Logging struct {
		Level string
	}
	Email struct {
		From string
		SMTP struct {
			Host     string
			Port     int
			Username string
			Password string
		}
	}
}

func Load() (*Config, error) {
	v := viper.New()

	v.SetEnvPrefix("APP")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	v.SetDefault("app.name", "bilio-backend")
	v.SetDefault("app.env", "development")
	v.SetDefault("app.port", 8080)

	v.SetDefault("cors.allowed_origins", []string{"http://localhost:3000"})

	v.SetDefault("server.readtimeout", "5s")
	v.SetDefault("server.writetimeout", "10s")
	v.SetDefault("server.idletimeout", "120s")
	v.SetDefault("server.shutdowntimeout", "15s")

	v.SetDefault("database.url", "postgresql://user:password@localhost:5432/bilio")

	v.SetDefault("logging.level", "info")

	v.SetDefault("email.from", "waitlist@billstack.com")
	v.SetDefault("email.smtp.host", "smtp.gmail.com")
	v.SetDefault("email.smtp.port", 587)
	v.SetDefault("email.smtp.username", "")
	v.SetDefault("email.smtp.password", "")

	bindings := map[string]string{
		"app.env":              "APP_ENV",
		"app.port":             "APP_PORT",
		"cors.allowed_origins": "APP_CORS_ALLOWED_ORIGINS",
		"database.url":         "DATABASE_URL",
		"email.from":           "APP_EMAIL_FROM",
		"email.smtp.username":  "EMAIL_USER",
		"email.smtp.password":  "EMAIL_PASSWORD",
	}
	for key, env := range bindings {
		if err := v.BindEnv(key, env); err != nil {
			return nil, fmt.Errorf("bind env %s: %w", env, err)
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	if envUser := os.Getenv("EMAIL_USER"); envUser != "" {
		cfg.Email.SMTP.Username = envUser
	}
	if envPass := os.Getenv("EMAIL_PASSWORD"); envPass != "" {
		cfg.Email.SMTP.Password = envPass
	}
	if envFrom := os.Getenv("APP_EMAIL_FROM"); envFrom != "" {
		cfg.Email.From = envFrom
	}

	if origins := os.Getenv("APP_CORS_ALLOWED_ORIGINS"); origins != "" {
		cfg.CORS.AllowedOrigins = parseCSV(origins)
	}
	if len(cfg.CORS.AllowedOrigins) == 0 {
		cfg.CORS.AllowedOrigins = []string{"http://localhost:3000"}
	}

	if cfg.Email.From == "" {
		cfg.Email.From = cfg.Email.SMTP.Username
	}

	if cfg.Email.SMTP.Host == "" {
		return nil, fmt.Errorf("email smtp host is required")
	}
	if cfg.Email.SMTP.Port == 0 {
		return nil, fmt.Errorf("email smtp port is required")
	}
	if cfg.Email.SMTP.Username == "" || cfg.Email.SMTP.Password == "" {
		return nil, fmt.Errorf("email smtp credentials are required; set EMAIL_USER and EMAIL_PASSWORD")
	}
	if cfg.Email.From == "" {
		return nil, fmt.Errorf("email from address is required")
	}

	return &cfg, nil
}

func parseCSV(value string) []string {
	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

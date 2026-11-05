package postgres

import (
	"errors"
	"fmt"
	"net"
	"net/url"
	"os"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Config holds connection settings. Use ConfigFromEnv to build one.
type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

func ConfigFromEnv() (Config, error) {
	cfg := Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     envOr("DB_PORT", "5432"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Name:     os.Getenv("DB_NAME"),
		SSLMode:  envOr("DB_SSLMODE", "disable"),
	}

	var missing []string
	for name, val := range map[string]string{
		"DB_HOST":     cfg.Host,
		"DB_USER":     cfg.User,
		"DB_PASSWORD": cfg.Password,
		"DB_NAME":     cfg.Name,
	} {
		if val == "" {
			missing = append(missing, name)
		}
	}
	if len(missing) > 0 {
		return Config{}, errors.New("missing required env vars: " + strings.Join(missing, ", "))
	}

	return cfg, nil
}

// DSN returns a postgres:// URL. Credentials are URL-escaped, so passwords
// containing spaces, quotes or '@' are safe.
func (c Config) DSN() string {
	u := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(c.User, c.Password),
		Host:     net.JoinHostPort(c.Host, c.Port),
		Path:     c.Name,
		RawQuery: url.Values{"sslmode": {c.SSLMode}}.Encode(),
	}
	return u.String()
}

// Connect opens the connection. gorm.Open pings the database, so a nil error
// means PostgreSQL is reachable.
func Connect(cfg Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("connect to postgres at %s: %w", cfg.Host, err)
	}
	return db, nil
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

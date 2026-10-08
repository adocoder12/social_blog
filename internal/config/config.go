package config

import (
	"fmt"
	"net"
	"net/url"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerConfig
	DatabaseConfig
	UploadFileConfig
	CorsConfig
}

type ServerConfig struct {
	SrvPort string
}
type DatabaseConfig struct {
	Host     string
	DBPort   string
	User     string
	Password string
	Name     string
	SSLMODE  string
	MaxConns int32
	MinConns int32
}
type UploadFileConfig struct {
	MaxSizeUpload  int64
	UploadProvider string
}
type CorsConfig struct {
	AllowedOrigin string
}

func LoadConfig() (*Config, error) {
	_ = godotenv.Load()
	maxConns, err := strconv.ParseInt(getEnv("DB_MAX_CONNS", "25"), 10, 32)
	if err != nil {
		return nil, fmt.Errorf("config: invalid DB_MAX_CONNS: %w", err)
	}
	minConns, err := strconv.ParseInt(getEnv("DB_MIN_CONNS", "5"), 10, 32)
	if err != nil {
		return nil, fmt.Errorf("config: invalid DB_MIN_CONNS: %w", err)
	}
	maxUploadSize, err := strconv.ParseInt(getEnv("MAX_UPLOAD_SIZE", "10485760"), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("config: invalid MAX_UPLOAD_SIZE: %w", err)
	}

	config := &Config{
		SrvPort: getEnv("PORT", "8080"),
		//db config
		Host:     getEnv("DB_HOST", "localhost"),
		DBPort:   getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "admin"),
		Password: getEnv("DB_PASSWORD", "password"),
		Name:     getEnv("DB_NAME", "golpher_social"),
		SSLMODE:  getEnv("DB_SSLMODE", "disable"),
		MaxConns: int32(maxConns),
		MinConns: int32(minConns),
		//file config
		MaxSizeUpload:  maxUploadSize,
		UploadProvider: getEnv("UPLOAD_PROVIDER", "local"),
		//Cors
		AllowedOrigin: getEnv("ALLOWED_ORIGIN", "http://localhost:3000"),
	}

	return config, nil
}
func (db *DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		db.Host, db.DBPort, db.User, db.Password, db.Name, db.SSLMODE,
	)
}

// MigrationDSN returns the postgres:// URL format required by golang-migrate.
// Built with net/url so special characters in the password are escaped.
func (db *DatabaseConfig) MigrationDSN() string {
	u := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(db.User, db.Password),
		Host:     net.JoinHostPort(db.Host, db.DBPort),
		Path:     "/" + db.Name,
		RawQuery: "sslmode=" + db.SSLMODE,
	}
	return u.String()
}

// helpers
func (config *Config) Validate() error {
	db := config.DatabaseConfig
	if db.Host == "" {
		return fmt.Errorf("config: DB_HOST is required")
	}
	if db.MaxConns < 1 {
		return fmt.Errorf("config: DB_MAX_CONNS must be at least 1")
	}
	if db.MinConns > db.MaxConns {
		return fmt.Errorf("config: DB_MIN_CONNS must not exceed DB_MAX_CONNS")
	}
	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

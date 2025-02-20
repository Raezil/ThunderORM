package ThunderORM

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// Config holds the configuration for the database connection.
type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// ORM encapsulates a database connection.
type ORM struct {
	DB *sql.DB
}

// NewORM creates a new ORM instance and verifies the database connection.
func NewORM(ctx context.Context, cfg Config) (*ORM, error) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	if err = db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	return &ORM{DB: db}, nil
}

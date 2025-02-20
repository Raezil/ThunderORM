package ThunderORM

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	host = "localhost"
	port = 5432
)

// ORM encapsulates a database connection.
type ORM struct {
	DB *sql.DB
}

// NewORM creates a new ORM instance and verifies the database connection.
func NewORM(ctx context.Context, user string, password string, dbname string) (*ORM, error) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	// Use PingContext to verify connectivity.
	if err = db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	return &ORM{DB: db}, nil
}

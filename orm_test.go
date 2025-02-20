package ThunderORM_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"ThunderORM"
)

type TestUser struct {
	Id   int
	Name string
}

// createTestMigrationFile writes a sample migration SQL file that creates the TestUser table.
func createTestMigrationFile(t *testing.T, dir string) {
	t.Helper()
	migrationSQL := `CREATE TABLE IF NOT EXISTS TestUser (
		"Id" INTEGER PRIMARY KEY,
		"Name" TEXT NOT NULL
	);`
	migrationFile := filepath.Join(dir, "001_create_testuser.sql")
	if err := os.WriteFile(migrationFile, []byte(migrationSQL), 0644); err != nil {
		t.Fatalf("Failed to write migration file: %v", err)
	}
}

func TestNewORM(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cfg := ThunderORM.Config{
		Host:     "localhost",
		Port:     5432,
		User:     "goc",
		Password: "password",
		DBName:   "ormtest",
		SSLMode:  "disable",
	}
	orm, err := ThunderORM.NewORM(ctx, cfg)
	if err != nil {
		t.Fatalf("Failed to create ORM: %v", err)
	}
	defer orm.DB.Close()
}

func TestAutoMigrate(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cfg := ThunderORM.Config{
		Host:     "localhost",
		Port:     5432,
		User:     "goc",
		Password: "password",
		DBName:   "ormtest",
		SSLMode:  "disable",
	}
	orm, err := ThunderORM.NewORM(ctx, cfg)
	if err != nil {
		t.Fatalf("Failed to create ORM: %v", err)
	}
	defer orm.DB.Close()

	tempDir, err := os.MkdirTemp("", "sql")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	createTestMigrationFile(t, tempDir)

	if err := orm.AutoMigrate(ctx, tempDir); err != nil {
		t.Fatalf("AutoMigrate failed: %v", err)
	}
}

func TestCRUDOperations(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cfg := ThunderORM.Config{
		Host:     "localhost",
		Port:     5432,
		User:     "goc",
		Password: "password",
		DBName:   "ormtest",
		SSLMode:  "disable",
	}
	orm, err := ThunderORM.NewORM(ctx, cfg)
	if err != nil {
		t.Fatalf("Failed to create ORM: %v", err)
	}
	defer orm.DB.Close()

	// Ensure the TestUser table exists.
	createTableQuery := `CREATE TABLE IF NOT EXISTS TestUser (
		"Id" INTEGER PRIMARY KEY,
		"Name" TEXT NOT NULL
	);`
	if _, err := orm.DB.ExecContext(ctx, createTableQuery); err != nil {
		t.Fatalf("Failed to create TestUser table: %v", err)
	}
	// Clean up after test.
	defer orm.DB.ExecContext(ctx, `DROP TABLE TestUser;`)

	// --- Test Insertion ---
	newUser := TestUser{
		Id:   1,
		Name: "Alice",
	}
	if err := orm.New(ctx, newUser); err != nil {
		t.Fatalf("Failed to insert new user: %v", err)
	}

	// --- Test Retrieval: All ---
	records, err := orm.All(ctx, TestUser{})
	if err != nil {
		t.Fatalf("Failed to retrieve records: %v", err)
	}
	if len(records) == 0 {
		t.Fatalf("Expected at least one record, got none")
	}

	// --- Test Retrieval: Find ---
	found, err := orm.Find(ctx, TestUser{}, "1")
	if err != nil {
		t.Fatalf("Failed to find record: %v", err)
	}
	if found == nil {
		t.Fatalf("Expected to find record with Id 1, but got nil")
	}
	user, ok := found.(*TestUser)
	if !ok {
		t.Fatalf("Expected *TestUser type, got %T", found)
	}
	if user.Name != "Alice" {
		t.Errorf("Expected Name to be 'Alice', got '%s'", user.Name)
	}

	// --- Test Update ---
	updatedUser := newUser
	updatedUser.Name = "Bob"
	if err := orm.Update(ctx, updatedUser); err != nil {
		t.Fatalf("Failed to update record: %v", err)
	}
	// Verify update.
	found, err = orm.Find(ctx, TestUser{}, "1")
	if err != nil {
		t.Fatalf("Error after update: %v", err)
	}
	if found == nil {
		t.Fatalf("Expected to find updated record, got nil")
	}
	user, ok = found.(*TestUser)
	if !ok {
		t.Fatalf("Expected *TestUser type after update, got %T", found)
	}
	if user.Name != "Bob" {
		t.Errorf("Expected Name to be 'Bob', got '%s'", user.Name)
	}

	// --- Test Deletion: Remove ---
	if err := orm.Remove(ctx, TestUser{}, "1"); err != nil {
		t.Fatalf("Failed to remove record: %v", err)
	}
	// Verify deletion.
	found, err = orm.Find(ctx, TestUser{}, "1")
	if err != nil {
		t.Fatalf("Error after deletion: %v", err)
	}
	if found != nil {
		t.Fatalf("Expected record to be deleted, but found one")
	}
}

// TestWhere validates the functionality of the Where method.
func TestWhere(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cfg := ThunderORM.Config{
		Host:     "localhost",
		Port:     5432,
		User:     "goc",
		Password: "password",
		DBName:   "ormtest",
		SSLMode:  "disable",
	}
	orm, err := ThunderORM.NewORM(ctx, cfg)
	if err != nil {
		t.Fatalf("Failed to create ORM: %v", err)
	}
	defer orm.DB.Close()

	// Ensure the TestUser table exists.
	createTableQuery := `CREATE TABLE IF NOT EXISTS TestUser (
		"Id" INTEGER PRIMARY KEY,
		"Name" TEXT NOT NULL
	);`
	if _, err := orm.DB.ExecContext(ctx, createTableQuery); err != nil {
		t.Fatalf("Failed to create TestUser table: %v", err)
	}
	// Clean up after test.
	defer orm.DB.ExecContext(ctx, `DROP TABLE TestUser;`)

	// Insert sample users.
	users := []TestUser{
		{Id: 1, Name: "Alice"},
		{Id: 2, Name: "Bob"},
		{Id: 3, Name: "Alice"},
	}
	for _, user := range users {
		if err := orm.New(ctx, user); err != nil {
			t.Fatalf("Failed to insert user: %v", err)
		}
	}

	// Retrieve users where Name = 'Alice'
	results, err := orm.Where(ctx, TestUser{}, `"Name" = $1`, "Alice")
	if err != nil {
		t.Fatalf("Where query failed: %v", err)
	}
	if len(results) != 2 {
		t.Errorf("Expected 2 records, got %d", len(results))
	}
}

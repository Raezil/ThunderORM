# ThunderORM

ThunderORM is a lightweight ORM library for Go that simplifies working with PostgreSQL databases. It provides an easy-to-use API for connecting to the database, running SQL migrations, and performing basic CRUD operations using struct metadata and reflection—all while supporting context-aware operations.

## Features

- **Encapsulated Connection:** Manage your database connection through the `ORM` struct.
- **Migrations:** Automatically scan and execute SQL migration files.
- **CRUD Operations:** Simple API for Create, Read, Update and Delete operations.
- **Context-Aware:** All operations accept a `context.Context` to support cancellations and timeouts.
- **Parameterized Queries:** Prevent SQL injection with secure query handling.

## Installation

To install ThunderORM, run:

```bash
go get github.com/Raezil/ThunderORM
```


## Usage Examples

Below are two common use cases that demonstrate how to work with ThunderORM.
### Example 1: Setting Up the ORM and Running Migrations

Create a file (e.g., main_migrate.go) to initialize the ORM and execute migrations stored as SQL files in a directory.

```
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Raezil/ThunderORM"
)

func main() {
	// Create a context with a timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Initialize the ORM with your database credentials.
	orm, err := ThunderORM.NewORM(ctx, "goc", "development", "development")
	if err != nil {
		log.Fatalf("Error initializing ORM: %v", err)
	}

	// Run all migrations from the "./migrations" directory.
	if err := orm.AutoMigrate(ctx, "./migrations"); err != nil {
		log.Fatalf("Error running migrations: %v", err)
	}

	fmt.Println("Migrations executed successfully.")
}
```

This example uses the AutoMigrate method to scan the specified directory for .sql files and execute them against your PostgreSQL database.
### Example 2: Performing CRUD Operations

Create another file (e.g., main_crud.go) to define a sample model and demonstrate basic CRUD operations.
```
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Raezil/ThunderORM"
)

// User represents a record in the "users" table.
type User struct {
	Id    int
	Name  string
	Email string
}

func main() {
	// Create a context with a timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Initialize the ORM.
	orm, err := ThunderORM.NewORM(ctx, "goc", "development", "development")
	if err != nil {
		log.Fatalf("Error initializing ORM: %v", err)
	}

	// Insert a new user.
	newUser := User{
		Id:    1,
		Name:  "Alice",
		Email: "alice@example.com",
	}
	if err := orm.New(ctx, newUser); err != nil {
		log.Fatalf("Error inserting new user: %v", err)
	}
	fmt.Println("Inserted new user.")

	// Retrieve all users.
	users, err := orm.All(ctx, User{})
	if err != nil {
		log.Fatalf("Error retrieving users: %v", err)
	}
	fmt.Println("Users in database:")
	for _, u := range users {
		// Type assertion: All returns pointers to User instances.
		user := u.(*User)
		fmt.Printf("%+v\n", user)
	}

	// Find a specific user by id.
	foundUser, err := orm.Find(ctx, User{}, "1")
	if err != nil {
		log.Fatalf("Error finding user: %v", err)
	}
	if foundUser == nil {
		fmt.Println("User not found.")
	} else {
		user := foundUser.(*User)
		fmt.Printf("Found user: %+v\n", user)
	}

	// Remove a user by id.
	if err := orm.Remove(ctx, User{}, "1"); err != nil {
		log.Fatalf("Error removing user: %v", err)
	}
	fmt.Println("User removed successfully.")
}
```

This example demonstrates how to:
  - Insert a new user record using the New method.
  - Retrieve all users with the All method.
  - Find a user by their ID using the Find method.
  - Delete a user with the Remove method.

## Contributing

Contributions are welcome! Please open an issue or submit a pull request with your suggestions or improvements.

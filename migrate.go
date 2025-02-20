package ThunderORM

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
)

// Data reads the file at the given path and returns its content.
func Data(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("cannot read file %s: %w", path, err)
	}
	return string(content), nil
}

// FindMigrations walks through the root directory to find all files with the given extension.
func FindMigrations(root, ext string) ([]string, error) {
	var migrations []string
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && filepath.Ext(d.Name()) == ext {
			migrations = append(migrations, path)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to walk directory: %w", err)
	}
	// Ensure migrations are sorted for predictable order.
	sort.Strings(migrations)
	return migrations, nil
}

// AutoMigrate executes all migration SQL files found in the specified root directory.
func (o *ORM) AutoMigrate(ctx context.Context, root string) error {
	migrations, err := FindMigrations(root, ".sql")
	if err != nil {
		return err
	}
	for _, migration := range migrations {
		query, err := Data(migration)
		if err != nil {
			return err
		}
		if _, err := o.DB.ExecContext(ctx, query); err != nil {
			return fmt.Errorf("failed executing migration %s: %w", migration, err)
		}
	}
	return nil
}

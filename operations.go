package ThunderORM

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strings"

	"github.com/fatih/structs"
)

// All retrieves all records from the table corresponding to the struct type of u.
func (o *ORM) All(ctx context.Context, u interface{}) ([]interface{}, error) {
	tableName := Name(u)
	query := fmt.Sprintf("SELECT * FROM %s", tableName)
	rows, err := o.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	var results []interface{}
	t := reflect.TypeOf(u)
	for rows.Next() {
		val := reflect.New(t).Interface()
		if err := rows.Scan(Scanning(val)...); err != nil {
			return nil, fmt.Errorf("failed scanning row: %w", err)
		}
		results = append(results, val)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}
	return results, nil
}

// Find retrieves a single record by id using a parameterized query.
func (o *ORM) Find(ctx context.Context, u interface{}, id string) (interface{}, error) {
	tableName := Name(u)
	query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1", tableName)
	row := o.DB.QueryRowContext(ctx, query, id)
	t := reflect.TypeOf(u)
	val := reflect.New(t).Interface()
	if err := row.Scan(Scanning(val)...); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Not found
		}
		return nil, fmt.Errorf("failed scanning row: %w", err)
	}
	return val, nil
}

// New inserts a new record using parameterized queries.
func (o *ORM) New(ctx context.Context, i interface{}) error {
	if !IsStruct(i) {
		return fmt.Errorf("provided value is not a struct")
	}
	tableName := Name(i)
	fieldsSlice := Fields(i)
	valuesSlice := Values(i)

	placeholders := make([]string, len(valuesSlice))
	for idx := range valuesSlice {
		placeholders[idx] = fmt.Sprintf("$%d", idx+1)
	}

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", tableName,
		strings.Join(fieldsSlice, ", "),
		strings.Join(placeholders, ", "))
	_, err := o.DB.ExecContext(ctx, query, valuesSlice...)
	if err != nil {
		return fmt.Errorf("failed to insert record: %w", err)
	}
	return nil
}

// Remove deletes a record by id using a parameterized query.
func (o *ORM) Remove(ctx context.Context, u interface{}, id string) error {
	tableName := Name(u)
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", tableName)
	res, err := o.DB.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete record: %w", err)
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to determine rows affected: %w", err)
	}
	if affected == 0 {
		return fmt.Errorf("no record found with id %s", id)
	}
	return nil
}

// Update updates an existing record using parameterized queries.
// It assumes that the primary key field is "Id" and does not update it.
func (o *ORM) Update(ctx context.Context, i interface{}) error {
	if !IsStruct(i) {
		return fmt.Errorf("provided value is not a struct")
	}
	tableName := Name(i)
	// Use structs.Map to obtain a map of field names to values.
	m := structs.Map(i)
	// Extract the primary key value, assumed to be "Id".
	id, ok := m["Id"]
	if !ok {
		return fmt.Errorf("struct does not have an Id field")
	}
	// Remove the primary key from the update set.
	delete(m, "Id")

	columns := make([]string, 0, len(m))
	values := make([]interface{}, 0, len(m))
	idx := 1
	for col, val := range m {
		columns = append(columns, fmt.Sprintf("%s = $%d", col, idx))
		values = append(values, val)
		idx++
	}
	// Append the primary key value as the last parameter.
	values = append(values, id)
	query := fmt.Sprintf("UPDATE %s SET %s WHERE id = $%d", tableName,
		strings.Join(columns, ", "), idx)
	_, err := o.DB.ExecContext(ctx, query, values...)
	if err != nil {
		return fmt.Errorf("failed to update record: %w", err)
	}
	return nil
}

// Where retrieves records from the table corresponding to the struct type of u
// that satisfy the provided condition.
// The condition should be a valid SQL WHERE clause (without the "WHERE" keyword),
// and args are the corresponding parameter values.
func (o *ORM) Where(ctx context.Context, u interface{}, condition string, args ...interface{}) ([]interface{}, error) {
	tableName := Name(u)
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s", tableName, condition)
	rows, err := o.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	var results []interface{}
	t := reflect.TypeOf(u)
	for rows.Next() {
		val := reflect.New(t).Interface()
		if err := rows.Scan(Scanning(val)...); err != nil {
			return nil, fmt.Errorf("failed scanning row: %w", err)
		}
		results = append(results, val)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}
	return results, nil
}

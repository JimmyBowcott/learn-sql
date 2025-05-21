package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"

	"github.com/lib/pq"
)

type CommandSet struct {
	Clauses   []string `json:"clauses"`
	Modifiers []string `json:"modifiers"`
	Operators []string `json:"operators"`
}

func (c *CommandSet) Scan(value any) error {
	// I want to point out that this will return null
	// if a field is not defined in the database.
	// This is totally not OK but I'm doing it anyway.

	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("expected []byte, got %T", value)
	}
	return json.Unmarshal(bytes, c)
}

func scanRows(rows *sql.Rows) ([]map[string]any, error) {
	res := []map[string]any{}

	for rows.Next() {
		row, err := scanRow(rows)
		if err != nil {
			return nil, err
		}
		res = append(res, row)
	}
	return res, nil
}

func scanRow(rows *sql.Rows) (map[string]any, error) {
	var (
		id          int
		description string
		tables      pq.Int32Array
		commands    CommandSet
	)

	fmt.Println(rows)

	err := rows.Scan(&id, &description, &tables, &commands)
	if err != nil {
		return nil, err
	}

	return map[string]any{"id": id, "description": description, "tables": tables, "commands": commands}, nil
}

func GetLevels() ([]map[string]any, error) {
	connStr := os.Getenv("DB_CONNECTION_STRING_2")
	empty := []map[string]any{}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return empty, err
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return empty, err
	}
	defer tx.Rollback()

	rows, err := tx.Query("SELECT * FROM level;")
	if err != nil {
		return empty, err
	}
	defer rows.Close()

	return scanRows(rows)
}

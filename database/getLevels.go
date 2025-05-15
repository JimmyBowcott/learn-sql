package database

import (
	"database/sql"
	"os"

	"github.com/lib/pq"
)

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
		id int
		description string
		tables pq.Int32Array
	)

	err := rows.Scan(&id, &description, &tables)
	if err != nil {
		return nil, err
	}

	return map[string]any{"id": id, "description": description, "tables": tables}, nil
}

func GetLevels() ([]map[string]any, error) {
	connStr := os.Getenv("DB_CONNECTION_STRING_2")
	empty := []map[string]any{}

	db, err := sql.Open("postgres", connStr)
	if err != nil { return empty, err }
	defer db.Close()

	tx, err := db.Begin()
	if err != nil { return empty, err }
	defer tx.Rollback()

	rows, err := tx.Query("SELECT * FROM level;")
	if err != nil { return empty, err }
	defer rows.Close()

	return scanRows(rows)
}

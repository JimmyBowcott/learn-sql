package database

import (
	"database/sql"
	"os"
	_"github.com/lib/pq"
)

func getCols(rows *sql.Rows, length int) ([]any, error) {
	colsPtr := make([]any, length)
	colsVal := make([]any, length)

	for i := range colsPtr {
		colsPtr[i] = &colsVal[i]
	}

	if err := rows.Scan(colsPtr...); err != nil {
		return []any{}, err
	}

	return colsVal, nil
}

func getRowData(cols []string, vals []any) map[string]any {
	row := make(map[string]any, len(cols))
	for i, col := range cols {
		val := vals[i]
		b, ok := val.([]byte)
		if ok {
			row[col] = string(b)
		} else {
			row[col] = val
		}
	}
	return row
}

func extractDataFromRows(rows *sql.Rows) ([]map[string]any, error) {
	res := []map[string]any{}

	cols, err := rows.Columns()
	if err != nil { return res, err }

	for rows.Next() {
		colVals, err := getCols(rows, len(cols))
		if err != nil { return res, err }
		
		row := getRowData(cols, colVals)
		res = append(res, row)
	}

	return res, nil
}

func ExecuteQuery(query string) ([]map[string]any, error) {
	connStr := os.Getenv("DB_CONNECTION_STRING")
	empty := []map[string]any{}

	db, err := sql.Open("postgres", connStr)
	if err != nil { return empty, err }
	defer db.Close()

	tx, err := db.Begin()
	if err != nil { return empty, err }
	defer tx.Rollback()

	rows, err := tx.Query(query)
	if err != nil { return empty, err }
	defer rows.Close()

	return extractDataFromRows(rows)
}

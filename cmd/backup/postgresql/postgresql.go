package postgresql

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func Client(username, password, host string, port int, dbName *string) (*sql.DB, error) {
	var database string
	if dbName == nil {
		database = "postgres"
	} else {
		database = *dbName
	}

	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host,
		port,
		username,
		password,
		database)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open sql connection: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping postgresql server: %w", err)
	}

	return db, nil
}

func Close(db *sql.DB) error {
	if db == nil {
		return fmt.Errorf("failed to close db connection: given connection handler is nil")
	}

	if err := db.Close(); err != nil {
		return fmt.Errorf("failed to close db connection: %w", err)
	}

	return nil
}

func GetDataDirectory(db *sql.DB) (string, error) {
	const query = "SHOW data_directory"

	rows, err := db.Query(query)
	if err != nil {
		return "", fmt.Errorf("failed to execute \"%s\" query: %w", query, err)
	}
	defer rows.Close()

	var result string

	if !rows.Next() {
		return "", fmt.Errorf("no date returned for \"%s\" query: %w", query, err)
	}

	if err := rows.Scan(&result); err != nil {
		return "", fmt.Errorf("failed to assign result to variable for the \"%s\" query: %w", query, err)
	}

	return result, nil
}

func GetCustomTablespaces(db *sql.DB) (map[string]string, error) {
	const query = "select spcname ,pg_tablespace_location(oid) from pg_tablespace"

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute \"%s\" query: %w", query, err)
	}
	defer rows.Close()

	result := map[string]string{}

	var spcName, location string
	for rows.Next() {
		if err := rows.Scan(&spcName, &location); err != nil {
			return nil, fmt.Errorf("failed to assign result to variable for the \"%s\" query: %w", query, err)
		}

		if location == "" {
			continue
		}
		result[spcName] = location
	}

	return result, nil
}

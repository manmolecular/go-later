package storage

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

const driverName = "sqlite3"

// LocalStorage defines local storage for records
type LocalStorage struct {
	db       *sql.DB
	filename string
}

// NewLocalStorage creates a new local storage object
func NewLocalStorage(filename string) (*LocalStorage, error) {
	if err := createDb(filename); err != nil {
		return nil, fmt.Errorf("can not prepare database: %s\n", err)
	}

	db, err := sql.Open(driverName, filename)
	if err != nil {
		return nil, fmt.Errorf("database connection with file '%s' can not be opened, error: %s\n", filename, err)
	}

	if err = createTable(db); err != nil {
		return nil, fmt.Errorf("table can not be created, error: %s\n", err)
	}

	return &LocalStorage{
		db:       db,
		filename: filename,
	}, nil
}

// CreateRecord creates a record in the storage
func (s *LocalStorage) CreateRecord(content string) error {
	query := `INSERT INTO record(content) VALUES (?)`

	stmt, err := s.db.Prepare(query)
	if err != nil {
		return fmt.Errorf("can not prepare add record statement, error: %s", err)
	}

	if _, err = stmt.Exec(content); err != nil {
		return fmt.Errorf("can not execute prepared add record statement, error: %s", err)
	}

	return nil
}

// GetRecordByID returns record content by its ID
func (s *LocalStorage) GetRecordByID(id int) (string, error) {
	var content string

	query := "SELECT content FROM record WHERE id = ?"

	stmt, err := s.db.Prepare(query)
	if err != nil {
		return content, fmt.Errorf("can not prepare show record statement, error: %s", err)
	}

	if err = stmt.QueryRow(id).Scan(&content); err != nil {
		return content, fmt.Errorf("can not get content for the record, error: %s", err)
	}

	return content, nil
}

// GetAllRecords returns all records
func (s *LocalStorage) GetAllRecords() ([]Record, error) {
	var records []Record

	rows, err := s.db.Query("SELECT * FROM record ORDER BY id DESC")
	if err != nil {
		return records, fmt.Errorf("can not get records, error: %s", err)
	}

	defer func() {
		if err = rows.Close(); err != nil {
			fmt.Printf("can not close rows for further enumeration")
		}
	}()

	for rows.Next() {
		rowRecord := Record{}
		if err = rows.Scan(&rowRecord.ID, &rowRecord.Content); err != nil {
			fmt.Printf("can not get content for the record, error: %s\n", err)
		}

		records = append(records, rowRecord)
	}

	return records, nil
}

// DeleteRecordByID deletes a record from the storage by its ID
func (s *LocalStorage) DeleteRecordByID(id int) error {
	query := `DELETE FROM record WHERE id = ?`

	stmt, err := s.db.Prepare(query)
	if err != nil {
		return fmt.Errorf("can not prepare delete record statement, error: %s", err)
	}

	if _, err = stmt.Exec(id); err != nil {
		return fmt.Errorf("can not execute prepared delete record statement, error: %s", err)
	}

	return nil
}

// DeleteLastRecord deletes the latest record from the storage
func (s *LocalStorage) DeleteLastRecord() error {
	query := `DELETE FROM record WHERE id = (SELECT id FROM record ORDER BY id DESC LIMIT 1)`

	if _, err := s.db.Exec(query); err != nil {
		return fmt.Errorf("can not execute delete last record statement, error: %s", err)
	}

	return nil
}

// Close closes the storage
func (s *LocalStorage) Close() {
	if err := s.db.Close(); err != nil {
		fmt.Printf("database connection with file '%s' can not be closed, error: %s\n", s.filename, err)
	}
}

// createDb creates database file for usage (if not exists yet)
func createDb(filename string) error {
	if _, err := os.Stat(filename); err == nil {
		return nil
	}

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("database file '%s' can not be created, error: %s", filename, err)
	}

	if err = file.Close(); err != nil {
		return fmt.Errorf("database file '%s' can not be closed, error: %s", filename, err)
	}

	return nil
}

// createTable creates table for record entities (if not exists yet)
func createTable(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS record (
		    "id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		    "content" TEXT
		);
	`

	stmt, err := db.Prepare(query)
	if err != nil {
		return fmt.Errorf("can not prepare create table statement, error: %s", err)
	}

	if _, err = stmt.Exec(); err != nil {
		return fmt.Errorf("can not execute prepared create table statement, error: %s", err)
	}

	return nil
}

package storage

import (
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
)

// LocalStorage defines local storage for records
type LocalStorage struct {
	db       *gorm.DB
	filename string
}

// Validate that structure satisfies the interface
var _ Storage = (*LocalStorage)(nil)

// NewLocalStorage creates a new local storage object
func NewLocalStorage(filename string) (*LocalStorage, error) {
	if err := createDb(filename); err != nil {
		return nil, fmt.Errorf("can not prepare database: %s", err)
	}

	db, err := gorm.Open(sqlite.Open(filename), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("database connection with file '%s' can not be opened, error: %s", filename, err)
	}

	if err = createTable(db); err != nil {
		return nil, fmt.Errorf("table can not be created, error: %s", err)
	}

	return &LocalStorage{
		db:       db,
		filename: filename,
	}, nil
}

// CreateRecord creates a record in the storage
func (s *LocalStorage) CreateRecord(content string) error {
	if err := s.db.Create(&Record{Content: content}).Error; err != nil {
		return fmt.Errorf("can not create record, error: %s", err)
	}

	return nil
}

// GetRecordByID returns record content by its ID
func (s *LocalStorage) GetRecordByID(id uint) (string, error) {
	var record Record

	if err := s.db.First(&record, id).Error; err != nil {
		return record.Content, fmt.Errorf("can not get record, error: %s", err)
	}

	return record.Content, nil
}

// GetAllRecords returns all records
func (s *LocalStorage) GetAllRecords() ([]Record, error) {
	var records []Record
	if err := s.db.Order("id DESC").Find(&records).Error; err != nil {
		return records, fmt.Errorf("can not get list of records, error: %s", err)
	}

	return records, nil
}

// DeleteRecordByID deletes a record from the storage by its ID
func (s *LocalStorage) DeleteRecordByID(id uint) error {
	if err := s.db.Delete(&Record{}, id).Error; err != nil {
		return fmt.Errorf("can not delete record, error: %s", err)
	}

	return nil
}

// DeleteLastRecord deletes the latest record from the storage
func (s *LocalStorage) DeleteLastRecord() error {
	query := `DELETE FROM records WHERE id = (SELECT id FROM records ORDER BY id DESC LIMIT 1)`

	db, err := s.db.DB()
	if err != nil {
		fmt.Printf("can not get database object, error: %s\n", err)
	}

	if _, err := db.Exec(query); err != nil {
		return fmt.Errorf("can not execute delete last record statement, error: %s", err)
	}

	return nil
}

// Close closes the storage
func (s *LocalStorage) Close() error {
	db, err := s.db.DB()
	if err != nil {
		return fmt.Errorf("can not get database object, error: %s\n", err)
	}

	if err := db.Close(); err != nil {
		return fmt.Errorf("database connection can not be closed, error: %s\n", err)
	}

	return nil
}

// CleanUp cleans up database
func (s *LocalStorage) CleanUp() error {
	if _, err := os.Stat(s.filename); err != nil {
		return fmt.Errorf("database file '%s' does not exist, error: %s", s.filename, err)
	}

	if err := os.Remove(s.filename); err != nil {
		return fmt.Errorf("database file '%s' can not be deleted, error: %s", s.filename, err)
	}

	return nil
}

// createDb creates database file for usage (if not exists yet)
func createDb(filename string) error {
	if _, err := os.Stat(filename); err == nil {
		return nil // file already exists
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

// createTable creates table for record entities
func createTable(db *gorm.DB) error {
	if err := db.AutoMigrate(&Record{}); err != nil {
		return fmt.Errorf("can not migrate the schema, error: %s", err)
	}

	return nil
}

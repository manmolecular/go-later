package storage

import (
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
	"path"
	"path/filepath"
)

const (
	defaultDbDir  = ".later"
	defaultDbFile = "later.db"
)

// LocalStorage defines local storage for records
type LocalStorage struct {
	db     *gorm.DB
	dbPath string
}

// Validate that structure satisfies the interface
var _ Storage = (*LocalStorage)(nil)

// NewCustomLocalStorage creates a new local storage with custom storage location
func NewCustomLocalStorage(baseDir, dbDir, dbName string) (*LocalStorage, error) {
	dbPath, err := createCustomStorage(baseDir, dbDir, dbName)
	if err != nil {
		return nil, fmt.Errorf("can not create custom storage, error: %s", err)
	}

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("database connection can not be established, error: %s", err)
	}

	if err = createTable(db); err != nil {
		return nil, fmt.Errorf("table can not be created, error: %s", err)
	}

	return &LocalStorage{
		db:     db,
		dbPath: dbPath,
	}, nil
}

// NewLocalStorage creates a new local storage with default configuration in current user's home directory
func NewLocalStorage() (*LocalStorage, error) {
	dbPath, err := createStorage()
	if err != nil {
		return nil, fmt.Errorf("can not create default storage, error: %s", err)
	}

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("database connection can not be established, error: %s", err)
	}

	if err = createTable(db); err != nil {
		return nil, fmt.Errorf("table can not be created, error: %s", err)
	}

	return &LocalStorage{
		db:     db,
		dbPath: dbPath,
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

// GetRecords returns all records
func (s *LocalStorage) GetRecords() ([]Record, error) {
	var records []Record
	if err := s.db.Order("id DESC").Find(&records).Error; err != nil {
		return records, fmt.Errorf("can not get list of records, error: %s", err)
	}

	return records, nil
}

// CountRecords counts all records
func (s *LocalStorage) CountRecords() (uint, error) {
	var count int64
	if err := s.db.Table("records").Count(&count).Error; err != nil {
		return uint(count), fmt.Errorf("can not count records, error: %s", err)
	}

	return uint(count), nil
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
	if _, err := os.Stat(s.dbPath); err != nil {
		return fmt.Errorf("database file does not exist, error: %s", err)
	}

	if err := os.Remove(s.dbPath); err != nil {
		return fmt.Errorf("database file can not be deleted, error: %s", err)
	}

	storageDir := filepath.Dir(s.dbPath)
	if err := os.RemoveAll(storageDir); err != nil {
		return fmt.Errorf("storage directory can not be deleted, error: %s", err)
	}

	return nil
}

// createCustomStorage creates a custom path storage
func createCustomStorage(baseDir, dbDir, dbName string) (string, error) {
	dbDirPath := path.Join(baseDir, dbDir)
	if err := os.MkdirAll(dbDirPath, 0700); err != nil {
		return "", fmt.Errorf("can not create storage directory, error: %s", err)
	}

	dbPath := path.Join(dbDirPath, dbName)
	if err := createDb(dbPath); err != nil {
		return "", fmt.Errorf("can not prepare database, error: %s", err)
	}

	return dbPath, nil
}

// createStorage creates a default storage in current user's home directory
func createStorage() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("can not locate home directory, error: %s", err)
	}

	dbPath, err := createCustomStorage(homeDir, defaultDbDir, defaultDbFile)
	if err != nil {
		return "", fmt.Errorf("can not create storage, error: %s", err)
	}

	return dbPath, nil
}

// createDb creates database file for usage (if not exists yet)
func createDb(dbPath string) error {
	if _, err := os.Stat(dbPath); err == nil {
		return nil // file already exists
	}

	file, err := os.Create(dbPath)
	if err != nil {
		return fmt.Errorf("database file can not be created, error: %s", err)
	}

	if err = file.Close(); err != nil {
		return fmt.Errorf("database file can not be closed, error: %s", err)
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

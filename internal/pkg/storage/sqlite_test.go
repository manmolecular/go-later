package storage

import (
	"os"
	"testing"
)

const (
	testBaseDir = "."
	testDbDir   = "test_db_dir"
	testDbName  = "test_db_name.db"
	testDbPath  = "./test_db_dir/test_db_name.db"
)

// createTestStorage creates test sqlite database using a custom test path
func createTestStorage() (*LocalStorage, error) {
	return NewCustomLocalStorage(testBaseDir, testDbDir, testDbName)
}

// TestNewCustomLocalStorage checks that sqlite database as a storage can be
// created and accessed using a custom path
func TestNewCustomLocalStorage(t *testing.T) {
	s, err := NewCustomLocalStorage(testBaseDir, testDbDir, testDbName)
	if err != nil {
		t.Errorf("custom local storage can not be created, unexpected error: %s", err)
	}

	defer func() {
		if err = s.CleanUp(); err != nil {
			t.Errorf("database file was created, but can not be deleted, unexpected error: %s", err)
		}
	}()

	if _, err := os.Stat(testDbPath); err != nil {
		t.Errorf("database file was not created, unexpected error: %s", err)
	}
}

// TestCreateRecord checks that record can be created and retrieved
func TestCreateRecord(t *testing.T) {
	s, err := createTestStorage()
	if err != nil {
		t.Errorf("test storage can not be created, unexpected error: %s", err)
	}

	defer func() {
		if err = s.CleanUp(); err != nil {
			t.Errorf("database file was created, but can not be deleted, unexpected error: %s", err)
		}
	}()

	testRecordContent := "test_record"

	if err = s.CreateRecord(testRecordContent); err != nil {
		t.Errorf("test record can not be created, unexpected error: %s", err)
	}

	records, err := s.GetAllRecords()
	if err != nil {
		t.Errorf("records can not be retrieved, unexpected error: %s", err)
	}

	if len(records) != 1 {
		t.Errorf("exactly 1 test record expected, got: %d", len(records))
	}

	if records[0].Content != testRecordContent {
		t.Errorf("expected test record content: %s, got: %s", testRecordContent, records[0].Content)
	}
}

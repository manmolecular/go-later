package storage

import "time"

// Record defines record format representation
type Record struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	Content   string
}

// Storage defines common interface for records management
type Storage interface {
	CreateRecord(content string) error
	GetRecordByID(id uint) (string, error)
	GetAllRecords() ([]Record, error)
	DeleteRecordByID(id uint) error
	DeleteLastRecord() error
	Close() error
	CleanUp() error
}

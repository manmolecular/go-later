package storage

// Record defines record format representation
type Record struct {
	ID      int
	Content string
}

// Storage defines common interface for records management
type Storage interface {
	CreateRecord(content string) error
	GetRecordByID(id int) (string, error)
	GetAllRecords() ([]Record, error)
	DeleteRecordByID(id int) error
	DeleteLastRecord() error
	Close()
}

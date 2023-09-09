package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/manmolecular/go-later/internal/pkg/storage"
)

const dbName = "later.db"

// Command implements command handler and router
type Command struct {
	storage storage.Storage
}

// NewCommand creates new Command object
func NewCommand(s storage.Storage) *Command {
	return &Command{storage: s}
}

// handle handles commands passed from the CLI
func (c *Command) handle(args []string) error {
	switch strings.ToLower(args[0]) {
	case "push": // content
		if len(args) < 2 {
			return errors.New("content is not provided")
		}
		record := strings.Join(args[1:], " ")
		if record == "" {
			return errors.New("no content to add")
		}
		if err := c.storage.CreateRecord(record); err != nil {
			return fmt.Errorf("record can not be added to the database, error: %s", err)
		}
	case "pop":
		if err := c.storage.DeleteLastRecord(); err != nil {
			return fmt.Errorf("last record can not be deleted, error: %s", err)
		}
	case "show": // by ID
		if len(args) < 2 {
			return errors.New("ID is not provided")
		}
		id, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("ID has invalid type, error: %s", err)
		}
		content, err := c.storage.GetRecordByID(uint(id))
		if err != nil {
			return fmt.Errorf("record can not be shown, error: %s", err)
		}
		fmt.Println(content)
	case "list":
		records, err := c.storage.GetAllRecords()
		if err != nil {
			return fmt.Errorf("records can not be displayed, error: %s", err)
		}
		for _, rowRecord := range records {
			fmt.Printf("%d. %s (created at: %s)\n", rowRecord.ID, rowRecord.Content, rowRecord.CreatedAt.Format("2006-01-02 15:04:05"))
		}
	case "delete": // by ID
		if len(args) < 2 {
			return errors.New("ID is not provided")
		}
		id, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("ID has invalid type, error: %s", err)
		}
		if err = c.storage.DeleteRecordByID(uint(id)); err != nil {
			return fmt.Errorf("record can not be deleted, error: %s", err)
		}
	case "clean":
		if err := c.storage.CleanUp(); err != nil {
			return fmt.Errorf("storage can not be cleaned up, error: %s", err)
		}
	default:
		return fmt.Errorf("command '%s' is unknown", args[0])
	}

	return nil
}

func main() {
	s, err := storage.NewLocalStorage(dbName)
	if err != nil {
		fmt.Printf("storage can not be accessed or created, error: %s\n", err)
		os.Exit(1)
	}

	defer func() {
		if err := s.Close(); err != nil {
			fmt.Printf("storage can not be closed, error: %s\n", err)
		}
	}()

	flag.Parse()
	args := flag.Args()

	if len(args) == 0 {
		fmt.Printf("later: \n - push <text>\n - pop\n - show <id>\n - list\n - delete <id>\n - clean\n")
		os.Exit(1)
	}

	command := NewCommand(s)
	if err := command.handle(args); err != nil {
		fmt.Printf("command error: %s\n", err)
		os.Exit(1)
	}
}

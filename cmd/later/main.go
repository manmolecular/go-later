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

const (
	cmdPush   = "push"
	cmdPop    = "pop"
	cmdShow   = "show"
	cmdList   = "list"
	cmdCount  = "count"
	cmdDelete = "delete"
	cmdClean  = "clean"
)

var cmdToDesc = map[string]string{
	cmdPush:   "add new task",
	cmdPop:    "delete the latest task",
	cmdShow:   "show the exact task by its ID",
	cmdList:   "list all tasks",
	cmdCount:  "count tasks",
	cmdDelete: "delete the exact task by its ID",
	cmdClean:  "clean the database",
}

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
	case cmdPush: // content
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
	case cmdPop:
		if err := c.storage.DeleteLastRecord(); err != nil {
			return fmt.Errorf("last record can not be deleted, error: %s", err)
		}
	case cmdShow: // by ID
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
	case cmdList:
		records, err := c.storage.GetRecords()
		if err != nil {
			return fmt.Errorf("records can not be displayed, error: %s", err)
		}
		for _, rowRecord := range records {
			fmt.Printf("%d. %s (created at: %s)\n", rowRecord.ID, rowRecord.Content, rowRecord.CreatedAt.Format("2006-01-02 15:04:05"))
		}
	case cmdCount:
		count, err := c.storage.CountRecords()
		if err != nil {
			return fmt.Errorf("records can not be counted, error: %s", err)
		}
		fmt.Println(count)
	case cmdDelete: // by ID
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
	case cmdClean:
		if err := c.storage.CleanUp(); err != nil {
			return fmt.Errorf("storage can not be cleaned up, error: %s", err)
		}
	default:
		return fmt.Errorf("command '%s' is unknown", args[0])
	}

	return nil
}

func main() {
	s, err := storage.NewLocalStorage()
	if err != nil {
		fmt.Printf("storage can not be accessed or created, error: %s\n", err)
		os.Exit(1)
	}

	defer func() {
		if err := s.Close(); err != nil {
			fmt.Printf("storage can not be closed, error: %s\n", err)
		}
	}()

	flag.Usage = func() {
		for cmd, desc := range cmdToDesc {
			fmt.Printf("- %s: %s\n", cmd, desc)
		}
	}

	flag.Parse()
	args := flag.Args()

	if len(args) == 0 {
		fmt.Println("no subcommands provided, list of supported subcommands:")
		flag.Usage()
		os.Exit(1)
	}

	command := NewCommand(s)
	if err := command.handle(args); err != nil {
		fmt.Printf("command error: %s\n", err)
		flag.Usage()
		os.Exit(1)
	}
}

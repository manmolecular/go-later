package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/manmolecular/go-later/internal/pkg/storage"
)

const dbName = "later.db"

func main() {
	db, err := storage.NewLocalStorage(dbName)
	if err != nil {
		fmt.Printf("storage can not be accessed or created, error: %s\n", err)
	}
	defer db.Close()

	flag.Parse()
	args := flag.Args()

	if len(args) == 0 {
		fmt.Printf("one of: push <text>, pop, show <id>, list, delete <id>\n")
		os.Exit(1)
	}

	switch strings.ToLower(args[0]) {
	case "push":
		record := strings.Join(args[1:], " ")
		if err := db.CreateRecord(record); err != nil {
			fmt.Printf("record can not be added to the database, error: %s\n", err)
			os.Exit(1)
		}
	case "pop":
		if err := db.DeleteLastRecord(); err != nil {
			fmt.Printf("last record can not be deleted, error: %s\n", err)
			os.Exit(1)
		}
	case "show":
		id, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Printf("ID has invalid type, error: %s\n", err)
			os.Exit(1)
		}
		content, err := db.GetRecordByID(id)
		if err != nil {
			fmt.Printf("record can not be shown, error: %s\n", err)
			os.Exit(1)
		}
		fmt.Println(content)
	case "list":
		records, err := db.GetAllRecords()
		if err != nil {
			fmt.Printf("records can not be displayed, error: %s\n", err)
			os.Exit(1)
		}
		for _, rowRecord := range records {
			fmt.Printf("%d. %s\n", rowRecord.ID, rowRecord.Content)
		}
	case "delete":
		id, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Printf("ID has invalid type, error: %s\n", err)
			os.Exit(1)
		}
		if err = db.DeleteRecordByID(id); err != nil {
			fmt.Printf("record can not be deleted, error: %s\n", err)
			os.Exit(1)
		}
	}
}

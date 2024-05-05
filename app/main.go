package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
)

// Usage: your_sqlite3.sh sample.db .dbinfo
func main() {
	databaseFilePath := os.Args[1]
	command := os.Args[2]

	switch command {
	case ".dbinfo":
		databaseFile, err := os.Open(databaseFilePath)
		if err != nil {
			log.Fatal(err)
		}
		defer databaseFile.Close()

		header, err := ReadDatafileHeader(databaseFile)
		if err != nil {
			log.Fatal(err)
		}

		leafTableHeader, payloads := ReadRootPage(databaseFile, header.pageSize)

		var numTables uint32
		tableSignature := []byte("table")
		for _, payload := range payloads {
			if bytes.Contains(payload, tableSignature) {
				numTables++
			}
		}

		fmt.Printf("database page size: %v\n", header.pageSize)
		fmt.Println("number of cells in the root page: ", leafTableHeader.NumCells)
		fmt.Printf("number of tables: %v\n", numTables)
	default:
		fmt.Println("Unknown command", command)
		os.Exit(1)
	}
}

package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	// Available if you need it!
	// "github.com/xwb1989/sqlparser"
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

		header := make([]byte, 100)

		_, err = databaseFile.Read(header)
		if err != nil {
			log.Fatal(err)
		}

		var headerPageSize uint16
		if err := binary.Read(bytes.NewReader(header[16:18]), binary.BigEndian, &headerPageSize); err != nil {
			fmt.Println("Failed to read integer:", err)
			return
		}
		var pageSize uint32
		if headerPageSize == 1 {
			pageSize = 65536
		} else {
			pageSize = uint32(headerPageSize)
		}
		// read the first page
		rootPage := make([]byte, pageSize)

		_, err = databaseFile.Read(rootPage)
		if err != nil {
			log.Fatal(err)
		}

		var numTables uint16
		if err := binary.Read(bytes.NewReader(rootPage[3:5]), binary.BigEndian, &numTables); err != nil {
			fmt.Println("Failed to read integer:", err)
			return
		}

		fmt.Printf("database page size: %v", pageSize)
		fmt.Printf("number of tables: %v", numTables)
	default:
		fmt.Println("Unknown command", command)
		os.Exit(1)
	}
}

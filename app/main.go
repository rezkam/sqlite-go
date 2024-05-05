package main

import (
	"bytes"
	"encoding/binary"
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

		// read the first page
		rootPage := make([]byte, header.pageSize)

		_, err = databaseFile.ReadAt(rootPage, 0)
		if err != nil {
			log.Fatal(err)
		}

		// For the first page, offset calculations must consider the database header
		pageOffset := 100 // Database header offset for the first page

		var pageType = rootPage[pageOffset]
		if pageType != 0x0D {
			log.Fatal("Root page is not a table b-tree page")
		}
		// Leaf Table B-Tree Page
		leafTableHeader := LeafTableHeader{
			FirstFreeblock:              binary.BigEndian.Uint16(rootPage[1+pageOffset : 3+pageOffset]),
			NumCells:                    binary.BigEndian.Uint16(rootPage[3+pageOffset : 5+pageOffset]),
			StartOfCellContentArea:      binary.BigEndian.Uint16(rootPage[5+pageOffset : 7+pageOffset]),
			NumberofFragmentedFreeBytes: rootPage[7+pageOffset],
		}

		cellPointerArray := make([]uint16, leafTableHeader.NumCells)
		for i := 0; i < int(leafTableHeader.NumCells); i++ {
			cellPointerArray[i] = binary.BigEndian.Uint16(rootPage[8+pageOffset+i*2 : 10+pageOffset+i*2])
		}

		var payloads [][]byte
		for _, pointer := range cellPointerArray {
			ptrOffset := int(pointer)
			payloadSize, payloadSizeLen := binary.Uvarint(rootPage[ptrOffset:])
			_, rowidLen := binary.Uvarint(rootPage[ptrOffset+payloadSizeLen:])

			payloads = append(payloads, rootPage[ptrOffset+payloadSizeLen+rowidLen:ptrOffset+int(payloadSize)])

		}

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

type LeafTableHeader struct {
	FirstFreeblock              uint16
	NumCells                    uint16
	StartOfCellContentArea      uint16
	NumberofFragmentedFreeBytes uint8
}

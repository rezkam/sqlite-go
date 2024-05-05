package main

import (
	"encoding/binary"
	"io"
	"log"
)

func ReadRootPage(databaseFile io.ReaderAt, pageSize uint32) (LeafTableHeader, [][]byte) {
	rootPage := make([]byte, pageSize)

	_, err := databaseFile.ReadAt(rootPage, 0)
	if err != nil {
		log.Fatal(err)
	}

	pageOffset := 100

	var pageType = rootPage[pageOffset]
	if pageType != 0x0D {
		log.Fatal("Root page is not a table b-tree page")
	}

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
	return leafTableHeader, payloads
}

type LeafTableHeader struct {
	FirstFreeblock              uint16
	NumCells                    uint16
	StartOfCellContentArea      uint16
	NumberofFragmentedFreeBytes uint8
}

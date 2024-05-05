package main

import (
	"bytes"
	"encoding/binary"
	"io"
)

type DatafileHeader struct {
	pageSize uint32
}

func ReadDatafileHeader(r io.Reader) (*DatafileHeader, error) {
	headerData := make([]byte, 100)

	_, err := r.Read(headerData)
	if err != nil {
		return nil, err
	}

	var headerPageSize uint16
	if err := binary.Read(bytes.NewReader(headerData[16:18]), binary.BigEndian, &headerPageSize); err != nil {
		return nil, err
	}
	var pageSize uint32
	if headerPageSize == 1 {
		pageSize = 65536
	} else {
		pageSize = uint32(headerPageSize)
	}

	return &DatafileHeader{
		pageSize: pageSize,
	}, nil
}

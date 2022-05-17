package client

import (
	"fmt"
	"hash"
	"hash/crc64"
	"io"
	"os"
	"path/filepath"
)

func NewHash() hash.Hash {
	return crc64.New(crc64.MakeTable(crc64.ISO))
}

func DoHash(ssc *SscClient, args *Arguments) error {
	if len(args.Directory) == 0 {
		return fmt.Errorf("no source path specified")
	}
	filename := ValueOrDefault(args.FileName, TEST_SOURCE_FILE)

	_, err := processHash(filepath.Join(args.Directory, filename))
	return err
}

func processHash(filePath string) ([]byte, error) {
	rIn, err := os.Open(filePath)
	if err != nil {
		return nil,  fmt.Errorf("Could not open %s to compute checksum \n%v\n", filePath, err)
	}
	defer rIn.Close()

	buff := make([]byte, 1024*100)
	hasher := NewHash()

	for {
		size, err := rIn.Read(buff)
		if size > 0 {
			_, crcErr := hasher.Write(buff[:size])
			if crcErr != nil {
				return nil, fmt.Errorf("failed to create checksum on test file (%s) %v\n", filePath, crcErr)
			}

		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read test file (%s) %v\n", filePath, err)
		}
	}

	sum := hasher.Sum(nil)
	fmt.Printf("Successfully created checksum %x on file %s\n", sum, filePath)

	return sum, nil
}
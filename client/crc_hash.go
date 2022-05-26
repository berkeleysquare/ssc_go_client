package client

import (
	"encoding/csv"
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

func PrintChecksumCsvHeader(csv *csv.Writer) error {
	var line = []string {"Key","Checksum"}
	return csv.Write(line)
}

func PrintChecksumList(csv *csv.Writer, key string, sum []byte) error {
	var line = []string {key, fmt.Sprintf("%x", sum)}
	return csv.Write(line)
}

func DoHash(ssc *SscClient, args *Arguments) error {
	if len(args.Directory) == 0 {
		return fmt.Errorf("no source path specified")
	}
	filename := ValueOrDefault(args.FileName, TEST_SOURCE_FILE)

	filePath :=filepath.Join(args.Directory, filename)
	sum, err := processHash(filePath)
	if err == nil {
		fmt.Printf("Successfully created checksum %x on file %s\n", sum, filePath)
	}
	return err
}

func HashDirectory(ssc *SscClient, args *Arguments) error {
	if len(args.Directory) == 0 {
		return fmt.Errorf("no source path specified")
	}
	outputFile := args.OutputFile
	wOut := os.Stdout
	if len(outputFile) > 0 {
		f, err := os.Create(outputFile)
		if err != nil {
			return fmt.Errorf("Could not create %s\n%v\n", outputFile, err)
		}
		defer f.Close()
		wOut = f
	}
	w := csv.NewWriter(wOut)
	defer w.Flush()
	err := PrintChecksumCsvHeader(w)
	if err != nil {
		return fmt.Errorf("Could not write to CSV\n%v\n", err)
	}

	limit := args.Count
	count := 0

	err = filepath.Walk(args.Directory,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				sum, err := processHash(path)
				if err != nil {
					return fmt.Errorf("could not generate checksum for %s, %v", path, err)
				}
				err = PrintChecksumList(w, path, sum)
				if err != nil {
					return fmt.Errorf("Failed writing checksum\n%v\n", err)
				}
				if count > limit {
					return io.EOF
				}
				count++
			}
			return nil
		})
	if err != nil && err != io.EOF {
		return fmt.Errorf("Could not walk filepath %s\n%v", args.Directory, err)
	}
	return nil
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
	return hasher.Sum(nil), nil
}

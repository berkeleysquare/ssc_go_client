package client

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/SpectraLogic/ssc_go_client/openapi"
	"html/template"
	"log"
	"os"
)

func writeBreadcrumbs(ssc *SscClient, args *Arguments) error {

	// fileName runs once, inputFile iterates through CSV
	if len(args.Job) == 0 {
		return fmt.Errorf("no job name specified")
	}

	verbose := args.Verbose

	var w *csv.Writer
	outputFile := args.OutputFile

	// list files
	wOut := os.Stdout
	if len(outputFile) > 0 {
		f, err := os.Create(outputFile)
		if err != nil {
			return fmt.Errorf("Could not create %s\n%v\n", outputFile, err)
		}
		defer f.Close()
		wOut = f
	}
	w = csv.NewWriter(wOut)
	defer w.Flush()

	err := PrintManifestCsvHeader(w)
	if err != nil {
		return fmt.Errorf("could not print manifest header %v", err)
	}

	/**
		wOut := os.Stdout
		w := bufio.NewWriter(wOut)
		defer w.Flush()
	**/

	tmpl := template.Must(template.ParseFiles(args.InputFile))

	// iterate through all pages and extract just names
	offset := int64(0)
	more := true
	if verbose {
		log.Printf("doGetManifest(%s)", args.Job)
	}
	for more {
		ret, isTruncated, err := doGetBogusManifest(ssc, args.Job, offset, verbose)
		if err != nil {
			return fmt.Errorf("get manifest %s failed %v\n", args.Job, err)
		}

		err = doBreadcrumbs(tmpl, ret)
		if err != nil {
			return fmt.Errorf("could not list search results %v\n", err)
		}
		offset += limit
		more = isTruncated
	}
	fmt.Printf("Successfully ran Command\n")
	return nil
}

func doGetBogusManifest(ssc *SscClient, job string, offset int64, verbose bool) ([]openapi.ApiManifestFile, bool, error) {

	bogus := `[
  {
    "path": "\\\\localhost\\c$\\StorCycle\\shares\\delayedAction\\coffeehouse\\jk\\IHaveDreamed.mp3",
    "name": "IHaveDreamed.mp3",
    "size": 1525449,
    "recordTime": "2023-08-25T20:56:18.143Z",
    "hash": "b2b256.9d7fa6ebee83b1ee72d7a9da8f4f02f35887f4c26d5d3992365a95e2c9982588"
  },
  {
    "path": "\\\\localhost\\c$\\StorCycle\\shares\\delayedAction\\FarAwayPlaces.jpg",
    "name": "FarAwayPlaces.jpg",
    "size": 62061,
    "recordTime": "2023-08-25T20:56:18.143Z",
    "hash": "b2b256.7b79177df8f840b2ebb06599af9fbbb6aef2d4431d63ec86cf187e1a3149c9d7"
  },
  {
    "path": "\\\\localhost\\c$\\StorCycle\\shares\\delayedAction\\coffeehouse\\IsYouIsOrIsYouAint.mp3",
    "name": "IsYouIsOrIsYouAint.mp3",
    "size": 2310657,
    "recordTime": "2023-08-25T20:56:18.144Z",
    "hash": "b2b256.bb0ebc9ade6d905decdaec5a4f03392b6e3d8ff5ca2220208cbeae6a0b2cd24b"
  }
]`
	isTruncated := false

	var files []openapi.ApiManifestFile

	err := json.Unmarshal([]byte(bogus), &files)
	if err != nil {
		return nil, isTruncated, err
	}

	if verbose {
		log.Printf("Manifest %s: TotalResults %d, offset %d, limit %d", job, len(files), offset, limit)
	}
	return files, isTruncated, nil
}

func doBreadcrumbs(tmpl *template.Template, files []openapi.ApiManifestFile) error {

	for fileIndex := range files {
		file := files[fileIndex]
		var f *os.File
		fpath := *file.Path + ".html"
		f, err := os.Create(fpath)
		if err != nil {
			return fmt.Errorf("Failed to create file %s\n%v", fpath, err)
		}
		tmpl.Execute(f, file)
		err = f.Close()
		if err != nil {
			return fmt.Errorf("Failed to close file %s\n%v", fpath, err)
		}
	}
	return nil
}

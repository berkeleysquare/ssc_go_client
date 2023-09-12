package client

import (
	"encoding/csv"
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

	tmpl := template.Must(template.ParseFiles(args.InputFile))

	// iterate through all pages and extract just names
	offset := int64(0)
	more := true
	if verbose {
		log.Printf("doGetManifest(%s)", args.Job)
	}
	for more {
		ret, isTruncated, err := doGetManifest(ssc, args.Job, offset, verbose)
		if err != nil {
			return fmt.Errorf("get manifest %s failed %v\n", args.Job, err)
		}

		err = doBreadcrumbs(tmpl, ret, verbose)
		if err != nil {
			return fmt.Errorf("could not list search results %v\n", err)
		}
		offset += limit
		more = isTruncated
	}
	fmt.Printf("Successfully ran Command\n")
	return nil
}

func doBreadcrumbs(tmpl *template.Template, files []openapi.ApiManifestFile, verbose bool) error {

	for fileIndex := range files {
		file := files[fileIndex]
		var f *os.File
		fpath := *file.Path + ".html"
		if verbose {
			log.Printf("makeBreadcrumb %s\n", fpath)
		}
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

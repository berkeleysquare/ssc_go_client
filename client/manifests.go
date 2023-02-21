package client

import (
	"encoding/csv"
	"fmt"
	"github.com/SpectraLogic/ssc_go_client/openapi"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	limit = int64(5000)
)

func PrintManifestCsvHeader(w *csv.Writer) error {
	var line = []string {"Name","Path","Size","Checksum"}
	return w.Write(line)
}

func getManifestFiles(ssc *SscClient, job string, fileName string, exts []string, verbose bool) ([]string, error) {
	ret := []string{}
	// search for all files including case number
	if len(job) == 0 {
		return nil, fmt.Errorf("no job name specified" )
	}
	if verbose {
		log.Printf("getManifestFiles(%s, %s)", job, fileName)
	}
	// iterate through all pages and extract just names
	offset := int64(0)
	more := true
	for more {
		if verbose {
			log.Printf("doGetManifest(%s, %d)", job, offset)
		}
		files, isTruncated, err := doGetManifest(ssc, job, offset, verbose)
		if err != nil {
			return nil, fmt.Errorf("get catalog for %s failed %v\n", job, err)
		}
		filtered := filterFiles(files, fileName, exts, verbose)
		for fileIndex := range filtered {
			file := files[fileIndex]
			ret = append(ret, *file.Path)
		}
		if verbose {
			log.Printf("%d matching files", len(filtered))
		}
		more = isTruncated
		offset += limit
	}
	return ret, nil
}

func doGetManifest(ssc *SscClient, job string, offset int64, verbose bool) ([]openapi.ApiManifestFile, bool, error) {
	// search for all files including case number
	if len(job) == 0 {
		return nil, false, fmt.Errorf("no job name specified" )
	}

	opts := openapi.ProjectApiListCatalogOpts{Limit: limit, Skip: offset}
	response, _, err := ssc.Client.ProjectApi.GetCatalog(*ssc.Context, job, &opts)
	if err != nil {
		return nil, false, fmt.Errorf("get catalog for %s failed %v\n", job, err)
	}
	isTruncated := response.TotalResults > offset + limit
	if verbose {
		log.Printf("Manifest %s: TotalResults %d, offset %d, limit %d", job, response.TotalResults, offset, limit)
	}
	return response.Data, isTruncated, nil
}

func filterFiles(files []openapi.ApiManifestFile, fileName string, exts []string, verbose bool) []openapi.ApiManifestFile {
	if len(fileName) == 0 && len(exts) == 0 {
		// no filters
		return files
	}

	matching := []openapi.ApiManifestFile{}
	if files != nil {
		for objectIndex := range files {
			objectName := files[objectIndex]
			if len(fileName) == 0 || strings.Contains(*objectName.Name, fileName) {
				for extIndex := range exts {
					if strings.HasSuffix(*objectName.Name, exts[extIndex]) {
						if verbose {
							log.Printf("Matched %s: job %s, ext: %s", *objectName.Name, fileName, exts[extIndex])
						}
						matching = append(matching, objectName)
					}
				}
			}
		}
	}
	return matching
}

func listCatalog(ssc *SscClient, args *Arguments) error {

	// fileName runs once, inputFile iterates through CSV
	if len(args.Job) == 0 {
		return fmt.Errorf("no job name specified" )
	}

	verbose := args.Verbose

	// output -- console, csv file, or none
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

		err = displayManifestFileObjects(w, ret)
		if err != nil {
			return fmt.Errorf("could not list search results %v\n", err)
		}
		offset += limit
		more = isTruncated
	}
	fmt.Printf("Successfully ran Command\n",)
	return nil
}

func displayManifestFileObjects(w *csv.Writer, files []openapi.ApiManifestFile) error {

	lines := [][]string{}
	for fileIndex := range files {
		file := files[fileIndex]
		path := ""
		if file.Path != nil {
			path = *file.Path
		}
		size := "0"
		if file.Size != nil {
			size = strconv.Itoa(int(*file.Size))
		}
		hash := ""
		if file.Checksum != nil {
			hash = *file.Checksum
		}
		lines = append(lines, []string {*file.Name, path, size, hash})
	}
	return w.WriteAll(lines)
}

package client

import (
	"fmt"
	"github.com/SpectraLogic/ssc_go_client/openapi"
	"html/template"
	"log"
	"os"
	"path/filepath"
)

type breadcrumbInfo struct {
	Name       *string `json:"name" xml:"name"`
	Job        *string `json:"job" xml:"job"`
	Path       *string `json:"path,omitempty" xml:"path"`
	Size       *int64  `json:"size,omitempty" xml:"size"`
	RecordTime *string `json:"recordTime,omitempty" xml:"recordTime"`
}

func makeBreadcrumb(job string, file *openapi.ApiManifestFile) *breadcrumbInfo {
	return &breadcrumbInfo{
		Name:       file.Name,
		Job:        &job,
		Path:       file.Path,
		Size:       file.Size,
		RecordTime: file.RecordTime,
	}
}

func writeBreadcrumbs(ssc *SscClient, args *Arguments) error {
	// required params
	if len(args.Job) == 0 {
		return fmt.Errorf("no job name (--job) specified")
	}
	if len(args.InputFile) == 0 {
		return fmt.Errorf("no HTML template (--in) specified")
	}
	verbose := args.Verbose

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

		err = doBreadcrumbs(tmpl, ret, args.Job, verbose)
		if err != nil {
			return fmt.Errorf("could not list search results %v\n", err)
		}
		offset += limit
		more = isTruncated
	}
	fmt.Printf("\nSuccessfully ran Command\n")
	return nil
}

func doBreadcrumbs(tmpl *template.Template, files []openapi.ApiManifestFile, job string, verbose bool) error {

	currentContainingDirectory := ""
	for fileIndex := range files {
		file := files[fileIndex]
		var f *os.File
		fullPath := *file.Path + ".html"
		// ensure directory exists on change
		directory := filepath.Dir(fullPath)
		if currentContainingDirectory != directory {
			if verbose {
				log.Printf("Create directories: %s", directory)
			}
			err := os.MkdirAll(directory, 0777)
			if err != nil {
				return fmt.Errorf("Failed to create directory %s\n%v", directory, err)
			}
			currentContainingDirectory = directory
		}

		info := makeBreadcrumb(job, &file)
		if verbose {
			log.Printf("makeBreadcrumb %s\n", fullPath)
		}
		f, err := os.Create(fullPath)
		if err != nil {
			return fmt.Errorf("Failed to create file %s\n%v", fullPath, err)
		}
		tmpl.Execute(f, info)
		err = f.Close()
		if err != nil {
			return fmt.Errorf("Failed to close file %s\n%v", fullPath, err)
		}
	}
	return nil
}

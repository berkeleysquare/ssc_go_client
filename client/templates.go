package client

import (
	"errors"
	"fmt"
	"github.com/SpectraLogic/ssc_go_client/openapi"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"sort"
)

const (
	mailConfigFile   = "mail_config.yaml"
	updateTokenEvery = 1_000_000
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

type HasWarningsError struct {
	message string
}

func (e *HasWarningsError) Error() string {
	return e.message
}

func writeBreadcrumbs(ssc *SscClient, args *Arguments) error {

	// required params
	if len(args.Job) == 0 {
		return fmt.Errorf("no job name (--job) specified")
	}
	if len(args.InputFile) == 0 {
		return fmt.Errorf("no HTML template (--in) specified")
	}

	err := breadcrumbsForOneProject(ssc, args.Job, args)
	if err != nil {
		return fmt.Errorf("could not write breadcrumbs %v\n", err)
	}
	fmt.Printf("\nSuccessfully ran Command\n")
	return nil
}

func breadcrumbsForOneProject(ssc *SscClient, job string, args *Arguments) error {
	// required params
	if len(job) == 0 {
		return fmt.Errorf("no job name specified")
	}
	var templateFile string
	// backward compatibility -- first version used --in for template in "write_breadcrumbs"
	// process_projects requires --in for the CSV and --template for the HTML template
	if args.Command == "write_breadcrumbs" {
		if len(args.InputFile) == 0 {
			return fmt.Errorf("no HTML template (--in) specified")
		}
		templateFile = args.InputFile
	} else {
		if len(args.TemplateFile) == 0 {
			return fmt.Errorf("no HTML template (--template) specified")
		}
		templateFile = args.TemplateFile
	}
	verbose := args.Verbose
	deleteDirCrumbs := args.DeleteDirCrumbs

	err := createStartFile(job)
	if err != nil {
		return fmt.Errorf("could not create start file %v\n", err)
	}
	tmpl := template.Must(template.ParseFiles(templateFile))

	// iterate through all pages and extract just names
	count := int64(0)
	offset := int64(args.Start)
	more := true
	if verbose {
		log.Printf("doGetManifest(%s)", job)
	}
	mySsc := ssc
	for more {
		ret, isTruncated, err := doGetManifest(mySsc, job, offset, verbose)
		if err != nil {
			return fmt.Errorf("get manifest %s failed %v\n", job, err)
		}

		err = doBreadcrumbs(tmpl, ret, job, args.Suffix, deleteDirCrumbs, verbose)
		if err != nil {
			if errors.As(err, &HasWarningsError{}) {
				// warnings, keep going
				_ = createWarningFile(job, err.Error())
			} else {
				return fmt.Errorf("could not write breadcrumbs %v\n", err)
			}
		}
		offset += limit
		count += int64(len(ret))
		more = isTruncated

		if more && offset%updateTokenEvery == 0 {
			// update token
			mySsc, err = mySsc.updateToken()
			if err != nil {
				return fmt.Errorf("could not update token after %d records %v\n", count, err)
			}
			if verbose {
				log.Print("Updated token")
			}
		}
	}
	err = createSuccessFile(job, count)
	if err != nil {
		return fmt.Errorf("could not create success file %v\n", err)
	}
	return nil
}

func doBreadcrumbs(tmpl *template.Template, files []openapi.ApiManifestFile,
	job string, suffix string, deleteDirCrumbs bool, verbose bool) error {

	warnings := 0
	currentContainingDirectory := ""
	for fileIndex := range files {
		file := files[fileIndex]
		fullPath := *file.Path + suffix
		if *file.IsDir {
			// dont create files for directories
			if deleteDirCrumbs {
				// old version left some? Remove them.
				if verbose {
					log.Printf("Delete directory crumbs: %s", fullPath)
				}
				err := os.Remove(fullPath)
				if err != nil {
					log.Printf("Failed to delete directory %s\n%v", fullPath, err)
				}
			}
			continue
		}
		var f *os.File
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
		if err == nil {
			tmpl.Execute(f, info)
			err = f.Close()
			if err != nil {
				// log and move on
				log.Printf("Failed to close file %s\n", fullPath)
			}
		} else {
			// increment warnings, log and move on
			warnings++
			log.Printf("ERROR: failed to create file %s\n%v", fullPath, err)
		}
	}
	if warnings > 0 {
		return &HasWarningsError{
			message: fmt.Sprintf("WARNING: %d files not created for job %s\n", warnings, job),
		}
	}
	return nil
}

func processProjects(ssc *SscClient, args *Arguments) error {
	verbose := args.Verbose
	projects, err := getJobsToProcess(ssc, args)
	if err != nil {
		return fmt.Errorf("could not process projects %v\n", err)
	}
	log.Printf("Processing %d projects\n", len(*projects))
	for projectName := range *projects {
		if verbose {
			log.Printf("Processing project: %s\n", projectName)
		}
		for _, jobName := range (*projects)[projectName] {
			if successFileExists(jobName) {
				if verbose {
					log.Printf("Job %s for %s has already completed\n", projectName, jobName)
				}
				break
			} else {
				if verbose {
					log.Printf("Processing job: %s\n", jobName)
				}
				err = breadcrumbsForOneProject(ssc, jobName, args)
				if err != nil {
					log.Printf("ERROR: could not process job %s %v\n", jobName, err)
					errorFileErr := createErrorFile(jobName, err)
					if errorFileErr != nil {
						log.Printf("ERROR: could not create error file %s %v\n", jobName, errorFileErr)
					}
				}
			}
		}
	}
	return nil
}

type ProjectJobMap map[string][]string

func getJobsToProcess(ssc *SscClient, args *Arguments) (*ProjectJobMap, error) {
	ret := make(ProjectJobMap)
	var projectNames []string
	var err error
	verbose := args.Verbose

	// fileName runs once, inputFile iterates through CSV
	if len(args.FileName) > 0 {
		projectNames = []string{args.FileName}
	} else if len(args.InputFile) > 0 {
		projectNames, err = loadFilenames(args.InputFile)
		if err != nil {
			return nil, fmt.Errorf("could not load project names from %s %v\n", args.InputFile, err)
		}
	} else {
		return nil, fmt.Errorf("no project name or input file specified")
	}
	if verbose {
		log.Printf("%d projects to search", len(projectNames))
	}

	for projectNameIndex := range projectNames {
		projectName := projectNames[projectNameIndex]
		jobNames := make([]string, 0)
		// update token
		mySsc, err := ssc.updateToken()
		if err != nil {
			return nil, fmt.Errorf("could not update token %v\n", err)
		}
		if verbose {
			log.Print("Updated token")
		}
		response, err := getJobsForProject(mySsc, projectName)
		if err != nil {
			return nil, fmt.Errorf("search objects for match %s failed %v\n", projectName, err)
		}
		jobs := response.Data
		sort.Slice(jobs, func(i, j int) bool {
			return jobs[i].CreatedTime.After(*jobs[j].CreatedTime)
		})
		for jobIndex := range jobs {
			job := jobs[jobIndex]
			jobNames = append(jobNames, *job.Name)
		}
		ret[projectName] = jobNames
	}
	return &ret, nil
}

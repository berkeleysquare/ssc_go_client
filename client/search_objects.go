package client

import (
	"encoding/csv"
	"fmt"
	"github.com/SpectraLogic/ssc_go_client/mongo_client"
	"github.com/SpectraLogic/ssc_go_client/openapi"
	"log"
	"os"
	"strings"
)

const (
	maxFilesPerJob = 1000
)

func PrintSearchCsvHeader(w *csv.Writer) error {
	var line = []string{"Key", "Job", "Project", "Share"}
	return w.Write(line)
}

func doSearch(ssc *SscClient, FileName string, exts []string, verbose bool) ([]openapi.ApiJob, error) {

	if verbose {
		log.Printf("doSearch(%s, %v)", FileName, exts)
	}

	// search for all files including case number
	if len(FileName) == 0 {
		return nil, fmt.Errorf("no match string specified")
	}
	if verbose {
		log.Printf("SearchObjects(%s)", FileName)
	}
	ctx, err := ssc.getContext(false)
	if err != nil {
		return nil, fmt.Errorf("getContext() failed %v\n", err)
	}
	response, _, err := ssc.Client.ProjectApi.SearchObjects(*ctx, FileName, nil)
	if err != nil {
		return nil, fmt.Errorf("search objects for match %s failed %v\n", FileName, err)
	}
	if verbose {
		log.Printf("SearchObjects(%s) returned %d jobs", FileName, len(response.Data))
	}

	// Try each supplied extension
	filesByProject := make(map[string][]string)
	for jobIndex := range response.Data {
		job := response.Data[jobIndex]
		project := *job.Project
		if verbose {
			log.Printf("getJobFiles(%s)", *job.Name)
		}
		jobFileNames, err := getJobFiles(ssc, &job, FileName, exts, verbose)
		if err != nil {
			return nil, fmt.Errorf("get files for job %s match %s failed %v\n", *job.Name, FileName, err)
		}
		filenames := *jobFileNames
		if verbose {
			log.Printf("getJobFiles(%s) returned %d files", *job.Name, len(filenames))
		}
		if filenames != nil {
			var matching []string
			if len(exts) == 0 {
				matching = filenames
			} else {
				matching = []string{}
				for objectIndex := range filenames {
					objectName := filenames[objectIndex]
					for extIndex := range exts {
						if strings.HasSuffix(objectName, exts[extIndex]) {
							matching = append(matching, objectName)
						}
					}
				}
			}
			if len(matching) > 0 {
				_, ok := filesByProject[project]
				if !ok {
					filesByProject[project] = []string{}
				}
				filesByProject[project] = append(filesByProject[project], matching...)
			}
			if verbose {
				log.Printf("getJobFiles(%s) returned %d matching files", *job.Name, len(matching))
			}
		}
	}
	if verbose {
		log.Printf("Total project matches for %s: %d", FileName, len(filesByProject))
	}
	ret := []openapi.ApiJob{}
	for project, projectFiles := range filesByProject {
		ret = append(ret, *openapi.MakeApiJobWithProject(project, projectFiles))
	}
	return ret, nil
}

func getJobFiles(ssc *SscClient, job *openapi.ApiJob, fileName string, exts []string, verbose bool) (*[]string, error) {
	if job == nil {
		return nil, fmt.Errorf("no job submitted to getJobFiles()")
	}

	// search returns the first 100 objects
	// if the job contains fewer than 100, just use those
	if job.Filenames != nil && len(*job.Filenames) < 100 {
		return job.Filenames, nil
	}

	// get manifest for job
	jobName := *job.Name
	if verbose {
		log.Printf("getManifestFiles(%s)", jobName)
	}
	filenames, err := getManifestFiles(ssc, jobName, fileName, exts, verbose)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve manifest for %s\n%v", jobName, err)
	}
	return &filenames, nil
}

func doRestore(ssc *SscClient, jobs []openapi.ApiJob, Share string, FileName string, Directory string, verbose bool) error {

	if verbose {
		log.Printf("doRestore() for %d jobs", len(jobs))
	}
	// create restore project(s) for each job
	for jobIndex := range jobs {
		job := jobs[jobIndex]
		if verbose {
			log.Printf("CreateSpecificFilesRestoreProject(%s)", *job.Name)
		}
		Filenames := *job.Filenames
		totalFilesThisJob := len(Filenames)
		filesSentThisJob := 0
		for totalFilesThisJob > filesSentThisJob {
			rightSlice := filesSentThisJob + maxFilesPerJob
			if rightSlice > totalFilesThisJob {
				rightSlice = totalFilesThisJob
			}
			thisJobSlice := Filenames[filesSentThisJob:rightSlice]
			if verbose {
				log.Printf("CreateSpecificFilesRestoreProject files %d-%d", filesSentThisJob, rightSlice)
			}
			if len(thisJobSlice) > 0 {
				err := CreateSpecificFilesRestoreProjectV4(ssc, Share, FileName, *job.Name, Directory, &thisJobSlice)
				if err != nil {
					return fmt.Errorf("failed to create restore job Restore_%s_%s %v\n", FileName, *job.Name, err)
				}
			}
			filesSentThisJob = rightSlice
		}

		if verbose {
			log.Printf("CreateSpecificFilesRestoreProject(%s) succeeded", *job.Name)
		}
	}

	log.Printf("Successfully ran Command\n")
	return nil
}
func executeSearch(ssc *SscClient, args *Arguments) error {

	var fileNames []string
	var err error
	verbose := args.Verbose
	var exts []string
	// no extensions or "*" means all files
	if args.Extensions != nil && len(args.Extensions) > 0 && args.Extensions[0] != "*" {
		exts = args.Extensions
	} else {
		exts = []string{}
	}
	// fileName runs once, inputFile iterates through CSV
	if len(args.FileName) > 0 {
		fileNames = []string{args.FileName}
	} else if len(args.InputFile) > 0 {
		fileNames, err = loadFilenames(args.InputFile)
		if err != nil {
			return fmt.Errorf("could not load file names from %s %v\n", args.InputFile, err)
		}
	} else {
		return fmt.Errorf("no match string or input file specified")
	}
	if verbose {
		log.Printf("%d files to search", len(fileNames))
	}

	if args.Command == "restore_objects" && len(args.Share) == 0 {
		return fmt.Errorf("no share specified for restore")
	}
	// output -- console, csv file, or none
	var w *csv.Writer
	outputFile := args.OutputFile
	writeOutput := args.Command == "search_objects" || len(outputFile) > 0

	if writeOutput {
		// list search results
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

		err := PrintSearchCsvHeader(w)
		if err != nil {
			return fmt.Errorf("could not print search results header %v", err)
		}
	}

	for fileNameIndex := range fileNames {
		fileName := fileNames[fileNameIndex]
		// update token
		mySsc, err := ssc.updateToken()
		if err != nil {
			return fmt.Errorf("could not update token %v\n", err)
		}
		if verbose {
			log.Print("Updated token")
		}
		ret, err := doSearch(mySsc, fileName, exts, verbose)
		if err != nil {
			return fmt.Errorf("search objects for match %s failed %v\n", fileName, err)
		}

		// list if command is search or if output file supplied to restore
		if writeOutput {
			err = displayJobObjectsV4(w, ret)
			if err != nil {
				return fmt.Errorf("could not list search results %v\n", err)
			}
		}

		if args.Command == "restore_objects" {
			err = doRestore(mySsc, ret, args.Share, fileName, args.Directory, verbose)
			if err != nil {
				return fmt.Errorf("failed to create restore jobs for %s %v\n", args.FileName, err)
			}
		}
	}
	log.Printf("Successfully ran Command\n")
	return nil
}

func loadFilenames(inputFile string) ([]string, error) {
	f, err := os.Open(inputFile)
	if err != nil {
		return nil, fmt.Errorf("could not open input file %s %v\n", inputFile, err)
	}
	defer f.Close()

	// read whole file using csv.Reader
	csvReader := csv.NewReader(f)
	data, err := csvReader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("could not read from input file %s %v\n", inputFile, err)
	}

	// expect numbers in first column starting on second line
	var ret []string
	for i, line := range data {
		if i > 0 {
			for j, field := range line {
				if j == 0 {
					ret = append(ret, field)
				}
			}
		}
	}
	return ret, nil
}

func displayJobObjects(w *csv.Writer, jobs []openapi.ApiJob) error {

	lines := [][]string{}
	for jobIndex := range jobs {
		job := jobs[jobIndex]
		jobName := job.Name
		project := ""
		share := ""
		fileNames := *job.Filenames
		for fileIndex := range fileNames {
			lines = append(lines, []string{fileNames[fileIndex], *jobName, project, share})
		}
	}
	return w.WriteAll(lines)
}

func displayJobObjectsV4(w *csv.Writer, jobs []openapi.ApiJob) error {

	lines := [][]string{}
	for jobIndex := range jobs {
		job := jobs[jobIndex]
		jobName := job.Name
		fileNames := *job.Filenames
		project := job.Project
		share := ""
		for fileIndex := range fileNames {
			lines = append(lines, []string{fileNames[fileIndex], *jobName, *project, share})
		}
	}
	return w.WriteAll(lines)
}

func loadFilesForRestore(inputFile string) ([]*mongo_client.SearchObject, error) {
	f, err := os.Open(inputFile)
	if err != nil {
		return nil, fmt.Errorf("could not open input file %s %v\n", inputFile, err)
	}
	defer f.Close()

	// read whole file using csv.Reader
	csvReader := csv.NewReader(f)
	data, err := csvReader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("could not read from input file %s %v\n", inputFile, err)
	}

	// expect numbers in first column starting on second line
	var ret []*mongo_client.SearchObject
	for i, line := range data {
		if i > 0 {
			if len(line) >= 3 {
				ret = append(ret, &mongo_client.SearchObject{Name: line[0], Manifest: line[1], Share: line[3]})
			}
		}
	}
	return ret, nil
}

package client

import (
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/SpectraLogic/ssc_go_client/mongo_client"
	"github.com/SpectraLogic/ssc_go_client/openapi"
	"log"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	objects_per_page = 500000
	FILE_OK          = "OK"
	FILE_MISSONG     = "ERROR -- not exist"
	FILE_STAT_ERROR  = "ERROR -- could not stat"
)

func executeDbSearch(ssc *SscClient, args *Arguments) error {

	var fileNames []string
	var err error
	exts := []string(args.Extensions)
	verbose := args.Verbose
	restore := args.Command == "restore_db_objects"

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
	if restore && len(args.Share) == 0 {
		return fmt.Errorf("no share specified")
	}
	if verbose {
		log.Printf("%d files to search", len(fileNames))
	}

	// output -- console, csv file, or none
	var w *csv.Writer
	outputFile := args.OutputFile
	writeOutput := args.Command == "search_db" || len(outputFile) > 0

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
			return fmt.Errorf("could not print db search results header %v", err)
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
		ret, err := doDbSearch(mySsc, fileName, exts, verbose)
		if err != nil {
			return fmt.Errorf("search db for match %s failed %v\n", fileName, err)
		}

		// list if command is search or if output file supplied to restore
		if writeOutput {
			err = mongo_client.DisplaySearchObjects(w, ret)
			if err != nil {
				return fmt.Errorf("could not list db search results %v\n", err)
			}
			if verbose && len(args.OutputFile) > 0 {
				log.Printf("Results written to %s", args.OutputFile)
			}
		}

		if restore {
			//package results in apiJobs. One per job containing all job objects
			jobObjects := makeJobObjects(ret)
			err = doRestore(mySsc, jobObjects, args.Share, fileName, args.Directory, verbose)
			if err != nil {
				return fmt.Errorf("failed to create restore jobs for %s %v\n", args.FileName, err)
			}
		}
	}

	log.Printf("Successfully ran Command\n")
	return nil
}

func executeDbProjectSearch(ssc *SscClient, args *Arguments) error {

	verbose := args.Verbose
	projectName := args.ProjectName
	share := args.Share
	path := ""
	// assert prints only error files -- verify prints all (with status)
	assert := args.Command == "assert_db_project"
	verify := assert || args.Command == "verify_db_project"
	pageSize := objects_per_page
	if args.Count < math.MaxInt {
		pageSize = args.Count
	}

	var err error
	var projectNames []string

	// fileName runs once, inputFile iterates through CSV
	if len(projectName) > 0 {
		projectNames = []string{projectName}
	} else if len(args.InputFile) > 0 {
		projectNames, err = loadFilenames(args.InputFile)
		if err != nil {
			return fmt.Errorf("could not load project names from %s %v\n", args.InputFile, err)
		}
	} else {
		return fmt.Errorf("no project name or input file specified")
	}

	for projectNameIndex := range projectNames {
		project := projectNames[projectNameIndex]
		if verbose {
			log.Printf("execute command %s on project %s", args.Command, project)
		}

		// update token
		mySsc, err := ssc.updateToken()
		if err != nil {
			return fmt.Errorf("could not update token %v\n", err)
		}
		if verbose {
			log.Print("Updated token")
		}

		if verify {
			if len(share) == 0 {
				return fmt.Errorf("must specify share (location where files were restored)")
			}
			location, _, err := getStorageLocation(mySsc, share)
			if err != nil {
				return fmt.Errorf("could not get information about share %s\n%v", share, err)
			}
			path = filepath.Join(*location.Path, args.Directory)
			if verbose {
				log.Printf("verify to path %s", path)
				if args.JobToPath {
					log.Printf("Job name (manifest) will be appended to share path")
				}
			}
		}

		// output -- console, csv file, or none
		var w *csv.Writer
		var f *os.File
		outputFile := args.OutputFile

		// set up pagination
		more := true
		offset := 0
		for more {
			// grab a page of data
			ret, err := doDbProjectSearch(mySsc, project, offset, pageSize, verbose)
			if err != nil {
				return fmt.Errorf("search db for project %s failed %v\n", project, err)
			}

			// open a new file for each page (Excel limits rows to 1,024,000)
			wOut := os.Stdout
			if len(outputFile) > 0 {
				outputFile = makeOutputFileName(args.OutputFile, project, offset)
				f, err = os.Create(outputFile)
				if err != nil {
					return fmt.Errorf("Could not create %s\n%v\n", outputFile, err)
				}
				wOut = f
			}
			w = csv.NewWriter(wOut)

			err = PrintSearchCsvHeader(w)
			if err != nil {
				return fmt.Errorf("could not print db search results header %v", err)
			}

			if verify {
				// check that all files exist on share
				err = verifyFilesExist(w, ret, path, !assert, args.JobToPath)
				if err != nil {
					return fmt.Errorf("verify project %s on path %s failed %v\n", project, path, err)
				}
			} else {
				err = mongo_client.DisplaySearchObjects(w, ret)
				if err != nil {
					return fmt.Errorf("could not list db search results %v\n", err)
				}
			}
			if verbose && len(outputFile) > 0 {
				log.Printf("Results written to %s", outputFile)
			}
			offset += pageSize
			more = pageSize == len(ret)
			w.Flush()
			if len(outputFile) > 0 {
				err = f.Close()
				if err != nil {
					return fmt.Errorf("could not close output file %s %v\n", outputFile, err)
				}
			}
		}
	}

	log.Printf("Successfully ran Command\n")
	return nil
}

func doDbSearch(ssc *SscClient, FileName string, exts []string, verbose bool) ([]*mongo_client.SearchObject, error) {

	if verbose {
		log.Printf("doDbSearch(%s, %v)", FileName, exts)
	}

	// search for all files including case number
	if len(FileName) == 0 {
		return nil, fmt.Errorf("no match string specified")
	}
	if verbose {
		log.Printf("SearchObjects(%s, %s)", FileName, exts)
	}

	response, err := mongo_client.RunQuery(FileName, exts)
	if err != nil {
		return nil, fmt.Errorf("search objects for match %s ext %s failed %v\n", FileName, exts, err)
	}

	if verbose {
		log.Printf("Total matches for %s: %d", FileName, len(response))
	}
	return response, nil
}

func doDbProjectSearch(ssc *SscClient, project string, offset int, limit int, verbose bool) ([]*mongo_client.SearchObject, error) {

	if verbose {
		log.Printf("doDbProjectSearch(%s, %d, %d)", project, offset, limit)
	}

	// search for all files including case number
	if len(project) == 0 {
		return nil, fmt.Errorf("no project name specified")
	}
	if verbose {
		log.Printf("SearchProjectObjects(%s)", project)
	}

	response, err := mongo_client.RunProjectQuery(project, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("search objects for project %s failed %v\n", project, err)
	}

	if verbose {
		log.Printf("Total matches for project %s: %d", project, len(response))
	}
	return response, nil
}

func makeJobObjects(searchObjs []*mongo_client.SearchObject) []openapi.ApiJob {
	ret := []openapi.ApiJob{}
	for objIndex := range searchObjs {
		obj := searchObjs[objIndex]
		// add the object under its job
		var manifest *openapi.ApiJob
		for jobIndex := range ret {
			job := ret[jobIndex]
			if *job.Name == obj.Manifest {
				manifest = &job
				break
			}
		}
		if manifest == nil {
			newnames := []string{}
			manifest = openapi.MakeApiJob(obj.Manifest, newnames)
			ret = append(ret, *manifest)
		}
		filenames := *manifest.Filenames
		filenames = append(filenames, obj.Name)
		*manifest.Filenames = filenames
	}
	return ret
}

func verifyFilesExist(w *csv.Writer, files []*mongo_client.SearchObject, path string, showAll bool, jobToPath bool) error {
	lines := [][]string{}
	for fileIndex := range files {
		file := files[fileIndex]
		var testPath string
		if jobToPath {
			testPath = filepath.Join(path, file.Manifest, file.Path)
		} else {
			testPath = filepath.Join(path, file.Path)
		}
		status := verifyFileExist(testPath)
		if showAll || status != FILE_OK {
			lines = append(lines, []string{file.Path, file.Manifest, strconv.Itoa(file.Size), file.Checksum, status})
		}
	}
	return w.WriteAll(lines)
}

func verifyFileExist(fullpath string) string {
	_, err := os.Stat(fullpath)
	if err == nil {
		return FILE_OK
	}
	// not there, keep waiting
	if errors.Is(err, os.ErrNotExist) {
		return FILE_MISSONG
	}
	return FILE_STAT_ERROR
}

func makeOutputFileName(name string, project string, start int) string {
	ret := strings.TrimSuffix(name, ".csv")
	ret += "_" + project
	if start > 0 {
		ret += "_" + strconv.Itoa(start)
	}
	return ret + ".csv"
}

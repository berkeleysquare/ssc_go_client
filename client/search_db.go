package client

import (
	"encoding/csv"
	"fmt"
	"github.com/SpectraLogic/ssc_go_client/mongo_client"
	"github.com/SpectraLogic/ssc_go_client/openapi"
	"log"
	"os"
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
		return fmt.Errorf("no match string or input file specified" )
	}
	if restore && len(args.Share) == 0 {
		return fmt.Errorf("no share specified" )
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
			if verbose && 	len(args.OutputFile) > 0 {
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

	log.Printf("Successfully ran Command\n",)
	return nil
}


func doDbSearch(ssc *SscClient, FileName string, exts []string, verbose bool) ([]*mongo_client.SearchObject, error) {

	if verbose {
		log.Printf("doDbSearch(%s, %v)", FileName, exts)
	}

	// search for all files including case number
	if len(FileName) == 0 {
		return nil, fmt.Errorf("no match string specified" )
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

func makeJobObjects(searchObjs []*mongo_client.SearchObject) []openapi.ApiJob {
	ret := []openapi.ApiJob {}
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
package client

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const TEST_SOURCE_NAME = "verify-test-source"
const TEST_SOURCE_DIRECTORY_NAME = "verify-test-source"
const TEST_TARGET_NAME = "verify-test-target"
const TEST_SOURCE_FILE = "verify-test-file.txt"
const TEST_FILE_HEADER = "Spectra Logic StorCycle test file.\n\nhttps://spectralogic.com/\n\n"
const TEST_MIGRATE_PROJECT = "verify-test-migrate"
const TEST_RESTORE_PROJECT = "verify-test-restore"
const TEST_BUCKET_NAME = "verify-test-bucket"

func MakeTimestamp() string {
	return string(time.Now().Format("060102150405"))
}

func CreateTestFile(ssc *SscClient, args *Arguments) error {

	if len(args.Directory) == 0 {
		return fmt.Errorf("no source path specified" )
	}

	outputFile := filepath.Join(args.Directory, ValueOrDefault(args.FileName, TEST_SOURCE_FILE))
	wOut, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("Could not create %s\n%v\n", outputFile, err)
	}
	defer wOut.Close()
	w := io.Writer(wOut)
	_, err = w.Write([]byte(TEST_FILE_HEADER))
	if err != nil {
		return fmt.Errorf("failed to write header to source file (%s) %v\n", outputFile, err)
	}
	timestamp := "Created:" + string(time.Now().Format("06-01-02-15-04-05")) + "\n"
	_, err = w.Write([]byte(timestamp))
	if err != nil {
		return fmt.Errorf("failed to write timestamp to source file (%s) %v\n", outputFile, err)
	}

	fmt.Printf("Successfully created test file %s\n", outputFile)
	return nil
}

func getTestSourceDirectory(rootDir string) (string, error) {
	sourceDir := filepath.Join(rootDir, TEST_SOURCE_DIRECTORY_NAME)

	_, err := os.Stat(sourceDir)
	if err == nil {
		// it's there
		return sourceDir, nil
	}
	if !errors.Is(err, os.ErrNotExist) {
		//something else wrong
		return "", fmt.Errorf("could not see if store directory %s exists\n%v\n", sourceDir, err)
	}
	// not there -- create dir
	err = os.Mkdir(sourceDir, os.ModeDir)
	return sourceDir, err
}

func CreateTestMigrate(ssc *SscClient, args *Arguments) error {

	args.ProjectName = ValueOrDefault(args.ProjectName, TEST_MIGRATE_PROJECT)
	args.Share = ValueOrDefault(args.Share, TEST_SOURCE_NAME)
	args.Target = ValueOrDefault(args.Target, TEST_TARGET_NAME)

	// does the project already exist?
	_, resp, err := ssc.Client.ProjectApi.GetProject(*ssc.Context, args.ProjectName)
	if err == nil && resp.StatusCode < 300 {
		fmt.Printf("Migrate project %s already exists\n", args.ProjectName)
		return RunNow(ssc, args)
	}

	// Nope, create it
	return CreateMigrateProject(ssc, args)
}

func CreateTestRestore(ssc *SscClient, args *Arguments) error {

	args.ProjectName = ValueOrDefault(args.ProjectName, TEST_RESTORE_PROJECT)
	args.Share = ValueOrDefault(args.Share, TEST_SOURCE_NAME)
	args.Target = ValueOrDefault(args.Target, TEST_TARGET_NAME)
	args.Directory = string(filepath.Separator) + "restore"

	// does the project already exist?
	_, resp, err := ssc.Client.ProjectApi.GetProject(*ssc.Context, args.ProjectName)
	if err == nil && resp.StatusCode < 300 {
		fmt.Printf("Restore project %s already exists\n", args.ProjectName)
		return RunNow(ssc, args)
	}

	// Nope, create it
	return CreateRestoreProject(ssc, args)
}

func waitForRestore(restoreFile string) error {
	fmt.Printf("Waiting for restore .")
	return tryRestore(restoreFile, 0, 1)
}

func tryRestore(restoreFile string, fib1 int, fib2 int) error {
	_, err := os.Stat(restoreFile)
	if err == nil {
		fmt.Printf("\n")
		return nil
	}
	// not there, keep waiting
	if errors.Is(err, os.ErrNotExist) {
		fmt.Printf(".")
	} else {
		//something else wrong
		return fmt.Errorf("could not see if restore file %s exists\n%v\n", restoreFile, err)
	}
	// wait
	fib3 := staggeredWait(fib1, fib2)
	// try again
	return tryRestore(restoreFile, fib2, fib3)
}

func listLatestJob(ssc *SscClient, args *Arguments) error {
	_, err := getLatestJob(ssc, args.ProjectName)
	if err != nil {
		return fmt.Errorf("search jobs for project name %s failed %v\n", args.ProjectName, err)
	}
	return nil
}

func getLatestJob(ssc *SscClient, project string) (string, error) {
	response, _, err := ssc.Client.ProjectApi.SearchJobs(*ssc.Context, project, nil)
	if err != nil {
		return "", fmt.Errorf("search jobs for project name %s failed %v\n", project, err)
	}
	count := len(response.Data)
	if count > 0 {
		job := *response.Data[count -1].Name
		fmt.Printf("Latest job for %s is %s\n", project, job)
		return job, nil
	}

	return "", fmt.Errorf("no jobs match %s\n", project)
}

func FullVerifyCycle(ssc *SscClient, args *Arguments) error {
	// Must provide a BP location to clone
	if len(args.Clone) == 0 {
		return fmt.Errorf("no target location specified to provide credentials to copy" )
	}

	/*
		SET UP SOURCE ENVIRONMENT
	 */
	timestamp := MakeTimestamp()
	sourceDirectory, err := getTestSourceDirectory(args.Directory)
	if err != nil {
		return fmt.Errorf("Get test source directory failed %v\n", err)
	}
	args.Directory = sourceDirectory
	sourceFileName := strings.Replace(TEST_SOURCE_FILE,  ".", timestamp + ".", 1 )
	args.FileName = sourceFileName

	err = CreateTestFile(ssc, args)
	if err != nil {
		return fmt.Errorf("CreateTestFile failed %v\n", err)
	}
	originalHash, err := processHash(filepath.Join(args.Directory, args.FileName))
	if err != nil {
		return fmt.Errorf("Create checksum failed %v\n", err)
	}

	/*
		CREATE SOURCE AND TARGET STORAGE LOCATIONS
	*/
	args.Share = TEST_SOURCE_NAME
	err = CreateSourceLocation(ssc, args)
	if err != nil {
		return fmt.Errorf("Create source location failed %v\n", err)
	}

	args.Target = TEST_TARGET_NAME
	err = CreateTargetLocation(ssc, args)
	if err != nil {
		return fmt.Errorf("Create target location failed %v\n", err)
	}

	 /*
		CREATE AND EXECUTE MIGRATE PROJECT
	 */
	migrateProjectName := TEST_MIGRATE_PROJECT
	args.ProjectName = migrateProjectName
	args.Directory = string(filepath.Separator)
	err = CreateTestMigrate(ssc, args)
	if err != nil {
		return fmt.Errorf("CreateTestMigrate failed %v\n", err)
	}

	// wait for placement in cache and get key as it appears on BP
	fullPath, err := getTestFileFullPath(args)
	if err != nil {
		return fmt.Errorf("getTestFileFullPath failed %v\n", err)
	}
	args.FileName = fullPath


	// get job name for restore (cannot edit/reuse restore jobs)
	migrateJob, err := getLatestJob(ssc, migrateProjectName)
	if err != nil {
		return fmt.Errorf("get latest job failed for project %s\n%v\n", migrateProjectName, err)
	}

	if !args.DontWaitForTape {
		fmt.Printf("Waiting for physical placement on %s ", fullPath)
		err = waitForPlacement(ssc, args)
		if err != nil {
			return fmt.Errorf("waitForPlacement failed %v\n", err)
		}
	}

	/*
		CREATE AND EXECUTE RESTORE PROJECT
 	*/
	args.ProjectName = TEST_RESTORE_PROJECT + timestamp
	args.Job = migrateJob
	err = CreateTestRestore(ssc, args)
	if err != nil {
		return fmt.Errorf("CreateTestRestore failed %v\n", err)
	}
	restoreDirectory := sourceDirectory + "/restore"

	err = waitForRestore(filepath.Join(restoreDirectory, sourceFileName))
	if err != nil {
		return fmt.Errorf("wait for restored file failed %v\n", err)
	}

	outputFile := filepath.Join(restoreDirectory, ValueOrDefault(sourceFileName, TEST_SOURCE_FILE))
	restoredHash, err := processHash(outputFile)
	if err != nil {
		return fmt.Errorf("Create checksum of restored file failed %v\n", err)
	}
	fmt.Printf("Verify succeeded\n")
	fmt.Printf("Original checksum: %x, final checksum: %x\n", originalHash, restoredHash)
	return nil
}


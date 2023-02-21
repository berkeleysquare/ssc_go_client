package client

import (
    "fmt"
)

type command func(*SscClient, *Arguments) error

var availableCommands = map[string]command {
    "full_verify": FullVerifyCycle,
    "get_locations": ListStorageLocations,
    "create_target": CreateTargetLocation,
    "create_source": CreateSourceLocation,
    "create_test_file": CreateTestFile,
    "run_project_now": RunNow,
    "checksum_test_file": DoHash,
    "directory_checksum": HashDirectory,
    "create_test_migrate": CreateTestMigrate,
    "create_test_restore": CreateTestRestore,
    "get_migrate_projects": ListMigrateProjects,
    "get_restore_projects": ListRestoreProjects,
    "list_jobs": listJobs,
    "get_job_status": GetJobStatus,
    "get_all_jobs": GetAllJobs,
    "get_all_restore_jobs": GetAllJobs,
    "get_active_restore_jobs": GetAllJobs,
    "restore_jobs_by_tag" : GetRestoreJobsByTag,
    "active_restore_jobs_by_tag" : GetRestoreJobsByTag,
    "wait_for_restore_jobs_by_tag" : WaitForRestoreJobsByTag,
    "search_objects": executeSearch,
    "restore_objects": executeSearch,
    "get_catalog": listCatalog,
    "latest_job": listLatestJob,
    "get_scan_projects": ListScanProjects,
    "create_scan": CreateScanProject,
    "create_archive": CreateMigrateProject,
    "create_restore": CreateRestoreProject,
    "physical_placement": getPhysicalPlacement,
    "wait_for_placement": waitForPlacement,
    "head_object": headObject,
    "inventory": listBucketContents,
    "search_db": executeDbSearch,
    "restore_db_objects": executeDbSearch,
}

// commands require only filesystem or BlackPearl access, no ssc client required
var noTokenRequired = map[string]bool {
    "checksum_test_file": true,
    "directory_checksum": true,
    "physical_placement": true,
    "wait_for_placement": true,
    "head_object": true,
    "inventory": true,
}

func CommandRequiresClientToken(args *Arguments) (bool, error) {
    // make sure command exists
    _, ok := availableCommands[args.Command]
    if ok {
        // return false if it is in noTokenRequired and true
        noToken, inTable:= noTokenRequired[args.Command]
        return !(noToken && inTable), nil
    } else {
        return false, fmt.Errorf("unsupported command: '%s'", args.Command)
    }
}

func RunCommand(ssc *SscClient, args *Arguments) error {
    cmd, ok := availableCommands[args.Command]
    if ok {
        return cmd(ssc, args)
    } else {
        return fmt.Errorf("unsupported command: '%s'", args.Command)
    }
}

func ListCommands(args *Arguments) error {
    fmt.Printf("Usage: ssc_cli --command <command>\n",)
    for key, _ := range availableCommands {
        fmt.Printf("%s\n", key)
    }
    return nil
}

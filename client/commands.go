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
    "create_test_migrate": CreateTestMigrate,
    "create_test_restore": CreateTestRestore,
    "get_migrate_projects": ListMigrateProjects,
    "get_restore_projects": ListRestoreProjects,
    "list_jobs": listJobs,
    "latest_job": listLatestJob,
    "get_scan_projects": ListScanProjects,
    "create_scan": CreateScanProject,
    "create_archive": CreateMigrateProject,
    "create_restore": CreateRestoreProject,
    "physical_placement": getPhysicalPlacement,
    "wait_for_placement": waitForPlacement,
    "head_object": headObject,
    "inventory": listBucketContents,
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

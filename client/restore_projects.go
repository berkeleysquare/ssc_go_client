package client

import (
	"fmt"
	"github.com/SpectraLogic/ssc_go_client/openapi"
	"github.com/antihax/optional"
	"strings"
	"time"
)

func displayRestoreProjects(projects openapi.ApiProjectPaginator) error {
	for projectIndex := range projects.Data {
		project := projects.Data[projectIndex]
		fmt.Printf("Project: %s, Share: %s, Manifest: %s\n", *project.Status.Name, *project.Share, *project.RestoreManifest)
	}
	return nil
}

func ListRestoreProjects(ssc *SscClient, args *Arguments) error {
	opts := &openapi.ProjectApiListProjectsOpts{
		Skip:		optional.NewInt64(int64(args.Start)),
		Limit:		optional.NewInt64(int64(args.Count)),
		FilterBy:   optional.NewInterface([]string{"Restore", "RestoreBreadcrumb"}),
	}

	projects, resp, err := getProjects(ssc, opts)
	if err != nil {
		return fmt.Errorf("could not retrieve projects (%d) %v\n", resp.StatusCode, err)
	}
	return displayRestoreProjects(projects)
}

func CreateRestoreProject(ssc *SscClient, args *Arguments) error {

	// create and run a Scan Project for "share"
	if len(args.Share) == 0 {
		return fmt.Errorf("no share specified in create_restore" )
	}
	if len(args.Job) == 0 {
		return fmt.Errorf("no job/manifest specified in create_restore" )
	}
	projectName := args.ProjectName
	if len(projectName) == 0 {
		projectName = fmt.Sprintf("Restore_%s%s", args.Share, strings.Replace(args.Directory, "\\", "_", -1))
	}
	policyType := "Restore"
	breadCrumbAction := ""
	if args.MakeLinks == "single" {
		policyType = "RestoreBreadcrumb"
		breadCrumbAction = "CreateHtmlLinkSingle"
	}
	if args.MakeLinks == "individual" {
		policyType = "RestoreBreadcrumb"
		breadCrumbAction = "CreateHtmlLink"
	}
	if args.MakeLinks == "symbolic" {
		policyType = "RestoreBreadcrumb"
		breadCrumbAction = "CreateSymbolicLink"
	}

	timestamp := string(time.Now().Format("06-01-02-15-04-05"))
	description := fmt.Sprintf("%s, Created by API %s", projectName, timestamp)
	active := true
	tags := []string{args.Share}

	restoreDefinition := &openapi.ApiProjectRestore{
		Description:      &description,
		Share:            &args.Share,
		WorkingDirectory: &args.Directory,
		Active:           &active,
		Tags:             &tags,
		Schedule:         *NowSchedule(),
		ProjectType:	  &policyType,
		RestoreManifest:  &args.Job,
	}
	if len(breadCrumbAction) > 0 {
		restoreDefinition.BreadCrumbAction = &breadCrumbAction
	}

	restore, resp, err := ssc.Client.ProjectApi.UpdateRestoreProject(*ssc.Context, projectName, *restoreDefinition)
	if err != nil {
		return fmt.Errorf("failed to create/update restore (%d) %v\n", resp.StatusCode, err)
	}
	fmt.Printf("Successfully created restore project %s\n", *restore.Status.Name)
	return nil
}

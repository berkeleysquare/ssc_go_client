package client

import (
	"fmt"
	"github.com/SpectraLogic/ssc_go_client/openapi"
	"github.com/antihax/optional"
	"net/http"
	"strings"
	"time"
)

func getProjects(ssc *SscClient, opts *openapi.ProjectApiListProjectsOpts) (openapi.ApiProjectPaginator, *http.Response, error) {
	return ssc.Client.ProjectApi.ListProjects(*ssc.Context, opts)
}

func displayMigrateProjects(projects openapi.ApiProjectPaginator) error {
	for projectIndex := range projects.Data {
		project := projects.Data[projectIndex]
		fmt.Printf("Project: %s, Share: %s, Targets: %v\n", *project.Status.Name, *project.Share, *project.Targets)
	}
	return nil
}

func ListMigrateProjects(ssc *SscClient, args *Arguments) error {
	opts := &openapi.ProjectApiListProjectsOpts{
		Skip:		optional.NewInt64(int64(args.Start)),
		Limit:		optional.NewInt64(int64(args.Count)),
		FilterBy:   optional.NewInterface([]string{"ScanAndArchive", "Archive"}),
	}

	projects, resp, err := getProjects(ssc, opts)
	if err != nil {
		return fmt.Errorf("could not retrieve projects (%d) %v\n", resp.StatusCode, err)
	}
	return displayMigrateProjects(projects)
}

func CreateMigrateProject(ssc *SscClient, args *Arguments) error {

	// create and run a Scan Project for "share"
	if len(args.Share) == 0 {
		return fmt.Errorf("no share specified in create_archive" )
	}
	if len(args.Target) == 0 {
		return fmt.Errorf("no target specified in create_archive" )
	}
	projectName := args.ProjectName
	if len(projectName) == 0 {
		projectName = fmt.Sprintf("Archive_%s%s", args.Share, strings.Replace(args.Directory, "\\", "_", -1))
	}

	timestamp := string(time.Now().Format("06-01-02-15-04-05"))
	description := fmt.Sprintf("%s, Created by API %s", projectName, timestamp)
	active := true
	policyType := "ScanAndArchive"
	tags := []string{args.Share}
	targets := []string{args.Target}

	breadCrumbAction := "KeepOriginal"
	if args.MakeLinks == "single" {
		breadCrumbAction = "CreateHtmlLinkSingle"
	}
	if args.MakeLinks == "individual" {
		breadCrumbAction = "CreateHtmlLink"
	}
	if args.MakeLinks == "symbolic" {
		breadCrumbAction = "CreateSymbolicLink"
	}

	archiveDefinition := &openapi.ApiProjectArchive{
		Description:      &description,
		Share:            &args.Share,
		WorkingDirectory: &args.Directory,
		Active:           &active,
		Tags:             &tags,
		Targets:   		  &targets,
		Schedule:         *NowSchedule(),
		ProjectType:	  &policyType,
		BreadCrumbAction: &breadCrumbAction,
	}

	if len(args.IncludeDirectory) > 0 || len(args.ExcludeDirectory) > 0 {
		filter := &openapi.ApiProjectFilter{}
		if len(args.IncludeDirectory) > 0 {
			filter.IncludeDirectories = &[]string{args.IncludeDirectory}
		}
		if len(args.ExcludeDirectory) > 0 {
			filter.ExcludeDirectories = &[]string{args.ExcludeDirectory}
		}
		archiveDefinition.Filter = *filter
	}

	store, resp, err := ssc.Client.ProjectApi.UpdateArchiveProject(*ssc.Context, projectName, *archiveDefinition)
	if err != nil {
		return fmt.Errorf("failed to create/update archive (%d) %v\n", resp.StatusCode, err)
	}
	fmt.Printf("Successfully created archive project %s\n", *store.Status.Name)
	return nil
}

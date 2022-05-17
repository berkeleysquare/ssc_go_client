package client

import (
	"fmt"
	"github.com/SpectraLogic/ssc_go_client/openapi"
	"github.com/antihax/optional"
	"time"
)


func displayScanProjects(projects openapi.ApiProjectPaginator) error {
	for projectIndex := range projects.Data {
		project := projects.Data[projectIndex]
		fmt.Printf("Project: %s, Share: %s\n", *project.Status.Name, *project.Share)
	}
	return nil
}

func ListScanProjects(ssc *SscClient, args *Arguments) error {
	opts := &openapi.ProjectApiListProjectsOpts{
		Skip:		optional.NewInt64(int64(args.Start)),
		Limit:		optional.NewInt64(int64(args.Count)),
		FilterBy:   optional.NewInterface([]string{"Scan"}),
	}

	projects, resp, err := getProjects(ssc, opts)
	if err != nil {
		return fmt.Errorf("could not retrieve projects (%d) %v\n", resp.StatusCode, err)
	}
	return displayScanProjects(projects)
}

func CreateScanProject(ssc *SscClient, args *Arguments) error {
	share := args.Share

	// create and run a Scan Project for "share"
	if len(share) > 0 {
		timestamp := string(time.Now().Format("06-01-02-15-04-05"))
		name := fmt.Sprintf("Scan_%s_%s", share, timestamp)
		description := fmt.Sprintf("%s, Created by API", name)
		active := true
		tags := []string{share}
		scanDefinition := &openapi.ApiProjectScan{
			Description:      &description,
			Share:            &share,
			WorkingDirectory: &args.Directory,
			Active:           &active,
			Tags:             &tags,
			Schedule:         *NowSchedule(),
		}

		scan, resp, err := ssc.Client.ProjectApi.UpdateScanProject(*ssc.Context, name, *scanDefinition)
		if err != nil {
			return fmt.Errorf("failed to create/update scan (%d) %v\n", resp.StatusCode, err)
		}
		fmt.Printf("Successfully created scan project %s\n", *scan.Status.Name)
		return nil
	}
	return fmt.Errorf("no share specified in create_scan" )
}

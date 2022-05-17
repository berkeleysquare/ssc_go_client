package client

import (
	"fmt"
	"github.com/SpectraLogic/ssc_go_client/openapi"
	"github.com/antihax/optional"
	"net/http"
	"strconv"
	"time"
)

func getStorageLocation(ssc *SscClient, name string) (openapi.ApiStorageLocation, *http.Response, error) {
	return ssc.Client.StorageApi.GetStorageLocation(*ssc.Context, name)
}

func getStorageLocations(ssc *SscClient, opts *openapi.StorageApiListStorageLocationsOpts) (openapi.ApiStorageLocationPaginator, *http.Response, error) {
	return ssc.Client.StorageApi.ListStorageLocations(*ssc.Context, opts)
}

func displayStorageLocations(locations openapi.ApiStorageLocationPaginator) error {
	for locationIndex := range locations.Data {
		location := locations.Data[locationIndex]
		fmt.Printf("Location: %s, type: %s, isTarget: %s\n", *location.Name, *location.Type, strconv.FormatBool(*location.IsTarget))
	}
	return nil
}

func ListStorageLocations(ssc *SscClient, args *Arguments) error {
	opts := &openapi.StorageApiListStorageLocationsOpts{
		Skip:	optional.NewInt64(int64(args.Start)),
		Limit:	optional.NewInt64(int64(args.Count)),
	}

	locations, resp, err := getStorageLocations(ssc, opts)
	if err != nil {
		return fmt.Errorf("could not retrieve storage locations (%d) %v\n", resp.StatusCode, err)
	}
	return displayStorageLocations(locations)
}

func CreateTargetLocation(ssc *SscClient, args *Arguments) error {
	// Clone a BP object endpoint
	if len(args.Clone) == 0 {
		return fmt.Errorf("no target location specified to provide credentials to copy" )
	}

	name := ValueOrDefault(args.Target, TEST_TARGET_NAME)
	// does the target already exist?
	_, _, err := getStorageLocation(ssc, name)
	if err == nil {
		fmt.Printf("Target storage location %s already exists\n", name)
		return nil
	}

	// get the endpoint to clone
	location, _, err := getStorageLocation(ssc, args.Clone)
	if err != nil {
		return fmt.Errorf("could not retrieve storage location to clone %s %v\n", args.Target, err)
	}

	mgmtEndpoint := location.SpectraMgmtEndpoint
	dataEndpoint := location.SpectraDataEndpoint
	timestamp := string(time.Now().Format("06-01-02-15-04-05"))
	description := fmt.Sprintf("%s, Created by Verify Test %s", name, timestamp)
	testBucketName := TEST_BUCKET_NAME
	locationType := "BlackPearl"
	locationIsTarget := true

	sourceDefinition := &openapi.ApiStorageLocationBlackPearl{
		SpectraMgmtEndpoint: mgmtEndpoint,
		SpectraDataEndpoint: dataEndpoint,
		Description:    	 &description,
		Bucket:				 &testBucketName,
		Type:				 &locationType,
		IsTarget:			 &locationIsTarget,
	}


	varOptions := &openapi.StorageApiUpdateStorageLocationOpts{
		CloneCredentials: optional.NewString(args.Clone),
	}

	_, _, err = ssc.Client.StorageApi.UpdateStorageLocation(*ssc.Context, name, *sourceDefinition, varOptions)
	if err != nil {
		return fmt.Errorf("failed to create target (%s) %v\n", name, err)
	}
	fmt.Printf("Successfully created target %s\n", name)
	return nil
}

func CreateSourceLocation(ssc *SscClient, args *Arguments) error {

	if len(args.Directory) == 0 {
		return fmt.Errorf("no source path specified" )
	}

	name := ValueOrDefault(args.Share, TEST_SOURCE_NAME)

	// does the location already exist?
	_, resp, err := getStorageLocation(ssc, name)
	if err == nil && resp.StatusCode < 300 {
		fmt.Printf("Source storage location %s already exists\n", name)
		return nil
	}


	timestamp := string(time.Now().Format("06-01-02-15-04-05"))
	description := fmt.Sprintf("%s, Created by Verify Test %s", name, timestamp)
	locationType := "NAS"
	locationIsTarget := false

	sourceDefinition := &openapi.ApiStorageLocationNas{
		Description:    	 &description,
		Path:				 &args.Directory,
		Type:				 &locationType,
		IsTarget:			 &locationIsTarget,
	}

	_, _, err = ssc.Client.StorageApi.UpdateNasStorageLocation(*ssc.Context, name, *sourceDefinition)
	if err != nil {
		return fmt.Errorf("failed to create source (%s) %v\n", name, err)
	}
	fmt.Printf("Successfully created source %s\n", name)
	return nil
}

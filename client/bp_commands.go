package client

import (
    "encoding/base64"
    "fmt"
    "github.com/SpectraLogic/ds3_go_sdk/ds3"
    "github.com/SpectraLogic/ds3_go_sdk/ds3/buildclient"
    "github.com/SpectraLogic/ds3_go_sdk/ds3/models"
    "github.com/SpectraLogic/ds3_go_sdk/ds3_cli/commands"
    "strconv"
    "strings"
    "time"
)

const MAX_POLLING_INTERVAL = 55

func makeSDKArgs(cli_args *Arguments) *commands.Arguments {
    var ret *commands.Arguments
    ret = new(commands.Arguments)
    ret.Endpoint =cli_args.Endpoint
    ret.Proxy = cli_args.Proxy
    ret.AccessKey = cli_args.AccessKey
    ret.SecretKey = cli_args.SecretKey
    return ret
}


func makeDs3Client (args *Arguments) (*ds3.Client, error) {
    // Build the client.
    client, err := buildclient.FromArgs(makeSDKArgs(args))
    if err != nil {
        return nil, fmt.Errorf("could not create BlackPearl client\n%v\n", err)
    }
    return client, nil
}

func headObject(ssc *SscClient, args *Arguments) error {
    // Validate the arguments.
    if args.Bucket == "" {
        return  fmt.Errorf("must specify a bucket")
    }
    bucket := args.Bucket
    key := args.FileName
    ds3Client, err := makeDs3Client(args)
    if err != nil {
        return  fmt.Errorf("cannot create client\n%v\n", err)
    }

    if key == "" {
        // do whole bucket
        request:= models.NewGetBucketRequest(bucket)
        response, err := ds3Client.GetBucket(request)
        if err != nil {
            return fmt.Errorf("Failed to get bucket %s\n%v", bucket, err)
        }

        for _, object := range response.ListBucketResult.Objects {
            err = processHeadObject(ds3Client, bucket, *object.Key)
            if err != nil {
                return fmt.Errorf("Failed head_object for Bucket %s, key %s\n%v",
                                    bucket, *object.Key, err)
            }
        }
        return nil
    } else {
        // single object
        return processHeadObject(ds3Client, bucket, key)
    }
}

func listBucketContents(ssc *SscClient, args *Arguments) error {
    bucket := ValueOrDefault(args.Bucket, TEST_BUCKET_NAME)
    ds3Client, err := makeDs3Client(args)
    if err != nil {
        return fmt.Errorf("cannot create client\n%v\n", err)
    }
    objects, err := processListBucket(ds3Client, bucket)
    for _, object := range objects {
        fmt.Println(*object.Key)
    }
    return err
}

func getTestFileFullPath(args *Arguments) (string, error) {
    fmt.Printf("Waiting for placement on cache .")
    return waitForCachePlacement(args, 0, 1)
}

func waitForCachePlacement(args *Arguments, fib1 int, fib2 int) (string, error) {
    bucket := ValueOrDefault(args.Bucket, TEST_BUCKET_NAME)
    fileName := ValueOrDefault(args.FileName, TEST_SOURCE_FILE)
    ds3Client, err := makeDs3Client(args)
    if err != nil {
        return "", fmt.Errorf("cannot create client\n%v\n", err)
    }
    objects, err := processListBucket(ds3Client, bucket)
    for _, object := range objects {
        if strings.Index(*object.Key, fileName) >= 0 {
            fmt.Printf("\nTarget file %s\n", *object.Key)
            return *object.Key, nil
        }
    }
    fmt.Printf(".")
    // wait
    fib3 := staggeredWait(fib1, fib2)
    // try again
    return waitForCachePlacement(args, fib2, fib3)
}

func processListBucket(client *ds3.Client, bucket string) ([]models.Contents, error) {
    request:= models.NewGetBucketRequest(bucket)
    response, err := client.GetBucket(request)
    if err != nil {
        return nil, fmt.Errorf("Failed to get bucket %s\n%v", bucket, err)
    }
    return response.ListBucketResult.Objects, nil
}

func processHeadObject(client *ds3.Client, bucket string, key string) error {
    // Build the request.
    request := models.NewHeadObjectRequest(bucket, key)

    // Perform the request.
    response, requestErr := client.HeadObject(request)
    if requestErr != nil {
        return requestErr
    }

    checksums := response.BlobChecksums
    checksumType := response.BlobChecksumType

    fmt.Printf("Bucket: %s\n", bucket)
    fmt.Printf("Key: %s\n", key)
    fmt.Printf("Type: %s\n", checksumType)
    for i, c := range checksums {
        fmt.Printf("%s: %s\n",strconv.Itoa(int(i)), c)
        raw, err := base64.StdEncoding.DecodeString(c)
        if err == nil {
            fmt.Printf("Decode: %s\n", string(raw))
        }
    }
    fmt.Printf("-----------------------------\n")

    return nil
}

func getPhysicalPlacement(ssc *SscClient, args *Arguments) error {
    // Validate the arguments.
    if args.Bucket == "" {
        return  fmt.Errorf("must specify a bucket")
    }
    if args.FileName == "" {
        return  fmt.Errorf("must specify a file anme")
    }

    ds3Client, err := makeDs3Client(args)
    if err != nil {
        return  fmt.Errorf("cannot create client\n%v\n", err)
    }
    _, err = processPhysicalPlacement(ds3Client, args.Bucket, args.FileName)
    return err
}

func waitForPlacement(ssc *SscClient, args *Arguments) error {
    // Validate the arguments.
    bucket := ValueOrDefault(args.Bucket, TEST_BUCKET_NAME)
    fileName := ValueOrDefault(args.FileName, TEST_SOURCE_FILE)

    ds3Client, err := makeDs3Client(args)
    if err != nil {
        return fmt.Errorf("cannot create client\n%v\n", err)
    }
    return tryPlacement(ds3Client, bucket, fileName, 0, 1)
}

func tryPlacement(client *ds3.Client, bucket string, key string, fib1 int, fib2 int) error {
    isPlaced, err := processPhysicalPlacement(client, bucket, key)
    if err != nil {
        // 404 is expected
        fmt.Printf(".")
    }
    if isPlaced {
        return nil
    }

    // wait
    fib3 := staggeredWait(fib1, fib2)
    // try again
    return tryPlacement(client, bucket, key, fib2, fib3)
}

func staggeredWait(fib1 int, fib2 int) int {
    // fibonacci up to maxInterval
    fib3 := fib1 + fib2
    if fib3 > MAX_POLLING_INTERVAL {
        fib3 = MAX_POLLING_INTERVAL
    }
    time.Sleep(time.Duration(fib3) * time.Second)
    return fib3
}

func processPhysicalPlacement(client *ds3.Client, bucket string, key string) (bool, error) {
    // Build the request.
    names := []string{key}
    request := models.NewVerifyPhysicalPlacementForObjectsWithFullDetailsSpectraS3Request(bucket, names)

    // Perform the request.
    response, requestErr := client.VerifyPhysicalPlacementForObjectsWithFullDetailsSpectraS3(request)
    if requestErr != nil {
        return false, requestErr
    }

    isPlaced := false
    placements := response.BulkObjectList.Objects

    for  _, placement:= range placements {
        for _, tape := range placement.PhysicalPlacement.Tapes {
            isPlaced = true
            fmt.Printf("\nOn tape: %s, barcode: %s\n", tape.Id, *tape.BarCode)
        }
    }

    return isPlaced, nil
}

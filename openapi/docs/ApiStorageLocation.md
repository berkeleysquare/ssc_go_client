# ApiStorageLocation

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**AccessKeyId** | Pointer to **string** | For use only with BlackPearl, S3 or MSAzure.  The access key ID of the credentials for BlackPearl, S3 or MS Azure. | [optional] 
**Active** | Pointer to **bool** | Used to re-activate a &#39;deleted&#39; storage location. | [optional] 
**ADGroup** | Pointer to **string** | Sources can have an adGroup name associated with them to restrict access to users not in the group. | [optional] 
**SpectraDataEndpoint** | Pointer to **string** | For use only with BlackPearl.  The URL of the BlackPearl server data port. | [optional] 
**SpectraMgmtEndpoint** | Pointer to **string** | For use only with BlackPearl or SpectraNAS.  The URL of the BlackPearl server management port. | [optional] 
**SpectraPassword** | Pointer to **string** | For use only with BlackPearl or SpectraNAS.  The password to login to the management port.  This is write-only. | [optional] 
**SpectraUsername** | Pointer to **string** | For use only with BlackPearl or SpectraNAS.  The username to login to the management port. | [optional] 
**Bucket** | Pointer to **string** | Used for BlackPearl and S3/Archive endpoints.  The bucket/vault to use for storage. | [optional] 
**BucketPolicyId** | Pointer to **string** | Used for specifying existing BlackPearl data policies. | [optional] 
**BytesTransferredIn** | Pointer to **int64** | Total number of bytes transferred to/from this location. This field is read-only. | [optional] [readonly] 
**BytesTransferredOut** | Pointer to **int64** | Total number of bytes transferred to/from this location. This field is read-only. | [optional] [readonly] 
**Connected** | Pointer to **bool** | True if location is currently connected to the server.  This field is read-only. | [optional] [readonly] 
**CostPerTB** | Pointer to **float32** | Cost per TB of storage in units of department currency.  Only valid when Department is set. | [optional] 
**Created** | Pointer to [**time.Time**](time.Time.md) | Time when storage location was created.  This field is read-only. | [optional] [readonly] 
**Department** | Pointer to **string** | ID of department responsible for location storage cost. | [optional] 
**Description** | Pointer to **string** | Description of the storage location. | [optional] 
**EnableEncryption** | Pointer to **bool** | Used to enable encryption that has been configured for the product. | [optional] 
**EnableSnapshots** | Pointer to **bool** | For use only with SpectraNAS. Used to enable snapshots on the SpectraNAS device. | [optional] 
**FilesTransferredIn** | Pointer to **int64** | Total number of bytes transferred to/from this location. This field is read-only. | [optional] [readonly] 
**FilesTransferredOut** | Pointer to **int64** | Total number of bytes transferred to/from this location. This field is read-only. | [optional] [readonly] 
**IsTarget** | Pointer to **bool** | Prevent location from use as a source. Set to true for BlackPearl. | 
**KeepReadOnly** | Pointer to **bool** | For use only with SpectraNAS. Used to keep the SpectraNAS device in a Read-Only state after a job completes. | [optional] 
**LastAccess** | Pointer to [**time.Time**](time.Time.md) | Time when storage location was last accessed.  This field is read-only. | [optional] [readonly] 
**LastScan** | Pointer to [**time.Time**](time.Time.md) | Time when storage location was last scanned for files.  This field is read-only. | [optional] [readonly] 
**Name** | Pointer to **string** | Storage location name.  This is the location ID. | [optional] [readonly] 
**OffPeakBytesPerSecondLimit** | Pointer to **float64** | The limit during off-peak hours, expressed in total number of bytes per second.  Specify 0 for no limit. | [optional] 
**OffPeakFilesPerSecondScanLimit** | Pointer to **float64** | The file scan limit during off-peak hours, expressed in total number of files per second. Specify 0 for no limit. | [optional] [default to 0]
**PackingType** | Pointer to **string** | How archived data should be packed.  This cannot be set on NAS locations. Default is None. | [optional] 
**Path** | Pointer to **string** | Used for NAS source/target and S3 source locations, this is the filesystem path to the NAS mount point and optionally the prefix for S3 path matching. | [optional] 
**PeakBytesPerSecondLimit** | Pointer to **float64** | The limit during peak hours, expressed in total number of bytes per second.  Specify 0 for no limit. | [optional] [default to 0]
**PeakFilesPerSecondScanLimit** | Pointer to **float64** | The file scan limit during peak hours, expressed in total number of files per second. Specify 0 for no limit. | [optional] [default to 0]
**PeakTimesSchedule** | Pointer to [**[]ApiPeakTimeSchedule**](api.peakTimeSchedule.md) | A list of the days of the week where rate limits are enforced during peak and off-peak times.  If a day of week is not in the list, it is assumed to be unlimited. | [optional] 
**RetentionDays** | Pointer to **int32** | The number of days to retain files on this target before scheduling files for deletion. This field is only applicable to targets. | 
**S3Endpoint** | Pointer to **string** | For use only with Amazon S3 or S3 emulators.  The URL of the endpoint to which the server should connect to access the storage. | [optional] 
**S3Region** | Pointer to **string** | For use only with Amazon S3 or S3 emulators.  The S3 region in which to connect.  If empty, the default region is assumed. | [optional] 
**IgnoreSSL** | Pointer to **bool** | For use only with S3 emulators.  Will ignore SSL certificate errors (self signed certificates). | [optional] 
**SecretAccessKey** | Pointer to **string** | For use only with BlackPearl, S3 or MSAzure.  The secret access key of the credentials for BlackPearl, S3 or MS Azure. This is write-only. | [optional] 
**Type** | Pointer to **string** |  | 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)



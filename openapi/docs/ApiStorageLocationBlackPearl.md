# ApiStorageLocationBlackPearl

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**AccessKeyId** | Pointer to **string** | For use only with BlackPearl.  The access key ID of the credentials for BlackPearl. | 
**Active** | Pointer to **bool** | Used to re-activate a &#39;deleted&#39; storage location. | [optional] 
**ADGroup** | Pointer to **string** | Sources can have an adGroup name associated with them to restrict access to users not in the group. | [optional] 
**SpectraDataEndpoint** | Pointer to **string** | For use only with BlackPearl.  The URL of the BlackPearl server data port. | 
**SpectraMgmtEndpoint** | Pointer to **string** | For use only with BlackPearl.  The URL of the BlackPearl server management port. | 
**SpectraPassword** | Pointer to **string** | For use only with BlackPearl.  The password to login to the management port.  This is write-only. | 
**SpectraUsername** | Pointer to **string** | For use only with BlackPearl.  The username to login to the management port. | 
**Bucket** | Pointer to **string** | Used for BlackPearl endpoints.  The bucket/vault to use for storage. | 
**BucketPolicyId** | Pointer to **string** | Used for specifying existing BlackPearl data policies. | [optional] 
**CostPerTB** | Pointer to **float32** | Cost per TB of storage in units of department currency.  Only valid when Department is set. | [optional] 
**Department** | Pointer to **string** | ID of department responsible for location storage cost. | [optional] 
**Description** | Pointer to **string** | Description of the storage location. | [optional] 
**EnableEncryption** | Pointer to **bool** | Used to enable encryption that has been configured for the product. | [optional] 
**IsTarget** | Pointer to **bool** | Prevent location from use as a source. Set to true for BlackPearl. | 
**OffPeakBytesPerSecondLimit** | Pointer to **float64** | The limit during off-peak hours, expressed in total number of bytes per second.  Specify 0 for no limit. | [optional] [default to 0]
**OffPeakFilesPerSecondScanLimit** | Pointer to **float64** | The file scan limit during off-peak hours, expressed in total number of files per second. Specify 0 for no limit. | [optional] [default to 0]
**PackingType** | Pointer to **string** | How archived data should be packed. Default is None. | [optional] 
**PeakBytesPerSecondLimit** | Pointer to **float64** | The limit during peak hours, expressed in total number of bytes per second.  Specify 0 for no limit. | [optional] [default to 0]
**PeakFilesPerSecondScanLimit** | Pointer to **float64** | The file scan limit during peak hours, expressed in total number of files per second. Specify 0 for no limit. | [optional] [default to 0]
**PeakTimesSchedule** | Pointer to [**[]ApiPeakTimeSchedule**](api.peakTimeSchedule.md) | A list of the days of the week where rate limits are enforced during peak and off-peak times.  If a day of week is not in the list, it is assumed to be unlimited. | [optional] 
**RetentionDays** | Pointer to **int32** | The number of days to retain files on this target before scheduling files for deletion. This field is only applicable to targets. | [optional] 
**SecretAccessKey** | Pointer to **string** | For use only with BlackPearl. The secret access key of the credentials for BlackPearl. This is write-only. | 
**Type** | Pointer to **string** |  | 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)



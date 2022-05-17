# ApiStorageLocationNas

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Active** | Pointer to **bool** | Used to re-activate a &#39;deleted&#39; storage location. | [optional] 
**ADGroup** | Pointer to **string** | Sources can have an adGroup name associated with them to restrict access to users not in the group. | [optional] 
**CostPerTB** | Pointer to **float32** | Cost per TB of storage in units of department currency.  Only valid when Department is set. | [optional] 
**Department** | Pointer to **string** | ID of department responsible for location storage cost. | [optional] 
**Description** | Pointer to **string** | Description of the storage location. | [optional] 
**IsTarget** | Pointer to **bool** | Prevent location from use as a source. | 
**EnableEncryption** | Pointer to **bool** | Used to enable encryption that has been configured for the product. | [optional] 
**OffPeakBytesPerSecondLimit** | Pointer to **float64** | The limit during off-peak hours, expressed in total number of bytes per second.  Specify 0 for no limit. | [optional] [default to 0]
**OffPeakFilesPerSecondScanLimit** | Pointer to **float64** | The file scan limit during off-peak hours, expressed in total number of files per second. Specify 0 for no limit. | [optional] [default to 0]
**Path** | Pointer to **string** | Used for NAS locations, this is the filesystem path to the NAS mount point. | 
**PeakBytesPerSecondLimit** | Pointer to **float64** | The limit during peak hours, expressed in total number of bytes per second.  Specify 0 for no limit. | [optional] [default to 0]
**PeakFilesPerSecondScanLimit** | Pointer to **float64** | The file scan limit during peak hours, expressed in total number of files per second. Specify 0 for no limit. | [optional] [default to 0]
**PeakTimesSchedule** | Pointer to [**[]ApiPeakTimeSchedule**](api.peakTimeSchedule.md) | A list of the days of the week where rate limits are enforced during peak and off-peak times.  If a day of week is not in the list, it is assumed to be unlimited. | [optional] 
**RetentionDays** | Pointer to **int32** | The number of days to retain files on this target before scheduling files for deletion. This field is only applicable to targets. | [optional] 
**Type** | Pointer to **string** |  | 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)



# ApiProjectScan

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Description** | Pointer to **string** | Description of the project. | [optional] 
**Share** | Pointer to **string** | The source share the project will scan. This is not required when updating the project. | 
**WorkingDirectory** | Pointer to **string** | The directory on the share to scan.  Default is the share root. | [optional] 
**Tags** | Pointer to **[]string** | Tags to attach to jobs executed by this project. | [optional] 
**Active** | Pointer to **bool** | Used to re-activate a deactivated project. | [optional] 
**Enabled** | Pointer to **bool** | Used to pause/resume a project. | [optional] 
**Filter** | [**ApiProjectFilter**](api.project.Filter.md) |  | [optional] 
**Schedule** | [**ApiProjectSchedule**](api.project.Schedule.md) |  | 
**Status** | [**ApiProjectStatus**](api.project.Status.md) |  | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)



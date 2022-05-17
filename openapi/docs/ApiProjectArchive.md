# ApiProjectArchive

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Description** | Pointer to **string** | Description of the project. | [optional] 
**Share** | Pointer to **string** | The source share the project will archive files off of. This is not required when updating the project. | 
**WorkingDirectory** | Pointer to **string** | The directory on the share to archive.  Default is the share root. | [optional] 
**Tags** | Pointer to **[]string** | Tags to attach to jobs executed by this project. | [optional] 
**Active** | Pointer to **bool** | Used to re-activate a deactivated project. | [optional] 
**Enabled** | Pointer to **bool** | Used to pause/resume a project. | [optional] 
**ProjectType** | Pointer to **string** | The type of archive to perform. | 
**BreadCrumbAction** | Pointer to **string** | The action that is performed on the original source file after successful archive. | [optional] [default to BREAD_CRUMB_ACTION_KEEP_ORIGINAL]
**Targets** | Pointer to **[]string** | The storage endpoints where files will be archived. | 
**Filter** | [**ApiProjectFilter**](api.project.Filter.md) |  | [optional] 
**Schedule** | [**ApiProjectSchedule**](api.project.Schedule.md) |  | 
**Status** | [**ApiProjectStatus**](api.project.Status.md) |  | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)



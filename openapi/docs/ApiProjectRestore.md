# ApiProjectRestore

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Description** | Pointer to **string** | Description of the project. | [optional] 
**Share** | Pointer to **string** | The source share the project will restore files to. This is not required when updating the project. | 
**WorkingDirectory** | Pointer to **string** | The directory on the share where the files will be restored. Default is the share root. | [optional] 
**Tags** | Pointer to **[]string** | Tags to attach to jobs executed by this project. | [optional] 
**Active** | Pointer to **bool** | Used to re-activate a deactivated project. | [optional] 
**Enabled** | Pointer to **bool** | Used to pause/resume a project. | [optional] 
**ProjectType** | Pointer to **string** | The type of restore to perform. | 
**BreadCrumbAction** | Pointer to **string** | The type of bread crumb to create. This is only valid with RestoreBreadcrumb projects. | [optional] 
**EmailOnComplete** | Pointer to **[]string** | A list of users from which email addresses may be found for sending notifications. | [optional] 
**RestoreManifest** | Pointer to **string** | The name of the manifest containing the files to be restored. This is synonymous with the original archive job that originally archived the files. | 
**RestoreVersions** | Pointer to **[]string** | The list of versions from the manifest to restore. If empty, the entire manifest is restored. | [optional] 
**Schedule** | [**ApiProjectSchedule**](api.project.Schedule.md) |  | 
**Status** | [**ApiProjectStatus**](api.project.Status.md) |  | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)



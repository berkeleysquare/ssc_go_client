# ApiProject

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Description** | Pointer to **string** | Description of the project. | [optional] 
**Share** | Pointer to **string** | The source share with which the project will work. This is not required when updating the project. | 
**WorkingDirectory** | Pointer to **string** | For restore projects, the location within the share where the files are restored. For archive projects, the starting subdirectory on the share to archive. For scan projects, the directory on the share to scan.  Default is the share root. | [optional] 
**Tags** | Pointer to **[]string** | Tags to attach to jobs executed by this project. | [optional] 
**Active** | Pointer to **bool** | Used to re-activate a deactivated project. | [optional] 
**Enabled** | Pointer to **bool** | Used to pause/resume a project. | [optional] 
**ProjectType** | Pointer to **string** | The type of project. | [optional] 
**BreadCrumbAction** | Pointer to **string** | For archive projects, the action that is performed on the original source file after successful archive. For restore breadcrumb projects, this specifies the type of bread crumb to create. | [optional] 
**EmailOnComplete** | Pointer to **[]string** | A list of users from which email addresses may be found for sending notifications. | [optional] 
**RestoreManifest** | Pointer to **string** | The name of the manifest containing the files to be restored. This is synonymous with the original archive job that originally archived the files. | [optional] 
**RestoreVersions** | Pointer to **[]string** | The list of versions from the manifest to restore. | [optional] 
**Targets** | Pointer to **[]string** | The storage endpoints where files will be archived. | [optional] 
**Filter** | [**ApiProjectFilter**](api.project.Filter.md) |  | [optional] 
**Schedule** | [**ApiProjectSchedule**](api.project.Schedule.md) |  | 
**Status** | [**ApiProjectStatus**](api.project.Status.md) |  | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)



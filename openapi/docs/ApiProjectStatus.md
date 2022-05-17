# ApiProjectStatus

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | The name of the project. This cannot be modified. | [optional] [readonly] 
**Created** | Pointer to [**time.Time**](time.Time.md) | Time when the project was created.  This is set automatically and is read-only. | [optional] [readonly] 
**CreatedBy** | Pointer to **string** | The user that created the project.  This is set automatically and is read-only. | [optional] [readonly] 
**Updated** | Pointer to [**time.Time**](time.Time.md) | The time the project information was last updated.  This is set automatically and is read-only. | [optional] [readonly] 
**UpdatedBy** | Pointer to **string** | The user that last updated the project.  This is set automatically and is read-only. | [optional] [readonly] 
**LastJobRunTime** | Pointer to [**time.Time**](time.Time.md) | Time when a job last ran for this project.  This is set automatically and is read-only. | [optional] [readonly] 
**LastSuccessfulJob** | Pointer to [**time.Time**](time.Time.md) | End time of most recent successful job for this project.  This is set automatically and is read-only. | [optional] [readonly] 
**ConflictDetected** | Pointer to [**time.Time**](time.Time.md) | If the last attempt to run the project ran into a conflict with an existing job, this will be set. This is set automatically and is read-only. | [optional] [readonly] 
**NextRunTime** | Pointer to [**time.Time**](time.Time.md) | The next time this project is scheduled to run. This is set automatically and is read-only. | [optional] [readonly] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)



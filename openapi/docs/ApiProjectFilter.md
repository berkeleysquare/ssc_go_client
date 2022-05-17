# ApiProjectFilter

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ExcludeDirectories** | Pointer to **[]string** | Directories to exclude. | [optional] 
**ExcludeFileTypes** | Pointer to **[]string** | File types to exclude. | [optional] 
**IncludeDirectories** | Pointer to **[]string** | Directories to include.  Specifying include directories will exclude any directories not explicitly identified here. | [optional] 
**IncludeFileTypes** | Pointer to **[]string** | File types to include. Specifying include file types will exclude any file types not explicitly identified here. | [optional] 
**MinimumAge** | Pointer to **string** | The minimum age of files with which this project will work.  Default is AnyAge. | [optional] 
**CustomAgeInDays** | Pointer to **int32** | The number of days old a file must equal or exceed for it to be processed by the job. This is only valid if CustomAge was selected for minimumAge. | [optional] 
**MinimumSize** | Pointer to **string** | The minimum size of files with which this project will work. Default is Any. | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)



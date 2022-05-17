# \ProjectApi

All URIs are relative to *https://localhost/openapi*

Method | HTTP request | Description
------------- | ------------- | -------------
[**GetProject**](ProjectApi.md#GetProject) | **Get** /projects/{projectName} | Retrieve a project
[**ListProjects**](ProjectApi.md#ListProjects) | **Get** /projects | List projects.
[**UpdateArchiveProject**](ProjectApi.md#UpdateArchiveProject) | **Put** /projects/archive/{projectName} | Create or update an archive project.
[**UpdateRestoreProject**](ProjectApi.md#UpdateRestoreProject) | **Put** /projects/restore/{projectName} | Create or update an restore project.
[**UpdateScanProject**](ProjectApi.md#UpdateScanProject) | **Put** /projects/scan/{projectName} | Create or update a scan project.



## GetProject

> ApiProject GetProject(ctx, projectName)

Retrieve a project

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**projectName** | **string**| Name of the project to retrieve | 

### Return type

[**ApiProject**](api.project.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListProjects

> ApiProjectPaginator ListProjects(ctx, optional)

List projects.

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***ListProjectsOpts** | optional parameters | nil if no parameters

### Optional Parameters

Optional parameters are passed through a pointer to a ListProjectsOpts struct


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **skip** | **optional.Int64**| Number of projects to skip before returning subsequent results. | [default to 0]
 **limit** | **optional.Int64**| Maximum number of projects to return. | [default to 500]
 **active** | **optional.Bool**| Set to true to list only the active projects. | 
 **sortBy** | **optional.String**| Specify which item to sort by. Default is to sort by created in descending order. Options: created|name | 
 **sortType** | **optional.String**| Specify sorting by project name as either ascending or descending. Default ascending. Options: ASC|DESC | 
 **filterBy** | [**optional.Interface of []string**](string.md)| Specify the type of job to list. Can specify this parameter multiple times with different values. If not specified, all job types are retrieved. Options: Scan|Archive|ScanAndArchive|Restore|RestoreBreadcrumb|BackupDb|ProjectDelete | 

### Return type

[**ApiProjectPaginator**](api.ProjectPaginator.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateArchiveProject

> ApiProjectArchive UpdateArchiveProject(ctx, projectName, body)

Create or update an archive project.

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**projectName** | **string**| Name of archive project to create or update | 
**body** | [**ApiProjectArchive**](ApiProjectArchive.md)|  | 

### Return type

[**ApiProjectArchive**](api.project.Archive.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateRestoreProject

> ApiProjectRestore UpdateRestoreProject(ctx, projectName, body)

Create or update an restore project.

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**projectName** | **string**| Name of restore project to create or update | 
**body** | [**ApiProjectRestore**](ApiProjectRestore.md)|  | 

### Return type

[**ApiProjectRestore**](api.project.Restore.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateScanProject

> ApiProjectScan UpdateScanProject(ctx, projectName, body)

Create or update a scan project.

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**projectName** | **string**| Name of scan project to create or update | 
**body** | [**ApiProjectScan**](ApiProjectScan.md)|  | 

### Return type

[**ApiProjectScan**](api.project.Scan.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


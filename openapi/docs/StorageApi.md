# \StorageApi

All URIs are relative to *https://localhost/openapi*

Method | HTTP request | Description
------------- | ------------- | -------------
[**GetStorageLocation**](StorageApi.md#GetStorageLocation) | **Get** /storage-locations/{storageLocationName} | Retrieve a storage location
[**ListStorageLocations**](StorageApi.md#ListStorageLocations) | **Get** /storage-locations/ | List all storage locations.
[**UpdateStorageLocation**](StorageApi.md#UpdateStorageLocation) | **Put** /storage-locations/{storageLocationName} | Create or update a storage location.



## GetStorageLocation

> ApiStorageLocation GetStorageLocation(ctx, storageLocationName)

Retrieve a storage location

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**storageLocationName** | **string**| Name of storage location to retrieve | 

### Return type

[**ApiStorageLocation**](api.storageLocation.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListStorageLocations

> ApiStorageLocationPaginator ListStorageLocations(ctx, optional)

List all storage locations.

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***ListStorageLocationsOpts** | optional parameters | nil if no parameters

### Optional Parameters

Optional parameters are passed through a pointer to a ListStorageLocationsOpts struct


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **skip** | **optional.Int64**| Number of storage locations to skip before returning subsequent results. | [default to 0]
 **limit** | **optional.Int64**| Maximum number of storage locations to return. | [default to 500]
 **includeAll** | **optional.Bool**| Whether to include storage locations that have been deactivated. | [default to false]
 **sortBy** | **optional.String**| Specify which item to sort by. Default is to sort by created in descending order. Options: created|name|locationType|isTarget|path | 
 **sortType** | **optional.String**| Specify sorting by project name as either ascending or descending. Default ascending. Options: ASC|DESC | 

### Return type

[**ApiStorageLocationPaginator**](api.StorageLocationPaginator.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateStorageLocation

> ApiStorageLocation UpdateStorageLocation(ctx, storageLocationName, body, optional)

Create or update a storage location.

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**storageLocationName** | **string**| Name of storage location to create or update | 
**body** | [**UNKNOWN_BASE_TYPE**](UNKNOWN_BASE_TYPE.md)|  | 
 **optional** | ***UpdateStorageLocationOpts** | optional parameters | nil if no parameters

### Optional Parameters

Optional parameters are passed through a pointer to a UpdateStorageLocationOpts struct


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **cloneCredentials** | **optional.String**| The source BlackPearl to use to get authentication credentials for the new storage location. | 

### Return type

[**ApiStorageLocation**](api.storageLocation.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


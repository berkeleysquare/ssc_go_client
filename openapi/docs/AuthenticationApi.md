# \AuthenticationApi

All URIs are relative to *https://localhost/openapi*

Method | HTTP request | Description
------------- | ------------- | -------------
[**Login**](AuthenticationApi.md#Login) | **Post** /tokens | Authenticate and authorize user access.



## Login

> ApiToken Login(ctx, body)

Authenticate and authorize user access.

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**body** | [**ApiCredentials**](ApiCredentials.md)|  | 

### Return type

[**ApiToken**](api.Token.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


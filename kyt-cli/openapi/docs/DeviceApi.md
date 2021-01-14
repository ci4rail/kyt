# \DeviceApi

All URIs are relative to *https://kyt.ci4rail.com/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**GetDevices**](DeviceApi.md#GetDevices) | **Get** /devices | List devices for a tenant



## GetDevices

> []Device GetDevices(ctx).Execute()

List devices for a tenant



### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "./openapi"
)

func main() {

    configuration := openapiclient.NewConfiguration()
    api_client := openapiclient.NewAPIClient(configuration)
    resp, r, err := api_client.DeviceApi.GetDevices(context.Background()).Execute()
    if err.Error() != "" {
        fmt.Fprintf(os.Stderr, "Error when calling `DeviceApi.GetDevices``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `GetDevices`: []Device
    fmt.Fprintf(os.Stdout, "Response from `DeviceApi.GetDevices`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiGetDevicesRequest struct via the builder pattern


### Return type

[**[]Device**](Device.md)

### Authorization

[api_key](../README.md#api_key)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

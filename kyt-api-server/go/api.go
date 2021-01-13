/*
 * KYT API Server
 *
 * This is the KYT API Server
 *
 * API version: 0.0.1
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

import (
	"context"
	"net/http"
)

// DeviceApiRouter defines the required methods for binding the api requests to a responses for the DeviceApi
// The DeviceApiRouter implementation should parse necessary information from the http request,
// pass the data to a DeviceApiServicer to perform the required actions, then write the service results to the http response.
type DeviceApiRouter interface {
	GetDevices(http.ResponseWriter, *http.Request)
}

// UserApiRouter defines the required methods for binding the api requests to a responses for the UserApi
// The UserApiRouter implementation should parse necessary information from the http request,
// pass the data to a UserApiServicer to perform the required actions, then write the service results to the http response.
type UserApiRouter interface {
	LoginUser(http.ResponseWriter, *http.Request)
	LogoutUser(http.ResponseWriter, *http.Request)
}

// DeviceApiServicer defines the api actions for the DeviceApi service
// This interface intended to stay up to date with the openapi yaml used to generate it,
// while the service implementation can ignored with the .openapi-generator-ignore file
// and updated with the logic required for the API.
type DeviceApiServicer interface {
	GetDevices(context.Context) (ImplResponse, error)
}

// UserApiServicer defines the api actions for the UserApi service
// This interface intended to stay up to date with the openapi yaml used to generate it,
// while the service implementation can ignored with the .openapi-generator-ignore file
// and updated with the logic required for the API.
type UserApiServicer interface {
	LoginUser(context.Context, string) (ImplResponse, error)
	LogoutUser(context.Context) (ImplResponse, error)
}

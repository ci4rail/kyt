/*
 * KYT API Server
 *
 * This is the KYT API Server
 *
 * API version: 0.0.1
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package main

import (
	"log"
	"net/http"

	openapi "github.com/GIT_USER_ID/GIT_REPO_ID/go"
)

func main() {
	log.Printf("Server started")

	DeviceApiService := openapi.NewDeviceApiService()
	DeviceApiController := openapi.NewDeviceApiController(DeviceApiService)

	UserApiService := openapi.NewUserApiService()
	UserApiController := openapi.NewUserApiController(UserApiService)

	router := openapi.NewRouter(DeviceApiController, UserApiController)

	log.Fatal(http.ListenAndServe(":8080", router))
}

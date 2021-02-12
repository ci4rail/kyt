package api

const (
	userFlow        = "b2c_1_signin_native"
	azureB2CTenant  = "ci4railtesting"
	azureB2CKeysURI = "https://" + azureB2CTenant + ".b2clogin.com/" + azureB2CTenant + ".onmicrosoft.com/" + userFlow + "/discovery/v2.0/keys"
)

/*
Copyright © 2021 Ci4Rail GmbH

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package token

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	jwt "github.com/form3tech-oss/jwt-go"
)

// Response is used as json output
type Response struct {
	Message string `json:"message"`
}

type jwks struct {
	Keys []jsonWebKeys `json:"keys"`
}

type jsonWebKeys struct {
	Kty string   `json:"kty"`
	Kid string   `json:"kid"`
	Use string   `json:"use"`
	N   string   `json:"n"`
	E   string   `json:"e"`
	X5c []string `json:"x5c"`
}

type customClaims struct {
	Scope   string   `json:"scope"`
	Tenants []string `json:"http://tenants"`
	jwt.StandardClaims
}

// CheckScope validates if the scope are present in the token
func CheckScope(scope string, tokenString string) bool {
	token, _ := jwt.ParseWithClaims(tokenString, &customClaims{}, func(token *jwt.Token) (interface{}, error) {
		cert, err := getPemCert(token)
		if err != nil {
			return nil, err
		}
		result, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
		return result, nil
	})

	claims, ok := token.Claims.(*customClaims)

	hasScope := false
	if ok && token.Valid {
		result := strings.Split(claims.Scope, " ")
		for i := range result {
			if result[i] == scope {
				hasScope = true
			}
		}
	}

	return hasScope
}

// GetTenants gets all tenants from token
func GetTenants(tokenString string) []string {
	token, _ := jwt.ParseWithClaims(tokenString, &customClaims{}, func(token *jwt.Token) (interface{}, error) {
		cert, err := getPemCert(token)
		if err != nil {
			return nil, err
		}
		result, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
		return result, nil
	})

	claims, ok := token.Claims.(*customClaims)

	if ok && token.Valid {
		return claims.Tenants
	}
	return nil
}

func getPemCert(token *jwt.Token) (string, error) {
	cert := ""
	resp, err := http.Get("https://edgefarm-dev.eu.auth0.com/.well-known/jwks.json")

	if err != nil {
		return cert, err
	}
	defer resp.Body.Close()

	var jwks = jwks{}
	err = json.NewDecoder(resp.Body).Decode(&jwks)

	if err != nil {
		return cert, err
	}

	for k := range jwks.Keys {
		if token.Header["kid"] == jwks.Keys[k].Kid {
			cert = "-----BEGIN CERTIFICATE-----\n" + jwks.Keys[k].X5c[0] + "\n-----END CERTIFICATE-----"
		}
	}

	if cert == "" {
		err := errors.New("Unable to find appropriate key")
		return cert, err
	}

	return cert, nil
}

func validationKeyGetter(token *jwt.Token) (interface{}, error) {
	// Verify 'iss' claim
	iss := "https://edgefarm-dev.eu.auth0.com/"
	checkIss := token.Claims.(jwt.MapClaims).VerifyIssuer(iss, false)
	if !checkIss {
		return token, errors.New("Invalid issuer")
	}

	cert, err := getPemCert(token)
	if err != nil {
		panic(err.Error())
	}

	result, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
	return result, nil
}

// CreateMiddleware create the middleware for Auth0 token handling
func CreateMiddleware(domain string) *jwtmiddleware.JWTMiddleware {
	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: validationKeyGetter,
		SigningMethod:       jwt.SigningMethodRS256,
	})
	return jwtMiddleware
}

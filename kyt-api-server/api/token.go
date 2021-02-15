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

package api

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
)

// downloads the current signing keys
// TODO: this can be cached and redownloaded once a token signature
// validation failed. For now it always downloads on every request.
func readFlowKey() (string, error) {
	// read all signing keys from Azure B2C for specific User Flow
	response, err := http.Get(azureB2CKeysURI)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	return string(responseData), nil
}

// keyArray and singleKey are used for unmarshalling the signing key response
type singleKey struct {
	Kid string `json:"kid"`
	Nbf int    `json:"nbf"`
	Use string `json:"use"`
	Kty string `json:"kty"`
	E   string `json:"e"`
	N   string `json:"n"`
}

type keyArray struct {
	Keys []singleKey `json:"keys"`
}

// convert n, e format to PublicKey
func tokenValidateConvert(js string) (*rsa.PublicKey, error) {
	k := &keyArray{}
	err := json.Unmarshal([]byte(js), &k)
	if err != nil {
		return &rsa.PublicKey{}, err
	}

	jwk := k.Keys[0]
	if jwk.Kty != "RSA" {
		log.Fatal("invalid key type:", jwk.Kty)
	}

	// decode the base64 bytes for n
	nb, err := base64.RawURLEncoding.DecodeString(jwk.N)
	if err != nil {
		return &rsa.PublicKey{}, err
	}

	e := 0
	// The default exponent is usually 65537, so just compare the
	// base64 for [1,0,1] or [0,1,0,1]
	if jwk.E == "AQAB" || jwk.E == "AAEAAQ" {
		e = 65537
	} else {
		// need to decode "e" as a big-endian int
		return &rsa.PublicKey{}, fmt.Errorf("need to deocde e:", jwk.E)
	}

	pk := &rsa.PublicKey{
		N: new(big.Int).SetBytes(nb),
		E: e,
	}
	return pk, nil
}

// extractor function for token signature validation.
// It looks up the current signing keys and returns the PublicKey for
// token signature validation
func tokenExtractor(token *jwt.Token) (interface{}, error) {
	keys, err := readFlowKey()
	if err != nil {
		return nil, err
	}
	cert, err := tokenValidateConvert(keys)
	if err != nil {
		return nil, err
	}
	return cert, nil
}

type claimsType struct {
	*jwt.StandardClaims
	Scopes string `json:"scp,omitempty"`
}

// checks if the token that came in with a request is valid.
// valid means:
// - token is not expired.
// - signature is validated to ensure that the token hasnt been changed since it was issued.
func validateToken(r *http.Request) (bool, []string, error) {
	token, err := request.ParseFromRequestWithClaims(r, request.OAuth2Extractor, &claimsType{}, tokenExtractor)
	if err != nil {
		fmt.Println(err)
		return false, nil, err
	}
	claims := tokenizeClaims(token.Claims.(*claimsType).Scopes)
	return token.Valid, claims, nil
}

func tokenizeClaims(claims string) []string {
	return strings.Split(claims, " ")
}

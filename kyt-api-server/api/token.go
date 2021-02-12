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
	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
)

func readFlowKey() (string, error) {
	url := "https://ci4railtesting.b2clogin.com/ci4railtesting.onmicrosoft.com/b2c_1_signin_native/discovery/v2.0/keys"
	response, err := http.Get(url)
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

func tokenValidateConvert(js string) (string, error) {
	k := &keyArray{}
	err := json.Unmarshal([]byte(js), &k)
	if err != nil {
		return "", err
	}

	jwk := k.Keys[0]
	if jwk.Kty != "RSA" {
		log.Fatal("invalid key type:", jwk.Kty)
	}

	// decode the base64 bytes for n
	nb, err := base64.RawURLEncoding.DecodeString(jwk.N)
	if err != nil {
		return "", err
	}

	e := 0
	// The default exponent is usually 65537, so just compare the
	// base64 for [1,0,1] or [0,1,0,1]
	if jwk.E == "AQAB" || jwk.E == "AAEAAQ" {
		e = 65537
	} else {
		// need to decode "e" as a big-endian int
		return "", fmt.Errorf("need to deocde e:", jwk.E)
	}

	pk := &rsa.PublicKey{
		N: new(big.Int).SetBytes(nb),
		E: e,
	}

	der, err := x509.MarshalPKIXPublicKey(pk)
	if err != nil {
		return "", err
	}

	block := &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: der,
	}

	var out bytes.Buffer
	err = pem.Encode(&out, block)
	if err != nil {
		return "", nil
	}
	fmt.Println(out.String())
	return out.String(), nil
}

func tokenExtractor(token *jwt.Token) (interface{}, error) {
	keys, err := readFlowKey()
	if err != nil {
		return nil, err
	}
	cert, err := tokenValidateConvert(keys)
	if err != nil {
		return nil, err
	}
	return []byte(cert), nil
}

// func validateToken(token string, endpoint string) (bool, error) {
func validateToken(r *http.Request) (bool, error) {
	token, err := request.ParseFromRequest(r, request.OAuth2Extractor, tokenExtractor)
	if err != nil {
		fmt.Println(err)
		return false, err
	}
	fmt.Println(token)
	return true, nil
	// var claims jwt.Claims = &jwt.MapClaims{}
	// jwtToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
	// 	return []byte("6OJ~NKqA17C27Lb9Zioark5te.5vPk__PZ"), nil
	// })
	// if err != nil {
	// 	fmt.Println(err)
	// 	return false, err
	// }
	// fmt.Println(jwtToken)
	// return true, nil
}

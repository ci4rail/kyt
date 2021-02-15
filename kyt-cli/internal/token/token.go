/*
Copyright © 2021 Ci4Rail GmbH <engineering@ci4rail.com>

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
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	common "github.com/ci4rail/kyt/kyt-cli/internal/common"
	configuration "github.com/ci4rail/kyt/kyt-cli/internal/configuration"
	e "github.com/ci4rail/kyt/kyt-cli/internal/errors"
	"github.com/dgrijalva/jwt-go"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

type accessTokenResponse struct {
	Type         string `json:"token_type"`
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	Scope        string `json:"scope"`
	IdToken      string `json:"id_token"`
	RefreshToken string `json:"refresh_token"`
}

func getScopes() string {
	scopes := "offline_access "
	for _, v := range configuration.GetConstScopes() {
		scopes += v + " "
	}
	return scopes
}

// CreateAccessTokenRequest creates the http request to obtain an access token
func CreateAccessTokenRequest(host string, cid string, uid string, upw string) (*http.Request, error) {
	data := url.Values{}
	data.Add("grant_type", "password")
	data.Add("username", uid)
	data.Add("password", upw)
	data.Add("client_id", cid)
	data.Add("scope", getScopes())
	req, err := http.NewRequest("POST", host, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(cid, "")
	return req, nil
}

// CreateRefreshTokenRequest creates the http request to obtain an refresh token
func CreateRefreshTokenRequest(host string, cid string, uid string, upw string) (*http.Request, error) {
	data := url.Values{}
	data.Add("grant_type", "refresh_token")
	data.Add("client_id", cid)
	data.Add("scope", getScopes())
	data.Add("refresh_token", viper.GetString("refresh_token"))

	req, err := http.NewRequest("POST", host, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(cid, "")
	return req, nil
}

// SendAccessTokenRequest sends the access token request
func SendAccessTokenRequest(req *http.Request) ([]byte, error) {
	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	if res.StatusCode == 400 {
		return nil, fmt.Errorf("invalid username or password\n")
	}

	if res.StatusCode != 200 {
		fmt.Fprintf(os.Stderr, "Error response from token endpoint (HTTP Status %d):\n", res.StatusCode)
		fmt.Fprintln(os.Stderr, string(body))
		return nil, err
	}

	return body, nil
}

// SendRefreshTokenRequest sends the refresh token request
func SendRefreshTokenRequest(req *http.Request) ([]byte, error) {
	res, err := http.DefaultClient.Do(req)

	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		panic(err)
	}

	// if res.StatusCode == 400 {
	// 	return nil, fmt.Errorf("invalid username or password\n")
	// }

	if res.StatusCode != 200 {
		fmt.Fprintf(os.Stderr, "Error response from token endpoint (HTTP Status %d):\n", res.StatusCode)
		fmt.Fprintln(os.Stderr, string(body))
		return nil, err
	}

	return body, nil
}

//GetTokenClaims retrieves the claims from a token
func GetTokenClaims(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, nil)
	if token == nil {
		return nil, err
	}
	claims, _ := token.Claims.(jwt.MapClaims)
	return claims, nil
}

// ExtractToken extracts the access token and the refresh token from the http response body
func ExtractToken(body []byte) (string, string, error) {
	// This intermediate step is needed, because `expires_in` is one time returned string and
	// the other time as int from:
	// `grant_type` == `password` and `token_refresh`
	raw := struct {
		Type         string      `json:"token_type"`
		AccessToken  string      `json:"access_token"`
		ExpiresIn    interface{} `json:"expires_in"`
		Scope        string      `json:"scope"`
		IdToken      string      `json:"id_token"`
		RefreshToken string      `json:"refresh_token"`
	}{}

	err := json.Unmarshal(body, &raw)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to parse response: \"")
		fmt.Fprintf(os.Stderr, "%s", err)
		fmt.Fprintf(os.Stderr, "\"\n")
		return "", "", err
	}
	atr := &accessTokenResponse{
		Type:         raw.Type,
		AccessToken:  raw.AccessToken,
		Scope:        raw.Scope,
		IdToken:      raw.IdToken,
		RefreshToken: raw.RefreshToken,
	}

	// Populate ExpiresIn by converting the value into an int
	// depending on the type of the value received
	switch v := raw.ExpiresIn.(type) {
	case int:
		atr.ExpiresIn = v
	case string:
		i, err := strconv.Atoi(v)
		if err != nil {
			e.Er(err)
		}
		atr.ExpiresIn = i
	}

	return atr.AccessToken, atr.RefreshToken, nil
}

// RefreshToken handles all stuff that is needed for refreshing the access token
func RefreshToken() error {
	req, err := CreateRefreshTokenRequest(configuration.TokenEndpoint, configuration.ClientId, common.Username, common.Password)
	if err != nil {
		e.Er(err)
	}
	resp, err := SendRefreshTokenRequest(req)
	if err != nil {
		e.Er(err)
	}
	token, refreshToken, err := ExtractToken(resp)
	if err != nil {
		e.Er(err)
	}
	WriteTokensToConfig(token, refreshToken)
	return nil
}

// WriteTokensToConfig stores access token and refresh token in the config file for later usage
func WriteTokensToConfig(token, refreshToken string) {
	viper.Set("token", token)
	viper.Set("refresh_token", refreshToken)

	home, err := homedir.Dir()
	if err != nil {
		e.Er(err)
	}
	err = viper.WriteConfigAs(fmt.Sprintf("%s/%s.%s", home, common.KytCliConfigFile, common.KytCliConfigFileType))
	if err != nil {
		e.Er(err)
	}
}

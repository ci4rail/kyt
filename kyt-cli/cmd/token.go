package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

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

func createAccessTokenRequest(host string, cid string, uid string, upw string) (*http.Request, error) {
	data := url.Values{}
	scope := fmt.Sprintf("%s %s", viper.GetString("scope"), "offline_access")
	data.Add("grant_type", "password")
	data.Add("username", uid)
	data.Add("password", upw)
	data.Add("client_id", cid)
	data.Add("scope", scope)
	req, err := http.NewRequest("POST", host, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(cid, "")
	return req, nil
}

func createRefreshTokenRequest(host string, cid string, uid string, upw string) (*http.Request, error) {
	data := url.Values{}
	scope := fmt.Sprintf("%s %s", viper.GetString("scope"), "offline_access")
	data.Add("grant_type", "refresh_token")
	data.Add("client_id", cid)
	data.Add("scope", scope)
	data.Add("refresh_token", viper.GetString("refresh_token"))

	req, err := http.NewRequest("POST", host, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(cid, "")
	return req, nil
}

func sendAccessTokenRequest(req *http.Request) ([]byte, error) {
	res, err := http.DefaultClient.Do(req)

	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		panic(err)
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

func sendRefreshTokenRequest(req *http.Request) ([]byte, error) {
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

func getTokenClaims(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, nil)
	if token == nil {
		return nil, err
	}
	claims, _ := token.Claims.(jwt.MapClaims)
	return claims, nil
}

func extractAccessToken(body []byte) (string, string, error) {
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
			er(err)
		}
		atr.ExpiresIn = i
	}

	return atr.AccessToken, atr.RefreshToken, nil
}

func RefreshToken() error {
	req, err := createRefreshTokenRequest(viper.GetString("token_endpoint"), viper.GetString("client_id"), username, password)
	if err != nil {
		er(err)
	}
	resp, err := sendRefreshTokenRequest(req)
	if err != nil {
		er(err)
	}
	token, refreshToken, err := extractAccessToken(resp)
	if err != nil {
		er(err)
	}
	writeTokensToConfig(token, refreshToken)
	return nil
}

func writeTokensToConfig(token, refreshToken string) {
	viper.Set("token", token)
	viper.Set("refresh_token", refreshToken)

	home, err := homedir.Dir()
	if err != nil {
		er(err)
	}
	err = viper.WriteConfigAs(fmt.Sprintf("%s/%s.%s", home, kytCliConfigFile, kytCliConfigFileType))
	if err != nil {
		er(err)
	}
}

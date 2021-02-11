package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
)

type accessTokenResponse struct {
	AccessToken  string `json:"access_token"`
	Type         string `json:"token_type"`
	ExpiresIn    string `json:"expires_in"`
	ExtExpiresIn string `json:"ext_expires_in"`
	ExpiresOn    string `json:"expires_on"`
	NotBefore    string `json:"not_before"`
	Resource     string `json:"resource"`
}

func createTokenRequest(host string, cid string, uid string, upw string) (*http.Request, error) {
	data := url.Values{}
	data.Add("grant_type", "password")
	data.Add("username", uid)
	data.Add("password", upw)
	data.Add("client_id", cid)
	data.Add("scope", viper.GetString("scope"))
	req, err := http.NewRequest("POST", host, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(cid, "")
	return req, nil
}

func createRrefreshTokenRequest(host string, cid string, uid string, upw string) (*http.Request, error) {
	data := url.Values{}
	data.Add("grant_type", "authorization_code")
	data.Add("client_id", cid)
	data.Add("username", uid)
	data.Add("password", upw)
	data.Add("scope", "offline_access")
	data.Add("code", viper.GetString("token"))
	req, err := http.NewRequest("POST", host, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(cid, "")
	return req, nil
}

func sendTokenRequest(req *http.Request) ([]byte, error) {
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

func extractAccessToken(body []byte) (string, error) {
	var atr accessTokenResponse
	err := json.Unmarshal(body, &atr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to parse response: \"")
		fmt.Fprintf(os.Stderr, "%s", err)
		fmt.Fprintf(os.Stderr, "\"\n")
		return "", err
	}
	return atr.AccessToken, nil
}

func RefreshToken() error {
	req, err := createRrefreshTokenRequest(viper.GetString("token_endpoint"), viper.GetString("client_id"), username, password)
	if err != nil {
		er(err)
	}
	resp, err := sendTokenRequest(req)
	if err != nil {
		er(err)
	}
	token, err := extractAccessToken(resp)
	if err != nil {
		er(err)
	}
	fmt.Printf("Refreshed token: %s\n", token)

	// apiClient, ctx := api.NewAPIWithToken(serverURL, viper.GetString("token"))
	// fmt.Println("Token expired. Refreshing...")
	// resp, openapierr := apiClient.AuthApi.AuthRefreshTokenGet(ctx).Execute()
	// if resp.StatusCode == 401 {
	// 	return fmt.Errorf("Unable to refresh access token. Please run `login` command again.")

	// } else if openapierr.Error() != "" {
	// 	er(fmt.Sprintf("Error calling RefreshApi.RefreshToken: %v\n", openapierr))
	// }

	// var data map[string]interface{}
	// err := json.NewDecoder(resp.Body).Decode(&data)
	// if err != nil {
	// 	return fmt.Errorf("Error: %e", err)
	// }

	// token := data["token"]
	// viper.Set("token", token)
	// // Find home directory.
	// home, err := homedir.Dir()
	// if err != nil {
	// 	er(err)
	// }
	// err = viper.WriteConfigAs(fmt.Sprintf("%s/%s.%s", home, kytCliConfigFile, kytCliConfigFileType))
	// if err != nil {
	// 	log.Println("Cannot save config file")
	// }
	return nil
}

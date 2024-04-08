package usecase

import (
	"bytes"
	"doc-translate-go/pkg/user/entity"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

var (
	authTokenEndpoint      = os.Getenv("AUTH_TOKEN_DOMAIN")
	authEndpoint           = os.Getenv("AUTH_DOMAIN")
	authIntrospectEndpoing = os.Getenv("AUTH_INTROSPECT_DOMAIN")
	authUserInfoEndpoint   = os.Getenv("AUTH_USERINFO_DOMAIN")
	clientId               = os.Getenv("CLIENT_ID")
	clientSecret           = os.Getenv("CLIENT_SECRET")
)

func newAuthHeader() string {
	cred := fmt.Sprintf("%s:%s", clientId, clientSecret)
	base64Cred := base64.StdEncoding.EncodeToString([]byte(cred))
	authzHeader := fmt.Sprintf("Basic %s", base64Cred)
	return authzHeader
}

func (uc *UserUseCase) RetrieveAccessToken(grantType string, token string) (map[string]any, error) {
	var codeType string

	switch grantType {
	case "refresh_token":
		codeType = "refresh_token"
	case "authorization_code":
		codeType = "code"
	default:
		return nil, fmt.Errorf("invalid grant type %s", grantType)
	}

	client := &http.Client{}
	data := url.Values{}
	data.Set("grant_type", grantType)
	data.Set(codeType, token)

	req, err := http.NewRequest("POST", authTokenEndpoint, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", newAuthHeader())

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var accessToken map[string]any
	if err := json.Unmarshal(body, &accessToken); err != nil {
		return nil, err
	}
	if _, exists := accessToken["error"]; exists {
		return nil, errors.New("failed to unmarshal")
	}

	return accessToken, nil
}

func Authorize() (string, error) {
	params := url.Values{}
	params.Add("response_type", "code")
	params.Add("client_id", clientId)
	params.Add("login_method", "form")

	endpoint := authEndpoint + "?" + params.Encode()

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Get(endpoint)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	redirectUrl := resp.Header.Get("Location")
	if redirectUrl == "" {
		return "", errors.New("location header not found")
	}

	return redirectUrl, nil
}

func IntrospectAccessToken(accessToken string) (string, error) {
	data := url.Values{}
	data.Set("token", accessToken)
	data.Set("token_type_hint", "access_token")

	req, err := http.NewRequest("POST", authIntrospectEndpoing, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return "", err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", newAuthHeader())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var introspect map[string]any
	if err := json.Unmarshal(body, &introspect); err != nil {
		return "", err
	}

	username, usernameExists := introspect["username"].(string)
	active, activeExists := introspect["active"].(bool)
	if !active || !usernameExists || !activeExists {
		return "", errors.New("invalid token")
	}

	return username, nil
}

func (uc *UserUseCase) RetrieveUserProfile(accessToken string) (*entity.UserProfile, error) {
	_, err := IntrospectAccessToken(accessToken)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", authUserInfoEndpoint, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var userProfile *entity.UserProfile
	if err := json.Unmarshal(body, userProfile); err != nil {
		return nil, err
	}

	return userProfile, nil
}

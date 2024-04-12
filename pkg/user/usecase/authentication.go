package usecase

import (
	"bytes"
	"doc-translate-go/pkg/config"
	"doc-translate-go/pkg/user/entity"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

func newAuthHeader(clientId string, clientSecret string) string {
	cred := fmt.Sprintf("%s:%s", clientId, clientSecret)
	base64Cred := base64.StdEncoding.EncodeToString([]byte(cred))
	authzHeader := fmt.Sprintf("Basic %s", base64Cred)
	return authzHeader
}

type AuthUseCase struct {
	authConfig *config.AuthConfig
}

func NewAuthUseCase(authConfig *config.AuthConfig) *AuthUseCase {
	return &AuthUseCase{authConfig}
}

func (uc *AuthUseCase) RetrieveAccessToken(grantType string, token string) (string, error) {
	var codeType string

	switch grantType {
	case "refresh_token":
		codeType = "refresh_token"
	case "authorization_code":
		codeType = "code"
	default:
		return "", fmt.Errorf("invalid grant type %s", grantType)
	}

	client := &http.Client{}
	data := url.Values{}
	data.Set("grant_type", grantType)
	data.Set(codeType, token)

	req, err := http.NewRequest("POST", uc.authConfig.TokenEndpoint, strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", newAuthHeader(uc.authConfig.ClientId, uc.authConfig.ClientSecret))

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var bod map[string]any
	if err := json.Unmarshal(body, &bod); err != nil {
		return "", err
	}
	if _, exists := bod["error"]; exists {
		return "", errors.New("failed to unmarshal")
	}

	accessToken, ok := bod["access_token"].(string)
	if !ok {
		return "", errors.New("failed to get access token")
	}

	return accessToken, nil
}

func (uc *AuthUseCase) Authorize() (string, error) {
	params := url.Values{}
	params.Add("response_type", "code")
	params.Add("client_id", uc.authConfig.ClientId)
	params.Add("login_method", "form")

	endpoint := uc.authConfig.AuthEndpoint + "?" + params.Encode()

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

func (uc *AuthUseCase) IntrospectAccessToken(accessToken string) (string, error) {
	data := url.Values{}
	data.Set("token", accessToken)
	data.Set("token_type_hint", "access_token")

	req, err := http.NewRequest("POST", uc.authConfig.IntrospectEndpoint, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return "", err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", newAuthHeader(uc.authConfig.ClientId, uc.authConfig.ClientSecret))

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

func (uc *AuthUseCase) RetrieveUserProfile(accessToken string) (*entity.UserProfile, error) {
	_, err := uc.IntrospectAccessToken(accessToken)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", uc.authConfig.UserInfoEndpoint, nil)
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

	var user *entity.UserProfile
	if err := json.Unmarshal(body, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (uc *AuthUseCase) ValidateDistributionListHasIsid(isid string) error {
	type Response struct {
		Count int `json:"count"`
	}

	var wg sync.WaitGroup
	found := false
	resultChan := make(chan bool, 1)
	semaphore := make(chan struct{}, 20)

	for _, dl := range uc.authConfig.DistributionList {
		wg.Add(1)

		go func(dl string) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			if found {
				return
			}

			url := fmt.Sprintf("%s/%s/members", uc.authConfig.DistributionListEndpoint, dl)
			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				return
			}

			query := req.URL.Query()
			query.Add("filter", fmt.Sprintf("isid=%s", isid))
			query.Add("includeNpa", "true")
			req.URL.RawQuery = query.Encode()

			req.Header.Add("Accept", "application/json")
			req.Header.Add("X-Merck-APIKey", uc.authConfig.ClientId)

			client := &http.Client{Timeout: 50 * time.Second}

			resp, err := client.Do(req)
			if err != nil {
				return
			}
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return
			}

			var directory Response
			if err := json.Unmarshal(body, &directory); err != nil {
				return
			}

			if directory.Count > 0 {
				select {
				case resultChan <- true:
				default:
				}
				found = true
			}
		}(dl)

		if found {
			break
		}
	}

	wg.Wait()
	close(resultChan)
	if !found {
		return errors.New("403: Forbidden")
	}

	return nil
}

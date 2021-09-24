package github

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/geisonbiazus/blog/internal/core/auth"
)

type HTTPClient interface {
	Get(url string) (resp *http.Response, err error)
}

type Client struct {
	httpClient HTTPClient
}

func NewClient(httpClient HTTPClient) *Client {
	return &Client{httpClient: httpClient}
}

func (c *Client) GetAuthenticatedUser() (auth.ProviderUser, error) {
	resp, err := c.requestCurrentUser()
	if err != nil {
		return auth.ProviderUser{}, err
	}

	defer resp.Body.Close()

	return c.parseResponse(resp)
}

func (c *Client) requestCurrentUser() (*http.Response, error) {
	resp, err := c.httpClient.Get("https://api.github.com/user")
	if err != nil {
		return nil, fmt.Errorf("error requesting user on github.Client: %w", err)
	}

	return resp, nil
}

func (c *Client) parseResponse(resp *http.Response) (auth.ProviderUser, error) {
	if resp.StatusCode != http.StatusOK {
		return c.errorResponse(resp)
	}

	user, err := c.decodeResponseBody(resp)
	if err != nil {
		return auth.ProviderUser{}, err
	}

	return auth.ProviderUser{
		ID:        strconv.Itoa(user.ID),
		AvatarURL: user.AvatarURL,
		Name:      user.Name,
		Email:     user.Email,
	}, nil
}

func (c *Client) errorResponse(resp *http.Response) (auth.ProviderUser, error) {
	body, _ := ioutil.ReadAll(resp.Body)
	err := fmt.Errorf("error requesting user. Status: %d. Response: %s", resp.StatusCode, body)
	return auth.ProviderUser{}, err
}

func (c *Client) decodeResponseBody(resp *http.Response) (*githubUser, error) {
	user := &githubUser{}
	err := json.NewDecoder(resp.Body).Decode(user)
	if err != nil {
		return nil, fmt.Errorf("error decoding JSON on github.Client: %w", err)
	}

	return user, nil
}

type githubUser struct {
	ID        int    `json:"id"`
	AvatarURL string `json:"avatar_url"`
	Type      string `json:"type"`
	Name      string `json:"name"`
	Email     string `json:"email"`
}

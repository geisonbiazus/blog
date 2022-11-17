package github_test

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/geisonbiazus/blog/internal/adapters/oauth2provider/github"
	"github.com/geisonbiazus/blog/internal/core/auth"
	"github.com/stretchr/testify/assert"
)

func TestClient(t *testing.T) {
	setup := func() (*github.Client, *HTTPClientSpy) {
		httpClient := &HTTPClientSpy{}
		client := github.NewClient(httpClient)

		return client, httpClient
	}

	t.Run("GetAuthenticatedUser", func(t *testing.T) {
		t.Run("It gets the GitHub user from the API", func(t *testing.T) {
			client, httpClient := setup()

			httpClient.GetReturnResponse = &http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(strings.NewReader(githubUserResponse)),
			}

			user, _ := client.GetAuthenticatedUser()

			assert.Equal(t, "https://api.github.com/user", httpClient.GetReceivedURL)
			assert.Equal(t, auth.ProviderUser{
				ID:        "1234",
				AvatarURL: "https://example.org/avatar",
				Name:      "User Name",
				Email:     "user@example.com",
			}, user)
		})

		t.Run("It returns error when responso is not 200", func(t *testing.T) {
			client, httpClient := setup()

			httpClient.GetReturnResponse = &http.Response{
				StatusCode: http.StatusNotFound,
				Body:       ioutil.NopCloser(strings.NewReader("Not found")),
			}

			user, err := client.GetAuthenticatedUser()

			assert.Equal(t, auth.ProviderUser{}, user)
			assert.Equal(t, errors.New("error requesting user. Status: 404. Response: Not found"), err)
		})
	})
}

var githubUserResponse = `
{
	"id": 1234,
	"avatar_url": "https://example.org/avatar",
	"type": "User",
	"name": "User Name",
	"email": "user@example.com"
}
`

type HTTPClientSpy struct {
	GetReceivedURL    string
	GetReturnResponse *http.Response
	GetReturnError    error
}

func (c *HTTPClientSpy) Get(url string) (resp *http.Response, err error) {
	c.GetReceivedURL = url
	return c.GetReturnResponse, c.GetReturnError
}

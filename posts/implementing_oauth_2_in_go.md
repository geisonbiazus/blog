title: Implementing OAuth 2.0 in Go
author: Geison Biazus
description: In this post, I show how to implement the OAuth 2.0 standard in Go to securely authenticate into applications using third-party providers.
image_path: /static/image/logo-small.png
time: 2021-11-04 09:00
--
User authentication in software development is a big topic. There are many ways to have a user authenticated in an application, and they vary in complexity based on the system's needs. It is a very important aspect of the application being one of the biggest security issues a system might have.

We can decide to have the regular "username" and "password" approach, but that means that our application will need to store the passwords in a database. These passwords need to be encrypted, so the credentials are not accessible in case of a data leak. This model also requires a user registration feature, which usually requires an account confirmation via email or another kind of message. One way to avoid the implementation of all this is to use OAuth 2.0.

OAuth (Open Authorization) is an open standard for access delegation. It allows users to use their credentials of a third-party application service to grant access to our system securely without the need to store passwords.

Replacing authentication is not the only use case for OAuth. Maybe our system has the traditional username and password method, but we want to provide an alternative way of registering and signing in. Oauth 2.0 is an option to allow the users to choose which provider they want to use to authenticate. OAuth 2.0 is implemented by many providers, such as Google, GitHub, Facebook, Twitter, among others. Allowing users to use their third-party accounts to authenticate is a very common usage of the OAuth 2.0 method.

Another situation where OAuth is useful is when our system needs access to the third-party application to do things on behalf of the user. When the users authenticate using OAuth 2.0, they can grant our application some pre-defined permissions allowing our system to access their account through the provider API. This makes it possible, for example, to publish posts on Facebook, read the repositories on Github, or manage YouTube videos, everything accordingly to the permissions granted by the users.

This post focuses on using OAuth 2.0 to replace the application sign-in / sign-up using GitHub as the Oauth provider. The examples are implemented using the Go programming language.

## Contents

- [OAuth 2.0 Flow](#oauth-20-flow)
- [Implementation](#implementation)
  - [Request OAuth 2.0](#request-oauth-20)
    - [Use Case](#use-case)
    - [Ports](#ports)
    - [Adapters](#adapters)
      - [OAuth2Provider](#oauth2provider)
      - [IDGenerator](#idgenerator)
      - [StateRepo](#staterepo)
    - [HTTP Handler](#http-handler)
  - [Confirm OAuth 2.0](#confirm-oauth-20)
    - [Use Case](#use-case-1)
    - [Ports](#ports-1)
    - [Entities](#entities)
    - [Adapters](#adapters-1)
      - [OAuth2Provider](#oauth2provider-1)
      - [StateRepo](#staterepo-1)
      - [UserRepo](#userrepo)
      - [TokenEncoder](#tokenencoder)
    - [HTTP Handler](#http-handler-1)
- [Final Thoughts](#final-thoughts)
- [References](#references)

## OAuth 2.0 Flow

The Oauth 2.0 flow involves three actors: the user that is trying to authenticate, the application that the user is trying to authenticate to, and the provider that takes care of authenticating the user granting access to the application.

The authentication flow is the following:

- The user accesses an endpoint in our application to sign in.
- The application generates a state. A random string that will be used to validate the authentication.
- The application redirects the user to the Auth 2.0 provider sending the state and a client ID as query parameters.
- Now in the provider application, the user signs in and grants permissions to our application.
- The provider redirects the user back to our application in a pre-configured callback URL. The redirect URL contains a code, generated by the provider, and the same state that before was given to the provider as an argument.
- The application validates the state and requests the user authentication token from the provider using the received code.
- The application is now authenticated in the provider API and can perform requests on it.

After being authenticated in the provider, we can get the user's email address or internal ID, and find or create a user record in our system. With our internal system's user, we can put it in a session using JWT, for example.

## Implementation

As always, we are using Clean Architecture concepts to implement the feature. I have written about it in the post [Applying Clean Architecture in Go](https://blog.geisonbiazus.com/posts/applying-clean-architecture-in-go).

The implementation consists of two use cases. The "Request Oauth 2.0" use case, starts with the user accessing the "/login/github" endpoint and ends with the user being redirected to the GitHub authentication page. The "Confirm Oauth 2.0" use case, starts with the user accessing the "/login/github/confirm" endpoint by being redirected from the GitHub page and ends with a session token being saved on the cookies.

### Request OAuth 2.0

The full implementation of a use case consists of the following parts or layers: Use case, entities, ports, adapters, UI. The `use case` is a module that contains the high-level policy of the feature. It orchestrates the other layers making them work together. This module is also known as "Service Layer" or "Application Layer" in other architectural patterns. `Entities` are managed by the use case and sometimes passed to and returned by the adapters, these entities can be simple data structures or more complex objects, they contain only pure business rules and are not aware of the application-specific business rules which are implemented in the use case layer. The use case interacts with external dependencies through `ports` that are expressed as interfaces in code. On these ports `adapters` are plugged in. These adapters implement the ports interfaces and adapt the external dependencies to our application's use cases.

#### Use Case

We first start with the `RequestOAuth2UseCase`. You can see the implementation below:

```go
// internal/core/auth/request_oauth_2_use_case.go

package auth

import "fmt"

type RequestOAuth2UseCase struct {
	provider  OAuth2Provider
	idGen     IDGenerator
	stateRepo StateRepo
}

func NewRequestOAuth2UseCase(
	provider OAuth2Provider,
	idGen IDGenerator,
	stateRepo StateRepo,
) *RequestOAuth2UseCase {
	return &RequestOAuth2UseCase{
		provider:  provider,
		idGen:     idGen,
		stateRepo: stateRepo,
	}
}

func (u *RequestOAuth2UseCase) Run() (string, error) {
	state := u.idGen.Generate()

	err := u.stateRepo.AddState(state)
	if err != nil {
		return "", fmt.Errorf("error saving state on RequestOAuth2UseCase: %w", err)
	}

	return u.provider.AuthURL(state), nil
}
```

This struct holds the high-level policy for requesting the authentication, the first step of the OAuth 2.0. The `Run` method first generates a state using an `IDGenerator`. Then this state is stored for later using a `StateRepo`. Finally, it requests and returns the authentication URL from an `Auth2Provider`. This is the URL for the user to be redirected to authenticate.

#### Ports

The three dependencies of this use case are interfaces defined in the `ports.go` file:

```go
// internal/core/auth/ports.go

package auth

type OAuth2Provider interface {
	AuthURL(state string) string
}

type IDGenerator interface {
	Generate() string
}

type StateRepo interface {
	AddState(state string) error
}
```

For each one of these interfaces, there is an adapter. Each interface is an abstraction and each adapter is concrete. Meaning we can replace the concrete implementation without changing the high-level rules of the use case.

#### Adapters

For each one of the ports, we have a concrete adapter implementation. For the `OAuth2Provider`. We have the `github.Provider` adapter. For the `IDGenerator`, the `uuid.Generator` concrete implementation. And for the `StateRepo`, we have the `memory.StateRepo`.

##### OAuth2Provider

Here is the implementation of the `github.Provider` module:

```go
// internal/adapters/oauth2provider/github/provider.go

package github

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

type Provider struct {
	config oauth2.Config
}

func NewProvider(clientID, clientSecret string) *Provider {
	config := oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint:     github.Endpoint,
	}

	return &Provider{config: config}
}

func (p *Provider) AuthURL(state string) string {
	return p.config.AuthCodeURL(state)
}
```

Go has support to OAuth 2.0 using the package `golang.org/x/oauth2`. This package abstracts away how to handle OAuth 2.0 with different providers, we only need to initialize the oauth2.Config properly. Part of this package is also several provider Endpoints available to use. We are using the GitHub implementation in this example. You can see a list of all the providers in the [oauth package documentation](https://pkg.go.dev/golang.org/x/oauth2#section-directories).

The `NewProvider` function receives a `clientID` and a `clientSecret` and creates an instance of the `oauth2.Config` with those values and the GitHub endpoints. The `AuthURL` method simply delegates to the config to generate the URL for the user to be redirected.

##### IDGenerator

The next adapter is the `uuid.Generator` which implements the `IDGenerator` interface:

```go
// internal/adapters/idgenerator/uuid/uuid.go

package uuid

import "github.com/google/uuid"

type Generator struct{}

func NewGenerator() *Generator {
	return &Generator{}
}

func (g *Generator) Generate() string {
	return uuid.New().String()
}
```

This package uses the `github.com/google/uuid` library. The `Generate` method generates a UUID version 4 and returns its string format.

##### StateRepo

The last adapter for this first use case is the `memory.StateRepo` which implements the `StateRepo` interface:

```go
// internal/adapters/staterepo/memory/state_repo.go

package memory

type StateRepo struct {
	states map[string]bool
}

func NewStateRepo() *StateRepo {
	return &StateRepo{
		states: make(map[string]bool),
	}
}

func (r *StateRepo) AddState(state string) error {
	r.states[state] = true
	return nil
}
```

This is an in-memory adapter for the `StateRepo` port. It stores all the states on a map that will be checked for existence later. As a state value doesn't live for too long, an in-memory version of this adapter can work in production, as long as the application is not scaled horizontally, and there are no frequent restarts due to deployments. A more stable version would be using a database like Redis or Postgres to store the states. For this example, the in-memory version is good enough.

#### HTTP Handler

Having all the adapters implemented, let's move to the UI layer handling the HTTP requests. For that we need an HTTP Handler that invokes the use case:

```go
// internal/ui/web/request_oauth_2_handler.go

package web

import "net/http"

type RequestOAuth2Handler struct {
	usecase  RequestOAuth2UseCase
	template *TemplateRenderer
}

func NewRequestOAuth2Handler(usecase RequestOAuth2UseCase, template *TemplateRenderer) *RequestOAuth2Handler {
	return &RequestOAuth2Handler{usecase: usecase, template: template}
}

func (h *RequestOAuth2Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	redirectURL, err := h.usecase.Run()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.template.Render(w, "500.html", nil)
		return
	}

	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}
```

The `RequestOAuth2Handler` receives the `RequestOAuth2UseCase` as a dependency. the `ServeHTTP` method calls the `Run` method from the use case retrieving the redirect URL. In case of error, it responds with a status code 500 and renders an error template. If there is no error it redirects the user to the retrieved redirect URL.

The router simply registers this handler to the `/login/github` path:

```go
// internal/ui/web/router.go

package web

import (
	"net/http"
)

func NewRouter(templatePath, staticFilesPath string, usecases *UseCases, baseURL string) http.Handler {
	templateRenderer := NewTemplateRenderer(templatePath, baseURL)

	mux := http.NewServeMux()

	// ...

	mux.Handle("/login/github", NewRequestOAuth2Handler(usecases.RequestOAuth2, templateRenderer))

	return mux
}

```

That ends the "Request OAuth 2.0" use case. By accessing the "/login/github" endpoint, users get redirected to the GitHub authentication page where they give permissions to our app to get their information from GitHub. After that, these users get redirected back to our application where the second use case starts.

### Confirm OAuth 2.0

The "Confirm OAuth 2.0" use case starts when the OAuth 2.0 provider redirects the user back to our application after the proper permissions were granted by this same user. Again the full implementation of this use case consists of the layers mentioned before: Use case, entities, ports, adapters, and UI.

#### Use Case

After the users authenticate on GitHub and provide access to our application, they get redirected to the configured redirect URL on GitHub. For this example, we are using the "/login/github/confirm" endpoint for that. To implement it we start with the `ConfirmOAuth2UseCase`. As this is a big module, we are breaking the implementation into parts:

```go
// internal/core/auth/confirm_oauth_2_use_case.go

package auth

import (
	"context"
	"errors"
	"fmt"
	"time"
)

type ConfirmOAuth2UseCase struct {
	provider     OAuth2Provider
	stateRepo    StateRepo
	userRepo     UserRepo
	idGen        IDGenerator
	tokenEncoder TokenEncoder
}

func NewConfirmOAuth2UseCase(
	provider OAuth2Provider,
	stateRepo StateRepo,
	userRepo UserRepo,
	idGen IDGenerator,
	tokenEncoder TokenEncoder,
) *ConfirmOAuth2UseCase {
	return &ConfirmOAuth2UseCase{
		provider:     provider,
		stateRepo:    stateRepo,
		userRepo:     userRepo,
		idGen:        idGen,
		tokenEncoder: tokenEncoder,
	}
}

// ...

```

We start with the struct definition and the constructor function. Here we can see that this use case has several dependencies. The `Auth2Provider` and the `StateRepo` are the same mentioned in the previous use case and are used here to validate the state and authenticate the user in the provider. The `IDGenerator` and the `UserRepo` are used to create or update the user representation in our system. Finally, the `TokenEncoder` is used to create a session token to authenticate the user.

Now that we have the use case struct let's see the `Run` method implementation:

```go
// internal/core/auth/confirm_oauth_2_use_case.go

// ...

func (u *ConfirmOAuth2UseCase) Run(ctx context.Context, state, code string) (string, error) {
	providerUser, err := u.processOAuth2Authentication(ctx, state, code)
	if err != nil {
		return "", err
	}

	return u.resolveUserAndGetToken(providerUser)
}

// ...

```

As the entry point of the use case, the `Run` method splits the flow into two private (or unexported in Go terms) method calls. The `processOAuth2Authentication` handles the communication with the OAuth 2.0 provider and the `resolveUserAndGetToken` method manages the user creation and authentication in our application based on the OAuth 2.0 response. You can see the implementation of the first flow next:

```go
// internal/core/auth/confirm_oauth_2_use_case.go

// ...

func (u *ConfirmOAuth2UseCase) processOAuth2Authentication(ctx context.Context, state, code string) (ProviderUser, error) {
	err := u.validateAndRemoveState(state)
	if err != nil {
		return ProviderUser{}, err
	}

	return u.getProviderAuthenticatedUser(ctx, code)
}

func (u *ConfirmOAuth2UseCase) validateAndRemoveState(state string) error {
	exists, err := u.stateRepo.Exists(state)
	if err != nil {
		return fmt.Errorf("error checking state on ConfirmOAuth2UseCase: %w", err)
	}

	if !exists {
		return ErrInvalidState
	}

	err = u.stateRepo.Remove(state)
	if err != nil {
		return fmt.Errorf("error authenticating user on ConfirmOAuth2UseCase: %w", err)
	}

	return nil
}

func (u *ConfirmOAuth2UseCase) getProviderAuthenticatedUser(ctx context.Context, code string) (ProviderUser, error) {
	providerUser, err := u.provider.AuthenticatedUser(ctx, code)
	if err != nil {
		return ProviderUser{}, fmt.Errorf("error authenticating user on ConfirmOAuth2UseCase: %w", err)
	}

	return providerUser, nil
}

// ...

```

The `processOAuth2Authentication` method splits the flow again into two private methods: `validateAndRemoveState` and `getProviderAuthenticatedUser`.

The `validateAndRemoveState` method checks if the state exists using the `StateRepo` dependency and, if it does, it is removed from the repository not allowing this confirmation to happen more than once. This state is the random string generated in the first use case "Request OAuth 2.0" that was sent to the provider and returned as an argument when the user got redirected back to our system.

Next is the `getProviderAuthenticatedUser` method. Here it gets the authenticated user from the `OAuth2Provider` dependency passing a context and a code. This code, like the state, also came from the provider in the redirection to our application. The `getProviderAuthenticatedUser` returns a `ProviderUser` which contains the user information needed to proceed with the authentication.

Going back to the `Run` method, having the `ProviderUser` value, it calls the `resolveUserAndGetToken` method:

```go
// internal/core/auth/confirm_oauth_2_use_case.go

// ...

func (u *ConfirmOAuth2UseCase) resolveUserAndGetToken(providerUser ProviderUser) (string, error) {
	user, err := u.createOrUpdateUser(providerUser)
	if err != nil {
		return "", err
	}

	return u.getAuthenticationToken(user)
}

// ...

```

This method splits the flow again into two private methods: `createOrUpdateUser` and `getAuthenticationToken`.

```go
// internal/core/auth/confirm_oauth_2_use_case.go

// ...

func (u *ConfirmOAuth2UseCase) createOrUpdateUser(providerUser ProviderUser) (User, error) {
	user, err := u.userRepo.FindUserByProviderUserID(providerUser.ID)

	if errors.Is(err, ErrUserNotFound) {
		return u.createNewUser(providerUser)
	}

	if err != nil {
		return User{}, fmt.Errorf("error finding user on ConfirmOAuth2UseCase: %w", err)
	}

	return u.updateExistingUser(user, providerUser)
}

func (u *ConfirmOAuth2UseCase) createNewUser(providerUser ProviderUser) (User, error) {
	user := User{
		ID:             u.idGen.Generate(),
		ProviderUserID: providerUser.ID,
		Email:          providerUser.Email,
		Name:           providerUser.Name,
		AvatarURL:      providerUser.AvatarURL,
	}

	err := u.userRepo.CreateUser(user)
	if err != nil {
		return User{}, fmt.Errorf("error creatinng user on ConfirmOAuth2UseCase: %w", err)
	}

	return user, nil
}

func (u *ConfirmOAuth2UseCase) updateExistingUser(user User, providerUser ProviderUser) (User, error) {
	user.Email = providerUser.Email
	user.Name = providerUser.Name
	user.AvatarURL = providerUser.AvatarURL

	err := u.userRepo.UpdateUser(user)
	if err != nil {
		return User{}, fmt.Errorf("error updating user on ConfirmOAuth2UseCase: %w", err)
	}

	return user, nil
}

// ...

```

The `createOrUpdateUser` method uses the `UserRepo` dependency to check if the user already exists. It does that by trying to find a user by the `ProviderUserID` field. This field is the internal user ID in the OAuth 2.0 provider.

If the user does not exist, it creates a new user by calling the `createNewUser` method, which will again use the `UserRepo` to persist the new user. This new user is created with an ID, generated by the `IDGenerator` dependency, and the email, name, and avatar URL provided by the OAuth 2.0 provider. Additionally, it also persists the `ProviderUserID` field.

If a user exists, `updateExistingUser` is called. This method updates the user email, name, and avatar URL to have the information always up to date when users login into our application.

Going back to the `resolveUserAndGetToken` method, the next method called is the `getAuthenticationToken`:

```go
// internal/core/auth/confirm_oauth_2_use_case.go

// ...

const TokenExpiration = 24 * time.Hour

func (u *ConfirmOAuth2UseCase) getAuthenticationToken(user User) (string, error) {
	token, err := u.tokenEncoder.Encode(user.ID, TokenExpiration)
	if err != nil {
		return "", fmt.Errorf("error encoding token on ConfirmOAuth2UseCase: %w", err)
	}

	return token, nil
}

```

This method uses the `TokenEncoder` dependency to encode the ID of the user in a token. This token is also the response of the use case and is the session token to be kept in the browser cookies as it contains the identity of the authenticated user.

As we can see, the `ConfirmOAuth2UseCase` type has multiple responsibilities. Each time it splits the flow into other methods could potentially be extracted into another type and facilitate code reuse. For example, if we need to authenticate a user using a method other than OAuth 2.0, the `getAuthenticationToken` method could become a standalone type. If we need to process the OAuth 2.0 confirmation but not proceed to the user creation and authentication, the `getProviderAuthenticatedUser` method could become another standalone type. For the moment, the current implementation is sufficient as these functionalities are only used in this module and we don't know yet what could potentially be reused in the future.

You can see the [full implementation](https://github.com/geisonbiazus/blog/blob/main/internal/core/auth/confirm_oauth_2_use_case.go) of `ConfirmOAuth2UseCase` on GitHub.

#### Ports

For this new use case, we need to extend some of the interface dependencies and add others. The interfaces changes are shown below:

```go
// internal/core/auth/ports.go

package auth

import (
	"context"
	"time"
)

type OAuth2Provider interface {
	AuthURL(state string) string
	AuthenticatedUser(ctx context.Context, code string) (ProviderUser, error)
}

type IDGenerator interface {
	Generate() string
}

type StateRepo interface {
	AddState(state string) error
	Exists(state string) (bool, error)
	Remove(state string) error
}

type UserRepo interface {
	CreateUser(user User) error
	UpdateUser(user User) error
	FindUserByProviderUserID(providerUserID string) (User, error)
}

type TokenEncoder interface {
	Encode(value string, expiresIn time.Duration) (string, error)
}
```

The `OAuth2Provider` has a new method `AuthenticatedUser`. `IdGenerator` is just reused without modification. On the `StateRepo`, two methods were added: `Exists` and `Remove`. And two new interfaces were added: `UserRepo` and `TokenEncoder`.

#### Entities

The `Oauth2Provider` now returns a `ProviderUser` and the `UserRepo` receives and returns instances of `User`. These entities are defined below:

```go
// internal/core/auth/entities.go

package auth

import "errors"

type ProviderUser struct {
	ID        string
	Email     string
	Name      string
	AvatarURL string
}

type User struct {
	ID             string
	ProviderUserID string
	Email          string
	Name           string
	AvatarURL      string
}

var ErrInvalidState = errors.New("invalid state error")
var ErrUserNotFound = errors.New("user not found")
var ErrTokenExpired = errors.New("token expired")
```

These two entities, although very similar, correspond to two different things. The `ProviderUser` structs is a representation of the user returned by the provider. The `User` struct is a representation of the internal user of our system. Although they are very similar, they can change for different reasons, so they should be separated. This file also defines the errors that are returned by the adapters and used by the use case.

#### Adapters

Now let's go to the implementation of the adapters. Let's start with the new method in the `OAuth2Provider`:

##### OAuth2Provider

```go
// internal/adapters/oauth2provider/github/provider.go

package github

import (
	"context"
	"fmt"
	"net/http"

	"github.com/geisonbiazus/blog/internal/core/auth"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

type Provider struct {
	config oauth2.Config
}

func NewProvider(clientID, clientSecret string) *Provider {
	config := oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint:     github.Endpoint,
	}

	return &Provider{config: config}
}

// ...

func (p *Provider) AuthenticatedUser(ctx context.Context, code string) (auth.ProviderUser, error) {
	httpClient, err := p.exchangeTokenAndGetClient(ctx, code)
	if err != nil {
		return auth.ProviderUser{}, err
	}

	return NewClient(httpClient).GetAuthenticatedUser()
}

func (p *Provider) exchangeTokenAndGetClient(ctx context.Context, code string) (*http.Client, error) {
	token, err := p.config.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("error exchanging token on github. Provider: %w", err)
	}

	tokenSource := p.config.TokenSource(ctx, token)
	return oauth2.NewClient(ctx, tokenSource), nil
}
```

The `AuthenticateUser` method calls the private method `exchangeTokenAndGetClient`. This method receives a context and a code and returns an instance of `http.Client` that is already authenticated on the GitHub API. The authentication on GitHub is done using `golang.org/x/oauth2` library. Using our instance of the `oauth2.Config`, the `Exchange` method is called returning an instance of the `oauth2.Token` type.

The `oauth2.Token` contains both an access token and a refresh token. It is serializable to JSON and can be saved in case our application needs to keep the connection to the provider persistent. So, for example, in the use case of authorizing our application to have access to GitHub to continuously collect metrics from the users' repositories, the users would give our application access to their GitHub account only once and the subsequent authentications would be done by either using the access token or the refresh token to get a new access token.

The `oauth2.Token` is used to get an instance of a `auth2.TokenSource` from the `oauth2.Config` instance, which is passed to the `oauth2.NewClient` method returning an `http.Client`. This client is already authenticated to the GitHub API and it also automatically refreshes the access token in case it expires.

Back to the `AuthenticatedUser` method, with the authenticated HTTP client in hand, the `github.NewClient` method is used to create a `github.Client` instance and then request the authenticated user information from GitHub using the `GetAuthenticatedUser` method.

You can see the `github.Client` implementation below:

```go
// internal/adapters/oauth2provider/github/client.go

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
```

The `github.Client` type uses has one public method called `GetAuthenticatedUser` that uses the GitHub REST API to get the authenticated user information. The API returns A JSON representation of the user, which is decoded and then used to build an `auth.ProviderUser`.

##### StateRepo

The `StateRepo` adapter needs to be enhanced with the `Exists` and `Remove` methods to satisfy the "Request Oauth 2.0" use case:

```go
// internal/adapters/staterepo/memory/state_repo.go

package memory

type StateRepo struct {
	states map[string]bool
}

func NewStateRepo() *StateRepo {
	return &StateRepo{
		states: make(map[string]bool),
	}
}

func (r *StateRepo) AddState(state string) error {
	r.states[state] = true
	return nil
}

// New methods bellow

func (r *StateRepo) Exists(state string) (bool, error) {
	_, ok := r.states[state]
	return ok, nil
}

func (r *StateRepo) Remove(state string) error {
	delete(r.states, state)
	return nil
}
```

##### UserRepo

The `UserRepo` adapter, like the `StateRepo`, uses an in-memory implementation for this post. Even though an "in-memory" implementation works well for testing the application, its data is lost every time the application is restarted, so for production, it is required a database or other kind of persistence implementation. As this is not the focus of this post, I'll display only the in-memory version here. You can see the implementation below:

```go
// internal/adapters/userrepo/memory/user_repo.go

package memory

import (
	"github.com/geisonbiazus/blog/internal/core/auth"
)

type UserRepo struct {
	users []auth.User
}

func NewUserRepo() *UserRepo {
	return &UserRepo{users: []auth.User{}}
}

func (r *UserRepo) CreateUser(user auth.User) error {
	r.users = append(r.users, user)
	return nil
}

func (r *UserRepo) UpdateUser(user auth.User) error {
	for i, existingUser := range r.users {
		if existingUser.ID == user.ID {
			r.users[i] = user
			return nil
		}
	}
	return auth.ErrUserNotFound
}

func (r *UserRepo) FindUserByProviderUserID(providerUserID string) (auth.User, error) {
	for _, user := range r.users {
		if user.ProviderUserID == providerUserID {
			return user, nil
		}
	}

	return auth.User{}, auth.ErrUserNotFound
}

```

##### TokenEncoder

The last atapter is the `TokenEncoder`. For this adapter, we have a JWT implementation to sign and encode the user reference to be kept in the user session. You can see the implementation below:

```go
// internal/adapters/tokenencoder/jwt/jwt.go

package jwt

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/geisonbiazus/blog/internal/core/auth"
)

type TokenEncoder struct {
	secret []byte
}

func NewTokenEncoder(secret string) *TokenEncoder {
	return &TokenEncoder{secret: []byte(secret)}
}

func (m *TokenEncoder) Encode(value string, expiresIn time.Duration) (string, error) {
	claims := newClaims(value, expiresIn)
	t := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	signedToken, err := t.SignedString(m.secret)
	if err != nil {
		return "", fmt.Errorf("error signing string on jwt.TokenManager: %w", err)
	}

	return signedToken, nil
}

type jwtClaims struct {
	jwt.StandardClaims
}

func newClaims(sub string, expiresIn time.Duration) *jwtClaims {
	return &jwtClaims{
		StandardClaims: jwt.StandardClaims{
			Subject:   sub,
			ExpiresAt: expiresAt(expiresIn),
		},
	}
}

func expiresAt(expiresIn time.Duration) int64 {
	return time.Now().Add(expiresIn).Unix()
}

```

The `jwt.TokenEncoder` uses the `github.com/dgrijalva/jwt-go` library to generate and sign the JWT. As in the content of this post we just encode the token, we are omitting the `Decode` method. You can see the full implementation on the [blog repository](https://github.com/geisonbiazus/blog/blob/main/internal/adapters/tokenencoder/jwt/jwt.go?ts=2) on GitHub. You can also learn more about JWT on the [JWT website](https://jwt.io)

##### HTTP Handler

Having all the adapters implemented, what is missing is an HTTP Handler to call our use case. You can see it below:

```go
// internal/ui/web/confirm_oauth_2_handler.go

package web

import (
	"errors"
	"net/http"

	"github.com/geisonbiazus/blog/internal/core/auth"
)

type ConfirmOAuth2Handler struct {
	usecase  ConfirmOAuth2UseCase
	template *TemplateRenderer
	baseURL  string
}

func NewConfirmOAuth2Handler(
	usecase ConfirmOAuth2UseCase,
	templateRenderer *TemplateRenderer,
	baseURL string,
) *ConfirmOAuth2Handler {
	return &ConfirmOAuth2Handler{
		usecase:  usecase,
		template: templateRenderer,
		baseURL:  baseURL,
	}
}

func (h *ConfirmOAuth2Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	state := r.URL.Query().Get("state")
	code := r.URL.Query().Get("code")

	token, err := h.usecase.Run(r.Context(), state, code)

	if err != nil {
		h.respondWithError(w, err)
		return
	}

	http.SetCookie(w, h.newSessionCookie(token))
	http.Redirect(w, r, h.baseURL, http.StatusSeeOther)
}

func (h *ConfirmOAuth2Handler) respondWithError(w http.ResponseWriter, err error) {
	if errors.Is(err, auth.ErrInvalidState) {
		w.WriteHeader(http.StatusNotFound)
		h.template.Render(w, "404.html", nil)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		h.template.Render(w, "500.html", nil)
	}
}

func (h *ConfirmOAuth2Handler) newSessionCookie(token string) *http.Cookie {
	return &http.Cookie{
		Name:  "_blog_session",
		Value: token,
		Path:  "/",
	}
}

```

The `ConfirmOAuth2Handler` serves the endpoint in which the users will be redirected from the OAuth 2.0 provider. The `ServeHTTP` method gets the state and the code from the URL query params and executes the `ConfirmOAuth2UseCase` by calling the `Run` method. If an error is returned, the appropriate error page is rendered. In case there is no error, it means that the authentication was successful and we have the user token. So this token is set to the browser cookies and the user is redirected to the "/" endpoint.

The handler is registered in the router to the `/login/github/confirm` path:

```go
// internal/ui/web/router.go

package web

import (
	"net/http"
)

func NewRouter(templatePath, staticFilesPath string, usecases *UseCases, baseURL string) http.Handler {
	templateRenderer := NewTemplateRenderer(templatePath, baseURL)

	mux := http.NewServeMux()

	// ...

	mux.Handle("/login/github/confirm", NewConfirmOAuth2Handler(usecases.ConfirmOAuth2, templateRenderer, baseURL))

	return mux
}
```

That ends the "Confirm OAuth 2.0" use case and the implementation of the OAuth 2.0 authentication. After that, we can get the value of the `_blog_session` cookie in the subsequent requests, decode it and use it to get the authenticated user.

## Final Thoughts

OAuth 2.0, being a standard that is followed by many providers, facilitates applications to authenticate and connect to many third-party systems in a secure way. There are many libraries that abstract this authentication process but this is just one portion of the full process, which also consists of managing the users and generating session tokens. The main challenge is to orchestrate this process in a testable and maintainable way, and for that, the Clean Architecture is a good fit.

You can see the full implementation in the [Blog Repository](https://github.com/geisonbiazus/blog) on GitHub.

## References

[jwt package - github.com/dgrijalva/jwt-go - pkg.go.dev](https://pkg.go.dev/github.com/dgrijalva/jwt-go)

[JWT website](https://jwt.io)

[OAuth - Wikipedia](https://en.wikipedia.org/wiki/OAuth)

[OAuth 2.0 - OAuth](https://oauth.net/2/)

[oauth2 package - golang.org/x/oauth2 - pkg.go.dev](https://pkg.go.dev/golang.org/x/oauth2)
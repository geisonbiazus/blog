title: The Different Types Of Mocks
author: Geison Biazus
description: In this post, I show the different types of mocks, when to use them, and how they can be implemented.
image_path: /static/image/logo-small.png
time: 2021-08-15 09:00
--

One of the most useful techniques when writing tests is the usage of test doubles. A test double is a special type of object or module, that is given as a dependency to the code under test. This object follows the same interface as the real dependency but allows us to manipulate its response to guide the code under test through the desired path.

Test doubles are more generally known as mocks, but there are several types of test doubles, each one is useful for a specific situation, and the mock is only one of these types. In this article, I'll show you how to implement and how to use each one of these types of test doubles.

To exemplify each kind of test double, I'm going to use the following use case. This is a very simplistic user creation implementation, but it allows us to use multiple types of test doubles on its different test cases. Here is a brief description of the use case:

- **Use case:** Create user
- **Arguments:** `email` and `password`
- **Return:** `User`, `ErrInvalidCrendentials`, or `ErrUserAlreadyExists`
- **Busines rules:**
  - If `email` or `password` are invalid, it returns `ErrInvalidCrendentials`
  - If there is already a user with the given `email`, it returns `ErrUserAlreadyExists`
  - Otherwise it creates, persists and returns the `User`

Here is the use case implementation. I used the Go programming language, but these concepts can be applied in any language:

```go
type User struct {
	Email    string
	Password string
}

type UserRepo interface {
	GetUserByEmail(email string) (User, error)
	CreateUser(user User) error
}

type CreateUserUseCase struct {
	UserRepo UserRepo
}

var ErrInvalidCrendentials = errors.New("invalid credentials")
var ErrUserAlreadyExists = errors.New("user already exists")
var ErrUserNotFound = errors.New("user not found")

func (u *CreateUserUseCase) Run(email, password string) (User, error) {
	if !u.validEmailAndPassword(email, password) {
		return User{}, ErrInvalidCrendentials
	}

	if u.userExists(email) {
		return User{}, ErrUserAlreadyExists
	}

	return u.createNewUser(email, password)
}

func (u *CreateUserUseCase) validEmailAndPassword(email, password string) bool {
	return email != "" && password != ""
}

func (u *CreateUserUseCase) userExists(email string) bool {
	_, err := u.UserRepo.GetUserByEmail(email)
	return err == nil
}

func (u *CreateUserUseCase) createNewUser(email, password string) (User, error) {
	user := User{
		Email:    email,
		Password: password,
	}

	err := u.UserRepo.CreateUser(user)
	if err != nil {
		return User{}, err
	}

	return user, nil
}
```

The code under test is the `CreateUserUseCase` struct. The use case receives a `UserRepo` as a dependency which is an interface with two methods: `GetUserByEmail` and `CreateUser`. This is the dependency that we are going to replace with test doubles. The use case is executed by calling the `Run` method.

In the next sessions, I explain each type of test double, how to implement and for which kind of situation it is useful.

## Dummy

The first and simplest is the dummy. It implements the interface but its methods return "empty" values like `null` or `0`. Its usage is simply to satisfy a dependency when testing a scenario that this dependency will not be called. From all test double types, this is the least useful as any other type can be used in its place.

Here is how to implement it:

```go
func TestInvalidEmailAndPasswordWithDummy(t *testing.T) {
	repo := UserRepoDummy{}
	usecase := &CreateUserUseCase{UserRepo: repo}

	_, err := usecase.Run("", "")

	if err != ErrInvalidCrendentials {
		t.Errorf("Did not return correct error: %v", err)
	}
}

type UserRepoDummy struct{}

func (r UserRepoDummy) GetUserByEmail(email string) (User, error) {
	return User{}, nil
}

func (r UserRepoDummy) CreateUser(user User) error {
	return nil
}
```

This test scenario only checks for the arguments validation, so the `UserRepo` dependency is never called. That's why a "dummy" implementation is given only to satisfy the dependency.

## Stub

The stub is one step ahead of the dummy. It has a similar implementation, but it returns pre-defined values. It is used when you want to drive the code into the desired path. You can have multiple stub implementations used by different test scenarios, each one returning a different value or throwing exceptions so you can test your code reacting to different values returned from your dependency.

Here is how to implement it:

```go
func TestUserAlreadyExistsWithStub(t *testing.T) {
	repo := UserRepoStub{}
	usecase := &CreateUserUseCase{UserRepo: repo}

	_, err := usecase.Run("user@example.com", "password")

	if err != ErrUserAlreadyExists {
		t.Errorf("Did not return correct error: %v", err)
	}
}

type UserRepoStub struct{}

func (r UserRepoStub) GetUserByEmail(email string) (User, error) {
	return User{Email: "user@example.com", Password: "password"}, nil
}

func (r UserRepoStub) CreateUser(user User) error {
	return nil
}
```

In this scenario, a stub is given to the use case. This stub will always return a user in the `GetUserByEmail` method meaning that the user "exists" with the given email. So it guides the production code to return the expected `ErrUserAlreadyExists` error.

## Spy

Here is where things start getting more interesting. The spy is again one step ahead of the stub. It is a configurable type of test double where you can specify what each method will return. It also records the calls to its methods and saves the received arguments allowing you to later verify if they were called with the correct arguments.

The spy is a very complete type of test double, and it can easily replace the previous types. The only downside is that it requires a little bit more code to implement it. Here is how it looks like:

```go
func TestCreateNewUserWithSpy(t *testing.T) {
	repo := &UserRepoSpy{}
	usecase := &CreateUserUseCase{UserRepo: repo}

	repo.GetUserByEmailReturnError = ErrUserNotFound
	repo.CreateUserReturnError = nil

	email := "user@example.com"
	password := "password"

	returnedUser, err := usecase.Run(email, password)

	expectedUser := User{Email: email, Password: password}

	if returnedUser != expectedUser {
		t.Errorf("Did not return correct user: %v", returnedUser)
	}

	if repo.CreateUserReceivedUser != expectedUser {
		t.Errorf("UserRepo.CreateUser did not receive correct user: %v", repo.CreateUserReceivedUser)
	}

	if err != nil {
		t.Errorf("Returned error is not nil: %v", err)
	}
}

type UserRepoSpy struct {
	GetUserByEmailReceivedEmail string
	GetUserByEmailReturnUser    User
	GetUserByEmailReturnError   error

	CreateUserReceivedUser User
	CreateUserReturnError  error
}

func (r *UserRepoSpy) GetUserByEmail(email string) (User, error) {
	r.GetUserByEmailReceivedEmail = email
	return r.GetUserByEmailReturnUser, r.GetUserByEmailReturnError
}

func (r *UserRepoSpy) CreateUser(user User) error {
	r.CreateUserReceivedUser = user
	return r.CreateUserReturnError
}
```

The `UserRepoSpy` allows us to configure the value that will be returned when its methods are called. Additionally, whenever those methods are executed, it remembers the received arguments. If we pick the `GetUserByEmail` as an example, we can configure what will be returned by setting values to the `GetUserByEmailReturnUser` and `GetUserByEmailReturnError` attributes. Whenever this method is called, we can assert if it received the correct argument by checking the value of the `GetUserByEmailReceivedEmail` attribute. The same behavior goes to the `CreateUser`method.

## Mock

With the popularization of the mocking libraries, the term mock became very common and it is used to refer to any test double. But the true mock, like the previous types, is one step ahead of the spy. It allows you to configure what the methods return, it also remembers what its methods receive, but the difference is that it knows how to verify if its methods were called with the correct arguments. The test itself doesn't check for anything, it just checks if everything went ok.

Here is how it can be implemented:

```go
func TestCreateNewUserWithMock(t *testing.T) {
	repo := &UserRepoMock{}
	usecase := &CreateUserUseCase{UserRepo: repo}

	email := "user@example.com"
	password := "password"
	expectedUser := User{Email: email, Password: password}

	repo.MockGetUserByEmail(email, User{}, ErrUserNotFound)
	repo.MockCreateUser(expectedUser, nil)

	returnedUser, _ := usecase.Run(email, password)

	if returnedUser != expectedUser {
		t.Errorf("Did not return correct user: %v", returnedUser)
	}

	repo.Verify(t)
}

type UserRepoMock struct {
	isGetUserByEmailMocked      bool
	getUserByEmailReceiveEmail  string
	getUserByEmailReceivedEmail string
	getUserByEmailReturnUser    User
	getUserByEmailReturnError   error

	isCreateUserMocked     bool
	createUserReceiveUser  User
	createUserReceivedUser User
	createUserReturnError  error
}

func (r *UserRepoMock) MockGetUserByEmail(
	receiveEmail string, returnUser User, returnError error,
) {
	r.isGetUserByEmailMocked = true
	r.getUserByEmailReceiveEmail = receiveEmail
	r.getUserByEmailReturnUser = returnUser
	r.getUserByEmailReturnError = returnError
}

func (r *UserRepoMock) GetUserByEmail(email string) (User, error) {
	r.getUserByEmailReceivedEmail = email
	return r.getUserByEmailReturnUser, r.getUserByEmailReturnError
}

func (r *UserRepoMock) MockCreateUser(receiveUser User, returnError error) {
	r.isCreateUserMocked = true
	r.createUserReceiveUser = receiveUser
	r.createUserReturnError = returnError
}

func (r *UserRepoMock) CreateUser(user User) error {
	r.createUserReceivedUser = user
	return r.createUserReturnError
}

func (r *UserRepoMock) Verify(t *testing.T) {
	if r.isGetUserByEmailMocked {
		if r.getUserByEmailReceiveEmail != r.getUserByEmailReceivedEmail {
			t.Errorf(
				"Did not receive the correct email on GetUserByEmail\nExpected: %v\nReceived: %v",
				r.getUserByEmailReceiveEmail, r.getUserByEmailReceivedEmail,
			)
		}
	}

	if r.isCreateUserMocked {
		if r.createUserReceiveUser != r.createUserReceivedUser {
			t.Errorf(
				"Did not receive the correct user on CreateUser\nExpected: %v\nReceived: %v",
				r.createUserReceiveUser, r.createUserReceivedUser,
			)
		}
	}
}
```

With `UserRepoMock` we can configure its methods behavior by calling the `MockGetUserByEmail` and `MockCreateUser` methods. These methods receive the expected arguments and what they should return. Then in the test, we call the `Verify` method which checks if each method was called with the correct arguments.

A mock has the benefit of making the test a lot cleaner compared to the spy. But the downside is that its implementation is more complex and verbose. This is where mocking libraries come to help with their dynamic behavior, but with the cost of coupling to the libraries. So, like always, it is a trade-off to use them or not.

## Fake

Here comes my favorite. The fake is not in the same category as the previous types. Different from the others, it has behavior and it simulates the real dependency. A fake can be used by the system in place of the original dependency without any problem. The most common examples of fakes are "in memory" repositories implementations.

Here is how it can be implemented:

```go
func TestCreateNewUserWithFake(t *testing.T) {
	repo := NewFakeUserRepo()
	usecase := &CreateUserUseCase{UserRepo: repo}

	email := "user@example.com"
	password := "password"

	returnedUser, _ := usecase.Run(email, password)

	expectedUser := User{Email: email, Password: password}

	if returnedUser != expectedUser {
		t.Errorf("Did not return correct user: %v", returnedUser)
	}

	createdUser, _ := repo.GetUserByEmail(email)

	if returnedUser != expectedUser {
		t.Errorf("Did not return create user: %v", createdUser)
	}
}

type FakeUserRepo struct {
	users map[string]User
}

func NewFakeUserRepo() *FakeUserRepo {
	return &FakeUserRepo{
		users: map[string]User{},
	}
}

func (r *FakeUserRepo) GetUserByEmail(email string) (User, error) {
	user, ok := r.users[email]
	if !ok {
		return User{}, ErrUserNotFound
	}

	return user, nil
}

func (r *FakeUserRepo) CreateUser(user User) error {
	r.users[user.Email] = user
	return nil
}
```

The `FakeUserRepo` works exactly how a relational database implementation of the `UserRepo` would work. It stores the created users in a map, and it returns the user when queried. The fake is great for integration tests to avoid external APIs or other dependencies being called.

The fake has one downside though. As it has behavior, it can become very complex and many of its behaviors need to be duplicated in the real implementation. Also because of its complexity, it might require tests for the fake itself. So my recommendation is to always keep the fake implementation very simple. It doesn't need to simulate the real thing completely, only the necessary to achieve its goal.

## Monkey Patch

One last type of "mocking" that is frequently used is monkey patching. This is very common in dynamic languages like Ruby, Python, or JavaScript. It changes the runtime implementation of the code during the test execution and reverts back to the original after the test is finished. In my opinion, it encourages the implementation of bad code and coupled code and it should be avoided at all costs. The only situation that it can have some value is when dealing with legacy code.

José Valin wrote about using mocks as "nouns" in his great article [Mocks and explicit contracts](http://blog.plataformatec.com.br/2015/10/mocks-and-explicit-contracts/).

## Mocking libraries

There are plenty of mocking libraries that can help you to create spies and mocks with less code than writing them by hand. But as always, every library that you put in your project is a new dependency that needs to be maintained, and these dependencies have their own dependencies and so on. So coupling your tests to a mocking library means that you cannot remove or replace the library easily. One option would be to wrap the mocks created by the library on your own modules where you have control, but this leads to having a more verbose solution, which compared to how easy it is to write your own mocks, might not be worth it.

That doesn't mean I don't use any mocking library. I think very focused and specific libraries are very useful. For example, a library that mocks HTTP requests, or a library that simulates a specific database. These mocks are really helpful to test adapters without the need for the real dependency. But they are used in only one place for a specific situation.

## Don't mock everything

Every time a test double is used, we are coupling our tests with the implementation. If we use them on every class and dependency, the tests become fragile making the code hard to refactor. So I try to use doubles in the following two situations:

1. When crossing an architectural boundary. We create our boundaries because we don't want to couple our system with libraries or external dependencies. In these situations, using the real implementation in the tests will make them break in case these dependencies change. So using a double as the dependency adapter keeps the code decoupled.

2. When dealing with randomness. Generating random numbers, random IDs, or getting the current date. Randon code is usually hard to test so using a stub can make the code predictable and testable.

## Final thoughts

In this post, I wrote about the multiple types of test doubles. It is important to know and understand them really well. Naming your doubles accordingly will facilitate the communication between the team, and the readability of the tests as you know what to expect when you read the code.

The implementations I showed you here are just some examples. In the end, a test double is a piece of code, and as such, you can design it to fulfill the needs of your tests.

## References

[Clean Code: Advanced TDD, Episode 23, Part 1](https://cleancoders.com/episode/clean-code-episode-23-p1)

[The Little Mocker - Robert C. Martin](https://blog.cleancoder.com/uncle-bob/2014/05/14/TheLittleMocker.html)

[Mocks Aren't Stubs - Martin Fowler](https://martinfowler.com/articles/mocksArentStubs.html)

[Mocks and explicit contracts - José Valin](http://blog.plataformatec.com.br/2015/10/mocks-and-explicit-contracts/)

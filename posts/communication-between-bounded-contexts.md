title: Communication between bounded countexts
author: Geison Biazus
description:
image_path: /static/image/logo-small.png
time: 2023-05-12 09:00
--

## Introduction

The [Bounded Context](https://martinfowler.com/bliki/BoundedContext.html) is one of the main patterns on Domain-Driven Design. It is used to simplify the domain model by spliting it into smaller models based on the [Ubiquitous Language](https://martinfowler.com/bliki/UbiquitousLanguage.html). Everything inside of a Bounded Context follows the domain language of that context and the context contains everything that is needed for it to make sense, to tell its story, but nothing more than that.

A Bounded Context is an isolated subset of the domain, but that doesn't mean it does not deppend on other contexts. It is very common for the same entity to be present in multiple contexts, even with different names to follow the internal Ubiquitous Language of the Bounded Context.

Let's say we have a blogging system where users can sign up to be able to write comments for the posts in the blog. We can define two bounded contexts in this application: Auth and Discussion.

The Auth context contains everything that is required to allow users to register and sign in into the application. It contains all the complexities about session management, password encryption, OAuth, email confirmation, among others. The User entity knows and holds data that is used to satisfy every authentication related use case.

The Discussion context contains everything that is required to allow users to comment blog posts, and other users to answer those comments.

## The problem

The same user that signs in into the application also has the hability to write comments in the blog posts. So we have here a clear dependency between the Auth and the Discussion bounded contexts. The straightforward way to model this is to add a direct dependency in the `discussion.Comment` to `auth.User`:

```go
// auth/user.go

package auth

type User struct {
  ID                 string
  Name               string
  Email              string
  AvatarURL          string
  Encrypted_password string
}

// discussion/comment.go

package discussion

type Comment struct {
  ID      string
  Content string
  User    *auth.User
}
```

The `discussion.Comment` contains a field called `user` of the type `auth.User`. Both entities reside in distinct packages. Although this approach works well for small applications, it has some downsides:

- Every change in the `auth.User` type can affect directly the `discussion.Comment`. E.g. if a field is renamed, all the references to that field should be adjusted.

- If any entity of the `auth` context ever depends on the `discussion` context, a cyclic dependency is created which, in Go, fails to compile. Some other languages support that, but it is usually seen as an anti-pattern to have it between bounded contexts or components as changes on any of the components can propagate into the whole dependency tree making higher and lower level components hard identify.

- If the system grows and we want to extract the `discussion` context into a standalone service, the `auth` context would need to be extracted together. If `auth` also depends on other bounded contexts, those would also need to be extracted making the standalone service be bigger than it should be, and containing components that are not in use.

## Isolating Bounded Contexts

One way of solving this problem is to completely isolate the domain model of both Bounded Contexts. To do so, the Discussion Context cannot deppend on any entity of the Auth Context, therfore the `User` dependency must be changed:

```go
// auth/user.go

package auth

type User struct {
  ID                 string
  Name               string
  Email              string
  AvatarURL          string
  Encrypted_password string
}

// discussion/comment.go

package discussion

type Comment struct {
  ID      string
  Content string
  Author    *Author
}

type Author struct {
  ID                 string
  Name               string
  Email              string
  AvatarURL          string
}
```

Instead of the `discussion.Comment` entity depend on the `auth.User`, we created another entity called `discussion.Author` inside of the Discussion Context. This entity represents the same person as the `auth.User` but in the context and the language of the Discussion Context. For that reason it is called "Author" instead of "User". Sometimes the same system "entity" is called by a different name in different Bounded Contexts to follow the Ubiquitous Language of that context, and that is what has happened here.

Now we have the domain model isolated, but that brings us another problem. This isolation was created to satisfy the different contexts of the domain model, but users of the system don't want to manually create the "user" twice. Instead, we need a way that whenever a user from the Auth Context is created or updated, this changes are reflected to the author from the Discussion Context. The next sections discribe a few ways that this can be achieved.

## Shared Database Table

We can have the domain model isolated into both contexts, but the data is shared by having both entities loaded from the same database table. This way we have a Repository in each one of the Bounded Contexts that hides this complexity:

```go
// auth/user_repository.go

package auth

type UserRepository interface {
  SaveUser(user *User) error
  GetUserByEmail(email String) (*User, error)
  GetUserByID(id string) (*User, error)
}

// discussion/author_repository.go

package discussion

type AuthorRepository interface {
  GetAuthorByEmail(email String) (*Author, error)
  GetAuthorById(id string) (*Author, error)
}

// auth/adapters/user_repository.go

package adapters

type UserRepository struct {}

func (r *UserRepository) SaveUser(user *auth.User) error {
  // Save User on "users" table
}

func (r *UserRepository) GetUserByEmail(email String) (*auth.User, error) {
  // Load User from the "users" table by email
}

func (r *UserRepository) GetUserByID(id String) (*auth.User, error) {
  // Load User from the "users" table by ID
}

// discussion/adapters/auth_repository.go

package adapters

type AuthorRepository struct {}

func (r *AuthorRepository) GetAuthorByEmail(email String) (*discussion.Author, error) {
  // Load Author from the "users" table by email
}

func (r *AuthorRepository) GetAuthorById(ID string) (*discussion.Author, error) {
  // Load Author from the "users" table by ID
}

```

The domain model of both Bounded Contexts depends only on their own repository interfaces. For each interface we have an adapter implementing it. This adapter, inside of the `adapters` package of each one of the contexts, shields the domain model from the place it is stored allowing both bounded contexts to have their own Ubiquitous Language and evolve independently.

Although the domain is isolated, we have some problems with this approach:

- Changes on the `users` database table will affect both repositories. If they are deployed in different services, coordinating the changes can be challenging.

- We can have separate services but the database is still shared, so if one service increases the overall load of the database, other services will be affected.

Despite the problems, sharing the database between contexts is a viable solution when dealing with legacy code where the domain model of a new Bounded Context can be isolated from the legacy system making this new context testable and maintainable.

## Load Dependencies Through an API

If we want to keep the Auth Context with full control and ownership of the `users` data, we can instead add the Auth Context as a dependency of the Discussion Context. But differently from the simply adding this dependency in the domain model, we keep the Discussion domain isolated by using the `disussion.Author` as a replacement for the `auth.User`, but we hide the dependency inside of the `adapters.AuthRepository`.

```go
// discussion/adapters/auth_repository.go

package adapters

type AuthorRepository struct {
  getUserByEmailUseCase auth.GetUserByEmailUseCase
  getUserByIDUseCase auth.GetUserByIDUseCase
}

func (r *AuthorRepository) GetAuthorByEmail(email String) (*discussion.Author, error) {
  user, err := u.getUserByEmailUseCase.run(email)
  if err := nil {
    return nil, err
  }
  return u.toAuthor(user)
}

func (r *AuthorRepository) GetAuthorByID(id string) (*discussion.Author, error) {
  user, err := u.getUserByIDUseCase.run(id)
  if err := nil {
    return nil, err
  }
  return u.toAuthor(user)
}

func (r *AuthorRepository) GetAuthorByID(user *auth.User) *discussion.Author {
  return &discussion.Author {
    ID:        user.ID,
    Name:      user.Name,
    Email:     user.Email,
    EvatarURL: user.EvatarURL,
  }
}
```

In this example, we load the user directly from the use cases of the Auth Context, but as it is hidden inside of the AuthorRepository implementation, it can be easily replaced. Let's say we now want to extract the Discussion context to a service. What we need to do is to implement a new AuthorRepository that instead of calling use cases, performs HTTP requests or gRPC calls. The database doesn't need to be shared between the contexts anymore.

Although this approach we solves some of the problems of both contexts sharing the database, we introduce a few others in case they are separate services:

- The Discussion Context has a runtime dependency on the Auth Context. If they are services and the Auth Service goes down, the Discussion service also stops working.

- Synchronizing deployments with API changes can be tricky. An older version of the Discussion Context might access a newer version of the Auth Context during the deployment of the services. To avoid that, we always need to have the newer and older version of the API working until we are sure the older version is not accessed anymore.

## Synchronously Create Author

Another solution is to invert the dependency between contexts. Instead of the Discussion Context load the Author from the Auth Context, we save the Author on its own database table inside of the Discussin Context and we make the Auth context call the Dicussion context to create or update the Author when a User is created or updated.

```go
// discussion/author_repository.go

package discussion

type AuthorRepository interface {
  SaveAuthor(author *Author) error
  GetAuthorByEmail(email String) (*Author, error)
  GetAuthorById(id string) (*Author, error)
}

// discussion/save_author_use_case.go

type SaveAuthorInput struct {
  Name      string
	Email     string
	AvatarURL string
}

type SaveAuthorUseCase struct {
	authorRepository AuthorRepository
}

func (u *SaveAuthorUseCase) Run(input SaveAuthorInput) (*Author, error) {
	author, err := u.authorRepository.GetAuthorByEmail(input.email)
	// If author is not foud initialize a new Author instance

	err = u.authorRepository.SaveAuthor(author)
	return author, nil
}

// auth/create_user_use_case.go

type CreateUserInput struct {
  Name      string
	Email     string
	AvatarURL string
}

type CreateUserUseCase struct {
	userRepository    UserRepository
  saveAuthorUseCase *discussion.SaveAUthorUseCase
}

func (u *CreateUserUseCase) Run(input CreateUserInput) (*User, error) {
	// Handle user creation

	u.saveAuthorUseCase.Run(duscussion.SaveAuthorInput{
    Name:      user.name,
    Email:     user.Email
    AvatarURL: user.AvatarURL,
   })
	return user, nil
}
```

## Event-Based Communication

## Final thoughts

## References

- [BoundedContext](https://martinfowler.com/bliki/BoundedContext.html)
- [UbiquitousLanguage](https://martinfowler.com/bliki/UbiquitousLanguage.html)

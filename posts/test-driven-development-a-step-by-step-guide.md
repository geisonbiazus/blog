title: Test-Driven Development: A Step-By-Step Guide
author: Geison Biazus
description: In this post, I explain the benefits of TDD and show how to step-by-step apply it with a real-world example.
image_path: /static/image/logo-small.png
time: 2021-07-03 08:30
--
I have been practicing Test-Driven Development (TDD) in my career for 13 years at the moment of this post, and I can say for sure that there is no better practice for developing software. It brings me confidence in my code, a better code design, and allows me to focus on a small thing at a time.

TDD was brought by Kent Back as part of the Extreme Programming (XP) practices. It consists of first writing a test and then writing the code to make that test pass. TDD itself is simple to apply, but it requires some degree of discipline and practice to master it.

## But why test first?

If you first write your production code, then you write your test, there is no guarantee that you tested everything your production code does. Also, there is no guarantee that your test will fail if you have a bug in the code. By following the TDD cycle and laws, explained below, you make sure that every behavior of your code has is covered by a test, and you can refactor with total confidence.

## TDD Cycle

The TDD cycle, also known as "red / green /refactor" is the following:

1. Write a failing test (red)
2. Write the production code to make the test pass (green)
3. Refactor your code (refactor)
4. Repeat the process with the next test

In the first step, you write your test code. It is going to fail since there is no production code to make it pass yet. That's why this phase is called "red".

The second step is to write the production code, only sufficient to make your test pass. Don't think about performance or clean code at this phase. This is the "green" phase.

Now you are confident that both your test and production code are correct since you saw the test fail and then pass. Now is the time to refactor the code, if needed. At this phase (refactor), both production and test code can be refactored, as long as you do one at a time and keep running your tests to check if everything is still working as expected. As Robert C. Martin (“Uncle Bob”) says, "as tests get more specific, the code gets more generic". So you aim to ways make your production code the more generic as possible to fulfill the scenarios, while still covering a lot of specific scenarios in the tests.

Finally, you repeat the process by writing the next test until the feature is done.

## Three laws of TDD

Robert C. Martin (“Uncle Bob”) created three basic rules for when you are practicing TDD called [The three laws of TDD](https://www.oreilly.com/library/view/modern-c-programming/9781941222423/f_0055.html). These are a set of rules that when strictly followed will guarantee that your code will be tested thoroughly. They are the following:

1. You must not write any production code without having a failing test.
2. You must not write more of a test than is sufficient to fail, and not compiling is a failure.
3. You must not write more production code than is sufficient to make the currently failing test pass.

By following these rules you will work in really small cycles of test and production code. At first, as you should always start with a test, the cycles start by mostly fixing compilation errors as the classes and data structures still don't exist in the production code. But after the code evolves, you start fixing the assertion errors, until you have a complete passing test scenario. You should not create even data structures and interfaces without having a failing test that uses them. By doing so, you guarantee that all your code is fully functioning and covered by tests.

## Benefits

The greatest benefit of working with TDD is confidence. You should always have a suite of tests that you can trust with your life. If your tests pass, you deploy. It is as simple as that.

Having this confidence in your test allows you to refactor. Without tests refactoring is impossible. Without refactoring, you cannot keep a simple design. Without a simple design, you cannot react to changes. So, in my opinion, TDD is the main practice and the base of agile.

Another benefit is the code design. First because in the refactor phase you are always looking to make your code cleaner. But also, when you are test driving, you will face a lot of complexities. These complexities are usually low-level details that you can separate from your high-level policy. TDD helps you to see these details and create abstractions leaving them to be tackled later on other modules with their own set of tests.

## Real world example

Usually, you can find TDD examples of some [Code Kata](<https://en.wikipedia.org/wiki/Kata_(programming)>) exercises or things that are more focused on algorithms. But instead, I'll show you something closer to a real-life scenario in the format of a use case. For this example, I'll show how to implement the `View Post` use case of this blog.

The `View Post` use case works as follows: A user accesses a post URL from the browser. Based on the path we try to load a post file in the file system written in Markdown. If the post is not found, we display a "not found" error. If the post is found, we convert the Markdown to HTML and present the post to the user in the browser.

As we are talking about the use case here, we are going to ignore the low-level details completely, so no web server, no storage, no file system. Instead, we will focus on the high-level policy using TDD to drive this code and design, creating "ports" where the low-level details can be plugged in as "adapters".

If you are interested in how this blog is architectured, I wrote a post about it and you can see it here: [Applying Clean Architecture in Go](https://blog.geisonbiazus.com/posts/applying-clean-architecture-in-go).

This example is done using the Go programming language but the flow is the same for any language. The Go standard library is used to write the tests but it doesn't have assertion functions. As I like the expressivity of assertions, I created some [helper functions](https://github.com/geisonbiazus/blog/blob/main/pkg/assert/assert.go) that I'll make use of in the tests.

### Starting from nothing (Red Phase)

We have nothing implemented yet, so the first thing to do is to create a test file. So let's create a file called `view_post_use_case_test.go` under `internal/core/blog`. The directory choice is because this is a use case feature, meaning that it belongs to the core layer of the system and the blog component.

The first thing we need is a module to hold our use case code, so let's create a test to check the existence of this module:

```go
// internal/core/blog/view_post_use_case_test.go

package blog_test

import (
	"testing"

	"github.com/geisonbiazus/blog/internal/core/blog"
)

func TestViewPostUseCase(t *testing.T) {
	t.Run("It initializes", func(t *testing.T) {
		var _ *blog.ViewPostUseCase = blog.NewViewPostUseCase()
	})
}
```

This test has no assertions but it checks if the `blog.NewViewPostUseCase` returns a `*blog.ViewPostUseCase` using the go compiler. If we run the tests, we have this output:

```sh
$ go test ./...
# github.com/geisonbiazus/blog/internal/core/blog_test [github.com/geisonbiazus/blog/internal/core/blog.test]
internal/core/blog/virew_post_use_case_test.go:11:10: undefined: blog.ViewPostUseCase
internal/core/blog/virew_post_use_case_test.go:11:33: undefined: blog.NewViewPostUseCase
FAIL    github.com/geisonbiazus/blog/internal/core/blog [build failed]
```

The second law of TDD says that not compiling is failing, so now it's time to make this test pass.

### Starting from nothing (Green Phase)

Let's create a new file in the same directory called `view_post_use_case.go` with the following content:

```go
// internal/core/blog/view_post_use_case.go

package blog

type ViewPostUseCase struct{}

func NewViewPostUseCase() *ViewPostUseCase {
	return &ViewPostUseCase{}
}

```

Running the test now, everything compiles and the test pass:

```sh
$ go test ./...
ok      github.com/geisonbiazus/blog/internal/core/blog   0.195s
```

### Running the use case (Red Phase)

Now we need to think about how we want to execute this use case. The "view post" use case will be invoked whenever a user accesses a post URL in the browser. Based on the post path, we need to load the post data, generate the HTML from the Markdown, and display it to the user. So the input of this use case is the post `path`, and the return is a `RenderedPost` and an optional `error` in case something goes wrong. So let's add a new test to define the `Run` method:

```go
// internal/core/blog/view_post_use_case_test.go

import "github.com/geisonbiazus/blog/pkg/assert"

// ...

func TestViewPostUseCase(t *testing.T) {
	// ...

	t.Run("It runs", func(t *testing.T) {
		usecase := blog.NewViewPostUseCase()
		path := "path"
		renderedPost, err := usecase.Run(path)

		assert.Equal(t, blog.RenderedPost{}, renderedPost)
		assert.Nil(t, err)
	})
})

```

In this new scenario, we instantiate the use case the same way as the previous scenario, then the `Run` method is called with the post `path`. A `blog.RenderedPost` and an `error` are returned. Running the tests, they fail because both the `Run` method and the `blog.RenderedPost` type don't exist at the moment:

```sh
$ go test ./...
# github.com/geisonbiazus/blog/internal/core/blog_test [github.com/geisonbiazus/blog/internal/core/blog.test]
internal/core/blog/virew_post_use_case_test.go:18:31: usecase.Run undefined (type *blog.ViewPostUseCase has no field or method Run)
internal/core/blog/virew_post_use_case_test.go:20:19: undefined: blog.RenderedPost
FAIL    github.com/geisonbiazus/blog/internal/core/blog [build failed]
```

### Running the use case (Green Phase)

Again we have a failure due to a compilation error, so to make this scenario pass, we first implement the `Run` method in the `ViewPostUseCase`. Just enough implementation to make the test pass:

```go
// internal/core/blog/view_post_use_case.go

// ...

func (u *ViewPostUseCase) Run(path string) (RenderedPost, error) {
	return RenderedPost{}, nil
}
```

And in a new file called `entities.go` we create the `RenderedPost` struct. Again just enough to make the test pass. That's why this struct is still empty:

```go
// internal/core/blog/entities.go

package blog

type RenderedPost struct{}

```

All the tests pass now:

```sh
$ go test ./...
ok      github.com/geisonbiazus/blog/internal/core/blog   0.193s
```

### Running the use case (Refactor Phase)

Now it is time for our first refactor. There's still not much in the production code so let's skip that for now. But we have a small duplication in the tests. This duplication is the call to the constructor function. Although right now it is just one line, keeping it duplicated on multiple tests will make us need to update multiple tests when we start adding dependencies to this constructor. So, to avoid that let's create a setup function that returns some fixture data where we can share the dependencies between all the tests:

```go
// internal/core/blog/view_post_use_case_test.go

// ...

type viewPostUseCaseFixture struct {
	usecase *blog.ViewPostUseCase
}

func TestViewPostUseCase(t *testing.T) {
	setup := func() *viewPostUseCaseFixture {
		usecase := blog.NewViewPostUseCase()
		return &viewPostUseCaseFixture{
			usecase: usecase,
		}
	}

	t.Run("It initializes", func(t *testing.T) {
		setup()
	})

	t.Run("It runs", func(t *testing.T) {
		f := setup()

		path := "path"
		renderedPost, err := f.usecase.Run(path)

		assert.Equal(t, blog.RenderedPost{}, renderedPost)
		assert.Nil(t, err)
	})
}
```

Here we created the `viewPostUseCaseFixture` struct containing the use case instance. This struct is a [test fixture](https://en.wikipedia.org/wiki/Test_fixture), and it contains data to be used by the tests. Next, the setup function returns this fixture to any test that calls it. This way, we can keep adding other dependencies to the fixture instance that can be shared between multiple tests without changing the tests themselves. As the setup function is called at the beginning of each test, the tests are still isolated having different instances on each one of them.

Running the tests we can see that they still pass.

```sh
$ go test ./...
ok      github.com/geisonbiazus/blog/internal/core/blog   0.172s
```

Still on this refactor phase, you can notice that the scenario `It initializes`, is not doing anything anymore. This scenario can now be removed as its behavior is being covered by the other test. This kind of scenario is called stair-step test. It exists with the purpose of writing a little bit of production code allowing us to write the next test. After it fulfilled its purpose, it can be safely removed. So let's remove this scenario and we'll end up having only one scenario again:

```go
// internal/core/blog/view_post_use_case_test.go

// ...

func TestViewPostUseCase(t *testing.T) {
	setup := func() *viewPostUseCaseFixture {
		// ...
	}

	t.Run("It runs", func(t *testing.T) {
		// ...
	})
}
```

### Post not found (Red Phase)

Now let's start implementing the behavior of the `Run` method. The first simplest scenario we can write is when the post is not found. So let's add a scenario like this:

```go
// internal/core/blog/view_post_use_case_test.go

// ...
func TestViewPostUseCase(t *testing.T) {
	// ...

	t.Run("It returns error when post is not found", func(t *testing.T) {
		f := setup()

		f.repo.ReturnError = blog.ErrPostNotFound

		renderedPost, err := f.usecase.Run("path")

		assert.Equal(t, "path", f.repo.ReceivedPath)
		assert.Equal(t, blog.RenderedPost{}, renderedPost)
		assert.Equal(t, blog.ErrPostNotFound, err)
	})
}
```

Again some new things here. Fist is the `f.repo.ReturnError`. This is a new dependency that will be added to the setup function and will return a test spy. Here we are saying that this spy should return the error `blog.ErrPostNotFound`. This is also a new type we need to create. Then, after calling the `Run` method, we again go to the repo spy to check if the "received path" is the same path given in the `Run` method. Finally, we check the method response and if the returned error is the same error returned by the `repo` spy. This will become more clear with the spy implementation.

Running the tests now, they will fail by compilation due to the lack of these new things:

```sh
$ go test ./...
# github.com/geisonbiazus/blog/internal/core/blog_test [github.com/geisonbiazus/blog/internal/core/blog.test]
internal/core/blog/virew_post_use_case_test.go:39:4: f.repo undefined (type *viewPostUseCaseFixture has no field or method repo)
internal/core/blog/virew_post_use_case_test.go:39:24: undefined: blog.ErrPostNotFound
internal/core/blog/virew_post_use_case_test.go:43:28: f.repo undefined (type *viewPostUseCaseFixture has no field or method repo)
internal/core/blog/virew_post_use_case_test.go:45:19: undefined: blog.ErrPostNotFound
FAIL    github.com/geisonbiazus/blog/internal/core/blog [build failed]
```

Still on the red phase, let's make it compile by creating the new required types:

```go
// internal/core/blog/view_post_use_case_test.go

// ...
type viewPostUseCaseFixture struct {
	usecase *blog.ViewPostUseCase
	repo    *PostRepoSpy
}

func TestViewPostUseCase(t *testing.T) {
	setup := func() *viewPostUseCaseFixture {
		repo := NewPostRepoSpy()
		usecase := blog.NewViewPostUseCase()

		return &viewPostUseCaseFixture{
			usecase: usecase,
			repo:    repo,
		}
	}
	// ...
}

type PostRepoSpy struct {
	ReturnError  error
	ReceivedPath string
}

func NewPostRepoSpy() *PostRepoSpy {
	return &PostRepoSpy{}
}

func (r *PostRepoSpy) GetPostByPath(path string) error {
	r.ReceivedPath = path
	return r.ReturnError
}
```

First we added the `repo` attribute to the `viewPostUseCaseFixture` struct of the type `*PostRepoSpy`. Then in the setup function, this repo is initialized using the `NewPostRepoSpy` constructor. At the end of the file, we added the `PostRepoSpy` implementation. It contains a method `GetPostByPath` that returns whatever is configured in this spy. It also stores the received `path` argument making it accessible by the test to assert its value later. A spy is one of the multiple types of [test doubles](https://en.wikipedia.org/wiki/Test_double) and its purpose is to be able to record all the received arguments that later can be checked in the test to guarantee correctness. I'll make a post about the multiple types of test doubles in the future.

We are now dealing with the first big design decision about the system architecture. We know that the post should be loaded from somewhere. In the case of this feature, it will be loaded from the file system. But for the use case implementation, the place where it is loaded is a low-level detail, and the use cases deal only with the high-level policies. This is where we are crossing an architectural boundary and we are going to plug in an adapter to fulfill this behavior. And the same goes for the tests. The adapter we are plugging in here is the `PostRepoSpy`.

The last thing required to make our code compile is the error type, so in the `entities.go` file we can add the following:

```go
// internal/core/blog/entities.go

import "errors"

var ErrPostNotFound = errors.New("post not found")
```

Now when we run the tests, and instead of failing by compilation error they fail because of the assertions:

```sh
$ go test ./...
--- FAIL: TestViewPostUseCase (0.00s)
    --- FAIL: TestViewPostUseCase/It_returns_error_when_post_is_not_found (0.00s)
        virew_post_use_case_test.go:47:
            expected: path
              actual:
        virew_post_use_case_test.go:49:
            expected: post not found
              actual: <nil>
FAIL
FAIL    github.com/geisonbiazus/blog/internal/core/blog   0.180s
```

### Post not found (Green Phase)

Now to make our test pass we need to try loading the post from a repository. In the test context, this repository is a spy that will always return an error but the production code cannot depend upon a spy. So what we are going to create here is an interface. This interface goes to a new file called `ports.go`. The purpose of this file is to hold all the interfaces of this `blog` component where an adapter will be plugged in. In this file we create the `PostRepo` interface like the following:

```go
// internal/core/blog/ports.go

package blog

type PostRepo interface {
	GetPostByPath(path string) error
}
```

Notice that the `GetPostByPath` method only returns an error and no post. This happens because we are still in an intermediate state and we will come back to fix it later.

And now we can do the changes in the `ViewPostUseCase` to satisfy the tests:

```go
// internal/core/blog/view_post_use_case.go

// ...
type ViewPostUseCase struct {
	postRepo PostRepo
}

func NewViewPostUseCase(postRepo PostRepo) *ViewPostUseCase {
	return &ViewPostUseCase{postRepo: postRepo}
}

func (u *ViewPostUseCase) Run(path string) (RenderedPost, error) {
	err := u.postRepo.GetPostByPath(path)
	return RenderedPost{}, err
}
```

First, we added the `PostRepo` interface to the struct and as a parameter in the constructor function. Then in the `Run` method, we call the `GetPostByPath` to get the error back and return it.

The last thing missing is to adjust the test setup function to pass the `PostRepoSpy` to the `blog.NewViewPostUseCase` constructor:

```go
// internal/core/blog/view_post_use_case_test.go

// ...
func TestViewPostUseCase(t *testing.T) {
	setup := func() *viewPostUseCaseFixture {
		repo := NewPostRepoSpy()
		usecase := blog.NewViewPostUseCase(repo)
		//...
	}
	// ...
}
```

Now when we run the tests, all of them pass again:

```sh
$ go test ./...
ok      github.com/geisonbiazus/blog/internal/core/blog   0.183s
```

### Post not found (Refactor Phase)

The same thing we did before with the `It initializes` test scenario, we can do with the `It runs`. This test is also a stair-step test and its behavior is covered by other scenarios, so it can be removed keeping only one scenario again:

```go
// internal/core/blog/view_post_use_case_test.go

// ...

func TestViewPostUseCase(t *testing.T) {
	setup := func() *viewPostUseCaseFixture {
		// ...
	}

	t.Run("It returns error when post is not found", func(t *testing.T) {
		// ...
	})
}
```

### The happy path (Red Phase)

The next scenario is the happy path. Given that a post is found in the repository, then we need to render this post and return the rendered post. So let's add the following test case:

```go
// internal/core/blog/view_post_use_case_test.go

import "time"

// ...

func TestViewPostUseCase(t *testing.T) {
	// ...
	t.Run("It returns a rendered post when post is found", func(t *testing.T) {
		f := setup()

		postTime, _ := time.Parse(time.RFC3339, "2021-04-03T00:00:00+00:00")
		post := blog.Post{
			Title:    "Title",
			Author:   "Author",
			Time:     postTime,
			Path:     "path",
			Markdown: "content",
		}

		f.repo.ReturnPost = post
		f.renderer.ReturnRenderedContent = "Rendered content"

		renderedPost, err := f.usecase.Run(post.Path)

		assert.Equal(t, post.Path, f.repo.ReceivedPath)
		assert.Equal(t, post.Markdown, f.renderer.ReceivedContent)
		assert.Nil(t, err)
		assert.Equal(t, blog.RenderedPost{
			Post: post,
			HTML: "Rendered content",
		}, renderedPost)
	})
}
```

Here we are introducing a new entity, the `blog.Post`, and telling the `PostRepoSpy` to return it instead of the error from the previous scenario. Next, we configure a new spy called `renderer` to return the rendered content. Then we run the use case normally and finally, we assert that the spy received the correct arguments and that the returned `blog.RenderedPost` contains the correct fields.

So to make it compile, let's first add the new `Post` entity and update the `RenderedPost` with the new fields in the `entities.go` file:

```go
// internal/core/blog/entities.go

import "time"

type Post struct {
	Title    string
	Author   string
	Time     time.Time
	Path     string
	Markdown string
}

type RenderedPost struct {
	Post Post
	HTML string
}

// ...
```

In the `view_post_use_case_test.go` file, let's update the `PostRepoSpy` with the `ReturnPost` field:

```go
// internal/core/blog/view_post_use_case_test.go

// ...

type PostRepoSpy struct {
	ReturnPost   blog.Post
	ReturnError  error
	ReceivedPath string
}

// ...
```

Still on the same file, we add the new `renderer` field to the text fixture and assign it in the setup function.

```go
// internal/core/blog/view_post_use_case_test.go

// ...
type viewPostUseCaseFixture struct {
	usecase  *blog.ViewPostUseCase
	repo     *PostRepoSpy
	renderer *RendererSpy
}

func TestViewPostUseCase(t *testing.T) {
	setup := func() *viewPostUseCaseFixture {
		repo := NewPostRepoSpy()
		renderer := NewRendererSpy()
		usecase := blog.NewViewPostUseCase(repo)

		return &viewPostUseCaseFixture{
			usecase:  usecase,
			repo:     repo,
			renderer: renderer,
		}
	}
	// ...
}
```

At the end of the file, let's add the `RendererSpy` implementation. It contains a `Render` method that receives a string and returns a rendered string and a possible error.

```go
// internal/core/blog/view_post_use_case_test.go

// ...

type RendererSpy struct {
	ReturnRenderedContent string
	ReturnError           error
	ReceivedContent       string
}

func NewRendererSpy() *RendererSpy {
	return &RendererSpy{}
}

func (r *RendererSpy) Render(content string) (string, error) {
	r.ReceivedContent = content
	return r.ReturnRenderedContent, r.ReturnError
}
```

Now when we run the tests they compile and we have the following failure:

```sh
$ go test ./...
--- FAIL: TestViewPostUseCase (0.00s)
    --- FAIL: TestViewPostUseCase/It_returns_a_rendered_post_when_post_is_found (0.00s)
        virew_post_use_case_test.go:74:
            expected: content
              actual:
        virew_post_use_case_test.go:76:
            expected: {{Title Author 2021-04-03 00:00:00 +0000 +0000 path content} Rendered content}
              actual: {{  0001-01-01 00:00:00 +0000 UTC  } }
FAIL
FAIL    github.com/geisonbiazus/blog/internal/core/blog   0.383s
```

### The happy path (Green Phase)

Now to make this test pass we do the following changes in the `ViewPostUseCase`:

```go
// internal/core/blog/view_post_use_case.go

type ViewPostUseCase struct {
	postRepo PostRepo
	renderer Renderer
}

func NewViewPostUseCase(postRepo PostRepo, renderer Renderer) *ViewPostUseCase {
	return &ViewPostUseCase{postRepo: postRepo, renderer: renderer}
}

func (u *ViewPostUseCase) Run(path string) (RenderedPost, error) {
	post, err := u.postRepo.GetPostByPath(path)

	if err != nil {
		return RenderedPost{}, err
	}

	html, _ := u.renderer.Render(post.Markdown)

	return RenderedPost{
		Post: post,
		HTML: html,
	}, nil
}
```

We added the `Renderer` dependency to the struct as well as to the constructor function. Then we updated the implementation of the `Run` method to first check for any error returned by the `PostRepo` to keep the previous test passing. Following from there, we call the `Render` method from the `Renderer` dependency, passing the post markdown as an argument. Finally, we build and return the `RenderedPost`.

These changes broke the compilation again, since we added a new dependency to the `ViewPostUseCase`, also a new value is returned from the `PostRepo`. So let's fix these errors now:

```go
// internal/core/blog/ports.go

type PostRepo interface {
	GetPostByPath(path string) (Post, error)
}

type Renderer interface {
	Render(content string) (string, error)
}
```

We updated the `PostRepo` interface to also return a `Post` in the `GetPostByPath` method. Also, we added the new `Renderer` interface.

Next, we need to do some fixes in the tests too:

```go
// internal/core/blog/view_post_use_case_test.go

// ...
func TestViewPostUseCase(t *testing.T) {
	setup := func() *viewPostUseCaseFixture {
		repo := NewPostRepoSpy()
		renderer := NewRendererSpy()
		usecase := blog.NewViewPostUseCase(repo, renderer)

		// ...
	}
	// ...
}
// ...

func (r *PostRepoSpy) GetPostByPath(path string) (blog.Post, error) {
	r.ReceivedPath = path
	return r.ReturnPost, r.ReturnError
}

// ...
```

We passed the renderer to the `ViewPostUseCase` constructor and updated the `PostRepoSpy`to to satisfy the changes to its interface.

Now when we run the tests, we have all of them passing again:

```sh
$ go test ./...
ok      github.com/geisonbiazus/blog/internal/core/blog   0.183s
```

### The happy path (Refactor Phase)

Now it is the refactor phase again. First, for the production code, we can make the `Run` method shorter by extracting the rendering code into a `renderPost` method like the following:

```go
// internal/core/blog/view_post_use_case.go

// ...
func (u *ViewPostUseCase) Run(path string) (RenderedPost, error) {
	post, err := u.postRepo.GetPostByPath(path)

	if err != nil {
		return RenderedPost{}, err
	}

	return u.renderPost(post)
}

func (u *ViewPostUseCase) renderPost(post Post) (RenderedPost, error) {
	html, _ := u.renderer.Render(post.Markdown)

	return RenderedPost{
		Post: post,
		HTML: html,
	}, nil
}
```

In the tests we can also make the test code shorter by extracting the building of the post into a separated function:

```go
// internal/core/blog/view_post_use_case_test.go

// ...
func TestViewPostUseCase(t *testing.T) {
	// ...
	t.Run("It returns a rendered post when post is found", func(t *testing.T) {
		f := setup()

		post := newPost()

		f.repo.ReturnPost = post
		f.renderer.ReturnRenderedContent = "Rendered content"

		renderedPost, err := f.usecase.Run(post.Path)

		assert.Equal(t, post.Path, f.repo.ReceivedPath)
		assert.Equal(t, post.Markdown, f.renderer.ReceivedContent)
		assert.Nil(t, err)
		assert.Equal(t, blog.RenderedPost{
			Post: post,
			HTML: "Rendered content",
		}, renderedPost)
	})
}

func newPost() blog.Post {
	postTime, _ := time.Parse(time.RFC3339, "2021-04-03T00:00:00+00:00")
	return blog.Post{
		Title:    "Title",
		Author:   "Author",
		Time:     postTime,
		Path:     "path",
		Markdown: "content",
	}
}
// ...
```

Running the tests now we can see that they keep passing and we are done with this refactor phase.

```sh
$ go test ./...
ok      github.com/geisonbiazus/blog/internal/core/blog   0.183s
```

### One last error check (Red Phase)

There is still one error we are ignoring. That is when the markdown content fails to render. So let's add the new test case to handle that situation:

```go
// internal/core/blog/view_post_use_case_test.go

// ...
func TestViewPostUseCase(t *testing.T) {
	// ...
	t.Run("It returns error when post fails to render", func(t *testing.T) {
		f := setup()

		post := newPost()
		f.repo.ReturnPost = post
		f.renderer.ReturnError = errors.New("render error")

		renderedPost, err := f.usecase.Run(post.Path)

		assert.Equal(t, f.renderer.ReturnError, err)
		assert.Equal(t, blog.RenderedPost{}, renderedPost)
	})
}
// ...
```

Here we tell the `RendererSpy` to return an error when it is called, and then we check if the return of the `Run` method returns the same error returned from the `RendererSpy`. The second assertion is just to check that on the error case we return an empty `blog.RenderedPost`. By running the tests now, we have the following failure:

```sh
$ go test ./...
--- FAIL: TestViewPostUseCase (0.00s)
    --- FAIL: TestViewPostUseCase/It_returns_error_when_post_fails_to_render (0.00s)
        virew_post_use_case_test.go:85:
            expected: render error
              actual: <nil>
        virew_post_use_case_test.go:86:
            expected: {{  0001-01-01 00:00:00 +0000 UTC  } }
              actual: {{Title Author 2021-04-03 00:00:00 +0000 +0000 path content} }
FAIL
FAIL    github.com/geisonbiazus/blog/internal/core/blog   0.263s
```

### One last error check (Green Phase)

We make this test pass by simply checking for an error and returning it in the `renderPost` method:

```go
// internal/core/blog/view_post_use_case.go

// ...
func (u *ViewPostUseCase) renderPost(post Post) (RenderedPost, error) {
	html, err := u.renderer.Render(post.Markdown)

	if err != nil {
		return RenderedPost{}, err
	}

	return RenderedPost{
		Post: post,
		HTML: html,
	}, nil
}
```

Running the tests now, we see that all of them pass again:

```sh
$ go test ./...
ok      github.com/geisonbiazus/blog/internal/core/blog   0.182s
```

### Final result

And that's it. We are done with the `ViewPostUseCase` implementation. You can see all the files we created on this post next.

The first file is the `view_post_use_case.go` where the `ViewPostUseCase` is implemented. It receives a `PostRepo` and a `Renderer` which are the dependencies used by fetching the post data and rendering it into HTML.

```go
// internal/core/blog/view_post_use_case.go

package blog

type ViewPostUseCase struct {
	postRepo PostRepo
	renderer Renderer
}

func NewViewPostUseCase(postRepo PostRepo, renderer Renderer) *ViewPostUseCase {
	return &ViewPostUseCase{postRepo: postRepo, renderer: renderer}
}

func (u *ViewPostUseCase) Run(path string) (RenderedPost, error) {
	post, err := u.postRepo.GetPostByPath(path)

	if err != nil {
		return RenderedPost{}, err
	}

	return u.renderPost(post)
}

func (u *ViewPostUseCase) renderPost(post Post) (RenderedPost, error) {
	html, err := u.renderer.Render(post.Markdown)

	if err != nil {
		return RenderedPost{}, err
	}

	return RenderedPost{
		Post: post,
		HTML: html,
	}, nil
}
```

The next file is `entities.go`. Here are the entities that represent the post and its rendered version.

```go
// internal/core/blog/entities.go
package blog

import (
	"errors"
	"time"
)

type Post struct {
	Title    string
	Author   string
	Time     time.Time
	Path     string
	Markdown string
}

type RenderedPost struct {
	Post Post
	HTML string
}

var ErrPostNotFound = errors.New("post not found")
```

The `ports.go` file contains the interfaces where adapters are plugged in. These are the interfaces that the `ViewPostUseCase` depends on and are supposed to be implemented by other layers of the application.

```go
// internal/core/blog/ports.go

package blog

type PostRepo interface {
	GetPostByPath(path string) (Post, error)
}

type Renderer interface {
	Render(content string) (string, error)
}
```

And the last file is the `view_post_use_case_test.go` which contains all the tests for the `ViewPostUseCase`. It guarantees that everything else exists and behaves correctly.

```go
// internal/core/blog/view_post_use_case_test.go

package blog_test

import (
	"errors"
	"testing"
	"time"

	"github.com/geisonbiazus/blog/internal/core/blog"
	"github.com/geisonbiazus/blog/pkg/assert"
)

type viewPostUseCaseFixture struct {
	usecase  *blog.ViewPostUseCase
	repo     *PostRepoSpy
	renderer *RendererSpy
}

func TestViewPostUseCase(t *testing.T) {
	setup := func() *viewPostUseCaseFixture {
		repo := NewPostRepoSpy()
		renderer := NewRendererSpy()
		usecase := blog.NewViewPostUseCase(repo, renderer)

		return &viewPostUseCaseFixture{
			usecase:  usecase,
			repo:     repo,
			renderer: renderer,
		}
	}

	t.Run("It returns error when post is not found", func(t *testing.T) {
		f := setup()

		f.repo.ReturnError = blog.ErrPostNotFound

		renderedPost, err := f.usecase.Run("path")

		assert.Equal(t, "path", f.repo.ReceivedPath)
		assert.Equal(t, blog.RenderedPost{}, renderedPost)
		assert.Equal(t, blog.ErrPostNotFound, err)
	})

	t.Run("It returns a rendered post when post is found", func(t *testing.T) {
		f := setup()

		post := newPost()

		f.repo.ReturnPost = post
		f.renderer.ReturnRenderedContent = "Rendered content"

		renderedPost, err := f.usecase.Run(post.Path)

		assert.Equal(t, post.Path, f.repo.ReceivedPath)
		assert.Equal(t, post.Markdown, f.renderer.ReceivedContent)
		assert.Nil(t, err)
		assert.Equal(t, blog.RenderedPost{
			Post: post,
			HTML: "Rendered content",
		}, renderedPost)
	})

	t.Run("It returns error when post fails to render", func(t *testing.T) {
		f := setup()

		post := newPost()
		f.repo.ReturnPost = post
		f.renderer.ReturnError = errors.New("render error")

		renderedPost, err := f.usecase.Run(post.Path)

		assert.Equal(t, f.renderer.ReturnError, err)
		assert.Equal(t, blog.RenderedPost{}, renderedPost)
	})
}

func newPost() blog.Post {
	postTime, _ := time.Parse(time.RFC3339, "2021-04-03T00:00:00+00:00")
	return blog.Post{
		Title:    "Title",
		Author:   "Author",
		Time:     postTime,
		Path:     "path",
		Markdown: "content",
	}
}

type PostRepoSpy struct {
	ReturnPost   blog.Post
	ReturnError  error
	ReceivedPath string
}

func NewPostRepoSpy() *PostRepoSpy {
	return &PostRepoSpy{}
}

func (r *PostRepoSpy) GetPostByPath(path string) (blog.Post, error) {
	r.ReceivedPath = path
	return r.ReturnPost, r.ReturnError
}

type RendererSpy struct {
	ReturnRenderedContent string
	ReturnError           error
	ReceivedContent       string
}

func NewRendererSpy() *RendererSpy {
	return &RendererSpy{}
}

func (r *RendererSpy) Render(content string) (string, error) {
	r.ReceivedContent = content
	return r.ReturnRenderedContent, r.ReturnError
}
```

## Final thoughts

In this post, we saw how we can implement and design our application in small steps using Test-Driven Development. These tests follow the laws of TDD that guarantee that 100% of our production code is covered by tests, giving us the confidence to trust them and if they pass, we deploy without fear.

TDD also helps us to identify the application boundaries and hide the low-level details behind interfaces and focus on the high-level policies of the application. Of course that the low-level details also have their tests but they are separated in a way that they can be developed and tested in isolation, and more important, extended or replaced if needed.

This post only shows the use case implementation. If you are interested in how the whole blog is architectured, you can see the post [Applying Clean Architecture in Go](https://blog.geisonbiazus.com/posts/applying-clean-architecture-in-go) and check out the [blog source code](https://github.com/geisonbiazus/blog) where all the layers are implemented and tested.

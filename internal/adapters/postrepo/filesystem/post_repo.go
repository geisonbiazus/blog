package filesystem

import (
	"io/ioutil"

	"github.com/geisonbiazus/blog/internal/core/posts"
)

type PostRepo struct {
	BasePath string
}

func NewPostRepo(basePath string) *PostRepo {
	return &PostRepo{BasePath: basePath}
}

func (r *PostRepo) GetPostByPath(path string) (posts.Post, error) {
	content, err := ioutil.ReadFile(r.BasePath + "/" + path + ".md")

	if err != nil {
		return posts.Post{}, posts.ErrPostNotFound
	}

	post, err := ParseFileContent(string(content))
	post.Path = path

	return post, err
}

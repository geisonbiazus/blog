package postrepo

import (
	"io/ioutil"

	"github.com/geisonbiazus/blog/internal/core/posts"
)

type FileSystemPostRepo struct {
	BasePath string
}

func NewFileSystemPostRepo(basePath string) *FileSystemPostRepo {
	return &FileSystemPostRepo{BasePath: basePath}
}

func (r *FileSystemPostRepo) GetPostByPath(path string) (posts.Post, error) {
	content, err := ioutil.ReadFile(r.BasePath + "/" + path + ".md")

	if err != nil {
		return posts.Post{}, posts.ErrPostNotFound
	}

	post, err := ParseFileContent(string(content))
	post.Path = path

	return post, err
}

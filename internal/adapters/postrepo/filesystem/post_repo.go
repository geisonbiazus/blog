package filesystem

import (
	"io/fs"
	"io/ioutil"
	"os"
	"sort"
	"strings"

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

func (r *PostRepo) GetAllPosts() ([]posts.Post, error) {
	postList := []posts.Post{}
	entries, err := os.ReadDir(r.BasePath)

	if err != nil {
		return postList, err
	}

	for _, entry := range entries {
		postList, err = r.maybeLoadPostFromFile(postList, entry)
	}

	sort.Slice(postList, func(i, j int) bool {
		return postList[i].Time.After(postList[j].Time)
	})

	return postList, err
}

func (r *PostRepo) maybeLoadPostFromFile(postList []posts.Post, entry fs.DirEntry) ([]posts.Post, error) {
	if !strings.HasSuffix(entry.Name(), ".md") {
		return postList, nil
	}

	fileName := strings.TrimSuffix(entry.Name(), ".md")
	post, err := r.GetPostByPath(fileName)

	if err != nil {
		return postList, err
	}

	return append(postList, post), nil
}

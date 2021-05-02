package filesystem

import (
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/geisonbiazus/blog/internal/core/blog"
)

type PostRepo struct {
	BasePath string
}

func NewPostRepo(basePath string) *PostRepo {
	return &PostRepo{BasePath: basePath}
}

func (r *PostRepo) GetPostByPath(path string) (blog.Post, error) {
	content, err := ioutil.ReadFile(filepath.Join(r.BasePath, path+".md"))

	if err != nil {
		return blog.Post{}, blog.ErrPostNotFound
	}

	post, err := ParseFileContent(string(content))
	post.Path = path

	return post, err
}

func (r *PostRepo) GetAllPosts() ([]blog.Post, error) {
	posts := []blog.Post{}
	entries, err := os.ReadDir(r.BasePath)

	if err != nil {
		return posts, err
	}

	for _, entry := range entries {
		posts = r.maybeLoadPostFromFile(posts, entry)
	}

	return r.sortPostsByTimeDesc(posts), err
}

func (r *PostRepo) maybeLoadPostFromFile(posts []blog.Post, entry fs.DirEntry) []blog.Post {
	if !strings.HasSuffix(entry.Name(), ".md") {
		return posts
	}

	fileName := strings.TrimSuffix(entry.Name(), ".md")
	post, err := r.GetPostByPath(fileName)

	if err != nil {
		log.Printf("WARNING: error loading post \"%s\": %v", fileName, err)

		return posts
	}

	return append(posts, post)
}

func (r *PostRepo) sortPostsByTimeDesc(posts []blog.Post) []blog.Post {
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Time.After(posts[j].Time)
	})

	return posts
}

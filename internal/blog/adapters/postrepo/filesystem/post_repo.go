package filesystem

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/geisonbiazus/blog/internal/blog/entities"
)

type PostRepo struct {
	BasePath string
}

func NewPostRepo(basePath string) *PostRepo {
	return &PostRepo{BasePath: basePath}
}

func (r *PostRepo) GetPostByPath(path string) (entities.Post, error) {
	content, err := os.ReadFile(filepath.Join(r.BasePath, path+".md"))

	if err != nil {
		return entities.Post{}, entities.ErrPostNotFound
	}

	post, err := ParseFileContent(string(content))
	post.Path = path

	return post, err
}

func (r *PostRepo) GetAllPosts() ([]entities.Post, error) {
	posts := []entities.Post{}
	entries, err := os.ReadDir(r.BasePath)

	if err != nil {
		return posts, err
	}

	for _, entry := range entries {
		posts = r.maybeLoadPostFromFile(posts, entry)
	}

	return r.sortPostsByTimeDesc(posts), err
}

func (r *PostRepo) maybeLoadPostFromFile(posts []entities.Post, entry fs.DirEntry) []entities.Post {
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

func (r *PostRepo) sortPostsByTimeDesc(posts []entities.Post) []entities.Post {
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Time.After(posts[j].Time)
	})

	return posts
}

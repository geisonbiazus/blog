package post_repo

import "github.com/geisonbiazus/blog/internal/core/posts"

type FileSystemPostRepo struct{}

func NewFileSystemPostRepo() *FileSystemPostRepo {
	return &FileSystemPostRepo{}
}

func (r *FileSystemPostRepo) GetPostByPath(path string) (posts.Post, error) {
	return posts.Post{}, posts.ErrPostNotFound
}

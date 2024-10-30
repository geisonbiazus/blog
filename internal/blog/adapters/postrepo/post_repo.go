package postrepo

import "github.com/geisonbiazus/blog/internal/blog/adapters/postrepo/filesystem"

func NewFileSystemPostRepo(basePath string) *filesystem.PostRepo {
	return filesystem.NewPostRepo(basePath)
}

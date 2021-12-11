package postrepo

import "github.com/geisonbiazus/blog/internal/adapters/postrepo/filesystem"

func NewFileSystemPostRepo(basePath string) *filesystem.PostRepo {
	return filesystem.NewPostRepo(basePath)
}

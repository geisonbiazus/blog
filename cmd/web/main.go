package main

import (
	"github.com/geisonbiazus/blog/internal/adapters/postrepo"
	"github.com/geisonbiazus/blog/internal/adapters/renderer"
	"github.com/geisonbiazus/blog/internal/core/posts"
	"github.com/geisonbiazus/blog/internal/ui/web"
)

func main() {
	renderer := renderer.NewGoldmarkRenderer()
	postRepo := postrepo.NewFileSystemPostRepo("internal/adapters/postrepo/test_posts")
	viewPostUseCase := posts.NewVewPostUseCase(postRepo, renderer)
	server := web.NewServer(3000, "web/template", viewPostUseCase)
	server.Start()
}

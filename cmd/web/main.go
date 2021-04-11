package main

import (
	"log"

	"github.com/geisonbiazus/blog/internal/adapters/postrepo/filesystem"
	"github.com/geisonbiazus/blog/internal/adapters/renderer/goldmark"
	"github.com/geisonbiazus/blog/internal/core/posts"
	"github.com/geisonbiazus/blog/internal/ui/web"
)

func main() {
	renderer := goldmark.NewGoldmarkRenderer()
	postRepo := filesystem.NewPostRepo("posts")
	viewPostUseCase := posts.NewVewPostUseCase(postRepo, renderer)
	server := web.NewServer(3000, "web/template", viewPostUseCase)
	log.Fatal(server.Start())
}

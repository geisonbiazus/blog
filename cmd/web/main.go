package main

import (
	"log"

	"github.com/geisonbiazus/blog/internal/app"
)

func main() {
	c := app.NewContext()
	log.Fatal(c.WebServer().Start())
}

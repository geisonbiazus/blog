package main

import (
	"log"

	"github.com/geisonbiazus/blog/internal/app"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	c := app.NewContext()
	log.Fatal(c.WebServer().Start())
}

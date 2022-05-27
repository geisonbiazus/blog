package main

import (
	"log"

	"github.com/geisonbiazus/blog/internal/app"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	c := app.NewContext()

	err := c.Migration().Up()
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(c.WebServer().Start())
}

package ports

import "github.com/geisonbiazus/blog/pkg/eventing/entities"

type Publisher interface {
	Publish(event entities.Event) error
}

package gen

import (
	"github.com/geisonbiazus/blog/pkg/gen/fake"
	"github.com/geisonbiazus/blog/pkg/gen/uuid"
)

type Generator interface {
	Generate() string
}

func NewUUIDGenerator() *uuid.Generator {
	return uuid.NewGenerator()
}

func NewFakeGenerator() *fake.Generator {
	return fake.NewGenerator()
}

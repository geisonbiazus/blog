package idgenerator

import (
	"github.com/geisonbiazus/blog/internal/adapters/idgenerator/fake"
	"github.com/geisonbiazus/blog/internal/adapters/idgenerator/uuid"
)

func NewUUIDGenerator() *uuid.Generator {
	return uuid.NewGenerator()
}

func NewFakeIDGenerator() *fake.IDGenerator {
	return fake.NewIDGenerator()
}

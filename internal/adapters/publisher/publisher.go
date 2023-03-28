package publisher

import "github.com/geisonbiazus/blog/internal/adapters/publisher/fake"

func NewFakePublisher() *fake.Publisher {
	return fake.NewPublisher()
}

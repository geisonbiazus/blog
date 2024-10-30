package publisher

import "github.com/geisonbiazus/blog/internal/shared/adapters/publisher/fake"

func NewFakePublisher() *fake.Publisher {
	return fake.NewPublisher()
}

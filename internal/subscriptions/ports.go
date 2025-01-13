package subscriptions

import (
	"context"

	"github.com/geisonbiazus/blog/internal/discussion"
	"github.com/geisonbiazus/blog/pkg/eventing"
)

type Subscriber interface {
	Subscribe(eventType string) chan eventing.Event
	NotifyError(event eventing.Event, err error)
	NotifySuccess(event eventing.Event)
}

type UseCases struct {
	SaveAuthor SaveAuthorUseCase
}

type SaveAuthorUseCase interface {
	Run(ctx context.Context, input discussion.SaveAuthorInput) (author *discussion.Author, err error)
}

package subscriptions

import (
	"context"

	"github.com/geisonbiazus/blog/internal/core/discussion"
	"github.com/geisonbiazus/blog/internal/core/shared"
)

type Subscriber interface {
	Subscribe(eventType string) chan shared.Event
	NotifyError(event shared.Event, err error)
	NotifySuccess(event shared.Event)
}

type UseCases struct {
	SaveAuthor SaveAuthorUseCase
}

type SaveAuthorUseCase interface {
	Run(ctx context.Context, input discussion.SaveAuthorInput) (author *discussion.Author, err error)
}

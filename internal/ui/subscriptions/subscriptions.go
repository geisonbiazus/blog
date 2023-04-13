package subscriptions

type Subscriptions struct {
	subscriber Subscriber
	usecases   *UseCases
}

func New(subscriber Subscriber, usecases *UseCases) *Subscriptions {
	return &Subscriptions{subscriber: subscriber, usecases: usecases}
}

func (s *Subscriptions) Start() {
	NewSaveAuthorSubscriber(s.usecases.SaveAuthor, s.subscriber).Start()
	NewUpdateAuthorSubscriber(s.usecases.SaveAuthor, s.subscriber).Start()
}

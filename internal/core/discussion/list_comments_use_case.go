package discussion

type ListCommentsUseCase struct{}

func NewListCommentsUseCase() *ListCommentsUseCase {
	return &ListCommentsUseCase{}
}

func (u *ListCommentsUseCase) Run(subjectId string) []RenderedComment {
	return []RenderedComment{}
}

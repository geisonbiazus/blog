package discussion

type ListCommentsUseCase struct {
	commentRepo CommentRepo
}

func NewListCommentsUseCase(commentRepo CommentRepo) *ListCommentsUseCase {
	return &ListCommentsUseCase{commentRepo}
}

func (u *ListCommentsUseCase) Run(subjectID string) ([]Comment, error) {
	return u.commentRepo.GetCommentsBySubjectID(subjectID)
}

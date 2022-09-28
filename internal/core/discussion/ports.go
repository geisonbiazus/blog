package discussion

type CommentRepo interface {
	GetCommentsBySubjectID(subjectID string) ([]Comment, error)
}

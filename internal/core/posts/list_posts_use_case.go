package posts

type ListPostsUseCase struct {
	postRepo PostRepo
}

func NewListPostsUseCase(postRepo PostRepo) *ListPostsUseCase {
	return &ListPostsUseCase{postRepo: postRepo}
}

func (u *ListPostsUseCase) Run() ([]Post, error) {
	return u.postRepo.GetAllPosts()
}

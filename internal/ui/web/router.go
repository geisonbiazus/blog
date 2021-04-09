package web

import "net/http"

func NewRouter(viewPostUseCase ViewPostUseCase) http.Handler {
	mux := http.NewServeMux()

	viewPostHandler := NewViewPostHandler(viewPostUseCase)

	mux.Handle("/", viewPostHandler)

	return mux
}

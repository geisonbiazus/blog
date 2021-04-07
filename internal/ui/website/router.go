package website

import "net/http"

func NewRouter() http.Handler {
	mux := http.NewServeMux()

	viewPostHandler := NewViewPostHandler()

	mux.Handle("/", viewPostHandler)

	return mux
}

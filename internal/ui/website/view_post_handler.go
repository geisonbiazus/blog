package website

import "net/http"

type ViewPostHandler struct{}

func NewViewPostHandler() *ViewPostHandler {
	return &ViewPostHandler{}
}

func (h *ViewPostHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusOK)
}

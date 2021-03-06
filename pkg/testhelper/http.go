package testhelper

import (
	"io"
	"net/http"
)

func ReadResponseBody(res *http.Response) string {
	body, err := io.ReadAll(res.Body)

	if err != nil {
		return ""
	}

	res.Body.Close()
	return string(body)
}

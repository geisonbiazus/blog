package web_test

import (
	"bytes"
	"log"
	"net/http"
	"testing"

	"github.com/geisonbiazus/blog/internal/ui/web"
	"github.com/geisonbiazus/blog/pkg/assert"
)

func TestLogHandler(t *testing.T) {
	t.Run("Given a request, it logs the verb, path and status code", func(t *testing.T) {
		buf := &bytes.Buffer{}
		logger := log.New(buf, "", 0)

		logHandler := web.NewLogHandler(logger, http.NotFoundHandler())

		doGetRequest(logHandler, "/log-path")

		assert.Equal(t, "GET /log-path 404\n", buf.String())
	})
}

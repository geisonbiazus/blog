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

		logHandler = web.NewLogHandler(logger, acceptedHandler())

		doGetRequest(logHandler, "/another-path")

		assert.Equal(t, ""+
			`{"type":"request","method":"GET","path":"/log-path","status":404}`+"\n"+
			`{"type":"request","method":"GET","path":"/another-path","status":202}`+"\n",
			buf.String(),
		)
	})
}

func acceptedHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
	})
}

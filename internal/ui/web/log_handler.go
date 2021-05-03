package web

import (
	"log"
	"net/http"
)

type LogHandler struct {
	Logger  *log.Logger
	Handler http.Handler
}

func NewLogHandler(logger *log.Logger, handler http.Handler) *LogHandler {
	return &LogHandler{Logger: logger, Handler: handler}
}

func (h *LogHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	loggingWritter := newLoggingResponseWriter(w)
	h.Handler.ServeHTTP(loggingWritter, r)
	h.Logger.Printf("%s %s %v", r.Method, r.URL.Path, loggingWritter.statusCode)
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK}
}

func (w *loggingResponseWriter) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

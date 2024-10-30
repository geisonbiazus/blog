package handlers

import (
	"encoding/json"
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
	lw := newLoggingResponseWriter(w)
	h.Handler.ServeHTTP(lw, r)
	h.Logger.Println(h.createRequestLogEntry(lw, r))
}

func (h *LogHandler) createRequestLogEntry(lw *loggingResponseWriter, r *http.Request) string {
	requestLog, err := json.Marshal(logEntry{
		Type:   "request",
		Method: r.Method,
		Path:   r.URL.Path,
		Status: lw.statusCode,
	})

	if err != nil {
		panic(err)
	}

	return string(requestLog)
}

type logEntry struct {
	Type   string `json:"type"`
	Method string `json:"method"`
	Path   string `json:"path"`
	Status int    `json:"status"`
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

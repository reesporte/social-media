package wrapper

import (
	"fmt"
	"log"
	"net/http"
)

// loggingResponseWriter is a struct that logs the statusCode associated with the response.
type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

// newLoggingResponseWriter creates a new loggingResponseWriter.
func newLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK}
}

// WriteHeader sends an HTTP response header with the provided status code and sets the status code in the loggingResponseWriter.
func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

// Log wraps a handler and logs request information.
func Log(h http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionid := "NO_SESSION_ID_FOUND"
		if s, err := r.Cookie("sessionid"); err == nil {
			sessionid = s.Value
		}

		lrw := newLoggingResponseWriter(w)
		h.ServeHTTP(lrw, r)
		status := lrw.statusCode
		log.Println(fmt.Sprint(status) + " " + r.Method + " " + r.URL.Path + " " + r.RemoteAddr + " " + sessionid)
	})
}

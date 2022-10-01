package rest

import (
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Error().Interface("error", err).Msg("panic")
			}
		}()

		start := time.Now()
		wrapped := wrapResponseWriter(w)
		next.ServeHTTP(wrapped, req)

		status := wrapped.status
		if status == 0 {
			status = http.StatusOK
		}

		log.Debug().
			Str("Method", req.Method).
			Str("URI", req.RequestURI).
			Str("HTTP ver", req.Proto).
			Str("IP", req.RemoteAddr).
			Time("Datetime", start).
			Dur("Latency", time.Since(start)).
			Int("Status code", status).
			Str("User-Agent", req.Header.Get("User-Agent")).
			Send()
	}

	return http.HandlerFunc(fn)
}

type responseWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func wrapResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w}
}

func (rw *responseWriter) Status() int {
	return rw.status
}

func (rw *responseWriter) WriteHeader(code int) {
	if rw.wroteHeader {
		return
	}

	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
	rw.wroteHeader = true
}

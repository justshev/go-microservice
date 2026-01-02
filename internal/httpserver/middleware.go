package httpserver

import (
	"net/http"
	"time"

	"github.com/rs/zerolog"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)
type statusRecorder struct {
	http.ResponseWriter
	status int
	bytes int
}


func (r *statusRecorder) WriteHeader(code int ){
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}
func (r *statusRecorder) Write(b []byte) (int,error){
if r.status == 0 {
	r.status = http.StatusOK
}
n, err := r.ResponseWriter.Write(b)
r.bytes += n 
return n,err
}


func SetupMiddleware(r chi.Router,baseLogger zerolog.Logger) {
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(RequestLogger(baseLogger))
}

func RequestLogger(base zerolog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			rec := &statusRecorder{ResponseWriter: w}

			next.ServeHTTP(rec, r)

			reqID := r.Context().Value(middleware.RequestIDKey)

			log := base.With().
				Interface("request_id", reqID).
				Str("method", r.Method).
				Str("path", r.URL.Path).
				Int("status", rec.status).
				Int("bytes", rec.bytes).
				Int64("duration_ms", time.Since(start).Milliseconds()).
				Logger()
			log.Info().Msg("http_request")
		})
	}
}

package httpserver

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
)

func NewRouter(baseLogger zerolog.Logger) http.Handler { 
	r := chi.NewRouter()
	SetupMiddleware(r, baseLogger)

	//routes 
	r.Get("/health",healthHandler)
	r.Route("/v1",func (r chi.Router){
		r.Get("/tasks",listTaskHandler)
	})

	return r

}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"ok"}`))
}

func listTaskHandler(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"tasks":[
		{"id":1,"name":"Task 1","completed":false},
		{"id":2,"name":"Task 2","completed":true}
	]}`))
}
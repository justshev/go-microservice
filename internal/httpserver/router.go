package httpserver

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func NewRouter() http.Handler { 
	r := chi.NewRouter()
	SetupMiddleware(r)

	//routes 
	r.Get("/health",healthHandler)
	r.Route("/v1",func (r chi.Router){
		r.Get("/tasks",liskTaskHandler)
	})

	return r

}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"ok"}`))
}

func liskTaskHandler(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"tasks":[
		{"id":1,"name":"Task 1","completed":false},
		{"id":2,"name":"Task 2","completed":true}
	]}`))
}
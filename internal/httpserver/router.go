package httpserver

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/justshev/go-micro-template/internal/db"
	"github.com/justshev/go-micro-template/internal/task"
	"github.com/rs/zerolog"
)

func NewRouter(baseLogger zerolog.Logger) http.Handler { 
	r := chi.NewRouter()
	SetupMiddleware(r, baseLogger)
	pg,err := db.NewPostgres(
		"postgres://taskuser:taskpass@localhost:5432/taskdb?sslmode=disable",
	)
	if err != nil {
		baseLogger.Fatal().Err(err).Msg("failed to connect to postgres")
	}
	taskRepo := task.NewPostgresRepo(pg)
	taskService := task.NewService(taskRepo)

	//routes 
	r.Get("/health",healthHandler)
	r.Route("/v1",func (r chi.Router){
		r.Get("/tasks",listTaskHandler(taskService))
		r.Post("/tasks",createTaskHandler(taskService))
	})

	return r

}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"ok"}`))
}

func listTaskHandler(svc *task.Service) http.HandlerFunc {
return func (W http.ResponseWriter, r *http.Request){
	tasks, err := svc.List(r.Context())
	if err != nil {
		http.Error(W, "failed to list tasks", http.StatusInternalServerError)
		return
	}
	writeJSON(W,http.StatusOK,map[string]any{
		"tasks": tasks,
	})

}
}

type createTaskRequest struct {
	Name string `json:"name"`
}

func createTaskHandler (svc *task.Service) http.HandlerFunc {
return func (W http.ResponseWriter, r *http.Request){
	var req createTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(W, "invalid request payload", http.StatusBadRequest)
		return
	}
	t, err := svc.Create(r.Context(),req.Name) 
	if err != nil { 
		if err == task.ErrInvalidName {
			writeJSON(W,http.StatusBadRequest,map[string]any{
				"error": err.Error(),
			})
			return
		}
		http.Error(W, "failed to create task", http.StatusInternalServerError)
		return

	}
	writeJSON(W,http.StatusCreated,map[string]any{
		"task": t,
	})
}
}


func writeJSON(W http.ResponseWriter, status int, v any){
	W.Header().Set("Content-Type", "application/json")
	W.WriteHeader(status)
	//serialize v to json and write to W
	_ = json.NewEncoder(W).Encode(v)
}
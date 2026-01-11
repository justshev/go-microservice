package httpserver

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/justshev/go-micro-template/internal/broker"
	"github.com/justshev/go-micro-template/internal/cache"
	"github.com/justshev/go-micro-template/internal/db"
	"github.com/justshev/go-micro-template/internal/task"
	"github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

func NewRouter(baseLogger zerolog.Logger) http.Handler { 
	rdb,err := cache.NewRedis("localhost:6379")
	if err != nil {
		panic(err)
	}
	r := chi.NewRouter()

	conn, err := broker.Connect("amqp://guest:guest@localhost:5672/")
	if err != nil {
		baseLogger.Fatal().Err(err).Msg("failed to connect to rabbitmq")
	}
	ch, err := conn.Channel()
	if err != nil {
		baseLogger.Fatal().Err(err).Msg("failed to open rabbitmq channel")
	}

	_,err = ch.QueueDeclare(
		"task.created.queue",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		baseLogger.Fatal().Err(err).Msg("failed to declare rabbitmq queue")
	}

	SetupMiddleware(r, baseLogger)
	r.Use(RateLimit(rdb,5,time.Minute))
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
		r.Get("/tasks",listTaskHandler(taskService,rdb,ch))
		r.Post("/tasks",createTaskHandler(taskService,rdb,ch, &baseLogger))
	})

	return r

}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"ok"}`))
}

func listTaskHandler(svc *task.Service, rdb *redis.Client, ch *amqp091.Channel) http.HandlerFunc {
return func (w http.ResponseWriter, r *http.Request){
	
	cacheKey := "tasks:all"

	val, err := rdb.Get(r.Context(),cacheKey).Result()
	if err == nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(val))
			return
	}
	// miss cache
	tasks, err := svc.List(r.Context())
	if err != nil { 
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "internal_error"})
			return
	}
	resp, _ := json.Marshal(map[string]any{
		"tasks": tasks,
	})
	// set cache
	_ = rdb.Set(r.Context(),cacheKey,string(resp),30*time.Second).Err()


	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}
}

type createTaskRequest struct {
	Name string `json:"name"`
}

func createTaskHandler (svc *task.Service, rdb *redis.Client, ch *amqp091.Channel, baseLogger *zerolog.Logger) http.HandlerFunc {
return func (W http.ResponseWriter, r *http.Request){
	var req createTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(W, "invalid request payload", http.StatusBadRequest)
		return
	}
	t, err := svc.Create(r.Context(),req.Name) 
	event := map[string]any {
		"event": "task.created",
		"task_id": t.ID,
		"task_name": t.Name,
		"created_at": t.CreatedAt,
	}
	body, _ := json.Marshal(event)
	pubErr := ch.Publish(
		"",
		"task.created.queue",
		false,
		false,
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        body,
			Timestamp: time.Now(),
		},
	)
	if pubErr != nil {
		baseLogger.Error().Err(pubErr).Msg("failed to publish task.created event")
	}
	
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

	_ = rdb.Del(r.Context(),"tasks:all").Err()
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
package task

import (
	"context"
	"sync"
	"time"
)

type Repository interface { 
	List(ctx context.Context)([]Task,error)
	Create(ctx context.Context,name string)(Task,error)
}

type MemoryRepo struct {
	mu sync.Mutex
	nextID int64
	tasks []Task

}

func NewMemoryRepo() *MemoryRepo { 
	return &MemoryRepo{
		nextID: 3,
		tasks: []Task{
			{ID: 1, Name: "Task 1", Completed: false, CreatedAt: time.Now()},
			{ID: 2, Name: "Task 2", Completed: true, CreatedAt: time.Now()},
		},

	}
}

func (r *MemoryRepo) List(ctx context.Context)([]Task,error){ 
	r.mu.Lock()
	defer r.mu.Unlock()

	out := make([]Task,len(r.tasks))
	copy(out,r.tasks)
	return out,nil
}


func (r *MemoryRepo) Create (ctx context.Context,name string)(Task,error){
	r.mu.Lock()
	defer r.mu.Unlock()

	t := Task {
		ID: r.nextID,
		Name: name,
		Completed: false,
		CreatedAt: time.Now(),
	}
	r.nextID++
	r.tasks = append(r.tasks,t)
	return t,nil
	
}
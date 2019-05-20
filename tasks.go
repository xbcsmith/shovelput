package main

import (
	"net/http"

	"github.com/go-chi/chi"
)

type tasksResource struct{}

func (rs tasksResource) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", rs.List)    // GET /tasks - read a list of tasks
	r.Post("/", rs.Create) // POST /tasks - create a new task and persist it
	r.Put("/", rs.Delete)

	r.Route("/{id}", func(r chi.Router) {
		r.Get("/", rs.Get)       // GET /tasks/{id} - read a single task by :id
		r.Put("/", rs.Update)    // PUT /tasks/{id} - update a single task by :id
		r.Delete("/", rs.Delete) // DELETE /tasks/{id} - delete a single task by :id
		r.Get("/sync", rs.Sync)
	})

	return r
}

func (rs tasksResource) List(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("tasks list of stuff.."))
}

func (rs tasksResource) Create(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("tasks create"))
}

func (rs tasksResource) Get(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("task get"))
}

func (rs tasksResource) Update(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("task update"))
}

func (rs tasksResource) Delete(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("task delete"))
}

func (rs tasksResource) Sync(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("task sync"))
}

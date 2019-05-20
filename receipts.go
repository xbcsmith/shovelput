package main

import (
	"net/http"

	"github.com/go-chi/chi"
)

type receiptsResource struct{}

// Routes creates a REST router for the tasks resource
func (rs receiptsResource) Routes() chi.Router {
	r := chi.NewRouter()
	// r.Use() // some middleware..

	r.Get("/", rs.List)    // GET /tasks - read a list of tasks
	r.Post("/", rs.Create) // POST /tasks - create a new task and persist it
	r.Put("/", rs.Delete)

	r.Route("/{id}", func(r chi.Router) {
		// r.Use(rs.TodoCtx) // lets have a tasks map, and lets actually load/manipulate
		r.Get("/", rs.Get)       // GET /tasks/{id} - read a single task by :id
		r.Put("/", rs.Update)    // PUT /tasks/{id} - update a single task by :id
		r.Delete("/", rs.Delete) // DELETE /tasks/{id} - delete a single task by :id
	})

	return r
}

func (rs receiptsResource) List(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("aaa list of stuff.."))
}

func (rs receiptsResource) Create(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("aaa create"))
}

func (rs receiptsResource) Get(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("aaa get"))
}

func (rs receiptsResource) Update(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("aaa update"))
}

func (rs receiptsResource) Delete(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("aaa delete"))
}

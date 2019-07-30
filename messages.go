package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
)

type Message struct {
	ID       string
	Name     string
	Version  string
	Action   string
	Artifact string
	Metadata string
}

var messages []Message

type messageResource struct{}

func (rs messageResource) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", rs.List)    // GET /message - read a list of message
	r.Post("/", rs.Create) // POST /message - create a new message and persist it
	r.Put("/", rs.Delete)

	r.Route("/{id}", func(r chi.Router) {
		r.Get("/", rs.Get)       // GET /message/{id} - read a single message by :id
		r.Put("/", rs.Update)    // PUT /message/{id} - update a single message by :id
		r.Delete("/", rs.Delete) // DELETE /message/{id} - delete a single message by :id
		r.Get("/sync", rs.Sync)
	})

	return r
}

func (rs messageResource) List(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("message list of stuff.."))
	json.NewEncoder(w).Encode(messages)
}

func (rs messageResource) Create(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("message create"))
	params := r.
	var message Message
	_ = json.NewDecoder(r.Body).Decode(&message)
	message.ID = params["id"]
	messages = append(messages, message)
	json.NewEncoder(w).Encode(messages)
}

func (rs messageResource) Get(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("message get"))
	params := mux.Vars(r)
	for _, item := range messages {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Message{})
}

func (rs messageResource) Update(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("message update"))
}

func (rs messageResource) Delete(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("message delete"))
}

func (rs messageResource) Sync(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("message sync"))
}

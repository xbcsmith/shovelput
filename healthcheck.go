package main

import (
	"net/http"

	"github.com/go-chi/chi"
)

type healthCheckResource struct{}

func (rs healthCheckResource) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/healthz", rs.Get)

	return r
}

func (rs healthCheckResource) Get(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`{"alive":true}`))
}

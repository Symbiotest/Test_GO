package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)


func NewHTTPHandler() http.Handler {
	db := NewDBClient()
	ac := &FakeAuthorizationsClient{}
	controller := NewCustomerController(ac, db)
	routeHandler := NewRouteHandler(controller)

	r := chi.NewRouter()
	r.Use(chiTraceMiddleware)

	HandlerFromMux(routeHandler, r, oapiTraceMiddleware)

	return httpServerEntry(r)
}

func httpServerEntry(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r = ensureTrace(r)
		addTrace(r.Context(), "http-server")
		next.ServeHTTP(w, r)
	})
}

func chiTraceMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		addTrace(r.Context(), "chi")
		next.ServeHTTP(w, r)
	})
}

func oapiTraceMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		addTrace(r.Context(), "oapi-middleware")
		next.ServeHTTP(w, r)
	})
}

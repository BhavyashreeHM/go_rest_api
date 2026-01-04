package router

import (
	"net/http"
	"rest_api_go/internal/api/handlers"
)

func StudentRouter() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /student", handlers.StudentsHandler)
	mux.HandleFunc("GET  /students", handlers.GetstudentsHandler)
	mux.HandleFunc("POST /students", handlers.AddstudentsHandler)
	mux.HandleFunc("PATCH /students", handlers.PatchstudentsHandler)
	mux.HandleFunc("DELETE /students", handlers.DeletestudentsHandler)

	mux.HandleFunc("GET  /students/{id}", handlers.GetstudentsByIdHandler)
	mux.HandleFunc("PUT /students/{id}", handlers.UpdatestudentsByIdHandler)
	mux.HandleFunc("PATCH /students/{id}", handlers.PatchstudentsByIdHandler)
	mux.HandleFunc("DELETE /students/{id}", handlers.DeletestudentsByIdHandler)
	return mux
}

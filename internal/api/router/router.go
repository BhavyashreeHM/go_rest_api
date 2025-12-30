package router

import (
	"net/http"
	"rest_api_go/internal/api/handlers"
)

func Router() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handlers.RootHandler)
	mux.HandleFunc("GET  /teachers", handlers.GetTeachersHandler)
	mux.HandleFunc("POST /teachers", handlers.AddteacherHandler)
	mux.HandleFunc("PATCH /teachers", handlers.PatchTeacherHandler)
	mux.HandleFunc("DELETE /teachers", handlers.DeleteTeacherHandler)

	mux.HandleFunc("GET  /teachers/{id}", handlers.GetTeacherByIdHandler)
	mux.HandleFunc("PUT /teachers/{id}", handlers.UpdateTeacherByIdHandler)
	mux.HandleFunc("PATCH /teachers/{id}", handlers.PatchTeacherByIdHandler)
	mux.HandleFunc("DELETE /teachers/{id}", handlers.DeleteTeacherByIdHandler)

	mux.HandleFunc("/students", handlers.StudentsHandler)
	mux.HandleFunc("/execs", handlers.ExecsHandler)
	return mux
}

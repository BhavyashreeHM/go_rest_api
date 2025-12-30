package handlers

import (
	"fmt"
	"net/http"
)

func StudentsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the Student's Page!\n")
	switch r.Method {
	case http.MethodGet:
		// Handle GET request to retrieve students
	case http.MethodPost:
		// Handle POST request to create a new student
	case http.MethodPut:
		// Handle PUT request to update a student
	case http.MethodDelete:
		// Handle DELETE request to remove a student
	case http.MethodPatch:
		// Handle PATCH request to partially update a student
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

}

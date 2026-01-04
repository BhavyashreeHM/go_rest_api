package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"rest_api_go/internal/models"
	"rest_api_go/internal/repository/sqlconnect"
	"strconv"
)

func StudentsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the Student's Page!\n")
}

func GetstudentsHandler(w http.ResponseWriter, r *http.Request) {
	var students []models.Student
	students, err := sqlconnect.GetstudentsDbHandler(students, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := struct {
		Staus string           `json:"status"`
		Count int              `json:"count"`
		Data  []models.Student `json:"data"`
	}{
		Staus: "success",
		Count: len(students),
		Data:  students,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}

func GetstudentsByIdHandler(w http.ResponseWriter, r *http.Request) {

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Invalid request Body", http.StatusBadRequest)
		return
	}

	student, err := sqlconnect.GetstudentsByIdDbHandler(id)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(student)
}

func AddstudentsHandler(w http.ResponseWriter, r *http.Request) {

	var newstudents []models.Student
	err := json.NewDecoder(r.Body).Decode(&newstudents)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println("New student Data:", newstudents)

	addedstudents, err := sqlconnect.AddstudentsDbHandler(newstudents)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	response := struct {
		Status string           `json:"status"`
		Count  int              `json:"count"`
		Data   []models.Student `json:"data"`
	}{
		Status: "sucess",
		Count:  len(addedstudents),
		Data:   addedstudents,
	}
	json.NewEncoder(w).Encode(response)

}

func UpdatestudentsByIdHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Invalid student ID", http.StatusInternalServerError)
		return
	}

	var Updates models.Student
	err = json.NewDecoder(r.Body).Decode(&Updates)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Invalid Request Payload", http.StatusInternalServerError)
		return
	}

	UpdatedstudentFromDb, err := sqlconnect.UpdatestudentsDbHandler(id, Updates)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "applicatuion/jsno")
	json.NewEncoder(w).Encode(UpdatedstudentFromDb)
	fmt.Printf("%d Updated succesfully", id)
}

func PatchstudentsHandler(w http.ResponseWriter, r *http.Request) {
	var Updates []map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&Updates)
	if err != nil {
		// fmt.Println(err)
		http.Error(w, "Invalid Request Payload", http.StatusBadRequest)
		return
	}
	err = sqlconnect.PatchstudentsDbHandler(Updates)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
	w.Header().Set("Content-Type", "applicatuion/jsno")
	json.NewEncoder(w).Encode(Updates)
	fmt.Printf(" Patch Operation Done Successfully")

}

func PatchstudentsByIdHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		// fmt.Println(err)
		http.Error(w, "Invalid student ID", http.StatusInternalServerError)
		return
	}

	var Updates map[string]interface{}
	err = json.NewDecoder(r.Body).Decode(&Updates)
	if err != nil {
		// fmt.Println(err)
		http.Error(w, "Invalid Request Payload", http.StatusInternalServerError)
		return
	}

	Updatedstudent, err := sqlconnect.PatchstudentsByIdDbHandler(id, Updates)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		// fmt.Println(err)
		return
	}
	w.Header().Set("Content-Type", "applicatuion/jsno")
	json.NewEncoder(w).Encode(Updatedstudent)
	fmt.Printf("%d Patch Operation Done", id)
}

// DeletestudentHandler
func DeletestudentsHandler(w http.ResponseWriter, r *http.Request) {

	var ids []int
	err := json.NewDecoder(r.Body).Decode(&ids)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid Request Payload", http.StatusBadRequest)
		return
	}
	deletedIds, err := sqlconnect.DeletestudentsDbHandler(ids)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content_Type", "application/json")
	response := struct {
		Status string `json:"status"`
		ID     []int  `json:"id"`
	}{
		Status: "student succesfully deleted",
		ID:     deletedIds,
	}
	json.NewEncoder(w).Encode(response)

}

func DeletestudentsByIdHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Invalid student ID", http.StatusInternalServerError)
		return
	}

	err = sqlconnect.DeletestudentsByIdDbHandler(id)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNotFound)

	w.Header().Set("Content_Type", "application/json")
	response := struct {
		Status string `json:"status"`
		ID     int    `json:"id"`
	}{
		Status: "student succesfully deleted",
		ID:     id,
	}
	json.NewEncoder(w).Encode(response)

}

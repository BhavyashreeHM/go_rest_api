package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"rest_api_go/internal/models"
	"rest_api_go/internal/repository/sqlconnect"
	_ "rest_api_go/internal/repository/sqlconnect"
	"strconv"
)

func GetTeachersHandler(w http.ResponseWriter, r *http.Request) {
	var teachers []models.Teacher
	teachers, err := sqlconnect.GetTeacherDbHandler(teachers, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := struct {
		Staus string           `json:"status"`
		Count int              `json:"count"`
		Data  []models.Teacher `json:"data"`
	}{
		Staus: "success",
		Count: len(teachers),
		Data:  teachers,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}

func GetTeacherByIdHandler(w http.ResponseWriter, r *http.Request) {

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Invalid request Body", http.StatusBadRequest)
		return
	}

	teacher, err := sqlconnect.GetTeacherByIdDbHandler(id)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(teacher)
}

func AddteacherHandler(w http.ResponseWriter, r *http.Request) {

	var newTeachers []models.Teacher
	err := json.NewDecoder(r.Body).Decode(&newTeachers)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println("New Teacher Data:", newTeachers)

	addedTeachers, err := sqlconnect.AddTeachersDbHandler(newTeachers)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	response := struct {
		Status string           `json:"status"`
		Count  int              `json:"count"`
		Data   []models.Teacher `json:"data"`
	}{
		Status: "sucess",
		Count:  len(addedTeachers),
		Data:   addedTeachers,
	}
	json.NewEncoder(w).Encode(response)

}

func UpdateTeacherByIdHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Invalid Teacher ID", http.StatusInternalServerError)
		return
	}

	var Updates models.Teacher
	err = json.NewDecoder(r.Body).Decode(&Updates)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Invalid Request Payload", http.StatusInternalServerError)
		return
	}

	UpdatedTeacherFromDb, err := sqlconnect.UpdateTeacherDbHandler(id, Updates)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "applicatuion/jsno")
	json.NewEncoder(w).Encode(UpdatedTeacherFromDb)
	fmt.Printf("%d Updated succesfully", id)
}

func PatchTeacherHandler(w http.ResponseWriter, r *http.Request) {
	var Updates []map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&Updates)
	if err != nil {
		// fmt.Println(err)
		http.Error(w, "Invalid Request Payload", http.StatusBadRequest)
		return
	}
	err = sqlconnect.PatchTeacherDbHandler(Updates)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
	w.Header().Set("Content-Type", "applicatuion/jsno")
	json.NewEncoder(w).Encode(Updates)
	fmt.Printf(" Patch Operation Done Successfully")

}

func PatchTeacherByIdHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		// fmt.Println(err)
		http.Error(w, "Invalid Teacher ID", http.StatusInternalServerError)
		return
	}

	var Updates map[string]interface{}
	err = json.NewDecoder(r.Body).Decode(&Updates)
	if err != nil {
		// fmt.Println(err)
		http.Error(w, "Invalid Request Payload", http.StatusInternalServerError)
		return
	}

	UpdatedTeacher, err := sqlconnect.PatchTeacherByIdDbHandler(id, Updates)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		// fmt.Println(err)
		return
	}
	w.Header().Set("Content-Type", "applicatuion/jsno")
	json.NewEncoder(w).Encode(UpdatedTeacher)
	fmt.Printf("%d Patch Operation Done", id)
}

// DeleteTeacherHandler
func DeleteTeacherHandler(w http.ResponseWriter, r *http.Request) {

	var ids []int
	err := json.NewDecoder(r.Body).Decode(&ids)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid Request Payload", http.StatusBadRequest)
		return
	}
	deletedIds, err := sqlconnect.DeleteTeacherDbHandler(ids)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content_Type", "application/json")
	response := struct {
		Status string `json:"status"`
		ID     []int  `json:"id"`
	}{
		Status: "Teacher succesfully deleted",
		ID:     deletedIds,
	}
	json.NewEncoder(w).Encode(response)

}

func DeleteTeacherByIdHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Invalid Teacher ID", http.StatusInternalServerError)
		return
	}

	err = sqlconnect.DeleteTeacherByIdDbHandler(id)
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
		Status: "Teacher succesfully deleted",
		ID:     id,
	}
	json.NewEncoder(w).Encode(response)

}

func GetStudentsByTeacherId(w http.ResponseWriter, r *http.Request) {
	teacherId := r.PathValue("id")

	var students []models.Student

	students, err := sqlconnect.GetStudentsByTeacherIdFromDb(teacherId, students)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := struct {
		Status string           `json:"status"`
		Count  int              `json:"count"`
		Data   []models.Student `json:"data"`
	}{
		Status: "success",
		Count:  len(students),
		Data:   students,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func GetStudentCountByTeacherId(w http.ResponseWriter, r *http.Request) {
	// admin, manager, exec
	// _, err := utils.AuthorizeUser(r.Context().Value(utils.ContextKey("role")).(string), "admin", "manager", "exec")
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	teacherId := r.PathValue("id")

	studentCount, err := sqlconnect.GetStudentCountByTeacherIdFromDb(teacherId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := struct {
		Status string `json:"status"`
		Count  int    `json:"count"`
	}{
		Status: "success",
		Count:  studentCount,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

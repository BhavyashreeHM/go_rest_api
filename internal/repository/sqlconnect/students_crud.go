package sqlconnect

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"rest_api_go/internal/models"
	"rest_api_go/pkg/utils"
	"strconv"
	"strings"
)

func GetstudentsDbHandler(students []models.Student, r *http.Request) ([]models.Student, error) {
	db, err := ConnectDB()
	if err != nil {
		return nil, utils.Errorhandler(err, "Error Connecting to database")
	} else {
		fmt.Println("Database connected from GetstudentHandler")
	}
	defer db.Close()

	query := "select id,first_name,last_name,email,class from students where 1=1"
	var args []interface{}
	query, args = utils.AddFilters(r, query, args)
	query = utils.AddSorting(r, query)

	rows, err := db.Query(query, args...)
	if err != nil {
		fmt.Println(err)
		return nil, utils.Errorhandler(err, "error retrieving database")
	}
	defer rows.Close()

	for rows.Next() {
		var student models.Student
		err := rows.Scan(&student.Id, &student.FirstName, &student.LastName, &student.Email, &student.Class)
		if err != nil {
			fmt.Println(err)
			// http.Error(w, "", http.StatusInternalServerError)
			return nil, utils.Errorhandler(err, "Error Scanning database error")
		}
		students = append(students, student)
	}
	return students, nil

}

func GetstudentsByIdDbHandler(id int) (models.Student, error) {
	db, err := ConnectDB()
	if err != nil {
		// http.Error(w, err.Error(), http.StatusInternalServerError)
		return models.Student{}, utils.Errorhandler(err, "Error Connecting to database")
	} else {
		fmt.Println("Database connected from GetstudentHandler")
	}
	defer db.Close()

	var student models.Student
	err = db.QueryRow("select id,first_name,last_name,email,class from students where id=?", id).Scan(&student.Id, &student.FirstName,
		&student.LastName, &student.Email, &student.Class)
	if err == sql.ErrNoRows {
		// http.Error(w, "student not found", http.StatusNotFound)
		return models.Student{}, utils.Errorhandler(err, "student not found")
	} else if err != nil {
		// http.Error(w, "Database query error", http.StatusInternalServerError)
		return models.Student{}, utils.Errorhandler(err, "Database query error")
	}
	return student, nil
}

func AddstudentsDbHandler(newstudents []models.Student) ([]models.Student, error) {
	db, err := ConnectDB()
	if err != nil {
		// http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil, utils.Errorhandler(err, "Database Connection error")
	} else {
		fmt.Println("Database connection established in AddstudentHandler")
	}
	defer db.Close()

	// stmt, err := db.Prepare("insert into students(first_name, last_name, email, class, subject) values (?, ?, ?, ?, ?)")
	stmt, err := db.Prepare(utils.GenerateInsertQuery("students", models.Student{}))
	if err != nil {
		// http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil, utils.Errorhandler(err, "error preparing stament")
	}
	defer stmt.Close()

	addedstudents := make([]models.Student, len(newstudents))
	for i, student := range newstudents {
		// res, err := stmt.Exec(student.FirstName, student.LastName, student.Email, student.Class, student.Subject)
		values := utils.GetStructValues(student)
		res, err := stmt.Exec(values...)
		if err != nil {
			// http.Error(w, err.Error(), http.StatusInternalServerError)
			fmt.Println("---------Error", err.Error())
			if strings.Contains(err.Error(), "a foreign key constraint fails (`school`.`students`, CONSTRAINT `students_ibfk_1` FOREIGN KEY (`class`) REFERENCES `teachers` (`class`))") {
				return nil, utils.Errorhandler(err, "class/class teacher does not exist")
			}
			return nil, utils.Errorhandler(err, "error in adding new student")
		}
		lastID, err := res.LastInsertId()
		if err != nil {
			// http.Error(w, err.Error(),http.StatusInternalServerError)
			// http.Error(w, "Error getting last insert ID", http.StatusInternalServerError)
			return nil, utils.Errorhandler(err, "Error getting last insert ID")
		}
		student.Id = int(lastID)
		addedstudents[i] = student

	}
	return addedstudents, nil
}

func UpdatestudentsDbHandler(id int, Updates models.Student) (models.Student, error) {
	db, err := ConnectDB()
	if err != nil {
		fmt.Println(err)
		// http.Error(w, "Unable to connect database", http.StatusInternalServerError)
		return models.Student{}, utils.Errorhandler(err, "Database Connection error")
	} else {
		fmt.Println("Database connected from PUTstudentHandler")
	}
	defer db.Close()

	var existingstudent models.Student
	err = db.QueryRow("SELECT id,first_name,last_name,email,class FROM students WHERE id =?", id).Scan(&existingstudent.Id,
		&existingstudent.FirstName, &existingstudent.LastName, &existingstudent.Email, &existingstudent.Class)

	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println(err)
			// http.Error(w, "student not found", http.StatusNotFound)
			return models.Student{}, utils.Errorhandler(err, "student not found")
		}
		fmt.Println(err)
		// http.Error(w, "Unable to retrieve data", http.StatusInternalServerError)
		return models.Student{}, utils.Errorhandler(err, "Unable to retrieve data")
	}
	Updates.Id = existingstudent.Id
	_, err = db.Exec("UPDATE students SET first_name =?, last_name=?,email=?, class=? WHERE id=?", Updates.FirstName,
		Updates.LastName, Updates.Email, Updates.Class, Updates.Id)
	if err != nil {
		fmt.Println(err)
		// http.Error(w, "Error while updating", http.StatusInternalServerError)
		return models.Student{}, utils.Errorhandler(err, "Error while updating")
	}
	return Updates, nil
}

func PatchstudentsDbHandler(Updates []map[string]interface{}) error {
	db, err := ConnectDB()
	if err != nil {
		// http.Error(w, "Unable to connect to database", http.StatusInternalServerError)
		return utils.Errorhandler(err, "Unable to connect to database")
	} else {
		fmt.Println("Database connected from PatchstudentHandlerr")
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		// print(err)
		// http.Error(w, "Error starting transaction", http.StatusInternalServerError)
		return utils.Errorhandler(err, "Error starting transaction")
	}

	for _, update := range Updates {
		idStr, ok := update["id"].(string)
		if !ok {
			tx.Rollback()
			// http.Error(w, "Invalid student ID in update", http.StatusBadRequest)
			return utils.Errorhandler(err, "Invalid student ID in update")
		}

		id, err := strconv.Atoi(idStr)
		if err != nil {
			tx.Rollback()
			// http.Error(w, "Invalid student ID", http.StatusInternalServerError)
			return utils.Errorhandler(err, "Invalid student ID")
		}
		var studentFromDb models.Student
		err = db.QueryRow("SELECT id,first_name,last_name,email,class FROM students WHERE id=?", id).Scan(&studentFromDb.Id,
			&studentFromDb.FirstName, &studentFromDb.LastName, &studentFromDb.Email, &studentFromDb.Class)
		if err != nil {
			tx.Rollback()
			if err == sql.ErrNoRows {
				print(err)
				// http.Error(w, "student Not Found", http.StatusNotFound)
				return utils.Errorhandler(err, "student Not Found")
			}
			// http.Error(w, "Error Retrieving student", http.StatusInternalServerError)
			return utils.Errorhandler(err, "Error Retrieving student")
		}

		studentVal := reflect.ValueOf(&studentFromDb).Elem()
		techerType := studentVal.Type()

		for k, v := range update {
			if k == "id" {
				continue
			}
			for i := 0; i < studentVal.NumField(); i++ {
				field := techerType.Field(i)
				if field.Tag.Get("json") == k+",omitempty" {
					fieldVal := studentVal.Field(i)
					if studentVal.CanSet() {
						val := reflect.ValueOf(v)
						if val.Type().ConvertibleTo(fieldVal.Type()) {
							fieldVal.Set(val.Convert(fieldVal.Type()))

						} else {
							tx.Rollback()
							log.Printf("Cannot convert %v to  %v", val.Type(), fieldVal.Type())
							return utils.Errorhandler(err, "Error starting transaction")
						}

					}
				}
			}
		}
		_, err = db.Exec("UPDATE students SET first_name =?, last_name=?,email=?, class=? WHERE id=?",
			studentFromDb.FirstName, studentFromDb.LastName, studentFromDb.Email, studentFromDb.Class, studentFromDb.Id)
		if err != nil {
			tx.Rollback()
			// http.Error(w, "Error while updating", http.StatusInternalServerError)
			return utils.Errorhandler(err, "Error while updating")
		}

	}
	err = tx.Commit()
	if err != nil {
		// http.Error(w, "Error comiting transaction", http.StatusInternalServerError)
		return utils.Errorhandler(err, "Error comiting transaction")
	}
	return nil
}

func PatchstudentsByIdDbHandler(id int, Updates map[string]interface{}) (models.Student, error) {
	db, err := ConnectDB()
	if err != nil {
		// fmt.Println(err)
		// http.Error(w, "Unable to connect database", http.StatusInternalServerError)
		return models.Student{}, utils.Errorhandler(err, "Unable to connect to database")
	} else {
		fmt.Println("Database connected from PatchstudentByIdHandlerr")
	}
	defer db.Close()

	var existingstudent models.Student
	err = db.QueryRow("SELECT id,first_name,last_name,email,class FROM students WHERE id =?", id).Scan(&existingstudent.Id,
		&existingstudent.FirstName, &existingstudent.LastName, &existingstudent.Email, &existingstudent.Class)

	if err != nil {
		if err == sql.ErrNoRows {
			// fmt.Println(err)
			// http.Error(w, "student not found", http.StatusNotFound)
			return models.Student{}, utils.Errorhandler(err, "student not found")
		}
		fmt.Println(err)
		// http.Error(w, "Unable to retrieve data", http.StatusInternalServerError)
		return models.Student{}, utils.Errorhandler(err, "Unable to retrieve data")
	}
	studentVal := reflect.ValueOf(&existingstudent).Elem()
	techerType := studentVal.Type()

	for k, v := range Updates {
		for i := 0; i < studentVal.NumField(); i++ {
			field := techerType.Field(i)
			if field.Tag.Get("json") == k+",omitempty" {
				if studentVal.Field(i).CanSet() {
					fieldVal := studentVal.Field(i)
					fmt.Println("fieldVal", fieldVal)
					fmt.Println("studentVal.Field(i).Type()", studentVal.Field(i).Type())
					fmt.Println("reflect.ValueOf(v)", reflect.ValueOf(v))
					fieldVal.Set(reflect.ValueOf(v).Convert(studentVal.Field(i).Type()))
				}
			}
		}
	}

	_, err = db.Exec("UPDATE students SET first_name =?, last_name=?,email=?, class=? WHERE id=?", existingstudent.FirstName,
		existingstudent.LastName, existingstudent.Email, existingstudent.Class, existingstudent.Id)
	if err != nil {
		// fmt.Println(err)
		// http.Error(w, "Error while updating", http.StatusInternalServerError)
		return models.Student{}, utils.Errorhandler(err, "Error while updating")
	}
	return existingstudent, nil
}

func DeletestudentsByIdDbHandler(id int) error {
	db, err := ConnectDB()
	if err != nil {
		// fmt.Println(err)
		// http.Error(w, "Unable to connect database", http.StatusInternalServerError)
		return utils.Errorhandler(err, "Unable to connect database")
	} else {
		fmt.Println("Database connected from DELETEByIdstudentHandler")
	}
	defer db.Close()

	result, err := db.Exec("DELETE from students WHERE id =?", id)
	if err != nil {
		// http.Error(w, "Error deleting student", http.StatusInternalServerError)
		return utils.Errorhandler(err, "Error deleting student")
	}

	fmt.Println(result.RowsAffected())
	rowsAffect, err := result.RowsAffected()
	if err != nil {
		// http.Error(w, "Error retrived deleted result", http.StatusInternalServerError)
		return utils.Errorhandler(err, "Error retrived deleted result")
	}
	if rowsAffect == 0 {
		// http.Error(w, "student not found", http.StatusNotFound)
		return utils.Errorhandler(err, "student not found")
	}
	return nil
}

func DeletestudentsDbHandler(ids []int) ([]int, error) {
	db, err := ConnectDB()
	if err != nil {
		// http.Error(w, "Unable to connect database", http.StatusInternalServerError)
		return nil, utils.Errorhandler(err, "Unable to connect database")
	} else {
		fmt.Println("Database Connected from DELETEHandler")
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		print(err)
		// http.Error(w, "Error starting transaction", http.StatusInternalServerError)
		return nil, utils.Errorhandler(err, "Error starting transaction")
	}
	stmt, err := tx.Prepare("DELETE FROM students WHERE id =?")
	if err != nil {
		log.Println(err)
		tx.Rollback()
		// http.Error(w, "Error retrived deleted result", http.StatusInternalServerError)
		return nil, utils.Errorhandler(err, "Error retrived deleted result")
	}
	defer stmt.Close()

	deletedIds := []int{}

	for _, id := range ids {
		result, err := stmt.Exec(id)
		if err != nil {
			// log.Println(err)
			tx.Rollback()
			// http.Error(w, "Error deleting student", http.StatusInternalServerError)
			return nil, utils.Errorhandler(err, "Error deleting student")
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			// log.Println(err)
			tx.Rollback()
			// http.Error(w, "Error Retrieving deleted student", http.StatusInternalServerError)
			return nil, utils.Errorhandler(err, "Error Retrieving deleted student")
		}
		if rowsAffected > 0 {
			deletedIds = append(deletedIds, id)
		}
		if rowsAffected > 1 {
			tx.Rollback()
			return nil, utils.Errorhandler(err, fmt.Sprintf("ID %d is not found", id))
		}

	}

	err = tx.Commit()
	if err != nil {
		// log.Println(err)
		// http.Error(w, "Error commiting transaction", http.StatusBadRequest)
		return nil, utils.Errorhandler(err, "Error commiting transaction")
	}
	if len(deletedIds) < 1 {
		// log.Println(err)
		// http.Error(w, "Id is not exist", http.StatusBadRequest)
		return nil, utils.Errorhandler(err, "Id is not exist")
	}
	return deletedIds, nil
}

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
)

func GetTeacherDbHandler(teachers []models.Teacher, r *http.Request) ([]models.Teacher, error) {
	db, err := ConnectDB()
	if err != nil {
		return nil, utils.Errorhandler(err, "Error Connecting to database")
	} else {
		fmt.Println("Database connected from GetTeacherHandler")
	}
	defer db.Close()

	query := "select id,first_name,last_name,email,class,subject from teachers where 1=1"
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
		var teacher models.Teacher
		err := rows.Scan(&teacher.Id, &teacher.FirstName, &teacher.LastName, &teacher.Email, &teacher.Class, &teacher.Subject)
		if err != nil {
			fmt.Println(err)
			// http.Error(w, "", http.StatusInternalServerError)
			return nil, utils.Errorhandler(err, "Error Scanning database error")
		}
		teachers = append(teachers, teacher)
	}
	return teachers, nil

}

func GetTeacherByIdDbHandler(id int) (models.Teacher, error) {
	db, err := ConnectDB()
	if err != nil {
		// http.Error(w, err.Error(), http.StatusInternalServerError)
		return models.Teacher{}, utils.Errorhandler(err, "Error Connecting to database")
	} else {
		fmt.Println("Database connected from GetTeacherHandler")
	}
	defer db.Close()

	var teacher models.Teacher
	err = db.QueryRow("select id,first_name,last_name,email,class,subject from teachers where id=?", id).Scan(&teacher.Id, &teacher.FirstName,
		&teacher.LastName, &teacher.Email, &teacher.Class, &teacher.Subject)
	if err == sql.ErrNoRows {
		// http.Error(w, "Teacher not found", http.StatusNotFound)
		return models.Teacher{}, utils.Errorhandler(err, "Teacher not found")
	} else if err != nil {
		// http.Error(w, "Database query error", http.StatusInternalServerError)
		return models.Teacher{}, utils.Errorhandler(err, "Database query error")
	}
	return teacher, nil
}

func AddTeachersDbHandler(newTeachers []models.Teacher) ([]models.Teacher, error) {
	db, err := ConnectDB()
	if err != nil {
		// http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil, utils.Errorhandler(err, "Database Connection error")
	} else {
		fmt.Println("Database connection established in AddTeacherHandler")
	}
	defer db.Close()

	// stmt, err := db.Prepare("insert into teachers(first_name, last_name, email, class, subject) values (?, ?, ?, ?, ?)")
	stmt, err := db.Prepare(utils.GenerateInsertQuery("teachers", models.Teacher{}))
	if err != nil {
		// http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil, utils.Errorhandler(err, "error preparing stament")
	}
	defer stmt.Close()

	addedTeachers := make([]models.Teacher, len(newTeachers))
	for i, teacher := range newTeachers {
		// res, err := stmt.Exec(teacher.FirstName, teacher.LastName, teacher.Email, teacher.Class, teacher.Subject)
		values := utils.GetStructValues(teacher)
		res, err := stmt.Exec(values...)
		if err != nil {
			// http.Error(w, err.Error(), http.StatusInternalServerError)
			return nil, utils.Errorhandler(err, "error in adding new teacher")
		}
		lastID, err := res.LastInsertId()
		if err != nil {
			// http.Error(w, err.Error(),http.StatusInternalServerError)
			// http.Error(w, "Error getting last insert ID", http.StatusInternalServerError)
			return nil, utils.Errorhandler(err, "Error getting last insert ID")
		}
		teacher.Id = int(lastID)
		addedTeachers[i] = teacher

	}
	return addedTeachers, nil
}

func UpdateTeacherDbHandler(id int, Updates models.Teacher) (models.Teacher, error) {
	db, err := ConnectDB()
	if err != nil {
		fmt.Println(err)
		// http.Error(w, "Unable to connect database", http.StatusInternalServerError)
		return models.Teacher{}, utils.Errorhandler(err, "Database Connection error")
	} else {
		fmt.Println("Database connected from PUTteacherHandler")
	}
	defer db.Close()

	var existingTeacher models.Teacher
	err = db.QueryRow("SELECT id,first_name,last_name,email,class,subject FROM teachers WHERE id =?", id).Scan(&existingTeacher.Id,
		&existingTeacher.FirstName, &existingTeacher.LastName, &existingTeacher.Email, &existingTeacher.Class, &existingTeacher.Subject)

	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println(err)
			// http.Error(w, "Teacher not found", http.StatusNotFound)
			return models.Teacher{}, utils.Errorhandler(err, "Teacher not found")
		}
		fmt.Println(err)
		// http.Error(w, "Unable to retrieve data", http.StatusInternalServerError)
		return models.Teacher{}, utils.Errorhandler(err, "Unable to retrieve data")
	}
	Updates.Id = existingTeacher.Id
	_, err = db.Exec("UPDATE teachers SET first_name =?, last_name=?,email=?, class=?, subject=? WHERE id=?", Updates.FirstName,
		Updates.LastName, Updates.Email, Updates.Class, Updates.Subject, Updates.Id)
	if err != nil {
		fmt.Println(err)
		// http.Error(w, "Error while updating", http.StatusInternalServerError)
		return models.Teacher{}, utils.Errorhandler(err, "Error while updating")
	}
	return Updates, nil
}

func PatchTeacherDbHandler(Updates []map[string]interface{}) error {
	db, err := ConnectDB()
	if err != nil {
		// http.Error(w, "Unable to connect to database", http.StatusInternalServerError)
		return utils.Errorhandler(err, "Unable to connect to database")
	} else {
		fmt.Println("Database connected from PatchTeacherHandlerr")
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
			// http.Error(w, "Invalid teacher ID in update", http.StatusBadRequest)
			return utils.Errorhandler(err, "Invalid teacher ID in update")
		}

		id, err := strconv.Atoi(idStr)
		if err != nil {
			tx.Rollback()
			// http.Error(w, "Invalid Teacher ID", http.StatusInternalServerError)
			return utils.Errorhandler(err, "Invalid Teacher ID")
		}
		var teacherFromDb models.Teacher
		err = db.QueryRow("SELECT id,first_name,last_name,email,class,subject FROM teachers WHERE id=?", id).Scan(&teacherFromDb.Id,
			&teacherFromDb.FirstName, &teacherFromDb.LastName, &teacherFromDb.Email, &teacherFromDb.Class, &teacherFromDb.Subject)
		if err != nil {
			tx.Rollback()
			if err == sql.ErrNoRows {
				print(err)
				// http.Error(w, "Teacher Not Found", http.StatusNotFound)
				return utils.Errorhandler(err, "Teacher Not Found")
			}
			// http.Error(w, "Error Retrieving Teacher", http.StatusInternalServerError)
			return utils.Errorhandler(err, "Error Retrieving Teacher")
		}

		teacherVal := reflect.ValueOf(&teacherFromDb).Elem()
		techerType := teacherVal.Type()

		for k, v := range update {
			if k == "id" {
				continue
			}
			for i := 0; i < teacherVal.NumField(); i++ {
				field := techerType.Field(i)
				if field.Tag.Get("json") == k+",omitempty" {
					fieldVal := teacherVal.Field(i)
					if teacherVal.CanSet() {
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
		_, err = db.Exec("UPDATE teachers SET first_name =?, last_name=?,email=?, class=?, subject=? WHERE id=?",
			teacherFromDb.FirstName, teacherFromDb.LastName, teacherFromDb.Email, teacherFromDb.Class, teacherFromDb.Subject, teacherFromDb.Id)
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

func PatchTeacherByIdDbHandler(id int, Updates map[string]interface{}) (models.Teacher, error) {
	db, err := ConnectDB()
	if err != nil {
		// fmt.Println(err)
		// http.Error(w, "Unable to connect database", http.StatusInternalServerError)
		return models.Teacher{}, utils.Errorhandler(err, "Unable to connect to database")
	} else {
		fmt.Println("Database connected from PatchTeacherByIdHandlerr")
	}
	defer db.Close()

	var existingTeacher models.Teacher
	err = db.QueryRow("SELECT id,first_name,last_name,email,class,subject FROM teachers WHERE id =?", id).Scan(&existingTeacher.Id,
		&existingTeacher.FirstName, &existingTeacher.LastName, &existingTeacher.Email, &existingTeacher.Class, &existingTeacher.Subject)

	if err != nil {
		if err == sql.ErrNoRows {
			// fmt.Println(err)
			// http.Error(w, "Teacher not found", http.StatusNotFound)
			return models.Teacher{}, utils.Errorhandler(err, "Teacher not found")
		}
		fmt.Println(err)
		// http.Error(w, "Unable to retrieve data", http.StatusInternalServerError)
		return models.Teacher{}, utils.Errorhandler(err, "Unable to retrieve data")
	}
	teacherVal := reflect.ValueOf(&existingTeacher).Elem()
	techerType := teacherVal.Type()

	for k, v := range Updates {
		for i := 0; i < teacherVal.NumField(); i++ {
			field := techerType.Field(i)
			if field.Tag.Get("json") == k+",omitempty" {
				if teacherVal.Field(i).CanSet() {
					fieldVal := teacherVal.Field(i)
					fmt.Println("fieldVal", fieldVal)
					fmt.Println("teacherVal.Field(i).Type()", teacherVal.Field(i).Type())
					fmt.Println("reflect.ValueOf(v)", reflect.ValueOf(v))
					fieldVal.Set(reflect.ValueOf(v).Convert(teacherVal.Field(i).Type()))
				}
			}
		}
	}

	_, err = db.Exec("UPDATE teachers SET first_name =?, last_name=?,email=?, class=?, subject=? WHERE id=?", existingTeacher.FirstName,
		existingTeacher.LastName, existingTeacher.Email, existingTeacher.Class, existingTeacher.Subject, existingTeacher.Id)
	if err != nil {
		// fmt.Println(err)
		// http.Error(w, "Error while updating", http.StatusInternalServerError)
		return models.Teacher{}, utils.Errorhandler(err, "Error while updating")
	}
	return existingTeacher, nil
}

func DeleteTeacherByIdDbHandler(id int) error {
	db, err := ConnectDB()
	if err != nil {
		// fmt.Println(err)
		// http.Error(w, "Unable to connect database", http.StatusInternalServerError)
		return utils.Errorhandler(err, "Unable to connect database")
	} else {
		fmt.Println("Database connected from DELETEByIdteacherHandler")
	}
	defer db.Close()

	result, err := db.Exec("DELETE from teachers WHERE id =?", id)
	if err != nil {
		// http.Error(w, "Error deleting teacher", http.StatusInternalServerError)
		return utils.Errorhandler(err, "Error deleting teacher")
	}

	fmt.Println(result.RowsAffected())
	rowsAffect, err := result.RowsAffected()
	if err != nil {
		// http.Error(w, "Error retrived deleted result", http.StatusInternalServerError)
		return utils.Errorhandler(err, "Error retrived deleted result")
	}
	if rowsAffect == 0 {
		// http.Error(w, "Teacher not found", http.StatusNotFound)
		return utils.Errorhandler(err, "Teacher not found")
	}
	return nil
}

func DeleteTeacherDbHandler(ids []int) ([]int, error) {
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
	stmt, err := tx.Prepare("DELETE FROM teachers WHERE id =?")
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
			// http.Error(w, "Error deleting teacher", http.StatusInternalServerError)
			return nil, utils.Errorhandler(err, "Error deleting teacher")
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			// log.Println(err)
			tx.Rollback()
			// http.Error(w, "Error Retrieving deleted teacher", http.StatusInternalServerError)
			return nil, utils.Errorhandler(err, "Error Retrieving deleted teacher")
		}
		if rowsAffected > 0 {
			deletedIds = append(deletedIds, id)
		}
		if rowsAffected < 1 {
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


func GetStudentsByTeacherIdFromDb(teacherId string, students []models.Student) ([]models.Student, error) {
	db, err := ConnectDB()
	if err != nil {
		return nil, utils.Errorhandler(err, "error retrieving data")
	}
	defer db.Close()

	query := `SELECT id, first_name, last_name, email, class FROM students WHERE class = (SELECT class from teachers WHERE id = ?)`
	rows, err := db.Query(query, teacherId)
	if err != nil {
		return nil, utils.Errorhandler(err, "error retrieving data")
	}
	defer rows.Close()

	for rows.Next() {
		var student models.Student
		err := rows.Scan(&student.Id, &student.FirstName, &student.LastName, &student.Email, &student.Class)
		if err != nil {
			return nil, utils.Errorhandler(err, "error retrieving data")
		}
		students = append(students, student)
	}
	err = rows.Err()
	if err != nil {
		return nil, utils.Errorhandler(err, "error retrieving data")
	}
	return students, nil
}

func GetStudentCountByTeacherIdFromDb(teacherId string) (int, error) {
	db, err := ConnectDB()
	if err != nil {
		return 0, utils.Errorhandler(err, "error retrieving data")
	}

	defer db.Close()

	query := `SELECT COUNT(*) FROM students WHERE class = (SELECT class FROM teachers WHERE id = ?)`
	var studentCount int
	err = db.QueryRow(query, teacherId).Scan(&studentCount)
	if err != nil {
		return 0, utils.Errorhandler(err, "error retrieving data")
	}
	return studentCount, nil
}
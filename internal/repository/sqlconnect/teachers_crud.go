package sqlconnect

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"rest_api_go/internal/models"
	"rest_api_go/pkg/util"
	"strconv"
	"strings"
)

func addFilters(r *http.Request, query string, args []interface{}) (string, []interface{}) {
	params := map[string]string{
		"first_name": "first_name",
		"last_name":  "last_name",
		"email":      "email",
		"class":      "class",
		"subject":    "subject",
	}
	for param, dbField := range params {
		value := r.URL.Query().Get(param)
		if value != "" {
			query += "AND" + dbField + "=?"
			args = append(args, value)
		}
	}
	return query, args
}

func AddSorting(r *http.Request, query string) string {
	sortParam := r.URL.Query()["sorted by"]
	if len(sortParam) > 0 {
		query += "ORDER BY"
		for i, param := range sortParam {
			parts := strings.Split(param, ":")
			if len(parts) != 2 {
				continue
			}
			field, order := parts[0], parts[1]
			if !isValidSortfield(field) || !isValidSortOrder(order) {
				continue
			}
			if i > 0 {
				query += ","
			}
			query += " " + field + " " + order

		}
	}
	return query

}

func isValidSortfield(field string) bool {
	validFields := map[string]bool{
		"first_name ": true,
		"last_name ":  true,
		"email ":      true,
		"class ":      true,
		"subject ":    true,
	}
	return validFields[field]
}
func isValidSortOrder(order string) bool {
	return order == "asc" || order == "desc"

}

func GetTeacherDbHandler(teachers []models.Teacher, r *http.Request) ([]models.Teacher, error) {
	db, err := ConnectDB()
	if err != nil {
		return nil, util.Errorhandler(err, "Error Connecting to database")
	} else {
		fmt.Println("Database connected from GetTeacherHandler")
	}
	defer db.Close()

	query := "select id,first_name,last_name,email,class,subject from teachers where 1=1"
	var args []interface{}
	query, args = addFilters(r, query, args)
	query = AddSorting(r, query)

	rows, err := db.Query(query, args...)
	if err != nil {
		fmt.Println(err)
		return nil, util.Errorhandler(err, "error retrieving database")
	}
	defer rows.Close()

	for rows.Next() {
		var teacher models.Teacher
		err := rows.Scan(&teacher.Id, &teacher.FirstName, &teacher.LastName, &teacher.Email, &teacher.Class, &teacher.Subject)
		if err != nil {
			fmt.Println(err)
			// http.Error(w, "", http.StatusInternalServerError)
			return nil, util.Errorhandler(err, "Error Scanning database error")
		}
		teachers = append(teachers, teacher)
	}
	return teachers, nil

}

func GetTeacherByIdDbHandler(id int) (models.Teacher, error) {
	db, err := ConnectDB()
	if err != nil {
		// http.Error(w, err.Error(), http.StatusInternalServerError)
		return models.Teacher{},util.Errorhandler(err, "Error Connecting to database")
	} else {
		fmt.Println("Database connected from GetTeacherHandler")
	}
	defer db.Close()

	var teacher models.Teacher
	err = db.QueryRow("select id,first_name,last_name,email,class,subject from teachers where id=?", id).Scan(&teacher.Id, &teacher.FirstName,
		&teacher.LastName, &teacher.Email, &teacher.Class, &teacher.Subject)
	if err == sql.ErrNoRows {
		// http.Error(w, "Teacher not found", http.StatusNotFound)
		return models.Teacher{}, util.Errorhandler(err, "Teacher not found")
	} else if err != nil {
		// http.Error(w, "Database query error", http.StatusInternalServerError)
		return models.Teacher{}, util.Errorhandler(err, "Database query error")
	}
	return models.Teacher{}, nil
}

func AddTeachersDbHandler(newTeachers []models.Teacher) ([]models.Teacher, error) {
	db, err := ConnectDB()
	if err != nil {
		// http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil, util.Errorhandler(err, "Database Connection error")
	} else {
		fmt.Println("Database connection established in AddTeacherHandler")
	}
	defer db.Close()

	// stmt, err := db.Prepare("insert into teachers(first_name, last_name, email, class, subject) values (?, ?, ?, ?, ?)")
	stmt, err := db.Prepare(GenerateInsertQuery("teachers", models.Teacher{}))
	if err != nil {
		// http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil, util.Errorhandler(err, "error preparing stament")
	}
	defer stmt.Close()

	addedTeachers := make([]models.Teacher, len(newTeachers))
	for i, teacher := range newTeachers {
		// res, err := stmt.Exec(teacher.FirstName, teacher.LastName, teacher.Email, teacher.Class, teacher.Subject)
		values := GetStructValues(teacher)
		res, err := stmt.Exec(values...)
		if err != nil {
			// http.Error(w, err.Error(), http.StatusInternalServerError)
			return nil, util.Errorhandler(err, "error in adding new teacher")
		}
		lastID, err := res.LastInsertId()
		if err != nil {
			// http.Error(w, err.Error(),http.StatusInternalServerError)
			// http.Error(w, "Error getting last insert ID", http.StatusInternalServerError)
			return nil, util.Errorhandler(err, "Error getting last insert ID")
		}
		teacher.Id = int(lastID)
		addedTeachers[i] = teacher

	}
	return addedTeachers, nil
}

func GenerateInsertQuery(tableName string, model interface{}) string {
	modelType := reflect.TypeOf(model)
	var columns, placeholders string
	for i := 0; i < modelType.NumField(); i++ {
		dbTag := modelType.Field(i).Tag.Get("db")
		fmt.Println("dbTag:", dbTag)
		dbTag = strings.TrimSuffix(dbTag, ",omitempty")
		if dbTag != "" && dbTag != "id" { // skip the ID field if it's auto increment
			if columns != "" {
				columns += ", "
				placeholders += ", "
			}
			columns += dbTag
			placeholders += "?"

		}
	}
	fmt.Printf("INSERT INTO %s (%s) VALUES (%s)\n", tableName, columns, placeholders)
	return fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", tableName, columns, placeholders)
}

func GetStructValues(model interface{}) []interface{} {
	modelValue := reflect.ValueOf(model)
	modelType := modelValue.Type()
	values := []interface{}{}
	for i := 0; i < modelType.NumField(); i++ {
		dbTag := modelType.Field(i).Tag.Get("db")
		if dbTag != "" && dbTag != "id,omitempty" {
			values = append(values, modelValue.Field(i).Interface())
		}
	}
	log.Println("Values:", values)
	return values
}

func UpdateTeacherDbHandler(id int, Updates models.Teacher) (models.Teacher, error) {
	db, err := ConnectDB()
	if err != nil {
		fmt.Println(err)
		// http.Error(w, "Unable to connect database", http.StatusInternalServerError)
		return models.Teacher{}, util.Errorhandler(err, "Database Connection error")
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
			return models.Teacher{}, util.Errorhandler(err, "Teacher not found")
		}
		fmt.Println(err)
		// http.Error(w, "Unable to retrieve data", http.StatusInternalServerError)
		return models.Teacher{}, util.Errorhandler(err, "Unable to retrieve data")
	}
	Updates.Id = existingTeacher.Id
	_, err = db.Exec("UPDATE teachers SET first_name =?, last_name=?,email=?, class=?, subject=? WHERE id=?", Updates.FirstName,
		Updates.LastName, Updates.Email, Updates.Class, Updates.Subject, Updates.Id)
	if err != nil {
		fmt.Println(err)
		// http.Error(w, "Error while updating", http.StatusInternalServerError)
		return models.Teacher{}, util.Errorhandler(err, "Error while updating")
	}
	return Updates, nil
}

func PatchTeacherDbHandler(Updates []map[string]interface{}) error {
	db, err := ConnectDB()
	if err != nil {
		// http.Error(w, "Unable to connect to database", http.StatusInternalServerError)
		return util.Errorhandler(err, "Unable to connect to database")
	} else {
		fmt.Println("Database connected from PatchTeacherHandlerr")
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		// print(err)
		// http.Error(w, "Error starting transaction", http.StatusInternalServerError)
		return util.Errorhandler(err, "Error starting transaction")
	}

	for _, update := range Updates {
		idStr, ok := update["id"].(string)
		if !ok {
			tx.Rollback()
			// http.Error(w, "Invalid teacher ID in update", http.StatusBadRequest)
			return util.Errorhandler(err, "Invalid teacher ID in update")
		}

		id, err := strconv.Atoi(idStr)
		if err != nil {
			tx.Rollback()
			// http.Error(w, "Invalid Teacher ID", http.StatusInternalServerError)
			return util.Errorhandler(err, "Invalid Teacher ID")
		}
		var teacherFromDb models.Teacher
		err = db.QueryRow("SELECT id,first_name,last_name,email,class,subject FROM teachers WHERE id=?", id).Scan(&teacherFromDb.Id,
			&teacherFromDb.FirstName, &teacherFromDb.LastName, &teacherFromDb.Email, &teacherFromDb.Class, &teacherFromDb.Subject)
		if err != nil {
			tx.Rollback()
			if err == sql.ErrNoRows {
				print(err)
				// http.Error(w, "Teacher Not Found", http.StatusNotFound)
				return util.Errorhandler(err, "Teacher Not Found")
			}
			// http.Error(w, "Error Retrieving Teacher", http.StatusInternalServerError)
			return util.Errorhandler(err, "Error Retrieving Teacher")
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
							return util.Errorhandler(err, "Error starting transaction")
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
			return util.Errorhandler(err, "Error while updating")
		}

	}
	err = tx.Commit()
	if err != nil {
		// http.Error(w, "Error comiting transaction", http.StatusInternalServerError)
		return util.Errorhandler(err, "Error comiting transaction")
	}
	return nil
}

func PatchTeacherByIdDbHandler(id int, Updates map[string]interface{}) (models.Teacher, error) {
	db, err := ConnectDB()
	if err != nil {
		// fmt.Println(err)
		// http.Error(w, "Unable to connect database", http.StatusInternalServerError)
		return models.Teacher{}, util.Errorhandler(err, "Unable to connect to database")
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
			return models.Teacher{}, util.Errorhandler(err, "Teacher not found")
		}
		fmt.Println(err)
		// http.Error(w, "Unable to retrieve data", http.StatusInternalServerError)
		return models.Teacher{}, util.Errorhandler(err, "Unable to retrieve data")
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
		return models.Teacher{}, util.Errorhandler(err, "Error while updating")
	}
	return existingTeacher, nil
}

func DeleteTeacherByIdDbHandler(id int) error {
	db, err := ConnectDB()
	if err != nil {
		// fmt.Println(err)
		// http.Error(w, "Unable to connect database", http.StatusInternalServerError)
		return util.Errorhandler(err, "Unable to connect database")
	} else {
		fmt.Println("Database connected from DELETEByIdteacherHandler")
	}
	defer db.Close()

	result, err := db.Exec("DELETE from teachers WHERE id =?", id)
	if err != nil {
		// http.Error(w, "Error deleting teacher", http.StatusInternalServerError)
		return util.Errorhandler(err, "Error deleting teacher")
	}

	fmt.Println(result.RowsAffected())
	rowsAffect, err := result.RowsAffected()
	if err != nil {
		// http.Error(w, "Error retrived deleted result", http.StatusInternalServerError)
		return util.Errorhandler(err, "Error retrived deleted result")
	}
	if rowsAffect == 0 {
		// http.Error(w, "Teacher not found", http.StatusNotFound)
		return util.Errorhandler(err, "Teacher not found")
	}
	return nil
}

func DeleteTeacherDbHandler(ids []int) ([]int, error) {
	db, err := ConnectDB()
	if err != nil {
		// http.Error(w, "Unable to connect database", http.StatusInternalServerError)
		return nil, util.Errorhandler(err, "Unable to connect database")
	} else {
		fmt.Println("Database Connected from DELETEHandler")
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		print(err)
		// http.Error(w, "Error starting transaction", http.StatusInternalServerError)
		return nil,util.Errorhandler(err, "Error starting transaction")
	}
	stmt, err := tx.Prepare("DELETE FROM teachers WHERE id =?")
	if err != nil {
		log.Println(err)
		tx.Rollback()
		// http.Error(w, "Error retrived deleted result", http.StatusInternalServerError)
		return nil,util.Errorhandler(err, "Error retrived deleted result")
	}
	defer stmt.Close()

	deletedIds := []int{}

	for _, id := range ids {
		result, err := stmt.Exec(id)
		if err != nil {
			// log.Println(err)
			tx.Rollback()
			// http.Error(w, "Error deleting teacher", http.StatusInternalServerError)
			return nil, util.Errorhandler(err, "Error deleting teacher")
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			// log.Println(err)
			tx.Rollback()
			// http.Error(w, "Error Retrieving deleted teacher", http.StatusInternalServerError)
			return nil, util.Errorhandler(err, "Error Retrieving deleted teacher")
		}
		if rowsAffected > 0 {
			deletedIds = append(deletedIds, id)
		}
		if rowsAffected < 1  {
			tx.Rollback()
			return nil, util.Errorhandler(err, fmt.Sprintf("ID %d is not found",id))
		}

	}

	err = tx.Commit()
	if err != nil {
		// log.Println(err)
		// http.Error(w, "Error commiting transaction", http.StatusBadRequest)
		return nil, util.Errorhandler(err, "Error commiting transaction")
	}
	if len(deletedIds) < 1 {
		// log.Println(err)
		// http.Error(w, "Id is not exist", http.StatusBadRequest)
		return nil, util.Errorhandler(err, "Id is not exist")
	}
	return deletedIds, nil
}

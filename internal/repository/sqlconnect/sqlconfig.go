package sqlconnect

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func ConnectDB() (*sql.DB, error) {

	// err := godotenv.Load()
	// if err !=nil{
	// 	return nil,err
	// }

	// dsn := fmt.Sprintf("root:dreamfight@tcp(localhost:3306)/%s", dbname)
	fmt.Println("Connecting to database:")

	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbname := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, host, port, dbname)
	// dsn := "root:dreamfight@tcp(localhost:3306)/" + dbname
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	// defer db.Close()

	// Test connection
	err = db.Ping()
	if err != nil {
	    log.Fatal("Error connecting to database:", err)
	}else {
		fmt.Println("Successfully connected to database")
	}

	// fmt.Println("Database connection establisdhed")
	return db, nil
}

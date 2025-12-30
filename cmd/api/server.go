package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"rest_api_go/internal/api/router"
	"rest_api_go/internal/repository/sqlconnect"

	"github.com/joho/godotenv"
	// "golang.org/x/net/http2"
)

func main() {
	// err := godotenv.Load()
	// if err != nil {
	// 	return nil, err
	// }
	err := godotenv.Load()
	if err != nil {
		return 
	}

	_, err = sqlconnect.ConnectDB()
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		return
	}

	port := os.Getenv("API_PORT")
	
	cert := "cert.pem"
	key := "key.pem"
	config := tls.Config{
		MinVersion: tls.VersionTLS12,
	}
	router := router.Router()

	server := &http.Server{
		Addr:      port,
		Handler:   router,
		TLSConfig: &config,
	}

	fmt.Printf("Server Up and Running at %v \n", port)
	err = server.ListenAndServeTLS(cert, key)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}

}

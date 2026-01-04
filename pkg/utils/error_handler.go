package utils

import (
	"fmt"
	"log"
	"os"
)

func Errorhandler(err error, message string) error {
	errorlogger := log.New(os.Stderr, "ERROR", log.Ldate|log.Ltime|log.Lshortfile)
	errorlogger.Println(message, err)
	return fmt.Errorf(message)

}

package handlers

import (
	"fmt"
	"net/http"
)

func GetExexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "This is execHandler")

}

func GetExecsHandler(w http.ResponseWriter, r *http.Request){

}
func AddExecsHandler(w http.ResponseWriter, r *http.Request){

}
func PatchExecsHandler(w http.ResponseWriter, r *http.Request){

}

func GetOneExecHandler(w http.ResponseWriter, r *http.Request){

}

func PatchOneExecHandler(w http.ResponseWriter, r *http.Request){

}

func DeleteOneExecHandler(w http.ResponseWriter, r *http.Request){

}

func UpdatePasswordHandler(w http.ResponseWriter, r *http.Request){

}

func LoginHandler(w http.ResponseWriter, r *http.Request){

}

func LogoutHandler(w http.ResponseWriter, r *http.Request){

}

func ForgotPasswordHandler(w http.ResponseWriter, r *http.Request){

}

func ResetPasswordHandler(w http.ResponseWriter, r *http.Request){

}

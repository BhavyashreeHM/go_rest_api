package router

import (
	"net/http"
)

func MainRouter() *http.ServeMux {

	tRouter := TeacherRouter()
	sRouter := StudentRouter()
	eRouter := ExecRouter()

	sRouter.Handle("/",eRouter)
	tRouter.Handle("/", sRouter)
	return tRouter
		// mux := http.NewServeMux()
		// mux.HandleFunc("GET /school", handlers.RootHandler)

	// 	mux.HandleFunc("GET /execs", handlers.ExecsHandler)
	// 	return mux
	//
}

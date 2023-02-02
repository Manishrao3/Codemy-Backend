package router

import (
	"RestfulApi/middleware"
	"github.com/gorilla/mux"
)

func Router() *mux.Router {

	router := mux.NewRouter()
	//fmt.Println(router)
	router.HandleFunc("/api/course/{id}", middleware.GetCourse).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/courses", middleware.GetAllCourses).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/newcourse", middleware.CreateCourse).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/course/{id}", middleware.UpdateCourse).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/deletecourse/{id}", middleware.DeleteCourse).Methods("DELETE", "OPTIONS")

	return router
}

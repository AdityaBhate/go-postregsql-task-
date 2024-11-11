package router

import (
	// _ "go-postgres/docs"
	"go-postgres/middleware"

	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/gorilla/mux"
)

func Router() *mux.Router {

	router := mux.NewRouter()

	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// Patient routes
	router.HandleFunc("/api/patient/{id}", middleware.GetPatient).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/newpatient", middleware.CreatePatient).Methods("POST", "OPTIONS")

	// Doctor routes
	router.HandleFunc("/api/doctor/{id}", middleware.GetDoctor).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/newdoctor", middleware.CreateDoctor).Methods("POST", "OPTIONS")

	return router
}

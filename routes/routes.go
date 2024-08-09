package routes

import (
	"cmdb-backend/controllers"
	"github.com/gorilla/mux"
)

func RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/api/assets", controllers.CreateAsset).Methods("POST")
	router.HandleFunc("/api/assets", controllers.GetAssets).Methods("GET")
}

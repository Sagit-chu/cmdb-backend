package routes

import (
	"cmdb-backend/controllers"

	"github.com/gorilla/mux"
)

func RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/api/assets", controllers.CreateAsset).Methods("POST")
	router.HandleFunc("/api/assets", controllers.GetAssets).Methods("GET")
	router.HandleFunc("/api/assets/{id:[0-9]+}", controllers.UpdateAsset).Methods("PUT")
	router.HandleFunc("/api/assets/{id:[0-9]+}", controllers.DeleteAsset).Methods("DELETE")
	router.HandleFunc("/api/assets/import", controllers.ImportAssets).Methods("POST")
}

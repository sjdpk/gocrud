package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sjdpk/gocrud/src/controllers"
)

func RegisterProductRoutes(router *mux.Router) {
	router.HandleFunc("/api/v1/product", controllers.CreateProduct).Methods(http.MethodPost)
	router.HandleFunc("/api/v1/product", controllers.GetAllProducts).Methods(http.MethodGet)
	router.HandleFunc("/api/v1/product/{id}", controllers.GetProduct).Methods(http.MethodGet)
	router.HandleFunc("/api/v1/product/{id}", controllers.UpdateProduct).Methods(http.MethodPut)
	router.HandleFunc("/api/v1/product/{id}", controllers.DeleteProduct).Methods(http.MethodDelete)
}

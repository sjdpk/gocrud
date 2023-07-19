package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sjdpk/gocrud/src/common"
	"github.com/sjdpk/gocrud/src/database"
	"github.com/sjdpk/gocrud/src/migrations"
	"github.com/sjdpk/gocrud/src/routes"
)

func main() {
	// Load Configurations
	common.LoadAppConfig()
	// Initialize Database
	database.Connect(common.AppConfig.DbConnectionString)
	migrations.Migration()

	// Initialize the router
	router := mux.NewRouter().StrictSlash(true)
	// Register Routers
	routes.RegisterProductRoutes(router)

	// start the server
	log.Println(fmt.Sprintf("Starting Server on port %s", common.AppConfig.Port))
	log.Fatal(http.ListenAndServe(fmt.Sprintf("localhost:%v", common.AppConfig.Port), router))

}

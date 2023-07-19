package main

import (
	"github.com/sjdpk/gocrud/src/common"
	"github.com/sjdpk/gocrud/src/database"
	"github.com/sjdpk/gocrud/src/migrations"
)

func main() {
	// Load Configurations
	common.LoadAppConfig()
	// Initialize Database
	database.Connect(common.AppConfig.DbConnectionString)
	migrations.Migration()

}

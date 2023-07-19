package migrations

import (
	"log"

	"github.com/sjdpk/gocrud/src/database"
	"github.com/sjdpk/gocrud/src/entities"
)

func Migration() {
	log.Println("ProductModel Migration Start...")
	database.Instance.AutoMigrate(&entities.Product{})
	log.Println("ProductModel Migration End")

}

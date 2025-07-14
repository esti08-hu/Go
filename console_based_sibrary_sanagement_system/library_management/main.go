package main

import (
	"library_management/controllers"
	"library_management/services"
)

func main() {
	service := services.NewLibraryService()
	controller := controllers.NewLibraryController(service)
	controller.Run()
}

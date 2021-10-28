package main

import (
	"Go-Echo-System/model"
	"Go-Echo-System/router"
)

func main() {
	model.InitModel()

	e := router.Router()
	e.Logger.Fatal(e.Start(":8080"))
}

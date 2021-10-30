package main

import (
	"Go-Echo-System/model"
	"Go-Echo-System/router"
)

func main() {
	model.InitModel()

	e := router.New()
	e.Logger.Fatal(e.Start(":8080"))
}

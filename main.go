package main

import (
	"log"

	"github.com/flowtr/ignite-api/pkg"
)

func main() {
	pkg.InitProviders()
	app := pkg.CreateApp()

	log.Fatal(app.Listen(":8008"))
}

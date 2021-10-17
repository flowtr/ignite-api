package pkg

import (
	"log"

	"github.com/flowtr/ignite-api/pkg/handler"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/weaveworks/ignite/pkg/config"
	"github.com/weaveworks/ignite/pkg/providers"
	"github.com/weaveworks/ignite/pkg/providers/ignite"
	"github.com/weaveworks/ignite/pkg/util"
)

func InitProviders() {
	util.GenericCheckErr(util.TestRoot())

	// Create the directories needed for running
	util.GenericCheckErr(util.CreateDirectories())

	if err := providers.Populate(ignite.Preload); err != nil {
		log.Println("Unable to populate ignite preload providers")
		log.Fatal(err)
	}

	if err := config.ApplyConfiguration(""); err != nil {
		log.Println("Unable to apply ignite configuration")
		log.Fatal(err)
	}

	if err := providers.Populate(ignite.Providers); err != nil {
		log.Println("Unable to populate ignite providers")
		log.Fatal(err)
	}
}

func CreateApp() *fiber.App {
	app := fiber.New()

	app.Use(recover.New())
	app.Use(cors.New())

	app.Get("/vm", handler.GetVMS)
	app.Get("/vm/{id}", handler.GetVM)
	app.Post("/vm", handler.CreateVM)

	return app
}

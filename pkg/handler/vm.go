package handler

import (
	"errors"
	"strings"

	igniteUtil "github.com/flowtr/ignite-api/pkg/util"
	"github.com/gofiber/fiber/v2"
	"github.com/weaveworks/ignite/pkg/providers"
	"github.com/weaveworks/libgitops/pkg/runtime"
)

func GetVMS(c *fiber.Ctx) error {
	vmList, err := providers.Client.VMs().List()

	if err != nil {
		return c.Status(500).JSON(err)
	}

	return c.JSON(vmList)
}

func GetVM(c *fiber.Ctx) error {
	id := c.Params("id")

	if strings.TrimSpace(id) == "" {
		return c.Status(400).JSON(errors.New(
			"invalid id parameter",
		))
	}

	vm, err := providers.Client.VMs().Get(runtime.UID(
		id,
	))

	if err != nil {
		return c.Status(500).JSON(err)
	}

	return c.JSON(vm)
}

// TODO: expose helper functions
func CreateVM(c *fiber.Ctx) error {
	var data igniteUtil.VMData

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	vm := igniteUtil.CreateVM(data)

	return c.JSON(vm)
}

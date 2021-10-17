package handler

import (
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/weaveworks/ignite/pkg/apis/ignite"
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

func withSpec(obj *ignite.VM, spec ignite.VMSpec) *ignite.VM {
	obj.Spec = spec
	return obj
}

func CreateVM(c *fiber.Ctx) error {
	var data ignite.VMSpec

	vm := withSpec(providers.Client.VMs().New(), data)

	err := providers.Client.VMs().Set(vm)

	if err != nil {
		return c.Status(500).JSON(err)
	} else {
		return c.JSON(vm)
	}
}

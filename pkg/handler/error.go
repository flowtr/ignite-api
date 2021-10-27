package handler

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

type ErrorResponse struct {
	err string
}

func HandleError(c *fiber.Ctx, err error, code int) bool {
	if err == nil {
		return false
	}

	log.Printf("err %s\n", err)

	c.Status(code).JSON(ErrorResponse{
		err: err.Error(),
	})
	return true
}

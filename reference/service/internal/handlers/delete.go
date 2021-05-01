package handlers

import (
	"fmt"

	"github.com/floppyzedolfin/sync/reference/client/reference"
	"github.com/floppyzedolfin/sync/reference/service/internal/endpoints/delete"
	"github.com/gofiber/fiber/v2"
)

// Delete parses the request and calls the endpoint
func Delete(ctx *fiber.Ctx) error {
	req := new(reference.DeleteRequest)
	if err := ctx.BodyParser(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("unable to parse body as a rqeuest: %s", err.Error())})
	}

	// call the endpoint
	res, err := delete.Delete(ctx, req)
	if err != nil {
		return ctx.Status(err.Code).JSON(fiber.Map{"error": err.Message})
	}
	return ctx.Status(fiber.StatusNoContent).JSON(res)
}

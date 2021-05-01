package handlers

import (
	"fmt"

	reference "github.com/floppyzedolfin/sync/reference/client/reference"
	patch "github.com/floppyzedolfin/sync/reference/service/internal/endpoints/patch"
	"github.com/gofiber/fiber/v2"
)

// Patch parses the request and calls the endpoint
func Patch(ctx *fiber.Ctx) error {
	req := new(reference.PatchRequest)
	if err := ctx.BodyParser(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("unable to parse body as a rqeuest: %s", err.Error())})
	}

	// call the endpoint
	res, err := patch.Patch(ctx, req)
	if err != nil {
		return ctx.Status(err.Code).JSON(fiber.Map{"error":err.Message})
	}
	return ctx.Status(fiber.StatusAccepted).JSON(res)
}

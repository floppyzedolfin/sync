package handlers

import (
	"fmt"

	reference2 "github.com/floppyzedolfin/sync/reference/client/reference"
	notify2 "github.com/floppyzedolfin/sync/reference/service/internal/endpoints/notify"
	"github.com/gofiber/fiber/v2"
)

func Notify(ctx *fiber.Ctx) error {
	req := new(reference2.NotifyRequest)
	if err := ctx.BodyParser(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("unable to parse body as a rqeuest: %s", err.Error())})
	}

	// call the endpoint
	res, err := notify2.Notify(ctx, req)
	if err != nil {
		return ctx.Status(err.Code).JSON(fiber.Map{"error":err.Message})
	}
	return ctx.Status(fiber.StatusAccepted).JSON(res)
}

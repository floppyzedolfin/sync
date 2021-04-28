package internal

import (
	"fmt"

	handlers2 "github.com/floppyzedolfin/sync/reference/service/internal/handlers"
	"github.com/gofiber/fiber/v2"
)

type Service struct {
	app *fiber.App
}

func New() *Service {
	var s Service
	s.app = fiber.New()
	s.registerEndpoints()
	return &s
}

func (s *Service) Listen(port int) {
	s.app.Listen(fmt.Sprintf(":%d", port))
}

func (s *Service) registerEndpoints() {
	s.app.Post("/notify", handlers2.Notify)
}

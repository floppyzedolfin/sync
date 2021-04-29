package internal

import (
	"fmt"

	handlers2 "github.com/floppyzedolfin/sync/reference/service/internal/handlers"
	"github.com/gofiber/fiber/v2"
)

type Service struct {
	app *fiber.App
}

// New returns a fully operation service, implementing all the required endpoints
func New() *Service {
	var s Service
	s.app = fiber.New()
	s.registerEndpoints()
	return &s
}

// Listen to a port for messages
func (s *Service) Listen(port int) {
	s.app.Listen(fmt.Sprintf(":%d", port))
}

// registerEndpoints adds each endpoint to the service
func (s *Service) registerEndpoints() {
	s.app.Post("/notify", handlers2.Notify)
}

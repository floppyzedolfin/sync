package internal

import (
	"fmt"

	handlers "github.com/floppyzedolfin/sync/reference/service/internal/handlers"
	"github.com/gofiber/fiber/v2"
)

// Service is our backend. It can listen to a port for messages.
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
	s.app.Post( , handlers.Patch)
	s.app.Post("/delete", handlers.Delete)
}

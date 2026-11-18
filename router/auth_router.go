package router

import (
	"tiket-event/controllers"
	"tiket-event/middleware"

	"github.com/gofiber/fiber/v3"
)

func NewRouter(server *fiber.App, controller controllers.AuthController, jwtsecret string) {
	// kelompok route dengan awalan /api/user
	api := server.Group("/api/user")

	// kelompokan router dangan awalan /api/user
	auth := api.Group("/group")
	auth.Post("/register", controller.Register)
	auth.Post("/login", controller.Login)

	api.Get("/me", middleware.AuthMiddleware(jwtsecret), func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"userId": c.Locals("userId"),
			"role":   c.Locals("role"),
		})
	})
}

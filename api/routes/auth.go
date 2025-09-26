package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lukashonok/micro-fiber-pet/api/handlers"
	"github.com/lukashonok/micro-fiber-pet/api/middleware"
	"github.com/lukashonok/micro-fiber-pet/internal/mq"
)

func AuthRouter(app fiber.Router, defaultServices mq.DefaultServices) {
	// secured
	protected := app.Group("/", middleware.JWTMiddleware(defaultServices.FirebaseAuth))
	protected.Get("/auth/me", handlers.MeHandler())
}

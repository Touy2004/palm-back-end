package main

import (
	"log"

	"github.com/Touy2004/palm-back-end/config"
	"github.com/Touy2004/palm-back-end/internal/handler"
	"github.com/Touy2004/palm-back-end/internal/middleware"
	"github.com/Touy2004/palm-back-end/internal/repository"
	"github.com/Touy2004/palm-back-end/internal/service"
	"github.com/Touy2004/palm-back-end/pkg/database"
	"github.com/Touy2004/palm-back-end/pkg/jwt"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Load config
	cfg := config.Load()

	// Connect to database
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatal("Error connecting to database: ", err)
	}

	// Init JWT
	jwtPkg := jwt.New(cfg.JWTSecret, cfg.JWTExpiry)

	// Init repositories
	userRepo := repository.NewUserRepository(db)

	// Init services
	authService := service.NewAuthService(userRepo, jwtPkg)

	// Init handlers
	h := &handler.Handler{
		Auth: handler.NewAuthHandler(authService),
	}

	// Init middlewares
	m := &handler.Middleware{
		Auth: middleware.NewAuthMiddleware(jwtPkg),
		Role: middleware.NewRoleMiddleware(),
	}

	// Init fiber
	app := fiber.New()

	// Setup routes
	handler.SetupRoutes(app, h, m)

	// Start server
	log.Fatal(app.Listen(":" + cfg.AppPort))
}

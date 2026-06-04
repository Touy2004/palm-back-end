package main

import (
	"fmt"
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
	// Load .env if it exists
	_ = godotenv.Load()

	// Load config
	cfg := config.Load()

	// Connect to database
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatal("Error connecting to database: ", err)
	}
	defer db.Close()

	// Init JWT
	jwtPkg := jwt.New(cfg.JWTSecret, cfg.JWTExpiry)

	// Init repositories
	userRepo := repository.NewUserRepository(db)
	deviceRepo := repository.NewDeviceRepository(db)
	attendanceRepo := repository.NewAttendanceRepository(db)
	pairingRepo := repository.NewPairingRepository(db)
	palmRepo := repository.NewPalmRepository(db)
	attemptRepo := repository.NewAttemptRepository(db)

	// Init Crypto Service (using dummy 32-byte key for dev)
	cryptoSvc, err := service.NewCryptoService("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=")
	if err != nil {
		log.Fatal("Error initializing crypto service: ", err)
	}

	// Init services
	authService := service.NewAuthService(userRepo, jwtPkg)
	adminService := service.NewAdminService(userRepo, deviceRepo, attendanceRepo)
	userService := service.NewUserService(userRepo, palmRepo, attendanceRepo)
	deviceService := service.NewDeviceService(deviceRepo, pairingRepo)
	pairingService := service.NewPairingService(pairingRepo)
	palmService := service.NewPalmService(palmRepo, pairingRepo, cryptoSvc, attemptRepo)
	attendanceService := service.NewAttendanceService(attendanceRepo, palmRepo, deviceRepo, userRepo, attemptRepo, cryptoSvc)

	// Init handlers
	h := &handler.Handler{
		Auth:   handler.NewAuthHandler(authService),
		Admin:  handler.NewAdminHandler(adminService),
		User:   handler.NewUserHandler(userService, pairingService),
		Device: handler.NewDeviceHandler(deviceService, palmService, attendanceService),
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
	for _, route := range app.GetRoutes() {
		fmt.Printf("%s %s\n", route.Method, route.Path)
	}
	// Start server
	log.Fatal(app.Listen(":" + cfg.AppPort))
}

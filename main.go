package main

import (
	"log"
	"os"
	"tiket-event/config"
	"tiket-event/controllers"
	"tiket-event/repository"
	"tiket-event/router"
	"tiket-event/services"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/joho/godotenv"
)

func main() {
	// baca file .env, kalau tidak ada (misal di production) cukup dicatat, tidak berhenti
	err := godotenv.Load()
	if err != nil {
		log.Println("file .env tidak ditemukan, pakai environment variable sistem")
	}

	port := os.Getenv("APP_PORT")
	jwtsecret := os.Getenv("JWT_SECRET")

	if jwtsecret == "" {
		log.Fatal("JWT Secret belum di isi")
	}

	db := config.NewDB()
	defer db.Close()

	authRepository := repository.NewRepository(db)
	authService := services.NewService(authRepository, jwtsecret)
	authController := controllers.NewController(authService)

	server := fiber.New()

	// middleware: catat setiap request di terminal
	server.Use(logger.New())
	// middleware: kalau ada panic di handler, app tidak mati, client dapat error 500
	server.Use(recover.New())

	router.NewRouter(server, authController, jwtsecret)

	err = server.Listen(":" + port)
	if err != nil {
		log.Fatal("Server behenti", err)
	}

}

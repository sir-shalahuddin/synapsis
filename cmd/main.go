package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sir-shalahuddin/synapsis/config"
	_ "github.com/sir-shalahuddin/synapsis/docs"
	db "github.com/sir-shalahuddin/synapsis/pkg/database"
	"github.com/sir-shalahuddin/synapsis/router"
)

// @title SYN STORE API
// @version 1.0
// @description bismillah lolos
// @BasePath /api/
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("no env provided")
	}

	AppConfig := config.AppConfig{Port: config.GetEnv("APP_PORT")}
	DBConfig := config.DBConfig{
		Host: config.GetEnv("DB_HOST"),
		Port: config.GetEnv("DB_PORT"),
		User: config.GetEnv("DB_USER"),
		Pass: config.GetEnv("DB_PASS"),
		Name: config.GetEnv("DB_NAME"),
	}
	JWTConfig := config.JWTConfig{
		Secret: config.GetEnv("JWT_SECRET"),
	}

	app := fiber.New()

	app.Use(cors.New())

	app.Get("/swagger/*", swagger.HandlerDefault)

	db, err := db.NewDB(DBConfig)
	if err != nil {
		panic(err)
	}

	router.SetupRoutes(app, db, JWTConfig)
	err = (app.Listen(":" + AppConfig.Port))
	if err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

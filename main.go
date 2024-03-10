package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/cuenobi/serviceDraft/config"
	"github.com/cuenobi/serviceDraft/middleware"
	"github.com/cuenobi/serviceDraft/service/entity"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/spf13/viper"
)

func init() {
	configFile := "./config.dev.yml"

	env := os.Getenv("APP_ENV")
	if env == "production" {
		log.Println(env)
		configFile = "./config.prod.yml"
	}

	viper.SetConfigFile(configFile)

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}
}

func main() {
	dbConn := config.InitDB()

	f := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
		// BodyLimit:   30 * 1024 * 1024,
	})

	corsAllowList := viper.GetStringSlice(`header.cors`)
	middL := config.CORS(corsAllowList)
	f.Use(middL)

	// timeFormat := "2006-01-02 15:04:05"
	loggerMiddleware := logger.New(logger.Config{
		TimeFormat: "2006-01-02 15:04:05",
		Format:     "${time} | ${status} | ${latency} | ${ips} | ${method} | ${path}\n",
	})
	f.Use(loggerMiddleware)

	f.Get("/ping", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "pong",
		})
	})

	type StaffRegisterBody struct {
		Username string `json:"username" validate:"required,min=6"`
		Password string `json:"password" validate:"required,min=6"`
	}
	f.Post("/post", func(c *fiber.Ctx) error {
		var input StaffRegisterBody

		// Parser input
		if err := c.BodyParser(&input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(err.Error())
		}

		// Validate input
		if err := middleware.Validate(input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(middleware.ErrorResponse(err, input))
		}

		user := &entity.User{
			Username:  &input.Username,
			Firstname: &input.Password,
		}

		result := dbConn.Create(user)
		if result.Error != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(result.Error)
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"username": user.Username,
			"password": user.Firstname,
		})
	})

	log.Fatal(f.Listen(fmt.Sprintf("localhost:%s", viper.GetString("server.port"))))
}

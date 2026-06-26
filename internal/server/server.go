package server

import (
	"github.com/gofiber/fiber/v2"

	"some-api/internal/database"
)

type FiberServer struct {
	*fiber.App

	db database.Service
}

func New() *FiberServer {
	server := &FiberServer{
		App: fiber.New(fiber.Config{
			ServerHeader: "some-api",
			AppName:      "some-api",
		}),

		db: database.New(),
	}

	return server
}

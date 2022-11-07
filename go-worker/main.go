package main

import (
	"go-worker/handler"

	"github.com/labstack/echo/v4"
)

func main() {

	server := echo.New()
	handler.RegisterRoutes(server)
	server.Logger.Fatal(server.Start(":8080"))

}

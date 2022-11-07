package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(server *echo.Echo) {
	server.GET("/", homeHandler)
}

func homeHandler(context echo.Context) error {
	return context.String(http.StatusOK, "Hello, World!")
}

package main

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type Handlers struct{}

func (h Handlers) GetStatus(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

package controllers

import (
	"net/http"

	"github.com/labstack/echo"
)

func Notfound(c echo.Context) error {
	return c.String(http.StatusOK, "404 Not found")
}

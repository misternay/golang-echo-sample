package controllers

import (
	"net/http"

	"github.com/labstack/echo"
)

func Index(c echo.Context) error {
	data := []interface{}{1, 2, "ok"}
	user := map[string]interface{}{"code": true, "data": data}
	return c.JSON(http.StatusOK, user)
}

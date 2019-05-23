package controllers

import (
	"net/http"

	"github.com/labstack/echo"
)

func Index(c echo.Context) error {
	type Data struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	type Response struct {
		Success bool   `json:"success"`
		Data    *Data  `json:"data,omitempty"`
		Message string `json:"message,omitempty"`
	}

	response := Response{
		Success: true,
		Data: &Data{
			Name: "John Due",
			Age:  18,
		},
		// Message: "sdlfj",
	}

	return c.JSON(http.StatusOK, response)
}

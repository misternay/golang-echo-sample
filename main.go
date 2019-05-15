package main

import (
	"net/http"

	"github.com/labstack/echo"
)

type user struct {
	Success bool   `json:"success"`
	Name    string `json:"name"`
}

func main() {
	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		user := &user{
			Success: true,
			Name:    "Hello",
		}
		return c.JSON(http.StatusOK, user)
	})

	e.Logger.Fatal(e.Start(":3000"))
}

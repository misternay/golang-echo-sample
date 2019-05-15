package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo"
)

type user struct {
	Success bool   `json:"success"`
	Name    string `json:"name"`
}

func getPort() string {
	var port = os.Getenv("PORT")
	if port == "" {
		port = "8080"
		fmt.Println("No Port In Heroku" + port)
	}
	return ":" + port
}

func main() {
	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		user := &user{
			Success: true,
			Name:    "Hello",
		}
		println("Hello naja heroku ja")
		return c.JSON(http.StatusOK, user)
	})

	e.Logger.Fatal(e.Start(getPort()))
}

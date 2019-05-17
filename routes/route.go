package routes

import (
	"os"

	"github.com/babyjazz/demo/controllers"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"gopkg.in/go-playground/validator.v9"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func getPort() string {
	var port = os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	return ":" + port
}

func Init() {
	e := echo.New()
	e.Pre(middleware.RemoveTrailingSlash())
	e.Validator = &CustomValidator{validator: validator.New()}

	api := e.Group("/api/v1")
	api.GET("", controllers.Index)
	api.POST("/login", controllers.Login)

	e.GET("/*", controllers.Notfound)

	e.Logger.Fatal(e.Start(getPort()))
}

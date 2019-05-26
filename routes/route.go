package routes

import (
	"os"

	"github.com/babyjazz/demo/config"
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

	auth := e.Group("/api/v1")
	auth.Use(middleware.JWT([]byte(config.Secret)))
	auth.POST("/register", controllers.RegisterChild)
	auth.GET("/child", controllers.GetChilds)
	auth.PATCH("/child", controllers.UpdateChild)
	auth.GET("/team", controllers.GetTeam)

	e.GET("/*", controllers.Notfound)

	e.Logger.Fatal(e.Start(getPort()))
}

package controllers

import (
	"net/http"

	"github.com/babyjazz/demo/handler"

	"github.com/babyjazz/demo/db"
	"github.com/babyjazz/demo/models"
	"github.com/labstack/echo"
)

func Login(c echo.Context) (err error) {
	type Response struct {
		Success bool          `json:"success"`
		Data    *models.Users `json:"data,omitempty"`
		Message string        `json:"message,omitempty"`
	}

	type Request struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	type Users struct {
		*models.Users
		Password string `json:"password"`
	}

	req := new(Request)
	if err = c.Bind(req); err != nil {
		return
	}

	pgdb := db.Connect()
	defer pgdb.Close()

	userModel := new(Users)

	err = pgdb.Model(userModel).Where("username=?", req.Username).First()
	if err != nil {
		res := &Response{
			Success: false,
			Message: "username or password is invalid",
		}
		return c.JSON(http.StatusUnauthorized, res)
	}
	if handler.CheckPasswordHash(req.Password, userModel.Password) == false {
		res := &Response{
			Success: false,
			Message: "username or password is invalid",
		}
		return c.JSON(http.StatusUnauthorized, res)
	}

	res := Response{
		Success: true,
	}

	return c.JSON(http.StatusOK, res)
}

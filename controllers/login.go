package controllers

import (
	"net/http"

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

	type Users struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	req := new(Users)
	if err = c.Bind(req); err != nil {
		return
	}

	pgdb := db.Connect()
	defer pgdb.Close()

	userModel := new(models.Users)

	err = pgdb.Model(userModel).Where("username=?", req.Username).Select()
	if err != nil {
		res := &Response{
			Success: false,
			Message: "username or password is invalid",
		}
		return c.JSON(http.StatusUnauthorized, res)
	}

	res := Response{
		Success: true,
		Data:    userModel,
	}

	return c.JSON(http.StatusOK, res)
}

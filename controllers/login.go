package controllers

import (
	"net/http"
	"strings"

	"github.com/babyjazz/demo/db"
	"github.com/babyjazz/demo/models"
	"github.com/labstack/echo"
	"gopkg.in/go-playground/validator.v9"
)

func Login(c echo.Context) (err error) {
	type Response struct {
		Success bool            `json:"success"`
		Data    *models.Members `json:"data,omitempty"`
		Message string          `json:"message,omitempty"`
	}

	type Body struct {
		Username string `json:"username" validate:"required,min=4,max=32"`
		Password string `json:"password" validate:"required,min=6,max=255"`
	}

	reqUser := new(Body)

	if err = c.Bind(reqUser); err != nil {
		return
	}

	if err = c.Validate(reqUser); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var fieldError string
			fieldError = strings.ToLower(err.Field() + "_" + err.ActualTag())
			if err.Param() != "" {
				fieldError = fieldError + "_" + err.Param()
			}
			res := &Response{
				Success: false,
				Message: fieldError,
			}
			return c.JSON(http.StatusBadRequest, res)
		}
	}

	pgdb := db.Connect()
	defer pgdb.Close()

	userModel := new(models.Members)

	err = pgdb.Model(userModel).Where("name=?", reqUser.Username).Select()
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

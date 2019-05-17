package controllers

import (
	"net/http"
	"strings"

	"github.com/labstack/echo"
	"gopkg.in/go-playground/validator.v9"
)

type body struct {
	Username string `json:"username" validate:"required,min=4,max=32"`
	Password string `json:"password" validate:"required,min=6,max=255"`
}

func Login(c echo.Context) (err error) {
	user := new(body)
	if err = c.Bind(user); err != nil {
		return
	}
	if err = c.Validate(user); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var fieldError string
			fieldError = strings.ToLower(err.Field() + "_" + err.ActualTag())
			if err.Param() != "" {
				fieldError = fieldError + "_" + err.Param()
			}
			res := map[string]interface{}{"success": false, "message": fieldError}
			return c.JSON(http.StatusBadRequest, res)
		}
	}

	return c.JSON(http.StatusOK, user)
}

package controllers

import (
	"net/http"

	"github.com/babyjazz/demo/db"
	"github.com/babyjazz/demo/handler"
	"github.com/labstack/echo"
	"gopkg.in/go-playground/validator.v9"
)

func Register(c echo.Context) (err error) {
	type (
		Response struct {
			Success bool   `json:"success"`
			Message string `json:"message,omitempty"`
		}

		Users struct {
			Username   string `json:"username" validate:"required,min=4,max=32"`
			Fullname   string `json:"fullname" validate:"required,min=4,max=32"`
			Password   string `json:"password" validate:"required,min=6,max=32"`
			Repassword string `json:"repassword" validate:"eqfield=Password"`
		}
	)

	req := new(Users)
	if err = c.Bind(req); err != nil {
		return
	}

	// Validate request
	trans := handler.TransValidator()
	if err = c.Validate(req); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			errorField, te := trans.T(err.Tag(), err.Field(), err.Param())
			if te != nil {
				errorField = "Invalid request"
			}
			res := Response{
				Success: false,
				Message: errorField,
			}
			return c.JSON(http.StatusBadRequest, res)
		}
	}

	pgdb := db.Connect()
	defer pgdb.Close()

	err = pgdb.Insert(req)
	if err != nil {
		res := &Response{
			Success: false,
		}
		return c.JSON(http.StatusUnauthorized, res)
	}

	response := Response{
		Success: true,
	}

	return c.JSON(http.StatusOK, response)
}

package controllers

import (
	"net/http"

	"github.com/babyjazz/demo/db"
	"github.com/babyjazz/demo/handler"
	"github.com/babyjazz/demo/models"
	"github.com/go-pg/pg"
	"github.com/labstack/echo"
	"gopkg.in/go-playground/validator.v9"
)

func GetChilds(c echo.Context) (err error) {
	type Request struct {
		Username string `json:"username" validate:"required"`
	}

	req := &Request{
		Username: c.Param("username"),
	}
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

	// Set parent
	userModel := new(models.Users)

	err = pgdb.Model(userModel).Where("username=?", req.Username).First()
	if err != nil {
		res := &Response{
			Success: true,
			Data:    userModel,
		}
		return c.JSON(http.StatusNotFound, res)
	}
	getChild(userModel, pgdb)

	res := &Response{
		Success: true,
		Data:    userModel,
	}

	return c.JSON(http.StatusOK, res)
}

func getChild(userModal *models.Users, pgdb *pg.DB) (err error) {
	childLeftNode := new(models.Users)
	childRightNode := new(models.Users)

	if userModal.ChildLeftId != 0 {
		err = pgdb.Model(childLeftNode).Where("id=?", userModal.ChildLeftId).First()
		userModal.ChildLeft = childLeftNode
		getChild(childLeftNode, pgdb)
	}

	if userModal.ChildRightId != 0 {
		err = pgdb.Model(childRightNode).Where("id=?", userModal.ChildRightId).First()
		userModal.ChildRight = childRightNode
		getChild(childRightNode, pgdb)
	}

	return nil
}

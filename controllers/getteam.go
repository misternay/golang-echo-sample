package controllers

import (
	"fmt"
	"net/http"

	"echo-sample/db"
	"echo-sample/models"
	"echo-sample/handler"
	"github.com/go-pg/pg"
	"github.com/labstack/echo"
	"gopkg.in/go-playground/validator.v9"
)

func GetTeam(c echo.Context) (err error) {
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
	var userModel []models.Users

	err = pgdb.Model(&userModel).Where("username=?", req.Username).Select()
	if err != nil {
		fmt.Println("User is not found")
	}
	userModel[0].Side = "parent" // Init side with left because tree will always go left first
	root := 0
	totalRoot := 0
	userModel[0].Root = &root
	getTeam(&userModel, 1, "left", &totalRoot, pgdb)

	res := &Response{
		Success: true,
		Data: map[string]interface{}{
			"total_root": totalRoot - 1, // Dec bcoz getTeam will go over dest of deepest child
			"childs":     userModel,
		},
	}

	return c.JSON(http.StatusOK, res)
}

func getTeam(userModal *[]models.Users, root int, direction string, totalRoot *int, pgdb *pg.DB) (err error) {
	childLeftNode := new(models.Users)
	childRightNode := new(models.Users)
	childLeft := (*userModal)[len(*userModal)-1].ChildLeftId
	childRight := (*userModal)[len(*userModal)-1].ChildRightId

	if root > *totalRoot {
		*totalRoot = root
	}

	if childLeft != 0 {
		err = pgdb.Model(childLeftNode).Where("id=?", childLeft).First()
		childLeftNode.Side = direction
		childLeftNode.Root = &root
		*userModal = append(*userModal, *childLeftNode)
		getTeam(userModal, root+1, direction, totalRoot, pgdb)
	}

	if childRight != 0 {
		err = pgdb.Model(childRightNode).Where("id=?", childRight).First()
		if root == 0 {
			childRightNode.Side = "right"
		} else {
			childRightNode.Side = direction
		}
		childRightNode.Root = &root
		*userModal = append(*userModal, *childRightNode)
		if root == 0 {
			getTeam(userModal, root+1, "right", totalRoot, pgdb)
		} else {
			getTeam(userModal, root+1, direction, totalRoot, pgdb)
		}
	}

	return nil
}

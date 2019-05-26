package controllers

import (
	"fmt"
	"net/http"

	"github.com/babyjazz/demo/db"
	"github.com/babyjazz/demo/models"
	"github.com/go-pg/pg"
	"github.com/labstack/echo"
)

func GetChilds(c echo.Context) (err error) {
	type Request struct {
		Username string `json:"username"`
	}

	req := new(Request)
	if err = c.Bind(req); err != nil {
		return
	}

	pgdb := db.Connect()
	defer pgdb.Close()

	// Set parent
	userModel := new(models.Users)

	err = pgdb.Model(userModel).Where("username=?", req.Username).First()
	if err != nil {
		fmt.Println("User is not found")
	}
	getChild(req.Username, userModel, pgdb)

	res := &Response{
		Success: true,
		Data:    userModel,
	}

	return c.JSON(http.StatusOK, res)
}

func getChild(username string, userModal *models.Users, pgdb *pg.DB) (err error) {
	childLeftNode := new(models.Users)
	childRightNode := new(models.Users)

	if userModal.ChildLeftId != 0 {
		err = pgdb.Model(childLeftNode).Where("id=?", userModal.ChildLeftId).First()
		userModal.ChildLeft = append(userModal.ChildLeft, childLeftNode)
		getChild(childLeftNode.Username, childLeftNode, pgdb)
	}

	if userModal.ChildRightId != 0 {
		err = pgdb.Model(childRightNode).Where("id=?", userModal.ChildRightId).First()
		userModal.ChildRight = append(userModal.ChildRight, childRightNode)
		getChild(childRightNode.Username, childRightNode, pgdb)
	}

	return nil
}

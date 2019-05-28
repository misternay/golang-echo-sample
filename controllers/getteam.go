package controllers

import (
	"fmt"
	"net/http"

	"echo-sample/db"
	"echo-sample/models"
	"github.com/go-pg/pg"
	"github.com/labstack/echo"
)

func GetTeam(c echo.Context) (err error) {
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
	var userModel []models.Users

	err = pgdb.Model(&userModel).Where("username=?", req.Username).Select()
	if err != nil {
		fmt.Println("User is not found")
	}
	userModel[0].Direction = "left" // Init side with left because tree will always go left first
	getTeam(&userModel, 0, "left", pgdb)

	res := &Response{
		Success: true,
		Data:    userModel,
	}

	return c.JSON(http.StatusOK, res)
}

func getTeam(userModal *[]models.Users, root int, direction string, pgdb *pg.DB) (err error) {
	childLeftNode := new(models.Users)
	childRightNode := new(models.Users)
	childLeft := (*userModal)[len(*userModal)-1].ChildLeftId
	childRight := (*userModal)[len(*userModal)-1].ChildRightId

	if childLeft != 0 {
		err = pgdb.Model(childLeftNode).Where("id=?", childLeft).First()
		childLeftNode.Direction = direction
		*userModal = append(*userModal, *childLeftNode)
		getTeam(userModal, root+1, direction, pgdb)
	}

	if childRight != 0 {
		err = pgdb.Model(childRightNode).Where("id=?", childRight).First()
		if root == 0 {
			childRightNode.Direction = "right"
		} else {
			childRightNode.Direction = direction
		}
		*userModal = append(*userModal, *childRightNode)
		if root == 0 {
			getTeam(userModal, root+1, "right", pgdb)
		} else {
			getTeam(userModal, root+1, direction, pgdb)
		}
	}

	return nil
}

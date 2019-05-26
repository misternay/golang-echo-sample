package controllers

import (
	"fmt"
	"net/http"

	"github.com/babyjazz/demo/db"
	"github.com/babyjazz/demo/models"
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
	getTeam(&userModel, 0, pgdb)

	res := &Response{
		Success: true,
		Data:    userModel,
	}

	return c.JSON(http.StatusOK, res)
}

func getTeam(userModal *[]models.Users, root int, pgdb *pg.DB) (err error) {
	childLeftNode := new(models.Users)
	childRightNode := new(models.Users)
	childLeft := (*userModal)[len(*userModal)-1].ChildLeftId
	childRight := (*userModal)[len(*userModal)-1].ChildRightId

	if childLeft != 0 {
		err = pgdb.Model(childLeftNode).Where("id=?", childLeft).First()
		*userModal = append(*userModal, *childLeftNode)
		getTeam(userModal, root+1, pgdb)
	}

	if childRight != 0 {
		err = pgdb.Model(childRightNode).Where("id=?", childRight).First()
		*userModal = append(*userModal, *childRightNode)
		getTeam(userModal, root+1, pgdb)
	}

	return nil
}

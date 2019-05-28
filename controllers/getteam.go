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

package controllers

import (
	"fmt"
	"net/http"

	"github.com/babyjazz/demo/db"
	"github.com/go-pg/pg"
	"github.com/labstack/echo"
)

type Users struct {
	Id           int      `json:"id,omitempty"`
	Username     string   `json:"username,omitempty"`
	Fullname     string   `json:"fullname,omitempty"`
	ChildLeftId  int      `json:"child_left_id,omitempty"`
	ChildRightId int      `json:"child_right_id,omitempty"`
	Child        []*Users `json:"child,omitempty" sql:"-"`
}

func GetChilds(c echo.Context) (err error) {
	type Request struct {
		Username string `json:"username"`
	}
	type Response struct {
		Success bool   `json:"success"`
		Data    *Users `json:"data,omitempty"`
		Message string `json:"message,omitempty"`
	}

	req := new(Request)
	if err = c.Bind(req); err != nil {
		return
	}

	pgdb := db.Connect()
	defer pgdb.Close()

	// Set parent
	userModel := new(Users)

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

func getChild(username string, userModal *Users, pgdb *pg.DB) (err error) {
	childLeftNode := new(Users)
	childRightNode := new(Users)

	if userModal.ChildLeftId != 0 {
		err = pgdb.Model(childLeftNode).Where("id=?", userModal.ChildLeftId).First()
		userModal.Child = append(userModal.Child, childLeftNode)
		getChild(childLeftNode.Username, childLeftNode, pgdb)
	}

	if userModal.ChildRightId != 0 {
		err = pgdb.Model(childRightNode).Where("id=?", userModal.ChildRightId).First()
		userModal.Child = append(userModal.Child, childRightNode)
		getChild(childRightNode.Username, childRightNode, pgdb)
	}

	return nil
}

package controllers

import (
	"fmt"
	"net/http"

	"github.com/babyjazz/demo/db"
	"github.com/go-pg/pg"
	"github.com/labstack/echo"
)

type Users struct {
	Id           int               `json:"id,omitempty"`
	Username     string            `json:"username,omitempty"`
	Fullname     string            `json:"fullname,omitempty"`
	ChildLeftId  int               `json:"child_left_id,omitempty"`
	ChildRightId int               `json:"child_right_id,omitempty"`
	Ok           map[string]*Users `json:",omitempty" sql:"-"`
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
	userModel.Ok = make(map[string]*Users)
	userModel.Ok["child"], err = getChild(req.Username, pgdb)

	res := &Response{
		Success: true,
		Data:    userModel,
	}

	return c.JSON(http.StatusOK, res)
}

func getChild(username string, pgdb *pg.DB) (user *Users, err error) {
	a := new(Users)
	s := new(Users)

	err = pgdb.Model(a).Where("username=?", username).First()
	if err != nil {
		fmt.Println("User is not found")
	}
	err = pgdb.Model(s).Where("id=?", a.ChildLeftId).First()
	if err != nil {
		fmt.Println("User is not found")
	}
	if s.ChildLeftId != 0 {
		s.Ok = make(map[string]*Users)
		s.Ok["child"], err = getChild(s.Username, pgdb)

	}
	return s, err
}

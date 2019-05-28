package controllers

import (
	"fmt"
	"net/http"

	"echo-sample/db"
	"echo-sample/models"
	"github.com/labstack/echo"
)

func UpdateChild(c echo.Context) (err error) {
	type Request struct {
		Username  string `json:"username"`
		Direction string `json:"direction"`
		ChildId   int    `json:"child_id"`
	}

	req := new(Request)
	if err = c.Bind(req); err != nil {
		return
	}

	pgdb := db.Connect()
	defer pgdb.Close()

	userModel := new(models.Users)

	// Find own id for update child id
	err = pgdb.Model(userModel).Where("username=?", req.Username).First()
	if err != nil || userModel.Id == 0 {
		res := &Response{
			Success: false,
			Message: "User is not found",
		}
		return c.JSON(http.StatusUnauthorized, res)
	}

	// Custom validate value
	if req.ChildId == userModel.Id {
		res := &Response{
			Success: false,
			Message: "child id must equal with user id",
		}
		return c.JSON(http.StatusUnauthorized, res)
	}

	// Update own child left id
	var updateField = fmt.Sprintf(`child_%s_id`, req.Direction)
	updated, err := pgdb.Model(userModel).Set(updateField+"=?", req.ChildId).Where("id=?", userModel.Id).Update()
	if err != nil {
		fmt.Println(err)
		res := &Response{
			Success: false,
			Message: "Update failed",
		}
		return c.JSON(http.StatusUnauthorized, res)
	}
	if updated.RowsAffected() == 0 {
		res := &Response{
			Success: false,
			Message: "User is not update",
		}
		return c.JSON(http.StatusUnauthorized, res)
	}

	res := &Response{
		Success: true,
		Data:    userModel,
	}

	return c.JSON(http.StatusOK, res)
}

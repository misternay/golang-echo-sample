package controllers

import (
	"net/http"
	"time"

	"github.com/babyjazz/demo/config"
	"github.com/babyjazz/demo/handler"
	"github.com/babyjazz/demo/models"
	"github.com/dgrijalva/jwt-go"

	"github.com/babyjazz/demo/db"
	"github.com/labstack/echo"
)

type (
	Response struct {
		Success     bool        `json:"success"`
		Data        interface{} `json:"data,omitempty"`
		Message     string      `json:"message,omitempty"`
		AccessToken string      `json:"accessToken,omitempty"`
		TotalRoot   *int        `json:"total_root,omitempty" sql:"-"`
	}
)

func Login(c echo.Context) (err error) {
	type Request struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	req := new(Request)
	if err = c.Bind(req); err != nil {
		return
	}

	pgdb := db.Connect()
	defer pgdb.Close()

	userModel := new(models.Users)

	err = pgdb.Model(userModel).Where("username=?", req.Username).First()
	if err != nil {
		res := &Response{
			Success: false,
			Message: "username or password is invalid",
		}
		return c.JSON(http.StatusUnauthorized, res)
	}
	if handler.CheckPasswordHash(req.Password, userModel.Password) == false {
		res := &Response{
			Success: false,
			Message: "username or password is invalid",
		}
		return c.JSON(http.StatusUnauthorized, res)
	}

	// Create jwt token
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["fullname"] = userModel.Fullname
	claims["username"] = userModel.Username
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	accessToken, err := token.SignedString([]byte(config.Secret))
	if err != nil {
		return err
	}

	res := &Response{
		Success:     true,
		Data:        userModel,
		AccessToken: accessToken,
	}

	return c.JSON(http.StatusOK, res)
}

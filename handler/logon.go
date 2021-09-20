package handler

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"

	"github.com/dpl10/monographia-server/util"
)

// JWTkey is the 256 bit signing key for HMAC-SHA256
var JWTkey []byte

// JWTlife is JWT lifetime in hours
var JWTlife time.Duration

// Logon handler function
func (h *Handler) Logon(c echo.Context) (err error) {
	username := c.FormValue("username")
	password := c.FormValue("password")

	// Throws unauthorized error
	if username != "jon" || password != "shhh!" {
		return echo.ErrUnauthorized
	}

	claims := &util.JWTclaims{
		"xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
		"dpl",
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * JWTlife).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString(JWTkey)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, echo.Map{
		"token": t,
	})
}

// test with:
// curl -X POST -d 'username=jon' -d 'password=shhh!' --insecure https://localhost:4420/api/logon
// response: {"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJQdWJsaWNLZXkiOiJ4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHgiLCJTY3JlZW5OYW1lIjoiZHBsIiwiZXhwIjoxNjEyOTM5MDc4fQ.qCcbg1lBb0Yxu-AllxwuDqdF4pWp6V9soNBT2CU4ILM"}
// curl https://localhost:4420/api/r/city --insecure -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJQdWJsaWNLZXkiOiJ4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHgiLCJTY3JlZW5OYW1lIjoiZHBsIiwiZXhwIjoxNjEyOTM5MDc4fQ.qCcbg1lBb0Yxu-AllxwuDqdF4pWp6V9soNBT2CU4ILM"

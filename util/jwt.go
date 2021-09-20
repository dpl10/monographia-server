package util

import (
	"github.com/dgrijalva/jwt-go"
)

type (
	// JWTclaims to hold Monographia specific data
	JWTclaims struct {
		PublicKey  string `json:"PublicKey"`
		ScreenName string `json:"ScreenName"`
		jwt.StandardClaims
	}
)

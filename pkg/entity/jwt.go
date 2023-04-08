package entity

import (
	"time"

	"github.com/go-chi/jwtauth"
)

func GenerateJWT(id string, jwt *jwtauth.JWTAuth, expiration int) (string, error) {
	_, token, err := jwt.Encode(map[string]interface{}{
		"sub": id,
		"exp": time.Now().Add(time.Second * time.Duration(expiration)).Unix(),
	})
	if err != nil {
		return "", err
	}
	return token, nil
}

package auth

import (
	"strconv"
	"time"
	"github.com/go-refresh-practice/go-refresh-course/config"
	"github.com/golang-jwt/jwt/v5"
)

func CreateJwt(secret []byte, userId int, userEmail, role string) (string, error) {
    expiration := time.Second * time.Duration(config.Envs.JWTExpirationInSeconds)

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "userId":    strconv.Itoa(userId),
        "userEmail": userEmail,
        "role":      role,    
        "expiredAt": time.Now().Add(expiration).Unix(),
    })

    tokenString, err := token.SignedString(secret)
    if err != nil {
        return "", err
    }

    return tokenString, nil
}

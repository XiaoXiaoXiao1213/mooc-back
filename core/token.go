package core

import (
	"fmt"
	"github.com/develop1024/jwt-go"
	log "github.com/sirupsen/logrus"
	"management/core/users"
	"time"
)

const secret = "management"

// 生成token
func GenerateToken(user users.User) (tokenString string, err error) {
	data := jwt.MapClaims{
		"phone":     user.Phone,
		"user_type": user.UserType,
		"user_id":   user.Id,
		"exp":       time.Now().Add(time.Hour * 3).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, data)
	log.Error(token)
	tokenString, err = token.SignedString([]byte(secret))
	return
}

// 解析token
func ParseToken(tokenString string) (*users.User, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		user := users.User{
			Id:       claims["user_id"].(int64),
			Phone:    claims["phone"].(string),
			UserType: claims["user_tpye"].(int),
		}

		return &user, nil
	} else {
		return nil, err
	}
}

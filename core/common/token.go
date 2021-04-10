package common

import (
	"errors"
	"fmt"
	"github.com/develop1024/jwt-go"
	"management/core/domain"
	"time"
)

const secret = "management"

// 生成token
func GenerateToken(user domain.User) (tokenString string, err error) {
	data := jwt.MapClaims{
		"phone":   user.Phone,
		"user_id": user.Id,
		"exp":     time.Now().Add(time.Hour * 3 * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, data)
	tokenString, err = token.SignedString([]byte(secret))
	return
}

// 解析token
func ParseToken(tokenString string) (*map[string]string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	user := make(map[string]string)
	if token == nil {
		return nil, errors.New("token 为空")
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		user["user_id"] = claims["user_id"].(string)
		user["phone"] = claims["phone"].(string)

		return &user, nil
	} else {
		return nil, err
	}
}

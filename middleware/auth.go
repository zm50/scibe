package middleware

import (
	"encoding/json"
	"errors"
	"scibe/global"
	"scibe/utils/resp"
	"scibe/utils/types"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func Auth(c *gin.Context) {
	token, err := c.Cookie("token")
	if err != nil {
		resp.AbortErr(c, err.Error())
		return
	}

	ctx, err := ParseToken(token)
	if err != nil {
		resp.AbortErr(c, err.Error())
		return
	}

	c.Set("ctx", ctx)
}

func GenToken(ctx types.Ctx) (string, error) {
	data, err := json.Marshal(ctx)
	if err != nil {
		return "", err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
		"iat":   time.Now().Unix(),
		"ctx":  string(data),
	})
	s, err := token.SignedString(global.GetJWTSecret())
	if err != nil {
		return "", err
	}
	return s, err
}

func ParseToken(token string) (*types.Ctx, error) {
	t, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return global.GetJWTSecret(), nil
	}, jwt.WithValidMethods([]string{"HS256"}))
	if err != nil {
		return nil, err
	}
	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("failed to parse claims")
	}
	infoData, ok := claims["ctx"].(string)
	if !ok {
		return nil, errors.New("failed to get info")
	}
	ctx := &types.Ctx{}
	err = json.Unmarshal([]byte(infoData), ctx)
	return ctx, err
}

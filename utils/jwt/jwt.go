package jwt

import (
	"encoding/json"
	"errors"
	"scibe/global"
	"scibe/utils/types"
	"time"

	"github.com/golang-jwt/jwt/v5"
)


func GenToken(prop types.Property) (string, error) {
	data, err := json.Marshal(prop)
	if err != nil {
		return "", err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
		"iat":   time.Now().Unix(),
		"nbf":   time.Now().Unix(),
		"prop":   string(data),
	})
	s, err := token.SignedString([]byte(global.Config().JWTSecret))
	if err != nil {
		return "", err
	}
	return s, err
}

func ParseToken(token string) (*types.Property, error) {
	t, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte(global.Config().JWTSecret), nil
	}, jwt.WithValidMethods([]string{"HS256"}))
	if err != nil {
		return nil, err
	}
	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("failed to parse claims")
	}
	infoData, ok := claims["prop"].(string)
	if !ok {
		return nil, errors.New("failed to get info")
	}
	prop := &types.Property{}
	err = json.Unmarshal([]byte(infoData), prop)
	return prop, err
}

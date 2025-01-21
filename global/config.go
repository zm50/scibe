package global

import "os"

const (
	JWT_SECRET = "JWT_SECRET"
)

func GetJWTSecret() string {
	secret, ok := os.LookupEnv(JWT_SECRET)
	if !ok {
		return "zuimojushi"
	}
	return secret
}

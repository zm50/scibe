package keys

import "fmt"

const (
	USER_TOKEN = "user.token"
)

func UserTokenKey(id uint) string {
	return fmt.Sprintf("%s.%d", USER_TOKEN, id)
}

package resp

import (
	"scibe/utils/types"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Ok(c *gin.Context, value interface{}) {
	c.JSON(http.StatusOK, types.Vo{Code: types.Success, Data: value})
}

func Err(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, types.Vo{Code: types.Failed, Message: message})
}

func AbortErr(c *gin.Context, message string) {
	c.Abort()
	c.JSON(http.StatusBadRequest, types.Vo{Code: types.Failed, Message: message})
}

func HACKER(c *gin.Context) {
	c.JSON(http.StatusBadRequest, types.Vo{Code: types.Failed, Message: "Hacker attempt!!!"})
}

func NotAuth(c *gin.Context, messages ...string) {
	if messages != nil {
		c.JSON(http.StatusUnauthorized, types.Vo{Code: types.NotAuthorized, Message: messages[0]})
	} else {
		c.JSON(http.StatusUnauthorized, types.Vo{Code: types.NotAuthorized, Message: "Not Authorized"})
	}
}

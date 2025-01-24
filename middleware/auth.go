package middleware

import (
	"scibe/utils/ret"
	"scibe/utils/jwt"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func Auth(c *gin.Context) {
	token, err := c.Cookie("token")
	if err != nil {
		log.Err(err).Msg("failed to get token")
		ret.AbortErr(c, err.Error())
		return
	}

	prop, err := jwt.ParseToken(token)
	if err != nil {
		log.Err(err).Msgf("failed to parse token: %v", token)
		ret.AbortErr(c, err.Error())
		return
	}

	c.Set("prop", prop)
}

package types

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Property struct {
	Uid uint `json:"uid"`
	Name string `json:"name"`
}

func GetProperty(c *gin.Context) *Property {
	return c.Value("prop").(*Property)
}

func (p *Property) Logger() zerolog.Logger {
	return log.With().Uint("uid", p.Uid).Str("name", p.Name).Logger()
}

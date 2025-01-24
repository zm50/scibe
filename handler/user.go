package handler

import (
	"scibe/global"
	"scibe/model"
	"scibe/utils/jwt"
	"scibe/utils/ret"
	"scibe/utils/types"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type LoginRequest struct {
	Name,
	Pass string
}

func Login(c *gin.Context) {
	req := LoginRequest{}
	err := c.BindJSON(&req)
	if err != nil {
		ret.AbortErr(c, err.Error())
		return
	}

	user := model.User{}
	if err := global.DB().Model(&model.User{}).Where("name = ? and pass = ?", req.Name, req.Pass).First(&user).Error; err != nil {
		ret.AbortErr(c, err.Error())
		return
	}

	prop := types.Property{
		Uid:  user.ID,
		Name: user.Name,
	}
	token, err := jwt.GenToken(prop)
	if err != nil {
		log.Err(err).Msg("failed to generate token")
		ret.AbortErr(c, err.Error())
		return
	}

	log.Info().Msgf("user %s login, id: %d", user.Name, user.ID)
	ret.Ok(c, gin.H{
		"token": token,
	})
}

type RegisterRequest struct {
	Name string `json:"name"`
	Pass string `json:"pass"`
}

func Register(c *gin.Context) {
	req := RegisterRequest{}
	err := c.BindJSON(&req)
	if err != nil {
		log.Err(err).Msg("failed to bind request")
		ret.AbortErr(c, err.Error())
		return
	}

	user := model.User{
		Name: req.Name,
		Pass: req.Pass,
	}

	if err := global.DB().Create(&user).Error; err != nil {
		ret.AbortErr(c, err.Error())
		return
	}

	ret.Ok(c, nil)
}

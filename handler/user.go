package handler

import (
	"scibe/global"
	"scibe/model"
	"scibe/utils/keys"
	"scibe/utils/resp"
	"scibe/utils/uuid"

	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Name,
	Pass string
}

func Login(c *gin.Context) {
	req := LoginRequest{}
	err := c.BindJSON(&req)
	if err != nil {
		resp.AbortErr(c, err.Error())
		return
	}

	user := model.User{}
	if err := global.DB().Model(&model.User{}).Where("name = ? && pass = ?", req.Name, req.Pass).First(&user).Error; err != nil {
		resp.AbortErr(c, err.Error())
		return
	}

	key := keys.UserTokenKey(user.ID)
	token := uuid.Gen()
	global.Set(key, token)

	resp.Ok(c, gin.H{
		"token": token,
	})
}

type RegisterRequest struct {
	Name,
	Pass string
}

func Register(c *gin.Context) {
	req := RegisterRequest{}
	err := c.BindJSON(&req)
	if err != nil {
		resp.AbortErr(c, err.Error())
		return
	}

	
	user := model.User{
		Name: req.Name,
		Pass: req.Pass,
	}

	if err := global.DB().Create(&user).Error; err != nil {
		resp.AbortErr(c, err.Error())
		return
	}

	resp.Ok(c, nil)
}

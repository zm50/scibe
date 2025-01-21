package handler

import (
	"crypto/md5"
	"fmt"
	"scibe/utils/io"
	sio "io"
	"os"
	"scibe/global"
	"scibe/model"
	"scibe/utils/resp"

	"github.com/gin-gonic/gin"
)

func UploadFile(c *gin.Context) {
	f, err := c.FormFile("file")
	if err != nil {
		resp.AbortErr(c, err.Error())
		return
	}

	fm := model.File{
		OriginalFilename: f.Filename,
	}

	if err := global.DB().Create(&fm).Error; err != nil {
		resp.AbortErr(c, err.Error())
		return
	}

	src, err := f.Open()
	if err != nil {
		resp.AbortErr(c, err.Error())
		return
	}
	defer src.Close()

	hash := md5.New()

	dst, err := os.Create(fmt.Sprintf("/xxx/%d", fm.ID))
	if err != nil {
		resp.AbortErr(c, err.Error())
		return
	}
	defer dst.Close()

	writer := io.NewMultiWriter(hash, dst)

	_, err = sio.Copy(writer, src)
	if err != nil {
		resp.AbortErr(c, err.Error())
		return
	}

	md5 := hash.Sum(nil)
	if err := global.DB().Model(&model.File{}).Updates(map[string]interface{}{
		"md5": md5,
	}).Error; err != nil {
		resp.AbortErr(c, err.Error())
		return
	}

	resp.Ok(c, gin.H{
		"id": fm.ID,
		"created_at": fm.CreatedAt,
	})
}

func Files(c *gin.Context) {
	uid, ok := c.Params.Get("uid")
	if !ok {
		resp.AbortErr(c, "not uid")
		return
	}
	fs := []*model.File{}
	if err := global.DB().Where("uid = ?", uid).Find(&fs).Error; err != nil {
		resp.AbortErr(c, err.Error())
		return
	}

	resp.Ok(c, gin.H{
		"files": fs,
	})
}

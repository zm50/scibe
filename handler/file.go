package handler

import (
	"crypto/md5"
	"errors"
	"fmt"
	sio "io"
	"os"
	"path/filepath"
	"scibe/global"
	"scibe/model"
	"scibe/utils/io"
	"scibe/utils/ret"
	"scibe/utils/types"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"github.com/gen2brain/go-fitz"
)

func UploadFile(c *gin.Context) {
	prop := types.GetProperty(c)
	f, err := c.FormFile("file")
	if err != nil {
		log.Err(err).Msg("failed to get file")
		ret.AbortErr(c, err.Error())
		return
	}

	ext := filepath.Ext(f.Filename)
	switch ext {
	case ".pdf":
	default:
		ret.AbortErr(c, "unsupported file type")
		return
	}

	filename := f.Filename[:len(f.Filename)-len(ext)]
	if len(ext) > 0 {
		ext = ext[1:] // remove dot
	}

	df := model.DocFile{
		Uid:  prop.Uid,
		Name: filename,
		Extension: ext,
	}

	result := global.DB().Where("uid = ? AND name = ? AND extension = ?", prop.Uid, filename, ext).FirstOrCreate(&df)
	if result.Error != nil {
		log.Err(err).Msg("failed to get file record")
		ret.AbortErr(c, err.Error())
		return
	}

	if result.RowsAffected == 0 {
		ret.AbortErr(c, "file already exists")
		return
	}

	src, err := f.Open()
	if err != nil {
		log.Err(err).Msg("failed to open file")
		ret.AbortErr(c, err.Error())
		return
	}
	defer src.Close()

	hash := md5.New()
	dst, err := os.Create(filePath(df.ID))
	if err != nil {
		log.Err(err).Msg("failed to create file")
		ret.AbortErr(c, err.Error())
		return
	}
	defer dst.Close()

	writer := io.NewMultiWriter(hash, dst)

	_, err = sio.Copy(writer, src)
	if err != nil {
		log.Err(err).Msg("failed to copy file")
		ret.AbortErr(c, err.Error())
		return
	}

	md5 := hash.Sum(nil)
	if err := global.DB().Model(&model.DocFile{}).Where("id = ?", df.ID).Updates(map[string]interface{}{
		"md5": md5,
	}).Error; err != nil {
		log.Err(err).Msg("failed to update file record")
		ret.AbortErr(c, err.Error())
		return
	}

	ret.Ok(c, gin.H{
		"id": df.ID,
		"created_at": df.CreatedAt,
	})
}

type FilesResponseItem struct {
	ID uint `json:"id"`
	FileName string `json:"file_name"`
	CreatedAt string `json:"created_at"`
}

func Files(c *gin.Context) {
	prop := types.GetProperty(c)

	fs := []*model.DocFile{}
	if err := global.DB().Where("uid = ?", prop.Uid).Find(&fs).Error; err != nil {
		log.Err(err).Msg("failed to get files")
		fmt.Println("failed to get files", err.Error())
		ret.AbortErr(c, err.Error())
		return
	}

	respItems := make([]*FilesResponseItem, 0, len(fs))

	for _, f := range fs {
		respItems = append(respItems, &FilesResponseItem{
			ID: f.ID,
			FileName: f.Name,
			CreatedAt: f.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	ret.Ok(c, gin.H{
		"files": respItems,
	})
}

type FileExtractTextRequest struct {
	Fid uint `json:"fid"`
}

type FileExtractTextResponse struct {
    Text string `json:"text"`
}

func FileExtractText(c *gin.Context) {
	prop := types.GetProperty(c)
	req := &FileExtractTextRequest{}
	if err := c.ShouldBindJSON(req); err != nil {
		log.Err(err).Msg("failed to bind request")
		ret.AbortErr(c, err.Error())
		return
	}

	docFile := &model.DocFile{}
	err := global.DB().Where("id = ? AND uid = ?", req.Fid, prop.Uid).First(docFile).Error
	if err != nil {
		log.Err(err).Msg("failed to get file")
		ret.AbortErr(c, err.Error())
		return
	}

	text, err := fileExtractText(filePath(docFile.ID))
	if err != nil {
		log.Err(err).Msg("failed to extract text")
		ret.AbortErr(c, err.Error())
		return
	}

	res := FileExtractTextResponse{
		Text: text,
	}

	ret.Ok(c, res)
}

func filePath(fid uint) string {
	return fmt.Sprintf("./file/%d", fid)
}

func fileExtractText(fp string) (string, error) {
	file, err := os.Open(fp)
	if err != nil {
		log.Err(err).Msg("failed to open file")
		return "", err
	}
	defer file.Close()

	doc, err := fitz.New(fp)
	if err != nil {
		log.Err(err).Msg("failed to open pdf")
		return "", err
	}
	defer doc.Close()

	if doc.NumPage() == 0 {
		log.Error().Msg("no pages in pdf")
		return "", errors.New("no pages in pdf")
	}

	text, err := doc.Text(0)
	if err != nil {
		log.Err(err).Msg("failed to extract text")
		return "", err
	}

	return text, nil
}

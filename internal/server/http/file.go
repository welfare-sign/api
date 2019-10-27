package http

import (
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"welfare-sign/internal/apicode"
)

// uploadFile 上传文件
// @Summary 上传文件
// @Description upload file
// @Tags 文件
// @Accept mpfd
// @Produce json
// @Param file formData file true "upload file"
// @Param type query string true "file type" default("avatar", "poster")
// @Success 200 {object} http.baseResponse	"{"status":true}"
// @Router /files/upload [post]
func uploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusOK, baseResponse{
			Status:  false,
			Code:    apicode.ErrInvalidParame,
			Message: apicode.MapZH[apicode.ErrInvalidParame],
			Error:   err.Error(),
		})
		return
	}
	spec := c.Query("type")
	if spec != "avatar" && spec != "poster" {
		c.JSON(http.StatusOK, baseResponse{
			Status:  false,
			Code:    apicode.ErrInvalidParame,
			Message: apicode.MapZH[apicode.ErrInvalidParame],
			Error:   "invalid type",
		})
		return
	}
	filePrefix := "upload/" + spec + "/"
	filename := strconv.FormatInt(time.Now().Unix(), 10) + ".png"
	if err := c.SaveUploadedFile(file, filePrefix+filename); err != nil {
		c.JSON(http.StatusOK, baseResponse{
			Status:  false,
			Code:    apicode.ErrUploadFile,
			Message: apicode.MapZH[apicode.ErrUploadFile],
			Error:   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, baseResponse{
		Status:  true,
		Code:    apicode.Success,
		Message: apicode.MapZH[apicode.Success],
		Data:    filename,
	})
}

// downloadFile 下载文件
// @Summary 下载文件
// @Description download file
// @Tags 文件
// @Accept json
// @Produce json
// @Param filename query string true "filename"
// @Param type query string true "file type" default("avatar", "poster")
// @Success 200 {object} http.baseResponse	"{"status":true}"
// @Router /files/download [get]
func downloadFile(c *gin.Context) {
	filename := c.Query("filename")
	spec := c.Query("type")
	if spec != "avatar" && spec != "poster" {
		c.JSON(http.StatusOK, baseResponse{
			Status:  false,
			Code:    apicode.ErrInvalidParame,
			Message: apicode.MapZH[apicode.ErrInvalidParame],
			Error:   "invalid type",
		})
		return
	}
	filepath := "upload/" + spec + "/" + filename
	file, err := os.Open(filepath)
	if err != nil {
		c.JSON(http.StatusOK, baseResponse{
			Status:  false,
			Code:    apicode.ErrDownloadFile,
			Message: apicode.MapZH[apicode.ErrDownloadFile],
			Error:   err.Error(),
		})
		return
	}
	defer file.Close()
	_, err = io.Copy(c.Writer, file)
	if err != nil {
		c.JSON(http.StatusOK, baseResponse{
			Status:  false,
			Code:    apicode.ErrDownloadFile,
			Message: apicode.MapZH[apicode.ErrDownloadFile],
			Error:   err.Error(),
		})
	}
}

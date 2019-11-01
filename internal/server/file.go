package server

import (
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"welfare-sign/internal/apicode"
	"welfare-sign/internal/pkg/wsgin"
)

// uploadFile 上传文件
// @Summary 上传文件
// @Description upload file
// @Tags 文件
// @Accept mpfd
// @Produce json
// @Param file formData file true "upload file"
// @Param type query string true "file type" default("avatar", "poster")
// @Success 200 {object} server.BaseResponse	"{"status":true}"
// @Router /files/upload [post]
func uploadFile(c *gin.Context) {

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusOK, BaseResponse{
			Status:  false,
			Code:    wsgin.APICodeInvalidParame,
			Message: wsgin.APICodeMapZH[wsgin.APICodeInvalidParame],
			Error:   err.Error(),
		})
		return
	}
	spec := c.Query("type")
	if spec != "avatar" && spec != "poster" {
		c.JSON(http.StatusOK, BaseResponse{
			Status:  false,
			Code:    wsgin.APICodeInvalidParame,
			Message: wsgin.APICodeMapZH[wsgin.APICodeInvalidParame],
			Error:   "invalid type",
		})
		return
	}
	filePrefix := "upload/" + spec + "/"
	filename := strconv.FormatInt(time.Now().Unix(), 10) + ".png"
	if err := c.SaveUploadedFile(file, filePrefix+filename); err != nil {
		c.JSON(http.StatusOK, BaseResponse{
			Status:  false,
			Code:    apicode.ErrUploadFile,
			Message: wsgin.APICodeMapZH[apicode.ErrUploadFile],
			Error:   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, BaseResponse{
		Status:  true,
		Code:    wsgin.APICodeSuccess,
		Message: wsgin.APICodeMapZH[wsgin.APICodeSuccess],
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
// @Success 200 {object} server.BaseResponse	"{"status":true}"
// @Router /files/download [get]
func downloadFile(c *gin.Context) {
	filename := c.Query("filename")
	spec := c.Query("type")
	if spec != "avatar" && spec != "poster" {
		c.JSON(http.StatusOK, BaseResponse{
			Status:  false,
			Code:    wsgin.APICodeInvalidParame,
			Message: wsgin.APICodeMapZH[wsgin.APICodeInvalidParame],
			Error:   "invalid type",
		})
		return
	}
	filepath := "upload/" + spec + "/" + filename
	file, err := os.Open(filepath)
	if err != nil {
		c.JSON(http.StatusOK, BaseResponse{
			Status:  false,
			Code:    apicode.ErrDownloadFile,
			Message: wsgin.APICodeMapZH[apicode.ErrDownloadFile],
			Error:   err.Error(),
		})
		return
	}
	defer file.Close()
	_, err = io.Copy(c.Writer, file)
	if err != nil {
		c.JSON(http.StatusOK, BaseResponse{
			Status:  false,
			Code:    apicode.ErrDownloadFile,
			Message: wsgin.APICodeMapZH[apicode.ErrDownloadFile],
			Error:   err.Error(),
		})
	}
}

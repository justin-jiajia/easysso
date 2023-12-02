package handler

import (
	"fmt"
	"log"
	"net/http"
	"path"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/justin-jiajia/easysso/api/config"
	"github.com/justin-jiajia/easysso/api/database"
	"github.com/justin-jiajia/easysso/api/middleware"
)

func AvatarUploadHandler(ctx *gin.Context) {
	allow_type := [...]string{".jpg", ".jpeg", ".gif", ".png"}
	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "上传失败，错误：" + err.Error()})
	} else {
		flag := false
		this_type := ""
		for _, noww := range allow_type {
			if strings.HasSuffix(file.Filename, noww) {
				flag = true
				this_type = noww
				break
			}
		}
		if !flag {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "上传失败，文件后缀应为 jpg jpeg gif png 中的一种"})
		} else {
			filename := fmt.Sprintf("%d", middleware.ID) + this_type
			err = ctx.SaveUploadedFile(file, path.Join(config.Config.AvatarSavePath, filename))
			if err != nil {
				log.Panicf("Error saving avatar. Error: %v", err.Error())
			}
			ctx.String(http.StatusNoContent, "\n")
			now_user := &database.User{}
			database.DB.First(&now_user, middleware.ID)
			now_user.AvatarFileName = filename
			database.DB.Save(&now_user)
		}
	}
}

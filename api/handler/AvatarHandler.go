package handler

import (
	"net/http"
	"path"
	"strconv"
	"strings"

	"github.com/justin-jiajia/easysso/api/config"
	"github.com/justin-jiajia/easysso/api/database"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

func AvatarHandler(ctx *gin.Context) {
	this_id_string := strings.Replace(ctx.Param("id"), "/", "", 2)
	this_id, err := strconv.Atoi(this_id_string)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID should be a number"})
		return
	}
	now_user := &database.User{}
	result := database.DB.First(&now_user, this_id)
	if result.Error == gorm.ErrRecordNotFound {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	} else if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "出错了：" + result.Error.Error()})
		return
	} else {
		filepath := path.Join(config.Config.AvatarSavePath, now_user.AvatarFileName)
		ctx.File(filepath)
	}
}

package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/justin-jiajia/easysso/api/database"
	"github.com/justin-jiajia/easysso/api/middleware"
	"github.com/justin-jiajia/easysso/api/utils"
)

type ChangePasswd struct {
	NewPasswd string `json:"new_passwd" binding:"required"`
	OldPasswd string `json:"old_passwd" binding:"required"`
}

func ChangePasswdHandler(ctx *gin.Context) {
	now_user := &database.User{}
	database.DB.First(&now_user, middleware.ID)
	var json ChangePasswd
	err := ctx.ShouldBindJSON(&json)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "出错了: " + err.Error()})
		return
	}
	if !utils.VerifyPasswd(json.OldPasswd, now_user.PasswordHash) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "原密码错误"})
		return
	}
	now_user.PasswordHash = utils.GetPasswdHash(json.NewPasswd)
	ctx.String(http.StatusNoContent, "\n")
	database.DB.Save(&now_user)
}

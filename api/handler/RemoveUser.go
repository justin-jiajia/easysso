package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/justin-jiajia/easysso/api/database"
	"github.com/justin-jiajia/easysso/api/middleware"
	"github.com/justin-jiajia/easysso/api/utils"
)

type RemoveUser struct {
	Password string `json:"passwd" binding:"required"`
}

func RemoveUserHandler(ctx *gin.Context) {
	now_user := &database.User{}
	database.DB.First(&now_user, middleware.ID)
	var json RemoveUser
	err := ctx.ShouldBindJSON(&json)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "出错了: " + err.Error()})
		return
	}
	if !utils.VerifyPasswd(json.Password, now_user.PasswordHash) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "密码错误"})
		return
	}
	ctx.String(http.StatusNoContent, "\n")
	database.DB.Delete(&now_user)
}

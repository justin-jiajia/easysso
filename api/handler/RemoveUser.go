package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/justin-jiajia/easysso/api/database"
	"github.com/justin-jiajia/easysso/api/middleware"
)

type RemoveUser struct {
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
	ctx.String(http.StatusNoContent, "\n")
	database.DB.Select("Credentials").Delete(&now_user)
}

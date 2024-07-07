package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/justin-jiajia/easysso/api/database"
	"github.com/justin-jiajia/easysso/api/utils"
)

type Information struct {
	ClientID     string `json:"client_id" binding:"required"`
	ClientSecret string `json:"client_secret" binding:"required"`
	Token        string `json:"token" binding:"required"`
}

func InformationHandler(ctx *gin.Context) {
	var json Information
	err := ctx.ShouldBindJSON(&json)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "出错了: " + err.Error()})
		return
	}
	uid, err := utils.VerifyOath2Token(json.Token, json.ClientID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	now_user := &database.User{}
	database.DB.First(&now_user, uid)
	ctx.JSON(http.StatusOK, gin.H{"id": uid, "username": now_user.UserName, "jointime": now_user.CreatedAt.Unix()})
}

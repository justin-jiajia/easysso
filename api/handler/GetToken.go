package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/justin-jiajia/easysso/api/config"
	"github.com/justin-jiajia/easysso/api/database"
	"gorm.io/gorm"
)

type GetToken struct {
	ClientSecret string `json:"client_secret" binding:"required"`
	ClientID     string `json:"client_id" binding:"required"`
	Code         string `json:"code" binding:"required"`
}

func GetTokenHandler(ctx *gin.Context) {
	var json GetToken
	err := ctx.ShouldBindJSON(&json)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "出错了: " + err.Error()})
		return
	}
	code := database.ServerCode2Token{}
	res := database.DB.First(&code, "code = ? AND client_id = ?", json.Code, json.ClientID)
	if res.Error == gorm.ErrRecordNotFound {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "错误的client_id,client_secret或code"})
		return
	} else if res.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "出错了:" + res.Error.Error()})
		return
	}
	if config.ClientsMap[json.ClientID].ClientSecret != json.ClientSecret {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "错误的client_id,client_secret或code"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"token": code.Token})
	database.DB.Delete(code)
}

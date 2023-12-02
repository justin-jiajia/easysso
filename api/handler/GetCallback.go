package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/justin-jiajia/easysso/api/config"
)

type VerifyServer struct {
	ClientID string `json:"client_id" binding:"required"`
}

func GetCallbackHandler(ctx *gin.Context) {
	var json VerifyServer
	err := ctx.ShouldBindJSON(&json)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "出错了: " + err.Error()})
		return
	}
	value, ok := config.ClientsMap[json.ClientID]
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "不正确的ClientID"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"callback": value.Callback, "name": value.Name})
}

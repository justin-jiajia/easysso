package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/justin-jiajia/easysso/api/config"
	"github.com/justin-jiajia/easysso/api/database"
	"github.com/justin-jiajia/easysso/api/middleware"
	"github.com/justin-jiajia/easysso/api/utils"
)

type GetCode struct {
	ClientID string `json:"client_id" binding:"required"`
}

func GetCodeHandler(ctx *gin.Context) {
	var json GetCode
	err := ctx.ShouldBindJSON(&json)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "出错了: " + err.Error()})
		return
	}
	c, ok := config.ClientsMap[json.ClientID]
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "不正确的ClientID"})
		return
	}
	code := uuid.New().String()
	nw := &database.ServerCode2Token{}
	nw.ClientID = json.ClientID
	nw.Exp = time.Now().Add(time.Minute * 5)
	nw.Token = utils.NewOAuth2Token(middleware.ID, json.ClientID, middleware.ExpTime)
	utils.NewUserLog(middleware.ID, ctx.Request.UserAgent(), ctx.ClientIP(), "登录"+c.Name)
	nw.Code = code
	database.DB.Create(nw)
	ctx.JSON(http.StatusOK, gin.H{"code": code})
}

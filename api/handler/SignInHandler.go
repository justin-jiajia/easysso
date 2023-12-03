package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/justin-jiajia/easysso/api/database"
	"github.com/justin-jiajia/easysso/api/utils"
	"gorm.io/gorm"
)

type SignIn struct {
	Username string `json:"username" binding:"required"`
	Passward string `json:"password" binding:"required"`
}

func SignInHandler(ctx *gin.Context) {
	var json SignIn
	err := ctx.ShouldBindJSON(&json)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "出错了: " + err.Error()})
		return
	}
	nowuser := database.User{}
	result := database.DB.Where(database.User{UserName: json.Username}).First(&nowuser)
	if result.Error == nil {
		if utils.VerifyPasswd(json.Passward, nowuser.PasswordHash) {
			token, exp := utils.NewUserToken(nowuser.ID)
			ctx.JSON(http.StatusOK, gin.H{"token": token, "id": nowuser.ID, "expire": exp.Unix()})
			return
		} else {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "密码不正确"})
			return
		}
	} else {
		if result.Error == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "用户不存在"})
			return
		} else {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "出错了: " + result.Error.Error()})
			return
		}
	}
}

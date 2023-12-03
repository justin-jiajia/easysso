package handler

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/justin-jiajia/easysso/api/database"
	"github.com/justin-jiajia/easysso/api/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SignUp struct {
	Username string `json:"username" binding:"required"`
	Passward string `json:"password" binding:"required"`
}

func SignUpHandler(ctx *gin.Context) {
	var json SignUp
	err := ctx.ShouldBindJSON(&json)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "出错了: " + err.Error()})
		return
	}
	test_user := database.User{}
	result := database.DB.Where(&database.User{UserName: json.Username}).First(&test_user)
	if result.Error == gorm.ErrRecordNotFound {
		new_user := database.User{UserName: json.Username, PasswordHash: utils.GetPasswdHash(json.Passward), UUID: uuid.New().String()}
		database.DB.Create(&new_user)
		token, exp := utils.NewUserToken(new_user.ID)
		ctx.JSON(http.StatusOK, gin.H{"token": token, "id": new_user.ID, "expire": exp.Unix()})
	} else {
		if result.Error == nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "用户名已存在，换个试试吧"})
		} else {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "出错了: " + result.Error.Error()})
		}
	}
}

package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/justin-jiajia/easysso/api/database"
	"github.com/justin-jiajia/easysso/api/middleware"
)

func RemoveTokenHandler(c *gin.Context) {
	nwtoken := database.UserToken{}
	database.DB.Model(&database.UserToken{}).Where(&database.UserToken{Token: middleware.Token}).First(&nwtoken)
	database.DB.Delete(&nwtoken)
	for _, v := range nwtoken.UserServerTokens {
		database.DB.Delete(&v)
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

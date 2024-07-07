package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/justin-jiajia/easysso/api/utils"
)

var ID uint
var ExpTime time.Time
var Token string

func Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		t, ok := ctx.Request.Header["Authorization"]
		if !ok {
			ctx.Abort()
			ctx.JSON(http.StatusForbidden, gin.H{"error": "Auth Required"})
		} else {
			token := t[0]
			if len(token) <= 7 || token[:7] != "Bearer " {
				ctx.Abort()
				ctx.JSON(http.StatusForbidden, gin.H{"error": "Invalid Token format"})
			}
			id, exp, err := utils.VerifyUserToken(token[7:])
			if err != nil {
				ctx.Abort()
				ctx.JSON(http.StatusForbidden, gin.H{"error": "Invalid Token. Error: " + err.Error()})
			}
			ID = id
			ExpTime = exp
			Token = token
		}
	}
}

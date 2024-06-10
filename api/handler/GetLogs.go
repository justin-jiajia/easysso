package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/justin-jiajia/easysso/api/database"
	"github.com/justin-jiajia/easysso/api/middleware"
)

type LogRes struct {
	Action     string `json:"action"`
	IP         string `json:"ip"`
	UserAgent  string `json:"useragent"`
	ActionTime int64  `json:"actiontime"`
}

func GetActions(c *gin.Context) {
	uid := middleware.ID
	user := &database.User{}
	res := database.DB.Preload("UserLogs").Where(database.User{ID: uid}).First(&user)
	if res.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "出错了: " + res.Error.Error()})
		return
	}
	rres := make([]LogRes, 0)
	for _, v := range user.UserLogs {
		rres = append(rres, LogRes{
			Action:     v.Action,
			IP:         v.IP,
			UserAgent:  v.UserAgent,
			ActionTime: v.ActionTime.UTC().Unix(),
		})
	}
	c.JSON(http.StatusOK, rres)
}

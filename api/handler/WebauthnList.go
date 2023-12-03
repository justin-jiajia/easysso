package handler

import (
	"encoding/base64"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/justin-jiajia/easysso/api/database"
	"github.com/justin-jiajia/easysso/api/middleware"
)

type CredentialRes struct {
	Name         string `json:"name"`
	LastUsed     int64  `json:"last_used"`
	Created      int64  `json:"created"`
	ID           string `json:"id"`
	UsernameLess bool   `json:"username_less"`
}

func WebauthnList(c *gin.Context) {
	uid := middleware.ID
	user := &database.User{}
	database.DB.Where(database.User{ID: uid}).First(&user)
	var res []CredentialRes
	var crs []database.Credential
	err := database.DB.Model(&user).Association("Credentials").Find(&crs)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "出错了: " + err.Error()})
		return
	}
	for _, v := range crs {
		res = append(res, CredentialRes{
			Name:         v.Name,
			LastUsed:     v.UpdatedAt.UTC().Unix(),
			Created:      v.CreatedAt.UTC().Unix(),
			ID:           base64.RawURLEncoding.EncodeToString(v.ID),
			UsernameLess: v.UsernameLess,
		})
	}
	c.JSON(http.StatusOK, res)
}

type EditQuery struct {
	Base64edID string `json:"id"`
	NewName    string `json:"new_name"`
}

func WebauthnEdit(c *gin.Context) {
	var json EditQuery
	err := c.ShouldBindJSON(&json)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "出错了: " + err.Error()})
		return
	}
	decodedID, err := base64.RawURLEncoding.DecodeString(json.Base64edID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "出错了: " + err.Error()})
		return
	}
	var cred database.Credential
	database.DB.Where(database.Credential{ID: decodedID}).First(&cred)
	cred.Name = json.NewName
	database.DB.Save(&cred)
	c.String(http.StatusNoContent, "\n")
}

type DeleteQuery struct {
	Base64edID string `json:"id"`
}

func WebauthnDelete(c *gin.Context) {
	var json DeleteQuery
	err := c.ShouldBindJSON(&json)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "出错了: " + err.Error()})
		return
	}

	decodedID, err := base64.RawURLEncoding.DecodeString(json.Base64edID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "出错了: " + err.Error()})
		return
	}

	var cred database.Credential
	database.DB.Where(database.Credential{ID: decodedID}).First(&cred)

	database.DB.Delete(&cred)
	c.String(http.StatusNoContent, "\n")
}

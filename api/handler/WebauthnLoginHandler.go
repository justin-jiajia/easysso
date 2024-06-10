package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-webauthn/webauthn/protocol"
	w "github.com/go-webauthn/webauthn/webauthn"
	"github.com/justin-jiajia/easysso/api/config"
	"github.com/justin-jiajia/easysso/api/database"
	"github.com/justin-jiajia/easysso/api/utils"
	"github.com/justin-jiajia/easysso/api/webauthn"
	"gorm.io/gorm"
)

type WSignIn struct {
	Username string `json:"username"`
}

func WebauthnLoginStartHandler(ctx *gin.Context) {
	var json WSignIn
	err := ctx.ShouldBindJSON(&json)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "出错了: " + err.Error()})
		return
	}
	nsession, errr := webauthn.Store.Get(ctx.Request, "webauthn")
	if errr != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "出错了: " + errr.Error()})
		return
	}
	var options *protocol.CredentialAssertion
	var session *w.SessionData
	if json.Username == "" {
		options, session, err = webauthn.WebAuthn.BeginDiscoverableLogin()
		nsession.Values["discoverable"] = true
	} else {
		user := database.User{}
		result := database.DB.Where(database.User{UserName: json.Username}).First(&user)
		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": "用户不存在"})
			} else {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": "出错了: " + result.Error.Error()})
			}
			return
		}
		nsession.Values["userrid"] = user.ID
		nsession.Values["discoverable"] = false
		options, session, err = webauthn.WebAuthn.BeginLogin(&user)
	}

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "出错了: " + err.Error()})
		return
	}
	nsession.Values["challenge"] = session.Challenge
	nsession.Values["user_verification"] = session.UserVerification
	nsession.Values["userid"] = session.UserID
	nsession.Values["exp"] = session.Expires
	err = nsession.Save(ctx.Request, ctx.Writer)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "出错了: " + err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"options": options,
	})
	// options.publicKey contain our registration options
}

func WebauthnLoginFinishHandler(ctx *gin.Context) {
	// Get the session data stored from the function above
	nsession, err := webauthn.Store.Get(ctx.Request, "webauthn")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "出错了: " + err.Error()})
		return
	}
	session := w.SessionData{
		Challenge:        nsession.Values["challenge"].(string),
		UserVerification: *nsession.Values["user_verification"].(*protocol.UserVerificationRequirement),
		UserID:           nsession.Values["userid"].([]byte),
		Expires:          *nsession.Values["exp"].(*time.Time),
	}
	user := &database.User{}
	var credential *w.Credential
	if nsession.Values["discoverable"].(bool) {
		credential, err = webauthn.WebAuthn.FinishDiscoverableLogin(database.DiscoverableUserHandler, session, ctx.Request)
		tmp := &database.Credential{}
		database.DB.Model(&database.Credential{}).Where(&database.Credential{ID: credential.ID}).Find(&tmp)
		user.ID = tmp.UserID
	} else {
		database.DB.Where(database.User{ID: nsession.Values["userrid"].(uint)}).First(&user)
		credential, err = webauthn.WebAuthn.FinishLogin(user, session, ctx.Request)
	}
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "出错了: " + err.Error()})
		return
	}
	if credential.Authenticator.CloneWarning {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "CloneWarning"})
		return
	}
	c := database.Credential{
		ID:              credential.ID,
		PublicKey:       credential.PublicKey,
		AttestationType: credential.AttestationType,
		Transport:       credential.Transport,
		Flags:           database.CredentialFlags(credential.Flags),
		Authenticator:   database.Authenticator(credential.Authenticator),
		UserID:          user.ID,
		UsernameLess:    nsession.Values["discoverable"].(bool),
	}
	database.DB.Model(&database.Credential{}).Where(&database.Credential{ID: c.ID}).Updates(c)

	token, exp := utils.NewUserToken(user.ID)
	ctx.SetCookie("webauthn", "", -1, "/", config.Config.RPID, false, true)
	utils.NewUserLog(user.ID, ctx.GetHeader("User-Agent"), ctx.ClientIP(), "修改密码")
	ctx.JSON(http.StatusOK, gin.H{"token": token, "id": user.ID, "expire": exp.Unix()})
}

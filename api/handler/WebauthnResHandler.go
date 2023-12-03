package handler

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-webauthn/webauthn/protocol"
	w "github.com/go-webauthn/webauthn/webauthn"
	"github.com/justin-jiajia/easysso/api/config"
	"github.com/justin-jiajia/easysso/api/database"
	"github.com/justin-jiajia/easysso/api/middleware"
	"github.com/justin-jiajia/easysso/api/webauthn"
)

func WebauthnResStartHandler(ctx *gin.Context) {
	uid := middleware.ID
	user := &database.User{}
	database.DB.Where(database.User{ID: uid}).First(&user)
	residentkeyrequirement := protocol.ResidentKeyRequirementDiscouraged
	if ctx.Query("discovered") != "" {
		residentkeyrequirement = protocol.ResidentKeyRequirementRequired
	}
	options, session, err := webauthn.WebAuthn.BeginRegistration(user, w.WithResidentKeyRequirement(residentkeyrequirement))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "出错了: " + err.Error()})
		return
	}
	// store the sessionData values to the client's cookie
	nsession, errr := webauthn.Store.Get(ctx.Request, "webauthn")
	if errr != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "出错了: " + errr.Error()})
		return
	}
	nsession.Values["discovered"] = ctx.Query("discovered") != ""
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

func WebauthnResFinishHandler(ctx *gin.Context) {
	uid := middleware.ID
	user := &database.User{}
	database.DB.Where(database.User{ID: uid}).First(&user)

	// Get the session data stored from the function above
	nsession, err := webauthn.Store.Get(ctx.Request, "webauthn")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "出错了: " + err.Error()})
		return
	}
	discovered := nsession.Values["discovered"].(bool)
	session := w.SessionData{
		Challenge:        nsession.Values["challenge"].(string),
		UserVerification: *nsession.Values["user_verification"].(*protocol.UserVerificationRequirement),
		UserID:           nsession.Values["userid"].([]byte),
		Expires:          *nsession.Values["exp"].(*time.Time),
	}

	credential, err := webauthn.WebAuthn.FinishRegistration(user, session, ctx.Request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "出错了: " + err.Error()})
		return
	}

	// If creation was successful, store the credential object
	// Pseudocode to add the user credential.
	c := database.Credential{
		ID:              credential.ID,
		PublicKey:       credential.PublicKey,
		AttestationType: credential.AttestationType,
		Transport:       credential.Transport,
		Flags:           database.CredentialFlags(credential.Flags),
		Authenticator:   database.Authenticator(credential.Authenticator),
		UserID:          user.ID,
		UsernameLess:    discovered,
	}
	log.Println(c)
	database.DB.Save(&c)
	ctx.SetCookie("webauthn", "", -1, "/", config.Config.RPID, false, true)
	ctx.String(http.StatusNoContent, "\n")
}

package utils

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/justin-jiajia/easysso/api/config"
	"github.com/justin-jiajia/easysso/api/database"
	"golang.org/x/crypto/bcrypt"
)

func GetPasswdHash(passwd string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(passwd), bcrypt.DefaultCost)
	if err != nil {
		fmt.Print(err)
		panic("Error getting password hash")
	}
	return string(hash)
}

func VerifyPasswd(passwd string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(passwd))
	return err == nil
}

func NewOAuth2Token(userid uint, clientid string, exp time.Time) string {
	s := uuid.New().String()
	database.DB.Create(&database.ServerToken{Token: s, Exp: exp, ClientID: clientid, UserID: userid})
	return s
}

func NewUserToken(userid uint) (string, time.Time) {
	s := uuid.New().String()
	timeexp := time.Now().Add(time.Duration(config.Config.TokenExpTime * int64(time.Second)))
	database.DB.Create(&database.UserToken{Token: s, Exp: timeexp, UserID: userid})
	return "Bearer " + s, timeexp
}

func NewUserLog(userid uint, useragent string, ip string, action string) {
	database.DB.Create(&database.UserLog{UserID: userid, UserAgent: useragent, IP: ip, ActionTime: time.Now(), Action: action})
}

func VerifyUserToken(token_string string) (uint, time.Time, error) {
	token := &database.UserToken{}
	database.DB.Model(&database.UserToken{}).Where(&database.UserToken{Token: token_string}).First(&token)
	if token.Exp.Unix() < time.Now().Unix() {
		database.DB.Delete(&token)
		return 0, time.Now(), errors.New("token expired")
	}
	return token.UserID, token.Exp, nil
}

func VerifyOath2Token(token_string string, clientid string) (uint, error) {
	token := &database.ServerToken{}
	database.DB.Model(&database.ServerToken{}).Where(&database.ServerToken{Token: token_string, ClientID: clientid}).First(&token)
	if token.Exp.Unix() < time.Now().Unix() {
		database.DB.Delete(&token)
		return 0, errors.New("token expired")
	}
	return token.UserID, nil
}

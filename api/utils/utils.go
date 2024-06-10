package utils

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
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
	t := jwt.NewWithClaims(jwt.SigningMethodHS512,
		jwt.RegisteredClaims{
			Issuer:    config.Config.TokenName,
			Subject:   fmt.Sprint(userid),
			ExpiresAt: jwt.NewNumericDate(exp),
			Audience:  jwt.ClaimStrings{clientid},
		})
	s, err := t.SignedString(config.Config.TokenKeyByte)
	if err != nil {
		fmt.Println(err)
		panic("Error sign token")
	}
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
	token, err := jwt.Parse(token_string,
		func(t *jwt.Token) (interface{}, error) {
			return config.Config.TokenKeyByte, nil
		},
		jwt.WithValidMethods([]string{"HS512"}),
		jwt.WithIssuer(config.Config.TokenName),
	)
	if err != nil {
		return 0, err
	}
	sub, err2 := token.Claims.GetSubject()
	if err2 != nil {
		return 0, err2
	}
	rid, err3 := strconv.Atoi(sub)
	if err3 != nil {
		return 0, err3
	}
	exp, err4 := token.Claims.GetExpirationTime()
	if err4 != nil {
		return 0, err4
	}
	if exp == nil {
		return 0, errors.New("no exp")
	}
	if int64(exp.Unix()) < time.Now().Unix() {
		return 0, errors.New("token expired")
	}
	aud, err5 := token.Claims.GetAudience()
	if err5 != nil {
		return 0, err5
	}
	audstr := aud[0]
	if string(audstr) != clientid {
		return 0, errors.New("clientid mismatch")
	}
	id := uint(rid)
	now_user := &database.User{}
	res := database.DB.First(&now_user, id)
	if res.Error != nil {
		return 0, res.Error
	}
	return id, nil
}

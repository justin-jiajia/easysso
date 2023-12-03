package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/justin-jiajia/easysso/api/config"
	"github.com/justin-jiajia/easysso/api/database"
	"github.com/justin-jiajia/easysso/api/router"
	"github.com/justin-jiajia/easysso/api/webauthn"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func SetupRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	config.ReadConfig()
	webauthn.InitSessionStore()
	webauthn.InitWebauthn()
	database.InitDB()
	database.Migrate()
	router.InitApi(r)
	return r
}

func GetUUID() string {
	return uuid.New().String()
}

func MapToJson(x map[string]interface{}) []byte {
	response, err := json.Marshal(x)
	if err != nil {
		panic("error!" + err.Error())
	}
	return response
}

func PerformPostRequest(r http.Handler, path, header string, body map[string]interface{}) *httptest.ResponseRecorder {
	req, _ := http.NewRequest("POST", path, bytes.NewBuffer(MapToJson(body)))
	req.Header.Set("Authorization", header)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func JSONToMap(w httptest.ResponseRecorder) map[string]interface{} {
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		panic("error!" + err.Error() + " | it were" + w.Body.String())
	}
	return response
}

func TestSignUpAndSignIn(t *testing.T) {
	r := SetupRouter()
	username := GetUUID()
	passwd := GetUUID()
	// Sign UP
	w := PerformPostRequest(r, "/api/user/sign_up/", "", gin.H{"username": username, "password": passwd})
	a := JSONToMap(*w)
	assert.Equal(t, w.Code, 200)
	assert.NotNil(t, a["token"])
	assert.NotNil(t, a["id"])
	id := a["id"]
	assert.NotNil(t, a["expire"])
	//Sign IN
	w = PerformPostRequest(r, "/api/user/sign_in/", "", gin.H{"username": username, "password": passwd})
	a = JSONToMap(*w)
	t.Log(a)
	assert.Equal(t, w.Code, 200)
	assert.NotNil(t, a["token"])
	token := a["token"].(string)
	assert.NotNil(t, a["id"])
	assert.Equal(t, a["id"], id)
	assert.NotNil(t, a["expire"])
	//Change Passwd
	passwd = uuid.New().String()
	w = PerformPostRequest(r, "/api/user/settings/change_passwd/", token, gin.H{"new_passwd": passwd})
	assert.Equal(t, w.Code, 204)
	//Sign IN Again
	w = PerformPostRequest(r, "/api/user/sign_in/", "", gin.H{"username": username, "password": passwd})
	a = JSONToMap(*w)
	assert.Equal(t, w.Code, 200)
	assert.NotNil(t, a["token"])
	assert.NotNil(t, a["id"])
	assert.Equal(t, a["id"], id)
	assert.NotNil(t, a["expire"])
	//Test OATH2
	//Callback
	w = PerformPostRequest(r, "/api/oath2/getcallback/", "", gin.H{"client_id": "test"})
	a = JSONToMap(*w)
	assert.Equal(t, w.Code, 200)
	assert.Equal(t, a["callback"], "http://example.com/")
	assert.Equal(t, a["name"], "test")
	//Get Code
	w = PerformPostRequest(r, "/api/user/settings/getcode/", token, gin.H{"client_id": "test"})
	a = JSONToMap(*w)
	assert.Equal(t, w.Code, 200)
	assert.NotNil(t, a["code"])
	code := a["code"]
	//Get Token
	w = PerformPostRequest(r, "/api/oath2/gettoken/", "", gin.H{"client_id": "test", "client_secret": "test", "code": code})
	a = JSONToMap(*w)
	assert.Equal(t, w.Code, 200)
	assert.NotNil(t, a["token"])
	otoken := a["token"]
	t.Log(otoken)
	//Try Again
	w = PerformPostRequest(r, "/api/oath2/gettoken/", "", gin.H{"client_id": "test", "client_secret": "test", "code": code})
	assert.Equal(t, w.Code, 400)
	//Get Inf
	w = PerformPostRequest(r, "/api/oath2/information/", "", gin.H{"client_id": "test", "client_secret": "test", "token": otoken})
	t.Log(w.Body.String())
	a = JSONToMap(*w)
	assert.Equal(t, w.Code, 200)
	assert.Equal(t, a["username"], username)
	assert.Equal(t, a["id"], id)
	//Remove User
	w = PerformPostRequest(r, "/api/user/settings/remove/", token, gin.H{"passwd": passwd})
	assert.Equal(t, w.Code, 204)
	//Try Again
	w = PerformPostRequest(r, "/api/user/sign_in/", "", gin.H{"username": username, "password": passwd})
	assert.Equal(t, w.Code, 400)
}

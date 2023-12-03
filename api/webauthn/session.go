package webauthn

import (
	"encoding/gob"
	"time"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/gorilla/sessions"
	"github.com/justin-jiajia/easysso/api/config"
)

var (
	Store *sessions.CookieStore
)
var userVerificationRequirement protocol.UserVerificationRequirement = "TEST"

func InitSessionStore() {
	Store = sessions.NewCookieStore(config.Config.TokenKeyByte)
	gob.Register(&time.Time{})
	gob.Register(&userVerificationRequirement)
}

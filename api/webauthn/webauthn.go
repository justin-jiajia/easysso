package webauthn

import (
	"log"

	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/justin-jiajia/easysso/api/config"
)

var (
	WebAuthn *webauthn.WebAuthn
	err      error
)

func InitWebauthn() {

	wconfig := &webauthn.Config{
		Debug:         config.Config.IsDev,
		RPDisplayName: config.Config.Displayname,        // Display Name for your site
		RPID:          config.Config.RPID,               // Generally the FQDN for your site
		RPOrigins:     []string{config.Config.RPOrigin}, // The origin URLs allowed for WebAuthn requests
	}

	if WebAuthn, err = webauthn.New(wconfig); err != nil {
		log.Panicln(err)
	}
}

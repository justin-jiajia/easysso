package database

import (
	"github.com/go-webauthn/webauthn/protocol"
)

type Authenticator struct {
	AAGUID       []byte
	SignCount    uint32
	CloneWarning bool
	Attachment   protocol.AuthenticatorAttachment
}

func SelectAuthenticator(att string, rrk *bool, uv string) protocol.AuthenticatorSelection {
	return protocol.AuthenticatorSelection{
		AuthenticatorAttachment: protocol.AuthenticatorAttachment(att),
		RequireResidentKey:      rrk,
		UserVerification:        protocol.UserVerificationRequirement(uv),
	}
}

func (a *Authenticator) UpdateCounter(authDataCount uint32) {
	if authDataCount <= a.SignCount && (authDataCount != 0 || a.SignCount != 0) {
		a.CloneWarning = true
		DB.Save(a)
		return
	}
	a.SignCount = authDataCount
	DB.Save(a)
}

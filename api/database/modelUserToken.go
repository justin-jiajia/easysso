package database

import (
	"time"

	"github.com/go-webauthn/webauthn/webauthn"
)

type User struct {
	ID             uint `gorm:"primaryKey"`
	UUID           string
	UserName       string
	PasswordHash   string
	AvatarFileName string
	CreatedAt      time.Time
	Credentials    []Credential
}

func (user *User) WebAuthnID() []byte {
	return []byte(user.UUID)
}

func (user *User) WebAuthnName() string {
	return user.UserName
}

func (user *User) WebAuthnDisplayName() string {
	return user.UserName
}

func (user *User) WebAuthnIcon() string {
	//TODO: replace it!
	return "https://pics.com/avatar.png"
}

func (user *User) WebAuthnCredentials() []webauthn.Credential {
	var credentials []Credential
	DB.Model(&user).Association("Credentials").Find(&credentials)
	var RCredentails []webauthn.Credential
	for _, cred := range credentials {
		RCredentails = append(RCredentails, webauthn.Credential{
			ID:              cred.ID,
			PublicKey:       cred.PublicKey,
			AttestationType: cred.AttestationType,
			Transport:       cred.Transport,
			Flags:           webauthn.CredentialFlags(cred.Flags),
			Authenticator:   webauthn.Authenticator(cred.Authenticator),
		})
	}
	return RCredentails
}

func DiscoverableUserHandler(rawid, userHandle []byte) (webauthn.User, error) {
	useruuid := string(userHandle)
	user := &User{}
	res := DB.Where(User{UUID: useruuid}).First(user)
	return user, res.Error
}

type Token struct {
	Token    string
	Exp      time.Time
	Code     string `gorm:"primaryKey"`
	ClientID string
}

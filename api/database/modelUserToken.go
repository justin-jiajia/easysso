package database

import (
	"time"

	"github.com/go-webauthn/webauthn/webauthn"
)

type User struct {
	ID               uint `gorm:"primaryKey"`
	UUID             string
	UserName         string
	PasswordHash     string
	AvatarFileName   string
	CreatedAt        time.Time
	Credentials      []Credential
	UserTokens       []UserToken
	UserLogs         []UserLog
	UserServerTokens []ServerToken
}

type UserToken struct {
	// it belongs to User
	UserID           uint
	Exp              time.Time
	Token            string `gorm:"primaryKey"`
	UserServerTokens []ServerToken
}

type UserLog struct {
	UserID     uint
	UserAgent  string
	IP         string
	ActionTime time.Time
	Action     string
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

type ServerCode2Token struct {
	Token    string
	Exp      time.Time
	Code     string `gorm:"primaryKey"`
	ClientID string
}

type ServerToken struct {
	Token       string
	Exp         time.Time
	UserID      uint
	ClientID    string
	UserTokenID string
}

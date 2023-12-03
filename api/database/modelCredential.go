package database

import (
	"time"

	"github.com/go-webauthn/webauthn/protocol"
)

// === above from https://github.com/go-webauthn/webauthn/blob/master/webauthn/credential.go ===
// Credential contains all needed information about a WebAuthn credential for storage.
type Credential struct {
	ID              []byte `gorm:"primaryKey"`
	PublicKey       []byte
	AttestationType string
	Transport       []protocol.AuthenticatorTransport `gorm:"serializer:json"`
	Flags           CredentialFlags                   `gorm:"serializer:json"`
	Authenticator   Authenticator                     `gorm:"serializer:json"`
	UserID          uint
	Name            string
	UsernameLess    bool
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type CredentialFlags struct {
	UserPresent    bool
	UserVerified   bool
	BackupEligible bool
	BackupState    bool
}

// Descriptor converts a Credential into a protocol.CredentialDescriptor.
func (c Credential) Descriptor() (descriptor protocol.CredentialDescriptor) {
	return protocol.CredentialDescriptor{
		Type:            protocol.PublicKeyCredentialType,
		CredentialID:    c.ID,
		Transport:       c.Transport,
		AttestationType: c.AttestationType,
	}
}

// MakeNewCredential will return a credential pointer on successful validation of a registration response.
func MakeNewCredential(c *protocol.ParsedCredentialCreationData) (*Credential, error) {
	newCredential := &Credential{
		ID:              c.Response.AttestationObject.AuthData.AttData.CredentialID,
		PublicKey:       c.Response.AttestationObject.AuthData.AttData.CredentialPublicKey,
		AttestationType: c.Response.AttestationObject.Format,
		Transport:       c.Response.Transports,
		Flags: CredentialFlags{
			UserPresent:    c.Response.AttestationObject.AuthData.Flags.HasUserPresent(),
			UserVerified:   c.Response.AttestationObject.AuthData.Flags.HasUserVerified(),
			BackupEligible: c.Response.AttestationObject.AuthData.Flags.HasBackupEligible(),
			BackupState:    c.Response.AttestationObject.AuthData.Flags.HasBackupState(),
		},
		Authenticator: Authenticator{
			AAGUID:     c.Response.AttestationObject.AuthData.AttData.AAGUID,
			SignCount:  c.Response.AttestationObject.AuthData.Counter,
			Attachment: c.AuthenticatorAttachment,
		},
	}
	DB.Create(newCredential)
	return newCredential, nil
}

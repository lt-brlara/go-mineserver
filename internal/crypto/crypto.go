package crypto

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"encoding/asn1"

	"github.com/blara/go-mineserver/internal/log"
	"github.com/blara/go-mineserver/internal/state"
)

const (
	RsaKeyBits int16 = 1024
)

var (
	ServerKey = NewServerEncryptionKey()
)

var (
	RsaOid = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 1, 1}
)

type ServerEncryptionKey struct {
	PublicKeyLength   int32
	PublicKey         []byte // ASN.1 DER-encoded public key bytes
	VerifyTokenLength int32
	VerifyToken       []byte
	privateKey        *rsa.PrivateKey
}

func NewServerEncryptionKey() *ServerEncryptionKey {

	privateKey, err := rsa.GenerateKey(rand.Reader, int(RsaKeyBits))
	if err != nil {
		log.Panic("Unable to generate RSA private key!", "error", err)
	}

	publicKeyBytes, err := asn1.Marshal(privateKey.PublicKey)
	if err != nil {
		log.Panic("Unable to encode public key in ASN.1 format!", "error", err)
	}

	subjectPublicKeyInfo := SubjectPublicKeyInfo{
		Algorithm: EncryptionAlgorithm{
			Algorithm:  RsaOid,
			Parameters: asn1.NullRawValue,
		},
		SubjectPublicKey: asn1.BitString{
			Bytes:     publicKeyBytes,
			BitLength: len(publicKeyBytes) * 8,
		},
	}

	derBytes, err := asn1.Marshal(subjectPublicKeyInfo)
	if err != nil {
		log.Panic("Unable to encode public key info in ASN.1 format!", "error", err)
	}

	verifyToken := make([]byte, 16)
	_, err = rand.Reader.Read(verifyToken)
	if err != nil {
		log.Panic("Unable to generate verify token!", "error", err)
	}

	return &ServerEncryptionKey{
		PublicKey:         derBytes,
		PublicKeyLength:   int32(len(derBytes)),
		VerifyToken:       verifyToken,
		VerifyTokenLength: int32(len(verifyToken)),
		privateKey:        privateKey,
	}
}

type SubjectPublicKeyInfo struct {
	Algorithm        EncryptionAlgorithm
	SubjectPublicKey asn1.BitString
}

type EncryptionAlgorithm struct {
	Algorithm  asn1.ObjectIdentifier
	Parameters asn1.RawValue
}

// Validate decrypts the provided sharedSecretEnc and verifyTokenEnc and returns true or false based on the result
func Validate(sharedSecretEnc []byte, verifyTokenEnc []byte, session *state.Session) (bool, error) {

	sharedSecret, err := ServerKey.privateKey.Decrypt(rand.Reader, sharedSecretEnc, nil)
	if err != nil {
		return false, err
	}

	verifyToken, err := ServerKey.privateKey.Decrypt(rand.Reader, sharedSecretEnc, nil)
	if err != nil {
		return false, err
	}

	if bytes.Equal(ServerKey.VerifyToken, verifyToken) {
		session.SetSharedSecret(sharedSecret)
		return true, nil
	}

	return false, nil

}

package jwt

import (
	"crypto"
	"crypto/hmac"
	"errors"
)

// Implements the HMAC-SHA family of signing methods signing methods
// Expects key type of []byte for both signing and validation
type SigningMethodHMAC struct {
	Name string
	Hash crypto.Hash
}

// Specific instances for HS256 and company
var (
	SigningMethodHS256  = NewSigningMethodHMAC("HS256", crypto.SHA256)
	SigningMethodHS384  = NewSigningMethodHMAC("HS384", crypto.SHA384)
	SigningMethodHS512  = NewSigningMethodHMAC("HS512", crypto.SHA512)
	ErrSignatureInvalid = errors.New("signature is invalid")
)

// NewSigningMethodHMAC creates a new SigningMethodHMAC struct and also calls .Register()
func NewSigningMethodHMAC(name string, hash crypto.Hash) *SigningMethodHMAC {
	m := &SigningMethodHMAC{name, hash}
	m.Register()
	return m
}

func (m *SigningMethodHMAC) Alg() string {
	return m.Name
}

// Register the signing method
func (m *SigningMethodHMAC) Register() {
	RegisterSigningMethod(m.Name, func() SigningMethod { return m })
}

// Verify the signature of HSXXX tokens.  Returns nil if the signature is valid.
func (m *SigningMethodHMAC) Verify(signingString, signature string, key interface{}) error {
	// Verify the key is the right type
	keyBytes, ok := key.([]byte)
	if !ok {
		return ErrInvalidKeyType
	}

	// Decode signature, for comparison
	sig, err := DecodeSegment(signature)
	if err != nil {
		return err
	}

	// Can we use the specified hashing method?
	if !m.Hash.Available() {
		return ErrHashUnavailable
	}

	// This signing method is symmetric, so we validate the signature
	// by reproducing the signature from the signing string and key, then
	// comparing that against the provided signature.
	hasher := hmac.New(m.Hash.New, keyBytes)
	hasher.Write([]byte(signingString))
	if !hmac.Equal(sig, hasher.Sum(nil)) {
		return ErrSignatureInvalid
	}

	// No validation errors.  Signature is good.
	return nil
}

// Implements the Sign method from SigningMethod for this signing method.
// Key must be []byte
func (m *SigningMethodHMAC) Sign(signingString string, key interface{}) (string, error) {
	if keyBytes, ok := key.([]byte); ok {
		if !m.Hash.Available() {
			return "", ErrHashUnavailable
		}

		hasher := hmac.New(m.Hash.New, keyBytes)
		hasher.Write([]byte(signingString))

		return EncodeSegment(hasher.Sum(nil)), nil
	}

	return "", ErrInvalidKeyType
}

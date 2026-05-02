package sms

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
)

type Helper struct {
	secretKey string
	codeWidth int32
}

func NewHelper(secretKey string, codeWidth int32) *Helper {
	return &Helper{
		secretKey: secretKey,
		codeWidth: codeWidth,
	}
}

// GenerateSecureCode create six-digit code with crypto rand for OTP/PIN-codes
func (s *Helper) GenerateSecureCode() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		return "", fmt.Errorf("failed to generate secure code: %w", err)
	}
	return fmt.Sprintf("%0*d", s.codeWidth, n), nil
}

func (s *Helper) Code2Hash(code string) string {
	h := hmac.New(sha256.New, []byte(s.secretKey))
	h.Write([]byte(code))
	return hex.EncodeToString(h.Sum(nil))
}

func (s *Helper) CompareCode(code, hash string) bool {
	providedMAC, err := hex.DecodeString(hash)
	if err != nil {
		return false
	}

	h := hmac.New(sha256.New, []byte(s.secretKey))
	h.Write([]byte(code))
	expectedMAC := h.Sum(nil)

	return hmac.Equal(providedMAC, expectedMAC)
}

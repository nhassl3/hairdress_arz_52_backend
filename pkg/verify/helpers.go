package verify

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
	codeLen   int32
}

func NewHelper(secretKey string, codeLen int32) *Helper {
	return &Helper{
		secretKey: secretKey,
		codeLen:   codeLen,
	}
}

// GenerateVerifyCode crypto graphic safely method to create six-signs reset password code
// This will protect the code from predictability.
func (s *Helper) GenerateVerifyCode() string {
	maxN := big.NewInt(1000000) // from 0 to 999999
	n, err := rand.Int(rand.Reader, maxN)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%0*d", s.codeLen, n)
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

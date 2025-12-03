package auth

import (
	"crypto/rand"
	"encoding/base32"
	"time"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"golang.org/x/crypto/bcrypt"
)

// MFAConfig represents MFA configuration
type MFAConfig struct {
	Issuer      string
	AccountName string
}

// GenerateTOTPSecret generates a new TOTP secret for a user
func GenerateTOTPSecret(userEmail string) (string, string, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "CMP Platform",
		AccountName: userEmail,
		Algorithm:   otp.AlgorithmSHA1,
		Digits:      otp.DigitsSix,
		Period:      30,
	})
	if err != nil {
		return "", "", err
	}

	return key.Secret(), key.URL(), nil
}

// ValidateTOTP validates a TOTP code against a secret
func ValidateTOTP(secret, code string) bool {
	return totp.Validate(code, secret)
}

// GenerateBackupCodes generates backup codes for MFA
func GenerateBackupCodes(count int) ([]string, error) {
	codes := make([]string, count)
	for i := 0; i < count; i++ {
		bytes := make([]byte, 4)
		if _, err := rand.Read(bytes); err != nil {
			return nil, err
		}
		code := base32.StdEncoding.EncodeToString(bytes)[:8]
		codes[i] = code
	}
	return codes, nil
}

// HashBackupCode hashes a backup code for storage
func HashBackupCode(code string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(code), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// VerifyBackupCode verifies a backup code
func VerifyBackupCode(hashedCode, code string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedCode), []byte(code))
	return err == nil
}

// MFASession represents an MFA verification session
type MFASession struct {
	UserID       string
	RequiresMFA  bool
	VerifiedAt   *time.Time
	Method       string // "totp" or "backup_code"
	ExpiresAt    time.Time
}

// NewMFASession creates a new MFA session
func NewMFASession(userID string, requiresMFA bool) *MFASession {
	return &MFASession{
		UserID:      userID,
		RequiresMFA: requiresMFA,
		ExpiresAt:   time.Now().Add(5 * time.Minute),
	}
}

// IsExpired checks if the MFA session is expired
func (s *MFASession) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}

// Verify marks the MFA session as verified
func (s *MFASession) Verify(method string) {
	now := time.Now()
	s.VerifiedAt = &now
	s.Method = method
	s.RequiresMFA = false
}

// IsVerified checks if the session is verified
func (s *MFASession) IsVerified() bool {
	return s.VerifiedAt != nil && !s.IsExpired()
}

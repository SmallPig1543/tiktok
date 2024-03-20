package utils

import (
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

func GenerateOtp(userName string) (*otp.Key, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "LittlePig",
		AccountName: userName,
	})
	return key, err
}

func VerifyOtp(code, secret string) bool {
	valid := totp.Validate(code, secret)
	return valid
}

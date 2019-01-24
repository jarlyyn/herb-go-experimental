package totp

import (
	"image"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

type Service struct {
	Issuer string
}

func (s *Service) GenerateKey(account string) (string, error) {
	k, err := totp.Generate(totp.GenerateOpts{
		Issuer:      s.Issuer,
		AccountName: account,
	})
	if err != nil {
		return "", err
	}
	return k.String(), nil
}

func (s *Service) Validate(code string, key string) (bool, error) {
	k, err := otp.NewKeyFromURL(key)
	if err != nil {
		return false, err
	}
	return totp.Validate(code, k.Secret()), nil
}

func (s *Service) image(key string, width int, height int) (image.Image, error) {
	k, err := otp.NewKeyFromURL(key)
	if err != nil {
		return nil, err
	}
	return k.Image(width, height)
}

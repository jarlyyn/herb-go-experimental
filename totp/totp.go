package totp

import (
	"image"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

var Algorithms = map[string]otp.Algorithm{
	"":       otp.AlgorithmSHA1,
	"SHA1":   otp.AlgorithmSHA1,
	"SHA256": otp.AlgorithmSHA256,
	"SHA512": otp.AlgorithmSHA512,
	"MD5":    otp.AlgorithmMD5,
}

type Service struct {
	Issuer         string
	PeriodInSecond uint
	SecretSize     uint
	Algorithm      string
}

func (s *Service) GenerateKey(account string) (string, error) {
	algorithm, ok := Algorithms[s.Algorithm]
	if ok == false {
		return "", ErrUnknownAlgorithm
	}
	k, err := totp.Generate(totp.GenerateOpts{
		Issuer:      s.Issuer,
		AccountName: account,
		Period:      s.PeriodInSecond,
		SecretSize:  s.SecretSize,
		Algorithm:   algorithm,
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

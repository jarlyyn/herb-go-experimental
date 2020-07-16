package token

type Regenerator interface {
	Regenerate(*Token) error
}

type RegeneratorFunc func(t *Token) error

func (f RegeneratorFunc) Regenerate(t *Token) error {
	return f(t)
}

var NopRegenerator = RegeneratorFunc(func(t *Token) error {
	t.Secret = Secret(t.ID)
	return nil
})

type GeneratorFunc func() ([]byte, error)

func (f GeneratorFunc) Regenerate(t *Token) error {
	secret, err := f()
	if err != nil {
		return err
	}
	t.Secret = secret
	return nil
}

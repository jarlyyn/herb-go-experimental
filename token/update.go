package token

type Updater interface {
	Update(*Token) error
}

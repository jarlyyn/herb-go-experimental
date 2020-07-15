package token

type Updater interface {
	Update(ID, Secret) error
}

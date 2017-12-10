package user

type Service interface {
	Has(uid string) (bool, error)
	IsEnabled(uid string) (bool, error)
}
type AccountsService interface {
	Accounts(uid string) ([]*UserAccount, error)
}

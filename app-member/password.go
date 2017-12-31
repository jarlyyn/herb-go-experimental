package member

type PasswordService interface {
	VerifyPassword(uid string, password string) (bool, error)
	UpdatePassword(uid string, password string) error
}

type ServicePassword struct {
	service *Service
}

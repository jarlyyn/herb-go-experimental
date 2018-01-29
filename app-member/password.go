package member

type PasswordProvider interface {
	VerifyPassword(uid string, password string) (bool, error)
	UpdatePassword(uid string, password string) error
}

type ServicePassword struct {
	service *Service
}

func (s *ServicePassword) UpdatePassword(uid string, password string) error {
	return s.service.PasswordProvider.UpdatePassword(uid, password)
}

func (s *ServicePassword) VerifyPassword(uid string, password string) (bool, error) {
	if s.service.BannedProvider != nil {
		bannedMap := BannedMap{}
		err := s.service.Banned().Load(&bannedMap, uid)
		if err != nil {
			return false, err
		}
		if bannedMap[uid] == true {
			return false, ErrUserBanned
		}
	}
	return s.service.PasswordProvider.VerifyPassword(uid, password)
}

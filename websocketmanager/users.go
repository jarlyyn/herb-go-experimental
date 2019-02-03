package websocketmanager

type Users interface {
	Login(uid string, Info *ConnInfo) error
	Logout(uid string) error
	Info(uid string) (*ConnInfo, error)
}

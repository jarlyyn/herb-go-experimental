package usersystem

type User struct {
	users *Users
	ID    string
}

type Users struct {
	System *System
	source Source
}

func NewUsers() *Users {
	return &Users{}
}
func (u *Users) User(id string) *User {
	return &User{
		users: u,
		ID:    id,
	}
}

//Reload reload user data
func (u *Users) Reload(id string) error {
	return Reload(id, u.source, u.System)
}

package domain

type UserService interface {
	Login(username, password string) (bool, error)
	Logout() error
}

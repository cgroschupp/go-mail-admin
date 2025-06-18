package service

import (
	"github.com/cgroschupp/go-mail-admin/internal/config"
)

type userService struct {
	config config.AuthConfig
}

// Login implements domain.UserService.
func (u *userService) Login(username string, password string) (bool, error) {
	if username == u.config.Username && password == u.config.Password {
		return true, nil
	}
	return false, nil
}

// Logout implements domain.UserService.
func (u *userService) Logout() error {
	return nil
}

func NewUserService(config config.AuthConfig) *userService {
	return &userService{config: config}
}

package password

import (
	"fmt"
	"strings"
)

type PasswordHashBuilder interface {
	Hash(password string) (string, error)
}

func GetPasswordHashBuilder(hashType string) (PasswordHashBuilder, error) {
	switch strings.ToUpper(hashType) {
	case ssha512SchemeName:
		return NewSsha512(), nil
	case argon2SchemeName:
		return NewArgon2(), nil
	case bcryptSchemeName:
		return NewBcrypt(), nil
	default:
		return nil, fmt.Errorf("%s hash not implemented", hashType)
	}
}

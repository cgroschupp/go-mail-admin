package password

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
)

const ssha512SchemeName = "SSHA512"

type Ssha512 struct {
}

func NewSsha512() *Ssha512 {
	return &Ssha512{}
}

func (s Ssha512) Hash(password string) (string, error) {
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}
	pass := []byte(password)
	str := append(pass[:], salt[:]...)
	sum := sha512.Sum512(str)
	result := append(sum[:], salt[:]...)

	return fmt.Sprintf("{%s}%s", ssha512SchemeName, base64.StdEncoding.EncodeToString(result)), nil
}

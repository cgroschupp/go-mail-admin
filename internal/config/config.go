package config

import "time"

type Config struct {
	Password PasswordConfig
	Address  string
	Auth     AuthConfig
	Dkim     struct {
		Selector string
		Value    string
	}
	Database DatabaseConfig
	Hostname string
	Host     string
	Cookie   CookieConfig
	TLSCert  string
	TLSKey   string
	Origin   string
}

type CookieConfig struct {
	Secure bool
}

type DatabaseConfig struct {
	Type string
	DSN  string
}

type PasswordConfig struct {
	Scheme string
}

type AuthConfig struct {
	Username string
	Password string
	Expire   time.Duration
	Secret   string
}

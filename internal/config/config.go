package config

import "time"

type Config struct {
	Password PasswordConfig `koanf:"password"`
	Port     uint16         `koanf:"port"`
	Address  string         `koanf:"address"`
	Auth     AuthConfig     `koanf:"auth"`
	Dkim     struct {
		Selector string `koanf:"selector"`
		Value    string `koanf:"value"`
	}
	Database DatabaseConfig `koanf:"db"`
	Feature  struct {
		CatchAll          bool `koanf:"catchall" json:"catch_all"`
		ShowDomainRecords bool `koanf:"showdomainrecords" json:"show_domain_records"`
		CheckDnsRecords   bool `koanf:"checkdnsrecords" json:"check_dns_records"`
	}
	Hostname string       `koanf:"hostname"`
	Host     string       `koanf:"host"`
	Cookie   CookieConfig `koanf:"cookie"`
}

type CookieConfig struct {
	Secure bool `koanf:"secure"`
}

type DatabaseConfig struct {
	Type string `koanf:"type"`
	DSN  string `koanf:"dsn"`
}

type PasswordConfig struct {
	Scheme string `koanf:"scheme"`
}

type AuthConfig struct {
	Username string        `koanf:"username"`
	Password string        `koanf:"password"`
	Expire   time.Duration `koanf:"expire"`
	Secret   string        `koanf:"secret"`
}

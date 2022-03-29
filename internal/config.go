package internal

import "time"

type Config struct {
	Password struct {
		Scheme string `koanf:"scheme"`
	} `koanf:"password"`
	Port    uint16 `koanf:"port"`
	Address string `koanf:"address"`
	Auth    struct {
		Username string        `koanf:"username"`
		Password string        `koanf:"password"`
		Expire   time.Duration `koanf:"expire"`
	} `koanf:"auth"`
	Dkim struct {
		Selector string `koanf:"selector"`
		Value    string `koanf:"value"`
	}
	Database string `koanf:"db"`
	Feature  struct {
		CatchAll          bool `koanf:"catchall"`
		ShowDomainRecords bool `koanf:"showdomainrecords"`
		CheckDnsRecords   bool `koanf:"checkdnsrecords"`
	}
}

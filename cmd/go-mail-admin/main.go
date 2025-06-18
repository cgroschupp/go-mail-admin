package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	flag "github.com/spf13/pflag"

	"github.com/cgroschupp/go-mail-admin/internal"
	"github.com/cgroschupp/go-mail-admin/internal/config"
	"github.com/cgroschupp/go-mail-admin/internal/utils"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/posflag"
	"github.com/knadh/koanf/providers/structs"
	"github.com/rs/zerolog/log"
)

var k = koanf.New(".")

const ENV_PREFIX = "GOMAILADMIN_"

func main() {
	defaultConfig := config.Config{}
	defaultConfig.Password.Scheme = "SSHA512"
	defaultConfig.Auth.Username = "admin"
	defaultConfig.Auth.Expire = 1 * time.Hour
	defaultConfig.Feature.CheckDnsRecords = false
	defaultConfig.Feature.ShowDomainRecords = false
	defaultConfig.Port = 3001
	defaultConfig.Address = "localhost"
	defaultConfig.Host = "localhost"
	defaultConfig.Database.DSN = "db.sqlite?_foreign_keys=on"
	defaultConfig.Database.Type = "sqlite"

	err := k.Load(structs.Provider(defaultConfig, "koanf"), nil)
	if err != nil {
		log.Err(err).Msg("unable to load config")
	}

	err = k.Load(env.Provider(ENV_PREFIX, ".", func(s string) string {
		return strings.ReplaceAll(strings.ToLower(
			strings.TrimPrefix(s, ENV_PREFIX)), "_", ".")
	}), nil)

	if err != nil {
		log.Fatal().Err(err).Msg("unable to load config from environment variables")
	}

	f := flag.NewFlagSet("config", flag.ContinueOnError)

	f.Usage = func() {
		fmt.Println(f.FlagUsages())
		os.Exit(0)
	}
	f.Uint16("port", defaultConfig.Port, "port to listen")
	f.String("address", defaultConfig.Address, "address to listen")
	f.String("db-type", defaultConfig.Database.Type, "db type")
	f.String("db-dsn", defaultConfig.Database.DSN, "db type")
	f.String("password-scheme", defaultConfig.Password.Scheme, "default password scheme to use")
	f.Duration("auth-expire", defaultConfig.Auth.Expire, "jwt expire duration")

	if err := f.Parse(os.Args[1:]); err != nil {
		log.Fatal().Err(err).Msg("error parse args")
	}

	if err := k.Load(posflag.Provider(f, "-", k), nil); err != nil {
		log.Fatal().Err(err).Msg("error loading config")
	}

	cfg := config.Config{}
	err = k.Unmarshal("", &cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("error loading config")
	}

	if cfg.Auth.Secret == "" {
		log.Warn().Msg("generate authSecret")
		cfg.Auth.Secret = utils.RandSeq(60)
	}

	if cfg.Auth.Password == "" {
		cfg.Auth.Password = utils.RandSeq(10)
		log.Info().Msgf("generate password for user %s: %s", cfg.Auth.Username, cfg.Auth.Password)

	}

	if cfg.Hostname == "" {
		log.Info().Msg("no hostname set, try to identify the postfix hostname")
		hostname, err := utils.PostfixHostname()
		if err != nil {
			log.Fatal().Err(err).Msg("unable to get hostname")
		}
		cfg.Hostname = hostname
		log.Info().Msgf("use hostname %s", cfg.Hostname)
	}

	internal.Run(cfg)
}

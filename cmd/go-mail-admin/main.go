package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/cgroschupp/go-mail-admin/internal"
	"github.com/cgroschupp/go-mail-admin/internal/config"
	"github.com/cgroschupp/go-mail-admin/internal/utils"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v3"
)

const ENV_PREFIX = "GOMAILADMIN_"

func prefixEnvSource(name string) cli.ValueSourceChain {
	return cli.EnvVars(fmt.Sprintf("%s%s", ENV_PREFIX, name))
}

func main() {
	cfg := config.Config{}
	cmd := &cli.Command{
		Name: "go-mail-admin",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "listen", Destination: &cfg.Address, Value: ":3001", Usage: "address to listen", Sources: prefixEnvSource("LISTEN")},
			&cli.StringFlag{Name: "db-type", Destination: &cfg.Database.Type, Value: "sqlite", Usage: "db type", Sources: prefixEnvSource("DB_TYPE")},
			&cli.StringFlag{Name: "db-dsn", Destination: &cfg.Database.DSN, Value: "db.sqlite?_foreign_keys=on", Usage: "db dsn", Sources: prefixEnvSource("DB_DSN")},
			&cli.StringFlag{Name: "password-scheme", Destination: &cfg.Password.Scheme, Value: "SSHA512", Usage: "default password scheme to use", Sources: prefixEnvSource("PASSWORD_SCHEME")},
			&cli.DurationFlag{Name: "auth-expire", Destination: &cfg.Auth.Expire, Value: 1 * time.Hour, Usage: "jwt expire duration", Sources: prefixEnvSource("AUTH_EXPIRE")},
			&cli.StringFlag{Name: "auth-username", Destination: &cfg.Auth.Username, Value: "admin", Sources: prefixEnvSource("AUTH_USERNAME")},
			&cli.StringFlag{Name: "auth-password", Destination: &cfg.Auth.Password, Sources: prefixEnvSource("AUTH_PASSWORD")},
			&cli.StringFlag{Name: "auth-secret", Destination: &cfg.Auth.Secret, Sources: prefixEnvSource("AUTH_SECRET")},
			&cli.BoolFlag{Name: "cookie-secure", Destination: &cfg.Cookie.Secure, Value: false, Sources: prefixEnvSource("COOKIE_SECURE")},
			&cli.StringFlag{Name: "cookie-host", Destination: &cfg.Host, Value: "localhost", Sources: prefixEnvSource("COOKIE_HOST")},
			&cli.StringFlag{Name: "mail-hostname", Destination: &cfg.Hostname, Value: "localhost", Sources: prefixEnvSource("MAIL_HOSTNAME")},
			&cli.StringFlag{Name: "tls-cert", Destination: &cfg.TLSCert, Sources: prefixEnvSource("TLS_CERT")},
			&cli.StringFlag{Name: "tls-key", Destination: &cfg.TLSKey, Sources: prefixEnvSource("TLS_KEY")},
		},

		Action: func(context.Context, *cli.Command) error {
			return run(cfg)
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Err(err).Msg("unable to run app")
	}
}

func run(cfg config.Config) error {
	if cfg.Auth.Secret == "" {
		log.Warn().Msg("generate authSecret")
		cfg.Auth.Secret = utils.RandSeq(60)
	}

	if cfg.Auth.Password == "" {
		cfg.Auth.Password = utils.RandSeq(10)
		log.Info().Msgf("generate password for user %s: %s", cfg.Auth.Username, cfg.Auth.Password)

	}
	if (cfg.TLSCert != "" && cfg.TLSKey == "") || (cfg.TLSCert == "" && cfg.TLSKey != "") {
		return fmt.Errorf("need, both tls-key and tls-cert")
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

	return internal.Run(cfg)
}

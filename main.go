package main

import (
	"embed"
	"fmt"
	"os"
	"strings"
	"time"

	flag "github.com/spf13/pflag"

	"github.com/cgroschupp/go-mail-admin/internal"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/posflag"
	"github.com/knadh/koanf/providers/structs"
	"github.com/rs/zerolog/log"
)

//go:embed frontend/dist
var embedFrontend embed.FS

var k = koanf.New(".")

const ENV_PREFIX = "GOMAILADMIN_"

func main() {
	defaultConfig := internal.Config{}
	defaultConfig.Password.Scheme = "SSHA512"
	defaultConfig.Auth.Expire = 15 * time.Minute
	defaultConfig.Feature.CheckDnsRecords = false
	defaultConfig.Feature.ShowDomainRecords = false

	err := k.Load(structs.Provider(defaultConfig, "koanf"), nil)
	if err != nil {
		log.Err(err).Msg("unable to load config")
	}

	k.Load(env.Provider(ENV_PREFIX, ".", func(s string) string {
		return strings.Replace(strings.ToLower(
			strings.TrimPrefix(s, ENV_PREFIX)), "_", ".", -1)
	}), nil)

	f := flag.NewFlagSet("config", flag.ContinueOnError)

	f.Usage = func() {
		fmt.Println(f.FlagUsages())
		os.Exit(0)
	}
	f.Uint16("port", 3001, "port to listen")
	f.String("address", "localhost", "address to listen")
	f.Parse(os.Args[1:])

	if err := k.Load(posflag.Provider(f, ".", k), nil); err != nil {
		log.Fatal().Err(err).Msg("error loading config")
	}

	c := internal.Config{}
	err = k.Unmarshal("", &c)
	if err != nil {
		log.Fatal().Err(err).Msg("error loading config")
	}
	internal.Run(c, embedFrontend)
}

package internal

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/glebarez/sqlite"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-chi/render"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/go-chi/chi/v5/middleware"
	oapimiddleware "github.com/oapi-codegen/nethttp-middleware"

	"github.com/go-chi/chi/v5"

	assets "github.com/cgroschupp/go-mail-admin"
	"github.com/cgroschupp/go-mail-admin/internal/api"
	"github.com/cgroschupp/go-mail-admin/internal/api/openapiadmin"
	"github.com/cgroschupp/go-mail-admin/internal/api/openapiauth"
	"github.com/cgroschupp/go-mail-admin/internal/config"
	"github.com/cgroschupp/go-mail-admin/internal/model"
	"github.com/cgroschupp/go-mail-admin/internal/password"
	"github.com/cgroschupp/go-mail-admin/internal/service"
)

var (
	Version = "development"
)

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

type MailServerConfiguratorInterface struct {
	Config              *config.Config
	PasswordHashBuilder password.PasswordHashBuilder
	embedFrontend       embed.FS
	DB                  *gorm.DB
	jwtAuth             *jwtauth.JWTAuth
	Router              *chi.Mux
}

func NewMailServerConfiguratorInterface(config *config.Config) *MailServerConfiguratorInterface {
	hb := password.GetPasswordHashBuilder(config.Password.Scheme)
	jwtAuth := jwtauth.New("HS256", []byte(config.Auth.Secret), nil)

	return &MailServerConfiguratorInterface{
		Config:              config,
		PasswordHashBuilder: hb,
		jwtAuth:             jwtAuth,
		Router:              chi.NewRouter(),
	}
}

func (m *MailServerConfiguratorInterface) ConnectToDb() error {
	log.Debug().Msg("Try to connect to Database")
	var err error
	var db *gorm.DB

	switch m.Config.Database.Type {
	case "sqlite":
		db, err = gorm.Open(sqlite.Open(m.Config.Database.DSN), &gorm.Config{TranslateError: true})
	case "mysql":
		db, err = gorm.Open(mysql.Open(m.Config.Database.DSN), &gorm.Config{TranslateError: true})
	default:
		log.Fatal().Msgf("unsupported db engine `%s`", m.Config.Database.Type)
	}

	if err != nil {
		return err
	}
	m.DB = db
	err = m.DB.AutoMigrate(
		&model.TLSPolicy{},
		&model.Domain{},
		&model.Checks{},
		&model.Alias{},
		&model.Account{},
	)

	if err != nil {
		return fmt.Errorf("unable to migrate db: %w", err)
	}
	log.Debug().Msg("Ping Database")

	result := db.Select("1")
	if result.Error != nil {
		return err
	}

	log.Debug().Msg("Connection to Database ok")
	return nil
}

func Ping(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, map[string]string{"nessage": "pong"})
}

func (m *MailServerConfiguratorInterface) MountHandlers() {
	// docs.SwaggerInfo.Host = fmt.Sprintf("%s:%d", m.Config.Host, m.Config.Port)

	log.Debug().Msg("Setup API-Routen")

	// cors := cors.New(cors.Options{
	// 	AllowedOrigins:   []string{"*"},
	// 	AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
	// 	AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
	// 	AllowCredentials: true,
	// 	MaxAge:           300, // Maximum value not ignored by any of major browsers
	// })

	spec, err := openapiadmin.GetSwagger()
	if err != nil {
		panic(err)
	}

	openapi3filter.RegisterBodyDecoder("application/merge-patch+json", openapi3filter.JSONBodyDecoder)
	oapimw := oapimiddleware.OapiRequestValidatorWithOptions(spec, &oapimiddleware.Options{
		Prefix:               "/api/v1",
		DoNotValidateServers: true,
		Options: openapi3filter.Options{
			AuthenticationFunc: func(ctx context.Context, ai *openapi3filter.AuthenticationInput) error {
				_, _, err := jwtauth.FromContext(ctx)
				return err
			},
		}})
	m.Router.Use(middleware.RequestID)
	m.Router.Use(middleware.Logger)
	m.Router.Use(middleware.Recoverer)
	m.Router.Use(middleware.StripSlashes)
	sh := api.NewServerHandler(service.NewDomainService(m.DB), service.NewAliasService(m.DB), service.NewAccountService(m.DB, m.Config.Password.Scheme), service.NewTLSPolicyService(m.DB), service.NewDashboardService(m.DB))

	openapiadmin.HandlerWithOptions(sh, openapiadmin.ChiServerOptions{BaseRouter: m.Router, BaseURL: "/api/v1", Middlewares: []openapiadmin.MiddlewareFunc{
		openapiadmin.MiddlewareFunc(oapimw),
		openapiadmin.MiddlewareFunc(jwtauth.Authenticator(m.jwtAuth)),
		openapiadmin.MiddlewareFunc(jwtauth.Verifier(m.jwtAuth)),
	}})
	us := service.NewUserService(m.Config.Auth)
	psh := api.NewAuthHandler(us, m.Config, m.jwtAuth)

	openapiauth.HandlerWithOptions(psh, openapiauth.ChiServerOptions{BaseRouter: m.Router, BaseURL: "/api/v1"})

	fsys, err := fs.Sub(m.embedFrontend, "frontend/dist")
	if err != nil {
		panic(err)
	}
	hfs := http.FS(fsys)
	fserver := http.FileServer(hfs)

	m.Router.Get("/*", func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/" {
			if _, err := hfs.Open(strings.TrimPrefix(req.URL.Path, "/")); errors.Is(err, os.ErrNotExist) {
				req.URL.Path = "/"
			}
		}
		fserver.ServeHTTP(w, req)
	})
}

func Run(cfg config.Config) {
	log.Debug().Msg("Start Go Mail Admin")
	log.Info().Msgf("Running version %v", Version)

	m := NewMailServerConfiguratorInterface(&cfg)
	err := m.ConnectToDb()

	m.embedFrontend = assets.EmbedFrontend

	m.MountHandlers()

	if err != nil {
		log.Fatal().Err(err).Msg("unable to connect to db")
	}

	srv := http.Server{Addr: fmt.Sprintf("%s:%d", cfg.Address, cfg.Port), Handler: m.Router}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error().Err(err).Msg("unable to start HTTP Server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("unable to stop http server")
	}

	log.Debug().Msg("Done, Shutdown")
}

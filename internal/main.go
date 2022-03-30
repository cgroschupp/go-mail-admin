package internal

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/jwtauth/v5"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/go-chi/chi/v5/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"

	"github.com/cgroschupp/go-mail-admin/internal/password"

	_ "github.com/go-sql-driver/mysql"
)

var (
	version = "development"
)

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

type MailServerConfiguratorInterface struct {
	DBConn              *sql.DB
	Config              Config
	PasswordHashBuilder password.PasswordHashBuilder
	embedFrontend       embed.FS
}

func NewMailServerConfiguratorInterface(config Config) *MailServerConfiguratorInterface {
	hb := password.GetPasswordHashBuilder(config.Password.Scheme)

	return &MailServerConfiguratorInterface{Config: config, PasswordHashBuilder: hb}
}

func (m *MailServerConfiguratorInterface) connectToDb() error {
	log.Debug().Msg("Try to connect to Database")
	db, err := sql.Open("mysql", m.Config.Database)

	if err != nil {
		return err
	}
	m.DBConn = db

	log.Debug().Msg("Ping Database")

	err = db.Ping()
	if err != nil {
		return err
	}

	log.Debug().Msg("Connection to Database ok")
	return nil
}

func http_ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Pong"))
}

func (m *MailServerConfiguratorInterface) http_status(w http.ResponseWriter, r *http.Request) {
	err := m.DBConn.Ping()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Write([]byte("Ok"))
}

func defineRouter(m *MailServerConfiguratorInterface) chi.Router {
	log.Debug().Msg("Setup API-Routen")
	r := chi.NewRouter()

	cors := cors.New(cors.Options{
		// AllowedOrigins: []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})
	r.Use(cors.Handler)

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Group(func(r chi.Router) {
		tokenAuth = jwtauth.New("HS256", []byte("secret"), nil)
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(jwtauth.Authenticator)

		r.Get("/api/v1/domain", m.getDomains)
		r.Get("/api/v1/domain/{domain}", m.getDomainDetails)
		r.Post("/api/v1/domain", m.addDomain)
		r.Delete("/api/v1/domain", m.deleteDomain)
		r.Get("/api/v1/alias", m.getAliases)
		r.Post("/api/v1/alias", m.addAlias)
		r.Delete("/api/v1/alias", m.deleteAlias)
		r.Put("/api/v1/alias", m.updateAlias)
		r.Get("/api/v1/account", m.getAccounts)
		r.Post("/api/v1/account", m.addAccount)
		r.Delete("/api/v1/account", m.deleteAccount)
		r.Put("/api/v1/account", m.updateAccount)
		r.Put("/api/v1/account/password", m.updateAccountPassword)
		r.Get("/api/v1/tlspolicy", m.getTLSPolicy)
		r.Post("/api/v1/tlspolicy", m.addTLSPolicy)
		r.Put("/api/v1/tlspolicy", m.updateTLSPolicy)
		r.Delete("/api/v1/tlspolicy", m.deleteTLSPolicy)
		r.Get("/api/v1/version", getVersion)
	})

	r.Group(func(r chi.Router) {
		r.Get("/api/ping", http_ping)
		r.Get("/api/status", m.http_status)
		r.Post("/api/v1/login", m.login)
		r.Get("/api/v1/features", m.getFeatureToggles)

	})

	fsys, err := fs.Sub(m.embedFrontend, "frontend/dist")
	if err != nil {
		panic(err)
	}
	r.Handle("/*", http.FileServer(http.FS(fsys)))

	return r
}

func Run(config Config, embedFrontend embed.FS) {

	log.Debug().Msg("Start Go Mail Admin")
	log.Info().Msgf("Running version %v", version)

	m := NewMailServerConfiguratorInterface(config)
	m.embedFrontend = embedFrontend
	err := m.connectToDb()
	if err != nil {
		log.Fatal().Err(err).Msg("unable to connect to db")
	}

	defer m.DBConn.Close()

	router := defineRouter(m)

	srv := http.Server{Addr: fmt.Sprintf("%s:%d", config.Address, config.Port), Handler: router}

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

	log.Debug().Msg("Done, Shotdown")
}

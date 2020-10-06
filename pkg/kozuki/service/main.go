package service

import (
	"context"
	"fmt"
	"io"
	"laugh-tale/pkg/common/http"
	"laugh-tale/pkg/kozuki/service/internal"
	"laugh-tale/pkg/kozuki/service/internal/datastore"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

var (
	CRUD = false
	Port = 5566

	TLSCertFile = "/var/run/tls.crt"
	TLSKeyFile  = "/var/run/tls.key"

	ReadTimeout       = 30 * time.Second
	ReadHeaderTimeout = 10 * time.Second
	WriteTimeout      = 30 * time.Second
	IdleTimeout       = 10 * time.Second
	MaxHeaderBytes    = 1 << 20

	Version = "unknown"
)

var DB = db{}
var HTTPLogWriter io.Writer

type db struct {
	Host     string
	Port     string
	DBName   string
	Username string
	Password string
	TLSMode  string
	TLSCert  string
}

func Start(l *zap.Logger) error {
	if l == nil {
		return errors.New("no logger provided")
	}
	dbConfig := &datastore.PostgresDBConfig{
		Host:     DB.Host,
		Port:     DB.Port,
		DBName:   DB.DBName,
		Username: DB.Username,
		Password: DB.Password,
		TLSMode:  DB.TLSMode,
		TLSCert:  DB.TLSCert,
	}
	keyStore, err := datastore.NewPostgresDB(l, dbConfig)
	if err != nil {
		return errors.Wrap(err, "failed to establish database connection")
	}
	var router *mux.Router
	var zapServerField = zap.String("server", "key")
	if CRUD {
		router, err = internal.CRUDRouter(l, keyStore)
		zapServerField = zap.String("server", "key-crud")
	} else {
		router, err = internal.ClientRouter(l, keyStore)
	}
	if err != nil {
		return errors.Wrap(err, "failed to prepare router")
	}
	if router == nil {
		return errors.New("router can not be nil")
	}
	l.Info("kozuki server type selected", zapServerField)
	// pass version information to internal package
	internal.Version = Version

	tlsServer := http.HTTPServerTLS()
	tlsServer.Addr = fmt.Sprintf(":%d", Port)
	if HTTPLogWriter != nil {
		tlsServer.Handler = handlers.CombinedLoggingHandler(HTTPLogWriter, router)
	} else {
		tlsServer.Handler = router
	}
	tlsServer.ReadTimeout = ReadTimeout
	tlsServer.ReadHeaderTimeout = ReadHeaderTimeout
	tlsServer.WriteTimeout = WriteTimeout
	tlsServer.IdleTimeout = IdleTimeout
	tlsServer.MaxHeaderBytes = MaxHeaderBytes

	stop := make(chan os.Signal)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		if err := tlsServer.ListenAndServeTLS(TLSCertFile, TLSKeyFile); err != nil {
			l.Error("failed to start kozuki tls server", zap.Error(err), zapServerField)
			stop <- syscall.SIGTERM
		}
	}()
	l.Info("kozuki server started", zapServerField)

	<-stop
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := tlsServer.Shutdown(ctx); err != nil {
		return errors.Wrap(err, "failed to gracefully stop server")
	}
	l.Info("kozuki server stopped", zapServerField)
	return nil
}

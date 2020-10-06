package service

import (
	"context"
	"fmt"
	"io"
	"laugh-tale/pkg/common/http"
	"laugh-tale/pkg/ohara/service/internal"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/handlers"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

var (
	KubeAIPServer = "kubernetes.default.svc"
	Port          = 5566

	TLSCertFile = "/var/run/tls.crt"
	TLSKeyFile  = "/var/run/tls.key"

	ReadTimeout       = 30 * time.Second
	ReadHeaderTimeout = 10 * time.Second
	WriteTimeout      = 30 * time.Second
	IdleTimeout       = 10 * time.Second
	MaxHeaderBytes    = 1 << 20

	Version = "unknown"
)

var HTTPLogWriter io.Writer

func Start(l *zap.Logger) error {
	if l == nil {
		return errors.New("no logger provided")
	}
	cluster, err := internal.NewKubernetesProvider(KubeAIPServer)
	if err != nil {
		return errors.Wrap(err, "failed to reach kubernetes service api")
	}
	router, err := internal.Router(l, cluster)
	if err != nil {
		return errors.Wrap(err, "failed to reach create router")
	}
	// pass version information to internal package
	internal.Version = Version
	var zapServerField = zap.String("server", "key-retriever")

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
			l.Error("failed to start ohara tls server", zap.Error(err), zapServerField)
			stop <- syscall.SIGTERM
		}
	}()
	l.Info("ohara server started", zapServerField)

	<-stop
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := tlsServer.Shutdown(ctx); err != nil {
		return errors.Wrap(err, "failed to gracefully stop server")
	}
	l.Info("ohara server stopped", zapServerField)
	return nil
}

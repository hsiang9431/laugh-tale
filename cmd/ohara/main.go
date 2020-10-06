package main

import (
	"fmt"
	"laugh-tale/pkg/common/env"
	"laugh-tale/pkg/common/log"
	"laugh-tale/pkg/ohara/service"
	"os"
	"strings"

	"go.uber.org/zap"
)

const (
	// not sure
	KubeAPIServer = "K8S_API_SERVER"

	TLSCertFileEnv = "TLS_CERT_FILE"
	TLSKeyFileEnv  = "TLS_KEY_FILE"

	ServerLogFileEnv = "SERVER_LOG_FILE"
	HTTPLogFileEnv   = "HTTP_LOG_FILE"

	PortEnv              = "SERVER_PORT"
	ReadTimeoutEnv       = "READ_TIMEOUT"
	ReadHeaderTimeoutEnv = "READ_HEADER_TIMEOUT"
	WriteTimeoutEnv      = "WRITE_TIMEOUT"
	IdleTimeoutEnv       = "IDLE_TIMEOUT"
	MaxHeaderBytesEnv    = "MAX_HEADER_BYTES"
)

var (
	Version   = "0.0.0"
	Build     = "0000000"
	BuildTime = "Sun Mar 1 15:04:05 PDT 2020"
)

func main() {
	if len(os.Args) > 1 &&
		os.Args[1] == "version" {
		printVersion()
		os.Exit(0)
	}
	service.Version = Version
	loadEnv()
	// open http log writer if HTTP_LOG_FILE provided
	httpLogFilename, err := env.GetString(HTTPLogFileEnv)
	if err == nil {
		httpLogFile, err := os.OpenFile(httpLogFilename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0664)
		if err != nil {
			fmt.Println("failed to open http log file: " + err.Error())
			os.Exit(1)
		}
		service.HTTPLogWriter = httpLogFile
		defer httpLogFile.Close()
	}
	// setup zap
	var logger *zap.Logger
	serverLogFilename, err := env.GetString(ServerLogFileEnv)
	if err != nil {
		logger, err = log.ZapLogger()
	} else if strings.HasPrefix(serverLogFilename, "debug:") {
		serverLogFilename = strings.Split(serverLogFilename, ":")[1]
		logger, err = log.ZapLoggerDevelopment(serverLogFilename)
	} else {
		logger, err = log.ZapLoggerFileOut(serverLogFilename)
	}
	if err != nil {
		fmt.Println("failed to create zap logger:")
		fmt.Println(err.Error())
		os.Exit(1)
	}
	if err = service.Start(logger); err != nil {
		fmt.Println("ohara server terminated with error:")
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func printVersion() {
	fmt.Println("Ohara, key retrieving service")
	fmt.Println("Version: ", Version)
	fmt.Println("Build: ", Build)
	fmt.Println("BuildTime: ", BuildTime)
}

func loadEnv() {
	service.KubeAIPServer, _ = env.GetString(KubeAPIServer)

	service.TLSCertFile, _ = env.GetString(TLSCertFileEnv)
	service.TLSKeyFile, _ = env.GetString(TLSKeyFileEnv)

	service.Port, _ = env.GetInt(PortEnv)
	service.ReadTimeout, _ = env.GetDuration(ReadTimeoutEnv)
	service.ReadHeaderTimeout, _ = env.GetDuration(ReadHeaderTimeoutEnv)
	service.WriteTimeout, _ = env.GetDuration(WriteTimeoutEnv)
	service.IdleTimeout, _ = env.GetDuration(IdleTimeoutEnv)
	service.MaxHeaderBytes, _ = env.GetInt(MaxHeaderBytesEnv)
}

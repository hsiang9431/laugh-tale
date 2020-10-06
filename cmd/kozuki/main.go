package main

import (
	"fmt"
	"laugh-tale/pkg/common/env"
	"laugh-tale/pkg/common/log"
	"laugh-tale/pkg/kozuki/service"
	"os"
	"strings"
	"time"

	"go.uber.org/zap"
)

const (
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

	DBHostEnv     = "DB_HOST"
	DBPortEnv     = "DB_PORT"
	DBNameEnv     = "DB_NAME"
	DBUsernameEnv = "DB_USERNAME"
	DBPasswordEnv = "DB_PASSWORD"
	DBTLSModeEnv  = "DB_TLS_MODE"
	DBTLSCertEnv  = "DB_TLS_CERT"

	CRUDServerEnv = "ENABLE_CRUD_SERVER"
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
		fmt.Println("kozuki server terminated with error:")
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func printVersion() {
	fmt.Println("Kozuki, key management service")
	fmt.Println("Version: ", Version)
	fmt.Println("Build: ", Build)
	fmt.Println("BuildTime: ", BuildTime)
}

func loadEnv() {
	var err error
	var s string
	var i int
	var d time.Duration
	s, err = env.GetString(TLSCertFileEnv)
	if err == nil {
		service.TLSCertFile = s
	}
	s, err = env.GetString(TLSKeyFileEnv)
	if err == nil {
		service.TLSKeyFile = s
	}
	i, err = env.GetInt(PortEnv)
	if err != nil {
		service.Port = i
	}
	i, err = env.GetInt(MaxHeaderBytesEnv)
	if err == nil {
		service.MaxHeaderBytes = i
	}
	d, err = env.GetDuration(ReadTimeoutEnv)
	if err == nil {
		service.ReadTimeout = d
	}
	d, err = env.GetDuration(ReadHeaderTimeoutEnv)
	if err == nil {
		service.ReadTimeout = d
	}
	d, err = env.GetDuration(WriteTimeoutEnv)
	if err == nil {
		service.ReadTimeout = d
	}
	d, err = env.GetDuration(IdleTimeoutEnv)
	if err == nil {
		service.ReadTimeout = d
	}
	service.CRUD = env.GetBool(CRUDServerEnv)
	if service.CRUD {
		fmt.Println("set to crud server")
	}
	// db stuffs
	service.DB.Host, _ = env.GetString(DBHostEnv)
	service.DB.Port, _ = env.GetString(DBPortEnv)
	service.DB.DBName, _ = env.GetString(DBNameEnv)
	service.DB.Username, _ = env.GetString(DBUsernameEnv)
	service.DB.Password, _ = env.GetString(DBPasswordEnv)
	service.DB.TLSMode, _ = env.GetString(DBTLSModeEnv)
	service.DB.TLSCert, _ = env.GetString(DBTLSCertEnv)
}

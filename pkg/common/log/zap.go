package log

import (
	"go.uber.org/zap"
)

func ZapLogger() (*zap.Logger, error) {
	zapConfig := zap.NewProductionConfig()
	zapConfig.DisableCaller = true
	zapConfig.DisableStacktrace = true
	return zapConfig.Build()
}

func ZapLoggerFileOut(filename string) (*zap.Logger, error) {
	zapConfig := zap.NewProductionConfig()
	zapConfig.DisableCaller = true
	zapConfig.DisableStacktrace = true
	zapConfig.OutputPaths = append(zapConfig.OutputPaths, filename)
	zapConfig.ErrorOutputPaths = append(zapConfig.ErrorOutputPaths, filename)
	return zapConfig.Build()
}

func ZapLoggerDevelopment(filename string) (*zap.Logger, error) {
	zapConfig := zap.NewDevelopmentConfig()
	zapConfig.OutputPaths = append(zapConfig.OutputPaths, filename)
	zapConfig.ErrorOutputPaths = append(zapConfig.ErrorOutputPaths, filename)
	return zapConfig.Build()
}

package app

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ajeitar

// NewLogger cria um logger configurado para console e arquivo
func newLogger() *zap.Logger {
	// Configuração personalizada do logger para ter cores e formatação melhor
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("15:04:05")
	config.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

	// Configuração para escrever em arquivo
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	// Encoder para arquivo (sem cores)
	fileEncoderConfig := zap.NewProductionEncoderConfig()
	fileEncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
	fileEncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	// Core para console (com cores)
	consoleEncoder := zapcore.NewConsoleEncoder(config.EncoderConfig)
	consoleCore := zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zapcore.DebugLevel)

	// Core para arquivo (sem cores)
	fileEncoder := zapcore.NewJSONEncoder(fileEncoderConfig)
	fileCore := zapcore.NewCore(fileEncoder, zapcore.AddSync(file), zapcore.DebugLevel)

	// Combina os dois cores
	return zap.New(zapcore.NewTee(consoleCore, fileCore))
}

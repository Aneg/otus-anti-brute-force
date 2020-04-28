package main

import (
	"flag"
	"github.com/Aneg/otus-anti-brute-force/internal/config"
	grpc2 "github.com/Aneg/otus-anti-brute-force/internal/web/grpc"
	"github.com/Aneg/otus-anti-brute-force/pkg/api"
	log2 "github.com/Aneg/otus-anti-brute-force/pkg/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

func main() {
	conf, err := config.GetConfigFromFile(getConfigDir())
	if err != nil {
		log.Fatal(err)
	}
	log2.Logger, err = getLogger(conf.LogFile, conf.LogLevel)
	if err != nil {
		log.Fatal(err)
	}

	lis, err := net.Listen("tcp", conf.HttpListen)
	if err != nil {
		log2.Logger.Fatal(err.Error())
	}
	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)

	server := grpc2.NewServer()
	api.RegisterAntiBruteForceServer(grpcServer, server)
	log2.Logger.Info("Начинаем слушать...")
	if err := grpcServer.Serve(lis); err != nil {
		log2.Logger.Fatal(err.Error())
	}
}

func getConfigDir() string {
	var configDir string
	flag.StringVar(&configDir, "config", "configs/config.yaml", "path to config file")
	flag.Parse()
	return configDir
}

func getLogger(logFile, logLevel string) (logger *zap.Logger, err error) {
	var level zapcore.Level
	switch logLevel {
	case "debug":
		level = zapcore.DebugLevel
		break
	case "info":
		level = zapcore.InfoLevel
		break
	case "warn":
		level = zapcore.WarnLevel
		break
	case "error":
		level = zapcore.ErrorLevel
	}

	return zap.Config{
		Encoding:    "json",
		Level:       zap.NewAtomicLevelAt(level),
		OutputPaths: []string{"stdout", logFile},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey: "message", // <--
		},
	}.Build()
}

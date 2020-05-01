package main

import (
	"flag"
	"github.com/Aneg/otus-anti-brute-force/internal/config"
	"github.com/Aneg/otus-anti-brute-force/internal/repositories/mysql"
	"github.com/Aneg/otus-anti-brute-force/internal/services"
	"github.com/Aneg/otus-anti-brute-force/internal/services/ip_guard"
	"github.com/Aneg/otus-anti-brute-force/internal/services/worker"
	grpc2 "github.com/Aneg/otus-anti-brute-force/internal/web/grpc"
	"github.com/Aneg/otus-anti-brute-force/pkg/api"
	"github.com/Aneg/otus-anti-brute-force/pkg/database"
	log2 "github.com/Aneg/otus-anti-brute-force/pkg/log"
	worker2 "github.com/Aneg/otus-anti-brute-force/pkg/worker"
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
	closeChan := make(chan bool)

	dbConn, err := database.MysqlOpenConnection(conf.DBUser, conf.DBPass, conf.DBHostPort, conf.DBName)
	if err != nil {
		log2.Logger.Fatal(err.Error())
	}

	maskRepository := mysql.NewMasksRepository(dbConn)
	whiteList := ip_guard.NewMemoryIpGuard(1, maskRepository)
	blackList := ip_guard.NewMemoryIpGuard(2, maskRepository)
	errorWorkerChan := make(chan error, 100)
	reloaderMasksWorker := worker.NewReloaderMasks([]services.IpGuard{whiteList, blackList}, errorWorkerChan)

	go func(errorChan chan error, closeChan chan bool) {
		for true {
			select {
			case err := <-errorChan:
				log2.Logger.Error(err.Error())
			case <-closeChan:
				return
			}
		}
	}(errorWorkerChan, closeChan)

	worker2.Start(reloaderMasksWorker)

	lis, err := net.Listen("tcp", conf.HttpListen)
	if err != nil {
		log2.Logger.Fatal(err.Error())
	}
	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)

	server := grpc2.NewServer(whiteList, blackList)
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

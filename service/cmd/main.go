package main

import (
	"context"
	"fmt"
	"github.com/barcodepro/go-service-template/service/app"
	"github.com/barcodepro/go-service-template/service/internal/log"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
	"os/signal"
	"syscall"
)

var (
	appName, gitCommit, gitBranch string
)

func main() {
	var (
		showVersion        = kingpin.Flag("version", "show version and exit").Default().Bool()
		logLevel           = kingpin.Flag("log-level", "set log level: debug, info, warn, error").Default("warn").Envar("LOG_LEVEL").String()
		listenAddress      = kingpin.Flag("listen-address", "Address to listen on for metrics").Default(":1080").Envar("SERVER_LISTEN_ADDRESS").TCP()
		corsAllowedOrigins = kingpin.Flag("cors-allowed-origins", "allowed origins for CORS").Default("http://127.0.0.1:*,http://localhost:*").Envar("SERVER_CORS_ALLOWED_ORIGINS").String()
		postgresURL        = kingpin.Flag("postgres-url", "postgresql url").Default("").Envar("POSTGRES_URL").String()
	)
	kingpin.Parse()
	log.SetLevel(*logLevel)
	log.SetApplication(appName)

	var config = &app.Config{
		ListenAddress:  **listenAddress,
		PostgresURL:    *postgresURL,
		AllowedOrigins: *corsAllowedOrigins,
	}

	if *showVersion {
		fmt.Printf("%s %s-%s\n", appName, gitCommit, gitBranch)
		os.Exit(0)
	}

	if err := config.Validate(); err != nil {
		log.Errorf("Cannot start %s, unable to validate config: %s", appName, err)
		os.Exit(1)
	}

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	var doExit = make(chan error, 2)
	go func() {
		doExit <- listenSignals()
		cancel()
	}()

	go func() {
		doExit <- app.Start(ctx, config)
		cancel()
	}()

	log.Warnf("shutdown: %s", <-doExit)
}

func listenSignals() error {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	return fmt.Errorf("got %s", <-c)
}

package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

const version = "1.0.0"

type application struct {
	config config
	logger *log.Logger
}

type config struct {
	port int
	env  string
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 7777, "api port")
	flag.StringVar(&cfg.env, "env", "dev", "environment (dev|staging|prod)")
	flag.Parse()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	app := &application{
		config: cfg,
		logger: logger,
	}

	server := &http.Server{
		Handler: app.Routes(),
		Addr:    fmt.Sprintf(":%d", app.config.port),
	}

	logger.Printf("%s server started at port %d", cfg.env, cfg.port)
	err := server.ListenAndServe()
	log.Fatal(err)
}

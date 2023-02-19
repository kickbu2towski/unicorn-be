package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/kickbu2towski/unicorn-be/cmd/api/internal/data"
	_ "github.com/lib/pq"
)

const version = "1.0.0"

type application struct {
	config config
	logger *log.Logger
	models data.Models
}

type config struct {
	port int
	env  string
	db   struct {
		dsn string
	}
}

func main() {
	var cfg config
	flag.IntVar(&cfg.port, "port", 7777, "API port")
	flag.StringVar(&cfg.env, "env", "dev", "Environment (dev|staging|prod)")
	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("UNICORN_DB_DSN"), "PosgreSQL DSN")
	flag.Parse()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	db, err := openDB(cfg)
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	logger.Print("established connection to database")

	app := &application{
		config: cfg,
		logger: logger,
		models: data.NewModels(db),
	}

	server := &http.Server{
		Handler: app.Routes(),
		Addr:    fmt.Sprintf(":%d", app.config.port),
	}

	logger.Printf("%s server started at port %d", cfg.env, cfg.port)
	err = server.ListenAndServe()
	logger.Fatal(err)
}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}

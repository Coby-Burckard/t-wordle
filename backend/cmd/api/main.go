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

	"github.com/coby-burckard/t-wordle/backend/internal/data"
	"github.com/coby-burckard/t-wordle/backend/internal/jsonlog"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

const version = "1.0.0"

type config struct {
	port int
	env string
	db struct {
		dsn string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime string
	}
}

type application struct {
	config config
	log *jsonlog.Logger
	models data.Models
}

func main() {
	var cfg config

	// initialize logger
	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	// load environment variables
	err := godotenv.Load("development.env")
	if err != nil {
		logger.PrintFatal(err, nil)
	}

	// read command line flags
	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Envirionment (development|staging|prodution)")
	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("WORDLE_DB_DSN"), "PostgreSQL DSN")
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.StringVar(&cfg.db.maxIdleTime, "db-max-idle-time", "15m", "PostgreSQL max connection idle time")
	flag.Parse()

	// test db connection
	db, err := openDB(cfg)
	if err != nil {
		logger.PrintFatal(err, nil)
	}
	defer db.Close()
	logger.PrintInfo("Database connection pool established", nil)

	// application struct
	app := &application{
		config: cfg,
		log: logger,
		models: data.NewModels(db),
	}
	
	// http server
	srv := &http.Server{
		Addr: fmt.Sprintf(":%d", cfg.port),
		Handler: app.routes(),
		ErrorLog: log.New(logger, "", 0),
		IdleTimeout: time.Minute,
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 20 * time.Second,
	}

	// start Http server
	logger.PrintInfo("Starting server", map[string]string{
		"address": srv.Addr,
		"environment": cfg.env,
	})
	err = srv.ListenAndServe()
	logger.PrintFatal(err, nil)
}

// openDB returns a sql.DB connection pool
func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	// configure db settings
	db.SetMaxOpenConns(cfg.db.maxOpenConns)
	db.SetMaxIdleConns(cfg.db.maxIdleConns)
	duration, err := time.ParseDuration(cfg.db.maxIdleTime)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(duration)

	// ensure a connection is established within 5 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
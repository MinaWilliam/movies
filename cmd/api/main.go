package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/MinaWilliam/movies/internal/data"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

const version string = "1.0.0"

type config struct {
	port string
	env  string
	db   struct {
		driver       string
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  string
	}
}

type app struct {
	config config
	logger *log.Logger
	models data.Models
}

func init() {
	initEnv()
}

func main() {
	startServer()
}

func initEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("error loading .env file")
	}
}

func startServer() {
	var config config

	config.db.driver = os.Getenv("DB_DRIVER")
	config.db.dsn = os.Getenv("DB_DSN")
	config.db.maxOpenConns, _ = strconv.Atoi(os.Getenv("DB_MAX_OPEN_CONNS"))
	config.db.maxIdleConns, _ = strconv.Atoi(os.Getenv("DB_MAX_IDLE_CONNS"))
	config.db.maxIdleTime = os.Getenv("DB_MAX_IDLE_TIME")

	flag.StringVar(&config.port, "port", os.Getenv("PORT"), "API server port")
	flag.StringVar(&config.env, "env", os.Getenv("ENV"), "(development|staging|production)")
	flag.Parse()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	db, err := openDB(config)
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	logger.Printf("database connection pool established")

	app := &app{
		config: config,
		logger: logger,
		models: data.NewModels(db),
	}

	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", config.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	logger.Printf("starting %s server on %s", config.env, server.Addr)
	err = server.ListenAndServe()
	logger.Fatal(err)
}

func openDB(config config) (*sql.DB, error) {
	db, err := sql.Open(config.db.driver, config.db.dsn)
	if err != nil {
		return nil, err
	}

	duration, err := time.ParseDuration(config.db.maxIdleTime)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(duration)
	db.SetMaxOpenConns(config.db.maxOpenConns)
	db.SetMaxIdleConns(config.db.maxIdleConns)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}

package main

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"moleben/assembly"
	"moleben/config"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib" // stdlib driver for goose

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pressly/goose/v3"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	// 1) Run goose migrations via database/sql (pgx stdlib)
	sqlDB, err := sql.Open("pgx", cfg.DB_DSN)
	if err != nil {
		log.Fatalf("sql open: %v", err)
	}
	defer sqlDB.Close()
	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatalf("goose dialect: %v", err)
	}

	ctx := context.Background()
	if err := goose.UpContext(ctx, sqlDB, cfg.MigrationsDir); err != nil {
		log.Fatalf("goose up: %v", err)
	}

	// 2) App pool for runtime
	pool, err := pgxpool.New(ctx, cfg.DB_DSN)
	if err != nil {
		log.Fatalf("pgxpool: %v", err)
	}
	defer pool.Close()

	app := assembly.Build(cfg, pool)
	srv := &http.Server{Addr: cfg.HTTPAddr, Handler: app.Router.Engine}

	go func() {
		log.Printf("listening on %s", cfg.HTTPAddr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	log.Println("shutting down...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = srv.Shutdown(shutdownCtx)
}

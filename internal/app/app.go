package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/indigowar/map-of-events/internal/config"
)

func Run(cfg *config.Config) {
	r := gin.New()

	url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.Name)

	p, err := pgxpool.New(context.Background(), url)

	if err != nil {
		log.Fatalln(err)
	}
	defer p.Close()

	server := &http.Server{
		Handler:        r,
		Addr:           ":" + cfg.HTTP.Port,
		WriteTimeout:   cfg.HTTP.WriteTimeout,
		ReadTimeout:    cfg.HTTP.ReadTimeout,
		MaxHeaderBytes: cfg.HTTP.MaxHeadersMegabytes << 20,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Fatalln(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	<-quit
	log.Println("Start the shutdown...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalln(err)
	}

	log.Println("Application has been stopped.")
}

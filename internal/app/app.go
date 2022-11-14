package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"

	"github.com/indigowar/map-of-events/internal/config"
	"github.com/indigowar/map-of-events/internal/delivery/http/json"
	"github.com/indigowar/map-of-events/internal/infra/adapters/storages"
	"github.com/indigowar/map-of-events/internal/services"
)

func Run(cfg *config.Config) {
	url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.Name)
	log.Println(url)
	conn, err := pgx.Connect(context.Background(), url)

	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close(context.Background())

	competitorStorage, err := storages.NewPostgresCompetitorStorage(conn)
	if err != nil {
		log.Fatalln(err)
	}
	organizerStorage, _ := storages.NewPostgresOrganizerStorage(conn)

	foundingRangeStorage := storages.NewFoundingRangePostgresStorage(conn)
	coFoundingRangeStorage := storages.NewCoFoundingRangePostgresStorage(conn)

	competitorService := services.NewCompetitorService(competitorStorage)
	organizerService, _ := services.NewOrganizerService(organizerStorage)
	foundingService := services.NewFoundingRangeService(foundingRangeStorage)
	coFoundingService := services.NewCoFoundingRangeService(coFoundingRangeStorage)

	r := gin.Default()

	r.Use(cors.Default())

	r.GET("api/v1/competitor", json.GetAllCompetitorsHandler(competitorService))
	r.POST("api/v1/competitor", json.CreateCompetitorHandler(competitorService))

	r.GET("/api/v1/founding_range/:id", json.GetByIDRangeHandler(foundingService))
	r.GET("/api/v1/founding_range", json.GetMaximumRangeHandler(foundingService))

	r.GET("/api/v1/co_founding_range/:id", json.GetByIDRangeHandler(coFoundingService))
	r.GET("/api/v1/co_founding_range", json.GetMaximumRangeHandler(coFoundingService))

	r.GET("api/v1/organizer_level", json.GetAllOrganizerLevelsHandler(organizerService))
	r.POST("api/v1/organizer_level", json.CreateOrganizerLevelHandler(organizerService))

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

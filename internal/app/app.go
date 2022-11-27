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
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/indigowar/map-of-events/internal/config"
	"github.com/indigowar/map-of-events/internal/infra/ports/delivery/http/v1/files"
	"github.com/indigowar/map-of-events/internal/infra/ports/delivery/http/v1/json"
	json2 "github.com/indigowar/map-of-events/internal/infra/ports/delivery/http/v2/json"
)

func Run(cfg *config.Config) {
	url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.Name)
	log.Println(url)

	postgresCPool, err := pgxpool.New(context.Background(), url)
	if err != nil {
		log.Fatalln(err)
	}
	defer postgresCPool.Close()

	services := initServices(postgresCPool)

	eventHandler := json.NewEventHandler(services)

	r := gin.Default()

	r.Use(cors.Default())

	v1 := r.Group("api/v1")
	{
		v1.GET("/competitor", json.GetAllCompetitorsHandler(services.Competitor))
		v1.POST("/competitor", json.CreateCompetitorHandler(services.Competitor))

		v1.GET("/founding_range/:id", json.GetByIDRangeHandler(services.FoundingRange))
		v1.GET("/founding_range", json.GetMaximumRangeHandler(services.FoundingRange))

		v1.GET("/co_founding_range/:id", json.GetByIDRangeHandler(services.CoFoundingRange))
		v1.GET("/co_founding_range", json.GetMaximumRangeHandler(services.CoFoundingRange))

		v1.GET("/organizer_level", json.GetAllOrganizerLevelsHandler(services.Organizer))
		v1.POST("/organizer_level", json.CreateOrganizerLevelHandler(services.Organizer))

		v1.GET("/organizer", json.GetAllOrganizersHandler(services.Organizer))
		v1.POST("/organizer", json.CreateOrganizerHandler(services.Organizer))
		v1.GET("/organizer/:id", json.GetByIDOrganizerHandler(services.Organizer))
		v1.PUT("/organizer/:id", json.UpdateOrganizerHandler(services.Organizer))
		v1.DELETE("/organizer/:id", json.DeleteOrganizerHandler(services.Organizer))

		v1.GET("/event", eventHandler.GetAllEvents)
		v1.POST("/event", eventHandler.Create)

		v1.GET("/event/:id", eventHandler.GetEventByID)
		v1.DELETE("/event/:id", eventHandler.DeleteEvent)
		v1.PUT("/event/:id", eventHandler.DeleteEvent)

		v1.GET("/minimal_event", eventHandler.GetAllAsMinimal)
		v1.GET("/minimal_event/:id", eventHandler.GetByIDMinimal)

		v1.POST("/image", files.UploadHandler(services.Image))
		v1.GET("/image/:link", files.RetrievingHandler(services.Image))
	}

	v2 := r.Group("/api/v2")
	{
		v2.GET("/organizer", json2.GetAllOrganizersHandler(services.Organizer, services.Image))
		v2.POST("/organizer", json2.CreateOrganizerHandler(services.Organizer, services.Image))
		v2.GET("/organizer/:id", json2.GetOrganizerByID(services.Organizer, services.Image))
	}

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

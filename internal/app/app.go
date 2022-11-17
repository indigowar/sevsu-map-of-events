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
	"github.com/indigowar/map-of-events/internal/infra/adapters/storages"
	json2 "github.com/indigowar/map-of-events/internal/infra/ports/delivery/http/json"
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
	subjectStorage := storages.NewPostgresSubjectStorage(conn)
	eventStorage := storages.NewPostgresEventStorage(conn, subjectStorage)

	competitorService := services.NewCompetitorService(competitorStorage)
	organizerService, _ := services.NewOrganizerService(organizerStorage)
	foundingService := services.NewFoundingRangeService(foundingRangeStorage)
	coFoundingService := services.NewCoFoundingRangeService(coFoundingRangeStorage)
	eventService := services.NewEventServices(eventStorage, subjectStorage, organizerService, foundingService, coFoundingService, competitorService)

	r := gin.Default()

	r.Use(cors.Default())

	v1 := r.Group("api/v1")
	{
		v1.GET("/competitor", json2.GetAllCompetitorsHandler(competitorService))
		v1.POST("/competitor", json2.CreateCompetitorHandler(competitorService))

		v1.GET("/founding_range/:id", json2.GetByIDRangeHandler(foundingService))
		v1.GET("/founding_range", json2.GetMaximumRangeHandler(foundingService))

		v1.GET("/co_founding_range/:id", json2.GetByIDRangeHandler(coFoundingService))
		v1.GET("/co_founding_range", json2.GetMaximumRangeHandler(coFoundingService))

		v1.GET("/organizer_level", json2.GetAllOrganizerLevelsHandler(organizerService))
		v1.POST("/organizer_level", json2.CreateOrganizerLevelHandler(organizerService))

		v1.GET("/organizer/", json2.GetAllOrganizersHandler(organizerService))
		v1.POST("/organizer/", json2.CreateOrganizerHandler(organizerService))
		v1.GET("/organizer/:id", json2.GetByIDOrganizerHandler(organizerService))
		v1.POST("/organizer/:id", json2.UpdateOrganizerHandler(organizerService)) // TODO: implement

		v1.GET("/event/", json2.GetAllEventHandler(eventService))
		v1.POST("/event/", json2.CreateEventHandler(eventService))

		v1.GET("/event/:id/", json2.GetByIDEventHandler(eventService))
		v1.DELETE("/event/:id", json2.DeleteEventHandler(eventService))
		v1.POST("/event/:id", json2.UpdateEventHandler(eventService)) // TODO: implement

		v1.GET("/minimal_event/", json2.GetAllAsMinimalHandler(eventService))
		v1.GET("/minimal_event/:id", json2.GetByIDAsMinimalHandler(eventService))
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

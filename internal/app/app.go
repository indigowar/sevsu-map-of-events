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
	"github.com/indigowar/map-of-events/internal/infra/adapters/storages"
	"github.com/indigowar/map-of-events/internal/infra/ports/delivery/http/v1/files"
	"github.com/indigowar/map-of-events/internal/infra/ports/delivery/http/v1/json"
	json2 "github.com/indigowar/map-of-events/internal/infra/ports/delivery/http/v2/json"
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

	postgresCPool, err := pgxpool.New(context.Background(), url)
	if err != nil {
		log.Fatalln(err)
	}
	defer postgresCPool.Close()

	competitorStorage := storages.NewPostgresCompetitorStorage(postgresCPool)
	organizerStorage := storages.NewPostgresOrganizerStorage(postgresCPool)
	foundingRangeStorage := storages.NewFoundingRangePostgresStorage(postgresCPool)
	coFoundingRangeStorage := storages.NewCoFoundingRangePostgresStorage(postgresCPool)
	subjectStorage := storages.NewPostgresSubjectStorage(postgresCPool)
	eventStorage := storages.NewPostgresEventStorage(postgresCPool)
	imageStorage := storages.NewPostgresImageStorage(postgresCPool)

	competitorService := services.NewCompetitorService(competitorStorage)
	organizerService, _ := services.NewOrganizerService(organizerStorage)
	foundingService := services.NewFoundingRangeService(foundingRangeStorage)
	coFoundingService := services.NewCoFoundingRangeService(coFoundingRangeStorage)
	subjectService := services.NewSubjectService(subjectStorage)
	eventService := services.NewEventServices(eventStorage, subjectService, organizerService, foundingService, coFoundingService, competitorService)

	imageService := services.NewImageService(imageStorage)

	r := gin.Default()

	r.Use(cors.Default())

	v1 := r.Group("api/v1")
	{
		v1.GET("/competitor", json.GetAllCompetitorsHandler(competitorService))
		v1.POST("/competitor", json.CreateCompetitorHandler(competitorService))

		v1.GET("/founding_range/:id", json.GetByIDRangeHandler(foundingService))
		v1.GET("/founding_range", json.GetMaximumRangeHandler(foundingService))

		v1.GET("/co_founding_range/:id", json.GetByIDRangeHandler(coFoundingService))
		v1.GET("/co_founding_range", json.GetMaximumRangeHandler(coFoundingService))

		v1.GET("/organizer_level", json.GetAllOrganizerLevelsHandler(organizerService))
		v1.POST("/organizer_level", json.CreateOrganizerLevelHandler(organizerService))

		v1.GET("/organizer/", json.GetAllOrganizersHandler(organizerService))
		v1.POST("/organizer/", json.CreateOrganizerHandler(organizerService))
		v1.GET("/organizer/:id", json.GetByIDOrganizerHandler(organizerService))
		v1.POST("/organizer/:id", json.UpdateOrganizerHandler(organizerService)) // TODO: implement

		v1.GET("/event/", json.GetAllEventHandler(eventService))
		v1.POST("/event/", json.CreateEventHandler(eventService))

		v1.GET("/event/:id/", json.GetByIDEventHandler(eventService))
		v1.DELETE("/event/:id", json.DeleteEventHandler(eventService))
		v1.POST("/event/:id", json.UpdateEventHandler(eventService)) // TODO: implement

		v1.GET("/minimal_event/", json.GetAllAsMinimalHandler(eventService))
		v1.GET("/minimal_event/:id", json.GetByIDAsMinimalHandler(eventService))

		v1.POST("/image/", files.UploadHandler(imageService))
		v1.GET("/image/:link", files.RetrievingHandler(imageService))
	}

	v2 := r.Group("/api/v2")
	{
		v2.GET("/organizer", json2.GetAllOrganizersHandler(organizerService, imageService))
		v2.POST("/organizer", json2.CreateOrganizerHandler(organizerService, imageService))
		v2.GET("/organizer/:id", json2.GetOrganizerByID(organizerService, imageService))
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

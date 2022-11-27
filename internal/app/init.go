package app

import (
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/indigowar/map-of-events/internal/domain/services"
	"github.com/indigowar/map-of-events/internal/infra/adapters/storages/postgres"
	svc "github.com/indigowar/map-of-events/internal/services"
)

func initServices(pool *pgxpool.Pool) services.Services {
	competitorStorage := postgres.NewPostgresCompetitorStorage(pool)
	organizerStorage := postgres.NewPostgresOrganizerStorage(pool)
	foundingRangeStorage := postgres.NewFoundingRangePostgresStorage(pool)
	coFoundingRangeStorage := postgres.NewCoFoundingRangePostgresStorage(pool)
	subjectStorage := postgres.NewPostgresSubjectStorage(pool)
	eventStorage := postgres.NewPostgresEventStorage(pool)
	imageStorage := postgres.NewPostgresImageStorage(pool)

	var s services.Services

	s.Subject = svc.NewSubjectService(subjectStorage)
	s.Image = svc.NewImageService(imageStorage)
	s.Organizer, _ = svc.NewOrganizerService(organizerStorage, s.Image)
	s.FoundingRange = svc.NewFoundingRangeService(foundingRangeStorage)
	s.CoFoundingRange = svc.NewCoFoundingRangeService(coFoundingRangeStorage)
	s.Competitor = svc.NewCompetitorService(competitorStorage)
	s.Event = svc.NewEventServices(eventStorage, s.Subject, s.Organizer, s.FoundingRange, s.CoFoundingRange, s.Competitor)

	return s
}

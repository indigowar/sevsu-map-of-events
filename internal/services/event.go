package services

import (
	"context"
	"errors"
	"log"

	"github.com/google/uuid"

	"github.com/indigowar/map-of-events/internal/domain/models"
	"github.com/indigowar/map-of-events/internal/domain/repos/adapters/storages"
	"github.com/indigowar/map-of-events/internal/domain/services"
)

type eventService struct {
	organizer        services.OrganizerService
	foundingRanges   services.RangeService
	coFoundingRanges services.RangeService
	competitors      services.CompetitorService

	eventStorage   storages.EventStorageRepository
	subjectStorage storages.SubjectStorageRepository
}

func (svc eventService) GetAll(ctx context.Context) ([]models.Event, error) {
	ids, err := svc.eventStorage.GetIDList(ctx)
	if err != nil {
		log.Println(err)
		return nil, errors.New("internal error")
	}
	events := make([]models.Event, len(ids))
	for i := range ids {
		events[i], err = svc.GetByID(ctx, ids[i])
		if err != nil {
			log.Println(err)
		}
	}
	return events, nil
}

func (svc eventService) GetByID(ctx context.Context, id uuid.UUID) (models.Event, error) {
	event, err := svc.eventStorage.GetByID(ctx, id)
	if err != nil {
		log.Println(err)
		return models.Event{}, errors.New("not found")
	}
	return event, nil
}

func (svc eventService) Create(ctx context.Context, info services.EventCreateUpdateInfo) (models.Event, error) {
	eventId := uuid.New()

	var undoError error = nil

	if _, err := svc.organizer.GetByID(ctx, info.Organizer); err != nil {
		return models.Event{}, nil
	}

	founding, err := svc.foundingRanges.Create(ctx, info.FoundingRangeLow, info.FoundingRangeHigh)
	if err != nil {
		return models.Event{}, errors.New("internal error")
	}
	defer func() {
		if undoError != nil {
			_ = svc.foundingRanges.Delete(context.TODO(), founding.ID)
		}
	}()

	coFounding, err := svc.coFoundingRanges.Create(ctx, info.CoFoundingRangeLow, info.CoFoundingRangeHigh)
	if err != nil {
		undoError = errors.New("range_svc to event_svc: " + err.Error())
		return models.Event{}, errors.New("internal error")
	}

	subjects := make([]uuid.UUID, 0)

	event := models.Event{
		ID:                  eventId,
		Title:               info.Title,
		Organizer:           info.Organizer,
		FoundingType:        info.FoundingType,
		FoundingRange:       founding.ID,
		CoFoundingRange:     coFounding.ID,
		SubmissionDeadline:  info.SubmissionDeadline,
		ConsiderationPeriod: info.ConsiderationPeriod,
		RealisationPeriod:   info.RealisationPeriod,
		Result:              info.Result,
		Site:                info.Site,
		Document:            info.Document,
		InternalContacts:    info.InternalContacts,
		TRL:                 info.TRL,
		Competitors:         subjects,
	}

	if err := svc.eventStorage.Create(ctx, event); err != nil {
		log.Println(err)
		undoError = errors.New("failed to create event")
		return models.Event{}, errors.New("internal error")
	}
	defer func() {
		if undoError != nil {
			_ = svc.Delete(context.TODO(), eventId)
		}
	}()

	for _, v := range info.Subjects {
		id := uuid.New()
		_ = svc.subjectStorage.Add(ctx, models.Subject{
			ID:   id,
			Name: v,
		})
		subjects = append(subjects, id)
	}

	return event, nil
}

func (svc eventService) Delete(ctx context.Context, id uuid.UUID) error {
	return svc.Delete(ctx, id)
}

func (svc eventService) GetAllAsMinimal(ctx context.Context) ([]services.EventMinimal, error) {
	events, err := svc.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]services.EventMinimal, len(events))
	for i, v := range events {
		result[i] = services.EventMinimal{
			ID:                 v.ID,
			Title:              v.Title,
			Organizer:          v.Organizer,
			SubmissionDeadline: v.SubmissionDeadline,
			TRL:                v.TRL,
		}
	}
	return result, nil
}

func (svc eventService) GetByIDAsMinimal(ctx context.Context, id uuid.UUID) (services.EventMinimal, error) {
	event, err := svc.GetByID(ctx, id)
	if err != nil {
		return services.EventMinimal{}, err
	}
	return services.EventMinimal{
		ID:                 event.ID,
		Title:              event.Title,
		Organizer:          event.Organizer,
		SubmissionDeadline: event.SubmissionDeadline,
		TRL:                event.TRL,
	}, nil
}

func (svc eventService) Update(_ context.Context, _ uuid.UUID, _ services.EventCreateUpdateInfo) (models.Event, error) {
	panic("not implemented")
}

func NewEventServices(storage storages.EventStorageRepository, subjects storages.SubjectStorageRepository,
	organizer services.OrganizerService,
	foundingRange, coFoundingRange services.RangeService,
	competitors services.CompetitorService) services.EventService {
	return &eventService{
		organizer:        organizer,
		foundingRanges:   foundingRange,
		coFoundingRanges: coFoundingRange,
		competitors:      competitors,
		eventStorage:     storage,
		subjectStorage:   subjects,
	}
}

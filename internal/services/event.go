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
	subjects         services.SubjectService

	eventStorage storages.EventStorageRepository
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

func (svc eventService) Create(ctx context.Context, info services.EventCreateInfo) (models.Event, error) {
	tx, err := svc.eventStorage.InvokeTransactionMechanism(ctx)
	if err != nil {
		log.Println(err)
	}
	defer func(eventStorage storages.EventStorageRepository, ctx context.Context, transaction interface{}) {
		_ = eventStorage.ShadowTransactionMechanism(ctx, transaction)
	}(svc.eventStorage, ctx, tx)

	eventId := uuid.New()

	serviceCtx := context.WithValue(ctx, "connection", tx)

	if _, err := svc.organizer.GetByID(serviceCtx, info.Organizer); err != nil {
		log.Println(err)
		return models.Event{}, errors.New("not found - organizer")
	}

	founding, err := svc.foundingRanges.Create(serviceCtx, info.FoundingRangeLow, info.FoundingRangeHigh)
	if err != nil {
		log.Println(err)
		return models.Event{}, errors.New("failed to create - founding range")
	}

	coFounding, err := svc.coFoundingRanges.Create(serviceCtx, info.CoFoundingRangeLow, info.CoFoundingRangeHigh)
	if err != nil {
		return models.Event{}, errors.New("failed to create - co founding range")
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
		Competitors:         info.Competitors,
	}
	if err := svc.eventStorage.Create(serviceCtx, event); err != nil {
		log.Println(err)
		return models.Event{}, errors.New("failed to create - event")
	}

	for _, v := range info.Subjects {
		s, _ := svc.subjects.Create(serviceCtx, eventId, v)
		subjects = append(subjects, s.ID)
	}

	return event, nil
}

func (svc eventService) Delete(ctx context.Context, id uuid.UUID) error {
	transaction, err := svc.eventStorage.InvokeTransactionMechanism(ctx)
	if err != nil {
		log.Println(err)
		return errors.New("failed to star a transaction")
	}
	defer func() {
		_ = svc.eventStorage.ShadowTransactionMechanism(ctx, transaction)
	}()
	serviceCtx := context.WithValue(ctx, "connection", transaction)

	event, err := svc.GetByID(serviceCtx, id)
	if err != nil {
		log.Println(err)
		return errors.New("not found - event")
	}

	if svc.foundingRanges.Delete(serviceCtx, event.FoundingRange) != nil {
		return errors.New("failed to delete - founding range")
	}

	if svc.coFoundingRanges.Delete(serviceCtx, event.CoFoundingRange) != nil {
		return errors.New("failed to delete - co founding range")
	}

	subjects, err := svc.subjects.GetAllForEvent(serviceCtx, event.ID)
	if err != nil {
		log.Println(err)
		return errors.New("failed to get - subjects")
	}

	for _, v := range subjects {
		err := svc.subjects.Delete(serviceCtx, v.ID)
		if err != nil {
			log.Println(err)
			return errors.New("failed to delete subject")
		}
	}

	err = svc.eventStorage.Delete(serviceCtx, event.ID)
	if err != nil {
		log.Println(err)
		return errors.New("failed to delete event")
	}
	return nil
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

func (svc eventService) Update(_ context.Context, _ uuid.UUID, _ services.EventCreateInfo) (models.Event, error) {
	panic("not implemented")
}

func NewEventServices(storage storages.EventStorageRepository,
	subjects services.SubjectService,
	organizer services.OrganizerService,
	foundingRange, coFoundingRange services.RangeService,
	competitors services.CompetitorService) services.EventService {
	return &eventService{
		organizer:        organizer,
		foundingRanges:   foundingRange,
		coFoundingRanges: coFoundingRange,
		competitors:      competitors,
		eventStorage:     storage,
		subjects:         subjects,
	}
}

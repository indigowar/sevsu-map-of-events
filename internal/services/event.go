package services

import (
	"context"
	"errors"
	"log"

	"github.com/google/uuid"

	"github.com/indigowar/map-of-events/internal/domain/models"
	"github.com/indigowar/map-of-events/internal/domain/repos/adapters/storages"
	"github.com/indigowar/map-of-events/internal/domain/services"
	"github.com/indigowar/map-of-events/internal/domain/validators"
)

type eventService struct {
	organizer        services.OrganizerService
	foundingRanges   services.RangeService
	coFoundingRanges services.RangeService
	competitors      services.CompetitorService
	subjects         services.SubjectService

	eventStorage storages.EventStorage
}

func (svc eventService) AllIDs(ctx context.Context) ([]uuid.UUID, error) {
	return svc.eventStorage.GetIDList(ctx)
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

func (svc eventService) validateCreationInfo(ctx context.Context, info services.EventCreateInfo) error {
	// The organizer should already exist
	// in the moment of creation it's event
	{
		existedOrganizers, err := svc.organizer.GetAllIDs(ctx)
		if err != nil {
			log.Println(err)
			return errors.New("internal error")
		}
		if !validators.IDExists(existedOrganizers, info.Organizer) {
			return errors.New("organizer does not exist")
		}
	}

	// Validate that both Founding Range and CoFoundingRange are valid
	{
		if err := validators.ValidateRange(models.RangeModel{Low: info.FoundingRangeLow, High: info.FoundingRangeHigh}); err != nil {
			return err
		}

		if err := validators.ValidatePercentRange(models.RangeModel{Low: info.CoFoundingRangeLow, High: info.CoFoundingRangeHigh}); err != nil {
			return err
		}
	}

	// Validate TRL
	// TRL should be between 0 and 9
	{
		if info.TRL <= 0 || info.TRL >= 10 {
			return errors.New("invalid trl")
		}
	}

	// All of used in this event
	// should already exist
	{
		competitors, err := svc.competitors.AllIDs(ctx)
		if err != nil {
			log.Println(err)
			return errors.New("internal error")
		}

		for _, competitor := range info.Competitors {
			if !validators.IDExists(competitors, competitor) {
				return errors.New("competitor does not exist")
			}
		}
	}

	return nil
}

func (svc eventService) Create(ctx context.Context, info services.EventCreateInfo) (models.Event, error) {
	if err := svc.validateCreationInfo(ctx, info); err != nil {
		return models.Event{}, err
	}

	tx, err := svc.eventStorage.BeginTransaction(ctx)
	if err != nil {
		log.Println(err)
	}
	defer func(eventStorage storages.EventStorage, ctx context.Context, transaction interface{}) {
		_ = eventStorage.CloseTransaction(ctx, transaction)
	}(svc.eventStorage, ctx, tx)

	serviceCtx := context.WithValue(ctx, "connection", tx)

	founding, err := svc.foundingRanges.Create(serviceCtx, info.FoundingRangeLow, info.FoundingRangeHigh)
	if err != nil {
		log.Println(err)
		return models.Event{}, errors.New("failed to create - founding range")
	}

	coFounding, err := svc.coFoundingRanges.Create(serviceCtx, info.CoFoundingRangeLow, info.CoFoundingRangeHigh)
	if err != nil {
		return models.Event{}, errors.New("failed to create - co founding range")
	}

	eventId := uuid.New()
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

	if err := svc.eventStorage.Add(serviceCtx, event); err != nil {
		log.Println(err)
		return models.Event{}, errors.New("failed to create - event")
	}

	subjects := make([]uuid.UUID, len(info.Subjects))
	for i, v := range info.Subjects {
		s, _ := svc.subjects.Create(serviceCtx, eventId, v)
		subjects[i] = s.ID
	}
	return event, nil
}

func (svc eventService) Delete(ctx context.Context, id uuid.UUID) error {
	transaction, err := svc.eventStorage.BeginTransaction(ctx)
	if err != nil {
		log.Println(err)
		return errors.New("failed to star a transaction")
	}
	defer func() {
		_ = svc.eventStorage.CloseTransaction(ctx, transaction)
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

	err = svc.eventStorage.Remove(serviceCtx, event.ID)
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

func (svc eventService) updateEventModel(e models.Event, i services.EventCreateInfo) models.Event {
	e.Title = i.Title
	e.Organizer = i.Organizer
	e.SubmissionDeadline = i.SubmissionDeadline
	e.ConsiderationPeriod = i.ConsiderationPeriod
	e.RealisationPeriod = i.RealisationPeriod
	e.Result = i.Result
	e.Site = i.Site
	e.Document = i.Document
	e.InternalContacts = i.InternalContacts
	e.TRL = i.TRL
	e.Competitors = i.Competitors
	return e
}

func (svc eventService) Update(ctx context.Context, id uuid.UUID, info services.EventCreateInfo) (models.Event, error) {
	storedEvent, err := svc.GetByID(ctx, id)
	if err != nil {
		log.Println(err)
		return models.Event{}, errors.New("event was not found")
	}

	if err := svc.validateCreationInfo(ctx, info); err != nil {
		log.Println(err)
		return models.Event{}, err
	}

	tx, err := svc.eventStorage.BeginTransaction(ctx)
	if err != nil {
		log.Println(err)
	}
	defer func(eventStorage storages.EventStorage, ctx context.Context, transaction interface{}) {
		_ = eventStorage.CloseTransaction(ctx, transaction)
	}(svc.eventStorage, ctx, tx)

	serviceCtx := context.WithValue(ctx, "connection", tx)

	storedEvent = svc.updateEventModel(storedEvent, info)

	if err := svc.eventStorage.Update(serviceCtx, storedEvent); err != nil {
		log.Println(err)
		return models.Event{}, errors.New("failed to update event")
	}

	foundingRange := models.FoundingRange{ID: storedEvent.FoundingRange, Low: info.FoundingRangeLow, High: info.FoundingRangeHigh}
	if _, err := svc.foundingRanges.Update(ctx, foundingRange); err != nil {
		log.Println(err)
		return models.Event{}, errors.New("failed to update founding range")
	}

	coFoundingRange := models.FoundingRange{ID: storedEvent.FoundingRange, Low: info.CoFoundingRangeLow, High: info.CoFoundingRangeHigh}
	if _, err := svc.coFoundingRanges.Update(ctx, coFoundingRange); err != nil {
		log.Println(err)
		return models.Event{}, errors.New("failed to update co-founding range")
	}

	// TODO: add here update of subjects

	return storedEvent, nil
}

func NewEventServices(storage storages.EventStorage,
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

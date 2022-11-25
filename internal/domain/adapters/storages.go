package adapters

import (
	"context"

	"github.com/google/uuid"

	"github.com/indigowar/map-of-events/internal/domain/models"
	"github.com/indigowar/map-of-events/pkg/errors"
)

// It's a set of error reasons
// This integer values are used to change logic of handling error.
/// When the repository returns an error, you can access to this reason as:
const (
	ErrReasonInternalStorageErr = iota
	ErrReasonObjectNotFoundErr
	ErrReasonObjectAlreadyExistsErr
)

// StorageWithTransaction is an ability of storage to create a transaction for other storages
type StorageWithTransaction interface {
	// BeginTransaction - creates a transaction and returns an object of transaction as interface{}
	BeginTransaction(ctx context.Context) (interface{}, error)
	// CloseTransaction is a closing method of transactions
	CloseTransaction(ctx context.Context, transaction interface{}) error
}

// CompetitorStorage - interface for storing models.Competitor
type CompetitorStorage interface {
	// AllIDs - return all ids of existing competitors
	AllIDs(ctx context.Context) ([]uuid.UUID, errors.Error)
	// Get - get a competitor with given id
	Get(ctx context.Context, id uuid.UUID) (models.Competitor, errors.Error)
	// GetAll - get all existing competitors
	GetAll(ctx context.Context) ([]models.Competitor, errors.Error)
	// Create - creates a new competitor in storage
	Create(ctx context.Context, competitor models.Competitor) errors.Error
	// Update - updates a competitor in storage
	Update(ctx context.Context, competitor models.Competitor) errors.Error
	// Delete - deletes competitor in storage with given ID
	Delete(ctx context.Context, id uuid.UUID) errors.Error

	StorageWithTransaction
}

// SubjectStorage - interface for storing models.Subject
type SubjectStorage interface {
	// GetByID - get subject with given ID
	GetByID(ctx context.Context, id uuid.UUID) (models.Subject, error)
	// GetByEvent - get subject by event id
	GetByEvent(ctx context.Context, id uuid.UUID) ([]uuid.UUID, error)
	// GetAll - get all existing subjects
	GetAll(ctx context.Context) ([]models.Subject, error)
	// Add - adds new subject to the storage
	Add(ctx context.Context, subject models.Subject) error
	// Delete - deletes subject with given id from the storage
	Delete(ctx context.Context, id uuid.UUID) error
	// Update - updates a subject in storage(will return an error, if it does not exist)
	Update(ctx context.Context, subject models.Subject) error

	StorageWithTransaction
}

// RangeStorage - a generic interface for storing models.RangeModel in the storages
type RangeStorage interface {
	// GetByID - get range by it's id
	GetByID(ctx context.Context, id uuid.UUID) (models.RangeModel, error)

	GetMaximumRange(ctx context.Context) (models.RangeModel, error)
	Create(ctx context.Context, foundingRange models.RangeModel) (models.RangeModel, error)
	Delete(ctx context.Context, id uuid.UUID) error

	StorageWithTransaction
}

// OrganizerStorage - interface for storing models.Organizer and their levels(models.OrganizerLevel)
type OrganizerStorage interface {
	// GetAllIDs - returns IDs of all organizers
	GetAllIDs(ctx context.Context) ([]uuid.UUID, error)
	// GetByID - returns organizer with given ID
	GetByID(ctx context.Context, id uuid.UUID) (models.Organizer, error)
	// GetAll - returns all organizers
	GetAll(ctx context.Context) ([]models.Organizer, error)
	// Add - adds new organizer to the storage
	Add(ctx context.Context, organizer models.Organizer) error
	// Remove - removes an organizer with given id from storage
	Remove(ctx context.Context, id uuid.UUID) error
	Update(ctx context.Context, organizer models.Organizer) error

	// GetLevelsIDs - returns IDs of all organizer levels
	GetLevelsIDs(ctx context.Context) ([]uuid.UUID, error)
	// GetLevels - returns all organizer levels from storage
	GetLevels(ctx context.Context) ([]models.OrganizerLevel, error)
	// AddLevel - adds a new organizer level to the storage
	AddLevel(ctx context.Context, level models.OrganizerLevel) error
	// RemoveLevel - removes level from the storage with given ID
	RemoveLevel(ctx context.Context, id uuid.UUID) error

	StorageWithTransaction
}

// ImageStorage - interface for storing models.Image
type ImageStorage interface {
	// GetAllLinks - returns links of all images
	GetAllLinks(ctx context.Context) ([]string, error)
	// Get - get image by it's link
	Get(ctx context.Context, link string) (models.StoredImage, error)
	// Add -  adds new image to the storage
	Add(ctx context.Context, image models.StoredImage) error
	// Remove - removes image from the storage
	Remove(ctx context.Context, link string) error
	// Update - updates image in the storage
	Update(ctx context.Context, image models.StoredImage) error
}

// EventStorage - interface for storing models.Event
type EventStorage interface {
	// GetIDList - returns ids of all events
	GetIDList(ctx context.Context) ([]uuid.UUID, error)
	// GetByID - returns models.Event with given ID
	GetByID(ctx context.Context, id uuid.UUID) (models.Event, error)
	// Add - adds a new event to the storage
	Add(ctx context.Context, event models.Event) error
	// Remove - removes event with given ID from the storage
	Remove(ctx context.Context, id uuid.UUID) error
	// Update - updates event in the storage
	Update(ctx context.Context, event models.Event) error

	// AddCompetitor - adds a competitor(competitorId) for event(id)
	AddCompetitor(ctx context.Context, id, competitorId uuid.UUID) error
	// RemoveCompetitor - removes competitor(competitorId) for event(id)
	RemoveCompetitor(ctx context.Context, id, competitorId uuid.UUID) error
	// GetCompetitors - get competitors for event with given id
	GetCompetitors(ctx context.Context, id uuid.UUID) ([]uuid.UUID, error)

	StorageWithTransaction
}

type UserStorage interface {
	GetByID(ctx context.Context, id uuid.UUID) (models.User, error)
	GetByName(ctx context.Context, name string) (models.User, error)
	Create(ctx context.Context, user models.User) error
	Delete(ctx context.Context, id uuid.UUID) error
	UpdateName(ctx context.Context, id uuid.UUID, name string) error
	UpdatePassword(ctx context.Context, id uuid.UUID, password string) error
}

type SessionStorage interface {
	GetByToken(ctx context.Context, token string) (models.TokenSession, error)
	GetByUser(ctx context.Context, id uuid.UUID) (models.TokenSession, error)

	Create(ctx context.Context, session models.TokenSession) error
	Delete(ctx context.Context, token string) error
}

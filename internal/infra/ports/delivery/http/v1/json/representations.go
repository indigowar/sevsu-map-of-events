package json

import (
	"github.com/google/uuid"

	"github.com/indigowar/map-of-events/internal/domain/models"
)

/*

This file does definition of representation structures,
which are used in the handlers.

They're specifies a json  serialisation of structure in different ways.

*/

type (
	rangeRepresentation struct {
		Low  int `json:"low"`
		High int `json:"high"`
	}

	rangeRepresentationWithID struct {
		ID uuid.UUID `json:"id"`
		rangeRepresentation
	}

	organizerLevelRepresentation struct {
		Name string `json:"name"`
		Code string `json:"code"`
	}

	organizerLevelRepresentationWithID struct {
		ID uuid.UUID `json:"id"`
		organizerLevelRepresentation
	}

	organizerRepresentation struct {
		Name  string    `json:"name"`
		Logo  string    `json:"logo"`
		Level uuid.UUID `json:"level"`
	}

	organizerRepresentationWithID struct {
		Id uuid.UUID `json:"id"`
		organizerRepresentation
	}
)

func (rr *rangeRepresentation) ToModel() models.RangeModel {
	return models.RangeModel{
		ID:   uuid.UUID{},
		Low:  rr.Low,
		High: rr.High,
	}
}

func (rr *rangeRepresentation) FromModel(m models.RangeModel) {
	rr.Low = m.Low
	rr.High = m.High
}

func (rr *rangeRepresentationWithID) ToModel() models.RangeModel {
	return models.RangeModel{
		ID:   rr.ID,
		Low:  rr.Low,
		High: rr.High,
	}
}

func (rr *rangeRepresentationWithID) FromModel(m models.RangeModel) {
	rr.Low = m.Low
	rr.High = m.High
	rr.ID = m.ID
}

func (olr *organizerLevelRepresentation) ToModel() models.OrganizerLevel {
	return models.OrganizerLevel{
		Name: olr.Name,
		Code: olr.Code,
	}
}

func (olr *organizerLevelRepresentation) FromModel(m models.OrganizerLevel) {
	olr.Name = m.Name
	olr.Code = m.Code
}

func (olr *organizerLevelRepresentationWithID) ToModel() models.OrganizerLevel {
	return models.OrganizerLevel{
		ID:   olr.ID,
		Name: olr.Name,
		Code: olr.Code,
	}
}

func (olr *organizerLevelRepresentationWithID) FromModel(m models.OrganizerLevel) {
	olr.ID = m.ID
	olr.Name = m.Name
	olr.Code = m.Code
}

func (r *organizerRepresentation) ToModel() models.Organizer {
	return models.Organizer{
		ID:    uuid.UUID{},
		Name:  r.Name,
		Logo:  r.Logo,
		Level: r.Level,
	}
}

func (r *organizerRepresentation) FromModel(m models.Organizer) {
	r.Name = m.Name
	r.Logo = m.Logo
	r.Level = m.Level
}

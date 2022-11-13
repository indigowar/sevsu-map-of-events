package models

import (
	"time"

	"github.com/google/uuid"
)

type (
	Event interface {
		ID() uuid.UUID
		Title() string
		Organizer() uuid.UUID
		FoundingType() string
		FoundingRange() uuid.UUID
		CoFoundingRange() uuid.UUID
		SubmissionDeadline() time.Time
		ConsiderationPeriod() string
		RealisationPeriod() string
		Result() string
		Site() string
		Document() string
		InternalContacts() string
		TRL() int
		Competitors() []uuid.UUID
		Subjects() []uuid.UUID
	}

	event struct {
		id                  uuid.UUID
		title               string
		organizer           uuid.UUID
		foundingType        string
		foundingRange       uuid.UUID
		coFoundingRange     uuid.UUID
		submissionDeadline  time.Time
		considerationPeriod string
		realisationPeriod   string
		result              string
		site                string
		document            string
		internalContacts    string
		trl                 int
		competitors         []uuid.UUID
		subjects            []uuid.UUID
	}
)

func (e event) ID() uuid.UUID {
	return e.id
}

func (e event) Title() string {
	return e.title
}

func (e event) Organizer() uuid.UUID {
	return e.organizer
}

func (e event) FoundingType() string {
	return e.foundingType
}

func (e event) FoundingRange() uuid.UUID {
	return e.foundingRange
}

func (e event) CoFoundingRange() uuid.UUID {
	return e.coFoundingRange
}
func (e event) SubmissionDeadline() time.Time {
	return e.submissionDeadline
}

func (e event) ConsiderationPeriod() string {
	return e.considerationPeriod
}

func (e event) RealisationPeriod() string {
	return e.realisationPeriod
}

func (e event) Result() string {
	return e.result
}

func (e event) Site() string {
	return e.site
}

func (e event) Document() string {
	return e.document
}

func (e event) InternalContacts() string {
	return e.internalContacts
}

func (e event) TRL() int {
	return e.trl
}

func (e event) Competitors() []uuid.UUID {
	return e.competitors
}

func (e event) Subjects() []uuid.UUID {
	return e.subjects
}

func NewEvent(title string,
	organizer uuid.UUID,
	foundingType string,
	foundingRange uuid.UUID,
	coFoundingRange uuid.UUID,
	submissionDeadline time.Time,
	considerationPeriod string,
	realisationPeriod string,
	result string,
	site string,
	document string,
	internalContacts string,
	tlr int,
	competitors []uuid.UUID,
	subjects []uuid.UUID,
) Event {
	return &event{
		id:                  uuid.New(),
		title:               title,
		organizer:           organizer,
		foundingType:        foundingType,
		foundingRange:       foundingRange,
		coFoundingRange:     coFoundingRange,
		submissionDeadline:  submissionDeadline,
		considerationPeriod: considerationPeriod,
		realisationPeriod:   realisationPeriod,
		result:              result,
		site:                site,
		document:            document,
		internalContacts:    internalContacts,
		trl:                 tlr,
		competitors:         competitors,
		subjects:            subjects,
	}
}

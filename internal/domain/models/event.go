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
		FoundingRange() uuid.UUID
		CoFoundingRange() uuid.UUID
		SubmissionDeadline() time.Time
		ConsiderationPeriod() time.Duration
		RealisationPeriod() time.Duration
		Result() string
		Site() string
		Document() string
		InternalContacts() string
		TLR() int
		Subjects() []uuid.UUID
		Competitors() []uuid.UUID
	}

	event struct {
		id                  uuid.UUID
		title               string
		organizer           uuid.UUID
		foundingRange       uuid.UUID
		coFoundingRange     uuid.UUID
		submissionDeadline  time.Time
		considerationPeriod time.Duration
		realisationPeriod   time.Duration
		result              string
		site                string
		document            string
		internalContacts    string
		tlr                 int
		subjects            []uuid.UUID
		competitors         []uuid.UUID
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

func (e event) FoundingRange() uuid.UUID {
	return e.foundingRange
}

func (e event) CoFoundingRange() uuid.UUID {
	return e.coFoundingRange
}
func (e event) SubmissionDeadline() time.Time {
	return e.submissionDeadline
}

func (e event) ConsiderationPeriod() time.Duration {
	return e.considerationPeriod
}

func (e event) RealisationPeriod() time.Duration {
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

func (e event) TLR() int {
	return e.tlr
}

func (e event) Subjects() []uuid.UUID {
	return e.subjects
}

func (e event) Competitors() []uuid.UUID {
	return e.competitors
}

func NewEvent(title string, organizer uuid.UUID, foundingRange uuid.UUID, coFoundingRange uuid.UUID, submissionDeadline time.Time, considerationPeriod time.Duration, realisationPeriod time.Duration, result string, site string, document string, internalContacts string, tlr int, subjects []uuid.UUID, competitors []uuid.UUID) Event {
	return &event{
		id:                  uuid.New(),
		title:               title,
		organizer:           organizer,
		foundingRange:       foundingRange,
		coFoundingRange:     coFoundingRange,
		submissionDeadline:  submissionDeadline,
		considerationPeriod: considerationPeriod,
		realisationPeriod:   realisationPeriod,
		result:              result,
		site:                site,
		document:            document,
		internalContacts:    internalContacts,
		tlr:                 tlr,
		subjects:            subjects,
		competitors:         competitors}
}

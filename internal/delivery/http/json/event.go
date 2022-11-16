package json

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/indigowar/map-of-events/internal/domain/models"
	"github.com/indigowar/map-of-events/internal/domain/services"
)

func GetAllEventHandler(svc services.EventService) gin.HandlerFunc {
	return func(c *gin.Context) {
		ids, err := svc.GetAll(c)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}
		c.JSON(http.StatusOK, ids)
	}
}

func GetByIDEventHandler(svc services.EventService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var id uuid.UUID
		{
			notParsedId := c.Param("id")
			i, err := uuid.Parse(notParsedId)
			id = i
			if err != nil {
				c.Status(http.StatusBadRequest)
				return
			}
		}

		event, err := svc.GetByID(c, id)
		if err != nil {
			c.Status(http.StatusNotFound)
			return
		}
		c.JSON(http.StatusOK, fromModel(event))
	}
}

func DeleteEventHandler(svc services.EventService) gin.HandlerFunc {
	return func(c *gin.Context) {
		notParsed := c.Param("id")
		id, err := uuid.Parse(notParsed)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}

		err = svc.Delete(c, id)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		c.Status(http.StatusAccepted)
	}
}

func GetAllAsMinimalHandler(svc services.EventService) gin.HandlerFunc {
	return func(c *gin.Context) {
		events, err := svc.GetAllAsMinimal(c)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}
		results := make([]eventMinimalBinding, len(events))
		for i, v := range events {
			results[i] = eventMinimalBinding{
				ID:                 v.ID,
				Title:              v.Title,
				Organizer:          v.Organizer,
				SubmissionDeadline: v.SubmissionDeadline,
				TRL:                v.TRL,
			}
		}
	}
}

func GetByIDAsMinimalHandler(svc services.EventService) gin.HandlerFunc {
	return func(c *gin.Context) {
		notParsed := c.Param("id")
		id, err := uuid.Parse(notParsed)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		event, err := svc.GetByIDAsMinimal(c, id)
		if err != nil {
			c.Status(http.StatusNotFound)
			return
		}

		c.JSON(http.StatusOK, eventMinimalBinding{
			ID:                 event.ID,
			Title:              event.Title,
			Organizer:          event.Organizer,
			SubmissionDeadline: event.SubmissionDeadline,
			TRL:                event.TRL,
		})
	}
}

type eventMinimalBinding struct {
	ID                 uuid.UUID `json:"ID"`
	Title              string    `json:"title"`
	Organizer          uuid.UUID `json:"organizer"`
	SubmissionDeadline time.Time `json:"submissionDeadline"`
	TRL                int       `json:"trl"`
}

type eventBinding struct {
	ID                  uuid.UUID   `json:"ID"`
	Title               string      `json:"title"`
	Organizer           uuid.UUID   `json:"organizer"`
	FoundingType        string      `json:"foundingType"`
	FoundingRange       uuid.UUID   `json:"foundingRange"`
	CoFoundingRange     uuid.UUID   `json:"coFoundingRange"`
	SubmissionDeadline  time.Time   `json:"submissionDeadline"`
	ConsiderationPeriod string      `json:"considerationPeriod"`
	RealisationPeriod   string      `json:"realisationPeriod"`
	Result              string      `json:"result"`
	Site                string      `json:"site"`
	Document            string      `json:"document"`
	InternalContacts    string      `json:"internalContacts"`
	TRL                 int         `json:"TRL"`
	Competitors         []uuid.UUID `json:"competitors"`
	Subjects            []uuid.UUID `json:"subjects"`
}

func fromModel(e models.Event) eventBinding {
	return eventBinding{
		ID:                  e.ID,
		Title:               e.Title,
		Organizer:           e.Organizer,
		FoundingType:        e.FoundingType,
		FoundingRange:       e.FoundingRange,
		CoFoundingRange:     e.CoFoundingRange,
		SubmissionDeadline:  e.SubmissionDeadline,
		ConsiderationPeriod: e.ConsiderationPeriod,
		RealisationPeriod:   e.RealisationPeriod,
		Result:              e.Result,
		Site:                e.Site,
		Document:            e.Document,
		InternalContacts:    e.InternalContacts,
		TRL:                 e.TRL,
		Competitors:         e.Competitors,
		Subjects:            e.Subjects,
	}
}

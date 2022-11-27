package json

import (
	"log"
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

func GetByIDEventHandler(svc services.EventService, subjects services.SubjectService) gin.HandlerFunc {
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
		subjects, err := subjects.GetAllForEvent(c, id)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusOK, fromModel(event))
			return
		}
		subjectsPresentation := make([]string, len(subjects))
		for i, v := range subjects {
			subjectsPresentation[i] = v.Name
		}

		result := fromModel(event)
		result.Subjects = subjectsPresentation

		c.JSON(http.StatusOK, result)
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

func CreateEventHandler(svc services.EventService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var info createUpdateInfoBinding
		if err := c.ShouldBindJSON(&info); err != nil {
			log.Println(err)
			c.Status(http.StatusBadRequest)
			return
		}

		result, err := svc.Create(c, toUpdateCreateInfo(info))
		if err != nil {
			log.Println(err)
			c.Status(http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusCreated, fromModel(result))
	}
}

func UpdateEventHandler(svc services.EventService) gin.HandlerFunc {
	return func(c *gin.Context) {
		paramId := c.Param("id")
		_, err := uuid.Parse(paramId)
		if err != nil {
			log.Println(err)
			c.Status(http.StatusBadRequest)
			return
		}

		var info createUpdateInfoBinding
		if err := c.ShouldBindJSON(&info); err != nil {
			log.Println(err)
			c.Status(http.StatusBadRequest)
			return
		}

		panic("not implemented")
	}
}

type eventMinimalBinding struct {
	ID                 uuid.UUID `json:"id"`
	Title              string    `json:"title"`
	Organizer          uuid.UUID `json:"organizer"`
	SubmissionDeadline time.Time `json:"submissionDeadline"`
	TRL                int       `json:"trl"`
}

type eventBinding struct {
	ID                  uuid.UUID   `json:"id"`
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
	TRL                 int         `json:"trl"`
	Competitors         []uuid.UUID `json:"competitors"`
	Subjects            []string    `json:"subjects"`
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
	}
}

type createUpdateInfoBinding struct {
	Title               string      `json:"title"`
	Organizer           uuid.UUID   `json:"organizer"`
	FoundingType        string      `json:"foundingType"`
	FoundingRangeLow    int         `json:"foundingRangeLow"`
	FoundingRangeHigh   int         `json:"foundingRangeHigh"`
	CoFoundingRangeLow  int         `json:"coFoundingRangeLow"`
	CoFoundingRangeHigh int         `json:"coFoundingRangeHigh"`
	SubmissionDeadline  time.Time   `json:"submissionDeadline"`
	ConsiderationPeriod string      `json:"considerationPeriod"`
	RealisationPeriod   string      `json:"realisationPeriod"`
	Result              string      `json:"result"`
	Site                string      `json:"site"`
	Document            string      `json:"document"`
	InternalContacts    string      `json:"internalContacts"`
	TRL                 int         `json:"trl"`
	Competitors         []uuid.UUID `json:"competitors"`
	Subjects            []string    `json:"subjects"`
}

func toUpdateCreateInfo(i createUpdateInfoBinding) services.EventCreateInfo {
	return services.EventCreateInfo{
		Title:               i.Title,
		Organizer:           i.Organizer,
		FoundingType:        i.FoundingType,
		FoundingRangeLow:    i.FoundingRangeLow,
		FoundingRangeHigh:   i.FoundingRangeHigh,
		CoFoundingRangeHigh: i.CoFoundingRangeHigh,
		CoFoundingRangeLow:  i.CoFoundingRangeLow,
		SubmissionDeadline:  i.SubmissionDeadline,
		ConsiderationPeriod: i.ConsiderationPeriod,
		RealisationPeriod:   i.RealisationPeriod,
		Result:              i.Result,
		Site:                i.Site,
		Document:            i.Document,
		InternalContacts:    i.InternalContacts,
		TRL:                 i.TRL,
		Competitors:         i.Competitors,
		Subjects:            i.Subjects,
	}
}

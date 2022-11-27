package json

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/indigowar/map-of-events/internal/domain/models"
	"github.com/indigowar/map-of-events/internal/domain/services"
)

type EventHandler struct {
	svc services.Services
}

func NewEventHandler(s services.Services) EventHandler {
	return EventHandler{svc: s}
}

func (h *EventHandler) GetAllEvents(c *gin.Context) {
	ids, err := h.svc.Event.GetAll(c)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, ids)
}

func (h *EventHandler) GetEventByID(c *gin.Context) {
	id, err := h.parseIDFromParam(c)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	result, status := h.getAndSerialize(c, id, h.buildView)
	if status != 0 {
		c.Status(status)
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *EventHandler) DeleteEvent(c *gin.Context) {
	id, err := h.parseIDFromParam(c)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	err = h.svc.Event.Delete(c, id)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	c.Status(http.StatusAccepted)
}

func (h *EventHandler) GetAllAsMinimal(c *gin.Context) {
	ids, err := h.svc.Event.AllIDs(c)
	if err != nil {
		log.Println(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	result := make([]interface{}, len(ids))

	for i, id := range ids {
		result[i], _ = h.getAndSerialize(c, id, h.buildMinimalView)
	}

	c.JSON(http.StatusOK, result)
}

func (h *EventHandler) GetByIDMinimal(c *gin.Context) {
	id, err := h.parseIDFromParam(c)
	if err != nil {
		log.Println(err)
		c.Status(http.StatusBadRequest)
		return
	}

	result, status := h.getAndSerialize(c, id, h.buildMinimalView)
	if status != 0 {
		c.Status(status)
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *EventHandler) Create(c *gin.Context) {
	var info createInfoView
	if err := c.ShouldBindJSON(&info); err != nil {
		log.Println(err)
		c.Status(http.StatusBadRequest)
		return
	}

	event, err := h.svc.Event.Create(c, h.createInfoFromView(info))
	if err != nil {
		log.Println(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	result, status := h.getAndSerialize(c, event.ID, h.buildView)
	if status != 0 {
		c.Status(status)
		return
	}

	c.JSON(http.StatusCreated, result)
}

func (h *EventHandler) Update(c *gin.Context) {
	id, err := h.parseIDFromParam(c)
	if err != nil {
		log.Println(err)
		c.Status(http.StatusBadRequest)
		return
	}

	var info createInfoView
	if err := c.ShouldBindJSON(&info); err != nil {
		log.Println(err)
		c.Status(http.StatusBadRequest)
		return
	}

	_, err = h.svc.Event.Update(c, id, h.createInfoFromView(info))
	if err != nil {
		log.Println(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	result, status := h.getAndSerialize(c, id, h.buildView)
	if status != 0 {
		c.Status(status)
		return
	}

	c.JSON(http.StatusAccepted, result)
}

type createInfoView struct {
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

func (h *EventHandler) parseIDFromParam(c *gin.Context) (uuid.UUID, error) {
	var id uuid.UUID
	notParsedId := c.Param("id")
	i, err := uuid.Parse(notParsedId)
	id = i
	if err != nil {
		return uuid.UUID{}, err
	}
	return id, nil
}

func (h *EventHandler) buildView(e models.Event, f, cf models.RangeModel, subs []models.Subject) interface{} {
	s := make([]string, len(subs))
	for i, v := range subs {
		s[i] = v.Name
	}

	return eventJSONView{
		ID:           e.ID,
		Title:        e.Title,
		Organizer:    e.Organizer,
		FoundingType: e.FoundingType,
		FoundingRange: rangeJSONView{
			Low:  f.Low,
			High: f.High,
		},
		CoFoundingRange: rangeJSONView{
			Low:  cf.Low,
			High: cf.High,
		},
		SubmissionDeadline:  e.SubmissionDeadline,
		ConsiderationPeriod: e.ConsiderationPeriod,
		RealisationPeriod:   e.RealisationPeriod,
		Result:              e.Result,
		Site:                e.Site,
		Document:            e.Document,
		InternalContacts:    e.InternalContacts,
		TRL:                 e.TRL,
		Competitors:         e.Competitors,
		Subjects:            s,
	}
}

func (h *EventHandler) getAndSerialize(ctx context.Context, id uuid.UUID, serializer serializerFunc) (interface{}, int) {
	// Get the event itself
	event, err := h.svc.Event.GetByID(ctx, id)
	if err != nil {
		log.Println(err)
		return nil, http.StatusNotFound
	}

	// load info about it's founding range
	foundingRange, err := h.svc.FoundingRange.GetByID(ctx, event.FoundingRange)
	if err != nil {
		log.Println(err)
		return nil, http.StatusInternalServerError
	}

	// load info about it's co-founding range
	coFoundingRange, err := h.svc.CoFoundingRange.GetByID(ctx, event.CoFoundingRange)
	if err != nil {
		log.Println(err)
		return nil, http.StatusInternalServerError
	}

	// load info about it's subjects
	subjects, err := h.svc.Subject.GetAllForEvent(ctx, event.ID)
	if err != nil {
		log.Println(err)
		return nil, http.StatusInternalServerError
	}

	return serializer(event, foundingRange, coFoundingRange, subjects), 0
}

func (h *EventHandler) buildMinimalView(e models.Event, f, cf models.RangeModel, subs []models.Subject) interface{} {
	s := make([]string, len(subs))
	for i, v := range subs {
		s[i] = v.Name
	}
	return eventMinimalJSONView{
		ID:           e.ID,
		Title:        e.Title,
		FoundingType: e.FoundingType,
		FoundingRange: rangeJSONView{
			Low:  f.Low,
			High: f.High,
		},
		CoFoundingRange: rangeJSONView{
			Low:  cf.Low,
			High: cf.High,
		},
		SubmissionDeadline: e.SubmissionDeadline,
		TRL:                e.TRL,
		Subjects:           s,
	}
}

func (h *EventHandler) createInfoFromView(i createInfoView) services.EventCreateInfo {
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

type serializerFunc func(event models.Event, f, cf models.RangeModel, subs []models.Subject) interface{}

type rangeJSONView struct {
	Low  int `json:"low"`
	High int `json:"high"`
}

type eventJSONView struct {
	ID                  uuid.UUID     `json:"id"`
	Title               string        `json:"title"`
	Organizer           uuid.UUID     `json:"organizer"`
	FoundingType        string        `json:"foundingType"`
	FoundingRange       rangeJSONView `json:"foundingRange"`
	CoFoundingRange     rangeJSONView `json:"coFoundingRange"`
	SubmissionDeadline  time.Time     `json:"submissionDeadline"`
	ConsiderationPeriod string        `json:"considerationPeriod"`
	RealisationPeriod   string        `json:"realisationPeriod"`
	Result              string        `json:"result"`
	Site                string        `json:"site"`
	Document            string        `json:"document"`
	InternalContacts    string        `json:"internalContacts"`
	TRL                 int           `json:"trl"`
	Competitors         []uuid.UUID   `json:"competitors"`
	Subjects            []string      `json:"subjects"`
}

type eventMinimalJSONView struct {
	ID                 uuid.UUID     `json:"id"`
	Title              string        `json:"title"`
	Organizer          uuid.UUID     `json:"organizer"`
	FoundingType       string        `json:"foundingType"`
	FoundingRange      rangeJSONView `json:"foundingRange"`
	CoFoundingRange    rangeJSONView `json:"coFoundingRange"`
	SubmissionDeadline time.Time     `json:"submissionDeadline"`
	TRL                int           `json:"trl"`
	Subjects           []string      `json:"subjects"`
}

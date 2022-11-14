package json

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/indigowar/map-of-events/internal/domain/services"
)

type organizerLevel struct {
	Id       uuid.UUID `json:"id"`
	Title    string    `json:"name"`
	CodeWord string    `json:"code"`
}

func GetAllOrganizerLevelsHandler(svc services.OrganizerService) func(c *gin.Context) {
	return func(c *gin.Context) {
		levels, err := svc.GetAllLevels(c)
		if err != nil {
			log.Println(err)

			c.Status(http.StatusInternalServerError)
			return
		}

		results := make([]organizerLevel, len(levels))
		for i, v := range levels {
			results[i] = organizerLevel{v.ID(), v.Name(), v.Code()}
		}

		c.JSON(http.StatusOK, results)
	}
}

type createOrganizerLevelRequest struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

func CreateOrganizerLevelHandler(svc services.OrganizerService) func(c *gin.Context) {
	return func(c *gin.Context) {
		var level createOrganizerLevelRequest
		if err := c.ShouldBindJSON(&level); err != nil {
			log.Println(err)
			c.Status(http.StatusBadRequest)
			return
		}

		o, err := svc.CreateLevel(c, level.Name, level.Code)
		if err != nil {
			log.Println(err)
			c.Status(http.StatusInternalServerError)
			return
		}
		c.JSON(http.StatusCreated, organizerLevel{o.ID(), o.Name(), o.Code()})
	}
}

func (l *organizerLevel) ID() uuid.UUID { return l.Id }
func (l *organizerLevel) Name() string  { return l.Title }
func (l *organizerLevel) Code() string  { return l.CodeWord }

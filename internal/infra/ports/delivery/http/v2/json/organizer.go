package json

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/indigowar/map-of-events/internal/domain/services"
	"github.com/indigowar/map-of-events/pkg/random"
)

type organizerBinding struct {
	Id    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Logo  string    `json:"logo"`
	Level uuid.UUID `json:"level"`
}

type createOrganizerInfo struct {
	Name  string    `json:"name"`
	Logo  string    `json:"logo"`
	Level uuid.UUID `json:"level"`
}

func GetAllOrganizersHandler(organizerSvc services.OrganizerService, imageSvc services.ImageService) gin.HandlerFunc {
	return func(c *gin.Context) {
		organizers, err := organizerSvc.GetAll(c)
		if err != nil {
			log.Println(err)
			c.Status(http.StatusBadRequest)
			return
		}
		result := make([]organizerBinding, len(organizers))
		for i, v := range organizers {
			result[i] = organizerBinding{Id: v.ID, Name: v.Name, Level: v.Level}
			image, err := imageSvc.Get(c, v.Logo)
			if err != nil {
				log.Println(err)
				continue
			}
			result[i].Logo = string(image.Value[:])
		}

		c.JSON(http.StatusOK, result)
	}
}

func GetOrganizerByID(organizerSvc services.OrganizerService, imageSvc services.ImageService) gin.HandlerFunc {
	return func(c *gin.Context) {
		parsedId := c.Param("id")
		id, err := uuid.Parse(parsedId)
		if err != nil {
			log.Println(err)
			c.Status(http.StatusBadRequest)
			return
		}
		organizer, err := organizerSvc.GetByID(c, id)
		if err != nil {
			log.Println(err)
			c.Status(http.StatusNotFound)
			return
		}
		result := organizerBinding{
			Id:    organizer.ID,
			Name:  organizer.Name,
			Level: organizer.Level,
		}

		image, err := imageSvc.Get(c, organizer.Logo)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusOK, result)
			return
		}
		result.Logo = string(image.Value[:])

		c.JSON(http.StatusOK, result)
	}
}

func CreateOrganizerHandler(orgSvc services.OrganizerService, imgSvc services.ImageService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var info createOrganizerInfo
		if err := c.ShouldBindJSON(&info); err != nil {
			log.Println(err)
			c.Status(http.StatusBadRequest)
			return
		}

		logo := random.RandStringRunes(10)

		image, err := imgSvc.Create(c, logo, []byte(info.Logo))
		if err != nil {
			log.Println(err)
			c.Status(http.StatusInternalServerError)
			return
		}

		organizer, err := orgSvc.Create(c, info.Name, image.Link, info.Level)
		if err != nil {
			log.Println(err)
			c.Status(http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusCreated, organizerBinding{
			Id:    organizer.ID,
			Name:  organizer.Name,
			Logo:  string(image.Value[:]),
			Level: organizer.Level,
		})
	}
}

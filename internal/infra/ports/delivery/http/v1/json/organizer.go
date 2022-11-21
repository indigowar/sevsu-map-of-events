package json

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/indigowar/map-of-events/internal/domain/services"
)

type orgLevelBinding struct {
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

		results := make([]orgLevelBinding, len(levels))
		for i, v := range levels {
			results[i] = orgLevelBinding{v.ID, v.Name, v.Code}
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
		c.JSON(http.StatusCreated, orgLevelBinding{o.ID, o.Name, o.Code})
	}
}

type organizerBinding struct {
	Id    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Logo  string    `json:"logo"`
	Level uuid.UUID `json:"level"`
}

func GetAllOrganizersHandler(svc services.OrganizerService) func(c *gin.Context) {
	return func(c *gin.Context) {
		objects, err := svc.GetAll(c)
		if err != nil {
			log.Println(err)
			c.Status(http.StatusInternalServerError)
			return
		}
		result := make([]organizerBinding, len(objects))
		for i, v := range objects {
			result[i] = organizerBinding{v.ID, v.Name, v.Logo, v.Level}
		}

		c.JSON(http.StatusOK, result)
	}
}

func GetByIDOrganizerHandler(svc services.OrganizerService) func(c *gin.Context) {
	return func(c *gin.Context) {
		stringId := c.Param("id")

		id, err := uuid.Parse(stringId)
		if err != nil {
			log.Println(err)
			c.Status(http.StatusBadRequest)
			return
		}

		organizer, err := svc.GetByID(c, id)
		if err != nil {
			log.Println(err)
			c.Status(http.StatusInternalServerError)
			return
		}
		c.JSON(http.StatusOK, organizerBinding{
			organizer.ID,
			organizer.Name,
			organizer.Logo,
			organizer.Level,
		})
	}
}

type createOrganizerRequest struct {
	Name  string    `json:"name"`
	Logo  string    `json:"logo"`
	Level uuid.UUID `json:"level"`
}

func CreateOrganizerHandler(svc services.OrganizerService) func(c *gin.Context) {
	return func(c *gin.Context) {
		var organizer createOrganizerRequest
		if err := c.ShouldBindJSON(&organizer); err != nil {
			log.Println(err)
			c.Status(http.StatusBadRequest)
			return
		}
		created, err := svc.Create(c, organizer.Name, organizer.Logo, organizer.Level)
		if err != nil {
			log.Println(err)
			c.Status(http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusCreated, organizerBinding{
			created.ID,
			created.Name,
			created.Logo,
			created.Level,
		})
	}
}

func UpdateOrganizerHandler(_ services.OrganizerService) func(c *gin.Context) {
	return func(c *gin.Context) {
		stringId := c.Param("id")

		_, err := uuid.Parse(stringId)
		if err != nil {
			log.Println(err)
			c.Status(http.StatusBadRequest)
			return
		}

		var input createOrganizerRequest
		if err := c.ShouldBindJSON(&input); err != nil {
			log.Println(err)
			c.Status(http.StatusBadRequest)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"msg": "unimplemented",
		})
	}
}

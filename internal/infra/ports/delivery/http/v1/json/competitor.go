package json

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/indigowar/map-of-events/internal/domain/services"
)

type Competitor struct {
	Id             uuid.UUID `json:"id"`
	CompetitorName string    `json:"name"`
}

func GetAllCompetitorsHandler(svc services.CompetitorService) func(c *gin.Context) {
	return func(c *gin.Context) {
		comps, err := svc.GetAll(c)

		if err != nil {
			switch err.Reason() {
			case services.ErrReasonInternalError:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": err.ShortErr(),
				})
				break
			}
			return
		}

		result := make([]Competitor, len(comps))
		for i, v := range comps {
			result[i] = Competitor{v.ID, v.Name}
		}
		c.JSON(http.StatusOK, result)
	}
}

func CreateCompetitorHandler(srv services.CompetitorService) func(c *gin.Context) {
	return func(c *gin.Context) {
		var name string
		if err := c.ShouldBindJSON(&name); err != nil {
			log.Println(err)
			c.Status(http.StatusBadRequest)
			return
		}

		obj, err := srv.Create(c, name)

		if err != nil {
			switch err.Reason() {
			case services.ErrReasonInternalError:
				log.Println(err.LongErr())
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": err.ShortErr(),
				})
				break
			case services.ErrReasonAlreadyExist:
				log.Println(err.LongErr())
				c.JSON(http.StatusConflict, gin.H{
					"msg": err.ShortErr(),
				})
				break
			}
			return
		}
		c.JSON(http.StatusAccepted, Competitor{obj.ID, obj.Name})
	}
}

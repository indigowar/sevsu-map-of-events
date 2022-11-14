package json

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/indigowar/map-of-events/internal/domain/services"
)

type rangeType struct {
	Id   uuid.UUID `json:"id"`
	Low  int       `json:"low"`
	High int       `json:"high"`
}

func GetByIDRangeHandler(srv services.RangeService) func(c *gin.Context) {
	return func(c *gin.Context) {
		stringId := c.Param("id")

		id, err := uuid.Parse(stringId)

		if err != nil {
			log.Println(err)
			c.Status(http.StatusBadRequest)
			return
		}

		r, err := srv.GetByID(c, id)
		if err != nil {
			log.Println(err)
			c.Status(http.StatusNotFound)
		}

		c.JSON(http.StatusOK, rangeType{r.ID(), r.Low(), r.High()})
	}
}

func GetMaximumRangeHandler(srv services.RangeService) func(c *gin.Context) {
	return func(c *gin.Context) {
		result, err := srv.GetMaximumRange(c)
		if err != nil {
			log.Println(err)
			c.Status(http.StatusInternalServerError)
			return
		}
		c.JSON(http.StatusOK, rangeType{result.ID(), result.Low(), result.High()})
	}
}

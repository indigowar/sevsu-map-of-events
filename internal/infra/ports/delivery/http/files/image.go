package files

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/indigowar/map-of-events/internal/domain/services"
)

func UploadHandler(svc services.ImageService) gin.HandlerFunc {
	return func(c *gin.Context) {
		image, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			log.Println(err)
			c.Status(http.StatusBadRequest)
			return
		}

		result, err := svc.Create(c, uuid.New().String(), image)

		c.JSON(http.StatusCreated, gin.H{
			"link": result.Link,
		})
	}
}

func RetrievingHandler(svc services.ImageService) gin.HandlerFunc {
	return func(c *gin.Context) {
		link := c.Param("link")

		m, err := svc.Get(c, link)
		if err != nil {
			log.Println(err)
			c.Status(http.StatusNotFound)
			return
		}

		c.JSON(http.StatusOK, m.Value)
	}
}

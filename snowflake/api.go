package snowflake

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type GeneratorAPI struct {
	Logger *zap.Logger
}

func NewGeneratorAPI(logger *zap.Logger) (app *GeneratorAPI) {
	app = &GeneratorAPI{Logger: logger}
	return app
}

func (app *GeneratorAPI) InitRoute(engine *gin.Engine, groupPath string) {
	group := engine.Group(groupPath)
	group.POST("/:key", app.idgen)
}

func (app *GeneratorAPI) idgen(c *gin.Context) {
	rawKey := c.Param("key")
	workerId, err := strconv.ParseInt(rawKey, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			gin.H{
				"error": "Key must be a string of integer",
			})
		return
	}
	err, idWorker := NewIdWorker(workerId)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			gin.H{
				"error": fmt.Sprintf("Key error: %v", err),
			})
		return
	}
	err, nextId := idWorker.NextId()
	if err != nil {
		c.JSON(http.StatusBadRequest,
			gin.H{
				"error": fmt.Sprintf("Service error: %v", err),
			})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"id": nextId,
	})
}

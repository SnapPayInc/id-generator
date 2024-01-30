package snowflake

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var SIdWorkers map[int]*IdWorker

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
	workerId, err := strconv.Atoi(rawKey)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			gin.H{
				"error": "Key must be a string of integer",
			})
		return
	}
	idWorker, ok := SIdWorkers[workerId]
	if ok {
		err, nextId := idWorker.NextId()
		if err != nil {
			c.JSON(http.StatusBadRequest,
				gin.H{
					"error": fmt.Sprintf("Service error: %v", err),
				})
			return
		}
		//fmt.Println(nextId)
		c.JSON(http.StatusOK, gin.H{
			"id": nextId,
		})
	} else {
		c.JSON(http.StatusBadRequest,
			gin.H{
				"error": fmt.Sprintf("Invalid key: %d", workerId),
			})
		return
	}
}

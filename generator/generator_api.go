package generator

import (
	"crypto/sha1"
	"fmt"
	"strconv"
	"id-generator/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
)

type GeneratorAPI struct {
	generatorStore *GeneratorDB
	Logger         *zap.Logger
}

func NewGeneratorAPI(generatorStore *GeneratorDB, logger *zap.Logger) (app *GeneratorAPI) {
	app = &GeneratorAPI{generatorStore: generatorStore, Logger: logger}
	return app
}

func (app *GeneratorAPI) InitRoute(engine *gin.Engine, groupPath string) {
	group := engine.Group(groupPath)
	group.PUT("/:key/tap", app.tap)
	group.POST("/:key/set", app.set)
}
func tableNameHash(name string) string {
	h := sha1.New()
	h.Write([]byte(name))
	bs := h.Sum(nil)
	return fmt.Sprintf("h%x", bs)
}

func (app *GeneratorAPI) tap(c *gin.Context) {
	key := tableNameHash(c.Param("key"))
	var count int
	count,err := strconv.Atoi(c.DefaultQuery("count", "1"))
	if err != nil {
		count = 1
	}
	utils.LogInfo("create table")
	err = app.generatorStore.CreateTableIfNotExist(key)
	if err != nil {
		app.Logger.Debug("Cannot create table")
		c.JSON(http.StatusBadRequest,
			gin.H{
				"error": "Cannot create table",
			})
		return
	}
	utils.LogInfo("create table done")

	utils.LogInfo("insert")
	seed, err := app.generatorStore.Insert(key, count)
	if err != nil {
		app.Logger.Debug("Insert error")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Insert error",
		})
		return
	}
	utils.LogInfo("insert done")
	if count > 1 {
		seeds := make([]int64, count)
		for i := range seeds {
			seeds[i] = seed - int64(i)
		}
		c.JSON(http.StatusOK, gin.H{
			"last_insert_ids": seeds,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"last_insert_id": seed,
		})
	}
	
}

func (app *GeneratorAPI) set(c *gin.Context) {
	key := tableNameHash(c.Param("key"))
	err := app.generatorStore.CreateTableIfNotExist(key)
	if err != nil {
		fmt.Println(1, err)
		app.Logger.Debug("Cannot create table")
		c.JSON(http.StatusBadRequest,
			gin.H{
				"error": "Cannot create table",
			})
		return
	}

	updateMap := make(map[string]interface{})
	err = c.ShouldBind(&updateMap)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Parse map error",
		})
		return
	}
	
	if val, ok := updateMap["value"]; ok {
		value := int64(val.(float64))
		seed, err := app.generatorStore.Set(key, value)
		if err != nil {
			app.Logger.Debug("Insert error")
			existed, _ := app.generatorStore.QueryLast(key)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("Insert error: the value must be larger then %v", existed),
			})
			return
		}

		if seed != value {
			fmt.Println(seed, value)
		}

		c.JSON(http.StatusOK, gin.H{
			"last_insert_id": seed,
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "No value in post body",
		})
		return
	}
}

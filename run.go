package main

import (
  "os"
  "fmt"
  "flag"
	"net/http"
  "gitlab.com/paiduay/queq-hospital-api/utils"
  "gitlab.com/paiduay/queq-hospital-api/config"
  "gitlab.com/paiduay/queq-hospital-api/models"
  "gitlab.com/paiduay/queq-hospital-api/middlewares"
  "gitlab.com/paiduay/queq-hospital-api/controllers"
	"github.com/gin-gonic/gin"
  _ "github.com/jinzhu/gorm/dialects/mssql"
)

var port = flag.Int("port", 7000, "default port is 7000")
var runEnv = flag.String("env", "dev", "default env is dev")

func init () {
  flag.Parse()
  config.RunPort = *port
  if *runEnv == "dev" {
    config.RunMode = "dev"
  } else if *runEnv == "prod" {
    config.RunMode = "prod"
  } else {
    config.RunMode = "dev"
  }

  fmt.Println("Running in " + *runEnv + " mode.")

  utils.ConfigLoader(*runEnv)
  if err := models.InitDatabase(); err != nil {
    fmt.Println("Cant connect to database due to: ")
    fmt.Println(err)
    fmt.Println("Program terminated.")
    os.Exit(3)
  } else {
    fmt.Println("Database Connected.")
  }
}

func main () {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
  router.Use(middlewares.CORSMiddleware())
	router.Use(middlewares.TokenCheck())

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, utils.ErrorMessage("endpoint not found", 404))
	})
	router.NoMethod(func(c *gin.Context) {
		c.JSON(http.StatusMethodNotAllowed, utils.ErrorMessage("method not allow", 405))
	})

  controllers.RegisterStaffEndpoints(router.Group("/staff"))
  controllers.RegisterStationEndpoints(router.Group("/station"))

	router.Run(fmt.Sprintf(":%v", *port))
}

package controllers

import "C"
import (
	// "fmt"
	"strconv"
	"net/http"
	// "mime/multipart"

	"github.com/gin-gonic/gin"

	"gitlab.com/paiduay/queq-hospital-api/middlewares"
	"gitlab.com/paiduay/queq-hospital-api/models"
	"gitlab.com/paiduay/queq-hospital-api/utils"
)

// RegisterStationEndpoints - to let main register these endpoints
func RegisterStationEndpoints(router *gin.RouterGroup) {
	stationRouter := router.Group("")
	// Annonymous Routes

	// Admin Routes
	stationRouter.Use(middlewares.LoginRequire())
	{
    stationRouter.GET("/highestWaitingTimeQueue/:sID", getHighestWaitingTimeQueue)
    stationRouter.GET("/check/:sCode", checkStationCode)
    stationRouter.GET("/checkStation", checkStationName)
    stationRouter.GET("/rooms/:sID", getStationInfo)
    stationRouter.GET("/rooms/:sID/withQueue", getStationInfo)
	}

	stationRouter.Use(middlewares.SuperAdminRequired())
	{
    stationRouter.POST("/add", addStation)
    stationRouter.POST("/edit/:sID", editStation)
    stationRouter.GET("/delete/:sID", deleteStation)
    stationRouter.GET("/info/:sID", getStationInfo)
    stationRouter.POST("/reorder", reorderStation)
    stationRouter.GET("/multistation/:hID", getStationByHospital)
	}
}

func getHighestWaitingTimeQueue (c *gin.Context) {
	var highestQ models.HighWaittime
	var date string
	var strSID string
	var sID uint64
	var err error

	strSID = c.Param("sID")
	date = c.Query("date")

	if sID, err = strconv.ParseUint(strSID, 10, 64); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, utils.ErrorMessage("bad request", http.StatusBadRequest))
		return
	}

	models.GetHighestWaitingTime(date, sID, &highestQ)

  c.AbortWithStatusJSON(http.StatusOK, utils.SuccessMessage(gin.H{
		"highestWaitingTimeQueue": highestQ,
	}))
}

func checkStationCode (c *gin.Context) {
	sCode := c.Param("sCode")

	available := models.CheckStationCodeAvailable(sCode)

  c.AbortWithStatusJSON(http.StatusOK, utils.SuccessMessage(gin.H{
		"available": available,
		"stationCode": sCode,
	}))
}

func checkStationName (c *gin.Context) {
	var sName string
	var strHID string
	var hID uint64
	var err error
	if sName = c.Query("name"); sName == "" {
		c.AbortWithStatusJSON(http.StatusOK, utils.ErrorMessage("bad request", http.StatusBadRequest))
		return
	}
	if strHID = c.Query("hospitalId"); strHID == "" {
		c.AbortWithStatusJSON(http.StatusOK, utils.ErrorMessage("bad request", http.StatusBadRequest))
		return
	}
	if hID, err = strconv.ParseUint(strHID, 10, 64); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, utils.ErrorMessage("bad request", http.StatusBadRequest))
		return
	}

	available := models.CheckStationNameAvailable(sName, hID)

  c.AbortWithStatusJSON(http.StatusOK, utils.SuccessMessage(gin.H{
		"available": available,
		"stationName": sName,
	}))
}

func addStation (c *gin.Context) {
	var station models.Station
	if err := c.BindJSON(&station); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, utils.ErrorMessage("bad request", http.StatusBadRequest))
		return
	}

	if err := models.CreateStation(&station); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, utils.ErrorMessage("internal error", http.StatusInternalServerError))
		return
	}

  c.AbortWithStatusJSON(http.StatusOK, utils.SuccessMessage(gin.H{
		"status": true,
	}))
}

func editStation (c *gin.Context) {
	var station models.Station
	if err := c.BindJSON(&station); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, utils.ErrorMessage("bad request", http.StatusBadRequest))
		return
	}

  c.AbortWithStatusJSON(http.StatusOK, utils.SuccessMessage(gin.H{
		"status": true,
	}))
}

func deleteStation (c *gin.Context) {
  c.AbortWithStatusJSON(http.StatusOK, utils.SuccessMessage(gin.H{
		"status": true,
	}))
}

func getStationInfo (c *gin.Context) {
  c.AbortWithStatusJSON(http.StatusOK, utils.SuccessMessage(gin.H{
		"status": true,
	}))
}

func reorderStation (c *gin.Context) {
  c.AbortWithStatusJSON(http.StatusOK, utils.SuccessMessage(gin.H{
		"status": true,
	}))
}

func getStationByHospital (c *gin.Context) {
  c.AbortWithStatusJSON(http.StatusOK, utils.SuccessMessage(gin.H{
		"status": true,
	}))
}

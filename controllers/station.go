package controllers

import "C"
import (
	// "fmt"
	"time"
	"net/http"
	"strconv"

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
		stationRouter.GET("/rooms/:sID", getRoomInStation)
		stationRouter.GET("/rooms/:sID/withQueue", getRoomInStationWithBrief)
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

func getHighestWaitingTimeQueue(c *gin.Context) {
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

func checkStationCode(c *gin.Context) {
	sCode := c.Param("sCode")

	available := models.CheckStationCodeAvailable(sCode)

	c.AbortWithStatusJSON(http.StatusOK, utils.SuccessMessage(gin.H{
		"available":   available,
		"stationCode": sCode,
	}))
}

func checkStationName(c *gin.Context) {
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
		"available":   available,
		"stationName": sName,
	}))
}

func getRoomInStation (c *gin.Context) {
	var rooms []models.Room
	sID := c.Param("sID")

	if exist := models.CheckStationExist(sID); !exist {
		c.AbortWithStatusJSON(http.StatusOK, utils.ErrorMessage("station not found", http.StatusNotFound))
		return
	}

	models.GetRoomsInStation(sID, &rooms)

	c.AbortWithStatusJSON(http.StatusOK, utils.SuccessMessage(gin.H{
		"status": true,
		"rooms": rooms,
	}))
}

func getRoomInStationWithBrief (c *gin.Context) {
	var rooms []models.Room
	sID := c.Param("sID")
	fromDate := c.Query("fromDate")
	toDate := c.Query("toDate")
	d := time.Now()
	if fromDate == "" || toDate == "" {
		fromDate = d.Format("2006-01-02")
		toDate = d.Format("2006-01-02")
	}

	models.GetRoomsInStationWithBrief(sID, fromDate, toDate, &rooms)

	c.AbortWithStatusJSON(http.StatusOK, utils.SuccessMessage(gin.H{
		"status": true,
		"rooms": rooms,
	}))
}

func addStation(c *gin.Context) {
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

func editStation(c *gin.Context) {
	var station models.Station
	strSID := c.Param("sID")

	if err := c.BindJSON(&station); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, utils.ErrorMessage("bad request", http.StatusBadRequest))
		return
	}

	if exist := models.CheckStationExist(strSID); !exist {
		c.AbortWithStatusJSON(http.StatusOK, utils.ErrorMessage("station not found", http.StatusNotFound))
		return
	}

	if err := models.EditStation(strSID, &station); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, utils.ErrorMessage("internal error", http.StatusInternalServerError))
		return
	}

	c.AbortWithStatusJSON(http.StatusOK, utils.SuccessMessage(gin.H{
		"status": true,
	}))
}

func deleteStation(c *gin.Context) {
	strSID := c.Param("sID")

	if exist := models.CheckStationExist(strSID); !exist {
		c.AbortWithStatusJSON(http.StatusOK, utils.ErrorMessage("station not found", http.StatusNotFound))
		return
	}

	if err := models.RemoveStation(strSID); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, utils.ErrorMessage("internal error", http.StatusInternalServerError))
		return
	}

	c.AbortWithStatusJSON(http.StatusOK, utils.SuccessMessage(gin.H{
		"status": true,
	}))
}

func getStationInfo(c *gin.Context) {
	var station models.Station
	sID := c.Param("sID")

	models.GetStationInfoByID(sID, &station)

	if station.StationName == nil {
		c.AbortWithStatusJSON(http.StatusOK, utils.ErrorMessage("station not found", http.StatusNotFound))
		return
	}

	c.AbortWithStatusJSON(http.StatusOK, utils.SuccessMessage(gin.H{
		"status": true,
		"station": station,
	}))
}

func reorderStation(c *gin.Context) {
	var reorder models.ReorderStationModel
	if err := c.BindJSON(&reorder); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, utils.ErrorMessage("bad request", http.StatusBadRequest))
		return
	}

	if err := models.ReorderStationInHos(reorder); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, utils.ErrorMessage("internal error", http.StatusInternalServerError))
		return
	}

	c.AbortWithStatusJSON(http.StatusOK, utils.SuccessMessage(gin.H{
		"status": true,
	}))
}

func getStationByHospital(c *gin.Context) {
	hID := c.Param("hID")
	var station []models.MultiStation

	models.GetMultiStationInHos(hID, &station)

	c.AbortWithStatusJSON(http.StatusOK, utils.SuccessMessage(gin.H{
		"status": true,
		"station": station,
	}))
}

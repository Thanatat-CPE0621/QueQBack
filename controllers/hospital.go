package controllers

import "C"
import (
	"fmt"
	"time"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"gitlab.com/paiduay/queq-hospital-api/middlewares"
	"gitlab.com/paiduay/queq-hospital-api/models"
	"gitlab.com/paiduay/queq-hospital-api/utils"
	"gitlab.com/paiduay/queq-hospital-api/classes"
)

// RegisterHospitalEndpoints - to let main register these endpoints
func RegisterHospitalEndpoints(router *gin.RouterGroup) {
	hospitalRouter := router.Group("")
	// Annonymous Routes

	// Admin Routes
	hospitalRouter.Use(middlewares.LoginRequire())
	{
    hospitalRouter.GET("/", getAllHospitalList)

    hospitalRouter.GET("/station/:hID", getStationInHospital)
    hospitalRouter.GET("/station/:hID/info", getStationInHospitalWithBriefInfo)

    hospitalRouter.GET("/queue/day/:hID", getHospitalQueue)
    hospitalRouter.GET("/queue/weekly/:hID", getWeeklyHospitalQueue)
    hospitalRouter.GET("/queue/monthly/:hID", getMonthlyHospitalQueue)
    hospitalRouter.GET("/queue/yearly/:hID", getYearlyHospitalQueue)
    hospitalRouter.GET("/queue/during/day/:hID", getHospitalQueueDuringTheDay)
    hospitalRouter.GET("/queue/avg/all/:hID", getHospitalAllAverageQueueTime)
    hospitalRouter.GET("/queue/avg/day/:hID", getHospitalAverageQueueTime)
    hospitalRouter.GET("/queue/range/day/:hID", getHospitalQueueInDateRange)
    hospitalRouter.GET("/queue/period/month/:hID", getQueuePeriodInMonth)
    hospitalRouter.GET("/queue/period/week/:hID", getQueuePeriodInWeek)
	}

	hospitalRouter.Use(middlewares.SuperAdminRequired())
	{
    hospitalRouter.GET("/info/:hID", getHospitalInfomation)
    hospitalRouter.GET("/check/name", checkHospitalNameAvailable)
    hospitalRouter.GET("/room/all/:hID", getAllRoomInHospital)
    hospitalRouter.POST("/add", createNewHospital)
    hospitalRouter.POST("/edit/:hID", editHospitalInfo)
    // hospitalRouter.POST("/upload/image", uploadHospitalImage)
    hospitalRouter.POST("/reorder", reorderHostpital)
	}
}

func getAllHospitalList (c *gin.Context) {
  var hospitals []models.Hospital
  if err := models.GetHospitalList(&hospitals); err != nil {
    fmt.Println(err)
    c.AbortWithStatusJSON(http.StatusOK, utils.ErrorMessage("internal error", http.StatusInternalServerError))
		return
  }

  c.AbortWithStatusJSON(http.StatusOK, utils.SuccessMessage(gin.H{
    "status": true,
    "hospitals": hospitals,
  }))
}

func getStationInHospital (c *gin.Context) {
  hID := c.Param("hID")
  var stations []models.Station
  models.GetStationsListInHospital(hID, &stations)

  c.AbortWithStatusJSON(http.StatusOK, utils.SuccessMessage(gin.H{
    "status": true,
    "stations": stations,
  }))
}

func getStationInHospitalWithBriefInfo (c *gin.Context) {
	now := time.Now()
	hID := c.Param("hID")
	var sModel models.Station
	fromDate := c.Query("fromDate")
	toDate := c.Query("toDate")
	if fromDate == "" {
		fromDate = now.Format("2006-01-02")
	}
	if toDate == "" {
		toDate = now.Format("2006-01-02")
	}
	stations, err := sModel.GetStationinHoswithInfo(fromDate, toDate, hID)
	if err != nil {
		fmt.Println(err)
    c.AbortWithStatusJSON(http.StatusOK, utils.ErrorMessage("internal error", http.StatusInternalServerError))
		return
	}
	c.AbortWithStatusJSON(http.StatusOK, utils.SuccessMessage(gin.H{
    "status": true,
    "stations": stations,
  }))
}

func getHospitalQueue (c *gin.Context) {
  hID, err := strconv.ParseUint(c.Param("hID"), 10, 32)
  fromDate := c.Query("fromDate")
  toDate := c.Query("toDate")
  hospital := models.Hospital{
    HospitalID: uint32(hID),
  }

  if fromDate == "" || toDate == "" || err != nil {
    c.AbortWithStatusJSON(http.StatusOK, utils.ErrorMessage("bad request", http.StatusBadRequest))
		return
  }

  models.GetQueueAmountInHospital(&hospital, fromDate, toDate)

  c.AbortWithStatusJSON(http.StatusOK, utils.SuccessMessage(gin.H{
    "status": true,
    "queues": hospital,
  }))
}

func getWeeklyHospitalQueue (c *gin.Context) {
	hID := c.Param("hID")
  year := c.Query("year")
  weekNumber := c.Query("weekNumber")
	if year == "" || weekNumber == "" {
		c.AbortWithStatusJSON(http.StatusOK, utils.ErrorMessage("bad request", http.StatusBadRequest))
		return
	}
	var days []map[string]interface{}

	models.GetWeeklyQueueInHos(&days, year, weekNumber, hID)

	c.AbortWithStatusJSON(http.StatusOK, utils.SuccessMessage(gin.H{
    "status": true,
    "days": days,
  }))
}

func getMonthlyHospitalQueue (c *gin.Context) {
	hID := c.Param("hID")
	year := c.Query("year")
	if year == "" {
		c.AbortWithStatusJSON(http.StatusOK, utils.ErrorMessage("bad request", http.StatusBadRequest))
		return
	}
	var months []map[string]interface{}
	latestMonth :=  make(map[string]interface{})

	models.GetMonthlyQueueInHos(&months, &latestMonth, hID, year)

	c.AbortWithStatusJSON(http.StatusOK, utils.SuccessMessage(gin.H{
    "status": true,
    "the_lastest_month": latestMonth,
		"months": months,
  }))
}

func getYearlyHospitalQueue (c *gin.Context) {
	hID := c.Param("hID")
  yearStart := c.Query("yearStart")
  yearEnd := c.Query("yearEnd")
	if yearStart == "" || yearEnd == "" {
		c.AbortWithStatusJSON(http.StatusOK, utils.ErrorMessage("bad request", http.StatusBadRequest))
		return
	}
	var years []map[string]interface{}

	models.GetYearlyQueueInHos(&years, yearStart, yearEnd, hID)

	c.AbortWithStatusJSON(http.StatusOK, utils.SuccessMessage(gin.H{
    "status": true,
    "years": years,
  }))
}

func getHospitalQueueDuringTheDay (c *gin.Context) {
  hID, err := strconv.ParseUint(c.Param("hID"), 10, 32)
  fromDate := c.Query("fromDate")
  toDate := c.Query("toDate")
  hospital := models.Hospital{
    HospitalID: uint32(hID),
  }

  if fromDate == "" || toDate == "" || err != nil {
    c.AbortWithStatusJSON(http.StatusOK, utils.ErrorMessage("bad request", http.StatusBadRequest))
		return
  }

  models.GetQueueDuringTheDay(&hospital, fromDate, toDate)

  c.AbortWithStatusJSON(http.StatusOK, utils.SuccessMessage(gin.H{
    "status": true,
    "queues": hospital,
  }))
}

func getHospitalAllAverageQueueTime (c *gin.Context) {
	now := time.Now()
	hID := c.Param("hID")
	date := c.Query("date")
	var day map[string]interface{}
	var dayWeek map[string]interface{}
	var dayMonth map[string]interface{}

	if date == "" {
		date = now.Format("2006-01-02")
	}

	models.GetAvgAllDay(hID, date, &day)
	models.GetAvgAllDayWeek(hID, date, &dayWeek)
	models.GetAvgAllDayMonth(hID, date, &dayMonth)

	c.AbortWithStatusJSON(http.StatusOK, utils.SuccessMessage(gin.H{
    "day": day,
    "week": dayWeek,
    "month": dayMonth,
  }))
}

func getHospitalAverageQueueTime (c *gin.Context) {
	now := time.Now()
	hID := c.Param("hID")
	var data map[string]interface{}
	fromDate := c.Query("fromDate")
	toDate := c.Query("toDate")
	if fromDate == "" {
		fromDate = now.Format("2006-01-02")
	}
	if toDate == "" {
		toDate = now.Format("2006-01-02")
	}

	models.GetAverageDayQueueingTime(hID, fromDate, toDate, &data)

	c.AbortWithStatusJSON(http.StatusOK, utils.SuccessMessage(gin.H{
    "time": data,
  }))
}

func getHospitalQueueInDateRange (c *gin.Context) {
	now := time.Now()
	var data classes.Week
	hID := c.Param("hID")
	fDate := c.Query("fromDate")
	tDate := c.Query("toDate")
	if fDate == "" || tDate == "" {
		fDate = now.Format("20016-01-02")
		tDate = now.Format("20016-01-02")
	}

	if err := models.GetHospitalQueueInDateRange(hID, fDate, tDate, &data); err != nil {
		fmt.Println(err)
		c.AbortWithStatusJSON(http.StatusOK, utils.ErrorMessage("internal error", http.StatusInternalServerError))
		return
	}

	c.AbortWithStatusJSON(http.StatusOK, utils.SuccessMessage(gin.H{
    "status": true,
    "week": data,
  }))
}

func getQueuePeriodInMonth (c *gin.Context) {
	now := time.Now()
	hID := c.Param("hID")
	date := c.Query("date")
	var hours models.QueueInterface
	if date == "" {
		date = now.Format("2006-01")
	}

	hours.GetperiodQMonth(hID, date)

	c.AbortWithStatusJSON(http.StatusOK, utils.SuccessMessage(gin.H{
    "hours": hours,
  }))
}

func getQueuePeriodInWeek (c *gin.Context) {
	hID := c.Param("hID")
	year := c.Query("year")
	weekNumber := c.Query("weekNumber")
	var hours models.QueueInterface
	if year == "" || weekNumber == "" {
		c.AbortWithStatusJSON(http.StatusOK, utils.ErrorMessage("bad request", http.StatusBadRequest))
		return
	}

	hours.GetperiodQWeek(hID, year, weekNumber)

	c.AbortWithStatusJSON(http.StatusOK, utils.SuccessMessage(gin.H{
    "hours": hours,
  }))
}

func getHospitalInfomation (c *gin.Context) {
	hID := c.Param("hID")
	var hospital models.Hospital
	if hID == "" {
		c.AbortWithStatusJSON(http.StatusOK, utils.ErrorMessage("bad request", http.StatusBadRequest))
		return
	}
	hospital.GetInformation(hID)
	c.AbortWithStatusJSON(http.StatusOK, utils.SuccessMessage(gin.H{
		"status": true,
    "hospital": hospital,
  }))
}

func checkHospitalNameAvailable (c *gin.Context) {
  var exist bool
  hName := c.Query("name")
  if hName == "" {
    c.AbortWithStatusJSON(http.StatusOK, utils.ErrorMessage("bad request", http.StatusBadRequest))
		return
  }

  if exist = models.CheckIfHospitalNameExist(hName); !exist {
    c.AbortWithStatusJSON(http.StatusOK, utils.SuccessMessage(gin.H{
      "hospital_name": hName,
  		"available": exist,
  	}))
    return
  }

  c.AbortWithStatusJSON(http.StatusOK, utils.SuccessMessage(gin.H{
    "hospital_name": hName,
		"available": exist,
	}))
}

func getAllRoomInHospital (c *gin.Context) {
    hID := c.Param("hID")
    var rooms []models.Room
    models.GetAllRoomInHospital(hID, &rooms)

    c.AbortWithStatusJSON(http.StatusOK, utils.SuccessMessage(gin.H{
  		"status": true,
  		"rooms": rooms,
  	}))
}

func createNewHospital (c *gin.Context) {
  var hospital models.Hospital
  if err := c.BindJSON(&hospital); err != nil {
    fmt.Println(err)
    c.AbortWithStatusJSON(http.StatusOK, utils.ErrorMessage("bad request", http.StatusBadRequest))
		return
  }

  if exist := models.CheckIfHospitalNameExist(*hospital.HospitalName); !exist {
    c.AbortWithStatusJSON(http.StatusOK, utils.ErrorMessage("hospital already exist", http.StatusBadRequest))
    return
  }

  if err := models.CreateNewHospital(&hospital); err != nil {
    fmt.Println(err)
		c.AbortWithStatusJSON(http.StatusOK, utils.ErrorMessage("internal error", http.StatusInternalServerError))
		return
  }

  c.AbortWithStatusJSON(http.StatusOK, utils.SuccessMessage(gin.H{
    "status": true,
    "hospital": hospital,
  }))
}

func editHospitalInfo (c *gin.Context) {
	hID := c.Param("hID")
	var hospital models.Hospital
	if err := c.BindJSON(&hospital); err != nil {
    fmt.Println(err)
    c.AbortWithStatusJSON(http.StatusOK, utils.ErrorMessage("bad request", http.StatusBadRequest))
		return
  }

	if err := hospital.EditHospitalInformation(hID); err != nil {
		fmt.Println(err)
		c.AbortWithStatusJSON(http.StatusOK, utils.ErrorMessage("internal error", http.StatusInternalServerError))
		return
	}

	c.AbortWithStatusJSON(http.StatusOK, utils.SuccessMessage(gin.H{
    "status": true,
  }))
}

func uploadHospitalImage (c *gin.Context) {

}

func reorderHostpital (c *gin.Context) {
	var hospital []models.ReorderHospitalList
	if err := c.BindJSON(&hospital); err != nil {
		fmt.Println(err)
		c.AbortWithStatusJSON(http.StatusOK, utils.ErrorMessage("bad request", http.StatusBadRequest))
		return
	}
	if err := models.ReorderHospitals(hospital); err != nil {
		fmt.Println(err)
		c.AbortWithStatusJSON(http.StatusOK, utils.ErrorMessage("internal error", http.StatusInternalServerError))
		return
	}
	c.AbortWithStatusJSON(http.StatusOK, utils.SuccessMessage(gin.H{
    "status": true,
  }))
}

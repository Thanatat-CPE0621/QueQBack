package controllers

import (
  "fmt"
  "strconv"
  "net/http"
  "github.com/gin-gonic/gin"

  "gitlab.com/paiduay/queq-hospital-api/utils"
  "gitlab.com/paiduay/queq-hospital-api/models"
  "gitlab.com/paiduay/queq-hospital-api/middlewares"
)

// RegisterStaffEndpoints - to let main register these endpoints
func RegisterStaffEndpoints(router *gin.RouterGroup) {
  userRouter := router.Group("")
  userRouter.Use(middlewares.TokenCheck())
  userRouter.Use(middlewares.LoginRequire())

	userRouter.GET("", getStaffList)
}

func getStaffList(c *gin.Context) {
	var staffs []models.Staff
  size, err := strconv.Atoi(c.Query("size"))
  if err != nil {
    c.AbortWithStatusJSON(http.StatusOK, utils.ErrorMessage("bad request", http.StatusBadRequest))
    return
  }
  page, err := strconv.Atoi(c.Query("page"))
  if err != nil {
    c.AbortWithStatusJSON(http.StatusOK, utils.ErrorMessage("bad request", http.StatusBadRequest))
    return
  }
  hID := c.Query("hid")
  rID := c.Query("rid")

  if err := models.GetStaffList(&staffs, size, page, hID, rID); err != nil {
    fmt.Println(err)
    c.AbortWithStatusJSON(http.StatusOK, utils.ErrorMessage("internal error", http.StatusInternalServerError))
    return
  }
  c.AbortWithStatusJSON(http.StatusOK, gin.H{
    "staffs": staffs,
  })
}

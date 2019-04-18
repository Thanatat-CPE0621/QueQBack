package middlewares

import (
	"fmt"
	"net/http"
	"github.com/gin-gonic/gin"
	"gitlab.com/paiduay/queq-hospital-api/models"
	"gitlab.com/paiduay/queq-hospital-api/utils"
)

// CORSMiddleware - a middleware to allow cors headers
func CORSMiddleware () gin.HandlerFunc {
	return func (c *gin.Context) {
		c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE, PATCH")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatusJSON(http.StatusOK, utils.ErrorMessage("unauthorized", http.StatusUnauthorized))
		} else {
			c.Next()
		}
	}
}

func TokenCheck () gin.HandlerFunc {
	return func (c *gin.Context) {
		auth := c.Request.Header.Get("Authorization")
		if auth != "" {
			if auth != "Basic JOl14hhlf0ia0W1fo4tlBZIBPv1WNuYtnSfD6oPF2piw8HYXuGokuTvA97PX24eWh9cgJrOvBC6mE1QgNyNbjEWQnAqF4MgtLRClLO644h4NtEo50W2MtNWuhex5JHC8" {
				c.AbortWithStatusJSON(http.StatusOK, utils.ErrorMessage("unauthorized", http.StatusUnauthorized))
			} else {
				c.Next()
			}
		} else {
			c.AbortWithStatusJSON(http.StatusOK, utils.ErrorMessage("unauthorized", http.StatusUnauthorized))
		}
	}
}

func LoginRequire () gin.HandlerFunc {
	return func (c *gin.Context) {
		userToken := c.Request.Header.Get("userToken")
		if userToken != "" {
			staff := models.Staff{}
			if err := models.GetHospitalIDbyStaffToken(&staff, userToken); err != nil {
				fmt.Println(err)
				c.AbortWithStatusJSON(http.StatusOK, utils.ErrorMessage("Not found!", http.StatusNotFound))
				return
			}
			if staff.StaffID != 0 {
				c.Next()
			} else {
				c.AbortWithStatusJSON(http.StatusOK, utils.ErrorMessage("unauthorized", http.StatusUnauthorized))
				return
			}
		} else {
			c.AbortWithStatusJSON(http.StatusOK, utils.ErrorMessage("unauthorized", http.StatusUnauthorized))
			return
		}
	}
}

func superAdminRequired () gin.HandlerFunc {
	return func (c *gin.Context) {
		userToken := c.Request.Header.Get("userToken")
		if userToken != "" {
			staff := models.Staff{}
			if err := models.GetHospitalIDbyStaffToken(&staff, userToken); err != nil {
				fmt.Println(err)
				c.AbortWithStatusJSON(http.StatusOK, utils.ErrorMessage("internal error", http.StatusInternalServerError))
			}
			if staff.StaffID != 0 {
				c.Next()
			}
		} else {
			c.AbortWithStatusJSON(http.StatusOK, utils.ErrorMessage("unauthorized", http.StatusUnauthorized))
		}
	}
}

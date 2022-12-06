package contreller

import (
	"github.com/abdumalik92/identification/internal/models"
	"github.com/abdumalik92/identification/internal/service"
	"github.com/abdumalik92/identification/internal/utils"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

func OtpServer(c *gin.Context) {
	var (
		request models.OtpRequest
	)

	if err := c.ShouldBindJSON(&request); err != nil {
		log.Println("OtpServer func error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"reason": "Что-то пошло не так", // something went wrong
		})
		return
	}

	// checks phoneNum length
	if err := utils.PhoneNumCheck(request.PhoneNum); err != nil {
		log.Println("phone num length error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"reason": err.Error(), // something went wrong
		})
		return
	}

	if err := service.OtpServer(request); err != nil {
		var statusCode int

		if err.Error() == "Количество запросов на подтверждение превышено" {
			statusCode = 400
		} else {
			if _, err := strconv.Atoi(err.Error()); err == nil {
				log.Println("is number")
				statusCode = 403
			} else {
				statusCode = 400
			}
		}
		log.Println("response statusCode:", statusCode, " error:", err.Error(), " code:", request.Code, " phone:", request.PhoneNum)
		c.JSON(statusCode, gin.H{
			"reason": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"reason": "success",
	})
}

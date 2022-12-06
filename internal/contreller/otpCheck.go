package contreller

import (
	"github.com/abdumalik92/identification/internal/models"
	"github.com/abdumalik92/identification/internal/service"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func OtpCheck(c *gin.Context) {
	var (
		request  models.OtpCheckRequest
		response models.OtpCheckResp
	)

	if err := c.ShouldBindJSON(&request); err != nil {
		log.Println("OtpCheck func bind error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"reason": "Что-то пошло не так",
		})
		return
	}
	// checking info sent to this route to be same as in hashSum

	if err := service.OtpCheck(request, &response); err != nil {
		log.Println("OtpCheck func service error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"reason": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

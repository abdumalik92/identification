package routes

import (
	"fmt"
	"github.com/abdumalik92/identification/internal/contreller"
	"github.com/abdumalik92/identification/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/natefinch/lumberjack"
	"net/http"

	"io"
	"log"
	"os"
	"time"
)

func RunRoutes() {
	r := gin.Default()
	r.Use(CORSMiddleware())

	//f, err := os.Create("gin.log")
	f, err := os.OpenFile(utils.AppSettings.AppParams.LogFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)

	if err != nil {
		fmt.Println("file create error", err.Error())
		return
	}

	logger := &lumberjack.Logger{
		Filename:   f.Name(),
		MaxSize:    10, // megabytes
		MaxBackups: 100,
		MaxAge:     28,   // days
		Compress:   true, // disabled by default
	}

	log.SetOutput(logger)

	gin.DefaultWriter = io.MultiWriter(logger, os.Stdout)

	// LoggerWithFormatter middleware will write the logs to gin.DefaultWriter
	// By default gin.DefaultWriter = os.Stdout
	r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {

		// your custom format
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))
	r.Use(gin.Recovery())

	r.POST("/otp", contreller.OtpServer)
	r.POST("/otp_check", contreller.OtpCheck)
	r.POST("/send_order", CheckToken, contreller.SendOrderToCFT)
	r.POST("/humo_online", contreller.SendOrderToCFT)

	_ = r.Run(utils.AppSettings.AppParams.PortRun)
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET,PUT,POST,DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "*")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
		}

		c.Next()
	}
}

func CheckToken(c *gin.Context) {

	token := c.GetHeader("Authorization")

	claims := utils.GetStruct(token, c.Writer)
	if claims == nil {
		log.Println("func CheckToken Token is expired")
		c.JSON(http.StatusUnauthorized, gin.H{"reason": "Ваша сессия устарела, пожалуйста попробуйте снова"})
		c.Abort()
		return
	}
	return
}

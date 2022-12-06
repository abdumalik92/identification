package service

import (
	"errors"
	"fmt"
	"github.com/abdumalik92/identification/internal/models"
	"github.com/abdumalik92/identification/internal/repository"
	"github.com/abdumalik92/identification/internal/utils"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func OtpServer(request models.OtpRequest) error {
	var otp int
	if err := repository.CheckNum(request); err != nil {
		return err
	}
	// generating random number (1000-9999) for OTP

	rand.Seed(time.Now().UnixNano())
	min := 1000
	max := 9999
	otp = rand.Intn(max-min+1) + min

	otp = 1000 + rand.Intn(8999)

	strOTP := strconv.Itoa(otp)

	if err := repository.OtpServer(request, otp); err != nil {
		return err
	}

	otpText := "Ваш код подтверждения: " + strOTP + ". Никому не говорите код!"

	//replace empty spaces. prepare url for [GET] request
	otpText = strings.Replace(otpText, " ", "%20", -1)

	log.Println(otp)

	phoneNum := request.Code + request.PhoneNum

	// API for sending (sms)OTP to user's phone number given in request  // number must have prefix i.e. +992, . . etc.
	urls := fmt.Sprintf(utils.AppSettings.SendSMSURL.Url, phoneNum, otpText)

	log.Println("sms url:", urls)
	client := &http.Client{}
	//proxyUrl, err := url.Parse("http://192.168.0.8:4480")
	//
	//if err != nil {
	//	log.Println("proxy error:", err.Error())
	//	return err
	//}
	//
	//client.Transport = &http.Transport{Proxy:http.ProxyURL(proxyUrl)}
	_, err := client.Get(urls)

	if err != nil {
		log.Println("OtpServer func sms send error", err.Error())
		return errors.New("Что-то пошло не так...")
	}

	return nil
}

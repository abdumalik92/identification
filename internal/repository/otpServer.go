package repository

import (
	"C"
	"context"
	"database/sql"
	"errors"
	"github.com/abdumalik92/identification/internal/db"
	"github.com/abdumalik92/identification/internal/models"
	"log"
	"time"
)

var sendOTP = "begin ibs.z$remote_identif_service_lib.remote_identif_send_otp(:phone, :otp, :errorcode,:errordescription); end;"

func OtpServer(request models.OtpRequest, otp int) error {
	var (
		otpResp       models.MainModel
		phoneWithCode = request.Code + request.PhoneNum
	)
	log.Println("send to CFT ", phoneWithCode, " ", otp)
	dbConn := db.GetOracleDB()
	ctx, _ := context.WithTimeout(context.Background(), 20*time.Second)

	_, err := dbConn.ExecContext(ctx, sendOTP,
		phoneWithCode,
		otp,
		sql.Out{Dest: &otpResp.ErrorCode},
		sql.Out{Dest: &otpResp.ErrorDescription},
	)
	if err != nil {
		log.Println("OtpServer error = ", err)
		return errors.New("Что-то пошло не так...")
	}

	if otpResp.ErrorCode != 0 {
		log.Println("OtpServer CFT ERROR_CODE =  ", otpResp.ErrorCode, " ERROR_DESCRIPTION = ", otpResp.ErrorDescription)
		return errors.New(otpResp.ErrorDescription)
	}

	return nil
}

var checkNum = "begin ibs.z$remote_identif_service_lib.remote_identif_send_otp(:phone, :errorcode, :errordescription); end;"

func CheckNum(request models.OtpRequest) error {
	var (
		otpResp       models.MainModel
		phoneWithCode = request.Code + request.PhoneNum
	)
	log.Println("check num in CFT ", phoneWithCode)
	dbConn := db.GetOracleDB()
	ctx, _ := context.WithTimeout(context.Background(), 20*time.Second)

	_, err := dbConn.ExecContext(ctx, sendOTP,
		phoneWithCode,
		sql.Out{Dest: &otpResp.ErrorCode},
		sql.Out{Dest: &otpResp.ErrorDescription},
	)
	if err != nil {
		log.Println("CheckNum error = ", err)
		return errors.New("Сервис временно не доступен")
	}

	if otpResp.ErrorCode != 0 {
		log.Println("CheckNum CFT ERROR_CODE =  ", otpResp.ErrorCode, "ERROR_DESCRIPTION = ", otpResp.ErrorDescription)
		return errors.New(otpResp.ErrorDescription)
	}

	return nil
}

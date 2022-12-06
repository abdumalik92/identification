package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/abdumalik92/identification/internal/db"
	"github.com/abdumalik92/identification/internal/models"
	"log"
	"time"
)

var checkOTP = "begin ibs.z$remote_identif_service_lib.remote_identif_check_otp(:phone, :otp, :errorcode, :errordescription); end;"

func OtpCheck(request models.OtpCheckRequest) error {

	var errorCode int
	var errorDescription string
	phoneNum := request.Code + request.PhoneNum

	dbConn := db.GetOracleDB()
	ctx, _ := context.WithTimeout(context.Background(), 20*time.Second)

	_, err := dbConn.ExecContext(ctx, checkOTP,
		phoneNum,
		request.Otp,
		sql.Out{Dest: &errorCode},
		sql.Out{Dest: &errorDescription},
	)

	if err != nil {
		log.Println("OtpCheck error = ", err)
		return errors.New("Что-то пошло не так")
	}

	if errorCode != 0 {
		log.Println("OtpCheck CFT ERROR_CODE =  ", errorCode, " ERROR_DESCRIPTION = ", errorDescription)
		return errors.New("Вы ввели неправильный или просроченный код подтверждения")
	}

	return nil
}

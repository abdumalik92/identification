package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/abdumalik92/identification/internal/db"
	"log"
	"time"
)

var sendOrderToCFTSQL = "begin ibs.z$remote_identif_service_lib.remote_identif_request(:product, :client_id, :file_link, :errorcode, :errordescription); end;"

func SendOrderToCFT(product string, clientID string, file_link string) error {
	dbConn := db.GetOracleDB()
	var errorCode int
	var errorDescription string

	ctx, _ := context.WithTimeout(context.Background(), 20*time.Second)

	_, err := dbConn.ExecContext(ctx, sendOrderToCFTSQL,
		product,
		clientID,
		file_link,
		sql.Out{Dest: &errorCode},
		sql.Out{Dest: &errorDescription},
	)

	if err != nil {
		log.Println("SendOrderToCFT error = ", err)
		return errors.New("Что-то пошло не так...")
	}

	if errorCode != 0 {
		log.Println("SendOrderToCFT CFT ERROR_CODE =  ", errorCode, " ERROR_DESCRIPTION = ", errorDescription)
		return errors.New(errorDescription)
	}
	return nil
}

func Ping() {
	dbConn := db.GetOracleDB()

	var _, err = dbConn.Exec("SELECT * FROM DUAL")

	if err != nil {
		log.Println("error on db ping ", err)
		return
	}
	log.Println("db ping success")
}

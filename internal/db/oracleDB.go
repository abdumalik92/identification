package db

import (
	"database/sql"
	"fmt"
	"github.com/abdumalik92/identification/internal/utils"
	_ "github.com/godror/godror"
	"log"
)

var (
	ODB *sql.DB
	err error
)

// InitOracleDB inits oracle db
func InitOracleDB() error {

	oracleDBParams := utils.AppSettings.OracleCFTDbParams

	oracleConnectionString := fmt.Sprintf("%s/%s@%s", oracleDBParams.User, oracleDBParams.Password, oracleDBParams.Server)
	fmt.Println("connectionString", oracleConnectionString)
	ODB, err = sql.Open("godror", oracleConnectionString)

	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

//GetOracleDB2 return connection to oracle db
func GetOracleDB() *sql.DB {
	return ODB
}

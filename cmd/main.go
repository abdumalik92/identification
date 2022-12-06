package main

import (
	"github.com/abdumalik92/identification/internal/db"
	"github.com/abdumalik92/identification/internal/repository"
	"github.com/abdumalik92/identification/internal/routes"
	"github.com/abdumalik92/identification/internal/utils"
	_ "github.com/godror/godror"
	"time"
)

func runScheduler() {
	duration := time.Duration(utils.AppSettings.OracleCFTDbParams.ConnectionCheckTime)
	timer := time.NewTicker(time.Minute * duration)
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			repository.Ping()
		}
	}
}

func main() {
	utils.ReadSettings()

	if err := db.InitOracleDB(); err != nil {
		panic("error on db initializing")
	}

	go runScheduler()

	routes.RunRoutes()
}

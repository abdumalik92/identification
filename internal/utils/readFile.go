package utils

import (
	"encoding/json"
	"fmt"
	"github.com/abdumalik92/identification/internal/models"
	"log"
	"os"
)

var (
	// have settings taken from json file
	AppSettings models.Settings
)

// Read Settings and assigns to AppSettings variable
func ReadSettings() {
	fmt.Println("starting reading settings file")
	configFile, err := os.Open("./internal/configs/settings.json")

	defer configFile.Close()

	if err != nil {
		log.Fatal("Couldn't open config file", err.Error())
	}

	fmt.Println("starting decoding settings file")
	err = json.NewDecoder(configFile).Decode(&AppSettings)

	if err != nil {
		log.Fatal("Couldn't decode settings json file", err.Error())
	}

	log.Println(AppSettings)
	return
}

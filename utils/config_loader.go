package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"

	"gitlab.com/paiduay/queq-hospital-api/config"
)

// ConfigLoader - Load config from configuration/json directory
// mode: prod, dev
func ConfigLoader(mode string) error {
	var configsFile config.GeneralConfigs
	var databaseConfigFile config.DatabaseConfiguration

	if err := databaseConfigurationJSONLoader(&databaseConfigFile); err != nil {
		return errors.New("failed to load database config")
	}
	if err := generalConfigsJSONLoader(&configsFile); err != nil {
		return errors.New("failed to load config")
	}

	config.SetDatabaseConfiguration(databaseConfigFile, mode)
	config.SetGeneralConfigs(configsFile, mode)

	return nil
}

// databaseConfigurationJSONLoader - Load database configuration from configuration/json/database.json
func databaseConfigurationJSONLoader(databaseConfigFile *config.DatabaseConfiguration) error {
	databaseJSONByte, err := ioutil.ReadFile("config/json/database.json")
	if err != nil {
		fmt.Println(err)
		return errors.New("parse json error")
	}
	if err := json.Unmarshal(databaseJSONByte, &databaseConfigFile); err != nil {
		fmt.Println(err)
		return errors.New("parse json error")
	}
	return nil
}

func generalConfigsJSONLoader(configFile *config.GeneralConfigs) error {
	jsonConfigByte, err := ioutil.ReadFile("config/json/general_conf.json")
	if err != nil {
		fmt.Println(err)
		return errors.New("parse json error")
	}
	if err := json.Unmarshal(jsonConfigByte, &configFile); err != nil {
		fmt.Println(err)
		return errors.New("parse json error")
	}
	return nil
}

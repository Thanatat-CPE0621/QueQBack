package config

// GeneralConfigs - for parsing general_configs.json
type GeneralConfigs struct {
	Development GeneralConfigBody `json:"development"`
	Production  GeneralConfigBody `json:"production"`
	Default     GeneralConfigBody `json:"default"`
}

// GeneralConfigBody is struct for store Configs using in app
type GeneralConfigBody struct {
	JWT             jwtConfigs `json:"jwt"`
	PasswordSalt    string     `json:"passwordSalt"`
	TokenSalt		    string     `json:"tokenSalt"`
	APIVersion      string     `json:"apiVersion"`
	StaticDirectory string     `json:"staticDirectory"`
}

type jwtConfigs struct {
	Secret   string `json:"secret"`
	Lifetime string `json:"lifetime"`
}

// Configs - current configs
var Configs GeneralConfigBody

// SetGeneralConfigs - set current configs
func SetGeneralConfigs(readConfig GeneralConfigs, mode string) {
	Configs = GeneralConfigBody{}
	if mode == "prod" {
		Configs = readConfig.Production
	} else if mode == "dev" {
		Configs = readConfig.Development
	} else {
		Configs = readConfig.Default
	}
}

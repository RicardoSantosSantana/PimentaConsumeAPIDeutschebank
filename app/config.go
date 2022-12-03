package main

import (
	"os"
	"strings"
)

type DbConnection struct {
	DatabaseName string `json:"databasename"`
	ServerName   string `json:"servername"`
	UserName     string `json:"username"`
	Password     string `json:"password"`
	Port         string `json:"port"`
}

type ApiAuth struct {
	Client_id     string `json:"client_id"`
	Client_secret string `json:"client_secret"`
	Code          string `json:"code"`
	Redirect_uri  string `json:"redirect_uri"`
}
type ApiEndPoint struct {
	Method string `json:"method"`
	Uri    string `json:"uri"`
}

type DataSettings struct {
	Database DbConnection `json:"database"`
	Auth     ApiAuth      `json:"auth"`
	GetToken ApiEndPoint  `json:"endpoint"`
}

// get variable content if exist
func environment(value_default string, env string) string {

	enviromentVariable := strings.Trim(os.Getenv(strings.ToUpper(env)), " ")

	if enviromentVariable != "" {
		return os.Getenv(env)
	}
	return value_default
}

func Settings() DataSettings {

	return DataSettings{
		DbConnection{
			DatabaseName: environment("bank_transactions", "MYSQL_DATABASE"),
			ServerName:   environment("127.0.0.1", "MYSQL_SERVER"),
			UserName:     environment("pimenta", "MYSQL_USER"),
			Password:     environment("pimenta001", "MYSQL_PASSWORD"),
			Port:         environment("3399", "MYSQL_PORT"),
		},
		ApiAuth{
			Client_id:     environment("39975680-2987-4cae-9463-8b93b39129a3", "CLIENT_ID"),
			Client_secret: environment("AK3BkpPgVzW9d3_qE-r_dWG1QPQLPc5SfV2PiQ-h4l0Z3iviqC_oITM7Z-wjtXmGmyts7uTORprAC_Gh9pb8sI4", "CLIENT_SECRET"),
			Code:          environment("VTE48m", "CODE"),
			Redirect_uri:  environment("http://pimenta:3000/", "REDIRECT_URL"),
		},
		ApiEndPoint{
			Method: "POST",
			Uri:    "https://simulator-api.db.com/gw/oidc/token",
		},
	}
}

package config

import "os"

type ConfigType struct {
	Uri string
	IsDebugMode bool
}

var Config ConfigType

func init(){
	Config = ConfigType{
		Uri: "8080",
		IsDebugMode: true,
	}
}

func getEnv(key string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return ""
}
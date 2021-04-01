package config

import "os"

type ConfigType struct {
	Uri string
	IsDebugMode bool
	Port string
}

var Config ConfigType

func init(){
	Config = ConfigType{
		Uri: getEnv("URI"),
		IsDebugMode: true,
		Port: getEnv("PORT"),
	}
}

func getEnv(key string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return ""
}
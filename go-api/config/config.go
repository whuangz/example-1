package config

import (
	"crypto/rsa"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"github.com/dgrijalva/jwt-go"
)

var (
	URI                string
	IS_DEBUG_MODE      bool
	PORT               string
	JWT_PRIVATE_KEY    *rsa.PrivateKey
	JWT_PUBLIC_KEY     *rsa.PublicKey
	JWT_REFRESH_SECRET string
	ACCESS_TOKEN_EXP   int64
	REFRESH_TOKEN_EXP  int64

	REDIS_HOST string
	REDIS_PORT string

	HANDLER_TIMEOUT int64
)

func init() {
	URI = getEnv("URI", "")
	IS_DEBUG_MODE = true
	PORT = getEnv("PORT", ":8080")

	var err error
	handlerTimeout := getEnv("HANDLER_TIMEOUT", "5")
	HANDLER_TIMEOUT, err = strconv.ParseInt(handlerTimeout, 0, 64)
	if err != nil {
		log.Fatalf("could not parse HANDLER_TIMEOUT as int: %v", err)
	}
	initJwtKey()
	initRedis()

}

func initJwtKey() {
	privKeyFile := getEnv("PRIV_KEY_FILE", "")
	priv, err := ioutil.ReadFile(privKeyFile)

	if err != nil {
		log.Fatalf("could not read private key pem file: %v\n", err)
	}

	JWT_PRIVATE_KEY, err = jwt.ParseRSAPrivateKeyFromPEM(priv)
	if err != nil {
		log.Fatalf("could not parse private key: %v\n", err)
	}

	pubKeyFile := getEnv("PUB_KEY_FILE", "")
	pub, err := ioutil.ReadFile(pubKeyFile)
	if err != nil {
		log.Fatalf("could not read public key pem file: %v\n", err)

	}

	JWT_PUBLIC_KEY, err = jwt.ParseRSAPublicKeyFromPEM(pub)
	if err != nil {
		log.Fatalf("could not parse public key: %v\n", err)
	}

	JWT_REFRESH_SECRET = getEnv("REFRESH_SECRET", "")

	AccTokenExp := getEnv("ACCESS_TOKEN_EXP", "")
	ACCESS_TOKEN_EXP, err = strconv.ParseInt(AccTokenExp, 0, 64)
	if err != nil {
		log.Fatalf("could not parse ID_TOKEN_EXP as int: %v", err)
	}

	refreshTokenExp := getEnv("REFRESH_TOKEN_EXP", "")
	REFRESH_TOKEN_EXP, err = strconv.ParseInt(refreshTokenExp, 0, 64)
	if err != nil {
		log.Fatalf("could not parse REFRESH_TOKEN_EXP as int: %v", err)
	}
}

func initRedis() {
	REDIS_HOST = getEnv("REDIS_HOST", "")
	REDIS_PORT = getEnv("REDIS_PORT", "")
}

func getEnv(key string, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}

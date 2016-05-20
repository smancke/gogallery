package service

import (
	"log"
	"os"
	"strconv"
)

var defaults = map[string]string{
	"address":                ":5005",
	"galleryDir":             "/tmp/gallery",
	"htmlDir":                "./html",
	"cookieName":             "okmsdc",
	"session_secret":         "secretsecretsecretsecretsecretse",
	"sessionLifetimeMinutes": "180",
}

func Cfg(key string) string {
	if env := os.Getenv(key); env != "" {
		return env
	}
	if value, exist := defaults[key]; exist {
		return value
	}
	log.Panicf("missing cfg value for key=%v", key)
	panic("not reachable")
}

func CfgInt(key string) int {
	intVal, error := strconv.Atoi(Cfg(key))
	if error != nil {
		log.Panicf("not an int value for key=%v: %v", key, Cfg(key))
	}
	return intVal
}

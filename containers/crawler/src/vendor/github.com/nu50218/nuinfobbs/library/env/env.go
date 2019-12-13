package env

import "os"

import "log"

func AssertAndGetEnv(key string) string {
	value, exist := os.LookupEnv(key)
	if !exist {
		log.Fatalln(key, "does not exist.")
	}
	return value
}

package tools

import (
	"github.com/joho/godotenv"
	"os"
	"log"
)

func Init_env() {
	if err := godotenv.Load(); err != nil {
		log.Fatalln(err)
	}
}

func GetEnvValue(key, def string) string {
	value, ok := os.LookupEnv(key)
	if ok {
		return value
	}
	log.Println(key, ": default value used.")
	return def
}
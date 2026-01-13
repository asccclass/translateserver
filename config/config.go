package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	CheckInterval int
	SmtpHost      string
	SmtpPort      string
	SmtpUser      string
	SmtpPass      string
)

func Init() {
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if err := godotenv.Load(currentDir + "/envfile"); err != nil {
		fmt.Println(err.Error())
		return
	}

	intervalStr := os.Getenv("CHECK_INTERVAL")
	if intervalStr == "" {
		CheckInterval = 30
	} else {
		CheckInterval, err = strconv.Atoi(intervalStr)
		if err != nil {
			CheckInterval = 30
		}
	}

	SmtpHost = os.Getenv("SMTP_HOST")
	SmtpPort = os.Getenv("SMTP_PORT")
	SmtpUser = os.Getenv("SMTP_USER")
	SmtpPass = os.Getenv("SMTP_PASS")
}

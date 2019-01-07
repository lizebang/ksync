package main

import (
	"os"

	gocron "github.com/robfig/cron"

	"github.com/lizebang/ksync/log"
)

func main() {
	client := NewClient()
	err := client.Init()
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	client.Run()

	cron := gocron.New()
	err = cron.AddJob("0 0 0 * * *", client)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	cron.Run()
}

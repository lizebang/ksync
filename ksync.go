package main

import (
	"os"

	"github.com/robfig/cron"

	"github.com/lizebang/ksync/pkg/log"
)

func main() {
	cl := NewClient()
	err := cl.Init()
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}

	cr := cron.New()
	err = cr.AddJob("0 0 0 * * *", cl)
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
	cr.Run()
}

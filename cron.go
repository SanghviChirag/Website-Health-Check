package main

import (
	"gopkg.in/robfig/cron.v2"
)

func setCron(sec int, checkLink func(Website), website Website) {
	c := cron.New()
	// c.AddFunc("*/"+string(sec)+" * * * * *", checkLink(website))
	c.Start()
}

package main

import (
	"fmt"
	"strconv"

	"gopkg.in/robfig/cron.v2"
)

func setCron(website Website) int {
	c := cron.New()
	interval := strconv.Itoa(website.CheckInterval)
	c.AddFunc("*/"+interval+" * * * * *", func() { checkLink(website) })
	c.Start()

	fmt.Println("CRON")
	return 1
}

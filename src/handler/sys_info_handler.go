package handler

import (
	"alming_backend/src/engine"
	"alming_backend/src/service"
	"fmt"
	"time"
)

func GetInfo(c *engine.Context) {
	sysInfo := service.GetSysInfo()
	c.WriteJSON(sysInfo)
}
func Timer(c *engine.Context) {
	ticker := time.NewTicker(time.Second * 1)
	go func() {
		for s := range ticker.C {
			fmt.Println(s.Format("2006-01-02 15:04:05"))
		}
	}()
	ticker.Stop()
	fmt.Println(ticker)
}

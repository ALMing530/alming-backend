package handler

import (
	"alming_backend/src/engine"
	"alming_backend/src/service"
)

func GetInfo(c *engine.Context) {
	sysInfo := service.GetSysInfo()
	c.WriteJSON(sysInfo)
}

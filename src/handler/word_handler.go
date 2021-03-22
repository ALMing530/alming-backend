package handler

import (
	"alming_backend/src/engine"
	"alming_backend/src/service"
)

func GetWords(c *engine.Context) {
	service.GetWords(c)
}

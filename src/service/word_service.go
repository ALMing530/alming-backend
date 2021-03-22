package service

import (
	"alming_backend/src/engine"
	"alming_backend/src/repository"
)

func GetWords(c *engine.Context) {
	words := repository.GetWords()
	c.WriteJSON(words)
}

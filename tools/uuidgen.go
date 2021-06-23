package tools

import (
	"github.com/google/uuid"
)

func UuidGen() string {
	Uuid := uuid.New().String()
	return Uuid
}

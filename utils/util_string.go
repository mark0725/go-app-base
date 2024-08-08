package utils

import (
	"github.com/google/uuid"
)

func GetUUID() string {
	// 生成新的UUID
	newUUID := uuid.New()
	return newUUID.String()
}

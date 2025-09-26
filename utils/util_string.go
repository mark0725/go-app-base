package utils

import (
	"strings"

	"github.com/google/uuid"
)

func GetUUID() string {
	// 生成新的UUID
	newUUID := uuid.New()
	return newUUID.String()
}

func GetObjectID(objectType string) string {
	// 生成新的UUID
	newUUID := uuid.New()
	id := objectType + "-" + strings.ReplaceAll(newUUID.String(), "-", "")
	return id
}

func Contains(arr []string, str string) bool {
	for _, v := range arr {
		if v == str {
			return true
		}
	}
	return false
}

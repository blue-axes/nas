package utils

import "github.com/google/uuid"

func UUIDString() string {
	uid := uuid.New()
	return uid.String()
}

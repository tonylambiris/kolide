package model

import (
	"github.com/nu7hatch/gouuid"
)

// UniqueKey returns a new v4 uuid
func UniqueKey() string {
	uuid, _ := uuid.NewV4()
	return uuid.String()
}

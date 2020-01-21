package utils

import "time"

func LoadLocation() time.Location {
	location, _ := time.LoadLocation("Europe/Moscow")
	return *location
}

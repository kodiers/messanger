package migrations

import "time"

type Migration struct {
	Id      int
	Name    string
	Created time.Time
}

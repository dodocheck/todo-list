package models

import "time"

type ActionLog struct {
	Action string
	Time   time.Time
}

package models

import "time"

type Response map[string]interface{}

type Spammer struct {
	Addr      string
	Count     int
	TimeLimit time.Time
}

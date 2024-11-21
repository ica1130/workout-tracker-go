package data

import "time"

type Member struct {
	ID        int64
	Email     string
	Name      string
	Height    int64
	Weight    int64
	CreatedAt time.Time
}

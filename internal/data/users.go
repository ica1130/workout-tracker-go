package data

import "time"

type User struct {
	ID        int64
	Email     string
	Name      string
	Height    int64
	Weight    int64
	CreatedAt time.Time
}

package model

import "time"

type User struct {
	UUID      string    `json:"uuid"`
	Login     string    `json:"login"`
	Balance   float64   `json:"balance"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
}

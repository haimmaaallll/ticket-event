package domain

import "time"

// merepsentasikan tabel user
type User struct {
	Id       string
	Name     string
	Email    string
	Password string
	Role     string
	CreateAt time.Time
}

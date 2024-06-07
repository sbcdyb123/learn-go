package domain

import "time"

type User struct {
	Id       int64
	Email    string
	Phone    string
	Password string
	Username string
	BirthDay int64
	Intro    string
	CTime    time.Time
	UTime    time.Time
}

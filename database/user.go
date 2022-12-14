package database

import "time"

type User struct {
	Id        *int64
	Username  string
	Password  string
	CreatedAt time.Time
	IsAdmin   bool
}

func (u User) FormattedTime() string {
	return u.CreatedAt.Format(time.RFC1123)
}

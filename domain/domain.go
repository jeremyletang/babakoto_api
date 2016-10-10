package domain

import "time"

type User struct {
	Id        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type AccessToken struct {
	Id        string    `json:"id"`
	UserId    string    `json:"user_id"`
	Ttl       int       `json:"ttl"`
	CreatedAt time.Time `json:"created_at"`
}

type UserSignupVerification struct {
	Id        string    `json:"id"`
	UserId    string    `json:"user_id"`
	Ttl       int       `json:"ttl"`
	CreatedAt time.Time `json:"created_at"`
}

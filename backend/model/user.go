package model

import "time"

type User struct {
	ID           int64     `json:"id" db:"id"`
	Email        string    `json:"email" db:"email"`
	Username     string    `json:"username" db:"username"`
	Name         string    `json:"name" db:"name"`
	GoogleID     string    `json:"google_id" db:"google_id"`
	PictureUrl   string    `json:"picture_url" db:"picture_url"`
	RefreshToken string    `json:"refresh_token" db:"refresh_token"`
	Expiry       time.Time `json:"expiry" db:"expiry"`
	CreatedAt    time.Time `json:"createt_at" db:"createt_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

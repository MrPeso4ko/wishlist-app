package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Login        string `gorm:"uniqueIndex;not null"`
	PasswordHash string `gorm:"not null"`
	Wishes       []Wish
}

type PublicUser struct {
	ID    uint   `json:"id"`
	Login string `json:"login"`
}

func (u *User) ToPublic() *PublicUser {
	return &PublicUser{
		ID:    u.ID,
		Login: u.Login,
	}
}

package models

import (
	"gorm.io/gorm"
)

type Wish struct {
	gorm.Model
	UserID   uint   `gorm:"not null"`
	Title    string `gorm:"not null"`
	Comment  string `gorm:"size:500"`
	ImageURL string
	Price    float64
	User     User `gorm:"foreignKey:UserID"`
}

type PublicWish struct {
	ID       uint       `json:"id"`
	Title    string     `json:"title"`
	Comment  string     `json:"comment,omitempty"`
	ImageURL string     `json:"image_url,omitempty"`
	Price    float64    `json:"price,omitempty"`
	User     PublicUser `json:"user"`
}

func (w *Wish) ToPublic() *PublicWish {
	return &PublicWish{
		ID:       w.ID,
		Title:    w.Title,
		Comment:  w.Comment,
		ImageURL: w.ImageURL,
		Price:    w.Price,
		User:     *w.User.ToPublic(),
	}
}

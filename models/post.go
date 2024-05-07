package models

import (
	"time"

	"github.com/google/uuid"
)

type Post struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null" json:"userID"`
	User      User      `gorm:"" json:"user"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	CreatedAt time.Time `gorm:"default:current_timestamp" json:"postedAt"`
	NoLikes   int       `gorm:"default:0" json:"noLikes"`
	Edited    bool      `gorm:"default:false" json:"edited"`
	Likes     []Like    `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE" json:"-"`
}

type Like struct {
	ID        uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	UserID    uuid.UUID  `gorm:"type:uuid;not null" json:"likedByID"`
	PostID    *uuid.UUID `gorm:"type:uuid" json:"postID"`
	CreatedAt time.Time  `gorm:"default:current_timestamp" json:"likedAt"`
}

package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type User struct {
	ID          uuid.UUID        `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	Name        string           `gorm:"type:text;not null" json:"name"`
	Username    string           `gorm:"type:text;unique;not null" json:"username"`
	Email       string           `gorm:"unique;not null" json:"-"`
	Password    string           `json:"-"`
	Bio         string           `json:"bio"`
	Links       pq.StringArray   `gorm:"type:text[]" json:"links"`
	Followers   []FollowFollower `gorm:"foreignKey:FollowerID;constraint:OnDelete:CASCADE" json:"-"`
	Following   []FollowFollower `gorm:"foreignKey:FollowedID;constraint:OnDelete:CASCADE" json:"-"`
	NoFollowing int              `gorm:"default:0" json:"noFollowing"`
	NoFollowers int              `gorm:"default:0" json:"noFollowers"`
	CreatedAt   time.Time        `gorm:"default:current_timestamp;index:idx_created_at,sort:desc" json:"-"`
}

type FollowFollower struct { //* follower follows followed
	FollowerID uuid.UUID `json:"followerID"`
	Follower   User      `gorm:"foreignKey:FollowerID" json:"follower"`
	FollowedID uuid.UUID `json:"followedID"`
	Followed   User      `gorm:"foreignKey:FollowedID" json:"followed"`
	CreatedAt  time.Time `gorm:"default:current_timestamp" json:"-"`
}

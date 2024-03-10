package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID           uint           `gorm:"primarykey;autoIncrement"`
	Uuid         string         `gorm:"primarykey;"`
	DeletedAt    gorm.DeletedAt `gorm:"index"`
	Prefix       *string        `gorm:"default:null"`
	Firstname    *string        `gorm:"default:null"`
	Lastname     *string        `gorm:"default:null"`
	Fullname     *string        `gorm:"default:null"`
	Phone        *string        `gorm:"default:null"`
	Email        *string        `gorm:"default:null"`
	Username     *string        `gorm:"default:null"`
	Password     *string        `gorm:"default:null"`
	Role         *string        `gorm:"default:null"`
	RefreshToken *string        `gorm:"default:null"`

	// Custom detail
}

func (user *User) BeforeCreate(db *gorm.DB) (err error) {
	uuid := uuid.New()
	user.Uuid = uuid.String()
	return nil
}

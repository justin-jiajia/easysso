package database

import (
	"time"
)

type User struct {
	ID             uint `gorm:"primaryKey"`
	UserName       string
	PasswordHash   string
	AvatarFileName string
	CreatedAt      time.Time
}

type Token struct {
	Token    string
	Exp      time.Time
	Code     string `gorm:"primaryKey"`
	ClientID string
}

func Migrate() {
	DB.AutoMigrate(&User{})
	DB.AutoMigrate(&Token{})
}

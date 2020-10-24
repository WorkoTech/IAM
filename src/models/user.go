package models

import (
	"strconv"

	"github.com/segmentio/fasthash/fnv1a"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type User struct {
    gorm.Model
    Email       string  `gorm:"unique;not null"`
    Password    string  `gorm:"not null"`
}

func (user *User) CheckPassword(plainPassword string) bool {
	encryptedPassword := strconv.FormatUint(fnv1a.HashString64(plainPassword), 32)
	return encryptedPassword == user.Password
}

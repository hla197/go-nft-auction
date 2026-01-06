package models

import "gorm.io/gorm"

type User struct {
	Id       uint64 `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	Address  string `json:"address" gorm:"uniqueIndex:uk_email;type:varchar(255)" `
	Password string `json:"password" gorm:"type:varchar(255)" `
	Username string `json:"username" gorm:"uniqueIndex:uk_username;type:varchar(255)" `
	*gorm.Model
}

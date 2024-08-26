package model

import "gorm.io/gorm"

type Message struct {
	gorm.Model
	Username string `json:"username"`
	Message  string `json:"message"`
}

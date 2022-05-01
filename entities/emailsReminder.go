package entities

import (
	"gorm.io/gorm"
)

type EmailsReminder struct {
	gorm.Model
	TaskId int
	Email  string
}

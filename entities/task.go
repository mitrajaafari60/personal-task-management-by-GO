package entities

import (
	"time"

	"gorm.io/gorm"
)

type Task struct {
	gorm.Model
	Description  string    `gorm:"unique;not null;type:varchar(100);default:null"`
	UserID       int       `json:"user_id",gorm:"not null"` // user ID which is assigned
	DeadLine     time.Time `json:"dead_line"`
	Reminder     uint      `gorm:"not null;type:int;default:24"` // time in hour
	RemindCounts int       // Update the counter after each remind
}

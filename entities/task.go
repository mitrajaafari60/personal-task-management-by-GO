package entities

import (
	"time"

	"gorm.io/gorm"
)

type Task struct {
	gorm.Model
	Description  string    `gorm:"not null;type:varchar(100);default:null"`
	UserID       int       `json:"user_id",gorm:"not null"` // user ID which is assigned
	StartTime    time.Time `json:"start_time",gorm:"not null"`
	EndTime      time.Time `json:"end_time",gorm:"not null"`
	Reminder     uint      `gorm:"not null;type:int;default:24"` // time in hour
	RemindCounts int       // Update the counter after each remind
}

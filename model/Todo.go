package model

import "time"

type Todo struct {
	ID        uint32    `gorm:"primary_key;auto_increment" json:"id"`
	UserID    uint32    `gorm:"not null;" json:"user_id"`
	Task      string    `gorm:"not null;" json:"task"`
	DueDate   time.Time `gorm:"not null;" json:"duedate"`
	Completed bool      `gorm:"not null;" json:"completed"`
}

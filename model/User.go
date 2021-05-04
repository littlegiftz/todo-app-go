package model

type User struct {
	ID       uint32 `gorm:"primary_key;auto_increment" json:"id"`
	Email    string `gorm:"not null;unique" json:"email"`
	Password string `gorm:"not null;" json:"password"`
}

package models

type GameSave struct {
	ID       uint   `gorm:"primaryKey"`
	GameName string `gorm:"unique;not null"`
	DirPath  string `gorm:"not null"`
}

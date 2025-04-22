package models

type Settings struct {
	ID     uint `gorm:"primaryKey"`
	APIKey string
	URL    string
}

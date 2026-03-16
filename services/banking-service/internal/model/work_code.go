package model

// WorkCode represents the industry classification (e.g. "10.1" = Food production)
type WorkCode struct {
	WorkCodeID  uint   `gorm:"primaryKey"`
	Code        string `gorm:"uniqueIndex;not null;size:10"` // "10.1"
	Description string `gorm:"not null"`
}

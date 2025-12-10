package domain

import "time"

type Event struct {
	BaseModel
	Name           string    `gorm:"type:varchar(100);not null" json:"name"`
	Location       string    `gorm:"type:text;not null" json:"location"`
	Date           time.Time `gorm:"not null" json:"date"`
	Price          float64   `gorm:"type:decimal(10,2);not null" json:"price"`
	Description    string    `gorm:"type:text" json:"description,omitempty"`
	Image          string    `gorm:"type:text" json:"image,omitempty"`
	TotalStock     int       `gorm:"not null" json:"total_stock"`
	AvailableStock int       `gorm:"not null;check:available_stock <= total_stock" json:"available_stock"`
}

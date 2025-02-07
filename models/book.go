package models

type Book struct {
	ID     uint    `gorm:"primaryKey" json:"id"`
	Title  string  `gorm:"uniqueIndex" json:"title"`
	Author string  `json:"author"`
	Price  float64 `json:"price"`
}

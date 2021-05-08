package model

import "gorm.io/gorm"

type Model interface {
	Load(g Getter)
	Validate() error
	Initialize(db *gorm.DB)
	Find(db *gorm.DB) ([]Model, error)
}

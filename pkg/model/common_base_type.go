package model

import (
	"strconv"

	"github.com/go-playground/validator"
	"gorm.io/gorm"
)

// CommonBaseType schema of the commonBaseType table
type CommonBaseType struct {
	ID            int64  `json:"common_base_type_id" gorm:"primary_key;auto_increment;not null"`
	BaseTypeTitle string `json:"base_type_title" validate:"required" gorm:"not null;type:varchar(800);UNIQUE"`
	BaseTypeCode  string `json:"base_type_code" gorm:"type:varchar(3);UNIQUE;<-:create"`
}

func (cbt *CommonBaseType) Load(g Getter) {
	cbt.ID, _ = strconv.ParseInt(g.Get("id"), 10, 64)
	cbt.BaseTypeTitle = g.Get("base_type_title")
	cbt.BaseTypeCode = g.Get("base_type_code")
}

func (cbt *CommonBaseType) Validate() error {
	return validator.New().Struct(cbt)
}

func (cbt *CommonBaseType) Initialize(db *gorm.DB) {
	cbt.BaseTypeCode = generateCode(3)
}

func (cbt *CommonBaseType) Find(db *gorm.DB) ([]Model, error) {
	result := []CommonBaseType{}
	if err := db.Find(&result, cbt).Error; err != nil {
		return nil, err
	}
	ret := make([]Model, len(result))
	for i := range result {
		ret[i] = &result[i]
	}
	return ret, nil
}

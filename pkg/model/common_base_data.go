package model

import (
	"strconv"

	"github.com/go-playground/validator"
	"gorm.io/gorm"
)

// CommonBaseData schema of the commonBaseData table
type CommonBaseData struct {
	ID               int64          `json:"common_base_data_id" gorm:"primary_key;auto_increment;not null"`
	BaseCode         string         `json:"base_code" validate:"required" gorm:"UNIQUE_INDEX:compositeindex;index;not null;type:varchar(6);<-:create"`
	BaseValue        string         `json:"base_value" validate:"required" gorm:"not null;type:varchar(800)"`
	CommonBaseTypeID int64          `json:"common_base_type_id" validate:"required" gorm:"UNIQUE_INDEX:compositeindex"`
	CommonBaseType   CommonBaseType `json:"common_base_type" validate:"-" gorm:"foreignKey:CommonBaseTypeID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (cbd *CommonBaseData) Load(g Getter) {
	cbd.ID, _ = strconv.ParseInt(g.Get("id"), 10, 64)
	cbd.CommonBaseTypeID, _ = strconv.ParseInt(g.Get("common_base_type_id"), 10, 64)
	cbd.BaseCode = g.Get("base_code")
	cbd.BaseValue = g.Get("base_value")
}

func (cbd *CommonBaseData) Validate() error {
	return validator.New().Struct(cbd)
}

func (cbd *CommonBaseData) Initialize(db *gorm.DB) {
	if cbd.CommonBaseType.BaseTypeCode == "" {
		db.Find(&cbd.CommonBaseType, &CommonBaseType{ID: cbd.CommonBaseTypeID})
	}
	cbd.BaseCode = cbd.CommonBaseType.BaseTypeCode + generateCode(3)
}

func (cbd *CommonBaseData) Find(db *gorm.DB) ([]Model, error) {
	result := []CommonBaseData{}
	if err := db.Preload("CommonBaseType").Find(&result, cbd).Error; err != nil {
		return nil, err
	}
	ret := make([]Model, len(result))
	for i := range result {
		ret[i] = &result[i]
	}
	return ret, nil
}

func (cbd *CommonBaseData) BeforeUpdate(tx *gorm.DB) (err error) {
	cbd.CommonBaseType = CommonBaseType{}
	return nil
}

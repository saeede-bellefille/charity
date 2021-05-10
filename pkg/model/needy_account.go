package model

import (
	"fmt"
	"strconv"

	"github.com/go-playground/validator"
	"gorm.io/gorm"
)

// NeedyAccount schema of the needyAccounts table
type NeedyAccount struct {
	ID            int64          `json:"needy_account_id" gorm:"primary_key;auto_increment;not null"`
	BankID        int64          `json:"bank_id" validate:"required" gorm:"not null"`
	Bank          CommonBaseData `json:"bank" validate:"required" gorm:"foreignKey:BankID;not null;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	NeedyID       int64          `json:"needy_id" gorm:"not null;UNIQUE_INDEX:compositeindex;index"`
	Needy         Personal       `json:"needy" validate:"required" gorm:"foreignKey:NeedyID;not null;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	OwnerName     string         `json:"owner_name" validate:"required" gorm:"not null;type:varchar(1000)"`
	CardNumber    string         `json:"card_number" gorm:"type:varchar(20)"`
	AccountNumber string         `json:"account_number" validate:"required" gorm:"not null;type:varchar(10);UNIQUE_INDEX:compositeindex;"`
	AccountName   string         `json:"account_name" gorm:"type:varchar(500)"`
	ShebaNumber   string         `json:"sheba_number" validate:"required" gorm:"not null;type:varchar(26);UNIQUE"`
}

func (na *NeedyAccount) Load(g Getter) {
	na.ID, _ = strconv.ParseInt(g.Get("id"), 10, 64)
	na.BankID, _ = strconv.ParseInt(g.Get("bank_id"), 10, 64)
}

func (na *NeedyAccount) Validate() error {
	if !validateSheba(na.ShebaNumber) {
		return fmt.Errorf("This sheba number is not valid: %s", na.ShebaNumber)
	}
	return validator.New().Struct(na)
}

func (na *NeedyAccount) Initialize(db *gorm.DB) {
}

func (na *NeedyAccount) Find(db *gorm.DB) ([]Model, error) {
	result := []NeedyAccount{}
	if err := db.Preload("Needy").Preload("Bank").Preload("Bank.CommonBaseType").Find(&result, na).Error; err != nil {
		return nil, err
	}
	ret := make([]Model, len(result))
	for i := range result {
		ret[i] = &result[i]
	}
	return ret, nil
}

package model

import (
	"fmt"
	"strconv"

	"github.com/go-playground/validator"
	"gorm.io/gorm"
)

// CharityAccount schema of the charityAccount table
type CharityAccount struct {
	ID            int64          `json:"charity_account_id" gorm:"primary_key;auto_increment;not null"`
	BankID        int64          `json:"bank_id" validate:"required" gorm:"not null"`
	Bank          CommonBaseData `json:"bank" validate:"-" gorm:"foreignKey:BankID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	BranchName    string         `json:"branch_name" validate:"required" gorm:"not null;type:varchar(500)"`
	OwnerName     string         `json:"owner_name" validate:"required" gorm:"not null;type:varchar(1000)"`
	CardNumber    string         `json:"card_number" gorm:"type:varchar(20)"`
	AccountNumber string         `json:"account_number" validate:"required" gorm:"not null;type:varchar(10);UNIQUE"`
	AccountName   string         `json:"account_name" gorm:"type:varchar(500)"`
}

func (ca *CharityAccount) Load(g Getter) {
	ca.ID, _ = strconv.ParseInt(g.Get("id"), 10, 64)
	ca.BankID, _ = strconv.ParseInt(g.Get("bank_id"), 10, 64)
	ca.BranchName = g.Get("branch_name")
	ca.OwnerName = g.Get("owner_name")
	ca.CardNumber = g.Get("card_number")
	ca.AccountNumber = g.Get("account_number")
	ca.AccountName = g.Get("account_name")
}

func (ca *CharityAccount) Validate() error {
	if !validateCardNumber(ca.CardNumber) {
		return fmt.Errorf("This number is not valid: %s", ca.CardNumber)
	}
	return validator.New().Struct(ca)
}

func (ca *CharityAccount) Initialize(db *gorm.DB) {
}

func (ca *CharityAccount) Find(db *gorm.DB) ([]Model, error) {
	result := []CharityAccount{}
	if err := db.Preload("Bank").Preload("Bank.CommonBaseType").Find(&result, ca).Error; err != nil {
		return nil, err
	}
	ret := make([]Model, len(result))
	for i := range result {
		ret[i] = &result[i]
	}
	return ret, nil
}

func (ca *CharityAccount) BeforeUpdate(tx *gorm.DB) (err error) {
	ca.Bank = CommonBaseData{}
	return nil
}

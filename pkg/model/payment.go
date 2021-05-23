package model

import (
	"fmt"
	"strconv"

	"github.com/go-playground/validator"
	"gorm.io/gorm"
)

// Payment schema of the payment table
type Payment struct {
	ID                     int64                `json:"payment_id" gorm:"primary_key;auto_increment;not null"`
	DonatorID              int64                `json:"donator_id"`
	Donator                Personal             `json:"donator" validate:"-" gorm:"foreignKey:DonatorID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CashAssistanceDetailID int64                `json:"cash_assistance_detail_id" validate:"required"`
	CashAssistanceDetail   CashAssistanceDetail `json:"cash_assistance_detail" validate:"-" gorm:"foreignKey:CashAssistanceDetailID;not null;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT"`
	PaymentPrice           float64              `json:"price" validate:"required" gorm:"not null;type:decimal(19,3)"`
	PaymentGatewayID       string               `json:"payment_gateway_id" gorm:"type:varchar(10)"`
	PaymentDate            string               `json:"payment_date" validate:"required" gorm:"not null;type:varchar(10)"`
	PaymentTime            string               `json:"payment_time" validate:"required" gorm:"not null;type:varchar(10)"`
	PaymentStatus          string               `json:"payment_status" validate:"required" gorm:"not null;type:varchar(500)"`
	SourceAccoutNumber     string               `json:"source_account_number" gorm:"type:varchar(10)"`
	TargetAccountNumber    string               `json:"target_account_number" validate:"required" gorm:"not null;type:varchar(10)"`
	CharityAccountID       *int64               `json:"charity_account_id" gorm:"DEFAULT:null"`
	CharityAccount         CharityAccount       `json:"charity_account" validate:"-" gorm:"foreignKey:CharityAccountID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	FollowCode             string               `json:"follow_code" validate:"required" gorm:"not null;type:varchar(10)"`
	NeedyID                int64                `json:"needy_id"`
	Needy                  Personal             `json:"needy" validate:"-" gorm:"foreignKey:NeedyID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	paymentSum             float64              `json:"-" gorm:"-"`
	settlementSum          float64              `json:"-" gorm:"-"`
}

func (p *Payment) Load(g Getter) {
	p.ID, _ = strconv.ParseInt(g.Get("id"), 10, 64)
	p.DonatorID, _ = strconv.ParseInt(g.Get("owner_donator_id"), 10, 64)
	p.CashAssistanceDetailID, _ = strconv.ParseInt(g.Get("cash_assistance_detail_id"), 10, 64)
	p.PaymentPrice, _ = strconv.ParseFloat(g.Get("price"), 64)
	p.PaymentGatewayID = g.Get("payment_gateway_id")
	p.PaymentTime = g.Get("payment_time")
	p.PaymentStatus = g.Get("payment_status")
	p.SourceAccoutNumber = g.Get("source_number")
	p.TargetAccountNumber = g.Get("target_number")
}

func (p *Payment) Validate() error {
	if p.CashAssistanceDetail.ID != 0 && p.settlementSum+p.PaymentPrice > p.CashAssistanceDetail.NeededPrice {
		return fmt.Errorf("Sum of settlement payments + price is more than needed")
	}
	if p.CashAssistanceDetail.ID != 0 && p.paymentSum+p.PaymentPrice > p.CashAssistanceDetail.NeededPrice {
		return fmt.Errorf("Sum of payments + price is more than needed")
	}
	return validator.New().Struct(p)
}

func (p *Payment) Initialize(db *gorm.DB) {
	if p.CashAssistanceDetail.ID == 0 && p.CashAssistanceDetailID != 0 {
		db.Find(&p.CashAssistanceDetail, &CashAssistanceDetail{ID: p.CashAssistanceDetailID})
	}
	if p.CharityAccount.ID == 0 && p.CharityAccountID != nil {
		db.Find(&p.CharityAccount, &CharityAccount{ID: *p.CharityAccountID})
	}
	db.Table("payments").
		Where(&Payment{CashAssistanceDetailID: p.CashAssistanceDetailID, PaymentStatus: "Success"}).
		Select("SUM(payment_price)").
		Row().
		Scan(&p.paymentSum)
	db.Table("payments").
		Where(&Payment{CashAssistanceDetailID: p.CashAssistanceDetailID, PaymentStatus: "Success", CharityAccountID: p.CharityAccountID}).
		Select("SUM(payment_price)").
		Row().
		Scan(&p.settlementSum)
}

func (p *Payment) Find(db *gorm.DB) ([]Model, error) {
	result := []Payment{}
	if err := db.Preload("Needy").
		Preload("Donator").
		Preload("CharityAccount").
		Preload("CharityAccount.Bank").
		Preload("CharityAccount.Bank.CommonBaseType").
		Preload("CashAssistanceDetail").
		Preload("CashAssistanceDetail.Plan").
		Preload("CashAssistanceDetail.AssignNeedyPlan").
		Preload("CashAssistanceDetail.AssignNeedyPlan.Needy").
		Preload("CashAssistanceDetail.AssignNeedyPlan.Plan").
		Find(&result, p).Error; err != nil {
		return nil, err
	}
	ret := make([]Model, len(result))
	for i := range result {
		ret[i] = &result[i]
	}
	return ret, nil
}

func (p *Payment) BeforeUpdate(tx *gorm.DB) (err error) {
	p.Donator = Personal{}
	p.CashAssistanceDetail = CashAssistanceDetail{}
	p.CharityAccount = CharityAccount{}
	p.Needy = Personal{}
	return nil
}

package model

import (
	"fmt"
	"strconv"
	"time"

	"github.com/go-playground/validator"
	"gorm.io/gorm"
)

// Payment schema of the payment table
type Payment struct {
	ID                     int64                `json:"payment_id" gorm:"primary_key;auto_increment;not null"`
	DonatorID              int64                `json:"donator_id"`
	Donator                Personal             `json:"donator" gorm:"foreignKey:DonatorID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CashAssistanceDetailID int64                `json:"cash_assistance_detail_id" validate:"required"`
	CashAssistanceDetail   CashAssistanceDetail `json:"cash_assistance_detail" gorm:"foreignKey:CashAssistanceDetailID;not null;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT"`
	PaymentPrice           float64              `json:"price" validate:"required" gorm:"not null;type:money"`
	PaymentGatewayID       string               `json:"payment_gateway_id" gorm:"type:varchar(10)"`
	PaymentDate            time.Time            `json:"payment_data" validate:"required" gorm:"not null;type:date"`
	PaymentTime            time.Time            `json:"payment_time" validate:"required" gorm:"not null;type:time"`
	PaymentStatus          string               `json:"payment_status" validate:"required" gorm:"not null;type:varchar(500)"`
	SourceAccoutNumber     string               `json:"source_account_number" gorm:"type:varchar(10)"`
	TargetAccountNumber    string               `json:"target_account_number" validate:"required" gorm:"not null;type:varchar(10)"`
	CharityAccountID       int64                `json:"charity_account_id" gorm:"not null"`
	CharityAccount         CharityAccount       `json:"charity_account" validate:"required" gorm:"foreignKey:CharityAccountID;not null;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	FollowCode             string               `json:"follow_code" validate:"required" gorm:"not null;type:varchar(10)"`
	NeedyID                int64                `json:"needy_id"`
	Needy                  Personal             `json:"needy" gorm:"foreignKey:NeedyID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	paymentSum             float64              `json:"-" gorm:"-"`
	settlementSum          float64              `json:"-" gorm:"-"`
}

func (p *Payment) Load(g Getter) {
	p.ID, _ = strconv.ParseInt(g.Get("id"), 10, 64)
	p.DonatorID, _ = strconv.ParseInt(g.Get("owner_donator_id"), 10, 64)
	p.CashAssistanceDetailID, _ = strconv.ParseInt(g.Get("cash_assistance_detail_id"), 10, 64)
	p.PaymentPrice, _ = strconv.ParseFloat(g.Get("price"), 64)
	p.PaymentGatewayID = g.Get("payment_gateway_id")
	p.PaymentTime, _ = time.Parse(time.RFC3339, g.Get("payment_time"))
	p.PaymentStatus = g.Get("payment_status")
	p.SourceAccoutNumber = g.Get("source_number")
	p.TargetAccountNumber = g.Get("target_number")
}

func (p *Payment) Validate() error {
	if p.settlementSum+p.PaymentPrice > p.CashAssistanceDetail.NeededPrice {
		return fmt.Errorf("Sum of payments + price is more than needed")
	}
	if p.paymentSum+p.PaymentPrice > p.CashAssistanceDetail.NeededPrice {
		return fmt.Errorf("Sum of payments + price is more than needed")
	}
	return validator.New().Struct(p)
}

func (p *Payment) Initialize(db *gorm.DB) {
	if p.CashAssistanceDetail.ID == 0 && p.CashAssistanceDetailID != 0 {
		db.Find(&p.CashAssistanceDetail, &CashAssistanceDetail{ID: p.CashAssistanceDetailID})
	}
	if p.CharityAccount.ID == 0 && p.CharityAccountID != 0 {
		db.Find(&p.CharityAccount, &CharityAccount{ID: p.CharityAccountID})
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
	if err := db.Preload("Personal").
		Preload("CharityAccount").
		Preload("CharityAccount.CommonBaseData").
		Preload("CharityAccount.CommonBaseData.CommonBaseType").
		Preload("CashAssistanceDetail").
		Preload("CashAssistanceDetail.Plan").
		Preload("CashAssistanceDetail.AssignNeedyToPlan").
		Preload("CashAssistanceDetail.AssignNeedyToPlan.Personal").
		Preload("CashAssistanceDetail.AssignNeedyToPlan.Plan").
		Find(&result, p).Error; err != nil {
		return nil, err
	}
	ret := make([]Model, len(result))
	for i := range result {
		ret[i] = &result[i]
	}
	return ret, nil
}

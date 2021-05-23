package model

import (
	"fmt"
	"strconv"

	"github.com/go-playground/validator"
	"gorm.io/gorm"
)

// CashAssistanceDetail schema of the cashAssistanceDetail table
type CashAssistanceDetail struct {
	ID                int64             `json:"cash_assistance_detail_id" gorm:"primary_key;auto_increment;not null"`
	AssignNeedyPlanID int64             `json:"assign_needy_plan_id" validate:"required" gorm:"UNIQUE_INDEX:compositeindex;index"`
	AssignNeedyPlan   AssignNeedyToPlan `json:"assign_needy_plan" validate:"-" gorm:"foreignKey:AssignNeedyPlanID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	PlanID            int64             `json:"plan_id" validate:"required" gorm:"UNIQUE_INDEX:compositeindex;not null"`
	Plan              Plan              `json:"plan" validate:"-" gorm:"foreignKey:PlanID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT"`
	NeededPrice       float64           `json:"needed_price" validate:"required" gorm:"not null;type:decimal(19,3)"`
	MinPrice          float64           `json:"min_price" gorm:"type:decimal(19,3)"`
	Description       string            `json:"description" gorm:"type:text"`
}

func (cad *CashAssistanceDetail) Load(g Getter) {
	cad.ID, _ = strconv.ParseInt(g.Get("id"), 10, 64)
	cad.AssignNeedyPlanID, _ = strconv.ParseInt(g.Get("assign_id"), 10, 64)
	cad.PlanID, _ = strconv.ParseInt(g.Get("plan_id"), 10, 64)
	cad.NeededPrice, _ = strconv.ParseFloat(g.Get("price"), 64)
	cad.MinPrice, _ = strconv.ParseFloat(g.Get("min_price"), 64)
	cad.Description = g.Get("description")
}

func (cad *CashAssistanceDetail) Validate() error {
	if cad.MinPrice > cad.NeededPrice {
		return fmt.Errorf("min price must be less than needed price!!")
	}
	return validator.New().Struct(cad)
}

func (cad *CashAssistanceDetail) Initialize(db *gorm.DB) {
}

func (cad *CashAssistanceDetail) Find(db *gorm.DB) ([]Model, error) {
	result := []CashAssistanceDetail{}
	if err := db.Preload("Plan").Preload("AssignNeedyPlan").Preload("AssignNeedyPlan.Needy").Preload("AssignNeedyPlan.Plan").Find(&result, cad).Error; err != nil {
		return nil, err
	}
	ret := make([]Model, len(result))
	for i := range result {
		ret[i] = &result[i]
	}
	return ret, nil
}

func (cad *CashAssistanceDetail) BeforeUpdate(tx *gorm.DB) (err error) {
	cad.AssignNeedyPlan = AssignNeedyToPlan{}
	cad.Plan = Plan{}
	return nil
}

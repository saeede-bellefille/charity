package model

import (
	"fmt"
	"strconv"
	"time"

	"github.com/go-playground/validator"
	"gorm.io/gorm"
)

// AssignNeedyToPlans schema of the assignNeedyToPlans table
type AssignNeedyToPlan struct {
	ID      int64     `json:"assign_needy_plan_id" gorm:"primary_key;auto_increment;not null"`
	NeedyID int64     `json:"needy_id" gorm:"not null;UNIQUE_INDEX:compositeindex;index"`
	Needy   Personal  `json:"needy" validate:"required" gorm:"foreignKey:NeedyID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	PlanID  int64     `json:"plan_id" gorm:"not null;UNIQUE_INDEX:compositeindex"`
	Plan    Plan      `json:"plan" validate:"required" gorm:"foreignKey:PlanID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Fdate   time.Time `json:"fdate" validate:"required" gorm:"not null;type:timestamp"`
	Tdate   time.Time `json:"tdate" validate:"required" gorm:"not null;type:timestamp"`
}

func (antp *AssignNeedyToPlan) Load(g Getter) {
	antp.ID, _ = strconv.ParseInt(g.Get("id"), 10, 64)
	antp.PlanID, _ = strconv.ParseInt(g.Get("plan_id"), 10, 64)
	antp.NeedyID, _ = strconv.ParseInt(g.Get("needy_id"), 10, 64)
	antp.Fdate, _ = time.Parse(time.RFC3339, g.Get("start_date"))
	antp.Tdate, _ = time.Parse(time.RFC3339, g.Get("end_date"))
}

func (antp *AssignNeedyToPlan) Validate() error {
	if antp.Fdate.After(antp.Tdate) {
		return fmt.Errorf("start_date must be less than end_date!!")
	}
	if antp.Fdate != antp.Plan.Fdate {
		return fmt.Errorf("fdate is not same as plan fdate!")
	}
	if antp.Tdate != antp.Plan.Tdate {
		return fmt.Errorf("tdate is not same as plan tdate!")
	}
	return validator.New().Struct(antp)
}

func (antp *AssignNeedyToPlan) Initialize(db *gorm.DB) {
}

func (antp *AssignNeedyToPlan) Find(db *gorm.DB) ([]Model, error) {
	result := []AssignNeedyToPlan{}
	if err := db.Preload("Personal").Preload("Plan").Find(&result, antp).Error; err != nil {
		return nil, err
	}
	ret := make([]Model, len(result))
	for i := range result {
		ret[i] = &result[i]
	}
	return ret, nil
}

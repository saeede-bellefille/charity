package model

import (
	"fmt"
	"strconv"
	"time"

	"github.com/go-playground/validator"
	"gorm.io/gorm"
)

// Plans schema of the plans table
type Plan struct {
	ID           int64     `json:"plan_id" gorm:"primary_key;auto_increment;not null"`
	PlanName     string    `json:"plan_name" validate:"required" gorm:"not null;type:varchar(1000);UNIQUE_INDEX:compositeindex;index"`
	Description  string    `json:"description" gorm:"type:text"`
	PlanNature   bool      `json:"plan_nature" validate:"required" gorm:"not null;type:boolean;default:true;UNIQUE_INDEX:compositeindex"`
	ParentPlanID int64     `json:"parent_plan_id" gorm:"UNIQUE_INDEX:compositeindex"`
	ParentPlan   *Plan     `json:"parent_plan" validate:"required" gorm:"foreignKey:ParentPlanID"`
	Icon         []byte    `json:"icon" gorm:"type:bytea"`
	Fdate        time.Time `json:"fdate" gorm:"type:timestamp"`
	Tdate        time.Time `json:"tdate" gorm:"type:timestamp"`
	NeededLogin  bool      `json:"needed_login" validate:"required" gorm:"not null;type:boolean;default:false"`
}

func (p *Plan) Load(g Getter) {
	p.ID, _ = strconv.ParseInt(g.Get("id"), 10, 64)
	p.ParentPlanID, _ = strconv.ParseInt(g.Get("parent_plan_id"), 10, 64)
	p.PlanName = g.Get("plan_name")
	p.Fdate, _ = time.Parse(time.RFC3339, g.Get("start_date"))
	p.Tdate, _ = time.Parse(time.RFC3339, g.Get("end_date"))
	p.PlanNature, _ = strconv.ParseBool(g.Get("plan_nature"))
	p.NeededLogin, _ = strconv.ParseBool(g.Get("needed_login"))
}

func (p *Plan) Validate() error {
	if p.Fdate.After(p.Tdate) {
		return fmt.Errorf("start_date must be less than end_date!!")
	}
	return validator.New().Struct(p)
}

func (p *Plan) Initialize(db *gorm.DB) {
}

func (p *Plan) Find(db *gorm.DB) ([]Model, error) {
	result := []Plan{}
	if err := db.Preload("Plan").Find(&result, p).Error; err != nil {
		return nil, err
	}
	ret := make([]Model, len(result))
	for i := range result {
		ret[i] = &result[i]
	}
	return ret, nil
}

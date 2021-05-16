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
	ID           int64  `json:"plan_id" gorm:"primary_key;auto_increment;not null"`
	PlanName     string `json:"plan_name" validate:"required" gorm:"not null;type:varchar(1000);UNIQUE_INDEX:compositeindex;index"`
	Description  string `json:"description" gorm:"type:text"`
	PlanNature   bool   `json:"plan_nature" validate:"-" gorm:"not null;type:boolean;default:true;UNIQUE_INDEX:compositeindex"`
	ParentPlanID int64  `json:"parent_plan_id" validate:"required" gorm:"UNIQUE_INDEX:compositeindex"`
	ParentPlan   *Plan  `json:"parent_plan" validate:"-" gorm:"foreignKey:ParentPlanID"`
	Icon         []byte `json:"icon" gorm:"type:bytea"`
	Fdate        string `json:"fdate" gorm:"type:varchar(10)"`
	Tdate        string `json:"tdate" gorm:"type:varchar(10)"`
	NeededLogin  bool   `json:"needed_login" validate:"-" gorm:"not null;type:boolean;default:false"`
}

func (p *Plan) Load(g Getter) {
	p.ID, _ = strconv.ParseInt(g.Get("id"), 10, 64)
	p.ParentPlanID, _ = strconv.ParseInt(g.Get("parent_plan_id"), 10, 64)
	p.PlanName = g.Get("plan_name")
	p.Fdate = g.Get("start_date")
	p.Tdate = g.Get("end_date")
	p.PlanNature, _ = strconv.ParseBool(g.Get("plan_nature"))
	p.NeededLogin, _ = strconv.ParseBool(g.Get("needed_login"))
}

func (p *Plan) Validate() error {
	fdate, _ := time.Parse(time.RFC3339, p.Fdate)
	tdate, _ := time.Parse(time.RFC3339, p.Tdate)
	if fdate.After(tdate) {
		return fmt.Errorf("start_date must be less than end_date!!")
	}
	return validator.New().Struct(p)
}

func (p *Plan) Initialize(db *gorm.DB) {
}

func (p *Plan) Find(db *gorm.DB) ([]Model, error) {
	result := []Plan{}
	if err := db.Preload("ParentPlan").Find(&result, p).Error; err != nil {
		return nil, err
	}
	ret := make([]Model, len(result))
	for i := range result {
		ret[i] = &result[i]
	}
	return ret, nil
}

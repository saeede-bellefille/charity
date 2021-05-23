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
	ID      int64    `json:"assign_needy_plan_id" gorm:"primary_key;auto_increment;not null"`
	NeedyID int64    `json:"needy_id" validate:"required" gorm:"not null;UNIQUE_INDEX:compositeindex;index"`
	Needy   Personal `json:"needy" validate:"-" gorm:"foreignKey:NeedyID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	PlanID  int64    `json:"plan_id" validate:"required" gorm:"not null;UNIQUE_INDEX:compositeindex;index"`
	Plan    Plan     `json:"plan" validate:"-" gorm:"foreignKey:PlanID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Fdate   string   `json:"fdate" validate:"-" gorm:"not null;type:varchar(10)"`
	Tdate   string   `json:"tdate" validate:"-" gorm:"not null;type:varchar(10)"`
}

type AssignNeedyToPlanMultiple struct {
	NeedyID []int64 `json:"needy_id"`
	PlanID  int64   `json:"plan_id"`
	Fdate   string  `json:"fdate"`
	Tdate   string  `json:"tdate"`
}

func (antp *AssignNeedyToPlan) Load(g Getter) {
	antp.ID, _ = strconv.ParseInt(g.Get("id"), 10, 64)
	antp.PlanID, _ = strconv.ParseInt(g.Get("plan_id"), 10, 64)
	antp.NeedyID, _ = strconv.ParseInt(g.Get("needy_id"), 10, 64)
	antp.Fdate = g.Get("start_date")
	antp.Tdate = g.Get("end_date")
}

func (antp *AssignNeedyToPlan) Validate() error {
	fdate, _ := time.Parse(time.RFC3339, antp.Fdate)
	tdate, _ := time.Parse(time.RFC3339, antp.Tdate)
	if fdate.After(tdate) {
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
	if antp.Plan.Fdate == "" || antp.Plan.Tdate == "" {
		db.Find(&antp.Plan, &Plan{ID: antp.PlanID})
	}
}

func (antp *AssignNeedyToPlan) Find(db *gorm.DB) ([]Model, error) {
	result := []AssignNeedyToPlan{}
	if err := db.Preload("Needy").Preload("Plan").Find(&result, antp).Error; err != nil {
		return nil, err
	}
	ret := make([]Model, len(result))
	for i := range result {
		ret[i] = &result[i]
	}
	return ret, nil
}

func (antpm *AssignNeedyToPlanMultiple) ToList() []*AssignNeedyToPlan {
	ret := make([]*AssignNeedyToPlan, 0, len(antpm.NeedyID))
	for _, needy := range antpm.NeedyID {
		ret = append(ret, &AssignNeedyToPlan{
			NeedyID: needy,
			PlanID:  antpm.PlanID,
			Fdate:   antpm.Fdate,
			Tdate:   antpm.Tdate,
		})
	}
	return ret
}

func (antp *AssignNeedyToPlan) BeforeUpdate(tx *gorm.DB) (err error) {
	antp.Needy = Personal{}
	antp.Plan = Plan{}
	return nil
}

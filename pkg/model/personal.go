package model

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"strconv"

	"github.com/go-playground/validator"
	"gorm.io/gorm"
)

//PersonType enum PersonType
type PersonType int64

const (
	personel PersonType = iota + 1
	needy
	charity
)

// Personal schema of the personal table
type Personal struct {
	ID           int64      `json:"person_id" gorm:"primary_key;auto_increment;not null"`
	Name         string     `json:"name" validate:"-" gorm:"UNIQUE_INDEX:compositeindex;not null;type:varchar(500)"`
	Family       string     `json:"family" validate:"-" gorm:"UNIQUE_INDEX:compositeindex;not null;type:varchar(500)"`
	NationalCode string     `json:"national_code" gorm:"UNIQUE_INDEX:compositeindex;type:varchar(10)"`
	IDNumber     string     `json:"id_number" gorm:"type:varchar(10)"`
	Sex          bool       `json:"sex" validate:"-" gorm:"not null;type:boolean;default:false"`
	BirthDate    string     `json:"birth_date" gorm:"type:varchar(10)"`
	BirthPlace   string     `json:"birth_place" gorm:"type:varchar(500)"`
	PersonType   PersonType `json:"person_type" validate:"-" gorm:"not null"`
	PersonPhoto  []byte     `json:"person_photo" gorm:"type:bytea"`
	SecretCode   string     `json:"secret_code" gorm:"type:varchar(20)"`
}

func (p *Personal) Load(g Getter) {
	p.ID, _ = strconv.ParseInt(g.Get("id"), 10, 64)
	p.Name = g.Get("name")
	p.Family = g.Get("family")
	p.NationalCode = g.Get("national_code")
	p.IDNumber = g.Get("id_number")
	p.Sex, _ = strconv.ParseBool(g.Get("sex"))
	pt, _ := strconv.ParseInt(g.Get("person_type"), 10, 64)
	p.PersonType = PersonType(pt)

	p.SecretCode = g.Get("secret_code")
}

func (p *Personal) Validate() error {
	if p.IDNumber != "" {
		if !validateID(p.NationalCode) {
			return fmt.Errorf("This id number is not valid: %s", p.NationalCode)
		}
	}
	return validator.New().Struct(p)
}

func (p *Personal) Initialize(db *gorm.DB) {
	if p.PersonType == needy {
		pass := p.Name + p.Family + p.NationalCode
		h := sha1.New()
		h.Write([]byte(pass))
		p.SecretCode = hex.EncodeToString(h.Sum(nil)[:10])
	}
}

func (p *Personal) Find(db *gorm.DB) ([]Model, error) {
	result := []Personal{}
	if err := db.Find(&result, p).Error; err != nil {
		return nil, err
	}
	ret := make([]Model, len(result))
	for i := range result {
		ret[i] = &result[i]
	}
	return ret, nil
}

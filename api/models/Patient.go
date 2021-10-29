package models

import (
	"fmt"
	"github.com/free-health/health24-gateway/parser"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	"strings"
	"time"
)

type Patient struct {
	ID           uint       `gorm:"primary_key;auto_increment" json:"id"`
	FirstName    string     `gorm:"size:255;not null" json:"first_name"`
	LastName     string     `gorm:"size:255;not null" json:"last_name"`
	MRN          string     `gorm:"size:255;not null;unique_index" json:"mrn"`
	CreatedAt    time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt    time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DischargedAt *time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"discharged_at"`
	Discharged   bool       `gorm:"default:false" json:"discharged"`
	Notes        string     `json:"notes"`
}

func (p *Patient) FromPatientMessage(patient *parser.Patient) {
	p.ID = 0
	p.MRN = patient.MRN
	p.FirstName = strings.ToUpper(patient.FirstName)
	p.LastName = strings.ToUpper(patient.LastName)
	p.Discharged = false
	p.DischargedAt = nil
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
	p.Notes = ""
}

func (p *Patient) SaveIfNotExists(db *gorm.DB) error {
	pp := &Patient{}
	result := db.Where("mrn = ?", p.MRN).First(&pp)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return fmt.Errorf("error finding record: %s", result.Error)
	}
	if result.Error == gorm.ErrRecordNotFound {
		// no record, we add it
		if db.Create(&p).Error != nil {
			return fmt.Errorf("error adding record: %s", result.Error)
		}
		log.Infof("patient added %s", p.MRN)
		return nil
	}

	log.Infof("%s patient exists", p.MRN)
	return nil
}

func (Patient) GetAll(db *gorm.DB, limit int, offset int) (*[]Patient, error) {
	var patients []Patient
	var err error

	err = db.Model(&Patient{}).Limit(limit).Offset(offset).Find(&patients).Error
	if err != nil {
		return &[]Patient{}, err
	}
	return &patients, nil
}

func (Patient) GetDischarged(db *gorm.DB, limit int, offset int) (*[]Patient, error) {
	var patients []Patient
	var err error

	err = db.Model(&Patient{}).Where("discharged = ?", true).Limit(limit).Offset(offset).Find(&patients).Error
	if err != nil {
		return &[]Patient{}, err
	}
	return &patients, nil
}

func (Patient) GetCurrentPatients(db *gorm.DB, limit int, offset int) (*[]Patient, error) {
	var patients []Patient
	var err error

	err = db.Model(&Patient{}).Where("discharged = ?", false).Limit(limit).Offset(offset).Find(&patients).Error
	if err != nil {
		return &[]Patient{}, err
	}
	return &patients, nil
}

func (p *Patient) Discharge(db *gorm.DB, id uint) error {
	return db.Model(&Patient{}).Where("id = ?", id).Update("discharged", true, "discharged_at", time.Now()).Error
}

func (p *Patient) Update(db *gorm.DB) error {
	return db.Save(&p).Error
}

func (p *Patient) Delete(db *gorm.DB, id uint) error {
	err := db.Model(&Patient{}).Where("id = ?", id).Take(&Patient{}).Delete(&Patient{}).Error
	if err != nil {
		return err
	}

	// delete alarms
	err = db.Model(&PhysiologicalAlarm{}).Where("mrn = ?", p.MRN).Delete(&PhysiologicalAlarm{}).Error
	if err != nil {
		return err
	}

	return nil
}

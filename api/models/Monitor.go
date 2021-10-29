package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

type Monitor struct {
	ID   uint   `gorm:"primary_key;auto_increment" json:"id"`
	Name string `gorm:"size:25;not null;" json:"name"`
	Room string `gorm:"size:25;not null;" json:"room"`
}

func (m *Monitor) Prepare(room string, name string) {
	m.Name = name
	m.Room = room
	m.ID = 0
}

func (m *Monitor) SaveIfNotExist(db *gorm.DB) error {
	mm := &Monitor{}
	result := db.Where("room = ? AND name = ?", m.Room, m.Name).First(&mm)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return result.Error
	}

	if result.Error == gorm.ErrRecordNotFound {
		// no record, we add it
		if db.Create(&m).Error != nil {
			return fmt.Errorf("error adding record: %s", result.Error)
		}
	}
	log.Infof("%s:%s exists", m.Room, m.Name)

	return nil
}

func (Monitor) GetAll(db *gorm.DB) (*[]Monitor, error) {
	var monitors []Monitor
	err := db.Model(&Monitor{}).Find(&monitors).Error
	if err != nil {
		return nil, err
	}
	return &monitors, nil
}

func (Monitor) Delete(db *gorm.DB, id uint) error {
	return db.Model(&Monitor{}).Where("id = ?", id).Take(&Monitor{}).Delete(&Monitor{}).Error
}

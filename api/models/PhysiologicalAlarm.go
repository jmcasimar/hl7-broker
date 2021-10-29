package models

import (
	"fmt"
	"github.com/free-health/health24-gateway/parser"
	"github.com/jinzhu/gorm"
	"strconv"
	"strings"
	"time"
)

type PhysiologicalAlarm struct {
	ID      uint      `gorm:"primary_key;auto_increment" json:"id"`
	MRN     string    `gorm:"size:255;not null;index" json:"mrn"`
	Level   string    `gorm:"size:10;" json:"level"`
	AlarmID uint32    `json:"alarm_id"`
	Message string    `json:"message"`
	Time    time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"time"`
}

func getLevel(l string) string {
	switch strings.TrimSpace(l) {
	case "1":
		return "high"
	case "2":
		return "medium"
	case "3":
		return "low"
	case "4":
		fallthrough
	default:
		return "message"
	}
}

func (p *PhysiologicalAlarm) FromPhyAlarmMessage(mrn string, alarm *parser.PhysiologicalAlarm) error {
	p.MRN = mrn
	p.ID = 0
	p.Time = alarm.Time
	alarmId, err := strconv.Atoi(alarm.AlarmID)
	if err != nil {
		return fmt.Errorf("error converting to integer: %s", err)
	}
	p.AlarmID = uint32(alarmId)
	p.Level = getLevel(alarm.Level)
	p.Message = alarm.Message

	return nil
}

func (p *PhysiologicalAlarm) Save(db *gorm.DB) error {
	var err error
	err = db.Create(&p).Error
	if err != nil {
		return err
	}
	return nil
}

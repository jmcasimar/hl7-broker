package realtime

import (
	"github.com/free-health/health24-gateway/api/db"
	"github.com/free-health/health24-gateway/api/models"
	"github.com/free-health/health24-gateway/parser"
	log "github.com/sirupsen/logrus"
)

var PersistPhysiologicalAlarms chan *parser.Record

func init() {
	PersistPhysiologicalAlarms = make(chan *parser.Record)

	go func() {
		for {
			select {
			case values := <-PersistPhysiologicalAlarms:
				mrn := values.MRN
				for _, v := range values.PhysiologicalAlarms {
					alarm := &models.PhysiologicalAlarm{}
					err := alarm.FromPhyAlarmMessage(mrn, v)
					if err != nil {
						log.Errorf("error processing alarm %s", err)
						continue
					}
					err = alarm.Save(db.DB)
					if err != nil {
						log.Errorf("error saving alarm %s", err)
						continue
					}
				}
			}
		}

	}()
}

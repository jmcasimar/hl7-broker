package parser

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"time"

	"github.com/free-health/hl7"
)

const (
	TimeFormat = "20060102150405"

	MessageTypeORU    = "ORU"
	MessageSubTypeR01 = "R01"

	ControlTypePeriodicParam            = "204"
	ControlTypeNIBP                     = "503" // Aperiodic
	ControlTypePhysiologicalAlarm       = "54"
	ControlTypeTechnicalAlarm           = "56"
	ControlTypeModuleLoading            = "11"
	ControlTypeAlarmLimit               = "51"
	ControlTypeAlarmLevel               = "58"
	ControlTypeModuleUnload             = "12"
	ControlTypeParameterLoadUnload      = "1202"
	ControlTypeEcho                     = "106"
	ControlTypePatientInformationChange = "103"
)

func strToTime(date string) (time.Time, error) {
	t, err := time.Parse(TimeFormat, date)

	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}

func ExtractPatientData(message hl7.Message) (*Patient, error) {
	p := &Patient{}

	// get mrn
	mrn, err := GetOBXIdentifierFieldValue(message, "2301")
	if err != nil {
		return nil, fmt.Errorf("error obtaining medical record number %s", err)
	}
	if len(mrn) == 0 {
		return nil, fmt.Errorf("empty MRN")
	}
	p.MRN = mrn

	// get firstname
	fnameQ, err := hl7.ParseQuery("PID-5-2")
	if err != nil {
		return nil, fmt.Errorf("error firstname: %s", err)
	}
	p.FirstName = fnameQ.GetString(message)

	// get lastname
	lnameQ, err := hl7.ParseQuery("PID-5-1")
	if err != nil {
		return nil, fmt.Errorf("error lasttname: %s", err)
	}
	p.LastName = lnameQ.GetString(message)

	return p, nil
}

func ExtractDeviceData(message hl7.Message) (port string, room string, monitorId string, err error) {
	portQ, err := hl7.ParseQuery("PV1-3-3-4")
	if err != nil {
		return "", "", "", fmt.Errorf("error reading TCP Port from PV1: %s", err)
	}
	port = portQ.GetString(message)
	if len(port) == 0 {
		return "", "", "", fmt.Errorf("invalid port %s", port)
	}

	roomQ, err := hl7.ParseQuery("PV1-3-3-1")
	if err != nil {
		return "", "", "", fmt.Errorf("error reading room from PV1: %s", err)
	}
	room = roomQ.GetString(message)
	if len(room) == 0 {
		return "", "", "", fmt.Errorf("invalid room %s", room)
	}

	val, err := GetOBXIdentifierFieldValue(message, "2304")
	if err != nil {
		return "", "", "", err
	}
	if len(val) == 0 {
		return "", "", "", fmt.Errorf("empty monitor id")
	}
	monitorId = val

	return
}

func GetOBXIdentifierFieldValue(message hl7.Message, field string) (string, error) {
	segs := message.Segments("OBX")
	for i := range segs {
		pos := i + 1
		obxTypeQ, err := hl7.ParseQuery(fmt.Sprintf("OBX(%d)-3-1", pos))
		if err != nil {
			return "", err
		}

		obxType := obxTypeQ.GetString(message)
		if obxType == field {
			obsValueQ, err := hl7.ParseQuery(fmt.Sprintf("OBX(%d)-5", pos))
			if err != nil {
				return "", err
			}

			return obsValueQ.GetString(message), nil
		}
	}

	return "", fmt.Errorf("OBX Field %s not found", field)
}

func ParseAperiodic(message hl7.Message, obxPos int) (*AperiodicParam, error) {
	p := &AperiodicParam{}

	obsValueQ, err := hl7.ParseQuery(fmt.Sprintf("OBX(%d)-5", obxPos))
	if err != nil {
		return nil, err
	}

	p.Value = obsValueQ.GetString(message)

	obsTimeQ, err := hl7.ParseQuery(fmt.Sprintf("OBX(%d)-14", obxPos))
	if err != nil {
		return nil, err
	}
	t, err := strToTime(obsTimeQ.GetString(message))
	if err != nil {
		return nil, err
	}
	p.Time = t

	obsIdentifierQ, err := hl7.ParseQuery(fmt.Sprintf("OBX(%d)-3-1", obxPos))
	if err != nil {
		return nil, err
	}
	p.ParameterID = obsIdentifierQ.GetString(message)

	obsNameQ, err := hl7.ParseQuery(fmt.Sprintf("OBX(%d)-3-2", obxPos))
	if err != nil {
		return nil, err
	}
	p.Parameter = obsNameQ.GetString(message)

	return p, nil
}

func ParsePeriodic(message hl7.Message, obxPos int) (*PeriodicParam, error) {
	p := &PeriodicParam{}

	obsValueQ, err := hl7.ParseQuery(fmt.Sprintf("OBX(%d)-5", obxPos))
	if err != nil {
		return nil, err
	}
	p.Value = obsValueQ.GetString(message)

	obsIdentifierQ, err := hl7.ParseQuery(fmt.Sprintf("OBX(%d)-3-1", obxPos))
	if err != nil {
		return nil, err
	}
	p.ParameterID = obsIdentifierQ.GetString(message)

	obsNameQ, err := hl7.ParseQuery(fmt.Sprintf("OBX(%d)-3-2", obxPos))
	if err != nil {
		return nil, err
	}
	p.Parameter = obsNameQ.GetString(message)

	return p, nil
}

func getAlarmLimitType(t string) string {
	switch t {
	case "2002":
		return "upper"
	case "2003":
		return "lower"
	default:
		return "unknown"
	}
}

func ParseAlarmLimit(message hl7.Message, obxPos int) (*AlarmLimit, error) {
	p := &AlarmLimit{}

	obsTypeQ, err := hl7.ParseQuery(fmt.Sprintf("OBX(%d)-3-1", obxPos))
	if err != nil {
		return nil, err
	}
	p.Type = getAlarmLimitType(obsTypeQ.GetString(message))

	obsParameterQ, err := hl7.ParseQuery(fmt.Sprintf("OBX(%d)-4", obxPos))
	if err != nil {
		return nil, err
	}
	p.Parameter = obsParameterQ.GetString(message)

	obsValueQ, err := hl7.ParseQuery(fmt.Sprintf("OBX(%d)-5", obxPos))
	if err != nil {
		return nil, err
	}
	p.Value = obsValueQ.GetString(message)

	return p, nil
}

func ParsePhyAlarm(message hl7.Message, obxPos int) (*PhysiologicalAlarm, error) {
	a := &PhysiologicalAlarm{}

	obsAlarmIDQ, err := hl7.ParseQuery(fmt.Sprintf("OBX(%d)-5-1", obxPos))
	if err != nil {
		return nil, err
	}
	id := obsAlarmIDQ.GetString(message)
	if id == "0" {
		return nil, fmt.Errorf("empty alram id")
	}
	a.AlarmID = id

	obsLevelQ, err := hl7.ParseQuery(fmt.Sprintf("OBX(%d)-3", obxPos))
	if err != nil {
		return nil, err
	}
	a.Level = obsLevelQ.GetString(message)

	obsAlarmQ, err := hl7.ParseQuery(fmt.Sprintf("OBX(%d)-5-2", obxPos))
	if err != nil {
		return nil, err
	}
	a.Message = obsAlarmQ.GetString(message)

	obsTimeQ, err := hl7.ParseQuery(fmt.Sprintf("OBX(%d)-14", obxPos))
	if err != nil {
		return nil, err
	}
	t, err := strToTime(obsTimeQ.GetString(message))
	if err != nil {
		return nil, err
	}
	a.Time = t

	return a, nil
}

func ParseTechAlarm(message hl7.Message, obxPos int) (*TechnicalAlarm, error) {
	a := &TechnicalAlarm{}
	obsAlarmIDQ, err := hl7.ParseQuery(fmt.Sprintf("OBX(%d)-5-1", obxPos))
	if err != nil {
		return nil, err
	}
	id := obsAlarmIDQ.GetString(message)
	if id == "0" {
		return nil, fmt.Errorf("empty id")
	}
	a.AlarmID = id

	obsLevelQ, err := hl7.ParseQuery(fmt.Sprintf("OBX(%d)-3", obxPos))
	if err != nil {
		return nil, err
	}
	a.Level = obsLevelQ.GetString(message)

	obsAlarmQ, err := hl7.ParseQuery(fmt.Sprintf("OBX(%d)-5-2", obxPos))
	if err != nil {
		return nil, err
	}
	a.Message = obsAlarmQ.GetString(message)

	return a, nil
}

func GetMessageType(message hl7.Message) (string, error) {
	msgTypeQ, err := hl7.ParseQuery("MSH-9-1")
	if err != nil {
		return "", fmt.Errorf("error parsing msh for message type: %s", err)
	}
	return msgTypeQ.GetString(message), nil
}

func GetMessageEventType(message hl7.Message) (string, error) {
	typeQ, err := hl7.ParseQuery("MSH-9-2")
	if err != nil {
		return "", fmt.Errorf("error parsing msh for message event type: %s", err)
	}
	return typeQ.GetString(message), nil
}

func GetControlType(message hl7.Message) (string, error) {
	msgTypeQ, err := hl7.ParseQuery("MSH-10")
	if err != nil {
		return "", fmt.Errorf("error parsing msh for control type: %s", err)
	}
	return msgTypeQ.GetString(message), nil
}

func ParseAlarmLimits(message hl7.Message) ([]*AlarmLimit, error) {
	obxMessages := message.Segments("OBX")
	params := make([]*AlarmLimit, 0)

	for i := range obxMessages {
		p, err := ParseAlarmLimit(message, i+1)
		if err != nil {
			log.Debugf("error parsing alarm limit: %s", err)
			continue
		}
		params = append(params, p)
	}

	return params, nil
}

func ParsePeriodicMessages(message hl7.Message) ([]*PeriodicParam, error) {
	obxMessages := message.Segments("OBX")
	params := make([]*PeriodicParam, 0)

	for i := range obxMessages {
		p, err := ParsePeriodic(message, i+1)
		if err != nil {
			log.Debugf("error parsing periodic param: %s", err)
			continue
		}
		params = append(params, p)
	}

	return params, nil
}

func ParseAperiodicMessages(message hl7.Message) ([]*AperiodicParam, error) {
	obxMessages := message.Segments("OBX")
	params := make([]*AperiodicParam, 0)

	for i := range obxMessages {
		p, err := ParseAperiodic(message, i+1)
		if err != nil {
			log.Debugf("error parsing aperiodic param: %s", err)
			continue
		}
		params = append(params, p)
	}

	return params, nil
}

func ParsePhysiologicalAlarmMessages(message hl7.Message) ([]*PhysiologicalAlarm, error) {
	obxMessages := message.Segments("OBX")
	params := make([]*PhysiologicalAlarm, 0)

	for i := range obxMessages {
		p, err := ParsePhyAlarm(message, i+1)
		if err != nil {
			log.Debugf("error parsing physiological alarm: %s", err)
			continue
		}
		params = append(params, p)
	}

	return params, nil
}

func ParseTechnicalAlarmMessages(message hl7.Message) ([]*TechnicalAlarm, error) {
	obxMessages := message.Segments("OBX")
	params := make([]*TechnicalAlarm, 0)

	for i := range obxMessages {
		p, err := ParseTechAlarm(message, i+1)
		if err != nil {
			log.Debugf("error parsing physiological alarm: %s", err)
			continue
		}
		params = append(params, p)
	}

	return params, nil
}

func Parse(message hl7.Message) (*Record, error) {
	pr := &Record{}
	pr.Time = time.Now()

	// get message type
	messageType, err := GetMessageType(message)
	if err != nil {
		return nil, err
	}

	// get message control type
	messageControlType, err := GetMessageEventType(message)
	if err != nil {
		return nil, err
	}

	if messageType != MessageTypeORU || messageControlType != MessageSubTypeR01 {
		log.Debugf("unsupported message type %s^%s", messageType, messageControlType)
		return nil, nil
	}

	// check for message type
	controlType, err := GetControlType(message)
	if err != nil {
		return nil, err
	}

	// we usually get messages grouped into a module
	switch controlType {
	case ControlTypeEcho:
		log.Debug("echo message received")
		return nil, nil

	case ControlTypePeriodicParam:
		log.Debug("Periodic param received")
		pp, err := ParsePeriodicMessages(message)
		if err != nil {
			return nil, err
		}
		if len(pp) == 0 {
			return nil, nil
		}
		pr.PeriodicParams = pp

	case ControlTypeNIBP:
		log.Debug("aperiodic param received(NIBP)")
		pp, err := ParseAperiodicMessages(message)
		if err != nil {
			return nil, err
		}
		if len(pp) == 0 {
			return nil, nil
		}
		pr.AperiodicParams = pp

	case ControlTypePhysiologicalAlarm:
		log.Debug("physiological alarm received")
		pp, err := ParsePhysiologicalAlarmMessages(message)
		if err != nil {
			return nil, err
		}
		if len(pp) == 0 {
			return nil, nil
		}
		pr.PhysiologicalAlarms = pp

	case ControlTypeTechnicalAlarm:
		log.Debug("technical alarm received")
		pp, err := ParseTechnicalAlarmMessages(message)
		if err != nil {
			return nil, err
		}
		if len(pp) == 0 {
			return nil, nil
		}
		pr.TechnicalAlarms = pp

	case ControlTypeAlarmLimit:
		log.Debug("alarm limit received")
		pp, err := ParseAlarmLimits(message)
		if err != nil {
			return nil, err
		}
		if len(pp) == 0 {
			return nil, nil
		}
		pr.AlarmLimits = pp

	case ControlTypePatientInformationChange:
		log.Debug("patient changed, disconnecting")
		return nil, fmt.Errorf("patient changed")

	default:
		log.Debugf("unsupported control packet type: %s", controlType)
		return nil, nil
	}

	return pr, nil
}

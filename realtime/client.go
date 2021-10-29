package realtime

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/free-health/health24-gateway/parser"
	"github.com/free-health/hl7"
	log "github.com/sirupsen/logrus"
)

func (d *Device) Handle(done chan error, room string, monitorId string, mrn string) {
	var sioRoom = d.DeviceID
	var sioChan = fmt.Sprintf("%s:data", sioRoom)
	for {
		m, err := d.ProcessHL7Packet()
		if err != nil {
			done <- err
		}

		log.Debug("processing record")
		pr, err := parser.Parse(m)
		if err != nil {
			log.Errorf("error building hl7: %s", err)
			continue
		}

		if pr != nil {
			pr.MRN = mrn
			pr.MonitorID = monitorId

			data, err := json.Marshal(pr)
			if err != nil {
				log.Errorf("error marshaling json: %s", err)
				continue
			}
			//fmt.Println("Parsed JSON:")
			//fmt.Println(string(data))
			SocketIO.BroadcastToRoom("/", sioRoom, sioChan, string(data))

			if pr.AlarmLimits != nil {
				// we get device data
				// store it for later purposes
				d.Meta.AddAlarms(pr.AlarmLimits)
			}

			if pr.PhysiologicalAlarms != nil {
				PersistPhysiologicalAlarms <- pr
			}
		}
	}
}

var cutset = string([]byte{0x0B, 0x1C, 0x0D})

func (d *Device) ProcessHL7Packet() (hl7.Message, error) {
	// read message start 0x0B
	b, err := d.ReadByte()
	if err != nil {
		return nil, fmt.Errorf("error reading start byte: %s", err)
	}
	if b != byte(0x0B) {
		return nil, fmt.Errorf("invalid header")
	}

	// read payload
	payloadWithDelimiter, err := d.ReadBytes(byte(0x1C))
	if err != nil {
		return nil, fmt.Errorf("error reading payload: %s", err)
	}

	b, err = d.ReadByte()
	if err != nil {
		return nil, fmt.Errorf("error reading end byte %s", err)
	}
	if b != byte(0x0D) {
		return nil, fmt.Errorf("invalid message end")
	}

	//fmt.Println("Payload with delim as a string:")
	//fmt.Println(strings.ReplaceAll(string(payloadWithDelimiter), "\r", "\n"))

	// skip last two bytes from the hl7 packet
	payload := payloadWithDelimiter[:len(payloadWithDelimiter)-1]
	log.Debugf("Length of payload %d\n", len(payload))
	//if payload[0] == byte(0x0B) {
	//	payload = payload[1:]
	//}

	payload = bytes.Trim(payload, cutset)

	//fmt.Println("Payload as a string:")
	//fmt.Println(strings.ReplaceAll(string(payload), "\r", "\n"))

	m, _, err := hl7.ParseMessage(payload)
	if err != nil {
		return nil, fmt.Errorf("error parsing hl7: %s\n", err)
	}

	fmt.Println("Decoded Payload:")
	fmt.Println(m)
	fmt.Println()
	return m, err
}

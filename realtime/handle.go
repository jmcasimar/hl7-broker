package realtime

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/free-health/health24-gateway/api/db"
	"github.com/free-health/health24-gateway/api/models"
	"github.com/free-health/health24-gateway/parser"
	"net"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	QryTimeFormat = "20060102150405"

	HL7EchoMessage = "\"MSH|^~\\\\&|||||||ORU^R01|106|P|2.3.1|\""
)

func BuildHL7Packet(msg string) *bytes.Buffer {
	buf := &bytes.Buffer{}
	buf.WriteByte(byte(0x0B))
	buf.WriteString(msg)
	buf.Write([]byte{byte(0x1C), byte(0x0D)})

	return buf
}

func buildHandshakeMessage() string {
	msg := fmt.Sprintf("MSH|^~\\&|||||||QRY^R02|1203|P|2.3.1\r"+
		"QRD|%s|R|I|Q895211|||||RES\r"+
		"QRF|MON||||0&0^1^1^1^\r"+
		"QRF|MON||||0&0^3^1^1^\r"+
		"QRF|MON||||0&0^4^1^1^\r", time.Now().Format(QryTimeFormat))
	return msg
}

func (d *Device) SendHl7Message(msg string) error {
	_, err := BuildHL7Packet(msg).WriteTo(d)
	if err != nil {
		return err
	}
	err = d.Flush()
	return err
}

func (d *Device) sendEcho(done chan error) {
	for {
		log.Debugf("sending echo to %s", d.Host)
		err := d.SendHl7Message(HL7EchoMessage)
		if err != nil {
			done <- fmt.Errorf("error sending echo to %s error: %s", d.Host, err)
			break
		}
		time.Sleep(time.Second * 1)
	}
}

func Handle(host string, port string, room string, monitorId string) {
	addr := net.JoinHostPort(host, port)
	log.WithField("host", host).Debugf("initiating connection with %s", addr)

	// add a nil record until we initialize connection
	Current.Add(host, nil)
	defer Current.Remove(host)

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		// we cannot init conn
		log.Errorf("error connecting to %s", addr)
		return
	}

	device := NewDevice(
		conn,
		bufio.NewReader(conn),
		bufio.NewWriter(conn),
		host,
		fmt.Sprintf("%s/%s", room, monitorId),
	)
	Current.Add(host, device)

	log.WithField("host", host).Debugf("waiting 2s before sending data")
	time.Sleep(time.Second * 2)

	err = device.SendHl7Message(buildHandshakeMessage())
	if err != nil {
		log.WithField("host", host).Errorf("error sending handshaking to %s", addr)
		return
	}

	// read patient change data
	msg, err := device.ProcessHL7Packet()
	if err != nil {
		log.WithField("host", host).Errorf("error processing hl7 packet: %s", err)
		return
	}
	patient, err := parser.ExtractPatientData(msg)
	if err != nil {
		log.WithField("host", host).Errorf("error processing patient data: %s", err)
		return
	}

	// register patient to db
	patientRecord := &models.Patient{}
	patientRecord.FromPatientMessage(patient)
	err = patientRecord.SaveIfNotExists(db.DB)
	if err != nil {
		log.WithField("host", host).Errorf("error saving patient data: %s", err)
		return
	}

	doneChan := make(chan error)

	// start sending echo
	go device.sendEcho(doneChan)

	// start handler
	go device.Handle(doneChan, room, monitorId, patient.MRN)

	select {
	case err := <-doneChan:
		log.WithField("host", host).Errorf("error: %s", err)
		_ = conn.Close()
		return
	}
}

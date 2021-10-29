package discovery

import (
	"github.com/free-health/health24-gateway/api/db"
	"github.com/free-health/health24-gateway/api/models"
	"github.com/free-health/health24-gateway/parser"
	"github.com/free-health/health24-gateway/realtime"
	"github.com/free-health/hl7"
	log "github.com/sirupsen/logrus"
	"net"
	"strings"
)

func HandleBroadcasts(pc net.PacketConn) {
	for {
		buf := make([]byte, 1024)
		n, addr, err := pc.ReadFrom(buf)
		if err != nil {
			log.Errorf("error reading udp packet: %s", err)
		}

		parts := strings.Split(addr.String(), ":")
		if realtime.Current.Has(parts[0]) {
			// device is already registered
			continue
		}

		if buf[0] != byte(0x0B) || buf[n-2] != byte(0x1C) || buf[n-1] != byte(0x0D) {
			log.Errorf("invalid packet type received from: %s", addr.String())
			continue
		}

		log.Infof("discovered device IP %s", parts[0])

		msg := buf[1 : n-2]
		m, _, err := hl7.ParseMessage(msg)
		if err != nil {
			log.Errorf("error parsing hl7 message: %s", err)
			continue
		}

		port, room, monitorId, err := parser.ExtractDeviceData(m)
		if err != nil {
			log.Errorf("error: %s", err)
			continue
		}

		// save to db
		mon := &models.Monitor{}
		mon.Prepare(room, monitorId)
		err = mon.SaveIfNotExist(db.DB)
		if err != nil {
			log.Errorf("error saving monitor: %s", err)
			continue
		}

		go realtime.Handle(parts[0], port, room, monitorId)
	}
}

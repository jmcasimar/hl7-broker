package discovery

import (
	"fmt"
	"net"

	log "github.com/sirupsen/logrus"
)

var packetConn net.PacketConn

func Start(host string) {
	pc, err := net.ListenPacket("udp4", fmt.Sprintf("%s:4600", host))
	packetConn = pc
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("listening on UDP %s:4600", host)

	go HandleBroadcasts(packetConn)
}

func Stop() {
	_ = packetConn.Close()
}

package realtime

import (
	"bufio"
	"github.com/free-health/health24-gateway/parser"
	"net"
	"sync"

	log "github.com/sirupsen/logrus"
)

type DeviceMeta struct {
	mux         *sync.RWMutex
	AlarmLimits []*parser.AlarmLimit `json:"alarm_limits"`
}

func (d *DeviceMeta) AddAlarms(limits []*parser.AlarmLimit) {
	d.mux.Lock()
	defer d.mux.Unlock()

	if d.AlarmLimits == nil {
		d.AlarmLimits = limits
		return
	}

	var newValues []*parser.AlarmLimit
	for _, newVal := range limits {
		found := false
		for _, old := range d.AlarmLimits {
			if newVal.Parameter == old.Parameter && newVal.Type == old.Type {
				old.Value = newVal.Value
				found = true
				break
			}
		}
		if !found {
			newValues = append(newValues, newVal)
		}
	}
	d.AlarmLimits = append(d.AlarmLimits, newValues...)
}

type Device struct {
	Connection net.Conn
	Host       string
	*bufio.Reader
	*bufio.Writer
	Meta     DeviceMeta
	DeviceID string // roomId/monitorId
}

func NewDevice(connection net.Conn, reader *bufio.Reader, writer *bufio.Writer, host string, deviceId string) *Device {
	return &Device{
		Connection: connection,
		Reader:     reader,
		Writer:     writer,
		Host:       host,
		Meta: DeviceMeta{
			AlarmLimits: nil,
			mux:         &sync.RWMutex{},
		},
		DeviceID: deviceId,
	}
}

type Devices struct {
	list map[string]*Device
	mux  sync.RWMutex
}

var Current Devices

func init() {
	Current = Devices{list: make(map[string]*Device)}
}

func (d *Devices) Add(ip string, device *Device) {
	d.mux.Lock()
	defer d.mux.Unlock()
	d.list[ip] = device
}

func (d *Devices) Get(ip string) (*Device, bool) {
	d.mux.RLock()
	defer d.mux.RUnlock()

	value, found := d.list[ip]
	return value, found
}

func (d *Devices) Remove(ip string) {
	d.mux.Lock()
	defer d.mux.Unlock()

	delete(d.list, ip)
	log.Debugf("%s disconnected", ip)
}

func (d *Devices) Has(ip string) bool {
	d.mux.RLock()
	defer d.mux.RUnlock()

	_, found := d.list[ip]
	return found
}

func (d *Devices) GetMetaById(id string) (DeviceMeta, bool) {
	d.mux.RLock()
	defer d.mux.RUnlock()

	for _, v := range d.list {
		if v.DeviceID == id {
			return v.Meta, true
		}
	}
	return DeviceMeta{}, false
}

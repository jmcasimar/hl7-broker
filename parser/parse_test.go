package parser

import (
	"encoding/json"
	"fmt"
	"github.com/free-health/hl7"
	"reflect"
	"strings"
	"testing"
)

// TODO: these aren't real tests ofc
// we'll write em later :P

//func Test_extractDeviceData(t *testing.T) {
//	const msg = `MSH|^~\&|||||||ADT^A01|101|P|2.3.1|
//EVN||00000000|
//PID|||14140f00-7bbc-0478-11122d2d02000000||rep^wet|
//PV1||I|^^icu&222&&4601&&1|||||||||||||||U|4294967040||
//OBX||ST|2304^MonitorName||M0001||||||F
//OBX||CE|2305^||0^||||||F
//OBX||CE|2306^||2^||||||F
//OBX||CE|4526^||1^||||||F
//OBX||CE|2307^||1^||||||F
//OBX||NM|2211^||0||||||F
//OBX||NM|4524^||0||||||F
//OBX||ST|2308^BedNoStr||222||||||F
//OBX||ST|4527^||000000000000000000000000||||||F
//OBX||CE|4528^||16^||||||F
//OBX||ST|4529^||001005000004||||||F
//OBX||CE|4530^||1^||||||F
//OBX||ST|2319^||001000015020020188123000000000000000000000000000||||||F
//OBX||CE|2320^||7^||||||F`
//
//	mk := strings.ReplaceAll(msg, "\n", "\r")
//
//	m, _, err := hl7.ParseMessage([]byte(mk))
//	if err != nil {
//		t.Error(err)
//	}
//
//	type args struct {
//		message hl7.Message
//	}
//	tests := []struct {
//		name          string
//		args          args
//		wantPort      string
//		wantMonitorId string
//		wantErr       bool
//	}{
//		{
//			name: "test all",
//			args: args{
//				message: m,
//			},
//			wantPort:      "4601",
//			wantMonitorId: "M0001",
//			wantErr:       false,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			gotPort,room, gotMonitorId, err := ExtractDeviceData(tt.args.message)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("extractDeviceData() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if gotPort != tt.wantPort {
//				t.Errorf("extractDeviceData() gotPort = %v, want %v", gotPort, tt.wantPort)
//			}
//			if gotMonitorId != tt.wantMonitorId {
//				t.Errorf("extractDeviceData() gotMonitorId = %v, want %v", gotMonitorId, tt.wantMonitorId)
//			}
//		})
//	}
//}

func TestParsePeriodic(t *testing.T) {
	msg := "MSH|^~\\&|||||||ORU^R01|204|P|2.3.1|\r" +
		"OBX||NM|101^HR|2101|60||||||F\r" +
		"OBX||NM|102^PVCs|2101|0||||||F\r" +
		"OBX||NM|105^I|2101|-100.00||||||F\r" +
		"OBX||NM|106^II|2101|-100.00||||||F\r" +
		"OBX||NM|107^III|2101|-100.00||||||F\r" +
		"OBX||NM|108^aVR|2101|-100.00||||||F\r" +
		"OBX||NM|109^aVL|2101|-100.00||||||F\r" +
		"OBX||NM|110^aVF|2101|-100.00||||||F\r" +
		"OBX||NM|117^ST-V|2101|-100.00||||||F\r"

	m, _, _ := hl7.ParseMessage([]byte(msg))

	r, err := Parse(m)
	if err != nil {
		t.Fatal(err)
	}

	b, _ := json.MarshalIndent(r, "", "  ")
	fmt.Println(string(b))
}

func TestParseAperiodic(t *testing.T) {
	msg := "MSH|^~\\&|||||||ORU^R01|503|P|2.3.1|\r" +
		"OBX||NM|171^Dia|2105|80||||||F||APERIODIC|20070106191915\r" +
		"OBX||NM|172^Mean|2105|93||||||F||APERIODIC|20070106191915\r" +
		"OBX||NM|170^Sys|2105|120||||||F||APERIODIC|20070106191915\r" +
		"OBX||NM|170^Sys|2105|-100||||||F||APERIODIC|00000000000000\r"

	m, _, _ := hl7.ParseMessage([]byte(msg))

	r, err := Parse(m)
	if err != nil {
		t.Fatal(err)
	}

	b, _ := json.MarshalIndent(r, "", "  ")
	fmt.Println(string(b))
}

func TestParsePhyAlarm(t *testing.T) {
	//msg := "MSH|^~\\&|||||||ORU^R01|54|P|2.3.1|\r" +
	//	"OBX||CE|2|1|10033^**SpO2 TOO HIGH||||||F|||20070106193145|\r" +
	//	"OBX||CE|2|1|10170^**Art-Sys TOO HIGH ||||||F|||20070106193145|\r" +
	//	"OBX||CE|2|1|10172^**Art-Mean TOO HIGH ||||||F|||20070106193145|\r" +
	//	"OBX||CE|2|1|10174^**Art-Dia TOO HIGH ||||||F|||20070106193145|\r" +
	//	"OBX||CE|2|1|10302^**CVP-Mean TOO HIGH ||||||F|||20070106193145|\r" +
	//	"OBX||CE|2|1|10002^**HR TOO LOW||||||F|||20070106193145|\r" +
	//	"OBX||CE|2|1|10044^**RR TOO LOW||||||F|||20070106193145|\r"

	msg2 := "MSH|^~\\&|||||||ORU^R01|54|P|2.3.1|\rOBX||CE|2|1|10076^**T1 Too Low||||||F|||20200418190201|\rOBX||CE|2040^||0^||||||F"

	m, _, _ := hl7.ParseMessage([]byte(msg2))

	r, err := Parse(m)
	if err != nil {
		t.Fatal(err)
	}

	b, _ := json.MarshalIndent(r, "", "  ")
	fmt.Println(string(b))
}

func TestParseTechAlarm(t *testing.T) {
	msg := "MSH|^~\\&|||||||ORU^R01|56|P|2.3.1|\r" +
		"OBX||CE|4||205^SpO2 Sensor Off||||||F|\r" +
		"OBX||CE|4||5^ECG Lead Off||||||F|\r" +
		"OBX||CE|2041^||0^||||||F\r" +
		"OBX||CE|2041^||0^||||||F\r"

	m, _, _ := hl7.ParseMessage([]byte(msg))

	r, err := Parse(m)
	if err != nil {
		t.Fatal(err)
	}

	b, _ := json.MarshalIndent(r, "", "  ")
	fmt.Println(string(b))
}

func TestExtractPatientData(t *testing.T) {
	const msg = `MSH|^~\&|||||||ORU^R01|103|P|2.3.1|
PID|||14140f00-7bbc-0478-11122d2d02000000||WEERASINGHE^KESHARA||19960714|M|
PV1||I|^^ICU&1&3232237756&4601&&1|||||||||||||||A|||
OBR||||Mindray Monitor|||0|
OBX||NM|52^||189.0||||||F
OBX||NM|51^||100.0||||||F
OBX||ST|2301^||2||||||F
OBX||CE|2302^Blood||0^N||||||F
OBX||CE|2303^Paced||2^||||||F
OBX||ST|2308^BedNoStr||BED-001||||||F`

	mk := strings.ReplaceAll(msg, "\n", "\r")

	m, _, err := hl7.ParseMessage([]byte(mk))
	if err != nil {
		t.Error(err)
	}

	type args struct {
		message hl7.Message
	}
	tests := []struct {
		name    string
		args    args
		want    *Patient
		wantErr bool
	}{
		{
			name: "should parse correct",
			args: args{
				message: m,
			},
			want: &Patient{
				FirstName: "KESHARA",
				LastName:  "WEERASINGHE",
				MRN:       "2",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExtractPatientData(tt.args.message)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExtractPatientData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ExtractPatientData() got = %v, want %v", got, tt.want)
			}
		})
	}
}

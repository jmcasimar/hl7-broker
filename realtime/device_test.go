package realtime

import (
	"fmt"
	"github.com/free-health/health24-gateway/parser"
	"testing"
)

// note: write the proper test
func TestDeviceMeta_AddAlarms(t *testing.T) {
	initial := &DeviceMeta{
		AlarmLimits: []*parser.AlarmLimit{
			{
				Type:      "high",
				Parameter: "1",
				Value:     "2",
			},
			{
				Type:      "low",
				Parameter: "2",
				Value:     "3",
			},
			{
				Type:      "low",
				Parameter: "4",
				Value:     "3",
			},
		},
	}

	newVals := &DeviceMeta{
		AlarmLimits: []*parser.AlarmLimit{
			{
				Type:      "low",
				Parameter: "1",
				Value:     "40",
			},
			{
				Type:      "low",
				Parameter: "2",
				Value:     "30",
			},
		},
	}

	initial.AddAlarms(newVals.AlarmLimits)

	fmt.Println(len(initial.AlarmLimits))
	for _, v := range initial.AlarmLimits {
		fmt.Println(v)
	}
}

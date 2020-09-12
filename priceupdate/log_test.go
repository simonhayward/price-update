package priceupdate

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"
)

func TestLog(t *testing.T) {
	testCases := []struct {
		log      Log
		expected []byte
	}{
		{
			log: Log{
				Message:  "boom",
				Severity: Error,
			},
			expected: []byte(`{"message":"boom","severity":"error"}`),
		},
		{
			log: Log{
				Message:  `no match for regex: Sell <span>\s+([.0-9]+)\s+GBX`,
				Severity: Critical,
			},
			expected: []byte(`{"message":"no match for regex: Sell \u003cspan\u003e\\s+([.0-9]+)\\s+GBX","severity":"critical"}`),
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("checking %s", tc.expected), func(t *testing.T) {
			log, err := json.Marshal(tc.log)
			if err != nil {
				t.Fatalf("%s", err)
			}
			if bytes.Equal(log, tc.expected) != true {
				t.Errorf("got %s; want %s", log, tc.expected)
			}
		})
	}

}

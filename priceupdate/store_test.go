package priceupdate

import (
	"fmt"
	"reflect"
	"testing"
)

func TestStore(t *testing.T) {
	testCases := []struct {
		body     []byte
		expected Rows
	}{
		{
			body: []byte(`1234.5,09342434,2006-01-02T15:04:05-0700
6542.98,095454,2006-01-02T15:04:05-0700
8754.32,010101,2006-01-02T15:04:05-0700`),
			expected: Rows{
				"09342434": &Row{
					Index:   0,
					Price:   "1234.5",
					Updated: "2006-01-02T15:04:05-0700",
				},
				"095454": &Row{
					Index:   1,
					Price:   "6542.98",
					Updated: "2006-01-02T15:04:05-0700",
				},
				"010101": &Row{
					Index:   2,
					Price:   "8754.32",
					Updated: "2006-01-02T15:04:05-0700",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("checking %v", tc.expected), func(t *testing.T) {
			result, err := GenerateRows(tc.body)
			if err != nil {
				t.Fatalf("unexpected nil %s", err)
			}
			if reflect.DeepEqual(result, tc.expected) == false {
				t.Errorf("got %v; want %v", result, tc.expected)
			}
		})
	}

}

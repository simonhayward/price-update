package priceupdate

import (
	"fmt"
	"reflect"
	"testing"
)

func TestSecurities(t *testing.T) {
	testCases := []struct {
		body     []byte
		expected Securities
	}{
		{
			body: []byte(`[
				{
					"isin": "01",
					"src": [
						{
							"url": "https://1.com/data?g=01",
							"p": "Price</span><span>([.0-9]+)</span>"
						},
						{
							"url": "https://2.com/01",
							"p": "buy/sell price</div><p>([.0-9]+)p</p>",
							"c": true
						}
					]
				},
				{
					"isin": "02",
					"src": [
						{
							"url": "https://1.com/data?g=02",
							"p": "Price</span><span>([.0-9]+)</span>"
						},
						{
							"url": "https://2.com/02",
							"p": "price</div><span>([.0-9]+)p</span>"
						}
					]
				}
			]`),
			expected: Securities{
				Security{
					ISIN: "01",
					Sources: Sources{
						Source{
							URL:     "https://1.com/data?g=01",
							Pattern: "Price</span><span>([.0-9]+)</span>",
						},
						Source{
							URL:     "https://2.com/01",
							Pattern: "buy/sell price</div><p>([.0-9]+)p</p>",
							Convert: true,
						},
					},
				},
				Security{
					ISIN: "02",
					Sources: Sources{
						Source{
							URL:     "https://1.com/data?g=02",
							Pattern: "Price</span><span>([.0-9]+)</span>",
						},
						Source{
							URL:     "https://2.com/02",
							Pattern: "price</div><span>([.0-9]+)p</span>",
						},
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("checking %v", tc.expected), func(t *testing.T) {
			result, err := GenerateSecurities(tc.body)
			if err != nil {
				t.Fatalf("unexpected nil %s", err)
			}
			if reflect.DeepEqual(*result, tc.expected) == false {
				t.Errorf("got %v; want %v", result, tc.expected)
			}
		})
	}

}

func TestSourceParse(t *testing.T) {
	testCases := []struct {
		source   Source
		expected string
		err      error
	}{
		{
			source: Source{
				Response: []byte("<page><price>10.1</price></page>"),
				Pattern:  "<page><price>([.0-9]+)</price></page>",
			},
			expected: "10.1",
		},
		{
			source: Source{
				Response: []byte("<page><price>7.30</price></page>"),
				Pattern:  "<page><price>([.0-9]+)</price></page>",
			},
			expected: "7.30",
		},
		{
			source: Source{
				Response: []byte("<td class='line text'>GBX 313.23</td>"),
				Pattern:  "GBX\xC2\xA0([.0-9]+)",
			},
			expected: "313.23",
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("checking %s", tc.expected), func(t *testing.T) {
			price, err := tc.source.Parse()
			if err != nil {
				t.Fatalf("%s", err)
			}
			if price != tc.expected {
				t.Errorf("got %s; want %s", tc.source.Price, tc.expected)
			}
		})
	}

}

func TestSourceCleanPrice(t *testing.T) {
	testCases := []struct {
		source   Source
		expected string
		err      error
	}{
		{
			source: Source{
				Price: "1,000.1",
			},
			expected: "1000.1",
		},
		{
			source: Source{
				Price: "7000,30",
			},
			expected: "700030",
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("checking %s", tc.expected), func(t *testing.T) {
			tc.source.CleanPrice()
			if tc.source.Price != tc.expected {
				t.Errorf("got %s; want %s", tc.source.Price, tc.expected)
			}
		})
	}

}

func TestSourceConvertPrice(t *testing.T) {
	testCases := []struct {
		source   Source
		expected string
		err      error
	}{
		{
			source: Source{
				Price:   "10.1",
				Convert: true,
			},
			expected: "1010.0",
		},
		{
			source: Source{
				Price: "7.30",
			},
			expected: "7.30",
		},
		{
			source: Source{
				Price:   ".12",
				Convert: true,
			},
			expected: "12.00",
		},
		{
			source: Source{
				Price:   "7.2",
				Convert: true,
			},
			expected: "720.0",
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("checking %s", tc.expected), func(t *testing.T) {
			tc.source.ConvertPrice()
			if tc.source.Price != tc.expected {
				t.Errorf("got %s; want %s", tc.source.Price, tc.expected)
			}
		})
	}

}

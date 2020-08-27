package priceupdate

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"strconv"
)

type Row struct {
	Index   int // Maintain order of row on output
	Price   string
	Updated string
}

type Rows map[string]*Row

// GetRows - from response body convert into rows
func GetRows(url string) (Rows, error) {
	body, err := GetResponse(url)
	if err != nil {
		return nil, err
	}
	return GenerateRows(body)
}

// GenerateRows - convert from CSV into rows
func GenerateRows(body []byte) (Rows, error) {
	lines, err := csv.NewReader(bytes.NewReader(body)).ReadAll()
	if err != nil {
		return nil, err
	}

	rows := make(map[string]*Row)
	for i, line := range lines {
		rows[line[1]] = &Row{Index: i, Price: line[0], Updated: line[2]}
	}

	return rows, nil
}

// IndexOrder - maintain the row order for CSV output
func (rs Rows) IndexOrder() []string {
	keys := make([]string, len(rs))
	for key, row := range rs {
		keys[row.Index] = key
	}
	return keys
}

// AsBytes - in CSV format
func (rs Rows) AsBytes() ([]byte, error) {
	b := &bytes.Buffer{}
	w := csv.NewWriter(b)

	for _, isin := range rs.IndexOrder() {
		line := []string{rs[isin].Price, isin, rs[isin].Updated}
		if err := w.Write(line); err != nil {
			return nil, fmt.Errorf("error writing to csv: %s", err)
		}
	}
	w.Flush()
	if err := w.Error(); err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

func (rs Rows) Save(url, token string, b []byte) error {
	// TODO:
	// The input is the raw json URL for the specific file `quotes.csv` in the gist.
	// Could be handled better, perhaps by unifying the input/output in a struct for
	// the `gist` structure and so sharing the same URL.
	jsonData := []byte(fmt.Sprintf(`{"files":{"quotes.csv":{"content": %s}}}`, strconv.Quote(fmt.Sprintf("%s", b))))

	var v interface{}
	if err := json.Unmarshal(jsonData, &v); err != nil {
		return fmt.Errorf("invalid json: %s", err)
	}

	_, err := PatchResponse(url, token, jsonData)
	if err != nil {
		return err
	}

	return nil
}

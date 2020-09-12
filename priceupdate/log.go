package priceupdate

import (
	"encoding/json"
	"fmt"
)

type LogLevel int

const (
	Critical LogLevel = iota
	Error
	Warning
	Info
	Debug
)

func (l LogLevel) String() string {
	return [...]string{"critical", "error", "warning", "info", "debug"}[l]
}

func (l LogLevel) MarshalJSON() ([]byte, error) {
	return []byte(`"` + l.String() + `"`), nil
}

type Log struct {
	Message  string   `json:"message"`
	Severity LogLevel `json:"severity"`
}

func LogOutput(l Log) {
	log, err := json.Marshal(l)
	if err != nil {
		panic(fmt.Sprintf("marshal failed: %s", err))
	}
	fmt.Println(string(log))
}

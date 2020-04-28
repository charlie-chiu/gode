package gode

import (
	"encoding/json"
	"time"
)

type Jackpot interface {
	Interval() time.Duration
	Fetch() json.RawMessage
}

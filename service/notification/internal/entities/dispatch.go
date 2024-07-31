package entities

import (
	"time"
)

type Dispatch struct {
	ID          string    `json:"id,omitempty"`
	BaseCcy     string    `json:"base_ccy,omitempty"`
	TargetCcies []string  `json:"target_ccies,omitempty"`
	Emails      []string  `json:"emails,omitempty"`
	SendAt      time.Time `json:"send_at,omitempty"`
}

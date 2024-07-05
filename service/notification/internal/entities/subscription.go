package entities

import "time"

type Subscription struct {
	DispatchID  string
	BaseCcy     string
	TargetCcies []string
	Email       string
	SendAt      time.Time
}

func (s Subscription) ToDispatch() *Dispatch {
	return &Dispatch{
		ID:          s.DispatchID,
		BaseCcy:     s.BaseCcy,
		TargetCcies: s.TargetCcies,
		Emails:      []string{s.Email},
		SendAt:      s.SendAt,
	}
}

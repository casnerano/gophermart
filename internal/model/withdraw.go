package model

import (
	"encoding/json"
	"time"
)

type Withdraw struct {
	UUID        string    `json:"-"`
	OrderNumber string    `json:"order"`
	Amount      float64   `json:"sum"`
	UserUUID    string    `json:"-"`
	ProcessedAt time.Time `json:"processed_at"`
}

func (w Withdraw) MarshalJSON() ([]byte, error) {
	type WithdrawAlias Withdraw
	return json.Marshal(&struct {
		WithdrawAlias
		ProcessedAt string `json:"processed_at"`
	}{
		WithdrawAlias: WithdrawAlias(w),
		ProcessedAt:   w.ProcessedAt.Format(time.RFC3339),
	})
}

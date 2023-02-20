package model

import (
	"encoding/json"
	"time"
)

type OrderStatus string

const (
	OrderStatusNew        OrderStatus = "NEW"
	OrderStatusProcessing OrderStatus = "PROCESSING"
	OrderStatusInvalid    OrderStatus = "INVALID"
	OrderStatusProcessed  OrderStatus = "PROCESSED"
)

type Order struct {
	UUID       string      `json:"-"`
	Number     string      `json:"number"`
	Status     OrderStatus `json:"status"`
	Accrual    float64     `json:"accrual,omitempty"`
	UserUUID   string      `json:"-"`
	UploadedAt time.Time   `json:"uploaded_at"`
}

func (o Order) MarshalJSON() ([]byte, error) {
	type OrderAlias Order
	return json.Marshal(&struct {
		OrderAlias
		UploadedAt string `json:"uploaded_at"`
	}{
		OrderAlias: OrderAlias(o),
		UploadedAt: o.UploadedAt.Format(time.RFC3339),
	})
}

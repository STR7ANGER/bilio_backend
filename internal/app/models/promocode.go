package models

import "time"

type Promocode struct {
	Code      string    `json:"code"`
	CreatedAt time.Time `json:"created_at"`
	UsedAt    *time.Time `json:"used_at,omitempty"`
}


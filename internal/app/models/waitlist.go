package models

import "time"

type WaitlistEntry struct {
	Email    string    `json:"email"`
	JoinedAt time.Time `json:"joined_at"`
}

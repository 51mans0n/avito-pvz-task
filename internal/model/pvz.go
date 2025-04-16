package model

import "time"

type PVZ struct {
	ID               string    `json:"id"`
	City             string    `json:"city"`
	RegistrationDate time.Time `json:"registrationDate"`
}

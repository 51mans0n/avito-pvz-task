package model

import "time"

type PVZWithReceptions struct {
	PVZ        *PVZResponse        `json:"pvz"`
	Receptions []ReceptionWithProd `json:"receptions"`
}

type PVZResponse struct {
	ID               string    `json:"id"`
	City             string    `json:"city"`
	RegistrationDate time.Time `json:"registrationDate"`
}

type ReceptionWithProd struct {
	Reception *ReceptionResponse `json:"reception"`
	Products  []ProductResponse  `json:"products"`
}

type ReceptionResponse struct {
	ID       string    `json:"id"`
	PVZID    string    `json:"pvzId"`
	DateTime time.Time `json:"dateTime"`
	Status   string    `json:"status"` // in_progress, close
}

type ProductResponse struct {
	ID          string    `json:"id"`
	DateTime    time.Time `json:"dateTime"`
	Type        string    `json:"type"` // электроника, одежда, обувь
	ReceptionID string    `json:"receptionId"`
}

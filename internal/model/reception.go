package model

import "time"

type Reception struct {
	ID       string    `json:"id" db:"id"`
	PVZID    string    `json:"pvzId" db:"pvz_id"`
	DateTime time.Time `json:"dateTime" db:"date_time"`
	Status   string    `json:"status" db:"status"`
}

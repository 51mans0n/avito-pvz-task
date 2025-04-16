package model

import "time"

type Product struct {
	ID          string    `db:"id"`
	ReceptionID string    `db:"reception_id"`
	DateTime    time.Time `db:"date_time"`
	Type        string    `db:"type"`
}

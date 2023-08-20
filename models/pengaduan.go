package models

import "time"

type Pengaduan struct {
	ID         int
	UserID     int
	Tanggal    time.Time
	NIK        string
	IsiLaporan string
	Status     string
	User       User `gorm:"foreignKey:UserID"`
}

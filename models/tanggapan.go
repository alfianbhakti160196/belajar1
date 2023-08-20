package models

import "time"

type Tanggapan struct {
	ID           int
	PengaduanID  int
	Tanggal      time.Time
	IsiTanggapan string
	PetugasID    int
	Pengaduan    Pengaduan `gorm:"foreignKey:PengaduanID"`
}

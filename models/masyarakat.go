package models

type Masarakat struct {
	IDPengaduan int
	NIK         string `gorm:"primaryKey"`
	Nama        string
	Username    string
	Password    string
	Telpon      string
}

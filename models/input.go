package models

type InputUser struct {
	Nama     string `form:"nama"`
	Username string `form:"username"`
	Password string `form:"password"`
}

type InputLogin struct {
	Username string `form:"username"`
	Password string `form:"password"`
}

type InputLaporan struct {
	Laporan string `form:"laporan"`
	UserID  int    `form:"userid"`
}

type InputTanggapan struct {
	PengaduanID  int    `form:"pengaduanid"`
	IsiTanggapan string `form:"isi"`
}

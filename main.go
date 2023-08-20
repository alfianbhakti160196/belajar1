package main

import (
	"blajar1/models"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
	"time"
)

func main() {
	// Koneksi Ke database
	dsn := "root:20122016@tcp(127.0.0.1:3306)/belajar1?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("error database")
	}

	// Auto Migrasi ke database
	err = db.AutoMigrate(&models.User{}, &models.Pengaduan{}, &models.Tanggapan{}, &models.Masarakat{})
	if err != nil {
		panic("error migrasi database")
	}

	// init server
	r := gin.Default()

	// Create Session
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("userSession", store))

	// Load File HTML
	r.LoadHTMLGlob("views/*")

	r.GET("/home", func(c *gin.Context) {
		c.HTML(http.StatusOK, "home.html", gin.H{
			"title": "Halaman Home",
		})
	})

	// Halaman Login
	r.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", gin.H{
			"title": "Halaman Login",
		})
	})

	// Halaman Register
	r.GET("/register", func(c *gin.Context) {
		c.HTML(http.StatusOK, "register.html", gin.H{
			"title": "Halaman Register",
		})
	})

	// Save Data User Register
	r.POST("/register/save", func(c *gin.Context) {
		// initialisasi variabel untuk menampung data dari form register user
		var input models.InputUser
		err := c.ShouldBind(&input)
		if err != nil {
			panic("error input from user")
		}
		//	Save Data to databse
		var newUser models.User
		newUser.Nama = input.Nama
		newUser.Username = input.Username
		newUser.Password = input.Password
		newUser.Roles = "masysarakat"
		err = db.Create(&newUser).Error
		if err != nil {
			panic("error save ")
		}

		c.Redirect(http.StatusMovedPermanently, "/login")

	})

	// Authtentication
	r.POST("/auth", func(c *gin.Context) {
		s := sessions.Default(c)
		var input models.InputLogin
		err := c.ShouldBind(&input)
		if err != nil {
			c.Redirect(http.StatusMovedPermanently, "/login")
		}

		fmt.Println("username:" + input.Username + " password:" + input.Password)

		var user models.User
		err = db.Where("username = ?", input.Username).Find(&user).Error
		if err != nil {
			c.Redirect(http.StatusMovedPermanently, "/login")
			return
		}

		fmt.Println(user.Roles)

		if user.Password != input.Password {
			c.Redirect(http.StatusMovedPermanently, "/login")
			return
		}

		s.Set("userID", user.ID)
		s.Save()
		if user.Roles == "masysarakat" {
			c.Redirect(http.StatusMovedPermanently, "/dashboard/masyarakat")
			return
		}

		if user.Roles == "petugas" {
			c.Redirect(http.StatusMovedPermanently, "/dashboard/petugas")
			return
		}

	})

	// Dashboard Masyarakat
	r.GET("dashboard/masyarakat", func(c *gin.Context) {
		s := sessions.Default(c)
		user := s.Get("userID")
		if user == 0 {
			c.Redirect(http.StatusMovedPermanently, "/login")
			return
		}
		var pengaduan []models.Pengaduan
		db.Where("user_id = ?", user).Find(&pengaduan)

		c.HTML(http.StatusOK, "dashboardMasyarakat.html", gin.H{
			"title":     "Dashboard Masyarakat",
			"userID":    user,
			"pengaduan": pengaduan,
		})
	})

	// Save Laporan
	r.POST("/laporan/save", func(c *gin.Context) {
		var input models.InputLaporan
		err := c.ShouldBind(&input)
		if err != nil {
			c.Redirect(http.StatusMovedPermanently, "/dashboard/masyarakat")
			return
		}

		// Save Laporan to DB
		var pengaduan models.Pengaduan
		pengaduan.NIK = "123"
		pengaduan.UserID = input.UserID
		pengaduan.IsiLaporan = input.Laporan
		pengaduan.Tanggal = time.Now()
		pengaduan.Status = "proses"

		db.Create(&pengaduan)

		c.Redirect(http.StatusMovedPermanently, "/dashboard/masyarakat")
	})

	// Dashboard Petugas
	r.GET("/dashboard/petugas", func(c *gin.Context) {
		var pengaduan []models.Pengaduan
		db.Preload("User").Find(&pengaduan)

		c.HTML(http.StatusOK, "dashboardPetugas.html", gin.H{
			"title":     "Dashboard Petugas",
			"pengaduan": pengaduan,
		})
	})

	// Logout
	r.GET("/logout", func(c *gin.Context) {
		session := sessions.Default(c)
		session.Delete("userSession")
		c.Redirect(http.StatusMovedPermanently, "/login")

	})

	// Save Tanggapan
	r.POST("/tanggapan/save", func(c *gin.Context) {
		s := sessions.Default(c)
		userID := s.Get("userID")
		var input models.InputTanggapan
		err := c.ShouldBind(&input)
		if err != nil {
			c.Redirect(http.StatusMovedPermanently, "/dashboard/petugas")
			return
		}

		var tanggapan models.Tanggapan
		tanggapan.Tanggal = time.Now()
		tanggapan.PetugasID = userID.(int)
		tanggapan.IsiTanggapan = input.IsiTanggapan
		tanggapan.PengaduanID = input.PengaduanID

		// Save Tanggapan to DB
		db.Create(&tanggapan)

		// Update satatus Tanggapan
		var pengaduan models.Pengaduan
		db.Model(&pengaduan).Where("id = ?", input.PengaduanID).Update("status", "selesai")

		c.Redirect(http.StatusMovedPermanently, "/dashboard/petugas")
	})

	r.Run()
}

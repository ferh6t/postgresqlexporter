package postgresqlexporter

import (
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type conn struct {
	db *gorm.DB
}

func InitDatabase(host string, user string, password string, dbname string, port string, sslmode string) *conn {
	dsn := "host=localhost" + host + " user=" + user + " password=" + password + " dbname=" + dbname + " port=" + port + " sslmode=" + sslmode
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	//db = db.Debug()

	db.AutoMigrate(&Sum{})
	db.AutoMigrate(&Summary{})
	db.AutoMigrate(&Gauge{})
	db.AutoMigrate(&Histogram{})

	dbObject := &conn{
		db: db,
	}
	return dbObject
}

func (a *conn) SaveSum(metric Sum) {
	a.db.Create(&metric)
}
func (a *conn) SaveSummary(metric Summary) {
	a.db.Create(&metric)
}
func (a *conn) SaveGauge(metric Gauge) {
	a.db.Create(&metric)
}
func (a *conn) SaveHistogram(metric Histogram) {
	a.db.Create(&metric)
}

type Sum struct {
	gorm.Model

	Signature  string
	Timestamp  time.Time
	Attributes string
	IntVal     int64
	DoubleVal  float64
	IsInt      bool
	Count      uint
}

type Summary struct {
	gorm.Model

	Signature  string
	Timestamp  time.Time
	Attributes string
	Sum        float64
}

type Gauge struct {
	gorm.Model

	Signature  string
	Timestamp  time.Time
	Attributes string
	IntVal     int64
	DoubleVal  float64
	IsInt      bool
}

type Histogram struct {
	gorm.Model

	Signature  string
	Timestamp  time.Time
	Attributes string
	IntVal     float64
	Count      uint64
}

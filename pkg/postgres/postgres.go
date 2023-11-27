package postgres

import (
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewConnPostgres(cfg *viper.Viper) *gorm.DB {
	host := cfg.GetString("PG_HOST")
	user := cfg.GetString("PG_USER")
	password := cfg.GetString("PG_PASSWORD")
	port := cfg.GetString("PG_PORT")
	db_name := cfg.GetString("DB_NAME")

	dsn := "host=" + host + " user=" + user + " password=" + password + " dbname=" + db_name + " port=" + port + " sslmode=disable TimeZone=UTC"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return db
}

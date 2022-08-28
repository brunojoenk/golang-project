package database

import (
	"fmt"

	"github/brunojoenk/golang-test/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Return new Postgresql db instance
func NewPsqlDB(c *config.Config) (*gorm.DB, error) {
	dataSourceName := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s",
		c.Postgres.PostgresqlHost,
		c.Postgres.PostgresqlUser,
		c.Postgres.PostgresqlPassword,
		c.Postgres.PostgresqlDbname,
		c.Postgres.PostgresqlPort)

	fmt.Println("Datasourcename: " + dataSourceName)

	db, err := gorm.Open(postgres.Open(dataSourceName), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	return db, err
}

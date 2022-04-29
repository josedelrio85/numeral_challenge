package numeral

import (
	"fmt"
	"log"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// Database represent the repository model
type Database struct {
	Db *gorm.DB
}

func CreateDbInstance() Database {
	return Database{}
}

// Open function opens a database connection using Database struct parameters
// Set the DB property of the struct
// Return error | nil
func (d *Database) Open() error {
	db, err := gorm.Open("sqlite3", "database.db")
	if err != nil {
		log.Fatalf("error opening sqlite connection. err: %v", err)
		return err
	}

	c := make(chan bool, 1)
	go func(db *gorm.DB) {
		err = db.DB().Ping()
		if err == nil {
			c <- true
		}
		time.Sleep(2 * time.Second)
	}(db)

	select {
	case res := <-c:
		fmt.Println("Database is ready %n", res)
	case <-time.After(20 * time.Second):
		fmt.Println("timeout 20")
	}

	d.Db = db

	return nil
}

func (d *Database) AutoMigrate() error {
	if err := d.Db.AutoMigrate(Payment{}).Error; err != nil {
		return err
	}
	return nil
}

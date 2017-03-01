package models

import (
	log "github.com/Sirupsen/logrus"
	"github.com/hkparker/Wave/helpers"
	"github.com/jinzhu/gorm"
)

func init() {
	if helpers.TestingCmd() && Orm == nil {
		var err error
		Orm, err = gorm.Open("sqlite3", ":memory:")
		if err != nil {
			log.WithFields(log.Fields{
				"at":    "models.init",
				"error": err.Error(),
			}).Fatal("unable to connect to testing database server")
		}
	}
	if Orm != nil {
		createTables()
	}
}

func Setup() {
	createTables()
	createAdmin()
}

func createTables() {
	if !Orm.HasTable(Alert{}) {
		Orm.CreateTable(Alert{})
		log.WithFields(log.Fields{
			"at": "models.createTables",
		}).Info("creating missing alert table")
	}

	if !Orm.HasTable(Collector{}) {
		Orm.CreateTable(Collector{})
		log.WithFields(log.Fields{
			"at": "models.createTables",
		}).Info("creating missing collector table")
	}

	if !Orm.HasTable(Device{}) {
		Orm.CreateTable(Device{})
		log.WithFields(log.Fields{
			"at": "models.createTables",
		}).Info("creating missing device table")
	}

	if !Orm.HasTable(Network{}) {
		Orm.CreateTable(Network{})
		log.WithFields(log.Fields{
			"at": "models.createTables",
		}).Info("creating missing network table")
	}

	if !Orm.HasTable(Session{}) {
		Orm.CreateTable(Session{})
		log.WithFields(log.Fields{
			"at": "models.createTables",
		}).Info("creating missing session table")
	}

	if !Orm.HasTable(TLS{}) {
		Orm.CreateTable(TLS{})
		log.WithFields(log.Fields{
			"at": "models.createTables",
		}).Info("creating missing tls configuration table")
	}

	if !Orm.HasTable(User{}) {
		Orm.CreateTable(User{})
		log.WithFields(log.Fields{
			"at": "models.createTables",
		}).Info("creating missing user table")
	}
}

func createAdmin() {
	var admins []User
	if err := Orm.Where("Admin = ?", true).Find(&admins).Error; err == nil {
		if len(admins) == 0 {
			var user User
			err = Orm.First(&user, "Username = ?", "root").Error
			if err == nil {
				Orm.Unscoped().Delete(&user)
			}
			admin := User{
				Username: "root",
				Admin:    true,
			}
			password := helpers.RandomString()
			err = admin.SetPassword(password)
			if err != nil {
				log.Fatal(err)
			}
			err = admin.Save()
			if err != nil {
				log.Fatal(err)
			}
			log.WithFields(log.Fields{
				"at":       "models.createAdmin",
				"username": "root",
				"password": password,
			}).Info("created_default_admin")
		}
	} else {
		log.Fatal(err)
	}
}

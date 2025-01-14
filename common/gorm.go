package common

import (
	"dating-be/app/domain/models"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var GormDB *gorm.DB

type Subscriber struct {
	FullName       string `json:"FullName"`
	Username       string `json:"Username"`
	Gender         string `json:"Gender"`
	Age            int8   `json:"Age"`
	Email          string `json:"Email"`
	Password       string `json:"Password"`
	PremiumPackage string `json:"PremiumPackage"`
}

type SubscriberList struct {
	Subscribers []Subscriber `json:"Subscribers"`
}

type GormContext struct {
	Driver   string
	Port     string
	Host     string
	Username string
	Password string
	DBName   string
}

type Gorm interface {
	Open() (*gorm.DB, error)
	Migrate() error
	Seeder() error
}

const (
	POSGRES_CONFIG    = "host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta"
	MYSQL_CONFIG      = "%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local"
	MYSQL_DRIVER      = "mysql"
	POSTGRESQL_DRIVER = "postgres"
)

func NewGormDB(model GormContext) Gorm {
	return GormContext{
		Driver:   model.Driver,
		Port:     model.Port,
		Host:     model.Host,
		Username: model.Username,
		Password: model.Password,
		DBName:   model.DBName,
	}
}

func (g GormContext) Open() (*gorm.DB, error) {

	db, err := g.openDB()
	if err != nil {
		return nil, err
	}

	GormDB = db

	return db, nil
}

func (g GormContext) Migrate() error {

	if err := GormDB.AutoMigrate(
		&models.Subscriber{},
		&models.UserView{}); err != nil {
		return err
	}

	return nil
}

func (g GormContext) Seeder() error {

	var data models.Subscriber
	if err := GormDB.First(&data).Error; err != nil {
		var datas []models.Subscriber
		if fmt.Sprintf("%v", err) == "record not found" {
			seederFiles := []string{
				"seeder/male-non-premium.json",
				"seeder/male-swipe-quota.json",
				"seeder/male-verified-label.json",
				"seeder/female-swipe-quota.json",
				"seeder/female-non-premium.json",
				"seeder/female-verified-label.json",
			}
			for _, seederFile := range seederFiles {
				subscriberList, err := fileToJson(seederFile)
				if err != nil {
					return err
				}

				for _, subscriber := range subscriberList.Subscribers {

					salt, err := GenerateSalt()
					if err != nil {
						return err
					}

					password, err := CreatePassword(subscriber.Password, salt)
					if err != nil {
						return err
					}
					data := models.Subscriber{
						Model:          gorm.Model{},
						FullName:       subscriber.FullName,
						Username:       subscriber.Username,
						Gender:         models.Gender(subscriber.Gender),
						Age:            subscriber.Age,
						Email:          subscriber.Email,
						Password:       password,
						PremiumPackage: models.PremiumPackage(subscriber.PremiumPackage),
						Salt:           base64.StdEncoding.EncodeToString(salt),
					}
					datas = append(datas, data)
				}
			}
			if err := GormDB.Create(&datas).Error; err != nil {
				return err
			}
		} else {
			return err
		}
	}

	return nil
}

func fileToJson(filePath string) (*SubscriberList, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fileContent, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var subscriberList SubscriberList
	if err := json.Unmarshal(fileContent, &subscriberList); err != nil {
		return nil, err
	}

	return &subscriberList, nil
}

func (g GormContext) openDB() (*gorm.DB, error) {

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Millisecond,
			LogLevel:                  logger.Silent,
			IgnoreRecordNotFoundError: true,
			ParameterizedQueries:      true,
			Colorful:                  true,
		},
	)

	config := gorm.Config{
		Logger: newLogger,
	}

	switch strings.ToLower(g.Driver) {
	case MYSQL_DRIVER:
		connectionUrl := fmt.Sprintf(
			MYSQL_CONFIG, g.Username, g.Password, g.Host, g.Port, g.DBName,
		)
		db, err := gorm.Open(mysql.Open(connectionUrl), &config)

		if err != nil {
			return nil, err
		}

		return db, nil
	case POSTGRESQL_DRIVER:
		connectionUrl := fmt.Sprintf(
			POSGRES_CONFIG, g.Host, g.Username, g.Password, g.DBName, g.Port,
		)
		db, err := gorm.Open(postgres.Open(connectionUrl), &config)
		if err != nil {
			return nil, err
		}

		return db, nil
	default:
		newError := errors.New("invalid db driver")
		return nil, newError
	}
}

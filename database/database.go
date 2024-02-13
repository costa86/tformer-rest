package database

import (
	"net/http"

	"github.com/costa86/tformer-rest/helper"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const dbFile = "test.db"

func GetDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dbFile), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

type WhoamiModel struct {
	gorm.Model
	Username string
}

type ProvisionModel struct {
	gorm.Model
	Username     string
	Name         string
	Workspace    string
	Organization string
	Message      string
}

func WhoamiCreate(username string) error {
	db, err := GetDB()

	if err != nil {
		return err
	}
	db.AutoMigrate(&WhoamiModel{})
	db.Create(&WhoamiModel{Username: username})
	return nil
}

func ProvisionCreate(username, name, workspace, organization, message string) error {
	db, err := GetDB()

	if err != nil {
		return err
	}
	db.AutoMigrate(&ProvisionModel{})
	db.Create(&ProvisionModel{
		Username:     username,
		Name:         name,
		Workspace:    workspace,
		Organization: organization,
		Message:      message,
	})
	return nil
}

func whoamiGet() ([]WhoamiModel, error) {
	db, err := GetDB()

	if err != nil {
		return nil, err
	}
	var records []WhoamiModel

	if err := db.Find(&records).Error; err != nil {
		return nil, err
	}
	return records, nil
}

func WhoamiGet(c *gin.Context) {
	_, err := helper.GetClient(helper.GetToken(c))

	if helper.IssueWasFound(c, "", http.StatusBadRequest, err) {
		return
	}

	records, err := whoamiGet()

	if helper.IssueWasFound(c, "", http.StatusBadRequest, err) {
		return
	}

	c.IndentedJSON(http.StatusOK, records)
}

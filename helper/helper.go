package helper

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hashicorp/go-tfe"
)

type Organization struct {
	Name string `form:"org"`
}

type OrganizationWorkspace struct {
	OrgName string `form:"org"`
	WsName  string `form:"ws"`
}

type Variable struct {
	Key          string `json:"key" binding:"required"`
	Value        string `json:"value" binding:"required"`
	Description  string `json:"description" binding:"required"`
	Category     string `json:"category" binding:"required"`
	HCL          bool   `json:"hcl" binding:"required"`
	Sensitive    bool   `json:"sensitive"`
	Workspace    string `json:"workspace" binding:"required"`
	Organization string `json:"organization" binding:"required"`
	Id           string `json:"id"`
}

type ConfigVersion struct {
	Id     string `json:"id"`
	Status string `json:"status"`
}

func IssueWasFound(c *gin.Context, message string, statusCode int, err error) bool {
	msg := message

	if err != nil {
		if message == "" {
			msg = err.Error()
		}
		c.JSON(statusCode, gin.H{"message": msg})
		return true
	}
	return false

}
func GetClient(token string) (*tfe.Client, error) {

	config := &tfe.Config{
		Token:             token,
		RetryServerErrors: true,
		Address:           "https://app.terraform.io",
	}

	client, err := tfe.NewClient(config)
	return client, err
}

func GetToken(c *gin.Context) string {
	return strings.Split(c.GetHeader("Authorization"), "Bearer ")[1]
}

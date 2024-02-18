package helper

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

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

const Address = "https://app.terraform.io" //cloud

// const Address = "https://tfe.d.bbg" //enterprise

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
		Address:           Address,
	}

	client, err := tfe.NewClient(config)
	return client, err
}

func GetToken(c *gin.Context) string {
	return strings.Split(c.GetHeader("Authorization"), "Bearer ")[1]
}
func generateRandomString(length int) string {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	_, err := rand.Read(b)

	if err != nil {
		fmt.Println("Error reading random bytes:", err)
		return ""
	}

	result := ""

	for _, num := range b {
		idx := big.NewInt(0).SetInt64(int64(num) % int64(len(charset)))
		result += string(charset[idx.Int64()])
	}

	return result
}

func createRun(
	client tfe.Client,
	ws tfe.Workspace,
	virtualFile,
	message,
	token,
	org,
	address string,
	autoApply,
	destroy bool,
) (string, error) {
	randomName := generateRandomString(10)
	configFileName := fmt.Sprintf("%s.tf", randomName)

	os.Mkdir(randomName, os.ModePerm)
	configFilePath := filepath.Join(randomName, configFileName)
	err := os.WriteFile(configFilePath, []byte(virtualFile), 0644)

	if err != nil {
		return "issue writing terraform file", err
	}

	ctx := context.Background()

	cv, err := client.ConfigurationVersions.Create(ctx, ws.ID, tfe.ConfigurationVersionCreateOptions{
		AutoQueueRuns: tfe.Bool(false),
		Speculative:   tfe.Bool(false),
	})

	if err != nil {
		return "issue creating configuration version", err
	}

	err = client.ConfigurationVersions.Upload(ctx, cv.UploadURL, randomName)

	if err != nil {
		return "issue uploading configuration version", err
	}

	result, err := client.Runs.Create(ctx, tfe.RunCreateOptions{
		AutoApply:            tfe.Bool(autoApply),
		Workspace:            &ws,
		ConfigurationVersion: cv,
		IsDestroy:            tfe.Bool(destroy),
		Message:              tfe.String(message)})

	if err != nil {
		return "issue creating run", err
	}

	os.RemoveAll(randomName)
	return result.ID, err
}

func GetCurrentTimestamp() string {
	return time.Now().String()[:16]
}

type BasicConfig struct {
	Organization,
	Ws,
	Message string
	AutoApply,
	Destroy bool
}

func getBool(queryValue string) bool {
	var result bool

	value, err := strconv.ParseBool(queryValue)
	if err != nil {
		result = false
	} else {
		result = value
	}
	return result
}

func getBasicConfig(c *gin.Context) BasicConfig {
	return BasicConfig{
		Organization: c.Query("org"),
		Ws:           c.Query("ws"),
		Message:      c.Query("message"),
		AutoApply:    getBool(c.Query("aa")),
		Destroy:      getBool(c.Query("destroy")),
	}
}

func Provision(c *gin.Context, terraformFile string) {
	ctx := context.Background()
	basicConfig := getBasicConfig(c)
	org := basicConfig.Organization
	wsName := basicConfig.Ws
	message := basicConfig.Message
	token := GetToken(c)
	client, err := GetClient(token)

	if IssueWasFound(c, "", http.StatusBadRequest, err) {
		return
	}

	ws, err := client.Workspaces.Read(ctx, org, wsName)

	if IssueWasFound(c, "", http.StatusBadRequest, err) {
		return
	}

	res, err := createRun(
		*client,
		*ws,
		terraformFile,
		message,
		token,
		org,
		Address,
		basicConfig.AutoApply,
		basicConfig.Destroy,
	)

	if IssueWasFound(c, "", http.StatusBadRequest, err) {
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": res,
	})
}

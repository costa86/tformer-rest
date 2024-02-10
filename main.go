package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hashicorp/go-tfe"
)

const port = 3000

func getClient(token string) (*tfe.Client, error) {

	config := &tfe.Config{
		Token:             token,
		RetryServerErrors: true,
		Address:           "https://app.terraform.io",
	}

	client, err := tfe.NewClient(config)
	return client, err
}

type Organization struct {
	Name string `form:"org"`
}

func GetToken(c *gin.Context) string {
	return strings.Split(c.GetHeader("Authorization"), "Bearer ")[1]
}

func ListWorkspaces(client tfe.Client, orgName string) (*tfe.WorkspaceList, error) {
	ctx := context.Background()
	result, err := client.Workspaces.List(ctx, orgName, nil)
	return result, err
}

type Ws struct {
	Name        string    `json:"name"`
	Id          string    `json:"id"`
	CreatedAt   time.Time `json:"createdAt"`
	Description string    `json:"description"`
	Locked      bool      `json:"locked"`
}

func GetOrg(client tfe.Client, orgName string) (*tfe.Organization, error) {
	ctx := context.Background()
	result, err := client.Organizations.Read(ctx, orgName)
	return result, err
}

func issueWasFound(c *gin.Context, message string, statusCode int, err error) bool {
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

func GetWorkspaces(c *gin.Context) {
	var orgStruct Organization
	c.Bind(&orgStruct)
	var wsList []Ws
	var err error

	client, err := getClient(GetToken(c))

	if issueWasFound(c, "", http.StatusBadRequest, err) {
		return
	}

	org, err := GetOrg(*client, orgStruct.Name)

	if issueWasFound(c, "organization not found", http.StatusNotFound, err) {
		return
	}

	workspaces, err := ListWorkspaces(*client, org.Name)

	if issueWasFound(c, "", http.StatusBadRequest, err) {
		return
	}
	for _, v := range workspaces.Items {
		wsList = append(wsList, Ws{v.Name, v.ID, v.CreatedAt, v.Description, v.Locked})
	}

	c.JSON(http.StatusOK, wsList)
}

func main() {
	r := gin.Default()
	r.GET("/workspaces", GetWorkspaces)
	r.Run(fmt.Sprintf(":%d", port))
}

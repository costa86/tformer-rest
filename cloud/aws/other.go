package aws

import (
	"context"
	"fmt"
	"net/http"

	"github.com/costa86/tformer-rest/database"
	"github.com/costa86/tformer-rest/helper"
	"github.com/gin-gonic/gin"
)

type Other struct {
	Name  string `json:"name" binding:"required"`
	Count uint   `json:"count" binding:"required"`
}

func ProvisionOther(c *gin.Context) {
	var resource Other
	ctx := context.Background()
	err := c.BindJSON(&resource)
	org := c.Query("org")
	wsName := c.Query("ws")
	message := c.Query("message")

	if helper.IssueWasFound(c, "", http.StatusBadRequest, err) {
		return
	}

	token := helper.GetToken(c)
	client, err := helper.GetClient(token)

	if helper.IssueWasFound(c, "", http.StatusBadRequest, err) {
		return
	}

	ws, err := client.Workspaces.Read(ctx, org, wsName)

	if helper.IssueWasFound(c, "", http.StatusBadRequest, err) {
		return
	}

	if resource.Count < 3 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "count too little",
		})
		return
	}

	terraformFile := fmt.Sprintf(`variable "name_one" {
		type    = string
		default = "%s"
	}
	  
	variable "name_count_one" {
		type    = number
		default = %d
	}
	
	resource "random_pet" "name_one" {
		prefix = var.name_one
		length = var.name_count_one
	}
	
	output "name_one" {
		value = random_pet.name_one
		}`, resource.Name, resource.Count)

	res, err := helper.ProvisionTerraform(terraformFile, ws.Name, message, token, org, helper.Address)

	if helper.IssueWasFound(c, "", http.StatusBadRequest, err) {
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": string(res),
	})

	user, err := client.Users.ReadCurrent(ctx)

	if helper.IssueWasFound(c, "", http.StatusBadRequest, err) {
		return
	}

	database.ProvisionCreate(user.Email, resource.Name, ws.Name, org, message)
}

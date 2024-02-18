package aws

import (
	"fmt"
	"net/http"

	"github.com/costa86/tformer-rest/helper"
	"github.com/gin-gonic/gin"
)

type Other struct {
	Name  string `json:"name" binding:"required"`
	Count uint   `json:"count" binding:"required"`
}

func ProvisionOther(c *gin.Context) {
	var resource Other

	err := c.BindJSON(&resource)

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

	helper.Provision(c, terraformFile)
}

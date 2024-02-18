package aws

import (
	"fmt"
	"net/http"

	"github.com/costa86/tformer-rest/helper"
	"github.com/gin-gonic/gin"
)

type Bucket struct {
	Name        string `json:"name" binding:"required"`
	Location    string `json:"location" binding:"required"`
	Environment string `json:"environment" binding:"required"`
}

func ProvisionBucket(c *gin.Context) {
	var resource Bucket

	err := c.BindJSON(&resource)

	if helper.IssueWasFound(c, "", http.StatusBadRequest, err) {
		return
	}

	terraformFile := fmt.Sprintf(`provider "google" {
		project = "12345"
		region  = "us-central1"
		access_token = "wfwfwfkmwfmwe"
	  }
	  
	  resource "google_storage_bucket" "example_bucket" {
		name     = "%s"
		location = "%s"
	  
		labels = {
		  environment = "%s"
		}
	  }`, resource.Name, resource.Location, resource.Environment)

	helper.Provision(c, terraformFile)
}

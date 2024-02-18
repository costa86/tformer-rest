package stateversion

import (
	"context"
	"net/http"
	"time"

	"github.com/costa86/tformer-rest/helper"
	"github.com/costa86/tformer-rest/organization"
	"github.com/hashicorp/go-tfe"

	"github.com/gin-gonic/gin"
)

type Run struct {
	ID      string `json:"id"`
	Message string `json:"message"`
}

type StateVersion struct {
	ID               string    `json:"id"`
	CreatedAt        time.Time `json:"createdAt"`
	Status           string    `json:"status"`
	DownloadURL      string    `json:"downloadUrl"`
	TerraformVersion string    `json:"terraformVersion"`
	Run              `json:"run"`
}

func listStateVersions(client tfe.Client, orgName, wsName string) (*tfe.StateVersionList, error) {
	ctx := context.Background()
	result, err := client.StateVersions.List(ctx, &tfe.StateVersionListOptions{Workspace: wsName, Organization: orgName})
	return result, err
}

func GetAll(c *gin.Context) {
	var orgStruct helper.OrganizationWorkspace
	c.Bind(&orgStruct)
	var stateVersionList []StateVersion
	var err error

	client, err := helper.GetClient(helper.GetToken(c))

	if helper.IssueWasFound(c, "", http.StatusBadRequest, err) {
		return
	}

	org, err := organization.Get(*client, orgStruct.OrgName)

	if helper.IssueWasFound(c, "organization not found", http.StatusNotFound, err) {
		return
	}

	stateVersions, err := listStateVersions(*client, org.Name, orgStruct.WsName)

	if helper.IssueWasFound(c, "", http.StatusBadRequest, err) {
		return
	}

	for _, v := range stateVersions.Items {
		stateVersionList = append(stateVersionList, StateVersion{
			ID:               v.ID,
			CreatedAt:        v.CreatedAt,
			TerraformVersion: v.TerraformVersion,
			Status:           string(v.Status),
			DownloadURL:      v.DownloadURL,
			Run:              Run{ID: v.Run.ID, Message: v.Run.Message},
		})
	}

	c.JSON(http.StatusOK, stateVersionList)
}

func GetById(c *gin.Context) {
	id := c.Param("id")

	ctx := context.Background()

	client, err := helper.GetClient(helper.GetToken(c))

	if helper.IssueWasFound(c, "", http.StatusBadRequest, err) {
		return
	}

	sv, err := client.StateVersions.Read(ctx, id)

	if helper.IssueWasFound(c, "", http.StatusBadRequest, err) {
		return
	}

	res := StateVersion{
		ID:               sv.ID,
		CreatedAt:        sv.CreatedAt,
		Status:           string(sv.Status),
		DownloadURL:      sv.DownloadURL,
		TerraformVersion: sv.TerraformVersion,
		Run:              Run{ID: sv.Run.ID, Message: sv.Run.Message},
	}

	c.IndentedJSON(http.StatusOK, res)
}

func GetLatest(c *gin.Context) {
	id := c.Param("id")

	ctx := context.Background()

	client, err := helper.GetClient(helper.GetToken(c))

	if helper.IssueWasFound(c, "", http.StatusBadRequest, err) {
		return
	}

	workspace, err := client.Workspaces.ReadByID(ctx, id)

	if helper.IssueWasFound(c, "", http.StatusBadRequest, err) {
		return
	}

	sv, err := client.StateVersions.ReadCurrent(ctx, workspace.ID)

	if helper.IssueWasFound(c, "", http.StatusBadRequest, err) {
		return
	}

	res := StateVersion{
		ID:               sv.ID,
		CreatedAt:        sv.CreatedAt,
		Status:           string(sv.Status),
		DownloadURL:      sv.DownloadURL,
		TerraformVersion: sv.TerraformVersion,
		Run:              Run{ID: sv.Run.ID, Message: sv.Run.Message},
	}

	c.IndentedJSON(http.StatusOK, res)
}

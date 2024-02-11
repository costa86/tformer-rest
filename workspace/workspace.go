package workspace

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/costa86/tformer-rest/helper"
	"github.com/costa86/tformer-rest/organization"
	"github.com/hashicorp/go-tfe"

	"github.com/gin-gonic/gin"
)

type Ws struct {
	Name        string    `json:"name"`
	Id          string    `json:"id"`
	CreatedAt   time.Time `json:"createdAt"`
	Description string    `json:"description"`
	Locked      bool      `json:"locked"`
}
type Variable struct {
	Key   string    `json:"key"`
	Id    string    `json:"id"`
	Value time.Time `json:"value"`
}

type WsLock struct {
	FieldA string `form:"field_a"`
}

type WsCreation struct {
	Name    string `json:"name"`
	OrgName string `json:"orgName"`
}

func listWorkspaces(client tfe.Client, orgName string) (*tfe.WorkspaceList, error) {
	ctx := context.Background()
	result, err := client.Workspaces.List(ctx, orgName, nil)
	return result, err
}

func GetAll(c *gin.Context) {
	var orgStruct helper.Organization
	c.Bind(&orgStruct)
	var wsList []Ws
	var err error

	client, err := helper.GetClient(helper.GetToken(c))

	if helper.IssueWasFound(c, "", http.StatusBadRequest, err) {
		return
	}

	org, err := organization.Get(*client, orgStruct.Name)

	if helper.IssueWasFound(c, "organization not found", http.StatusNotFound, err) {
		return
	}

	workspaces, err := listWorkspaces(*client, org.Name)

	if helper.IssueWasFound(c, "", http.StatusBadRequest, err) {
		return
	}
	for _, v := range workspaces.Items {
		wsList = append(wsList, Ws{v.Name, v.ID, v.CreatedAt, v.Description, v.Locked})
	}

	c.JSON(http.StatusOK, wsList)
}

func GetById(c *gin.Context) {
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

	ws := Ws{
		Name:        workspace.Name,
		Id:          workspace.ID,
		CreatedAt:   workspace.CreatedAt,
		Description: workspace.Description,
		Locked:      workspace.Locked,
	}

	c.IndentedJSON(http.StatusOK, ws)
}

func create(client tfe.Client, name, organization string) (tfe.Workspace, error) {
	ctx := context.Background()
	result, err := client.Workspaces.Create(ctx, organization, tfe.WorkspaceCreateOptions{Name: tfe.String(name)})
	return *result, err
}

func getByName(client tfe.Client, org, name string) (*tfe.Workspace, error) {
	ws, err := client.Workspaces.Read(context.Background(), org, name)
	return ws, err
}

func Create(c *gin.Context) {
	var newWorkspace WsCreation
	err := c.BindJSON(&newWorkspace)

	if helper.IssueWasFound(c, "", http.StatusBadRequest, err) {
		return
	}

	client, err := helper.GetClient(helper.GetToken(c))

	if helper.IssueWasFound(c, "", http.StatusBadRequest, err) {
		return
	}

	existingWs, err := getByName(*client, newWorkspace.OrgName, newWorkspace.Name)

	if helper.IssueWasFound(c, "", http.StatusBadRequest, err) {
		return
	}

	fmt.Println(existingWs)

	if existingWs != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("workspace already exists: %s", existingWs.Name),
		})
		return
	}

	res, err := create(*client, newWorkspace.Name, newWorkspace.OrgName)
	fmt.Println(res)

	if helper.IssueWasFound(c, "", http.StatusBadRequest, err) {
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": fmt.Sprintf("workspace created: %s", res.Name),
	})
}

func DeleteById(c *gin.Context) {
	id := c.Param("id")

	ctx := context.Background()

	client, err := helper.GetClient(helper.GetToken(c))

	if helper.IssueWasFound(c, "", http.StatusBadRequest, err) {
		return
	}

	err = client.Workspaces.DeleteByID(ctx, id)

	if helper.IssueWasFound(c, "", http.StatusBadRequest, err) {
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("workspace deleted: %s", id),
	})
}

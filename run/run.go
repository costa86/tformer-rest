package run

import (
	"context"
	"net/http"
	"time"

	"github.com/costa86/tformer-rest/helper"
	"github.com/hashicorp/go-tfe"

	"github.com/gin-gonic/gin"
)

type Run struct {
	ID        string    `json:"id"`
	Message   string    `json:"message"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
	IsDestroy bool      `json:"isDestroy"`
}

func listRuns(client tfe.Client, wsId string) (*tfe.RunList, error) {
	ctx := context.Background()
	result, err := client.Runs.List(ctx, wsId, &tfe.RunListOptions{})
	return result, err
}

func GetAll(c *gin.Context) {
	id := c.Param("id")
	var runList []Run
	ctx := context.Background()

	client, err := helper.GetClient(helper.GetToken(c))

	if helper.IssueWasFound(c, "", http.StatusBadRequest, err) {
		return
	}

	workspace, err := client.Workspaces.ReadByID(ctx, id)

	if helper.IssueWasFound(c, "", http.StatusBadRequest, err) {
		return
	}
	runs, err := listRuns(*client, workspace.ID)

	if helper.IssueWasFound(c, "", http.StatusBadRequest, err) {
		return
	}

	for _, v := range runs.Items {
		runList = append(runList, Run{
			ID:        v.ID,
			Message:   v.Message,
			Status:    string(v.Status),
			CreatedAt: v.CreatedAt,
			IsDestroy: v.IsDestroy,
		})
	}

	c.JSON(http.StatusOK, runList)
}

func GetById(c *gin.Context) {
	id := c.Param("id")

	ctx := context.Background()

	client, err := helper.GetClient(helper.GetToken(c))

	if helper.IssueWasFound(c, "", http.StatusBadRequest, err) {
		return
	}

	run, err := client.Runs.Read(ctx, id)

	if helper.IssueWasFound(c, "", http.StatusBadRequest, err) {
		return
	}

	res := Run{
		ID:        run.ID,
		Message:   run.Message,
		Status:    string(run.Status),
		IsDestroy: run.IsDestroy,
		CreatedAt: run.CreatedAt,
	}

	c.IndentedJSON(http.StatusOK, res)
}

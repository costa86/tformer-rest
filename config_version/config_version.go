package configversion

import (
	"context"
	"net/http"

	"github.com/costa86/tformer-rest/helper"
	"github.com/gin-gonic/gin"
	"github.com/hashicorp/go-tfe"
)

type ConfigVersion struct {
	WorkspaceId string `form:"ws-id"`
}

func GetAll(c *gin.Context) {
	var cv ConfigVersion
	c.Bind(&cv)
	ctx := context.Background()
	var cvListResult []helper.ConfigVersion

	client, err := helper.GetClient(helper.GetToken(c))

	if helper.IssueWasFound(c, "", http.StatusBadRequest, err) {
		return
	}

	cvList, err := client.ConfigurationVersions.List(ctx, cv.WorkspaceId, &tfe.ConfigurationVersionListOptions{})

	if helper.IssueWasFound(c, "", http.StatusBadRequest, err) {
		return
	}

	for _, v := range cvList.Items {
		cvListResult = append(cvListResult, helper.ConfigVersion{
			Id:     v.ID,
			Status: string(v.Status)})
	}

	c.JSON(http.StatusOK, cvListResult)
}

func GetById(c *gin.Context) {
	id := c.Param("id")
	ctx := context.Background()

	client, err := helper.GetClient(helper.GetToken(c))

	if helper.IssueWasFound(c, "", http.StatusBadRequest, err) {
		return
	}

	cv, err := client.ConfigurationVersions.Read(ctx, id)

	if helper.IssueWasFound(c, "", http.StatusBadRequest, err) {
		return
	}

	c.JSON(http.StatusOK, cv)
}

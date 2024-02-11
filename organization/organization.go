package organization

import (
	"context"
	"net/http"
	"time"

	"github.com/costa86/tformer-rest/helper"
	"github.com/gin-gonic/gin"
	"github.com/hashicorp/go-tfe"
)

type Organization struct {
	Name      string    `form:"name"`
	CreatedAt time.Time `form:"createdAt"`
}

func Get(client tfe.Client, orgName string) (*tfe.Organization, error) {
	ctx := context.Background()
	result, err := client.Organizations.Read(ctx, orgName)
	return result, err
}

func getAll(client tfe.Client) (*tfe.OrganizationList, error) {
	ctx := context.Background()

	result, err := client.Organizations.List(ctx, nil)
	return result, err

}

func GetAll(c *gin.Context) {
	var orgStruct helper.Organization
	c.Bind(&orgStruct)
	var resultList []Organization
	var err error

	client, err := helper.GetClient(helper.GetToken(c))

	if helper.IssueWasFound(c, "", http.StatusBadRequest, err) {
		return
	}

	queryList, err := getAll(*client)

	if helper.IssueWasFound(c, "", http.StatusBadRequest, err) {
		return
	}
	for _, v := range queryList.Items {
		resultList = append(resultList, Organization{v.Name, v.CreatedAt})
	}

	c.JSON(http.StatusOK, resultList)
}
func GetByName(c *gin.Context) {
	name := c.Param("name")

	client, err := helper.GetClient(helper.GetToken(c))

	if helper.IssueWasFound(c, "", http.StatusBadRequest, err) {
		return
	}

	queryResult, err := Get(*client, name)

	if helper.IssueWasFound(c, "", http.StatusBadRequest, err) {
		return
	}

	c.IndentedJSON(http.StatusOK, queryResult)
}

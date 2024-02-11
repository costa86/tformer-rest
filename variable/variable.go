package variable

import (
	"context"
	"net/http"

	"github.com/costa86/tformer-rest/helper"
	"github.com/gin-gonic/gin"
	"github.com/hashicorp/go-tfe"
)

type VariableResponse struct {
	Key         string `json:"key"`
	Value       string `json:"value"`
	Description string `json:"description"`
	Category    string `json:"category"`
	HCL         bool   `json:"hcl"`
	Sensitive   bool   `json:"sensitive"`
	Workspace   string `json:"workspace"`
	Id          string `json:"id"`
}

func create(client tfe.Client, variable helper.Variable) (tfe.Variable, error) {
	ctx := context.Background()

	result, err := client.Variables.Create(ctx, variable.Workspace, tfe.VariableCreateOptions{
		Type:        "vars",
		Key:         tfe.String(variable.Key),
		Value:       tfe.String(variable.Value),
		Description: tfe.String(variable.Description),
		Category:    tfe.Category(tfe.CategoryType(variable.Category)),
		HCL:         tfe.Bool(variable.HCL),
		Sensitive:   tfe.Bool(variable.Sensitive)})

	return *result, err
}

func Create(c *gin.Context) {
	var newVariable helper.Variable

	err := c.BindJSON(&newVariable)
	if helper.IssueWasFound(c, "", http.StatusBadRequest, err) {
		return
	}

	client, err := helper.GetClient(helper.GetToken(c))

	if helper.IssueWasFound(c, "", http.StatusBadRequest, err) {
		return
	}

	res, err := create(*client, newVariable)

	if helper.IssueWasFound(c, "", http.StatusBadRequest, err) {
		return
	}
	c.IndentedJSON(http.StatusCreated, res)
}

func getAll(client tfe.Client, workspaceId string) (*tfe.VariableList, error) {
	ctx := context.Background()
	result, err := client.Variables.List(ctx, workspaceId, nil)
	return result, err
}

func GetAll(c *gin.Context) {
	wsId := c.Param("ws-id")
	var variableList []helper.Variable

	client, err := helper.GetClient(helper.GetToken(c))

	if helper.IssueWasFound(c, "", http.StatusBadRequest, err) {
		return
	}

	variables, err := getAll(*client, wsId)

	if helper.IssueWasFound(c, "", http.StatusBadRequest, err) {
		return
	}

	for _, v := range variables.Items {
		variableList = append(variableList, helper.Variable{
			Key:         v.Key,
			Value:       v.Value,
			Description: v.Description,
			Id:          v.ID,
		})
	}

	c.JSON(http.StatusOK, variableList)
}

type VariableParam struct {
	VarId string `form:"var-id"`
	WsId  string `form:"ws-id"`
}

func GetById(c *gin.Context) {
	var variable VariableParam
	c.Bind(&variable)
	ctx := context.Background()

	client, err := helper.GetClient(helper.GetToken(c))

	if helper.IssueWasFound(c, "", http.StatusBadRequest, err) {
		return
	}

	varFound, err := client.Variables.Read(ctx, variable.WsId, variable.VarId)

	if helper.IssueWasFound(c, "", http.StatusBadRequest, err) {
		return
	}

	c.JSON(http.StatusOK, VariableResponse{
		Id:          varFound.ID,
		Key:         varFound.Key,
		Value:       varFound.Value,
		Description: varFound.Description,
		Category:    string(varFound.Category),
		HCL:         varFound.HCL,
		Sensitive:   varFound.Sensitive,
		Workspace:   varFound.Workspace.Name,
	})
}

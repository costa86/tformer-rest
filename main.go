package main

import (
	"fmt"

	aws "github.com/costa86/tformer-rest/cloud/aws"
	gcp "github.com/costa86/tformer-rest/cloud/gcp"
	"github.com/costa86/tformer-rest/run"

	cv "github.com/costa86/tformer-rest/config_version"
	"github.com/costa86/tformer-rest/database"
	"github.com/costa86/tformer-rest/organization"
	sv "github.com/costa86/tformer-rest/state_version"
	"github.com/costa86/tformer-rest/user"
	"github.com/costa86/tformer-rest/variable"
	"github.com/costa86/tformer-rest/workspace"

	"github.com/gin-gonic/gin"
)

const port = 3000
const cvRoute = "config-versions"
const varRoute = "variables"
const wsRoute = "workspaces"
const orgRoute = "organizations"
const stateVersionRoute = "state-versions"
const runRoute = "runs"

func main() {
	r := gin.Default()

	//aws
	r.POST("/aws/other", aws.ProvisionOther)
	//gcp
	r.POST("/gcp/bucket", gcp.ProvisionBucket)

	//workspace
	r.GET(fmt.Sprintf("/%s", wsRoute), workspace.GetAll)
	r.GET(fmt.Sprintf("/%s/:id", wsRoute), workspace.GetById)
	r.POST(fmt.Sprintf("/%s", wsRoute), workspace.Create)
	r.DELETE(fmt.Sprintf("/%s/:id", wsRoute), workspace.DeleteById)

	//config version
	r.GET(fmt.Sprintf("/%s", cvRoute), cv.GetAll)
	r.GET(fmt.Sprintf("/%s/:id", cvRoute), cv.GetById)

	//variable
	r.GET(fmt.Sprintf("/%s/:ws-id", varRoute), variable.GetAll)
	r.GET(fmt.Sprintf("/%s", varRoute), variable.GetById)
	r.POST(fmt.Sprintf("/%s", varRoute), variable.Create)

	//organization
	r.GET(fmt.Sprintf("/%s", orgRoute), organization.GetAll)
	r.GET(fmt.Sprintf("/%s/:name", orgRoute), organization.GetByName)

	//user
	r.GET("/whoami", user.Get)

	//state version
	r.GET(fmt.Sprintf("/%s", stateVersionRoute), sv.GetAll)
	r.GET(fmt.Sprintf("/%s/:id", stateVersionRoute), sv.GetById)
	r.GET(fmt.Sprintf("/%s/workspaces/:id", stateVersionRoute), sv.GetLatest)

	//run
	r.GET(fmt.Sprintf("/%s/workspaces/:id", runRoute), run.GetAll)
	r.GET(fmt.Sprintf("/%s/:id", runRoute), run.GetById)

	//database
	r.GET("/records/whoami", database.WhoamiGet)

	r.Run(fmt.Sprintf(":%d", port))
}

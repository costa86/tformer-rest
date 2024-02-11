package main

import (
	"fmt"

	cv "github.com/costa86/tformer-rest/config_version"
	"github.com/costa86/tformer-rest/organization"
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

func main() {
	r := gin.Default()

	//workspace
	r.GET(fmt.Sprintf("/%s", wsRoute), workspace.GetAll)
	r.GET(fmt.Sprintf("/%s/:id", wsRoute), workspace.GetById)
	r.POST(fmt.Sprintf("/%s", wsRoute), workspace.Create)
	r.DELETE(fmt.Sprintf("/%s/:id", wsRoute), workspace.DeleteById)

	//config version
	r.GET(fmt.Sprintf("/%s", cvRoute), cv.GetAll)
	r.GET(fmt.Sprintf("/%s/:id", cvRoute), cv.GetById)

	//run
	// r.POST("/runs", run.Create)

	//variable
	r.GET(fmt.Sprintf("/%s/:ws-id", varRoute), variable.GetAll)
	r.GET(fmt.Sprintf("/%s", varRoute), variable.GetById)
	r.POST(fmt.Sprintf("/%s", varRoute), variable.Create)

	//organization
	r.GET(fmt.Sprintf("/%s", orgRoute), organization.GetAll)
	r.GET(fmt.Sprintf("/%s/:name", orgRoute), organization.GetByName)

	//user
	r.GET("/whoami", user.Get)

	r.Run(fmt.Sprintf(":%d", port))
}

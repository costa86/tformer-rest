package user

import (
	"context"
	"net/http"

	"github.com/costa86/tformer-rest/database"

	"github.com/costa86/tformer-rest/helper"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

func Get(c *gin.Context) {
	ctx := context.Background()
	client, err := helper.GetClient(helper.GetToken(c))

	if helper.IssueWasFound(c, "", http.StatusBadRequest, err) {
		return
	}

	user, err := client.Users.ReadCurrent(ctx)

	if helper.IssueWasFound(c, "", http.StatusBadRequest, err) {
		return
	}

	c.IndentedJSON(http.StatusOK, user)
	database.WhoamiCreate(user.Username)
}

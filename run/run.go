package run

// func Create(c *gin.Context) {
// 	var newVariable helper.Variable
// 	ctx := context.Background()

// 	err := c.BindJSON(&newVariable)
// 	if helper.IssueWasFound(c, "", http.StatusBadRequest, err) {
// 		return
// 	}

// 	client, err := helper.GetClient(helper.GetToken(c))

// 	if helper.IssueWasFound(c, "", http.StatusBadRequest, err) {
// 		return
// 	}

// 	res, err := client.Runs.Create(ctx, tfe.RunCreateOptions{
// 		AutoApply:            tfe.Bool(autoApply),
// 		Workspace:            ws,
// 		ConfigurationVersion: cv,
// 		IsDestroy:            tfe.Bool(isDestroy),
// 		Message:              tfe.String(message)})

// 	if helper.IssueWasFound(c, "", http.StatusBadRequest, err) {
// 		return
// 	}
// 	c.IndentedJSON(http.StatusCreated, res)
// }

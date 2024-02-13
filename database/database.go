package database

import (
	"context"
	"fmt"
	"net/http"

	"github.com/costa86/tformer-rest/helper"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	Username  string
	Timestamp string
}

const mongoUsername = ""
const mongoPassword = ""

func getAllUsers(filter interface{}, collection mongo.Collection) ([]*User, error) {
	var records []*User
	ctx := context.Background()

	cur, err := collection.Find(ctx, filter)
	if err != nil {
		return records, err
	}

	for cur.Next(ctx) {
		var t User
		err := cur.Decode(&t)
		if err != nil {
			return records, err
		}

		records = append(records, &t)
	}

	if err := cur.Err(); err != nil {
		return records, err
	}
	cur.Close(ctx)

	if len(records) == 0 {
		return records, mongo.ErrNoDocuments
	}

	return records, nil
}
func SaveToMongoDb(username, timestamp string) error {
	ctx := context.Background()
	client, err := getMongo(mongoUsername, mongoPassword)

	if err != nil {
		return err
	}

	collection := client.Database("mydb").Collection("mycollection")

	record := User{Username: username, Timestamp: timestamp}

	_, err = collection.InsertOne(ctx, record)

	if err != nil {
		return err
	}

	return nil
}

func WhoamiCreate(username string) error {
	SaveToMongoDb(username, helper.GetCurrentTimestamp())
	return nil
}

func WhoamiGet(c *gin.Context) {
	_, err := helper.GetClient(helper.GetToken(c))

	if helper.IssueWasFound(c, "", http.StatusBadRequest, err) {
		return
	}

	client, err := getMongo(mongoUsername, mongoPassword)

	if err != nil {
		return
	}

	collection := client.Database("mydb").Collection("mycollection")

	filter := bson.D{{}}
	res, err := getAllUsers(filter, *collection)

	if err != nil {
		return
	}

	c.IndentedJSON(http.StatusOK, res)
}

func getMongo(username, password string) (*mongo.Client, error) {
	ctx := context.Background()
	url := fmt.Sprintf("mongodb+srv://%s:%s@main.on0grqm.mongodb.net/?retryWrites=true&w=majority", username, password)
	clientOptions := options.Client().ApplyURI(url)
	client, err := mongo.Connect(ctx, clientOptions)
	return client, err
}

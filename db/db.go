// Made by @reeshuxd
package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var db_url = `mongodb+srv://paid:reeshu1234@cluster1.janyh50.mongodb.net/?retryWrites=true&w=majority` // Your mOngoDb url here....

var (
	ctx, _ = context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	client, _ = mongo.Connect(
		ctx,
		options.Client().ApplyURI(db_url),
	)
	database = client.Database("ReferBot")
)

// Made by @reeshuxd

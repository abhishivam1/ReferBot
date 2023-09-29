// Made by @reeshuxd
package db

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var rdb = database.Collection("Refers")

type Refers struct {
	Refers int64 `json:"refers"`
}

func GetRefers(user_id int64) int64 {
	find := rdb.FindOne(context.TODO(), bson.M{"user_id": user_id})
	if find.Err() == nil {
		var Result Refers
		find.Decode(&Result)
		return Result.Refers
	} else {
		return 0
	}
}

type User struct {
	UserID int64 `bson:"user_id"`
	Refers int64 `bson:"refers"`
}

func GetUsersByRefersAscending() ([]User, error) {
	pipeline := []bson.M{
		{
			"$sort": bson.M{"refers": 1}, // Sort by "refers" in ascending order
		},
		{
			"$group": bson.M{
				"_id":       "$user_id",
				"maxRefers": bson.M{"$max": "$refers"},
			},
		},
	}

	cursor, err := rdb.Aggregate(context.TODO(), pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var userInfo []User

	for cursor.Next(context.TODO()) {
		var result struct {
			UserID    int64 `bson:"_id"`
			MaxRefers int64 `bson:"maxRefers"`
		}
		if err := cursor.Decode(&result); err != nil {
			return nil, err
		}

		userInfo = append(userInfo, User{UserID: result.UserID, Refers: result.MaxRefers})
	}

	return userInfo, nil
}

func Refer_Update(user_id int64, mode string) {
	var point int64
	points := GetRefers(user_id)
	if mode == "e" {
		point = points + 1
	} else if mode == "d" {
		point = points - 1
	}

	// fmt.Println(point)
	rdb.UpdateOne(
		context.TODO(),
		bson.M{"user_id": user_id},
		bson.M{"$set": bson.M{"refers": point}},
		options.Update().SetUpsert(true),
	)
}

func SetCustom(user_id int64, refers int64) {
	go rdb.UpdateOne(
		context.TODO(),
		bson.M{"user_id": user_id},
		bson.M{"$set": bson.M{"refers": refers}},
		options.Update().SetUpsert(true),
	)
}

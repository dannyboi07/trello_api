package db

import (
	"trelloBE/model"
	"trelloBE/schema"
	"trelloBE/util"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func InsertUser(user schema.UserRegister, hashedPw string) error {
	var coll *mongo.Collection = db.Collection(userCollection)
	var dbUser model.User = model.User{
		Email:        *user.Email,
		PasswordHash: hashedPw,
	}

	ctx, cancel := getDbContext()
	defer cancel()
	_, err := coll.InsertOne(ctx, dbUser)

	if err != nil {
		util.Log.Println("Failed to insert a new user into DB, err:", err)
	}

	return err
}

func InsertBoard(userId primitive.ObjectID) (*model.Board, error) {
	var coll *mongo.Collection = db.Collection(boardCollection)
	var dbBoard = model.Board{
		UserId:    userId,
		MemberIds: make([]primitive.ObjectID, 0),
	}

	ctx, cancel := getDbContext()
	defer cancel()

	result, err := coll.InsertOne(ctx, dbBoard)

	if err != nil {
		util.Log.Println("Failed to insert a new board into DB, err:", err)
		return nil, err
	}

	dbBoard.Id = result.InsertedID.(primitive.ObjectID)

	return &dbBoard, nil
}

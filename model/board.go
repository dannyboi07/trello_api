package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Board struct {
	Id        primitive.ObjectID   `bson:"_id,omitempty"`
	UserId    primitive.ObjectID   `bson:"user_id,omitempty"`
	MemberIds []primitive.ObjectID `bson:"member_ids,omitempty"`
}

type Column struct {
	Id      primitive.ObjectID `bson:"_id,omitempty"`
	BoardId primitive.ObjectID `bson:"board_id,omitempty"`
	Title   string             `bson:"title,omitempty"`
}

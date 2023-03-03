package schema

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Board struct {
	Id        primitive.ObjectID   `json:"id"`
	UserId    primitive.ObjectID   `json:"user_id"`
	MemberIds []primitive.ObjectID `json:"member_ids"`
}

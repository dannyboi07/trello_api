package db

import (
	"go.mongodb.org/mongo-driver/mongo"
)

func IsDuplicateKeyError(err error) bool {
	return mongo.IsDuplicateKeyError(err)
}

func IsNotFoundErr(err error) bool {
	return err == mongo.ErrNoDocuments
}

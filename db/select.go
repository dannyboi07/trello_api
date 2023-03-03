package db

import (
	"trelloBE/model"
	"trelloBE/util"
)

// bool return indicates whether user was found or not
func SelectUserByEmail(email string) (model.User, bool, error) {
	ctx, cancel := getDbContext()
	defer cancel()

	var user model.User
	err := db.Collection(userCollection).FindOne(ctx, model.User{
		Email: email,
	}).Decode(&user)

	if err != nil {
		if IsNotFoundErr(err) {
			return user, false, nil
		}

		util.Log.Println("Failed to findOne user from db, err:", err)
		return user, true, err
	}

	return user, true, nil
}

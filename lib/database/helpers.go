package database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateUser(user User) primitive.ObjectID {
	newUser, _ := Users.InsertOne(context.TODO(), user)
	return newUser.InsertedID.(primitive.ObjectID)
}

func UpdateUser(user primitive.M, update bson.M) User {
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	newUser := NewUser()
	Users.FindOneAndUpdate(context.TODO(), user, update, opts).Decode(&newUser)
	return newUser
}

func FindUserViaId(userId primitive.ObjectID) primitive.M {
	user := bson.M{}
	Users.FindOne(context.TODO(), bson.M{"_id": userId}).Decode(&user)
	return user
}

func InsertIntoUserGeneric(user primitive.M, collection string, generic any) User {
	update := bson.M{"$push": bson.M{collection: generic}}
	return UpdateUser(user, update)
}

func InsertSessionTokenIntoUser(user primitive.M, token SessionToken) int {
	userAfter := InsertIntoUserGeneric(user, "sessiontokens", token)
	return len(userAfter.SessionTokens)
}

func RemoveFromUserViaIdGeneric(user primitive.M, collection string, id int) User {
	update := bson.M{"$pull": bson.M{collection: user[collection].(primitive.A)[id]}}
	return UpdateUser(user, update)
}

func RemoveFromUserGeneric(user primitive.M, collection string, generic any) User {
	update := bson.M{"$pull": bson.M{collection: generic}}
	return UpdateUser(user, update)
}

func RemoveSessionTokenFromUser(user primitive.M, token primitive.M) {
	RemoveFromUserGeneric(user, "sessiontokens", token)
}

func RemoveSessionTokenViaIdFromUser(user primitive.M, id int) {
	RemoveFromUserViaIdGeneric(user, "sessiontokens", id)
}

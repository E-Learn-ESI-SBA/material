package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	UserId   string             `json:"userId" bson:"userId"`
	Avatar   string             `json:"avatar" bson:"avatar"`
	Email    string             `json:"email" bson:"email"`
	Group    string             `bson:"group" json:"group"`
	Role     string             `json:"role" bson:"role"`
	Username string             `bson:"username" json:"username"`
	ID       primitive.ObjectID `json:"id" bson:"_id"`
}

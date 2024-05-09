package models

type User struct {
	ID       string `json:"id" bson:"userId"`
	Avatar   string `json:"avatar" bson:"avatar"`
	Email    string `json:"email" bson:"email"`
	Group    string `bson:"group" json:"group"`
	role     string `json:"role" bson:"role"`
	Username string `bson:"username" json:"username"`
}

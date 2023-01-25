package users

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name      string             `bson:"name,omitempty" json:"name"`
	Email     string             `bson:"email,omitempty" json:"email"`
	Password  string             `bson:"password,omitempty" json:"password,omitempty"`
	IsDeleted bool               `bson:"is_deleted,omitempty" json:"is_deleted"`
	LastLogin primitive.DateTime `bson:"last_login,omitempty" json:"last_login"`
	CreatedAt primitive.DateTime `bson:"created_at,omitempty" json:"created_at" default:"now()"`
	UpdatedAt primitive.DateTime `bson:"updated_at,omitempty" json:"updated_at" default:"now()"`
	DeleteAt  primitive.DateTime `bson:"delete_at,omitempty" json:"delete_at"`
}

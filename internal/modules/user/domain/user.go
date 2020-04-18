package domain

import (
	"crypto/sha1"

	p "github.com/wuriyanto48/go-pbkdf2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User model
type User struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	FullName  string             `json:"fullName" bson:"fullName"`
	Email     string             `json:"email" bson:"email"`
	Password  string             `json:"password" bson:"password"`
	Salt      string             `json:"salt" bson:"salt"`
	CreatedAt string             `json:"createdAt" bson:"createdAt"`
	UpdatedAt string             `json:"updatedAt" bson:"updatedAt"`
	Photo     string             `json:"photo" bson:"photo"`
}

// IsValidPassword function
func (u *User) IsValidPassword(password string) bool {
	pass := p.NewPassword(sha1.New, 8, 32, 1000)

	isValid := pass.VerifyPassword(password, u.Password, u.Salt)
	return isValid
}

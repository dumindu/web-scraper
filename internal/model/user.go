package model

import "github.com/google/uuid"

type User struct {
	Model
	Email           string `gorm:"default:null"`
	Auth            *UserAuth
	ActivationToken *UserActivationToken
}

func NewUser(email string, hashedPassword string) *User {
	userID := uuid.New()

	return &User{
		Model: Model{
			ID: userID,
		},
		Email:           email,
		ActivationToken: NewUserActivationToken(userID),
		Auth:            NewUserAuth(userID, hashedPassword),
	}
}

package controller_users

import (
	"errors"
	testTable "serverless-aws-cdk/internal/db/tables"
	"serverless-aws-cdk/utils"
	"time"

	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/google/uuid"
)

type User struct {
	PK        string `json:"pk,omitempty" dynamodbav:"pk,omitempty"`
	ID        string `json:"id,omitempty" dynamodbav:"sk,omitempty"`
	Name      string `json:"name" dynamodbav:"name,omitempty"`
	Email     string `json:"email" dynamodbav:"email,omitempty"`
	Password  string `json:"password" dynamodbav:"password,omitempty"`
	IsActive  int8   `json:"isActive" dynamodbav:"isActive,omitempty"`
	CreatedAt int64  `json:"createdAt" dynamodbav:"createdAt,omitempty"`
	UpdatedAt int64  `json:"updatedAt" dynamodbav:"updatedAt,omitempty"`
}

const PK = "USERS"

func GetUser(id string) (User, error) {
	item, err := testTable.GetItem(PK, id)
	if err != nil {
		return User{}, err
	}

	if item == nil {
		return User{}, nil
	}

	user := User{}
	if err := dynamodbattribute.UnmarshalMap(item, &user); err != nil {
		return User{}, err
	}

	return user, nil
}

func GetAllUsers() ([]User, error) {
	var users []User

	items, err := testTable.QueryItems([]map[string]interface{}{
		{
			"pk": PK,
		},
	})
	if err != nil {
		return nil, err
	}

	if len(items) == 0 {
		return []User{}, nil
	}

	for _, item := range items {
		user := User{}
		if err := dynamodbattribute.UnmarshalMap(item, &user); err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func CreateUser(name, email, password string) error {

	hashedPass, err := utils.HashPassword(password)
	if err != nil {
		return err
	}

	now := time.Now().Unix()

	item := User{
		PK:        PK,
		ID:        uuid.New().String(),
		Name:      name,
		Email:     email,
		Password:  hashedPass,
		IsActive:  1,
		CreatedAt: now,
		UpdatedAt: now,
	}

	return testTable.PutItem(item)
}

func UpdateUser(id, currPass, name, newPass string) error {
	user, err := GetUser(id)
	if err != nil {
		return err
	}

	if user == (User{}) {
		return errors.New("user not found")
	}

	if passMatch := utils.VerifyPassword(currPass, user.Password); !passMatch {
		return errors.New("current password does not match")
	}
	if passMatch := utils.VerifyPassword(newPass, user.Password); !passMatch {
		return errors.New("new password cannot be the same as the current password")
	}

	if name == "" && newPass == "" {
		return errors.New("no fields to update")
	}

	hashedPass, err := utils.HashPassword(newPass)
	if err != nil {
		return err
	}

	item := User{
		Name:      name,
		Password:  hashedPass,
		UpdatedAt: time.Now().Unix(),
	}

	return testTable.UpdateItem(PK, id, utils.StructToMap(item))
}

func DeleteUser(id, password string) error {
	user, err := GetUser(id)
	if err != nil {
		return err
	}

	if user == (User{}) {
		return errors.New("user not found")
	}

	if passMatch := utils.VerifyPassword(password, user.Password); !passMatch {
		return errors.New("current password does not match")
	}

	return testTable.DeleteItem(PK, id)
}

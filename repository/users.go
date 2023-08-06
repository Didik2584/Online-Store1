package repository

import (
	"a21hc3NpZ25tZW50/db"
	"a21hc3NpZ25tZW50/model"
	"encoding/json"
)

type UserRepository struct {
	db db.DB
}

func NewUserRepository(db db.DB) UserRepository {
	return UserRepository{db}
}

func (u *UserRepository) ReadUser() ([]model.Credentials, error) {
	records, err := u.db.Load("users")
	if err != nil {
		return nil, err
	}

	var listUser []model.Credentials
	err = json.Unmarshal(records, &listUser)
	if err != nil {
		return nil, err
	}

	return listUser, nil
}

func (u *UserRepository) AddUser(creds model.Credentials) error {
	var users []model.Credentials

	jsonData, err := u.db.Load("users")
	if err != nil {
		return err
	}

	err = json.Unmarshal(jsonData, &users)
	if err != nil {
		return err
	}

	users = append(users, creds)

	data, err := json.Marshal(users)
	if err != nil {
		return err
	}

	err = u.db.Save("users", data)
	return err
}

func (u *UserRepository) ResetUser() error {
	err := u.db.Reset("users", []byte("[]"))
	if err != nil {
		return err
	}

	return nil
}

func (u *UserRepository) LoginValid(list []model.Credentials, req model.Credentials) bool {
	for _, element := range list {
		if element.Username == req.Username && element.Password == req.Password {
			return true
		}
	}
	return false
}

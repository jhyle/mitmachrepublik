package mmr

import (
	"errors"
	"labix.org/v2/mgo/bson"
)

type (
	UserService struct {
		database  Database
		tablename string
	}
)

func NewUserService(database Database, tablename string) (*UserService, error) {

	err := database.Table(tablename).EnsureIndices("name", "email", "approved", "categories", "addr.city", "addr.pcode")
	if err != nil {
		return nil, err
	} else {
		return &UserService{database, tablename}, nil
	}
}

func (users *UserService) table() Table {

	return users.database.Table(users.tablename)
}

func (users *UserService) Count(place string, categoryIds []int) (int, error) {

	query := buildQuery(place, nil, categoryIds)
	return users.table().Count(query)
}

func (users *UserService) Validate(user *User) error {

	var result []User
	err := users.table().Find(bson.M{"email": user.Email}, &result)
	if err != nil {
		return err
	}

	for i := 0; i < len(result); i++ {
		if result[i].Email == user.Email && result[i].Id != user.Id {
			return errors.New("Email address is already in use.")
		}
	}

	return nil
}

func (users *UserService) Load(id bson.ObjectId) (*User, error) {

	var user User
	err := users.table().LoadById(id, &user)
	return &user, err
}

func (users *UserService) Search(place string, categoryIds []int, page, pageSize int, sort string) (*OrganizerSearchResult, error) {

	var result OrganizerSearchResult
	query := buildQuery(place, nil, categoryIds)
	err := users.table().Search(query, page*pageSize, pageSize, &result, "name")
	return &result, err
}

func (users *UserService) Store(user *User) error {

	_, err := users.table().UpsertById(user.Id, user)
	return err
}

func (users *UserService) Delete(id bson.ObjectId) error {

	return users.table().DeleteById(id)
}

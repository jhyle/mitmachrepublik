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

	err := database.Table(tablename).DropIndices()
	if err != nil {
		return nil, err
	}

	err = database.Table(tablename).EnsureIndices([][]string{
		{"approved", "categories", "addr.city", "addr.pcode"},
		{"name"},
		{"email"},
	})
	if err != nil {
		return nil, err
	}

	return &UserService{database, tablename}, nil
}

func (users *UserService) table() Table {

	return users.database.Table(users.tablename)
}

func (users *UserService) buildQuery(place string, categoryIds []int) bson.M {

	query := make([]bson.M, 0, 3)

	query = append(query, bson.M{"approved": true})

	if len(place) > 0 {
		postcodes := Postcodes(place)
		placesQuery := make([]bson.M, len(postcodes)+1)
		for i, postcode := range postcodes {
			placesQuery[i] = bson.M{"addr.pcode": postcode}
		}
		placesQuery[len(postcodes)] = bson.M{"addr.city": place}
		query = append(query, bson.M{"$or": placesQuery})
	}

	if len(categoryIds) > 0 && categoryIds[0] != 0 {
		categoriesQuery := make([]bson.M, len(categoryIds))
		for i, categoryId := range categoryIds {
			categoriesQuery[i] = bson.M{"categories": categoryId}
		}
		query = append(query, bson.M{"$or": categoriesQuery})
	}

	return bson.M{"$and": query}
}

func (users *UserService) Count(place string, categoryIds []int) (int, error) {

	return users.table().Count(users.buildQuery(place, categoryIds))
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
	err := users.table().Search(users.buildQuery(place, categoryIds), page*pageSize, pageSize, &result, "name")
	return &result, err
}

func (users *UserService) Store(user *User) error {

	_, err := users.table().UpsertById(user.Id, user)
	return err
}

func (users *UserService) Delete(id bson.ObjectId) error {

	return users.table().DeleteById(id)
}

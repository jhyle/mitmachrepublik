package mmr

import (
	"encoding/json"
	"errors"
	"github.com/pilu/traffic"
	"io/ioutil"
	"labix.org/v2/mgo/bson"
)

type (
	Request struct {
		*traffic.Request
	}
)

func (r *Request) ReadBinary() ([]byte, error) {

	return ioutil.ReadAll(r.Request.Body)
}

func (r *Request) ReadJson(v interface{}) error {

	buffer, err := r.ReadBinary()
	if err != nil {
		return err
	}
	return json.Unmarshal(buffer, v)
}

func (r *Request) ReadEmailAndPwd() (*emailAndPwd, error) {

	var form emailAndPwd
	err := r.ReadJson(&form)
	if err != nil {
		return nil, err
	}

	return &form, nil
}

func (r *Request) ReadUser() (*User, error) {

	var data User
	user := &data
	err := r.ReadJson(user)
	return user, err
}

func (r *Request) ReadSessionId() (bson.ObjectId, error) {

	cookie, err := r.Cookie("SESSIONID")
	if err != nil || !bson.IsObjectIdHex(cookie.Value) {
		return "", errors.New("Failed to read session id.")
	}

	return bson.ObjectIdHex(cookie.Value), nil
}

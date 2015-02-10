package model

import (
	"bytes"
	"crypto/sha1"
	"time"
	"golang.org/x/crypto/ssh"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	TableUser = "user"
)

type User struct {
	CreateAt time.Time `bson:"create_at"`

	Name string `bson:"name"`

	Keys     map[string][]byte `bson:"keys"`
	Password []byte            `bson:"password"`

	Permissions ssh.Permissions
}

func (user *User) CheckPassword(password []byte) bool {
	return bytes.Compare(password, user.Password) == 0
}

func (user *User) CheckPublicKey(key []byte) bool {
	for _, userKey := range user.Keys {
		if bytes.Compare(key, userKey) == 0 {
			return true
		}
	}
	return false
}

func HashPassword(password []byte) (hash []byte) {
	tmp := sha1.Sum(password)
	return tmp[:]
}

func GetUser(session *mgo.Session, db, name string) (*User, error) {
	search := bson.M{
		"name": name,
	}
	var user User
	q := session.DB(db).C(TableUser).Find(search)
	if err := q.One(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func PutUser(session *mgo.Session, db string, user *User) error {
	user.CreateAt = time.Now()
	return session.DB(db).C(TableUser).Insert(user)
}
